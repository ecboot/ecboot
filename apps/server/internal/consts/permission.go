package consts

// 权限点全集（模块:资源:动作）——RBAC 权限树种子与路由权限校验共用。
// 清单权威：specs/004-api-surface/data-model.md §三；修订须同步契约与种子迁移。
const (
	PermProductCategoryRead   = "product:category:read"
	PermProductCategoryCreate = "product:category:create"
	PermProductCategoryUpdate = "product:category:update"
	PermProductCategoryDelete = "product:category:delete"
	PermProductBrandRead      = "product:brand:read"
	PermProductBrandCreate    = "product:brand:create"
	PermProductBrandUpdate    = "product:brand:update"
	PermProductBrandDelete    = "product:brand:delete"
	PermProductSpuRead        = "product:spu:read"
	PermProductSpuCreate      = "product:spu:create"
	PermProductSpuUpdate      = "product:spu:update"
	PermProductSpuDelete      = "product:spu:delete"
	PermProductSkuRead        = "product:sku:read"
	PermProductSkuCreate      = "product:sku:create"
	PermProductSkuUpdate      = "product:sku:update"
	PermProductSkuDelete      = "product:sku:delete"

	PermInventoryRead   = "inventory:read"
	PermInventoryAdjust = "inventory:adjust"

	PermOrderRead    = "order:read"
	PermOrderDeliver = "order:deliver"
	PermOrderCancel  = "order:cancel"
	PermOrderUpdate  = "order:update"

	PermAfterSaleRead   = "aftersale:read"
	PermAfterSaleAudit  = "aftersale:audit"
	PermAfterSaleRefund = "aftersale:refund"

	PermPromotionCouponRead       = "promotion:coupon:read"
	PermPromotionCouponCreate     = "promotion:coupon:create"
	PermPromotionCouponUpdate     = "promotion:coupon:update"
	PermPromotionCouponDelete     = "promotion:coupon:delete"
	PermPromotionFullReductionAll = "promotion:fullreduction:manage"
	PermPromotionGroupBuyAll      = "promotion:groupbuy:manage"
	PermPromotionFlashSaleAll     = "promotion:flashsale:manage"
	PermPromotionBargainAll       = "promotion:bargain:manage"
	PermPromotionAssistAll        = "promotion:assist:manage"

	PermOperationBannerRead   = "operation:banner:read"
	PermOperationBannerManage = "operation:banner:manage"
	PermOperationFloorRead    = "operation:floor:read"
	PermOperationFloorManage  = "operation:floor:manage"

	PermStoreManageRead   = "store:manage:read"
	PermStoreManageUpdate = "store:manage:update"
	PermStoreManageCreate = "store:manage:create"
	PermStoreManageDelete = "store:manage:delete"

	PermLogisticsRead   = "logistics:company:read"
	PermLogisticsManage = "logistics:company:manage"

	PermDistributionRead          = "distribution:read"
	PermDistributionAudit         = "distribution:audit"
	PermDistributionRuleRead      = "distribution:rule:read"
	PermDistributionRuleManage    = "distribution:rule:manage"
	PermDistributionWithdrawRead  = "distribution:withdraw:read"
	PermDistributionWithdrawAudit = "distribution:withdraw:audit"
	PermDistributionWithdrawPay   = "distribution:withdraw:pay"

	PermMemberRead   = "member:read"
	PermMemberUpdate = "member:update"

	PermRiskRuleRead     = "risk:rule:read"
	PermRiskRuleManage   = "risk:rule:manage"
	PermRiskRecordRead   = "risk:record:read"
	PermRiskRecordAppeal = "risk:record:appeal"

	PermSystemRoleRead     = "system:role:read"
	PermSystemRoleManage   = "system:role:manage"
	PermSystemRoleAssign   = "system:role:assign"
	PermSystemAdminRead    = "system:admin:read"
	PermSystemAdminManage  = "system:admin:manage"
	PermSystemAdminAssign  = "system:admin:assign"
	PermSystemAuditRead    = "system:audit:read"
	PermSystemConfigRead   = "system:config:read"
	PermSystemConfigUpdate = "system:config:update"

	PermDashboardRead = "dashboard:read"
)
