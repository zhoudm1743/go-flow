package system

import (
	"fmt"
	"sync"

	"github.com/zhoudm1743/go-flow/app/models"
	"github.com/zhoudm1743/go-flow/core/database"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// TenantDatabaseService 租户数据库连接管理服务
type TenantDatabaseService interface {
	GetTenantDB(tenantID uint) (*gorm.DB, error)
	CreateTenantDB(tenant *models.SystemTenant) error
	TestTenantDBConnection(tenant *models.SystemTenant) error
	InitTenantTables(tenantID uint) error
	RemoveTenantDB(tenantID uint)
}

type tenantDatabaseService struct {
	db          database.Database
	tenantDBs   map[uint]*gorm.DB
	tenantMutex sync.RWMutex
}

// GetTenantDB 获取租户数据库连接
func (t *tenantDatabaseService) GetTenantDB(tenantID uint) (*gorm.DB, error) {
	t.tenantMutex.RLock()
	if db, exists := t.tenantDBs[tenantID]; exists {
		t.tenantMutex.RUnlock()
		return db, nil
	}
	t.tenantMutex.RUnlock()

	// 如果不存在，则从数据库获取租户信息并创建连接
	var tenant models.SystemTenant
	if err := t.db.GetDB().Where("id = ? AND status = 1", tenantID).First(&tenant).Error; err != nil {
		return nil, fmt.Errorf("租户不存在或已禁用: %v", err)
	}

	return t.createAndCacheTenantDB(&tenant)
}

// CreateTenantDB 创建租户数据库连接
func (t *tenantDatabaseService) CreateTenantDB(tenant *models.SystemTenant) error {
	_, err := t.createAndCacheTenantDB(tenant)
	return err
}

// TestTenantDBConnection 测试租户数据库连接
func (t *tenantDatabaseService) TestTenantDBConnection(tenant *models.SystemTenant) error {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		tenant.DatabaseUser,
		tenant.DatabasePassword,
		tenant.DatabaseHost,
		tenant.DatabasePort,
		tenant.DatabaseName,
	)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("数据库连接失败: %v", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("获取数据库实例失败: %v", err)
	}
	defer sqlDB.Close()

	if err := sqlDB.Ping(); err != nil {
		return fmt.Errorf("数据库ping失败: %v", err)
	}

	return nil
}

// InitTenantTables 初始化租户表结构
func (t *tenantDatabaseService) InitTenantTables(tenantID uint) error {
	tenantDB, err := t.GetTenantDB(tenantID)
	if err != nil {
		return err
	}

	// 自动迁移租户相关的表
	err = tenantDB.AutoMigrate(
		&models.SystemAdmin{},
		&models.SystemMenu{},
		&models.SystemTenantRole{},
		&models.SystemTenantPermission{},
		&models.SystemTenantRolePermission{},
		&models.SystemTenantUserRole{},
		&models.SystemCasbinRule{},
	)

	if err != nil {
		return fmt.Errorf("租户表初始化失败: %v", err)
	}

	return nil
}

// RemoveTenantDB 移除租户数据库连接（用于租户删除或禁用时）
func (t *tenantDatabaseService) RemoveTenantDB(tenantID uint) {
	t.tenantMutex.Lock()
	defer t.tenantMutex.Unlock()

	if db, exists := t.tenantDBs[tenantID]; exists {
		if sqlDB, err := db.DB(); err == nil {
			sqlDB.Close()
		}
		delete(t.tenantDBs, tenantID)
	}
}

// createAndCacheTenantDB 创建并缓存租户数据库连接
func (t *tenantDatabaseService) createAndCacheTenantDB(tenant *models.SystemTenant) (*gorm.DB, error) {
	t.tenantMutex.Lock()
	defer t.tenantMutex.Unlock()

	// 双重检查锁定
	if db, exists := t.tenantDBs[tenant.ID]; exists {
		return db, nil
	}

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		tenant.DatabaseUser,
		tenant.DatabasePassword,
		tenant.DatabaseHost,
		tenant.DatabasePort,
		tenant.DatabaseName,
	)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("创建租户数据库连接失败: %v", err)
	}

	// 设置连接池
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("获取数据库实例失败: %v", err)
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)

	// 缓存连接
	t.tenantDBs[tenant.ID] = db

	return db, nil
}

// NewTenantDatabaseService 创建租户数据库服务
func NewTenantDatabaseService(db database.Database) TenantDatabaseService {
	return &tenantDatabaseService{
		db:        db,
		tenantDBs: make(map[uint]*gorm.DB),
	}
}
