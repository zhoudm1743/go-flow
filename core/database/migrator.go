package database

import (
	"github.com/zhoudm1743/go-flow/core/logger"
)

// Migrator 数据库迁移器
type Migrator struct {
	db     Database
	logger logger.Logger
}

// NewMigrator 创建新的迁移器
func NewMigrator(db Database, log logger.Logger) *Migrator {
	return &Migrator{
		db:     db,
		logger: log,
	}
}

// AutoMigrate 自动迁移所有模型
func (m *Migrator) AutoMigrate() error {
	models := GetAllModels()

	m.logger.WithField("models", len(models)).Info("开始数据库自动迁移")

	if err := m.db.AutoMigrate(models...); err != nil {
		m.logger.WithError(err).Error("数据库迁移失败")
		return err
	}

	m.logger.Info("数据库迁移完成")
	return nil
}
