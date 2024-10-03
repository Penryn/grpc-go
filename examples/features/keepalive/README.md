# Keepalive

此示例说明如何设置客户端保持活动 ping 和服务器端保持活动 ping 强制和连接空闲设置。有关这些设置的更多详细信息，请参阅[完整文档](https://github.com/grpc/grpc-go/tree/master/Documentation/keepalive.md)。

```
go run server/main.go
```

```
GODEBUG=http2debug=2 go run client/main.go
```

