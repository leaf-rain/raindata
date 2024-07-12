package config

import (
	"flag"
	"fmt"
	"os"
)

type CmdArgs struct {
	ConfigFile *string
	Mode       *string
}

var Args = &CmdArgs{
	ConfigFile: flag.String("configFile", "./configs/config.yaml", "配置文件路径"),
	Mode:       flag.String("mode", "pro", "环境"),
}

func NewCmdArgs() *CmdArgs {
	fmt.Print("work path:")
	fmt.Println(os.Getwd())
	flag.Parse()
	return Args
}
