# Encryption

示例中的加密部分包括 TLS 和 ALTS 两种加密机制的独立示例。

## 试用

在每个示例的子目录中：

```
go run server/main.go
```

```
go run client/main.go
```

## 解释

### TLS

TLS 是一种常用的加密协议，用于提供端到端的通信安全。在示例中，我们展示了如何设置服务器认证的 TLS 连接来传输 RPC。

在我们的 `grpc/credentials` 包中，我们提供了几个方便的方法来基于 TLS 创建 grpc
[`credentials.TransportCredentials`](https://godoc.org/google.golang.org/grpc/credentials#TransportCredentials)。
详情请参考 [godoc](https://godoc.org/google.golang.org/grpc/credentials)。

在我们的示例中，我们使用了预先创建的公钥/私钥：
* "server_cert.pem" 包含服务器证书（公钥）。
* "server_key.pem" 包含服务器私钥。
* "ca_cert.pem" 包含可以验证服务器证书的证书（证书颁发机构）。

在服务器端，我们提供 "server.pem" 和 "server.key" 的路径来配置 TLS，并使用
[`credentials.NewServerTLSFromFile`](https://godoc.org/google.golang.org/grpc/credentials#NewServerTLSFromFile)
创建服务器凭证。

在客户端，我们提供 "ca_cert.pem" 的路径来配置 TLS，并使用
[`credentials.NewClientTLSFromFile`](https://godoc.org/google.golang.org/grpc/credentials#NewClientTLSFromFile)
创建客户端凭证。注意，我们将服务器名称覆盖为 "x.test.example.com"，因为服务器证书对 *.test.example.com 有效，但对 localhost 无效。这仅仅是为了方便示例。

一旦在两端创建了凭证，我们可以使用刚创建的服务器凭证启动服务器（通过调用
[`grpc.Creds`](https://godoc.org/google.golang.org/grpc#Creds)），并让客户端使用创建的客户端凭证拨号到服务器（通过调用
[`grpc.WithTransportCredentials`](https://godoc.org/google.golang.org/grpc#WithTransportCredentials)）。

最后，我们通过创建的 `grpc.ClientConn` 进行 RPC 调用，以测试基于 TLS 的安全连接是否成功建立。

### ALTS
注意：ALTS 目前需要在 GCP 上获得特殊的早期访问权限。您可以在 https://groups.google.com/forum/#!forum/grpc-io 中询问详细过程。

ALTS 是 Google 的应用层传输安全协议，支持相互认证和传输加密。注意，ALTS 目前仅在 Google Cloud Platform 上支持，因此您只能在 GCP 环境中成功运行该示例。在我们的示例中，我们展示了如何启动基于 ALTS 的安全连接。

与 TLS 不同，ALTS 使证书/密钥管理对用户透明。因此，设置更为简单。

在服务器端，首先调用
[`alts.DefaultServerOptions`](https://godoc.org/google.golang.org/grpc/credentials/alts#DefaultServerOptions)
获取 alts 的配置，然后将配置提供给
[`alts.NewServerCreds`](https://godoc.org/google.golang.org/grpc/credentials/alts#NewServerCreds)
以基于 alts 创建服务器凭证。

在客户端，首先调用
[`alts.DefaultClientOptions`](https://godoc.org/google.golang.org/grpc/credentials/alts#DefaultClientOptions)
获取 alts 的配置，然后将配置提供给
[`alts.NewClientCreds`](https://godoc.org/google.golang.org/grpc/credentials/alts#NewClientCreds)
以基于 alts 创建客户端凭证。

接下来，与 TLS 相同，使用服务器凭证启动服务器，并让客户端使用客户端凭证拨号到服务器。

最后，进行 RPC 调用，以测试基于 ALTS 的安全连接是否成功建立。

### mTLS

在双向 TLS（mTLS）中，客户端和服务器相互认证。gRPC 允许用户在连接级别配置双向 TLS。

在普通 TLS 中，服务器只需向客户端展示服务器证书以供验证。在双向 TLS 中，服务器还会加载一组受信任的 CA 文件列表，以验证客户端提供的证书。这是通过设置
[`tls.Config.ClientCAs`](https://pkg.go.dev/crypto/tls#Config.ClientCAs)
为受信任的 CA 文件列表，并设置
[`tls.config.ClientAuth`](https://pkg.go.dev/crypto/tls#Config.ClientAuth)
为 [`tls.RequireAndVerifyClientCert`](https://pkg.go.dev/crypto/tls#RequireAndVerifyClientCert) 来完成的。

在普通 TLS 中，客户端只需通过一个或多个受信任的 CA 文件来验证服务器。在双向 TLS 中，客户端还会向服务器展示其客户端证书以供认证。这是通过设置
[`tls.Config.Certificates`](https://pkg.go.dev/crypto/tls#Config.Certificates) 来完成的。

