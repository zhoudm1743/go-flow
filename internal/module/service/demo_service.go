package service

import (
	"github.com/zhoudm1743/go-frame/internal/module/model"
	"github.com/zhoudm1743/go-frame/internal/module/repository"
	"github.com/zhoudm1743/go-frame/internal/module/schemas/req"
	"github.com/zhoudm1743/go-frame/pkg/response"
)

// DemoService 示例服务
type DemoService struct {
	repo *repository.DemoRepository
}

// NewDemoService 创建示例服务
func NewDemoService(repo *repository.DemoRepository) *DemoService {
	return &DemoService{
		repo: repo,
	}
}

// GetAll 获取所有示例
func (s *DemoService) GetAll() ([]*model.Demo, error) {
	return s.repo.FindAll()
}

// GetByID 根据ID获取示例
func (s *DemoService) GetByID(id uint) (*model.Demo, error) {
	return s.repo.FindByID(id)
}

// Create 创建示例
func (s *DemoService) Create(req *req.DemoCreateReq) (*model.Demo, error) {
	// 将请求转换为模型
	demo := &model.Demo{
		Name:        req.Name,
		Description: req.Description,
		Status:      req.Status,
		// 根据需要添加更多字段
	}

	// 创建记录
	return s.repo.Create(demo)
}

// Update 更新示例
func (s *DemoService) Update(id uint, req *req.DemoUpdateReq) (*model.Demo, error) {
	// 查找记录
	demo, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}

	// 更新字段
	demo.Name = req.Name
	demo.Description = req.Description
	demo.Status = req.Status
	// 根据需要添加更多字段

	// 保存更新
	return s.repo.Update(demo)
}

// Delete 删除示例
func (s *DemoService) Delete(id uint) error {
	return s.repo.Delete(id)
}

// GetPage 分页查询示例
func (s *DemoService) GetPage(req *req.PageReq) (*response.PageResult[*model.Demo], error) {
	return s.repo.FindPage(req)
}
