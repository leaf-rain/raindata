OS_NAME = $(shell uname -s)
BUILD_NAME = dxl_server_go

# 获取当前执行的Makefile的绝对路径
current_mkfile := $(abspath $(lastword $(MAKEFILE_LIST)))
# 如果需要获取不含文件名的目录路径，可以使用`dir`和`patsubst`函数进一步处理
current_dir_ := $(patsubst %/,%,$(dir $(current_mkfile)))
current_dir := $(patsubst %/,%,$(dir $(current_dir_)))
ORIG_PATH := $(PATH)
PATH := $(current_dir)/github.com/leaf-rain/raindata/common/third_party/bin:$(ORIG_PATH)
export PATH

.PHONY: build_linux
build_linux:
	rm -f $(BUILD_NAME)
	@echo $(BUILD_NAME)编译版本$(VERSION)
	@if [ -n "$(GO)" ] ; then \
			CGO_ENABLED=0 GOOS=linux GOARCH=amd64 GOMAXPROCS=2 GOPROXY=https://goproxy.cn GOPATH=/DevOps/go/gopath GOCACHE=/DevOps/go/cache/go-build GOENV=/DevOps/go/env $(GO) build -buildvcs=false -a -ldflags '-extldflags "-static"' -o $(BUILD_NAME) ./cmd/chat_svr/; \
  		else \
  			CGO_ENABLED=0 GOOS=linux GOARCH=amd64  GOPROXY=https://goproxy.cn go build -a -ldflags '-extldflags "-static"' -o $(BUILD_NAME) ./cmd/;  \
  		fi
	@chmod +x ./$(BUILD_NAME)
	@echo "编译成功 $$(date +"%Y-%m-%d %H:%M:%S")"

.PHONY: swagger
swagger:
	swag init -g ./internal/adapter/http.go -o ./api

.PHONY: generate
generate:
	cd ./cmd && go generate
	cd ./internal/adapter && go generate
	cd ./internal/application && go generate
	cd ./internal/domain && go generate
	cd ./internal/infrastructure && go generate

.PHONY: proto
proto:
	cd api/grpc/proto && protoc --go_out=../ --go_opt=paths=source_relative \
	--go-grpc_out=../ --go-grpc_opt=paths=source_relative \
	./*.proto

.PHONY: help
help:
	@echo help:"\n  --查看帮助文档"
	@echo build:"\n  --编译$(BUILD_NAME)文件"

