// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// ProductSkuDao is the data access object for the table product_sku.
type ProductSkuDao struct {
	table    string             // table is the underlying table name of the DAO.
	group    string             // group is the database configuration group name of the current DAO.
	columns  ProductSkuColumns  // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler // handlers for customized model modification.
}

// ProductSkuColumns defines and stores column names for the table product_sku.
type ProductSkuColumns struct {
	Id        string // SKU ID
	SkuNo     string // SKU业务编码(全局唯一)
	SpuId     string // 所属SPU ID
	Name      string // SKU名称=SPU名+规格串(冗余生成,展示用)
	Specs     string // 规格组合快照:{"颜色":"黑","尺码":"M"}(同SPU内组合唯一由应用层保证)
	SpecsHash string // 规格组合哈希(同SPU内组合唯一;应用序列化稳定时生效)
	Price     string // 现售价(元)
	LinePrice string // 划线价/原价(营销展示,可为空)
	CostPrice string // 成本价(后台可见,毛利分析预留)
	Image     string // SKU图URL(空则用SPU主图)
	Weight    string // 重量(克,运费按重计费时使用)
	Barcode   string // 条形码(预留)
	Sort      string // 展示排序
	Status    string // 状态:1启用 0禁用(断码/失效单品)
	Deleted   string // 软删除:0否 1是
	CreatedAt string // 创建时间
	UpdatedAt string // 更新时间
}

// productSkuColumns holds the columns for the table product_sku.
var productSkuColumns = ProductSkuColumns{
	Id:        "id",
	SkuNo:     "sku_no",
	SpuId:     "spu_id",
	Name:      "name",
	Specs:     "specs",
	SpecsHash: "specs_hash",
	Price:     "price",
	LinePrice: "line_price",
	CostPrice: "cost_price",
	Image:     "image",
	Weight:    "weight",
	Barcode:   "barcode",
	Sort:      "sort",
	Status:    "status",
	Deleted:   "deleted",
	CreatedAt: "created_at",
	UpdatedAt: "updated_at",
}

// NewProductSkuDao creates and returns a new DAO object for table data access.
func NewProductSkuDao(handlers ...gdb.ModelHandler) *ProductSkuDao {
	return &ProductSkuDao{
		group:    "default",
		table:    "product_sku",
		columns:  productSkuColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *ProductSkuDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *ProductSkuDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *ProductSkuDao) Columns() ProductSkuColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *ProductSkuDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *ProductSkuDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *ProductSkuDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
