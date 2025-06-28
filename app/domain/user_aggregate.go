package domain

import (
	"fmt"
	"time"

	"github.com/zhoudm1743/go-flow/core/domain"
	"github.com/zhoudm1743/go-flow/pkg/types"
)

// UserAggregate 用户聚合根（DDD风格示例）
type UserAggregate struct {
	domain.BaseAggregateRoot
	Username    string      `gorm:"uniqueIndex;size:100;not null" json:"username"`
	Password    string      `gorm:"size:255;not null" json:"-"`
	Status      int8        `gorm:"default:1" json:"status"`
	Email       string      `gorm:"uniqueIndex;size:200;not null" json:"email"`
	Phone       types.Phone `gorm:"uniqueIndex;size:20;not null" json:"phone"`
	Avatar      string      `gorm:"size:255" json:"avatar"`
	LastLoginAt int64       `json:"last_login_at"`
	LastLoginIP string      `gorm:"size:45" json:"last_login_ip"`
}

// NewUserAggregate 创建新的用户聚合根
func NewUserAggregate(username, email string, phone types.Phone) *UserAggregate {
	user := &UserAggregate{
		BaseAggregateRoot: *domain.NewBaseAggregateRoot("User"),
		Username:          username,
		Email:             email,
		Phone:             phone,
		Status:            1, // 激活状态
	}

	// 生成用户ID
	user.SetID(fmt.Sprintf("user_%d", time.Now().UnixNano()))

	// 添加用户创建事件
	user.AddDomainEvent(domain.NewEntityCreatedEvent(
		user.GetAggregateID(),
		user.GetAggregateType(),
		map[string]interface{}{
			"username": username,
			"email":    email,
			"phone":    phone,
			"status":   user.Status,
		},
	))

	return user
}

// ChangePassword 修改密码（业务方法）
func (u *UserAggregate) ChangePassword(newPassword string) error {
	if newPassword == "" {
		return fmt.Errorf("密码不能为空")
	}

	oldPassword := u.Password
	u.Password = newPassword

	// 添加密码修改事件
	u.AddDomainEvent(&PasswordChangedEvent{
		BaseDomainEvent: *domain.NewBaseDomainEvent(
			"user.password_changed",
			u.GetAggregateID(),
			u.GetAggregateType(),
			map[string]interface{}{
				"username":   u.Username,
				"changed_at": time.Now(),
			},
		),
		OldPasswordHash: oldPassword,
		NewPasswordHash: newPassword,
	})

	return nil
}

// Login 用户登录（业务方法）
func (u *UserAggregate) Login(ip string) error {
	if u.Status != 1 {
		return fmt.Errorf("用户状态异常，无法登录")
	}

	u.LastLoginAt = time.Now().Unix()
	u.LastLoginIP = ip

	// 添加用户登录事件
	u.AddDomainEvent(&UserLoginEvent{
		BaseDomainEvent: *domain.NewBaseDomainEvent(
			"user.login",
			u.GetAggregateID(),
			u.GetAggregateType(),
			map[string]interface{}{
				"username":   u.Username,
				"login_ip":   ip,
				"login_time": u.LastLoginAt,
			},
		),
		LoginIP:   ip,
		LoginTime: time.Unix(u.LastLoginAt, 0),
	})

	return nil
}

// Deactivate 停用用户（业务方法）
func (u *UserAggregate) Deactivate(reason string) error {
	if u.Status == 0 {
		return fmt.Errorf("用户已经是停用状态")
	}

	u.Status = 0

	// 添加用户停用事件
	u.AddDomainEvent(&UserDeactivatedEvent{
		BaseDomainEvent: *domain.NewBaseDomainEvent(
			"user.deactivated",
			u.GetAggregateID(),
			u.GetAggregateType(),
			map[string]interface{}{
				"username":       u.Username,
				"reason":         reason,
				"deactivated_at": time.Now(),
			},
		),
		Reason: reason,
	})

	return nil
}

// CanBeDeleted 实现聚合根接口
func (u *UserAggregate) CanBeDeleted() bool {
	// 只有停用状态的用户才能被删除
	return u.Status == 0
}

// Validate 验证聚合根
func (u *UserAggregate) Validate() error {
	if err := u.BaseAggregateRoot.Validate(); err != nil {
		return err
	}

	if u.Username == "" {
		return fmt.Errorf("用户名不能为空")
	}

	if u.Email == "" {
		return fmt.Errorf("邮箱不能为空")
	}

	return nil
}

// IsActive 检查用户是否激活
func (u *UserAggregate) IsActive() bool {
	return u.Status == 1
}

// GetDisplayName 获取显示名称
func (u *UserAggregate) GetDisplayName() string {
	if u.Username != "" {
		return u.Username
	}
	return u.Email
}

// PasswordChangedEvent 密码修改事件
type PasswordChangedEvent struct {
	domain.BaseDomainEvent
	OldPasswordHash string `json:"old_password_hash"`
	NewPasswordHash string `json:"new_password_hash"`
}

// UserLoginEvent 用户登录事件
type UserLoginEvent struct {
	domain.BaseDomainEvent
	LoginIP   string    `json:"login_ip"`
	LoginTime time.Time `json:"login_time"`
}

// UserDeactivatedEvent 用户停用事件
type UserDeactivatedEvent struct {
	domain.BaseDomainEvent
	Reason string `json:"reason"`
}
