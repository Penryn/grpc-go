# Health

gRPC 提供了一个健康检查库，用于向客户端传达系统的健康状况。
它通过 [health/v1](https://github.com/grpc/grpc-proto/blob/master/grpc/health/v1/health.proto) API 提供服务定义。

通过使用健康检查库，客户端可以在服务器遇到问题时优雅地避免使用这些服务器。
大多数语言都提供了开箱即用的实现，使其在系统之间具有互操作性。

## 试用

```
go run server/main.go -port=50051 -sleep=5s
go run server/main.go -port=50052 -sleep=10s
```

```
go run client/main.go
```

## 解释

### 客户端

客户端有两种方式监控服务器的健康状况。
他们可以使用 `Check()` 探测服务器的健康状况，或者使用 `Watch()` 观察变化。

在大多数情况下，客户端不需要直接检查后端服务器。
相反，当在 [服务配置](https://github.com/grpc/proposal/blob/master/A17-client-side-health-checking.md#service-config-changes) 中指定 `healthCheckConfig` 时，他们可以透明地执行此操作。
此配置指示在建立连接时应检查哪个后端 `serviceName`。
空字符串 (`""`) 通常表示应报告服务器的整体健康状况。

```go
// 导入 grpc/health 以启用透明的客户端检查
import _ "google.golang.org/grpc/health"

// 设置适当的服务配置
serviceConfig := grpc.WithDefaultServiceConfig(`{
  "loadBalancingPolicy": "round_robin",
  "healthCheckConfig": {
    "serviceName": ""
  }
}`)

conn, err := grpc.Dial(..., serviceConfig)
```

有关更多详细信息，请参阅 [A17 - 客户端健康检查](https://github.com/grpc/proposal/blob/master/A17-client-side-health-checking.md)。

### 服务器

服务器控制其服务状态。
他们通过检查依赖系统，然后相应地更新自己的状态来实现这一点。
健康服务器可以返回四种状态之一：`UNKNOWN`、`SERVING`、`NOT_SERVING` 和 `SERVICE_UNKNOWN`。

`UNKNOWN` 表示当前状态尚不清楚。
这种状态通常在服务器实例启动时看到。

`SERVING` 表示系统健康并准备好处理请求。
相反，`NOT_SERVING` 表示系统此时无法处理请求。

`SERVICE_UNKNOWN` 表示客户端请求的 `serviceName` 未被服务器识别。
此状态仅由 `Watch()` 调用报告。

服务器可以使用 `healthServer.SetServingStatus("serviceName", servingStatus)` 切换其健康状态。

