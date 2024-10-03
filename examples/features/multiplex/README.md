# Multiplex

一个 `grpc.ClientConn` 可以被两个存根共享，两个服务可以共享一个 `grpc.Server`。这个例子展示了如何执行这两种共享。

```
go run server/main.go
```

```
go run client/main.go
```

