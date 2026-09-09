# 更新日志

本项目的所有重要更改都将记录在此文件中。

格式基于 [Keep a Changelog](https://keepachangelog.com/zh-CN/1.1.0/)，版本号遵循 [语义化版本](https://semver.org/lang/zh-CN/)。

## [未发布]

### 新增

- 日本、巴西数据中心
- Device Token 换取 Registration ID 与 App VIP 状态 API

### 变更

- 补齐 Push、Schedule、Status、Plan 和 Group Push 协议字段
- Voice 改为官网 multipart 文件协议，OPPO Image 改为官网 JSON URL 协议
- Tag 计数/配额及 Plan Detail 查询参数改为官网协议

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
