# gRPC xDS example

xDS 是最初由 Envoy 使用的协议，正在演变为服务网格的通用数据平面 API。

xDS 示例是一个 Hello World 客户端/服务器，能够使用 XDS 管理协议进行配置。开箱即用，它的行为与[我们的其他 hello world 示例](https://github.com/grpc/grpc-go/tree/master/examples/helloworld)相同。服务器回复包含其主机名的响应。

## xDS 环境设置

此示例不包括设置 xDS 环境的说明。请参考您的 xDS 管理服务器的特定文档。示例将稍后添加。

客户端还需要一个引导文件。请参阅 [gRFC A27](https://github.com/grpc/proposal/blob/master/A27-xds-global-load-balancing.md#xdsclient-and-bootstrap-file) 了解引导文件格式。

## 客户端

客户端应用程序需要导入 xDS 包以安装解析器和负载均衡器：

```go
_ "google.golang.org/grpc/xds" // 安装 xds 解析器和负载均衡器。
```

然后，使用 `xds` 目标方案为 ClientConn。

```
$ export GRPC_XDS_BOOTSTRAP=/path/to/bootstrap.json
$ go run client/main.go "xDS world" xds:///target_service
```

