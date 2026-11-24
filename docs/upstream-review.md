# 上游审查记录

2026-11-24，Windows Codex；对应 [#1](https://github.com/LazyBuddyDev/openhongguo/issues/1) 和 [#2](https://github.com/LazyBuddyDev/openhongguo/issues/2)。

源基线为 `015a388112ee148a7d2dfb7ecfafe171bbab957f`，本地定制基线为 `7aa0121`。
[#1 的比较范围](https://github.com/beilunyang/openhongguo/compare/015a388112ee148a7d2dfb7ecfafe171bbab957f...33c888d3a8a2db3ea1ff2f1dce5c479c15bf7626)有 25 个提交；[#2 的增量](https://github.com/beilunyang/openhongguo/compare/33c888d3a8a2db3ea1ff2f1dce5c479c15bf7626...cac5dbc0e1150c0f93e309d8de68bb483f50f499)有 10 个提交。
#2 包含 #1，两条合计审查 35 个提交。工作流 `REVIEWED` 推进到 `cac5dbc0e1150c0f93e309d8de68bb483f50f499`，表示已作取舍，不表示全部合入或已通过线上验收。

## 选择性合入

| 上游提交 | 结果与验证 |
| --- | --- |
| `8214376` | 过滤终端输出中的控制及双向文字控制字符，避免弹幕写入终端转义序列；先运行原测试观察失败，再应用实现并通过。 |
| `7fa618c` | 演员作品先查询实际标签布局，分页游标保留标签，兼容旧游标；同步两种布局样例、SDK/契约测试及各绑定类型注释。原 SDK 测试复现首次请求数和无作品标签行为错误，修复后通过。 |
| `6919947` | 搜索联想前缀由 `query` 改为 `q`；原回归测试确认缺少 `q`，修改后通过。上游关于实时接口兼容性的描述未在本机重新验证。 |
| `df0fa78` | 同步演员榜不再列出子榜单的样例与断言；解析行为未改动。服务端变化来自上游记录，本机只验证样例。 |

保留现有账号收藏/进度接口、固定 player GET 取源、清晰度选择、禁止静默降级和错集拒绝逻辑。未改动 hongguo-service 的 vendored 源码、已安装模块或宿主。

## 不合入的更新

| 范围与上游提交 | 决定与依据 |
| --- | --- |
| V5 播放：`788103f`、合并提交 `999f57b`、`2a261a8`、`88182cd` | 新增 V5 POST、空响应内部重试、旧接口回退及模型/地址处理。`sdk/media.go` 的重试循环不受 `RetryAttempts=1` 控制，且当前服务调用的是定制 `BriefDockPlayback`，直接合入并不能更新其取源路径。整组不合入；后续若接入，应另行确定端点、重试及清晰度边界并验证实际媒体推进。 |
| 远程签名：`955ffe1`、`56dd71a`、`b87d5ce`、`228cf0d`、`6377f3f`、`7546c73`、`20f5906` | 引入 Java/unidbg 服务和新的请求出口。远程签名 transport 将完整请求 URL 与 body 转送到配置服务；当前服务使用受限传输且无需此新增运行时。整组不合入，不自动部署公共或私有签名器。 |
| #2 签名/研究：`e7a90a6`、`89fb6a5`、`1cb153f`、`8134661`、`50baf58`、`979160b`、`08ab6af` | 改写普通请求的 Medusa/Helios/Gorgon、加入设备/进程状态模拟和签名研究工具。普通请求签名会尝试调整 `_rticket` 以命中部分已实现分支；离线向量不能证明账号请求兼容。#2 无独立、已验证可用于当前服务的修复，保持当前签名实现。 |
| 发布/CI：`c1d42d0`、`bf8df22`、`33c888d` | 属于上游自己的多平台发布工作流；本分支只保留自有上游通知工作流，不恢复上游发布渠道。 |
| 版本、发布说明、二维码：`1d3354a`、`2769fb6`、`590a6ab`、`6433b4d`、`d2277dc`、`4ea5e43`、`cd27b9e`、`c6de145`、`32d2a8c`、`cac5dbc` | 未完整升级到对应版本，不修改包版本或引入上游宣传素材。 |

上述不合入项是本次审查的范围决定。它们不构成现有播放空响应的已验证根因或修复；已暂停的实播排查保持原状态。

## 验证

使用红果服务已有 Go 工具链，网络请求由测试样例/本地测试服务器替代，未使用真实账号。

以下检查通过：

```bash
go test ./internal/protocol/... ./sdk/... ./internal/contract ./bindings/internal/... -count=1
go test ./internal/cliapp -skip TestCLIAnonymousReadsUseShippedClient -count=1
go vet ./internal/protocol/... ./sdk/... ./internal/contract ./bindings/internal/... ./internal/cliapp
```

覆盖两种演员页布局、无作品标签、游标标签/旧格式、异人游标拒绝、搜索参数、控制字符过滤，以及现有播放解析和定制清晰度测试。

修改前 `go test ./...` 已失败：Windows 不提供测试期待的 POSIX `0600`/符号链接行为，媒体测试使用 Unix 假执行脚本，现有生成文件也与定制 SDK 不同步（包括 docgen 未登记的 BriefDock 方法和新增清晰度字段）。本次不扩大范围修改这些既有问题。CLI 定向测试仅跳过基线已确认依赖 POSIX 文件权限的 `TestCLIAnonymousReadsUseShippedClient`。未声称全量测试通过。

未运行真实接口、实际媒体或其他平台绑定运行时验收；未构建/发布红果模块。本记录不能证明用户原有播放问题已解决。
