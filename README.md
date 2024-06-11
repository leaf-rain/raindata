# raindata

数据上报解析服务

### 1. 技术栈：

前后端通信协议：websocket、protobuf (gorilla)

网关协议：websocket、rpc、protobuf（go-zero、gorilla）

服务协议：grpc、protobuf（go-zero）

数据库: mysql(sqlx)、redis(go-redis)

### 2. 整体架构

todo:

### 2. 目录树

[common 公共库](./github.com/leaf-rain/raindata/common/README.md)：存放公共库代码

[app_basicsdata 元数据服务](./github.com/leaf-rain/raindata/app_basicsdata/README.md)：用户管理&字段属性管理

[app_report 数据上报服务](./github.com/leaf-rain/raindata/app_report/README.md)：数据上报

推荐使用项目目录下的makefile进行脚本编译。

使用：make help


