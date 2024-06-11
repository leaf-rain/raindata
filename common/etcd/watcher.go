package etcd

import (
	"bytes"
	"context"
	"github.com/leaf-rain/raindata/common/config"
	"github.com/leaf-rain/raindata/common/config/encoding"
	"github.com/leaf-rain/raindata/common/str"
	"github.com/spf13/viper"
	"go.etcd.io/etcd/api/v3/mvccpb"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.uber.org/zap"
	"path/filepath"
	"reflect"
	"strings"
	"sync"
	"unicode"
)

type WatcherConfig struct {
	cache    sync.Map
	logger   *zap.Logger
	etcdCli  *clientv3.Client
	pathRoot string
	once     sync.Once
	locker   sync.Locker

	ch        clientv3.WatchChan
	closeChan chan struct{}
	ctx       context.Context
	cancel    context.CancelFunc
	encoding  *encoding.Encoding
}

type value struct {
	Info     interface{}
	fileType string
}

var defaultConfigType = ".yaml"

func NewFileConfig(ctx context.Context, logger *zap.Logger, etcdCli *clientv3.Client) (*WatcherConfig, func(), error) {
	if etcdCli == nil {
		return nil, func() {}, nil
	}
	w := &WatcherConfig{
		cache:    sync.Map{},
		logger:   logger.Named("etcd_watcher_config"),
		etcdCli:  etcdCli,
		once:     sync.Once{},
		encoding: encoding.NewEncoding(),
	}
	newCtx, cancelFunc := context.WithCancel(ctx)
	w.ctx = newCtx
	w.cancel = cancelFunc
	w.closeChan = make(chan struct{})
	return w, w.Close, nil
}

func NewFileConfigByCatalogue(ctx context.Context, logger *zap.Logger, etcdCli *clientv3.Client, rootPath string, marshal func(name string) interface{}) (*WatcherConfig, func(), error) {
	if etcdCli == nil {
		return nil, func() {}, nil
	}
	f := &WatcherConfig{
		cache:    sync.Map{},
		logger:   logger.Named("etcd_watcher_config"),
		etcdCli:  etcdCli,
		once:     sync.Once{},
		encoding: encoding.NewEncoding(),
	}
	newCtx, cancelFunc := context.WithCancel(ctx)
	f.ctx = newCtx
	f.cancel = cancelFunc
	f.pathRoot = rootPath
	f.closeChan = make(chan struct{})
	if !strings.HasSuffix(f.pathRoot, "/") {
		f.pathRoot = f.pathRoot + "/"
		rootPath += "/"
	}
	datas, err := f.etcdCli.Get(f.ctx, rootPath, clientv3.WithPrefix())
	if err != nil {
		f.logger.Error("[RegisterCatalogue] cli.Get failed", zap.String("rootPath", rootPath), zap.Error(err))
		return nil, func() {}, nil
	}
	var pathKey string
	var fileType string
	for _, item := range datas.Kvs {
		fileType = filepath.Ext(string(item.Key))
		pathKey = str.RemoveRootPrefix(string(item.Key), rootPath)
		var info = marshal(string(item.Key))
		f.cache.Store(pathKey, value{
			Info:     info,
			fileType: fileType,
		})
		f.storage(item)
	}
	f.ch = f.etcdCli.Watch(f.ctx, f.pathRoot, clientv3.WithPrefix())
	go f.Watcher()
	return f, f.Close, nil
}

func (f *WatcherConfig) Register(cs ...config.ConfigInfo) {
	if f.pathRoot == "" {
		ps := make([]string, len(cs))
		for index := range ps {
			ps[index] = cs[index].Name
		}
		f.pathRoot = longestCommonPath(ps)
		f.ch = f.etcdCli.Watch(f.ctx, f.pathRoot, clientv3.WithPrefix())
		go f.Watcher()
	}
	var pathKey, fileType string
	for _, item := range cs {
		pathKey = item.GetKey()
		fileType = filepath.Ext(item.Name)
		if fileType == "" {
			fileType = defaultConfigType
		}
		f.cache.Store(pathKey, value{
			Info:     item.Info,
			fileType: fileType,
		})
		result, err := f.etcdCli.Get(f.ctx, f.pathRoot+item.Name+fileType)
		if err != nil {
			f.logger.Error("[Register] etcd get error", zap.Error(err))
			continue
		}
		if len(result.Kvs) == 0 {
			err = f.Storage(item)
			if err != nil {
				f.logger.Error("[Register] Storage error", zap.Error(err))
				continue
			}
		} else {
			for _, kv := range result.Kvs {
				f.storage(kv)
			}
		}
	}
}

// removeAfterLast finds the last occurrence of a separator in a string and returns a new string without the part after it.
func removeAfterLast(str, sep string) string {
	// Find the index of the last occurrence of the separator
	lastSepIndex := strings.LastIndex(str, sep)

	// If the separator is not found, return the original string
	if lastSepIndex == -1 {
		return str
	}

	// Return the substring up to the last occurrence of the separator
	return str[:lastSepIndex]
}

func longestCommonPath(paths []string) string {
	if len(paths) == 0 {
		return "/"
	}
	if len(paths) == 1 {
		return removeAfterLast(paths[0], "/")
	}
	base := paths[0]
	for i := 1; i < len(paths); i++ {
		if strings.HasPrefix(paths[i], base) {
			continue
		} else {
			base = removeAfterLast(base, "/")
		}
	}
	if base == "" {
		base = "/"
	}
	return base
}

func (f *WatcherConfig) Watcher() {
	for {
		select {
		case <-f.ctx.Done():
			return
		case <-f.closeChan:
			return
		case kv, ok := <-f.ch:
			if !ok {
				return
			}
			for _, v := range kv.Events {
				f.storage(v.Kv)
			}
		}
	}
}

func (f *WatcherConfig) storage(item *mvccpb.KeyValue) {
	var info interface{}
	var ok bool
	var vip = viper.New()
	var fileType string
	var err error
	k := string(item.Key)
	pathKey := str.RemoveRootPrefix(k, f.pathRoot)
	info, ok = f.cache.Load(pathKey)
	if !ok {
		return
	}
	fileType = info.(value).fileType
	if fileType == "" {
		//f.logger.Error("[storage] fileType is empty", zap.String("name", k))
		//return
		fileType = defaultConfigType
	} else {
		fileType = fileType[1:]
	}
	vip.SetConfigType(fileType)
	err = vip.ReadConfig(bytes.NewReader(item.Value))
	if err != nil {
		f.logger.Error("[storage] ReadConfig failed", zap.String("name", k), zap.Error(err))
		return
	}
	err = vip.Unmarshal(info.(value).Info)
	if err != nil {
		f.logger.Error("[storage] vip.Unmarshal failed", zap.String("name", k), zap.Error(err))
		return
	}
}

func (f *WatcherConfig) Close() {
	close(f.closeChan)
}

func (f *WatcherConfig) GetByKey(src, key string) interface{} {
	key = str.RemoveRootPrefix(key, f.pathRoot)
	key = src + "/" + key
	info, exist := f.cache.Load(key)
	if !exist {
		return nil
	}
	return info.(value).Info
}

func (f *WatcherConfig) Storage(info config.ConfigInfo) error {
	var ok bool
	var err error
	var cache interface{}
	cache, ok = f.cache.Load(info.GetKey())
	if !ok {
		f.Register(info)
		cache, _ = f.cache.Load(info.GetKey())
	}
	var vip = viper.New()
	fileType := cache.(value).fileType
	if fileType == "" {
		//f.logger.Error("[Storage] fileType is empty", zap.String("name", info.GetKey()))
		//return err
		fileType = defaultConfigType
	}
	fileType = fileType[1:]
	vip.SetConfigType(fileType)
	setConfigFromStruct(vip, info.Info)
	var byteSlice []byte
	byteSlice, err = f.encoding.Encode(fileType, vip.AllSettings())
	if err != nil {
		f.logger.Error("[Storage] encoding.Encode failed", zap.String("app", info.App), zap.String("name", info.Name), zap.String("file type", fileType), zap.Error(err))
		return err
	}
	if len(byteSlice) > 0 {
		var fileName = f.pathRoot + "/" + info.App + "/" + info.Name + "." + fileType
		fileName = strings.ReplaceAll(fileName, "//", "/")
		_, err = f.etcdCli.Put(f.ctx, fileName, string(byteSlice))
		if err != nil {
			f.logger.Error("[Storage] etcdCli.Put failed", zap.String("name", info.GetKey()), zap.String("file type", fileType), zap.Error(err))
			return err
		}
	}
	return nil
}

func (f *WatcherConfig) Range() map[string]interface{} {
	var result = make(map[string]interface{})
	f.cache.Range(func(key, v interface{}) bool {
		result[key.(string)] = v.(value).Info
		return true
	})
	return result
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
