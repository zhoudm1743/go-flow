package domain

import (
	"github.com/zhoudm1743/go-flow/core/domain"
	"gorm.io/gorm"
)

// UserRepository 用户仓储接口
type UserRepository interface {
	domain.Repository[*UserAggregate]

	// 扩展方法
	FindByUsername(username string) (*UserAggregate, error)
	FindByEmail(email string) (*UserAggregate, error)
	FindActiveUsers() ([]*UserAggregate, error)
}

// userRepository 用户仓储实现
type userRepository struct {
	*domain.BaseRepository[*UserAggregate]
}

// NewUserRepository 创建用户仓储
func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{
		BaseRepository: domain.NewBaseRepository[*UserAggregate](db),
	}
}

// FindByUsername 根据用户名查找用户
func (r *userRepository) FindByUsername(username string) (*UserAggregate, error) {
	spec := domain.NewSpecification[*UserAggregate](
		func(user *UserAggregate) bool {
			return user.Username == username
		},
		"username = ?", username,
	)

	users, err := r.FindBySpec(spec)
	if err != nil {
		return nil, err
	}

	if len(users) == 0 {
		return nil, nil
	}

	return users[0], nil
}

// FindByEmail 根据邮箱查找用户
func (r *userRepository) FindByEmail(email string) (*UserAggregate, error) {
	spec := domain.NewSpecification[*UserAggregate](
		func(user *UserAggregate) bool {
			return user.Email == email
		},
		"email = ?", email,
	)

	users, err := r.FindBySpec(spec)
	if err != nil {
		return nil, err
	}

	if len(users) == 0 {
		return nil, nil
	}

	return users[0], nil
}

// FindActiveUsers 查找所有激活用户
func (r *userRepository) FindActiveUsers() ([]*UserAggregate, error) {
	spec := domain.NewSpecification[*UserAggregate](
		func(user *UserAggregate) bool {
			return user.Status == 1
		},
		"status = ?", 1,
	)

	return r.FindBySpec(spec)
}

// UserSpecifications 用户规约集合
type UserSpecifications struct{}

// ActiveUsers 激活用户规约
func (UserSpecifications) ActiveUsers() domain.Specification[*UserAggregate] {
	return domain.NewSpecification[*UserAggregate](
		func(user *UserAggregate) bool { return user.Status == 1 },
		"status = ?", 1,
	)
}

// InactiveUsers 非激活用户规约
func (UserSpecifications) InactiveUsers() domain.Specification[*UserAggregate] {
	return domain.NewSpecification[*UserAggregate](
		func(user *UserAggregate) bool { return user.Status != 1 },
		"status != ?", 1,
	)
}

// UsernameContains 用户名包含指定字符串的规约
func (UserSpecifications) UsernameContains(keyword string) domain.Specification[*UserAggregate] {
	return domain.NewSpecification[*UserAggregate](
		func(user *UserAggregate) bool {
			return len(user.Username) > 0 && user.Username[:1] == keyword[:1]
		},
		"username LIKE ?", "%"+keyword+"%",
	)
}

// EmailDomain 指定邮箱域名的规约
func (UserSpecifications) EmailDomain(emailDomain string) domain.Specification[*UserAggregate] {
	return domain.NewSpecification[*UserAggregate](
		func(user *UserAggregate) bool {
			return len(user.Email) > len(emailDomain) && user.Email[len(user.Email)-len(emailDomain):] == emailDomain
		},
		"email LIKE ?", "%"+emailDomain,
	)
}
