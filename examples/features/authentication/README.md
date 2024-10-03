# Authentication
在 gRPC 中，认证被抽象为 [`credentials.PerRPCCredentials`](https://godoc.org/google.golang.org/grpc/credentials#PerRPCCredentials)。它通常也包含授权。用户可以在每个连接或每个调用的基础上进行配置。

当前的认证示例包括一个使用 OAuth2 与 gRPC 的示例。

## 试用

```
go run server/main.go
```

```
go run client/main.go
```

## 解释

### OAuth2

OAuth 2.0 协议是当今广泛使用的认证和授权机制。gRPC 提供了方便的 API 来配置 OAuth 以与 gRPC 一起使用。详情请参考 godoc：https://godoc.org/google.golang.org/grpc/credentials/oauth。

#### 客户端

在客户端，用户应首先获取一个有效的 OAuth 令牌，然后初始化一个实现了 `credentials.PerRPCCredentials` 的 [`oauth.TokenSource`](https://godoc.org/google.golang.org/grpc/credentials/oauth#TokenSource)。接下来，如果用户希望对同一连接上的所有 RPC 调用应用单个 OAuth 令牌，则使用 `DialOption` [`WithPerRPCCredentials`](https://godoc.org/google.golang.org/grpc#WithPerRPCCredentials) 配置 gRPC `Dial`。或者，如果用户希望对每个调用应用 OAuth 令牌，则使用 `CallOption` [`PerRPCCredentials`](https://godoc.org/google.golang.org/grpc#PerRPCCredentials) 配置 gRPC RPC 调用。

请注意，OAuth 需要底层传输是安全的（例如 TLS 等）。

在 gRPC 内部，提供的令牌前缀为令牌类型和一个空格，然后附加到键为 "authorization" 的元数据中。

### 服务器

在服务器端，用户通常在拦截器中获取令牌并验证它。要获取令牌，请在给定的上下文中调用 [`metadata.FromIncomingContext`](https://godoc.org/google.golang.org/grpc/metadata#FromIncomingContext)。它返回元数据映射。接下来，使用键 "authorization" 获取相应的值，这是一个字符串切片。对于 OAuth，该切片应仅包含一个元素，即格式为 `<token-type> + " " + <token>` 的字符串。用户可以通过解析字符串轻松获取令牌，然后验证其有效性。

如果令牌无效，则返回错误代码 `codes.Unauthenticated` 的错误。

如果令牌有效，则调用方法处理程序开始处理 RPC。
