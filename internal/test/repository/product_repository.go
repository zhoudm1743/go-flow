package repository

import (
	"github.com/zhoudm1743/go-flow/internal/test/model"
	"github.com/zhoudm1743/go-flow/internal/test/schemas/req"
	"github.com/zhoudm1743/go-flow/pkg/response"
	"gorm.io/gorm"
)

// ProductRepository 产品仓库
type ProductRepository struct {
	db *gorm.DB
}

// NewProductRepository 创建产品仓库
func NewProductRepository(db *gorm.DB) *ProductRepository {
	// 自动迁移数据库模型
	_ = db.AutoMigrate(&model.Product{})

	return &ProductRepository{db: db}
}

// FindAll 查询所有产品
func (r *ProductRepository) FindAll() ([]*model.Product, error) {
	var items []*model.Product
	err := r.db.Find(&items).Error
	return items, err
}

// FindByID 根据ID查询产品
func (r *ProductRepository) FindByID(id uint) (*model.Product, error) {
	var item model.Product
	err := r.db.First(&item, id).Error
	return &item, err
}

// Create 创建产品
func (r *ProductRepository) Create(item *model.Product) (*model.Product, error) {
	err := r.db.Create(item).Error
	return item, err
}

// Update 更新产品
func (r *ProductRepository) Update(item *model.Product) (*model.Product, error) {
	err := r.db.Save(item).Error
	return item, err
}

// Delete 删除产品
func (r *ProductRepository) Delete(id uint) error {
	return r.db.Delete(&model.Product{}, id).Error
}

// FindPage 分页查询产品
func (r *ProductRepository) FindPage(req *req.PageReq) (*response.PageResult[*model.Product], error) {
	var items []*model.Product
	var total int64

	// 查询总数
	query := r.db.Model(&model.Product{})
	err := query.Count(&total).Error
	if err != nil {
		return nil, err
	}

	// 查询数据
	err = query.Offset(req.GetOffset()).Limit(req.GetLimit()).Find(&items).Error
	if err != nil {
		return nil, err
	}

	return &response.PageResult[*model.Product]{
		Total: total,
		Items: items,
	}, nil
}
