package repository

import (
	"errors"

	"github.com/zhoudm1743/go-flow/internal/user/model"
	"github.com/zhoudm1743/go-flow/pkg/log"
	"gorm.io/gorm"
)

// AccountRepository 用户仓库接口
type AccountRepository interface {
	FindAll() ([]model.Account, error)
	FindByID(id uint) (*model.Account, error)
	Find(query map[string]interface{}) ([]model.Account, error)
	Paginate(page, pageSize int, query map[string]interface{}) ([]model.Account, error)
	Create(user *model.Account) (*model.Account, error)
	Update(user *model.Account) (*model.Account, error)
	Delete(id uint) error
}

// accountRepository 仓库实现
type accountRepository struct {
	db     *gorm.DB
	logger log.Logger
}

// NewAccountRepository 创建仓库
func NewAccountRepository(db *gorm.DB, logger log.Logger) AccountRepository {
	db.AutoMigrate(&model.Account{})
	return &accountRepository{
		db:     db,
		logger: logger,
	}
}

// Find 按条件查找
func (r *accountRepository) Find(query map[string]interface{}) ([]model.Account, error) {
	chain := r.db.Model(&model.Account{})
	for key, value := range query {
		chain = chain.Where(key, value)
	}
	var items []model.Account
	if err := chain.Find(&items).Error; err != nil {
		r.logger.Errorf("查询account失败: %v", err)
		return nil, err
	}
	return items, nil
}

// Paginate 分页查询
func (r *accountRepository) Paginate(page int, pageSize int, query map[string]interface{}) ([]model.Account, error) {
	chain := r.db.Model(&model.Account{})
	for key, value := range query {
		chain = chain.Where(key, value)
	}
	var items []model.Account
	if err := chain.Offset((page - 1) * pageSize).Limit(pageSize).Find(&items).Error; err != nil {
		r.logger.Errorf("分页查询account失败: %v", err)
		return nil, err
	}
	return items, nil
}

// FindAll 查询所有记录
func (r *accountRepository) FindAll() ([]model.Account, error) {
	var items []model.Account
	if err := r.db.Find(&items).Error; err != nil {
		r.logger.Errorf("查询所有account失败: %v", err)
		return nil, err
	}
	return items, nil
}

// FindByID 按ID查找
func (r *accountRepository) FindByID(id uint) (*model.Account, error) {
	var item model.Account
	if err := r.db.First(&item, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			r.logger.Warnf("account不存在 ID=%d", id)
			return nil, errors.New("account不存在")
		}
		r.logger.Errorf("查询account失败 ID=%d: %v", id, err)
		return nil, err
	}
	return &item, nil
}

// Create 创建记录
func (r *accountRepository) Create(item *model.Account) (*model.Account, error) {
	if err := r.db.Create(item).Error; err != nil {
		r.logger.Errorf("创建account失败: %v", err)
		return nil, err
	}
	return item, nil
}

// Update 更新记录
func (r *accountRepository) Update(item *model.Account) (*model.Account, error) {
	if err := r.db.Save(item).Error; err != nil {
		r.logger.Errorf("更新account失败 ID=%d: %v", item.ID, err)
		return nil, err
	}
	return item, nil
}

// Delete 删除记录
func (r *accountRepository) Delete(id uint) error {
	if err := r.db.Delete(&model.Account{}, id).Error; err != nil {
		r.logger.Errorf("删除account失败 ID=%d: %v", id, err)
		return err
	}
	return nil
}
