package config

import (
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"testing"
)

var f *FileConfig
var f2 *FileConfig

func TestMain(m *testing.M) {
	config := zap.NewProductionConfig()
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder // 自定义时间格式
	logger, _ := config.Build()
	f = NewFileConfig(logger)
	f2 = NewFileConfigByCatalogue(logger, "./test/", "", func(name string) interface{} {
		return &T{}
	})
	m.Run()
}

type T struct {
	Name string  `json:"name"`
	Age  int64   `json:"age"`
	List []int64 `json:"list"`
}

type T2 struct {
	Name string `yaml:"name" mapstructure:"name"`
	Age  int64  `yaml:"age" mapstructure:"age"`
}

func TestFileConfig_Register2(t *testing.T) {
	content := f2.GetByKey("t1")
	fmt.Printf("%+v\n", content)
	f2.Register(ConfigInfo{
		Name: "t2",
		Info: &T{
			Name: "t2",
			Age:  2,
			List: []int64{2, 3, 4, 5},
		},
	})
	content = f2.GetByKey("t2")
	fmt.Printf("%+v\n", content)
	err := f2.Storage(ConfigInfo{
		Name: "t3",
		Info: &T{
			Name: "t3",
			Age:  3,
			List: []int64{3, 4, 5},
		},
	})
	content = f2.GetByKey("t3")
	fmt.Printf("%+v\n", content)
	if err != nil {
		t.Error(err)
	}
}

func TestFileConfig_Range(t *testing.T) {
	result := f2.Range()
	for _, item := range result {
		fmt.Printf("%+v\n", item)
	}
}
