// Package errcode 业务错误码登记（契约总则见 specs/004-api-surface/contracts/conventions.md）。
// 分段：1xxxx 通用；2xxxx 用户域；3xxxx 商品域；4xxxx 交易域；5xxxx 促销域；
// 6xxxx 分销域；7xxxx 门店域；8xxxx 后台治理域。
package errcode

import (
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
)

// 通用（1xxxx）
const (
	CodeOK           = 0
	CodeInvalidParam = 10001 // 参数校验失败
	CodeSystemError  = 10002 // 系统错误（附追踪号）
	CodeUnauthorized = 10003 // 未登录/凭证失效
	CodeTooFrequent  = 10004 // 操作过于频繁
	CodeForbidden    = 10005 // 无权限
)

// 用户域（2xxxx）
const (
	CodeCaptchaError   = 20001 // 验证码错误或已过期
	CodeLocked         = 20002 // 尝试次数超限已锁定
	CodeUserDisabled   = 20003 // 账号已禁用
	CodeWxBindConflict = 20004 // 微信身份绑定冲突
)

// 商品域（3xxxx）
const (
	CodeProductNotFound = 30001 // 商品不存在/已下架
	CodeSkuUnsellable   = 30002 // SKU 不可售
	CodeProductInvalid  = 30003 // 商品失效（结算项）
	CodeCategoryInUse   = 30004 // 分类有商品引用
	CodeNoEnableSku     = 30005 // 无启用 SKU 禁止上架
	CodeSpecDup         = 30006 // 规格组合重复
	CodeInventoryAdjust = 30007 // 库存调整非法
)

// 交易域（4xxxx）
const (
	CodeStockInsufficient = 40001 // 库存不足
	CodeCouponUnusable    = 40002 // 券不可用
	CodeSoldOut           = 40003 // 已抢完
	CodeBargainUnpayable  = 40004 // 砍价单不可下单
	CodeOrderNotFound     = 40005 // 订单不存在
	CodeStatusNotAllowed  = 40006 // 状态不允许该操作
	CodePayCreateFailed   = 40007 // 支付单创建失败
	CodeAfterSaleDenied   = 40008 // 不可售后
	CodeAfterSaleNotFound = 40009 // 售后单不存在
	CodeAlreadyReviewed   = 40010 // 已评价
	CodeExtraReviewDenied = 40011 // 追评违规（已追评/超期）
)

// 促销域（5xxxx）
const (
	CodeCouponSoldOut    = 50001 // 券已领完/超限
	CodeActivityInvalid  = 50002 // 活动无效
	CodeActivityNotFound = 50003 // 活动/参与单不存在
	CodeAlreadyCut       = 50004 // 已帮砍过
	CodeFloorReached     = 50005 // 已到底价
	CodeAssistUsedUp     = 50006 // 发起次数用尽
	CodeAlreadyAssisted  = 50007 // 已助力
	CodeLadderDup        = 50008 // 档位门槛重复
)

// 分销域（6xxxx）
const (
	CodeDistAlreadyApplied  = 60001 // 已申请/已是推广员
	CodeBalanceInsufficient = 60002 // 余额不足
	CodeWithdrawOngoing     = 60003 // 存在进行中的提现单
)

// 门店域（7xxxx）
const (
	CodeStoreUnavailable = 70001 // 门店不存在/已歇业
)

// 后台治理域（8xxxx）
const (
	CodeAdminBadCredential = 80001 // 后台凭证错误
	CodeOldPasswordWrong   = 80002 // 原密码错误
)

// New 构造带契约错误码的业务错误（统一响应中间件映射为三段式）。
func New(code int, msg string) error {
	return gerror.NewCode(gcode.New(code, msg, nil))
}

