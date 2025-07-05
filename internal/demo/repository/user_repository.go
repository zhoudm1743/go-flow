package repository

import (
	"errors"

	"github.com/zhoudm1743/go-flow/internal/demo/model"
	"github.com/zhoudm1743/go-flow/pkg/log"
	"gorm.io/gorm"
)

// UserRepository 用户仓库接口
type UserRepository interface {
	FindAll() ([]model.User, error)
	FindByID(id uint) (*model.User, error)
	Find(query map[string]interface{}) ([]model.User, error)
	Paginate(page, pageSize int, query map[string]interface{}) ([]model.User, error)
	Create(user *model.User) (*model.User, error)
	Update(user *model.User) (*model.User, error)
	Delete(id uint) error
}

// userRepository 用户仓库实现
type userRepository struct {
	db     *gorm.DB
	logger log.Logger
}

// Find implements UserRepository.
func (r *userRepository) Find(query map[string]interface{}) ([]model.User, error) {
	chain := r.db.Model(&model.User{})
	for key, value := range query {
		chain = chain.Where(key, value)
	}
	var users []model.User
	if err := chain.Find(&users).Error; err != nil {
		r.logger.Errorf("查询用户失败: %v", err)
		return nil, err
	}
	return users, nil
}

// Paginate implements UserRepository.
func (r *userRepository) Paginate(page int, pageSize int, query map[string]interface{}) ([]model.User, error) {
	chain := r.db.Model(&model.User{})
	for key, value := range query {
		chain = chain.Where(key, value)
	}
	var users []model.User
	if err := chain.Offset((page - 1) * pageSize).Limit(pageSize).Find(&users).Error; err != nil {
		r.logger.Errorf("分页查询用户失败: %v", err)
		return nil, err
	}
	return users, nil
}

// NewUserRepository 创建用户仓库
func NewUserRepository(db *gorm.DB, logger log.Logger) UserRepository {
	db.AutoMigrate(&model.User{})
	return &userRepository{
		db:     db,
		logger: logger,
	}
}

// FindAll 查找所有用户
func (r *userRepository) FindAll() ([]model.User, error) {
	var users []model.User
	if err := r.db.Find(&users).Error; err != nil {
		r.logger.Errorf("查询所有用户失败: %v", err)
		return nil, err
	}
	return users, nil
}

// FindByID 根据ID查找用户
func (r *userRepository) FindByID(id uint) (*model.User, error) {
	var user model.User
	if err := r.db.First(&user, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			r.logger.Warnf("用户不存在 ID=%d", id)
			return nil, errors.New("用户不存在")
		}
		r.logger.Errorf("查询用户失败 ID=%d: %v", id, err)
		return nil, err
	}
	return &user, nil
}

// Create 创建用户
func (r *userRepository) Create(user *model.User) (*model.User, error) {
	if err := r.db.Create(user).Error; err != nil {
		r.logger.Errorf("创建用户失败: %v", err)
		return nil, err
	}
	return user, nil
}

// Update 更新用户
func (r *userRepository) Update(user *model.User) (*model.User, error) {
	if err := r.db.Save(user).Error; err != nil {
		r.logger.Errorf("更新用户失败 ID=%d: %v", user.ID, err)
		return nil, err
	}
	return user, nil
}

// Delete 删除用户
func (r *userRepository) Delete(id uint) error {
	if err := r.db.Delete(&model.User{}, id).Error; err != nil {
		r.logger.Errorf("删除用户失败 ID=%d: %v", id, err)
		return err
	}
	return nil
}
