# RBAC authorization

此示例使用 `google.golang.org/grpc/authz` 包中的 `StaticInterceptor`。它使用基于头部的 RBAC 策略将每个 gRPC 方法匹配到所需的角色。为简单起见，context 中注入了包含所需角色的模拟元数据，但这些数据应根据经过身份验证的 context 从适当的服务中获取。

## 试一试

服务器要求经过身份验证的用户具有以下角色才能授权使用这些方法：

- `UnaryEcho` 需要角色 `UNARY_ECHO:W`
- `BidirectionalStreamingEcho` 需要角色 `STREAM_ECHO:RW`

收到请求后，服务器首先检查是否提供了令牌，解码它并检查是否正确设置了密钥（为简单起见，硬编码为 `super-secret`，在生产环境中应使用适当的 ID 提供者）。

如果上述步骤成功，它会使用令牌中的用户名来设置适当的角色（为简单起见，如果用户名匹配 `super-user`，则硬编码为上述两个所需角色，这些角色也应从外部提供）。

启动服务器：

```
go run server/main.go
```

客户端实现展示了如何使用有效令牌（设置用户名和密钥）调用每个端点将成功返回。它还示例了使用错误令牌将导致服务返回 `codes.PermissionDenied`。

启动客户端：

```
go run client/main.go
```

