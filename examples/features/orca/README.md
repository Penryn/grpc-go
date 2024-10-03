# ORCA Load Reporting

ORCA 是一种在服务器和客户端之间报告负载的协议。这个示例展示了如何从客户端和服务器端实现它。更多详情，请参见 [gRFC A51](https://github.com/grpc/proposal/blob/master/A51-custom-backend-metrics.md)。

## 试一试

```
go run server/main.go
```

```
go run client/main.go
```

## 解释

gRPC ORCA 支持提供了两种不同的方式从服务器向客户端报告负载数据：带外和每个 RPC。带外指标会在某个间隔定期通过流报告，而每个 RPC 指标会在调用结束时与尾随元数据一起报告。这两种机制都是可选的，并且独立工作。

完整的 ORCA API 文档可在此处找到：
https://pkg.go.dev/google.golang.org/grpc/orca

### 带外指标

服务器注册一个 ORCA 服务用于带外指标。它通过使用 `orca.Register()` 并在返回的 `orca.Service` 上使用其方法设置指标来实现这一点。

客户端通过 LB 策略接收带外指标。它通过在 `SubConn` 上注册监听器来接收回调，使用 `orca.RegisterOOBListener`。

### 每个 RPC 指标

服务器设置在其 RPC 处理程序中报告查询成本指标。要报告每个 RPC 指标，必须使用 `orca.CallMetricsServerOption()` 选项创建 gRPC 服务器，并通过调用从 `orca.CallMetricRecorderFromContext()` 返回的 `orca.CallMetricRecorder` 的方法来设置指标。

客户端每秒执行一个 RPC。每个 RPC 指标可通过从 LB 策略的选择器返回的 `Done()` 回调获取。

