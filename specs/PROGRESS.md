# 后端业务层收口进度总表（PROGRESS）

> **本文件是"service 补齐 + controller 连线"工作的唯一任务台账**。任何一轮 AI 执行会话开工前必读本文件；
> 台账之外的任务一律不执行。修订本文件须走 §五 变更记录。
>
> 基线：2026-09-21，commit `c482eb0`。分层契约见 `docs/layer-contracts.md`（评审依据，优先级高于本表）。

## 一、执行协议（每轮会话的固定动作）

1. **开工**：读本表 → 取状态为 ⬜ 的最小编号批（或 🔶 批续作）→ 读该批 `specs/0NN-<slug>/` 下的 tasks.md。
2. **批内**：按 tasks.md 任务序推进；service 一律 TDD（红→绿），测试与被测包同目录。
3. **收尾**：跑 §四 DoD 全项 → **更新本表状态（与代码同一 commit）** → commit 格式 `feat(0NN-<slug>): 中文描述`。
4. **遇到范围/依赖问题**：先在 §五 记变更，再动代码；禁止静默扩大范围。

## 二、基线快照（2026-09-21，commit c482eb0）

- **端点总量 218**（admin 129 / shop 42 / user 39 / common 8），API req/res 已全量定义（specs/004-api-surface）。
- **已连线 8 个**（认证纵切片，不进批次）：
  - user：`me` / `logout` / `sms_login` / `wx_login` / `token_refresh`
  - common：`get_captcha` / `verify_captcha` / `get_sms_code`
- **未连线桩 210 个**（`CodeNotImplemented`），分 13 批，见 §三。
- **service 现状**：user/shop/system 三域接口已定义（含业务规则注释）；已实现并有测试：shop 域 `browse / product / cart / order / promotion_calc`（交易链路）与 user 域认证流；其余为"接口已有、impl 待补"，各批启动时在 tasks.md 做精确盘点。
- **生成物现状**：dao/model 全量已生成（69 表，勿手改）；controller 桩全量已生成（`gf gen ctrl`）。
- **验收工具**：`apps/server` 下 `make check-stub` 输出各渠道桩数（无 make 环境可直接执行 Makefile 内等价 shell）。

## 三、批次总表

状态图例：⬜ 未开始 / 🔶 进行中 / ✅ 完成（本批桩清零且 DoD 全过）

| # | 特性目录 | 域 | 端点数 | 依赖 | service 起点 | 状态 | 完成 commit |
|---|---|---|---|---|---|---|---|
| 01 | `specs/007-admin-base` | 后台账户与系统配置（登录/登出/刷新/资料/改密、admin_user、RBAC 角色、config、ping/短信调试桩） | 22 | — | 已实现（含会话 audience 安全修复） | ✅ | `45cd314` |
| 02 | `specs/008-store` | 门店域（自提/核销载体：admin 管理 + common 查询） | 7 | 01 | 接口已有，impl 待补 | ⬜ | |
| 03 | `specs/009-logistics-ops` | 物流公司与运营装修（admin logistics/banner/floor + shop banner/floor） | 15 | 01 | 接口已有，impl 待补 | ⬜ | |
| 04 | `specs/010-product-admin` | 商品目录后台与 C 端浏览收口（admin spu/sku/类目/品牌/库存 + shop 浏览连线） | 27 | 01 | **impl 大半已有**（005 遗产），补库存后台 | ⬜ | |
| 05 | `specs/011-member-center` | 会员中心（资料/地址/收藏/足迹/消息/积分/通知偏好/登录日志/邀请记录） | 21 | 01 | 接口已有，impl 待补 | ⬜ | |
| 06 | `specs/012-trade-wiring` | 交易闭环连线（shop 购物车/订单/支付 + admin 订单发货/取消/备注 + user 优惠券） | 22 | 02 03 04 | **shop 侧 impl 已实现**，重连线 + admin 侧补 impl | ⬜ | |
| 07 | `specs/013-after-sale` | 售后域（shop 申请/撤销售后 + admin 审核/收货/重退款） | 11 | 06 | 接口已有，impl 待补 | ⬜ | |
| 08 | `specs/014-review` | 评价域（shop 发表/追加/我的评价/商品评价列表） | 4 | 06 | 接口已有，impl 待补 | ⬜ | |
| 09 | `specs/015-marketing-c` | 营销 C 端（首页聚合 + 秒杀/拼团/砍价/助力/满减列表与玩法动作） | 12 | 06 | 接口已有，impl 待补 | ⬜ | |
| 10 | `specs/016-marketing-admin` | 营销后台（券/秒杀/拼团/砍价/助力/满减管理） | 30 | 04 09 | 接口已有，impl 待补 | ⬜ | |
| 11 | `specs/017-distribution-fund` | 分销与资金（分销账户/申请/关系/佣金规则、提现、邀请记录、分享归因；两级红线已锁表结构） | 23 | 06 | 接口已有，impl 待补 | ⬜ | |
| 12 | `specs/018-member-audit-risk` | 会员管理与审计风控（admin member/登录与操作日志/风控规则与记录） | 13 | 01 11 | 接口已有，impl 待补 | ⬜ | |
| 13 | `specs/019-dashboard-final` | 看板与收口终验（dashboard 三看板 + 全量桩清零 + 全量回归） | 3 | 全部 | 接口已有，impl 待补 | ⬜ | |
| — | （已连线基线） | 认证纵切片 + 验证码/短信 | 8 | — | ✅ 已实现 | ✅ | `ae7acf5` 等 |

合计：210（待做）+ 8（已连线）= 218 ✓（建表时已按批次映射规则对账：13 批覆盖全部 210 桩，无遗漏无重复）

依赖要点：06（交易）依赖 02/03/04 是因为结算试算消费门店（自提）、运费（物流）、商品与库存数据；
09/10 营销拆两批是因为玩法上下文被下单引用（09 先行），后台管理复用同一批 service（10 随后）。

## 四、防偏离条款（执行宪法）

1. **唯一事实源**：本表 + 各批 tasks.md 是唯一任务清单；每轮会话只做当前批清单内的任务，台账外任务不做（发现遗漏 → 走 §五 记录后再做）。
2. **状态与代码同 commit**：批次状态/勾选更新与对应代码进同一个 commit，杜绝表滞后于代码。
3. **批次边界 = 文件边界**：本批只允许改动本批端点对应的 controller 文件、该批 service/测试文件及本表与该批 specs 目录；需要动他批文件 → 停下，先在 §五 登记并调整批次，再继续。
4. **变更先记账**：任何范围变化（增删端点、调整批次/依赖）先写 §五 变更记录，再动代码。
5. **验收数字化**：每批 DoD 以 `make check-stub` 数字为准——本批覆盖端点桩数归零；禁止以"AI 自评完成"替代。
6. **每批 DoD 全项**：
   - [ ] 本批覆盖端点 `CodeNotImplemented` 桩数 = 0（`make check-stub`）
   - [ ] `make test` 全绿（含认证/交易既有测试回归）
   - [ ] `make lint` 通过
   - [ ] 该批 `specs/0NN-<slug>/tasks.md` 任务全勾
   - [ ] 本表状态列更新为 ✅ 并与收尾代码同 commit，commit 前缀 `feat(0NN-<slug>):`

## 五、变更记录

| 日期 | 批次 | 变更内容 | 原因 | commit |
|---|---|---|---|---|
| 2026-09-21 | — | 建表：13 批规划对账（210 桩 + 8 已连线 = 218），新增 `make check-stub` | 初始规划 | `688c08a` |
| 2026-09-21 | 01 | spec 假设修订：允许新增纯数据种子迁移 000035_admin_seed（超管账号 1 个） | 勘察发现无种子管理员，后台无可登录入口；plan/research D3 | `2b0df43` |
| 2026-09-21 | 01 | api 契约微扩：AdminLogoutReq 增加 refreshToken 可选字段 | FR-004 双凭证同毁需要；对齐 user 渠道 LogoutReq 既有形态 | `51b0b6b` |
| 2026-09-21 | 01 | IRBACLogic 微扩 AdminUserDetail 方法 | api 有账号详情端点而接口漏定义；同 D6 模式 | `1ef20f5` |
| 2026-09-21 | 01 | 评审修复轮（With fixes）：C1 访问日志脱敏（accesslog 既有文件，因本批登录/改密端点引入泄露面）；I1 HasPermission 补角色状态过滤（data-model §三 图纸同步勘误） | 代码评审发现，Critical/Important 合并前必修 | |
