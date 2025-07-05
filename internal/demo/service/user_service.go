package service

import (
	"github.com/zhoudm1743/go-flow/internal/demo/model"
	"github.com/zhoudm1743/go-flow/internal/demo/repository"
	"github.com/zhoudm1743/go-flow/pkg/log"
)

// UserService 用户服务接口
type UserService interface {
	GetAllUsers() ([]model.User, error)
	GetUserByID(id uint) (*model.User, error)
	CreateUser(name, email string) (*model.User, error)
	UpdateUser(id uint, name, email string) (*model.User, error)
	DeleteUser(id uint) error
}

// userService 用户服务实现
type userService struct {
	userRepo repository.UserRepository
	logger   log.Logger
}

// NewUserService 创建用户服务
func NewUserService(userRepo repository.UserRepository, logger log.Logger) UserService {
	return &userService{
		userRepo: userRepo,
		logger:   logger,
	}
}

// GetAllUsers 获取所有用户
func (s *userService) GetAllUsers() ([]model.User, error) {
	s.logger.Debug("服务层: 获取所有用户")
	return s.userRepo.FindAll()
}

// GetUserByID 根据ID获取用户
func (s *userService) GetUserByID(id uint) (*model.User, error) {
	s.logger.Debugf("服务层: 获取用户 ID=%d", id)
	return s.userRepo.FindByID(id)
}

// CreateUser 创建用户
func (s *userService) CreateUser(name, email string) (*model.User, error) {
	s.logger.Debugf("服务层: 创建用户 Name=%s, Email=%s", name, email)

	user := &model.User{
		Name:  name,
		Email: email,
	}

	return s.userRepo.Create(user)
}

// UpdateUser 更新用户
func (s *userService) UpdateUser(id uint, name, email string) (*model.User, error) {
	s.logger.Debugf("服务层: 更新用户 ID=%d", id)

	user, err := s.userRepo.FindByID(id)
	if err != nil {
		return nil, err
	}

	// 只更新提供的字段
	if name != "" {
		user.Name = name
	}

	if email != "" {
		user.Email = email
	}

	return s.userRepo.Update(user)
}

// DeleteUser 删除用户
func (s *userService) DeleteUser(id uint) error {
	s.logger.Debugf("服务层: 删除用户 ID=%d", id)
	return s.userRepo.Delete(id)
}
