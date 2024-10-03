# Metadata interceptor example

这个示例展示了如何从服务器上的一元和流拦截器更新元数据。
请参阅
[grpc-metadata.md](https://github.com/grpc/grpc-go/blob/master/Documentation/grpc-metadata.md)
获取更多信息。

## 试一试

```
go run server/main.go
```

```
go run client/main.go
```

## 解释

#### 一元拦截器

拦截器可以从传递给它的 RPC 上下文中读取现有的元数据。
由于 Go 上下文是不可变的，拦截器必须创建一个带有更新元数据的新上下文并将其传递给提供的处理程序。

```go
func SomeInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
    // 从 RPC 上下文中获取传入的元数据，并添加一个新的键值对。
    md, ok := metadata.FromIncomingContext(ctx)
    md.Append("key1", "value1")

    // 创建一个带有新元数据的上下文并将其传递给处理程序。
    ctx = metadata.NewIncomingContext(ctx, md)
    return handler(ctx, req)
}
```

#### 流拦截器

`grpc.ServerStream` 不提供修改其 RPC 上下文的方法。因此，流拦截器需要实现 `grpc.ServerStream` 接口并返回一个带有更新元数据的上下文。

最简单的方法是创建一个嵌入 `grpc.ServerStream` 接口的类型，并且只重写 `Context()` 方法以返回一个带有更新元数据的上下文。然后流拦截器将这个包装的流传递给提供的处理程序。

```go
type wrappedStream struct {
    grpc.ServerStream
    ctx context.Context
}

func (s *wrappedStream) Context() context.Context {
    return s.ctx
}

func SomeStreamInterceptor(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
    // 从 RPC 上下文中获取传入的元数据，并添加一个新的键值对。
    md, ok := metadata.FromIncomingContext(ctx)
    md.Append("key1", "value1")

    // 创建一个带有新元数据的上下文并将其传递给处理程序。
    ctx = metadata.NewIncomingContext(ctx, md)

    return handler(srv, &wrappedStream{ss, ctx})
}
```

