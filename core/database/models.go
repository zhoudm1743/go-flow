package database

import (
	"time"

	"gorm.io/gorm"
)

// BaseModel 基础模型
type BaseModel struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at"`
}

// User 用户模型示例
type User struct {
	BaseModel
	Username string `gorm:"uniqueIndex;size:100;not null" json:"username"`
	Email    string `gorm:"uniqueIndex;size:200;not null" json:"email"`
	Password string `gorm:"size:255;not null" json:"-"`
	Status   int    `gorm:"default:1" json:"status"`
}

// TableName 指定表名
func (User) TableName() string {
	return "users"
}

// Post 文章模型示例
type Post struct {
	BaseModel
	Title   string `gorm:"size:200;not null" json:"title"`
	Content string `gorm:"type:text" json:"content"`
	UserID  uint   `gorm:"not null;index" json:"user_id"`
	User    User   `gorm:"foreignKey:UserID" json:"user"`
	Status  int    `gorm:"default:1" json:"status"`
}

// TableName 指定表名
func (Post) TableName() string {
	return "posts"
}

// GetAllModels 获取所有需要迁移的模型
func GetAllModels() []interface{} {
	return []interface{}{
		&User{},
		&Post{},
	}
}
