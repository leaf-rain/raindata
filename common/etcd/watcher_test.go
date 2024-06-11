package etcd

import (
	"fmt"
	"github.com/leaf-rain/raindata/common/config"
	"go.uber.org/zap"
	"testing"
	"time"
)

type Tame struct {
	Name string `json:"name,omitempty"`
	Age  int    `json:"age,omitempty"`
	List []int  `json:"list,omitempty"`
}

func marshal(name string) interface{} {
	return &Tame{}
}

func TestNewFileConfig(t *testing.T) {
	watcher, cancel, err := NewFileConfigByCatalogue(ctx, zap.NewExample(), cli, "/test/", marshal)
	defer cancel()
	if err != nil {
		t.Fatal(err)
	}
	infos := []config.ConfigInfo{
		{
			App:  "test",
			Name: "t1",
			Info: &Tame{},
		},
		{
			App:  "test",
			Name: "a2",
			Info: &Tame{},
		},
	}
	watcher.Register(infos...)
	//for {
	//	time.Sleep(time.Second * 3)
	//	fmt.Printf("a1:%v\n", watcher.GetByKey("/test/t1.yaml"))
	//	fmt.Printf("a2:%v\n", watcher.GetByKey("/test/a2.yaml"))
	//}
	watcher.Storage(config.ConfigInfo{
		App:  "test",
		Name: "a3",
		Info: &Tame{
			Name: "a3",
			Age:  3,
		},
	})
	for {
		time.Sleep(time.Second * 3)
		fmt.Printf("a1:%v\n", watcher.GetByKey("test", "t1"))
		fmt.Printf("a2:%v\n", watcher.GetByKey("test", "a2"))
		fmt.Printf("a3:%v\n", watcher.GetByKey("test", "a3"))
	}
}

func TestWatcherConfig_Range(t *testing.T) {
	watcher, cancel, err := NewFileConfigByCatalogue(ctx, zap.NewExample(), cli, "/test/", marshal)
	defer cancel()
	if err != nil {
		t.Fatal(err)
	}
	result := watcher.Range()
	for _, item := range result {
		fmt.Printf("%v\n", item)
	}
}
