package service

import (
	"github.com/zhoudm1743/go-flow/internal/test/model"
	"github.com/zhoudm1743/go-flow/internal/test/repository"
	"github.com/zhoudm1743/go-flow/internal/test/schemas/req"
	"github.com/zhoudm1743/go-flow/pkg/response"
)

// ProductService 产品服务
type ProductService struct {
	repo *repository.ProductRepository
}

// NewProductService 创建产品服务
func NewProductService(repo *repository.ProductRepository) *ProductService {
	return &ProductService{
		repo: repo,
	}
}

// GetAll 获取所有产品
func (s *ProductService) GetAll() ([]*model.Product, error) {
	return s.repo.FindAll()
}

// GetByID 根据ID获取产品
func (s *ProductService) GetByID(id uint) (*model.Product, error) {
	return s.repo.FindByID(id)
}

// Create 创建产品
func (s *ProductService) Create(req *req.ProductCreateReq) (*model.Product, error) {
	// 将请求转换为模型
	product := &model.Product{
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		Stock:       req.Stock,
		CategoryID:  req.CategoryID,
		Status:      req.Status,
	}

	// 创建记录
	return s.repo.Create(product)
}

// Update 更新产品
func (s *ProductService) Update(id uint, req *req.ProductUpdateReq) (*model.Product, error) {
	// 查找记录
	product, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}

	// 更新字段
	product.Name = req.Name
	product.Description = req.Description
	product.Price = req.Price
	product.Stock = req.Stock
	product.CategoryID = req.CategoryID
	product.Status = req.Status

	// 保存更新
	return s.repo.Update(product)
}

// Delete 删除产品
func (s *ProductService) Delete(id uint) error {
	return s.repo.Delete(id)
}

// GetPage 分页查询产品
func (s *ProductService) GetPage(req *req.PageReq) (*response.PageResult[*model.Product], error) {
	return s.repo.FindPage(req)
}
