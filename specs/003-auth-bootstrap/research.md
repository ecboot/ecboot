# Research: 认证引导纵切片（Phase 0）

2026-09-20。spec 无 NEEDS CLARIFICATION；本文件定 8 项实现决策。

## D1 图形验证码：AWT 自绘，零第三方

- **Decision**: `ImageCaptchaGenerator` 用 AWT 绘制 4 位字符 + 干扰线/噪点，输出 Base64 PNG；不引入验证码库。
- **Rationale**: 简单优先/零新依赖（新第三方须改 BOM+契约）；AWT 自绘 60 行内覆盖需求；组件接口抽象，未来可无痛换 aj-captcha 等。
- **Alternatives**: easy-captcha/aj-captcha（被拒：当前需求不值得引依赖；接口已抽象，替换成本仅实现类）。

## D2 验证码存储与一次性语义：Redis + GETDEL

- **Decision**: 图形码 `captcha:img:{ticket}`（值=小写答案，TTL=图形有效期）；短信码 `captcha:sms:{phone}`（值=6 位码）。校验一律 `GETDEL`（原子取删=天然一次性）；频控计数 `captcha:limit:{phone}`（INCR+TTL）。
- **Rationale**: GETDEL 一条命令解决"校验+作废"原子性，无需 Lua；Redis 已在栈内（infra-core 补 starter-data-redis 依赖，契约登记）。
- **Alternatives**: DB 存验证码（被拒：高频短生命周期数据不宜入库）；先 GET 再 DEL（被拒：非原子，并发下凭证可复用）。

## D3 会话：opaque token + Redis，滑动续期

- **Decision**: 登录发放 64 位 hex token（SecureRandom）；`session:{token}` → userId（TTL 默认 7 天）；认证过滤器校验时 EXPIRE 续期（滑动）；登出 DEL。
- **Rationale**: 不透明令牌服务端全可控（登出即失效/封号即失效），免 JWT 签名/续期/黑名单复杂度；JWT 为演进路径（多端/无状态扩展时）。
- **Alternatives**: JWT（被拒 YAGTI：单端小程序场景服务端会话最简）；Spring Session（被拒：其核心价值是分布式 servlet 会话，我们是 API 令牌）。

## D4 手机号加密：AES-256-GCM + 盐化 SHA-256 哈希

- **Decision**: `PhoneCipher`（infra-core）：密文=AES-256-GCM（随机 IV 前置，密钥经 `ecboot.security.phone-key` Base64 配置注入）；哈希=SHA-256(密钥派生盐 + 手机号) 落 `phone_hash`。
- **Rationale**: GCM 认证加密防篡改；哈希用密钥派生盐（非静态盐）防彩虹表/字典；密钥只在配置（不入库不入日志）。
- **Alternatives**: AES-CBC（被拒：无认证）；手机号整体 SHA-256 无加密（被拒：后台展示需可还原）。

## D5 微信登录：WxClient 接口 + Mock 实现（配置开关）

- **Decision**: `WxClient` 接口（code2session 返回 openid/unionid；getPhoneNumber）两实现：`MockWxClient`（`ecboot.weixin.mock-enabled=true` 时装配：openid=`mock-{code}`，手机号=命令中显式传入的开发字段）；真实现留空位。
- **Rationale**: 归并规则（手机号优先）与渠道解耦——mock 只喂身份数据，规则全真；真实微信接入另立特性时只补一个实现类。
- **Alternatives**: 直接接微信 SDK（被拒：需真实 appid/网络，开发态不可测）。

## D6 模块落位（enforcer 边界内逐一核验）

- **Decision**（表：能力→落位→边界核验）：
  | 能力 | 落位 | 核验 |
  |---|---|---|
  | ApiResponse/ErrorCode/异常 | common | 零容器依赖 ✓ |
  | 验证码/短信/加密/会话组件 | infra-core | 白名单{common} ✓（Redis 依赖为技术栈） |
  | 全局异常/TraceId/@CurrentUser/安全白名单 | api-webmvc | 禁服务 ✓（认证过滤器只调 infra 组件） |
  | 验证码端点 | api-common | 禁服务 ✓（只调 infra-core） |
  | 登录应用服务/实体/仓储 | service-user | DDD 四层 ✓ |
  | auth 端点 | api-user | 白名单含 service-user/webmvc ✓ |
  | 配置（密钥/白名单/mock） | start application.yaml | 配置收敛 ✓ |
- **Rationale**: 与既有契约矩阵零冲突；登录端点归 api-user（会员域），api-common 只放无业务归属的验证码。
- **Alternatives**: auth 端点放 api-common（被拒：登录是会员业务，api-common 禁服务而登录需要 service-user）。

## D7 错误码分段

- **Decision**: `0`=成功；`10001+` 通用（参数/系统/未登录/频控）；`20001+` 用户域（验证码错误/手机号锁定/账号禁用/微信绑定冲突/归并引导）。分域段位使新域扩展不冲突（商品域 3xxxx 预留）。
- **Rationale**: 错误码是前后端契约，段位即命名空间；与统一响应（code/message/data）配套。

## D8 mock 短信码的测试取码通道

- **Decision**: mock 模式下 `MockSmsSender` 将码写入 `mock:sms:{phone}`（TTL 同验证码）；api-common 提供 `GET /api/captcha/sms/mock-latest`（**仅 `ecboot.weixin.mock-enabled` 或独立 `ecboot.mock-mode` 为 true 时条件装配**）。
- **Rationale**: 端到端测试/联调需要取码，日志解析不可靠；条件装配保证生产（mock 关闭）端点物理不存在。
- **Alternatives**: 响应体带码（被拒：污染正式契约）；测试直连 Redis（被拒：quickstart 无法给外部联调者用）。

## 关键事实核对

- user 表结构（V1+V11+V19+V27/V28）：phone 密文/phone_hash/wx_openid/wx_unionid/status/last_login_at/last_active_at/register_channel——登录链路所需列全部就绪，零迁移 ✓
- system_config 已有（V30）：阈值读此表（30 分钟生效=应用层缓存 TTL ≤30min）
- enforcer：api-common 禁 service-*、api-user 禁兄弟渠道、start 禁服务直连——D6 落位全部合法 ✓
