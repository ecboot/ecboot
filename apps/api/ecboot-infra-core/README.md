# ecboot-infra-core（技术组件容器）

基础层的技术组件容器——**横切技术能力的实现与适配**，业务模块通过它获得基础设施而不直连第三方 SDK。

## 功能内容

| 功能域 | 内容 | 承接方 |
|---|---|---|
| 验证码组件 | 图形验证码生成/校验；短信验证码生成/存储/校验（Redis TTL，防刷计数） | `ecboot-api-common` 渠道暴露为 REST |
| 短信组件 | 短信发送适配（V1 对接云短信服务；SDK 版本在 BOM 登记，接口抽象便于替换厂商） | 同上 |
| 缓存组件 | Redis 缓存封装（键规范/过期策略）、分布式锁（呼应库存 Redis 挡板方案） | 全部领域服务 |
| 幂等组件 | 幂等键（Redis `SET NX`）——下单 token、回调去重的公共实现 | 交易/支付链路 |
| 序列化配置 | 统一 JSON 序列化（日期/Long→String 防前端精度丢失等全局规则） | 全渠道 |
| 对象存储 | OSS/CDN 文件上传适配（商品图/评价图，只回传 URL） | 商品/评价域 |

## 职责边界

- **做**：技术实现与第三方适配；自动配置基座（`spring-boot-starter` + `configuration-processor`）
- **不做**：不含业务语义（验证码是技术能力，"注册发码"流程在领域服务）；不依赖任何领域服务/渠道（enforcer 全禁，白名单仅 `ecboot-common`）

## 依赖关系

- 白名单：`ecboot-common`
- 被依赖：api-webmvc / service-* / api-common 渠道

相关文档：`../README.md`、ADR-0001（库存 Redis 挡板）、`docs/schema-design.md`
