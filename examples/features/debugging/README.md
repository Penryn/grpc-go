# Debugging

目前，gRPC 提供了两种主要工具来帮助用户调试问题，分别是日志记录和 Channelz。

## 日志
gRPC 在 gRPC 的关键路径上放置了大量的日志记录工具，以帮助用户调试问题。
[日志级别](https://github.com/grpc/grpc-go/blob/master/Documentation/log_levels.md) 文档描述了每个日志级别在 gRPC 上下文中的含义。

要打开调试日志，请使用以下环境变量运行代码：
`GRPC_GO_LOG_VERBOSITY_LEVEL=99 GRPC_GO_LOG_SEVERITY_LEVEL=info`。

## Channelz
我们还提供了一个运行时调试工具 Channelz，以帮助用户进行实时调试。

请参阅此处的 Channelz 博客文章 ([链接](https://grpc.io/blog/a-short-introduction-to-channelz/)) 了解如何使用 Channelz 服务调试实时程序的详细信息。

## 试一试
该示例能够展示日志记录和 Channelz 如何帮助调试。有关完整说明，请参阅上面链接的 Channelz 博客文章。

```
go run server/main.go
```

```
go run client/main.go
```

