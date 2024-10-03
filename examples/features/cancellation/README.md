# Cancellation

此示例展示了客户端如何通过取消传递给 RPC 调用的上下文来取消正在进行的 RPC。客户端将收到状态代码为 `Canceled` 的响应，并且服务处理程序的上下文也将被取消。

```
go run server/main.go
```

```
go run client/main.go
```

