// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// ProductSpuDao is the data access object for the table product_spu.
type ProductSpuDao struct {
	table    string             // table is the underlying table name of the DAO.
	group    string             // group is the database configuration group name of the current DAO.
	columns  ProductSpuColumns  // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler // handlers for customized model modification.
}

// ProductSpuColumns defines and stores column names for the table product_spu.
type ProductSpuColumns struct {
	Id                string // SPU ID
	SpuNo             string // SPU业务编码(全局唯一,对外展示用)
	Name              string // 商品名称(SPU级)
	SubTitle          string // 副标题/卖点
	CategoryId        string // 所属分类ID(三级分类)
	BrandId           string // 品牌ID,可为空(无品牌商品)
	SellerId          string // 所属商家(1=自营;多商家预留)
	FreightTemplateId string // 运费模板ID(NULL=包邮)
	Images            string // 主图+轮播图URL数组,按序存储
	Description       string // 图文详情(富文本)
	VideoUrl          string // 主图视频URL(预留)
	SpecDefinitions   string // 规格定义:[{"name":"颜色","values":["黑","白"]}]
	Attributes        string // 非销售属性键值对(材质/产地等)
	SaleRestrictCodes string // 禁售省级区划代码列表(黑名单,NULL=全国可售;下单时地址省代码∈列表即拒绝;列表页过滤由ES承载)
	Status            string // 上架状态:0下架 1上架(默认下架)
	SaleCount         string // 累计销量(异步冗余,支付成功事件累加,排序用)
	PriceMin          string // 启用SKU最低价(冗余,SKU变更时刷新;全部禁用为NULL)
	PriceMax          string // 启用SKU最高价(冗余)
	Deleted           string // 软删除:0否 1是(删除后历史订单不受影响,依赖订单项快照)
	CreatedAt         string // 创建时间
	UpdatedAt         string // 更新时间
}

// productSpuColumns holds the columns for the table product_spu.
var productSpuColumns = ProductSpuColumns{
	Id:                "id",
	SpuNo:             "spu_no",
	Name:              "name",
	SubTitle:          "sub_title",
	CategoryId:        "category_id",
	BrandId:           "brand_id",
	SellerId:          "seller_id",
	FreightTemplateId: "freight_template_id",
	Images:            "images",
	Description:       "description",
	VideoUrl:          "video_url",
	SpecDefinitions:   "spec_definitions",
	Attributes:        "attributes",
	SaleRestrictCodes: "sale_restrict_codes",
	Status:            "status",
	SaleCount:         "sale_count",
	PriceMin:          "price_min",
	PriceMax:          "price_max",
	Deleted:           "deleted",
	CreatedAt:         "created_at",
	UpdatedAt:         "updated_at",
}

// NewProductSpuDao creates and returns a new DAO object for table data access.
func NewProductSpuDao(handlers ...gdb.ModelHandler) *ProductSpuDao {
	return &ProductSpuDao{
		group:    "default",
		table:    "product_spu",
		columns:  productSpuColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *ProductSpuDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *ProductSpuDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *ProductSpuDao) Columns() ProductSpuColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *ProductSpuDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *ProductSpuDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *ProductSpuDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
