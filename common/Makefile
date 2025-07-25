GOHOSTOS:=$(shell go env GOHOSTOS)
GOPATH:=$(shell go env GOPATH)
VERSION=$(shell git describe --tags --always)

ifeq ($(GOHOSTOS), windows)
	#the `find.exe` is different from `find` in bash/shell.
	#to see https://docs.microsoft.com/en-us/windows-server/administration/windows-commands/find.
	#changed to use git-bash.exe to run find cli or other cli friendly, caused of every developer has a Git.
	#Git_Bash= $(subst cmd\,bin\bash.exe,$(dir $(shell where git)))
	Git_Bash=$(subst \,/,$(subst cmd\,bin\bash.exe,$(dir $(shell where git))))
	INTERNAL_PROTO_FILES=$(shell $(Git_Bash) -c "find . -name *.proto")
else
	INTERNAL_PROTO_FILES=$(shell find . -name *.proto)
endif

.PHONY: build
# generate internal proto
build:
	protoc --proto_path=. \
 	       --go_out=paths=source_relative:. \
	       $(INTERNAL_PROTO_FILES)


.PHONY: help
help:
	@echo help:"\n  --查看帮助文档"
	@echo build:"\n  --编译$(BUILD_NAME)文件"

