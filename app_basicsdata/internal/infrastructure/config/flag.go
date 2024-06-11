package config

import (
	"flag"
	"fmt"
	"os"
)

type CmdArgs struct {
	ConfigFile string
	Mode       string
}

func NewCmdArgs() *CmdArgs {
	var result = new(CmdArgs)
	flag.StringVar(&result.Mode, "mode", "pro", "环境")
	flag.StringVar(&result.ConfigFile, "configFile", "./configs/config.yaml", "配置文件路径")
	flag.Parse()
	exePath, _ := os.Getwd()
	fmt.Println("getting executable path:", exePath)
	return result
}
