# ecboot-api-admin（管理后台渠道）

**管理后台** REST API——运营/管理员的全域管理端点，含完整 RBAC 与审计。

## 功能内容（REST 端点规划）

| 功能组 | 端点（规划） | 后端域 |
|---|---|---|
| 认证与权限 | 后台登录（登录审计含失败尝试）、登出；RBAC（角色/权限/账号-角色、权限树管理） | admin_user 系列（V10） |
| 商品管理 | 分类/品牌维护、SPU/SKU 创建编辑、两级上下架、运费模板管理 | service-shop · 目录 |
| 库存管理 | 库存调整（留流水）、库存查询、预警 | service-shop · 库存 |
| 订单管理 | 订单查询/详情、发货（物流公司字典）、取消、导出 | service-shop · 交易/物流 |
| 售后管理 | 待审核列表、同意/拒绝、确认收货退款、重试打款 | service-shop · 售后 |
| 促销管理 | 优惠券模板/发放、满减活动（档位/范围）、拼团/秒杀活动配置 | service-shop · 促销 |
| 分销管理 | 推广员审核/冻结、佣金规则配置、佣金记录查询、提现审核与打款 | service-user · 分销 |
| 会员管理 | 会员查询（手机号**仅精确检索**）、禁用/启用、改绑手机、注销处理 | service-user |
| 风控管理 | 黑名单/规则维护、风控事件查询、申诉处理 | service-shop · 风控 |
| 审计查询 | 操作审计日志（模块/操作人/耗时）、登录审计查询 | admin_operation_log / admin_login_log |
| 运营看板 | 交易/会员/商品基础统计 | 跨域只读查询 |

## 职责边界

- **做**：管理端 REST 暴露；依赖双领域服务（全域管理）；写操作经操作审计切面留痕
- **不做**：与 api-common/user/shop **互不编译依赖**；C 端端点不在本渠道

## 依赖关系

- 白名单：`ecboot-api-webmvc`、`ecboot-service-user`、`ecboot-service-shop`、`ecboot-common`、`ecboot-infra-core`
- 被依赖：仅 `ecboot-start`

相关文档：`../README.md`、V10 迁移（RBAC 七表）
