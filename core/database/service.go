package database

import (
	"github.com/zhoudm1743/go-flow/core/logger"
	"go.uber.org/fx"
)

// UserService 用户服务示例
type UserService struct {
	db     Database
	logger logger.Logger
}

// NewUserService 创建用户服务
func NewUserService(db Database, log logger.Logger) *UserService {
	return &UserService{
		db:     db,
		logger: log,
	}
}

// CreateUser 创建用户
func (s *UserService) CreateUser(user *User) error {
	if err := s.db.GetDB().Create(user).Error; err != nil {
		s.logger.WithError(err).Error("创建用户失败")
		return err
	}

	s.logger.WithFields(map[string]interface{}{
		"user_id":  user.ID,
		"username": user.Username,
		"email":    user.Email,
	}).Info("用户创建成功")

	return nil
}

// GetUserByID 根据ID获取用户
func (s *UserService) GetUserByID(id uint) (*User, error) {
	var user User
	if err := s.db.GetDB().First(&user, id).Error; err != nil {
		s.logger.WithError(err).WithField("user_id", id).Error("获取用户失败")
		return nil, err
	}

	return &user, nil
}

// GetUserByUsername 根据用户名获取用户
func (s *UserService) GetUserByUsername(username string) (*User, error) {
	var user User
	if err := s.db.GetDB().Where("username = ?", username).First(&user).Error; err != nil {
		s.logger.WithError(err).WithField("username", username).Error("获取用户失败")
		return nil, err
	}

	return &user, nil
}

// UpdateUser 更新用户
func (s *UserService) UpdateUser(user *User) error {
	if err := s.db.GetDB().Save(user).Error; err != nil {
		s.logger.WithError(err).WithField("user_id", user.ID).Error("更新用户失败")
		return err
	}

	s.logger.WithField("user_id", user.ID).Info("用户更新成功")
	return nil
}

// DeleteUser 删除用户（软删除）
func (s *UserService) DeleteUser(id uint) error {
	if err := s.db.GetDB().Delete(&User{}, id).Error; err != nil {
		s.logger.WithError(err).WithField("user_id", id).Error("删除用户失败")
		return err
	}

	s.logger.WithField("user_id", id).Info("用户删除成功")
	return nil
}

// ListUsers 获取用户列表
func (s *UserService) ListUsers(limit, offset int) ([]User, error) {
	var users []User
	if err := s.db.GetDB().Limit(limit).Offset(offset).Find(&users).Error; err != nil {
		s.logger.WithError(err).Error("获取用户列表失败")
		return nil, err
	}

	return users, nil
}

// ServiceModule 服务模块
var ServiceModule = fx.Options(
	fx.Provide(NewUserService),
)
