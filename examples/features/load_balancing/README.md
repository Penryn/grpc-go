# Load balancing

这个示例展示了 `ClientConn` 如何选择不同的负载均衡策略。

注意：为了展示负载均衡器的效果，本示例中安装了一个示例解析器来获取后端地址。建议在阅读本示例之前先阅读名称解析器示例。

## 试一试

```
go run server/main.go
```

```
go run client/main.go
```

## 解释

两个回声服务器分别在 ":50051" 和 ":50052" 上提供服务。它们将在响应中包含其服务地址。因此，位于 ":50051" 的服务器将回复 RPC `this is examples/load_balancing (from :50051)`。

创建了两个客户端，连接到这两个服务器（它们从名称解析器获取两个服务器地址）。

每个客户端选择不同的负载均衡器（使用 `grpc.WithDefaultServiceConfig`）：`pick_first` 或 `round_robin`。（这两种策略在 gRPC 中默认支持。要添加自定义的负载均衡策略，请实现 https://godoc.org/google.golang.org/grpc/balancer 中定义的接口。）

注意，负载均衡器也可以通过服务配置进行切换，这允许服务所有者（而不是客户端所有者）选择要使用的负载均衡器。服务配置文档可在 https://github.com/grpc/grpc/blob/master/doc/service_config.md 获取。

### pick_first

第一个客户端配置为使用 `pick_first`。`pick_first` 尝试连接到第一个地址，如果连接成功则使用它进行所有 RPC，如果失败则尝试下一个地址（并继续这样做直到一个连接成功）。因此，所有的 RPC 都将发送到同一个后端。接收到的响应都显示相同的后端地址。

```
this is examples/load_balancing (from :50051)
this is examples/load_balancing (from :50051)
this is examples/load_balancing (from :50051)
this is examples/load_balancing (from :50051)
this is examples/load_balancing (from :50051)
this is examples/load_balancing (from :50051)
this is examples/load_balancing (from :50051)
this is examples/load_balancing (from :50051)
this is examples/load_balancing (from :50051)
this is examples/load_balancing (from :50051)
```

### round_robin

第二个客户端配置为使用 `round_robin`。`round_robin` 连接到它看到的所有地址，并按顺序将每个 RPC 发送到每个后端。例如，第一个 RPC 将发送到后端-1，第二个 RPC 将发送到后端-2，第三个 RPC 将再次发送到后端-1。

```
this is examples/load_balancing (from :50051)
this is examples/load_balancing (from :50051)
this is examples/load_balancing (from :50052)
this is examples/load_balancing (from :50051)
this is examples/load_balancing (from :50052)
this is examples/load_balancing (from :50051)
this is examples/load_balancing (from :50052)
this is examples/load_balancing (from :50051)
this is examples/load_balancing (from :50052)
this is examples/load_balancing (from :50051)
```

注意，可能会看到连续的两个 RPC 发送到同一个后端。这是因为 `round_robin` 只选择准备好进行 RPC 的连接。因此，如果两个连接中的一个由于某种原因未准备好，所有 RPC 将发送到准备好的连接。

