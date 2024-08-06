OS_NAME = $(shell uname -s)
BUILD_NAME = dxl_server_go



.PHONY: help
help:
	@echo help:"\n  --查看帮助文档"
	@echo build:"\n  --编译$(BUILD_NAME)文件"

