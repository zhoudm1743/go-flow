package service

import (
	"github.com/zhoudm1743/go-flow/internal/shop/model"
	"github.com/zhoudm1743/go-flow/internal/shop/repository"
	"github.com/zhoudm1743/go-flow/internal/shop/schemas/req"
	"github.com/zhoudm1743/go-flow/pkg/response"
)

// CategoryService 分类服务
type CategoryService struct {
	repo *repository.CategoryRepository
}

// NewCategoryService 创建分类服务
func NewCategoryService(repo *repository.CategoryRepository) *CategoryService {
	return &CategoryService{
		repo: repo,
	}
}

// GetAll 获取所有分类
func (s *CategoryService) GetAll() ([]*model.Category, error) {
	return s.repo.FindAll()
}

// GetByID 根据ID获取分类
func (s *CategoryService) GetByID(id uint) (*model.Category, error) {
	return s.repo.FindByID(id)
}

// Create 创建分类
func (s *CategoryService) Create(req *req.CategoryCreateReq) (*model.Category, error) {
	// 将请求转换为模型
	category := &model.Category{
		Name:        req.Name,
		Description: req.Description,
		ParentID:    req.ParentID,
		Level:       req.Level,
		Sort:        req.Sort,
		Status:      req.Status,
	}

	// 创建记录
	return s.repo.Create(category)
}

// Update 更新分类
func (s *CategoryService) Update(id uint, req *req.CategoryUpdateReq) (*model.Category, error) {
	// 查找记录
	category, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}

	// 更新字段
	category.Name = req.Name
	category.Description = req.Description
	category.ParentID = req.ParentID
	category.Level = req.Level
	category.Sort = req.Sort
	category.Status = req.Status

	// 保存更新
	return s.repo.Update(category)
}

// Delete 删除分类
func (s *CategoryService) Delete(id uint) error {
	return s.repo.Delete(id)
}

// GetPage 分页查询分类
func (s *CategoryService) GetPage(req *req.PageReq) (*response.PageResult[*model.Category], error) {
	return s.repo.FindPage(req)
}
