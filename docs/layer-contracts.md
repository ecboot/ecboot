# 分层契约（Go 后端）

后端 `apps/server` 采用 GoFrame 惯用单体结构。本文件是**分层方向、数据流与代码生成**的权威约定（评审依据），修订须同步 §六 的触发器记录。

## 一、分层方向

```text
main.go / internal/app（装配与启动，bootstrap 注册驱动）
  → api/{common,user,shop,admin}/v1        接口定义（req/res 结构体，四渠道目录隔离）
  → internal/controller/{common,user,shop,admin}   控制器（参数绑定 → 只调 service 和 manager → 响应组装）
  → internal/service/{user,shop,system}    业务逻辑（接口定义 + 实现; 领域兄弟互禁, 经领域事件协作）
  → internal/manager                        第三方服务/通用管理器（非 dao 数据源适配）
  → internal/dao（gf gen dao 生成）→ internal/model/entity + do（gf gen dao 生成）
横切：internal/middleware（认证/统一响应/恢复/TraceId）、internal/library（技术组件：验证码/短信/加密/会话）、
      internal/consts（常量/权限点）、internal/errcode（错误码）、internal/model（公共 DTO，见 §三）
```

## 二、数据流与强类型约定（禁止 map 传参）

| 层 | 查询条件（Where） | 写入数据（Data） | 查询返回 | 说明 |
|---|---|---|---|---|
| controller | v1 Req 结构体 | — | v1 Res 结构体 | 只组装，不写业务 |
| service | **`*do` 对象 / `model.PageReq`** | **`*do` 对象** | **`entity` 结构体** → 转换为 `model` DTO 返回 | **禁止 `g.Map`/`g.List` 传参** |
| dao（生成物） | entity | do | entity | 勿手改 |

- **查询返回一律 entity 接收**（gf Scan 到 `model/entity` 结构体），service 层负责 entity → model DTO 的转换。
- **Data/Where 参数一律 do 对象**：`dao.User.Ctx(ctx).Data(do.User{...}).Where(dao.User.Columns().PhoneHash, hash)`；禁止 `g.Map{...}` 与裸 SQL 传参（聚合统计等确需 Raw 的，封装为 dao 层具名方法）。
- **强类型替代弱类型的例外**：无（动态列场景应建表或建视图，不做 g.Map 逃生门）。

## 三、DTO 归属（internal/model 统一）

- **`internal/model` 是全部公共 DTO 的唯一归属**：
  - `model.go` — 契约分页结构（`PageReq`/`PageRes`）与内部容器（`PageResult[T]`）；
  - `dto_*.go`（按域分文件）— service 层出入参 DTO（**类型名带域前缀**：`UserAddressItem`/`AdminRoleItem`，避免平铺冲突）；
  - `entity/` 与 `do/` — gf gen dao 生成物（表映射），禁止手改。
- **api 层不再维护本地别名**（原 `api/base` 已删除）：`api/{渠道}/v1` 直接 `model.PageReq` / `model.PageRes` 嵌入与引用。
- **service 中的 struct 一律放 model**：接口签名出入参、跨层 DTO 全部 `model.` 前缀；service 包内只保留接口（`IXxxLogic`）与实现。
- entity（表映射）与 model DTO 的**转换发生在 service 层**（entity 是存储形态、DTO 是接口形态，不共享结构）。

## 四、调用关系约定

| 层 | 允许调用 | 禁止 |
|---|---|---|
| controller | **service（接口）、manager** | dao、entity/do 直查、其他 controller、跨渠道 |
| service | **dao、manager**、其他 service 仅经领域事件/接口 | controller、跨渠道、裸 g.DB() map 传参 |
| manager | dao、library、外部 SDK | service、controller |

- **controller 原则上只调用 service 和 manager**：参数绑定 → 调用 → 组装 v1 Res；禁止出现 SQL、业务规则、跨渠道调用。
- **service 仅调用 dao 和 manager**：领域逻辑与事务编排在 service；禁止绕过 dao 直查（聚合统计封装为 dao 具名方法）。
- **manager** 职责：第三方平台适配（短信/微信/OSS/支付渠道）、跨 dao 的通用数据装配——与 library（纯技术组件，无数据语义）区分。

## 五、测试约定（红绿循环）

- **service 方法一律 TDD**：先写失败测试（红）→ 实现 → 转绿；测试文件与被测包同目录（`*_test.go`）。
- 测试配置**确定性注入**（gdb.SetConfig/gredis.SetConfig 指向 compose 基线），禁止依赖 manifest/config 的本地差异。
- 验收断言到行为级（错误码/状态迁移/恒等式），不以"函数被调用"为通过标准。

## 六、生成物治理

- **`gf gen ctrl`**：`internal/controller/*` 与 `api/{渠道}/{渠道}.go` 聚合接口——修改 v1 后运行同步；`*_new.go` 可手工调整，`*_v1_*.go` 桩由业务特性填充。
- **`gf gen dao`**：`internal/dao`、`internal/model/entity`、`internal/model/do`——表结构变更后运行；**禁止手改**。
- `internal/model/model.go` 与 `dto_*.go` 为手工维护公共模型，不属于生成范围。

## 七、包组织与拆分触发器（service 层）

- 现行三领域包：`service/{user,shop,system}`，包内平铺文件。
- 命名纪律：**包内导出类型必须带子域前缀**（`DistXxx`/`PointXxx`），禁止裸通用名。
- 拆分触发器（到点机械执行）：单包文件 > 15 或 > 5k 行；某子域需独立测试装配/独立依赖；user 分销落地文件 > 5 → 拆 `service/distribution`。

## 系统定位

纯 B2C + 多门店社交电商：平台统一经营商品与交易，门店（store）为线下载体（自提/核销/附近门店）。无多租户、无多商户/商家概念；B2B2C 需求另立独立项目。
