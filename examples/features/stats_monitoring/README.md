# Stats Monitoring Handler

此示例演示了如何使用 [`stats`](https://pkg.go.dev/google.golang.org/grpc/stats) 包来报告各种网络和 RPC 统计信息。  
_请注意，所有字段都是只读的，并且 `stats` 包的 API 是实验性的_。

## 试一试

```
go run server/main.go
```

```
go run client/main.go
```

## 解释

gRPC 提供了一种机制，通过 [`stats.Handler`](https://pkg.go.dev/google.golang.org/grpc/stats#Handler) 接口挂钩到请求-响应网络周期的各个事件（阶段）。要访问这些事件，需要将实现 `stats.Handler` 的具体类型传递给客户端的 `grpc.WithStatsHandler()` 和服务器端的 `grpc.StatsHandler()`。

`stats.Handler` 上的 `HandleRPC(context.Context, RPCStats)` 方法在请求-响应周期内被多次调用，并且各种事件统计信息会传递给其 `RPCStats` 参数（一个接口）。实现此接口的具体类型有：`*stats.Begin`、`*stats.InHeader`、`*stats.InPayload`、`*stats.InTrailer`、`*stats.OutHeader`、`*stats.OutPayload`、`*stats.OutTrailer` 和 `*stats.End`。这些事件的顺序在客户端和服务器端有所不同。

同样，`stats.Handler` 上的 `HandleConn(context.Context, ConnStats)` 方法在连接开始时调用一次（使用 `*stats.ConnBegin`），在连接结束时调用一次（使用 `*stats.ConnEnd`）。

[`stats.Handler`](https://pkg.go.dev/google.golang.org/grpc/stats#Handler) 接口还提供了 `TagRPC(context.Context, *RPCTagInfo) context.Context` 和 `TagConn(context.Context, *ConnTagInfo) context.Context` 方法。这些方法主要用于将网络相关信息附加到给定的上下文中。

`TagRPC(context.Context, *RPCTagInfo) context.Context` 方法返回一个上下文，从该上下文派生出 RPC 的其余生命周期中使用的上下文。此行为在 gRPC 客户端和服务器之间是一致的。

从 `TagConn(context.Context, *ConnTagInfo) context.Context` 返回的上下文具有不同的生命周期：

- 在 gRPC 客户端中：  
  RPC 的其余生命周期中使用的上下文不会从此上下文派生。因此，附加到此上下文的信息只能由 `HandleConn(context.Context, ConnStats)` 方法使用。
- 在 gRPC 服务器中：  
  RPC 的其余生命周期中使用的上下文将从此上下文派生。

注意：[stats](https://pkg.go.dev/google.golang.org/grpc/stats) 包应仅用于网络监控目的，而不是作为 [interceptors](https://github.com/grpc/grpc-go/blob/master/examples/features/interceptor) 的替代品。

