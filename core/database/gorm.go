package database

import (
	"fmt"
	"time"

	"github.com/zhoudm1743/go-flow/core/config"
	"github.com/zhoudm1743/go-flow/core/logger"
	"go.uber.org/fx"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"
)

// Database 数据库接口
type Database interface {
	GetDB() *gorm.DB
	Close() error
	Ping() error
	AutoMigrate(dst ...interface{}) error
}

// GormDatabase GORM 数据库实现
type GormDatabase struct {
	db     *gorm.DB
	logger logger.Logger
}

// GetDB 获取 GORM 数据库实例
func (d *GormDatabase) GetDB() *gorm.DB {
	return d.db
}

// Close 关闭数据库连接
func (d *GormDatabase) Close() error {
	sqlDB, err := d.db.DB()
	if err != nil {
		return fmt.Errorf("获取底层数据库实例失败: %w", err)
	}
	return sqlDB.Close()
}

// Ping 检查数据库连接
func (d *GormDatabase) Ping() error {
	sqlDB, err := d.db.DB()
	if err != nil {
		return fmt.Errorf("获取底层数据库实例失败: %w", err)
	}
	return sqlDB.Ping()
}

// AutoMigrate 自动迁移数据库表结构
func (d *GormDatabase) AutoMigrate(dst ...interface{}) error {
	return d.db.AutoMigrate(dst...)
}

// Module fx模块
var Module = fx.Options(
	fx.Provide(NewDatabase),
)

// NewDatabase 创建新的数据库实例
func NewDatabase(cfg *config.Config, log logger.Logger) (Database, error) {
	// 构建 MySQL 连接字符串
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.Database.Username,
		cfg.Database.Password,
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.Database,
	)

	// 配置 GORM 日志器
	gormLogLevel := gormLogger.Silent
	switch cfg.Log.Level {
	case "debug":
		gormLogLevel = gormLogger.Info
	case "info":
		gormLogLevel = gormLogger.Warn
	case "warn", "warning":
		gormLogLevel = gormLogger.Error
	case "error":
		gormLogLevel = gormLogger.Silent
	}

	// 自定义 GORM 日志器，使用我们的 logger
	customLogger := gormLogger.New(
		&GormLogWriter{logger: log},
		gormLogger.Config{
			SlowThreshold: time.Second, // 慢查询阈值
			LogLevel:      gormLogLevel,
			Colorful:      true, // 启用颜色
		},
	)

	// 连接数据库
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: customLogger,
	})
	if err != nil {
		return nil, fmt.Errorf("连接数据库失败: %w", err)
	}

	// 获取底层的 sql.DB 实例进行连接池配置
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("获取底层数据库实例失败: %w", err)
	}

	// 配置连接池
	sqlDB.SetMaxIdleConns(10)           // 设置空闲连接池中连接的最大数量
	sqlDB.SetMaxOpenConns(100)          // 设置打开数据库连接的最大数量
	sqlDB.SetConnMaxLifetime(time.Hour) // 设置了连接可复用的最大时间

	// 测试连接
	if err := sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("数据库连接测试失败: %w", err)
	}

	log.WithFields(map[string]interface{}{
		"host":     cfg.Database.Host,
		"port":     cfg.Database.Port,
		"database": cfg.Database.Database,
		"username": cfg.Database.Username,
	}).Info("数据库连接成功")

	return &GormDatabase{
		db:     db,
		logger: log,
	}, nil
}

// GormLogWriter GORM 日志写入器，将 GORM 日志写入我们的 logger
type GormLogWriter struct {
	logger logger.Logger
}

// Printf 实现 gorm logger Writer 接口
func (w *GormLogWriter) Printf(format string, args ...interface{}) {
	w.logger.Infof(format, args...)
}
