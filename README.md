# EngageLab AppPush Go SDK

EngageLab AppPush REST API 的 Go 语言 SDK，零第三方依赖，仅使用 Go 标准库。

## 安装

```bash
go get github.com/engagelab/engagelab-go
```

## 快速开始

```go
package main

import (
    "context"
    "fmt"
    "log"

    engagelab "github.com/engagelab/engagelab-go"
)

func main() {
    // 创建客户端（默认新加坡数据中心）
    client := engagelab.NewClient("your-app-key", "your-master-secret")

    // 或指定数据中心
    client = engagelab.NewClient("your-app-key", "your-master-secret",
        engagelab.WithDataCenter(engagelab.HongKong),
    )

    ctx := context.Background()

    // 发送推送
    result, err := client.Push.Send(ctx, &engagelab.PushParam{
        From: "push",
        To:   "all",
        Body: &engagelab.PushBody{
            Platform: "all",
            Notification: &engagelab.NotificationMessage{
                Alert: "Hello from Go SDK!",
                Android: &engagelab.AndroidNotification{
                    Alert: "Hello Android!",
                    Title: "Test Push",
                },
                IOS: &engagelab.IOSNotification{
                    Alert: "Hello iOS!",
                },
            },
        },
    })
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("Push sent: msg_id=%s\n", result.MsgID)
}
```

## 数据中心

| 常量 | 区域 | Base URL |
|------|------|----------|
| `Singapore` | 新加坡 (默认) | `https://pushapi-sgp.engagelab.com` |
| `HongKong` | 香港 | `https://pushapi-hk.engagelab.com` |
| `Virginia` | 美国弗吉尼亚 | `https://pushapi-usva.engagelab.com` |
| `Frankfurt` | 德国法兰克福 | `https://pushapi-defra.engagelab.com` |

## API 模块

### Push — 推送

```go
client.Push.Send(ctx, param)           // 创建推送
client.Push.SendRaw(ctx, rawBody)      // 自定义 JSON 推送
client.Push.Validate(ctx, param)       // 推送校验
client.Push.Withdraw(ctx, msgID)       // 消息撤回
client.Push.BatchByRegID(ctx, param)   // 按 Registration ID 批量推送
client.Push.BatchByAlias(ctx, param)   // 按 Alias 批量推送
```

### Group Push — 应用分组推送

```go
// Group Push 使用独立认证: group-{GroupKey}:{GroupMasterSecret}
groupClient := engagelab.NewGroupPushClient("group-key", "group-master-secret",
    engagelab.WithDataCenter(engagelab.Singapore),
)
groupClient.Send(ctx, param)
```

### Device — 设备

```go
client.Device.Get(ctx, registrationID)         // 查询设备信息
client.Device.Set(ctx, registrationID, param)  // 设置设备标签/别名
client.Device.Delete(ctx, registrationID)      // 删除设备
client.Device.GetStatus(ctx, param)            // 查询设备在线状态
```

### Tag — 标签

```go
client.Tag.List(ctx)                                    // 获取标签列表
client.Tag.Set(ctx, tag, param)                         // 添加/移除标签设备
client.Tag.Delete(ctx, tag, platforms)                   // 删除标签
client.Tag.GetCount(ctx, tags, platforms)                // 查询标签设备数
client.Tag.GetDeviceStatus(ctx, tag, registrationID)     // 查询设备标签绑定状态
client.Tag.GetQuota(ctx, tags, platforms)                // 查询标签配额
```

### Alias — 别名

```go
client.Alias.Get(ctx, alias, platforms)     // 查询别名设备
client.Alias.Delete(ctx, alias, platforms)  // 删除别名
```

### Schedule — 定时任务

```go
client.Schedule.Create(ctx, param)          // 创建定时推送
client.Schedule.Update(ctx, id, param)      // 更新定时推送
client.Schedule.Delete(ctx, id)             // 删除定时推送
client.Schedule.Get(ctx, id)                // 获取定时推送详情
client.Schedule.List(ctx, page)             // 获取定时推送列表
client.Schedule.GetMsgIDs(ctx, id)          // 获取定时推送消息 ID
```

### Status — 统计

```go
client.Status.Users(ctx, timeUnit, start, duration)         // 用户统计
client.Status.MessageDetail(ctx, messageIDs)                 // 消息送达统计
client.Status.MessageLifecycle(ctx, msgID, registrationIDs)  // 消息生命周期
client.Status.BatchMessageDetail(ctx, messageIDs)            // 批量消息统计
client.Status.PlanDetail(ctx, planID, messageIDs)            // 推送计划统计
```

### Plan — 推送计划

```go
client.Plan.CreateOrUpdate(ctx, param)                         // 创建/更新推送计划
client.Plan.List(ctx, pageIndex, pageSize, sendSource, desc)   // 查询推送计划列表
client.Plan.QueryMsg(ctx, planIDs, startDate, endDate)         // 查询计划消息 ID
client.Plan.Delete(ctx, planID)                                // 删除推送计划
client.Plan.BatchDelete(ctx, planIDs)                          // 批量删除推送计划
```

### Voice — 语音/TTS

```go
client.Voice.Create(ctx, param)       // 创建语音模板
client.Voice.List(ctx)                // 获取语音模板列表
client.Voice.Get(ctx, language)       // 获取语音模板
client.Voice.Delete(ctx, language)    // 删除语音模板
```

### Image — 图片

```go
client.Image.UploadOppo(ctx, filePath)                    // 上传 OPPO 大图 (文件路径)
client.Image.UploadOppoFromReader(ctx, filename, reader)  // 上传 OPPO 大图 (io.Reader)
```

## 错误处理

SDK 使用 `*engagelab.ApiError` 类型返回 API 错误：

```go
result, err := client.Push.Send(ctx, param)
if err != nil {
    var apiErr *engagelab.ApiError
    if errors.As(err, &apiErr) {
        fmt.Printf("API Error: status=%d code=%d message=%s\n",
            apiErr.StatusCode, apiErr.ErrorBody.Code, apiErr.ErrorBody.Message)
    } else {
        fmt.Printf("Network Error: %v\n", err)
    }
}
```

## 客户端配置

```go
// 自定义 HTTP 客户端
client := engagelab.NewClient("key", "secret",
    engagelab.WithHTTPClient(&http.Client{
        Timeout: 60 * time.Second,
        Transport: &http.Transport{
            MaxIdleConns: 100,
        },
    }),
)

// 仅修改超时
client = engagelab.NewClient("key", "secret",
    engagelab.WithTimeout(60 * time.Second),
)

// 自定义 Base URL
client = engagelab.NewClient("key", "secret",
    engagelab.WithBaseURL("https://custom-api.example.com"),
)
```

## 认证方式

| API 模块 | 认证 |
|---------|------|
| Push / Device / Tag / Alias / Schedule / Status / Plan / Voice / Image | Basic Auth: `appKey:masterSecret` |
| Group Push | Basic Auth: `group-{GroupKey}:GroupMasterSecret` |

## 更新日志

详见 [CHANGELOG.md](CHANGELOG.md)。

## 许可证

MIT License
