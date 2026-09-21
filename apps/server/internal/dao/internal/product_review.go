// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// ProductReviewDao is the data access object for the table product_review.
type ProductReviewDao struct {
	table    string               // table is the underlying table name of the DAO.
	group    string               // group is the database configuration group name of the current DAO.
	columns  ProductReviewColumns // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler   // handlers for customized model modification.
}

// ProductReviewColumns defines and stores column names for the table product_review.
type ProductReviewColumns struct {
	Id           string // 评价ID
	OrderItemId  string // 订单项ID(一项一评)
	OrderNo      string // 订单号(冗余,免联查)
	UserId       string // 评价用户ID
	SpuId        string // SPU ID(聚合统计用)
	SkuId        string // SKU ID
	SpuName      string // SPU名称(下单时快照,商品软删后自洽展示)
	SkuSpecs     string // 规格组合(下单时快照)
	Score        string // 评分:1-5星
	Content      string // 评价内容
	Images       string // 评价图片URL数组
	IsAnonymous  string // 匿名:0否 1是
	AuditStatus  string // 审核状态:0待审核 1通过 2驳回
	ExtraContent string // 追评内容(仅一次,90天内)
	ExtraTime    string // 追评时间
	ReplyContent string // 商家回复(仅一次)
	ReplyTime    string // 回复时间
	Deleted      string // 软删除:0否 1是(用户删除自己的评价)
	CreatedAt    string // 创建时间
	UpdatedAt    string // 更新时间
}

// productReviewColumns holds the columns for the table product_review.
var productReviewColumns = ProductReviewColumns{
	Id:           "id",
	OrderItemId:  "order_item_id",
	OrderNo:      "order_no",
	UserId:       "user_id",
	SpuId:        "spu_id",
	SkuId:        "sku_id",
	SpuName:      "spu_name",
	SkuSpecs:     "sku_specs",
	Score:        "score",
	Content:      "content",
	Images:       "images",
	IsAnonymous:  "is_anonymous",
	AuditStatus:  "audit_status",
	ExtraContent: "extra_content",
	ExtraTime:    "extra_time",
	ReplyContent: "reply_content",
	ReplyTime:    "reply_time",
	Deleted:      "deleted",
	CreatedAt:    "created_at",
	UpdatedAt:    "updated_at",
}

// NewProductReviewDao creates and returns a new DAO object for table data access.
func NewProductReviewDao(handlers ...gdb.ModelHandler) *ProductReviewDao {
	return &ProductReviewDao{
		group:    "default",
		table:    "product_review",
		columns:  productReviewColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *ProductReviewDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *ProductReviewDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *ProductReviewDao) Columns() ProductReviewColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *ProductReviewDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *ProductReviewDao) Ctx(ctx context.Context) *gdb.Model {
	model := dao.DB().Model(dao.table)
	for _, handler := range dao.handlers {
		model = handler(model)
	}
	return model.Safe().Ctx(ctx)
}

// Transaction wraps the transaction logic using function f.
// It rolls back the transaction and returns the error if function f returns a non-nil error.
// It commits the transaction and returns nil if function f returns nil.
//
// Note: Do not commit or roll back the transaction in function f,
// as it is automatically handled by this function.
func (dao *ProductReviewDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
