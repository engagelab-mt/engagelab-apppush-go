# 更新日志

本项目的所有重要更改都将记录在此文件中。

格式基于 [Keep a Changelog](https://keepachangelog.com/zh-CN/1.1.0/)，版本号遵循 [语义化版本](https://semver.org/lang/zh-CN/)。

## [0.1.1]

### 新增

- `DataCenter`：新增 `Japan`和`Brazil`，用于连接部署在日本、巴西数据中心的应用。
- `Device.RegisterToken`：新增厂商 Token 注册能力；`DeviceTokenRegisterResult`可获取每个 Token 对应的 Registration ID、是否首次创建、错误码和错误信息。
- `App.GetVIPStatus`：新增应用 VIP 状态和到期时间查询，通过`Client.App`调用。

### 变更

- `Push.Send`、`Push.Validate`、Batch Push、Group Push 和 Schedule 推送参数：新增`PushBody.VoIP`、Android`BadgeSetNum/IsFold`、Message`TestMessage/ReceiptID`、Options`AutoTruncation`。
- `NotificationMessage`、`AndroidNotification`和`CustomMessage`：通知内容及自定义消息内容支持字符串或结构化 JSON 对象。
- `Push.BatchByRegID`和`Push.BatchByAlias`：`BatchPushResult`新增逐目标错误及限流信息，可识别请求成功时返回的部分失败。
- `GroupPushClient.Send`：`GroupPushResult.Successes`和`Errors`可按 AppKey 获取各应用结果，同时保留`GroupMsgID`。
- `Device.Set`：`DeviceSetParam.Tags`使用强类型`*DeviceSetTags`；设置`ClearTags: true`可清空设备全部标签。
- `Schedule.Create`和`Schedule.Update`：新增`TriggerIntelligent.BackupTime`智能定时配置。
- `Tag.GetDeviceStatus`：改为返回`TagStatusGetResult.Result`，用于判断设备是否拥有指定标签。
- `Tag.GetCount`和`Tag.GetQuota`：参数改为`[]string tags`和单个`platform`。
- `Status.MessageDetail`：结果新增 Plan ID、推送内容、Live Activity、VoIP、应用内消息、HMOS 小计及 HMOS 渠道统计。
- `Status.MessageLifecycle`：结果新增`ErrorCode`、`ITime`和`Channel`。
- `Status.BatchMessageDetail`：返回值改为`[]MessageLifecycleGetResult`，每项包含`MessageID`和`RegistrationID`。
- `Status.PlanDetail`：参数改为`[]string planIDs, string startDate, string endDate`。
- `Plan.List`：`PushPlanInfo`使用`PlanID`并新增`EntityTag`；创建时间和最后使用时间继续使用毫秒时间戳。
- `Voice.Create`：参数改为`language, filePath`并上传本地语音文件；`Voice.List`返回`[]VoiceResult`，`Voice.Get`结果新增`FileURL`。
- `Image.UploadOppo`：参数改为`*OppoImageParam`，通过大图或小图 URL 上传，结果返回对应的图片 ID；不符合要求的参数通过既有`ApiError`返回。
- `Image.UploadOppoFromReader`：已移除，请改用`Image.UploadOppo`。
- `Device.RegisterToken`：数量、平台和 APNs 条件不在客户端提前拦截，服务端错误继续通过`ApiError`返回。

### 修复

- `Tag.GetCount`、`Tag.GetQuota`和`Status.PlanDetail`：修正参数编码，避免服务端收到错误的查询条件。
- `Device.Set`：通过 JSON 构造`DeviceSetParam`时，正确解析标签增删对象和`tags:""`清空标签形式。
- `Status.MessageDetail`：正确解析通知、自定义消息、Live Activity、VoIP 和应用内消息统计中的`target/click`及`targets/clicks`字段。
- `Voice.Create`、`Voice.List`和`Image.UploadOppo`：修正请求或响应格式不一致导致的调用失败、字段丢失问题。
- Tests：补充公开方法的请求序列化、Group Push 多应用结果、Batch Push 部分限流、响应解析及错误响应测试。

## [0.1.0] - 2026-03-20

### 新增

- `NewClient` 客户端初始化，支持 `WithDataCenter`、`WithBaseURL`、`WithHTTPClient`、`WithTimeout` 等选项
- 四大数据中心常量：`Singapore`（默认）、`HongKong`、`Virginia`、`Frankfurt`
- **Push** — 推送服务：`Send`、`SendRaw`、`Validate`、`Withdraw`、`BatchByRegID`、`BatchByAlias`
- **Group Push** — 应用分组推送：独立 `NewGroupPushClient` + `Send`
- **Device** — 设备管理：`Get`、`Set`、`Delete`、`GetStatus`
- **Tag** — 标签管理：`List`、`Set`、`Delete`、`GetCount`、`GetDeviceStatus`、`GetQuota`
- **Alias** — 别名管理：`Get`、`Delete`
- **Schedule** — 定时推送：`Create`、`Update`、`Delete`、`Get`、`List`、`GetMsgIDs`
- **Status** — 统计查询：`Users`、`MessageDetail`、`MessageLifecycle`、`BatchMessageDetail`、`PlanDetail`
- **Plan** — 推送计划：`CreateOrUpdate`、`List`、`QueryMsg`、`Delete`、`BatchDelete`
- **Voice** — 语音/TTS 模板：`Create`、`List`、`Get`、`Delete`
- **Image** — 图片上传：`UploadOppo`、`UploadOppoFromReader`
- `*ApiError` 结构化错误类型，包含 HTTP 状态码和业务错误码
- 完整的单元测试（65 个测试用例，使用 `httptest` 模拟，零外部依赖）

<!--
## 模板

## [x.y.z] - yyyy-mm-dd

### 新增
- 新功能说明

### 变更
- 对已有功能的修改说明

### 废弃
- 即将移除的功能说明

### 移除
- 已移除的功能说明

### 修复
- Bug 修复说明

### 安全
- 安全相关的修复说明
-->
