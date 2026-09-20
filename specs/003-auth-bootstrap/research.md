# Research: 认证引导纵切片（Phase 0）

2026-09-20。8 项实现决策：

## D1 图形验证码：Go 自绘（x/image 字体），唯一新依赖

- **Decision**: `image` 包自绘 4 位字符 + 干扰线，`image/png` 编码 Base64；字体用 `golang.org/x/image/font/basicfont`。
- **Rationale**: 零第三方商业库；x/image 为半官方扩展库（登记 go.mod，用途=验证码字体）；接口抽象，未来可换 base64Captcha 等。
- **Alternatives**: github.com/mojocn/base64Captcha（被拒：当前需求不值得引第三方；实现类替换成本低）。

## D2 验证码存储：gredis + GETDEL 原子一次性

- **Decision**: 键模型不变（`captcha:img:{ticket}`/`captcha:sms:{phone}`/频控计数）；校验用 `gredis.Do("GETDEL", key)`；重发间隔 `SET NX EX`。
- **Rationale**: GoFrame 自带 gredis（栈内），不引 go-redis；GETDEL 原子性语义同 Redis 命令。
- **Alternatives**: 引 go-redis（被拒：gredis 覆盖需求，避免双客户端）。

## D3 会话：opaque token + gredis 滑动续期

- **Decision**: `crypto/rand` 32 字节 hex token；`session:{token}` → userId，TTL 默认 7 天（system_config）；认证中间件校验时 `EXPIRE` 续期；登出 `DEL`。
- **Rationale**: 不透明令牌服务端全可控；JWT 为演进路径。
- **Alternatives**: JWT（被拒 YAGNI）。

## D4 手机号加密：AES-256-GCM + 盐化 SHA-256（标准库实现）

- **Decision**: `crypto/aes` GCM（随机 nonce 前置）；哈希 `crypto/sha256`（密钥派生盐+手机号）；密钥经 GoFrame 配置 `security.phoneKey`（Base64）注入。
- **Rationale**: Go 标准库全量覆盖，无新依赖；GCM 认证加密防篡改。

## D5 微信登录：WxClient 接口 + Mock

- **Decision**: `service/user` 内定义 `WxClient` 接口；`MockWxClient`（配置 `wechat.mockEnabled=true` 装配：openid=`mock-{code}`，手机号取命令开发字段）。
- **Rationale**: 归并规则与渠道解耦；真实微信接入另立特性只补实现。

## D6 模块落位（按工程实际 GoFrame 结构）

| 能力 | 落位 |
|---|---|
| 验证码/登录接口定义 | `api/common/v1`、`api/user/v1`（req/res 结构体） |
| 控制器（参数绑定/组装） | `internal/controller/{common,user}` |
| 业务实现（登录/归并/休眠核身/WxClient mock） | `internal/service/user`（+`internal/logic` 按需） |
| 验证码/短信/加密/会话组件 | `internal/library` 或 `internal/service` 公共部分（实现时定） |
| 认证中间件/统一响应 | `internal/middleware` |
| 路由注册 | `internal/routes`（或 cmd 内） |
| 配置 | `manifest/config` + `hack/config.yaml`（gf gen） |

## D7 错误码分段

`0` 成功；`10001+` 通用（参数/系统/未登录/频控）；`20001+` 用户域——见 contracts 错误码表。

## D8 mock 短信取码通道（条件路由）

- **Decision**: mock 模式 `MockSmsSender` 写 `mock:sms:{phone}`；`api/commonc` 在 `mockEnabled` 配置为真时**注册**取码路由（生产物理不存在该路由）。
- **Rationale**: GoFrame 路由按条件注册天然支持；联调/端到端取码必需。

## 关键事实核对

- 70 表就绪（迁移 000001~000031 重放实证）；user 表登录链路所需列全齐（密文/哈希/openid/休眠时间戳/share_code）
- system_config（000030）已有读取基础；本切片新增 7 配置种子（000032）
- 落位遵循工程既有目录（api 定义/controller/service/repository/middleware），分层由目录约定与评审约束
- 现代规范（use-modern-go 已加载）：JSON 字段 omitzero、errors.Is/As、any、context 传递等在实现中遵循
