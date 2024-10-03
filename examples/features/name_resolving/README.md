# Name resolving

此示例展示了 `ClientConn` 如何选择不同的名称解析器。

## 什么是名称解析器

名称解析器可以看作是 `map[service-name][]backend-ip`。它接受一个服务名称，并返回后端 IP 列表。一个常用的名称解析器是 DNS。

在此示例中，创建了一个解析器，将 `resolver.example.grpc.io` 解析为 `localhost:50051`。

## 试一试

```
go run server/main.go
```

```
go run client/main.go
```

## 解释

回声服务器在 ":50051" 上提供服务。创建了两个客户端，一个拨号到 `passthrough:///localhost:50051`，另一个拨号到 `example:///resolver.example.grpc.io`。它们都可以连接到服务器。

名称解析器是根据目标字符串中的 `scheme` 选择的。有关目标语法，请参见 https://github.com/grpc/grpc/blob/master/doc/naming.md。

第一个客户端选择了 `passthrough` 解析器，它接受输入并将其用作后端地址。

第二个客户端连接到服务名称 `resolver.example.grpc.io`。如果没有合适的名称解析器，这将失败。在示例中，它选择了我们安装的 `example` 解析器。`example` 解析器可以通过返回后端地址正确处理 `resolver.example.grpc.io`。因此，即使在创建 ClientConn 时未设置后端 IP，也会创建到正确后端的连接。

