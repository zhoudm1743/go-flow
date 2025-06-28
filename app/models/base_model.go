package models

// BaseModel 基础模型
type BaseModel struct {
	ID        uint  `gorm:"primarykey" json:"id"`
	CreatedAt int64 `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt int64 `json:"updated_at" gorm:"autoUpdateTime"`
}
