package types

import "gorm.io/plugin/soft_delete"

type GormModel struct {
	ID        uint  `gorm:"primaryKey"`
	CreatedAt int64 `grom:"autoCreateTime"`
	UpdatedAt int64 `grom:"autoUpdateTime"`
}

type SoftDeleteModel struct {
	DeletedAt soft_delete.DeletedAt
}
