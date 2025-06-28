package database

import "github.com/zhoudm1743/go-flow/app/models"

// GetAllModels 获取所有需要迁移的模型
func GetAllModels() []interface{} {
	return []interface{}{
		&models.SystemAdmin{},
		&models.SystemMenu{},
	}
}
