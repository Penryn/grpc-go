# Wait for ready example

此示例展示了如何在 RPC 调用中启用“等待就绪”。

此代码启动一个延迟 2 秒的服务器。如果未启用“等待就绪”，则 RPC 会立即失败并返回 `Unavailable` 代码（情况 1）。如果启用了“等待就绪”，则 RPC 会等待服务器。如果在服务器可用之前上下文失效，则会返回 `DeadlineExceeded` 失败（情况 3）。否则，它会成功（情况 2）。

## 运行示例

```
go run main.go
```

