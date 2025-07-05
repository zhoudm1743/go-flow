package repository

import (
	"github.com/zhoudm1743/go-flow/internal/shop/model"
	"github.com/zhoudm1743/go-flow/internal/shop/schemas/req"
	"github.com/zhoudm1743/go-flow/pkg/response"
	"gorm.io/gorm"
)

// CategoryRepository 分类仓库
type CategoryRepository struct {
	db *gorm.DB
}

// NewCategoryRepository 创建分类仓库
func NewCategoryRepository(db *gorm.DB) *CategoryRepository {
	// 自动迁移数据库模型
	_ = db.AutoMigrate(&model.Category{})

	return &CategoryRepository{
		db: db,
	}
}

// FindAll 查询所有分类
func (r *CategoryRepository) FindAll() ([]*model.Category, error) {
	var items []*model.Category
	err := r.db.Find(&items).Error
	return items, err
}

// FindByID 根据ID查询分类
func (r *CategoryRepository) FindByID(id uint) (*model.Category, error) {
	var item model.Category
	err := r.db.First(&item, id).Error
	return &item, err
}

// Create 创建分类
func (r *CategoryRepository) Create(item *model.Category) (*model.Category, error) {
	err := r.db.Create(item).Error
	return item, err
}

// Update 更新分类
func (r *CategoryRepository) Update(item *model.Category) (*model.Category, error) {
	err := r.db.Save(item).Error
	return item, err
}

// Delete 删除分类
func (r *CategoryRepository) Delete(id uint) error {
	return r.db.Delete(&model.Category{}, id).Error
}

// FindPage 分页查询分类
func (r *CategoryRepository) FindPage(req *req.PageReq) (*response.PageResult[*model.Category], error) {
	var items []*model.Category
	var total int64

	// 查询总数
	query := r.db.Model(&model.Category{})
	err := query.Count(&total).Error
	if err != nil {
		return nil, err
	}

	// 查询数据
	err = query.Offset(req.GetOffset()).Limit(req.GetLimit()).Find(&items).Error
	if err != nil {
		return nil, err
	}

	return &response.PageResult[*model.Category]{
		Total: total,
		Items: items,
	}, nil
}
