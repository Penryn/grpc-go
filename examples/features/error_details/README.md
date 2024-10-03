# Description

此示例演示了在 grpc 错误中使用状态详细信息。

# 运行示例代码

运行服务器：

```sh
$ go run ./server/main.go
```

然后在另一个终端中运行客户端：

```sh
$ go run ./client/main.go
```

它应该成功并打印从服务器收到的问候语。
然后再次运行客户端：

```sh
$ go run ./client/main.go
```

这次，它应该通过打印从服务器收到的错误状态详细信息而失败。

