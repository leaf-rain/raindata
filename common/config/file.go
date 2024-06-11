package config

import (
	"github.com/fsnotify/fsnotify"
	"github.com/leaf-rain/raindata/common/config/encoding"
	"github.com/leaf-rain/raindata/common/slice"
	"github.com/leaf-rain/raindata/common/str"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"os"
	"path/filepath"
	"reflect"
	"sync"
	"unicode"
)

type FileConfig struct {
	cache           sync.Map
	vipMap          sync.Map
	logger          *zap.Logger
	rootPath        string
	defaultFileType string
}

func NewFileConfig(logger *zap.Logger) *FileConfig {
	return &FileConfig{cache: sync.Map{}, vipMap: sync.Map{}, logger: logger.Named("file_config")}
}

func NewFileConfigByCatalogue(logger *zap.Logger, rootPath, defaultFileType string, marshal func(name string) interface{}) *FileConfig {
	if defaultFileType == "" {
		defaultFileType = "yaml"
	}
	var f = &FileConfig{cache: sync.Map{}, vipMap: sync.Map{}, logger: logger.Named("file_config"), rootPath: rootPath, defaultFileType: defaultFileType}
	var ok bool
	var xt string
	var registerFunc = func(path string, info os.FileInfo, err error) error {
		pathKey := str.RemoveRootPrefix(path, rootPath)
		if pathKey == "" {
			return nil
		}
		// 判断是否配置文件
		xt = filepath.Ext(path)
		if len(xt) > 0 {
			xt = xt[1:]
		} else {
			xt = defaultFileType
		}
		if slice.ContainsElementString(encoding.AllConfigType, xt) == -1 {
			return nil
		}
		if _, ok = f.cache.Load(pathKey); ok {
			return nil
		}
		data := marshal(pathKey)
		f.cache.Store(pathKey, data)
		vip := viper.New()
		vip.SetConfigFile(path)
		vip.SetConfigType(xt)
		f.vipMap.Store(pathKey, vip)
		err = vip.ReadInConfig()
		if err != nil {
			f.logger.Error("[NewFileConfigByCatalogue] config file vip.ReadInConfig failed", zap.String("name", path), zap.Error(err))
			return nil
		}
		err = vip.Unmarshal(data)
		if err != nil {
			f.logger.Error("[NewFileConfigByCatalogue] config file vip.Unmarshal failed", zap.String("name", path), zap.Error(err))
		}
		f.logger.Info("[NewFileConfigByCatalogue] config file vip.Unmarshal success", zap.String("name", path), zap.Error(err))
		go f.dynamicConfig(vip, data)
		return nil
	}
	err := filepath.Walk(rootPath, registerFunc)
	if err != nil {
		f.logger.Error("[NewFileConfigByCatalogue] filepath.Walk failed", zap.String("rootPath", rootPath), zap.Error(err))
	}
	return f
}

func (f *FileConfig) SetDefaultFileType(ty string) error {
	if slice.ContainsElementString(encoding.AllConfigType, ty) == -1 {
		return UnsupportedType
	}
	f.defaultFileType = ty
	return nil
}

func (f *FileConfig) Register(cs ...ConfigInfo) {
	var ok bool
	var err error
	for _, item := range cs {
		var pathKey = str.GetFileName(item.Name)
		if _, ok = f.cache.Load(pathKey); ok {
			continue
		}
		if slice.ContainsElementString(encoding.AllConfigType, item.FileType) == -1 {
			item.FileType = f.defaultFileType
		}
		f.cache.Store(pathKey, item.Info)
		vip := viper.New()
		vip.SetConfigFile(f.rootPath + pathKey + "." + f.defaultFileType)
		vip.SetConfigType(f.defaultFileType)
		f.vipMap.Store(pathKey, vip)
		err = vip.ReadInConfig()
		if err != nil {
			f.logger.Error("[Register] config file vip.ReadInConfig failed", zap.String("name", item.Name), zap.Error(err))
			continue
		}
		err = vip.Unmarshal(item.Info)
		if err != nil {
			f.logger.Error("[Register] config file vip.Unmarshal failed", zap.String("name", item.Name), zap.Error(err))
		}
		err = vip.WriteConfig()
		if err != nil {
			f.logger.Error("[Storage] vip.WriteConfig failed", zap.String("name", item.Name), zap.Error(err))
			continue
		}
		go f.dynamicConfig(vip, item.Info)
	}
}

func (f *FileConfig) dynamicConfig(vc *viper.Viper, cfg interface{}) {
	vc.OnConfigChange(func(event fsnotify.Event) {
		f.logger.Info("[dynamicConfig] config file vip.Unmarshal", zap.String("name", vc.ConfigFileUsed()))
		if err := vc.Unmarshal(cfg); err != nil {
			f.logger.Error("[dynamicConfig] config file vip.Unmarshal failed", zap.String("name", vc.ConfigFileUsed()), zap.Error(err))
		}
	})
	vc.WatchConfig()
}

func (f *FileConfig) GetByKey(src, key string) interface{} {
	key = str.GetFileName(key)
	key = src + "/" + key
	info, exist := f.cache.Load(key)
	if !exist {
		return nil
	}
	return info
}

func (f *FileConfig) Storage(info ConfigInfo) error {
	var ok bool
	var err error
	var pathKey = str.GetFileName(info.Name)
	if slice.ContainsElementString(encoding.AllConfigType, info.FileType) == -1 {
		info.FileType = f.defaultFileType
	}
	if _, ok = f.cache.Load(pathKey); !ok {
		err = f.EnsurePathExists(f.rootPath+info.Name, info.FileType)
		if err != nil {
			f.logger.Error("[Storage] EnsurePathExists failed", zap.String("name", info.Name), zap.Error(err))
			return err
		}
		f.Register(info)
	}
	var vipI interface{}
	vipI, ok = f.vipMap.Load(pathKey)
	if !ok {
		f.logger.Error("[Storage] vipMap.Load failed", zap.String("name", info.Name))
		return err
	}
	var vip *viper.Viper
	vip, ok = vipI.(*viper.Viper)
	if !ok {
		f.logger.Error("[Storage] vipI type failed", zap.String("name", info.Name))
		return err
	}
	setConfigFromStruct(vip, info.Info)
	err = vip.WriteConfig()
	if err != nil {
		f.logger.Error("[Storage] vip.WriteConfig failed", zap.String("name", info.Name), zap.Error(err))
		return err
	}
	return nil
}

// EnsurePathExists 检查给定的路径是否存在，如果不存在则尝试创建它
func (f *FileConfig) EnsurePathExists(path, ty string) error {
	if filepath.Ext(path) == "" {
		path += "." + ty
	}
	// 使用filepath.EvalSymlinks检查路径是否存在，这会处理符号链接
	if _, err := os.Stat(path); os.IsNotExist(err) {
		// 路径不存在，尝试创建
		var file *os.File
		file, err = os.Create(path)
		if err != nil {
			f.logger.Error("[EnsurePathExists] dont creat path", zap.String("path", path), zap.Error(err))
			return err
		}
		defer file.Close()
	} else if err != nil {
		f.logger.Error("[EnsurePathExists] check path failed", zap.String("path", path), zap.Error(err))
		return err
	}
	return nil
}

// setConfigFromStruct 遍历结构体并将字段值设置到Viper中
func setConfigFromStruct(v *viper.Viper, config interface{}) {
	val := reflect.ValueOf(config)
	// 检查是否是指针类型并解指针
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}
	typ := val.Type()
	var name string
	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		name = typ.Field(i).Name
		if name != "" && field.Type().Kind() != reflect.Func && unicode.IsUpper([]rune(name)[0]) {
			v.Set(name, field.Interface())
		}
	}
}
func (f *FileConfig) Range() map[string]interface{} {
	var result = make(map[string]interface{})
	f.cache.Range(func(key, value interface{}) bool {
		result[key.(string)] = value
		return true
	})
	return result
}
func (f *FileConfig) Close() {

}
