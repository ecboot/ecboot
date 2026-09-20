package repository

import (
	"ecboot/internal/model"

	"github.com/gogf/gf/v2/database/gdb"
)

// IRepository 是一个通用的泛型接口，用于获取单个实体
// 使用 [T any] 来定义一个泛型类型参数 T
type IRepository[T, D any] interface {
	Model() *gdb.Model
	Create(do *D) error
	CreateBatch(dos []*D) error
	InsertAndGetId(do *D) (int64, error)
	DeleteById(id any) error
	DeleteByIds(ids []any) error
	UpdateById(do *D, id any) error
	GetById(id any) (*T, error)
	GetOne(where interface{}, args ...interface{}) (*T, error)
	ListByIds(ids []any) ([]*T, error)
	List(where interface{}, args ...interface{}) ([]*T, error)
	Count(where interface{}, args ...interface{}) (int, error)
	Page(req model.PageReq, where interface{}, args ...interface{}) (*model.PageResult[T], error)
}

// Repository 是 IRepository 接口的一个泛型实现
type Repository[T, D any] struct {
	model *gdb.Model
}

// New 创建并返回一个 Repository 的泛型实例
func New[T, D any](model *gdb.Model) *Repository[T, D] {
	return &Repository[T, D]{
		model: model,
	}
}

func (s *Repository[T, D]) Model() *gdb.Model {
	return s.model
}

func (s *Repository[T, D]) Create(do *D) error {
	if _, err := s.Model().Insert(do); err != nil {
		return err
	}
	return nil
}

func (s *Repository[T, D]) CreateBatch(dos []*D) error {
	if len(dos) == 0 {
		return nil
	}
	if _, err := s.Model().Insert(dos); err != nil {
		return err
	}
	return nil
}

func (s *Repository[T, D]) InsertAndGetId(do *D) (int64, error) {
	return s.Model().InsertAndGetId(do)
}

func (s *Repository[T, D]) DeleteById(id any) error {
	if _, err := s.Model().WherePri(id).Delete(); err != nil {
		return err
	}
	return nil
}

func (s *Repository[T, D]) DeleteByIds(ids []any) error {
	if len(ids) == 0 {
		return nil
	}
	if _, err := s.Model().WherePri(ids).Delete(); err != nil {
		return err
	}
	return nil
}

func (s *Repository[T, D]) UpdateById(do *D, id any) error {
	if _, err := s.Model().WherePri(id).Update(do); err != nil {
		return err
	}
	return nil
}

func (s *Repository[T, D]) GetById(id any) (entity *T, err error) {
	err = s.Model().WherePri(id).Scan(&entity)
	if err != nil {
		return nil, err
	}
	return entity, nil
}

func (s *Repository[T, D]) GetOne(where interface{}, args ...interface{}) (entity *T, err error) {
	err = s.Model().Where(where, args...).Scan(&entity)
	if err != nil {
		return nil, err
	}
	return entity, nil
}

func (s *Repository[T, D]) ListByIds(ids []any) (entities []*T, err error) {
	if len(ids) == 0 {
		return entities, nil
	}
	err = s.Model().WherePri(ids).Scan(&entities)
	if err != nil {
		return nil, err
	}
	return entities, nil
}

func (s *Repository[T, D]) List(where interface{}, args ...interface{}) (entities []*T, err error) {
	err = s.Model().Where(where, args...).Scan(&entities)
	if err != nil {
		return nil, err
	}
	return entities, nil
}

func (s *Repository[T, D]) Count(where interface{}, args ...interface{}) (int, error) {
	return s.Model().Where(where, args...).Count()
}

// Page 通用分页查询：入参复用契约 PageReq（page/pageSize 单一语义），
// 防御性兜底与契约默认值对齐（默认 10、上限 100）。
func (s *Repository[T, D]) Page(req model.PageReq, where interface{}, args ...interface{}) (result *model.PageResult[T], err error) {
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 10
	}
	if req.PageSize > 100 {
		req.PageSize = 100
	}

	var (
		total   int
		records []T
	)

	total, err = s.Count(where, args...)
	if err != nil {
		return nil, err
	}

	if total == 0 {
		return &model.PageResult[T]{List: records, Total: 0}, nil
	}

	offset := (req.Page - 1) * req.PageSize
	err = s.Model().Where(where, args...).Offset(offset).Limit(req.PageSize).Scan(&records)
	if err != nil {
		return nil, err
	}

	return &model.PageResult[T]{
		List:  records,
		Total: int64(total),
	}, nil
}
