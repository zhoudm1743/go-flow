package user

import (
	"github.com/zhoudm1743/go-flow/internal/user/model"
	"github.com/zhoudm1743/go-flow/internal/user/repository"
	"github.com/zhoudm1743/go-flow/internal/user/schemas/req"
	"github.com/zhoudm1743/go-flow/internal/user/schemas/resp"
	"github.com/zhoudm1743/go-flow/pkg/log"
	"github.com/zhoudm1743/go-flow/pkg/response"
	"gorm.io/gorm"
)

// AccountService 服务接口
type AccountService interface {
	GetAll() ([]resp.AccountResp, error)
	GetByID(id uint) (*resp.AccountResp, error)
	Create(req *req.AccountCreateReq) (*resp.AccountResp, error)
	Update(id uint, req *req.AccountUpdateReq) (*resp.AccountResp, error)
	Delete(id uint) error
	GetPage(pageReq *req.PageReq) (resp.PageResp, error)
}

// accountService 服务实现
type accountService struct {
	repo   repository.AccountRepository
	logger log.Logger
	db     *gorm.DB
}

// NewAccountService 创建服务
func NewAccountService(
	repo repository.AccountRepository,
	logger log.Logger,
	db *gorm.DB,
) AccountService {
	return &accountService{
		repo:   repo,
		logger: logger,
		db:     db,
	}
}

// GetAll 获取所有记录
func (s *accountService) GetAll() ([]resp.AccountResp, error) {
	items, err := s.repo.FindAll()
	if err != nil {
		return nil, err
	}
	var respItems []resp.AccountResp
	response.Copy(&respItems, items)
	return respItems, nil
}

// GetByID 按ID获取
func (s *accountService) GetByID(id uint) (*resp.AccountResp, error) {
	item, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}
	var respItem resp.AccountResp
	response.Copy(&respItem, item)
	return &respItem, nil
}

// Create 创建记录
func (s *accountService) Create(req *req.AccountCreateReq) (*resp.AccountResp, error) {
	// 创建模型
	var item model.Account
	response.Copy(req, &item)
	// 如果状态为空，设置默认值
	if item.Status == 0 {
		item.Status = 1 // 默认为启用状态
	}

	s.logger.Debugf("创建账户: %s", item.Username)
	res, err := s.repo.Create(&item)
	if err != nil {
		return nil, err
	}
	var respItem resp.AccountResp
	response.Copy(&respItem, res)
	return &respItem, nil
}

// Update 更新记录
func (s *accountService) Update(id uint, req *req.AccountUpdateReq) (*resp.AccountResp, error) {
	// 获取现有记录
	item, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}

	response.Copy(req, item)

	s.logger.Debugf("更新账户: ID=%d", id)
	res, err := s.repo.Update(item)
	if err != nil {
		return nil, err
	}
	var respItem resp.AccountResp
	response.Copy(&respItem, res)
	return &respItem, nil
}

// Delete 删除记录
func (s *accountService) Delete(id uint) error {
	return s.repo.Delete(id)
}

// GetPage 分页查询
func (s *accountService) GetPage(pageReq *req.PageReq) (resp.PageResp, error) {
	// 构建查询条件
	query := map[string]interface{}{}

	// 执行分页查询
	items, err := s.repo.Paginate(pageReq.PageNo, pageReq.PageSize, query)
	if err != nil {
		return resp.PageResp{}, err
	}

	// 获取总数
	var total int64
	s.db.Model(&model.Account{}).Count(&total)

	var respItems []resp.AccountResp
	response.Copy(items, &respItems)

	// 返回分页结果
	return resp.PageResp{
		PageNo:   pageReq.PageNo,
		PageSize: pageReq.PageSize,
		Total:    total,
		Data:     respItems,
	}, nil
}
