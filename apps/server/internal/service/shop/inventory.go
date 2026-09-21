// inventory.go 库存域——表: inventory / inventory_log（V3）。
// 规则(ADR-0001): 三段式 total/locked, 可售=total-locked 推导;
// 全部变更为单行原子条件 UPDATE + 同事务 inventory_log（前后值快照, 对账唯一依据）;
// 热点 SKU 可前置 Redis 挡板（DB 仍唯一事实源）; 调整留痕 operator。
package shop

import "context"

// IInventoryLogic 库存。
type IInventoryLogic interface {
	// List 库存列表（可售推导）。
	List(ctx context.Context, skuId int64, keyword string, page PageQuery) (*PageResult[InventoryItem], error)
	// Warnings 预警列表（可售≤warn_count）。
	Warnings(ctx context.Context, page PageQuery) (*PageResult[InventoryItem], error)
	// Adjust 后台调整（delta 有符号; 变更前后快照写 log; operator=后台账号）。
	Adjust(ctx context.Context, skuId int64, delta int, operator string, remark string) (int, error)
	// Lock 下单锁定（核心: UPDATE ... SET locked=locked+n WHERE total-locked>=n; affected=0 → 库存不足 40001）。
	// items 需在同一事务（调用方编排订单创建）。
	Lock(ctx context.Context, tx any, skuId int64, qty int, orderNo string) error
	// Deduct 支付核销（total-n, locked-n）。
	Deduct(ctx context.Context, tx any, skuId int64, qty int, orderNo string) error
	// Release 取消/超时释放（locked-n）。
	Release(ctx context.Context, tx any, skuId int64, qty int, orderNo string) error
	// Restock 售后回补（total+n, bizType=6）。
	Restock(ctx context.Context, skuId int64, qty int, orderNo string) error
}

type InventoryItem struct {
	SkuId     int64  `json:"skuId"`
	SkuNo     string `json:"skuNo"`
	SkuName   string `json:"skuName"`
	Total     int    `json:"total"`
	Locked    int    `json:"locked"`
	Available int    `json:"available" dc:"total-locked"`
	WarnCount int    `json:"warnCount"`
}
