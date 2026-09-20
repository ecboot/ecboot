# 分层契约（Go 后端）

后端 `apps/server` 采用 GoFrame 惯用单体结构，分层方向与职责如下（描述性契约，评审依据）：

## 分层方向

```
main.go / internal/app（装配与启动，bootstrap 初始化）
  → api/{common,user,shop,admin}/v1   接口定义（req/res 结构体，四渠道目录隔离）
  → internal/controller/{common,user,shop,admin}   控制器（参数绑定 → service 调用 → 响应组装）
  → internal/service/*                业务接口与实现（领域逻辑，user/shop 兄弟域经事件协作）
  → internal/repository → internal/dao|model（gf gen dao 生成物，勿手改）
横切：internal/middleware（认证/统一响应/恢复）、internal/library（技术组件：验证码/短信/加密/会话）
```

## 规则

- **渠道隔离**：`api/{common,user,shop,admin}` 四目录互不引用；公共渠道（common）不依赖任何领域业务。
- **领域兄弟互禁**：`service/user` 与 `service/shop` 互不 import，经领域事件协作。
- **生成物治理**：
  - **控制器约定**：`internal/controller/*` 与 `api/{渠道}/{渠道}.go` 聚合接口由 **`gf gen ctrl`** 生成——修改 `api/{渠道}/v1` 接口定义后运行该命令同步；`*_new.go` 允许手工调整，`*_v1_*.go` 桩由业务特性填充实现。
  - **数据访问约定**：`internal/dao`、`internal/model/entity`、`internal/model/do` 由 **`gf gen dao`** 更新（配置 `hack/config.yaml`）——表结构变更后运行；生成物禁止手改。
  - `internal/model/model.go` 为手工维护的公共模型（含 api 契约分页结构 PageReq/PageRes），不属于 dao 生成范围。
- 依赖方向若需机器强制（如 depguard），作为演进项在 `.golangci.yml` 落地并回填本文件。

## 系统定位

纯 B2C + 多门店社交电商：平台统一经营商品与交易，门店（store）为线下载体（自提/核销/附近门店）。无多租户、无多商户/商家概念；B2B2C 需求另立独立项目。
