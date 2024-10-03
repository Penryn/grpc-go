# Description

此示例演示了 gRPC 中基本的 RPC 错误处理。

# 运行示例代码

运行服务器，如果 RPC 请求的 `Name` 字段为空，则返回错误。

```sh
$ go run ./server/main.go
```

然后在另一个终端中运行客户端，客户端会进行两次请求：一次请求的 Name 字段为空，另一次请求的 Name 字段填充为 os/user 提供的当前用户名。

```sh
$ go run ./client/main.go
```

它应该打印从服务器接收到的状态代码。

