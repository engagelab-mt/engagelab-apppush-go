# 更新日志

本项目的所有重要更改都将记录在此文件中。

格式基于 [Keep a Changelog](https://keepachangelog.com/zh-CN/1.1.0/)，版本号遵循 [语义化版本](https://semver.org/lang/zh-CN/)。

## [0.1.1]

### 新增

- 新增日本、巴西数据中心常量，可按应用所属接入点选择对应 AppPush 地址。
- 新增 `Device.RegisterToken`，适配 `POST /v4/devices/token/registration_id`；请求支持 `platform`、`tokens`、`apns_production`，结果支持每个 Token 的 `registration_id`、`is_new`、`code`、`message`。
- 新增 `App.GetVIPStatus` 并挂载到 `Client.App`，映射 `vip_status` 和 `vip_end_time`。

### 变更

- Push 请求模型补齐 `body.voip`、Android `badge_set_num/is_fold`、Message `test_message/receipt_id`、Options `auto_truncation`。
- `notification.alert`、Android `alert` 和 `message.msg_content`支持官网定义的字符串或 JSON Object；VoIP 与厂商扩展继续使用动态 map。
- Batch Push 响应补齐逐目标 `error.code/error.message`及顶层`rate_limit_info`，可以识别 HTTP 200 下的部分失败和限流。
- Group Push 按顶层动态 AppKey 解析成功与失败结果，分别写入 `Successes`、`Errors`，并保留 `GroupMsgID`。
- Device 标签更新恢复强类型 `Tags *DeviceSetTags`；`ClearTags`为 `true` 时序列化为`tags:""`以清空全部标签。
- Schedule 新增 `TriggerIntelligent.backup_time`，创建和更新定时任务均复用补齐后的 Push 模型。
- Tag-device 查询改为返回`TagStatusGetResult.Result`；Tag 计数和配额接口改为`[]string tags + 单个 platform`，tags 使用 repeated-key query。
- Status Detail 补齐 `plan_id`、`pushContent`、`live_activity`、`voip`、`inapp_message`、`sub_hmos`及 HMOS 渠道数据；消息生命周期补齐`error_code/itime/channel`。
- Batch Message Lifecycle 返回类型修正为`[]MessageLifecycleGetResult`，并补齐`message_id`和`registration_id`。
- Plan Detail 查询参数修正为`plan_ids/start_date/end_date`；Push Plan 列表字段由错误的`push_id`改为`plan_id`并增加`entity_tag`。
- Voice `Create`由 JSON 文本参数改为官网`language + file` multipart 上传；列表改为数组响应，单项结果补齐`file_url`。
- OPPO Image 由 multipart 文件上传改为 JSON URL 请求，使用`big_picture_url/small_picture_url`并映射`big_picture_id/small_picture_id`；移除不符合官网协议的`UploadOppoFromReader`。
- Device Token 和 OPPO Image 的数量、平台、条件及二选一业务约束交由服务端校验，SDK 沿用`ApiError`返回错误。

### 修复

- 修正 Tag、Status、Plan、Voice、Image 和 Group Push 中与官网不一致的 query、body 与响应层级。
- 补充请求序列化、动态 AppKey、部分限流、响应解析及错误响应测试，并通过`go test ./...`和`go vet ./...`。

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
