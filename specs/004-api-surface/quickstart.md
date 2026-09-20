# Quickstart: 全角色 API 接口面验证指南

本特性交付契约层（接口定义与路由注册），验证 = 结构正确 + 路由完备 + 旅程走查。

## 前提

```bash
cd apps/server && docker compose up -d && make migrate-up
```

## 验证一：契约层编译（宪法 IV）

```bash
make build    # go build ./... 零错误 → api 定义结构合法
make test     # go test ./...
```

## 验证二：路由注册自检

启动应用（`make run`），访问路由调试输出（GoFrame 启动日志打印路由表）核对：

- 四前缀分组存在：`/common`、`/user`、`/shop`、`/admin`
- 路由总数与 contracts 四文件端点行数一致（约 140；±条件注册的 mock 端点）
- mock 取码路由在生产配置（mock 关闭）下**不存在**

## 验证三：三旅程契约走查（对照 contracts 文件）

以 curl/接口工具按序调用（业务实现未落地时以"路由可达 + 统一响应结构"为通过标准；实现落地后升级为语义断言）：

1. **游客逛**：`/common/stores` → `/shop/categories` → `/shop/products` → `/shop/products/{id}` → `/shop/activities/flash-sales` → `/shop/banners`——全部公开可达、响应三段式。
2. **会员买**（需 003 登录链路可用）：登录取 token → `/shop/cart/items` → `/shop/cart/checkout` → `/shop/orders` → `/shop/pay` → `/shop/orders/{no}` → `/shop/after-sales`——每步鉴权生效（无 token 返回 10003）。
3. **管理员管**：`/admin/login` → `/admin/products`（POST）→ `/admin/products/{id}/status` → `/admin/orders` → `/admin/withdraws/{no}/audit`——无权限点账号访问返回 10005。

## 验证四：契约一致性抽检

- 同一业务状态 C 端与管理端约束一致：会员取消已完成订单（40006）= 管理员对已完成订单发货（40006）。
- 金额/ID/时间类型抽检：任一列表响应中金额为 string、ID 为 string、时间为 RFC3339。
- 脱敏抽检：`/admin/members` 响应手机号含 `****`。
