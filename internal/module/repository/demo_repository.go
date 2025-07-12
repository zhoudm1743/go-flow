package repository

import (
	"github.com/zhoudm1743/go-frame/internal/module/model"
	"github.com/zhoudm1743/go-frame/internal/module/schemas/req"
	"github.com/zhoudm1743/go-frame/pkg/response"
	"gorm.io/gorm"
)

// DemoRepository 示例仓库
type DemoRepository struct {
	db *gorm.DB
}

// NewDemoRepository 创建示例仓库
func NewDemoRepository(db *gorm.DB) *DemoRepository {
	// 自动迁移数据库模型
	_ = db.AutoMigrate(&model.Demo{})
	
	return &DemoRepository{
		db: db,
	}
}

// FindAll 查询所有示例
func (r *DemoRepository) FindAll() ([]*model.Demo, error) {
	var items []*model.Demo
	err := r.db.Find(&items).Error
	return items, err
}

// FindByID 根据ID查询示例
func (r *DemoRepository) FindByID(id uint) (*model.Demo, error) {
	var item model.Demo
	err := r.db.First(&item, id).Error
	return &item, err
}

// Create 创建示例
func (r *DemoRepository) Create(item *model.Demo) (*model.Demo, error) {
	err := r.db.Create(item).Error
	return item, err
}

// Update 更新示例
func (r *DemoRepository) Update(item *model.Demo) (*model.Demo, error) {
	err := r.db.Save(item).Error
	return item, err
}

// Delete 删除示例
func (r *DemoRepository) Delete(id uint) error {
	return r.db.Delete(&model.Demo{}, id).Error
}

// FindPage 分页查询示例
func (r *DemoRepository) FindPage(req *req.PageReq) (*response.PageResult[*model.Demo], error) {
	var items []*model.Demo
	var total int64

	// 查询总数
	query := r.db.Model(&model.Demo{})
	err := query.Count(&total).Error
	if err != nil {
		return nil, err
	}

	// 查询数据
	err = query.Offset(req.GetOffset()).Limit(req.GetLimit()).Find(&items).Error
	if err != nil {
		return nil, err
	}

	return &response.PageResult[*model.Demo]{
		Total: total,
		Items: items,
	}, nil
}
