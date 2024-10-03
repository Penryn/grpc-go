# Interceptor

gRPC 提供了简单的 API 来在每个 ClientConn/Server 基础上实现和安装拦截器。拦截器会拦截每个 RPC 调用的执行。用户可以使用拦截器来进行日志记录、身份验证/授权、指标收集以及许多其他可以跨 RPC 共享的功能。

## 试一试

```
go run server/main.go
```

```
go run client/main.go
```

## 解释

在 gRPC 中，拦截器可以根据它们拦截的 RPC 调用类型分为两类。第一类是 **一元拦截器**，它拦截一元 RPC 调用。另一类是 **流拦截器**，它处理流式 RPC 调用。有关一元 RPC 和流式 RPC 的解释，请参见[这里](https://grpc.io/docs/guides/concepts.html#rpc-life-cycle)。客户端和服务器各自有自己的类型的一元和流拦截器。因此，gRPC 中总共有四种不同类型的拦截器。

### 客户端

#### 一元拦截器

[`UnaryClientInterceptor`](https://godoc.org/google.golang.org/grpc#UnaryClientInterceptor) 是客户端一元拦截器的类型。它本质上是一种函数类型，签名为：`func(ctx context.Context, method string, req, reply interface{}, cc *ClientConn, invoker UnaryInvoker, opts ...CallOption) error`。一元拦截器的实现通常可以分为三个部分：预处理、调用 RPC 方法和后处理。

在预处理中，用户可以通过检查传入的参数获取当前 RPC 调用的信息，例如 RPC 上下文、方法字符串、要发送的请求和配置的 CallOptions。通过这些信息，用户甚至可以修改 RPC 调用。例如，在示例中，我们检查 CallOptions 列表，看看是否配置了调用凭证。如果没有，则配置它以使用令牌 "some-secret-token" 的 oauth2 作为后备。在我们的示例中，我们故意省略配置每个 RPC 凭证以求后备。

预处理完成后，用户可以通过调用 `invoker` 来调用 RPC 调用。

一旦调用者返回回复和错误，用户可以对 RPC 调用进行后处理。通常，这是处理返回的回复和错误。在示例中，我们记录了 RPC 的时间和错误信息。

要在 ClientConn 上安装一元拦截器，请使用 `DialOption` 配置 `Dial`，即 [`WithUnaryInterceptor`](https://godoc.org/google.golang.org/grpc#WithUnaryInterceptor)。

#### 流拦截器

[`StreamClientInterceptor`](https://godoc.org/google.golang.org/grpc#StreamClientInterceptor) 是客户端流拦截器的类型。它是一种函数类型，签名为：`func(ctx context.Context, desc *StreamDesc, cc *ClientConn, method string, streamer Streamer, opts ...CallOption) (ClientStream, error)`。流拦截器的实现通常包括预处理和流操作拦截。

对于预处理，它类似于一元拦截器。

然而，与在之后进行 RPC 方法调用和后处理不同，流拦截器拦截用户对流的操作。首先，拦截器调用传入的 `streamer` 以获取 `ClientStream`，然后包装 `ClientStream` 并通过拦截逻辑重载其方法。最后，拦截器返回包装的 `ClientStream` 供用户操作。

在示例中，我们定义了一个新的结构 `wrappedStream`，它嵌入了一个 `ClientStream`。然后，我们在 `wrappedStream` 上实现（重载）`SendMsg` 和 `RecvMsg` 方法，以拦截这些对嵌入的 `ClientStream` 的操作。在示例中，我们记录了消息类型信息和时间信息以进行拦截。

要为 ClientConn 安装流拦截器，请使用 `DialOption` 配置 `Dial`，即 [`WithStreamInterceptor`](https://godoc.org/google.golang.org/grpc#WithStreamInterceptor)。

### 服务器端

服务器端拦截器类似于客户端拦截器，但提供的信息略有不同。

#### 一元拦截器

[`UnaryServerInterceptor`](https://godoc.org/google.golang.org/grpc#UnaryServerInterceptor) 是服务器端一元拦截器的类型。它是一种函数类型，签名为：`func(ctx context.Context, req interface{}, info *UnaryServerInfo, handler UnaryHandler) (resp interface{}, err error)`。

有关详细的实现解释，请参阅客户端一元拦截器部分。

要为服务器安装一元拦截器，请使用 `ServerOption` 配置 `NewServer`，即 [`UnaryInterceptor`](https://godoc.org/google.golang.org/grpc#UnaryInterceptor)。

#### 流拦截器

[`StreamServerInterceptor`](https://godoc.org/google.golang.org/grpc#StreamServerInterceptor) 是服务器端流拦截器的类型。它是一种函数类型，签名为：`func(srv interface{}, ss ServerStream, info *StreamServerInfo, handler StreamHandler) error`。

有关详细的实现解释，请参阅客户端流拦截器部分。

要为服务器安装流拦截器，请使用 `ServerOption` 配置 `NewServer`，即 [`StreamInterceptor`](https://godoc.org/google.golang.org/grpc#StreamInterceptor)。

