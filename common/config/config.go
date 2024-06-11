package config

type ConfigInfo struct {
	App      string
	Name     string
	Info     interface{}
	FileType string
}

func (conf ConfigInfo) GetKey() string {
	return conf.App + "/" + conf.Name
}

type ConfigInterface interface {
	Register(cs ...ConfigInfo)
	GetByKey(src, key string) interface{}
	Storage(info ConfigInfo) error
	Range() map[string]interface{}
	Close()
}
