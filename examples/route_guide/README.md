# Description
路由指南服务器和客户端演示了如何使用 gRPC Go 库来执行一元、客户端流、服务器流和全双工 RPC。

请参阅 [gRPC 基础：Go](https://grpc.io/docs/tutorials/basic/go.html) 以获取更多信息。

请参阅 `routeguide/route_guide.proto` 中的路由指南服务定义。

# 运行示例代码
要编译和运行服务器，假设您在 `route_guide` 文件夹的根目录下，即 `.../examples/route_guide/`，只需：

```sh
$ go run server/server.go
```

同样，要运行客户端：

```sh
$ go run client/client.go
```

# 可选命令行标志
服务器和客户端都接受可选的命令行标志。例如，客户端和服务器默认情况下不启用 TLS。要启用 TLS：

```sh
$ go run server/server.go -tls=true
```

以及

```sh
$ go run client/client.go -tls=true
```
