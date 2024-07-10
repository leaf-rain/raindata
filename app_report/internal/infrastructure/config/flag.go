package config

import "flag"

type CmdArgs struct {
	ConfigFile string
	Mode       string
}

func NewCmdArgs() *CmdArgs {
	var result = new(CmdArgs)
	flag.StringVar(&result.Mode, "mode", "pro", "环境")
	flag.StringVar(&result.ConfigFile, "configFile", "../../configs/config.yaml", "配置文件路径")
	flag.Parse()
	return result
}
