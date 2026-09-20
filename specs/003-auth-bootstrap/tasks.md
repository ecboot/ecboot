---
description: "Task list for 认证引导纵切片（Go 版实现）"
---

# Tasks: 认证引导纵切片（auth-bootstrap）

**Input**: Design documents from `/specs/003-auth-bootstrap/`（spec、plan/research/data-model/quickstart 均为 Go 版、contracts/api-contracts.md）

**Prerequisites**: 特性 004 已合入——契约定义（api/user/v1/auth.go、api/common/v1/captcha.go|sms.go）、路由分组、Auth 白名单占位、errcode 常量、统一响应中间件**均已就绪**；本特性交付其背后的业务实现。

**Tests**: quickstart 四场景为验收（端到端 curl/gtest 序列）；组件级单测（验证码一次性/频控/加密往返）随任务交付。

**Organization**: 按 spec 用户故事 US1~US5；落位遵循工程既有 GoFrame 结构——业务在 `internal/service/user`，技术组件在 `internal/library/{captcha,sms,security}`（新建），中间件在 `internal/middleware`。

## Format: `[ID] [P?] [Story] Description`

---

## Phase 1: Setup (Shared Infrastructure)

- [x] T001 编写迁移 `migrations/000033_auth_config.up.sql`：system_config 种子 7 行（captcha.image.ttl_seconds=300 / captcha.sms.ttl_seconds=300 / captcha.sms.resend_seconds=60 / captcha.fail.max=5 / captcha.fail.lock_seconds=1800 / session.ttl_days=7 / dormant.tier1.days=90，见 data-model §四）+ down 占位
- [x] T002 [P] 新建 `internal/config/`（或 bootstrap 注入）运行时配置读取：security.phone-key（Base64 AES 密钥，经 env JWT_SECRET 同款机制注入）、wechat.mock-enabled、captcha 相关默认值兜底；`application.yaml`/config.yaml 增对应键

**Checkpoint**: 配置与迁移就绪

---

## Phase 2: Foundational (Blocking Prerequisites)

- [x] T003 实现统一错误出口：`internal/library/err`（或 service 内）BusinessException（code+message+ wraps error）；全局恢复中间件（panic→10002+追踪号）、gf 校验失败→10001（message 含字段）、BusinessException→其 code 的映射，挂载于根路由组（替代桩的 gcode.CodeNotImplemented 裸露）
- [x] T004 [P] TraceId 中间件：请求生成/透传 X-Trace-Id，写入响应头与日志上下文

**Checkpoint**: 三类错误注入响应契约符合（spec US1 验收）

---

## Phase 3: User Story 1 - 统一 Web 基座收尾 (Priority: P1) 🎯 MVP

**Goal**: US1 基座（响应结构/异常映射/追踪号）在 004 骨架上完成最后实现

**Independent Test**: quickstart——任一桩端点返回统一结构；构造 panic 返回 10002+追踪号

- [x] T005 [US1] 三类错误注入走查并修复偏差：业务异常/参数校验/panic 各自 code 与 message 断言（gtest 或 curl 序列留证）

**Checkpoint**: 统一契约零例外（SC-002 基线）

---

## Phase 4: User Story 2 - 图形验证码 (Priority: P1)

**Goal**: 图形验证码获取与一次性校验

**Independent Test**: 获取→正确答案通过→同 key 复用失败；过期失败

- [x] T006 [P] [US2] `internal/library/captcha`：image.go 自绘 4 位字符+干扰线（`image/png` Base64 输出，x/image 字体）；service.go Get(key)/Verify(原子 GETDEL)、TTL 读 system_config
- [x] T007 [US2] 实现 `api/common/v1` GetCaptcha 与 VerifyCaptcha 控制器→服务接线（替换桩）；单元测试：一次性/过期/忽略大小写

**Checkpoint**: SC-003 图形码部分（一次性/过期正反用例）

---

## Phase 5: User Story 3 - 手机号注册登录 (Priority: P1) 🎯 核心

**Goal**: 短信码发送（mock）+ 注册即登录全链路

**Independent Test**: quickstart 场景一全链路 + 场景二防刷用例

- [x] T008 [P] [US3] `internal/library/sms`：Sender 接口 + MockSmsSender（写 `mock:sms:{phone}` + 日志输出；真实渠道另立特性）
- [x] T009 [P] [US3] `internal/library/security/phone_cipher.go`：AES-256-GCM 密文（随机 nonce 前置）+ 盐化 SHA-256 哈希（密钥派生盐）；单测：往返/密文随机性/哈希稳定
- [x] T010 [US3] GetSmsCode 实现：手机号格式入口校验（非法不入流程）→图形码校验→60s 重发锁（SET NX EX）→6 位码生成落 Redis→Mock 发送
- [x] T011 [US3] `internal/service/user` DDD 落位：domain（User/UserLoginLog 实体映射既有表、Repository 接口）+ infrastructure（gf dao 之上封装 FindByPhoneHash/FindById/CreateLogin/TouchLoginTimes/AppendLoginLog）+ repository 接线
- [x] T012 [US3] AuthAppService.SmsLogin：决策树（data-model §三）——码校验(GETDEL)→hash 查询→禁用拒绝(20003)→休眠分级（≥阈值且非完整核验→拒绝）→更新时间戳+写日志→发会话；并发注册唯一索引兜底转登录
- [x] T013 [US3] SmsLogin 控制器接线 + quickstart 场景一（1-5 步）走通 + 场景二防刷（一次性/60s 频控/5 次锁定）留证

**Checkpoint**: SC-001/SC-003/SC-005/SC-007（库内零明文 SQL 断言）

---

## Phase 6: User Story 4 - 会话凭证生命周期 (Priority: P1)

**Goal**: 不透明令牌双凭证会话 + 受保护端点真实鉴权

**Independent Test**: quickstart 场景一第 6 步（me 可达/登出后 10003）

- [x] T014 [P] [US4] `internal/library/security/session.go`：SessionManager——Create（crypto/rand 64hex）→`session:{token}`→userId（TTL=配置）；Refresh（refreshToken 换新双凭证，旧全失效）；Validate（存在即续期滑动）；Destroy（登出）。refreshToken 独立键 `session:refresh:{token}` 长效
- [x] T015 [US4] `internal/middleware/auth.go` 真实化：会员组 Validate+续期+userId 注入 ctx；管理组占位（后台登录属后续）；公开白名单维持
- [x] T016 [US4] Me/Logout/TokenRefresh 控制器接线 + quickstart 场景一第 6 步断言（me→logout→me=10003）

**Checkpoint**: SC-001 全链路 + 会话三态（有效/登出/过期）

---

## Phase 7: User Story 5 - 微信登录（mock 归并） (Priority: P2)

**Goal**: mock 微信身份登录 + 手机号优先归并

**Independent Test**: quickstart 场景三（归并同账号/直接登录/绑定冲突）

- [x] T017 [US5] `internal/service/user/wx.go`：WxClient 接口 + MockWxClient（wechat.mock-enabled=true 时 openid=`mock-{wxCode}`，手机号取命令开发字段）；归并决策树（data-model §三微信分支：命中直接登录/绑定/20004 冲突引导人工）
- [x] T018 [US5] WxLogin 控制器接线 + quickstart 场景三三用例（归并同号、二次直登、冲突 20004）留证；DB 断言无双账号（SC-004）

**Checkpoint**: US5 独立闭环

---

## Phase 8: Polish & Cross-Cutting Concerns

- [x] T019 休眠分级断言（quickstart 场景四）：SQL 拨表 91 天→微信静默路径被拒引导短信通道；阈值改 system_config 后生效（含 30 分钟缓存语义说明）
- [x] T020 quickstart 全场景终验 + 宪法 IV（`make build && make test`）输出留证到本目录 verification.md
- [x] T021 文档收尾：CONTEXT.md/模块 README 如有术语/行为补充则同步；errors/权限点清单核对

---

## Dependencies & Execution Order

- Setup（T001/T002）→ Foundational（T003/T004）→ US1（T005）→ US2（T006~T007）→ US3（T008~T013）→ US4（T014~T016）→ US5（T017~T018）→ Polish（T019~T021）
- US2 可与 T008/T009 并行（图形码与短信/加密组件互不依赖）；US4 依赖 T012（登录发会话）；US5 依赖 T012（归并复用登录路径）
- [P] 任务：不同文件可并行

### Parallel Opportunities

- T006/T008/T009 三个技术组件互不依赖可同时开工；T012 依赖 T009/T011

---

## Implementation Strategy

### MVP First

Setup → Foundational → US1 → US2 → US3（T008~T013）= 用户可注册登录的最小闭环（无会话持久验证前的临时 token 形态由 T014 立即补齐）。

### Incremental Delivery

US3+US4 合并交付即"用户能进门且保持进门"；US5 补双通道；每故事一次中文提交。

---

## Notes

- **终验（2026-09-20）**：TestAuthFlow 端到端 PASS——图形码一次性→发码→注册即登录（库内零明文+会话有效）→微信归并同账号→绑定冲突 20004→重复登录幂等→休眠拨表 91 天微信静默被拒+短信通道放行。运行时 /common/captcha 真实输出 PNG Base64。修复记录：contrib 驱动注册（mysql/redis）、phone key 32 字节、gvalid/gerror API 适配。

- 手机号明文只存在于请求与内存；落库一律经 PhoneCipher（SC-007）
- 所有阈值走 system_config，代码默认值兜底（配置缺失不阻断）
- middleware/auth.go 的管理组鉴权占位维持（后台登录属 0010 后台特性）
- 验证以输出为证（quickstart 场景 + 单测）；每任务组一次提交
