package facades

import (
	"context"

	"gorm.io/gorm"
)

// Model 获取模型的查询构建器
func (db *DBFacade) Model(value interface{}) *gorm.DB {
	return GetGormDB().Model(value)
}

// Create 创建记录
func (db *DBFacade) Create(value interface{}) *gorm.DB {
	return GetGormDB().Create(value)
}

// Save 保存记录
func (db *DBFacade) Save(value interface{}) *gorm.DB {
	return GetGormDB().Save(value)
}

// First 获取第一条记录
func (db *DBFacade) First(dest interface{}, conds ...interface{}) *gorm.DB {
	return GetGormDB().First(dest, conds...)
}

// Find 查找记录
func (db *DBFacade) Find(dest interface{}, conds ...interface{}) *gorm.DB {
	return GetGormDB().Find(dest, conds...)
}

// Delete 删除记录
func (db *DBFacade) Delete(value interface{}, conds ...interface{}) *gorm.DB {
	return GetGormDB().Delete(value, conds...)
}

// Where 条件查询
func (db *DBFacade) Where(query interface{}, args ...interface{}) *gorm.DB {
	return GetGormDB().Where(query, args...)
}

// Transaction 在事务中执行函数
func (db *DBFacade) Transaction(fc func(tx *gorm.DB) error) error {
	return GetGormDB().Transaction(fc)
}

// WithContext 设置上下文
func (db *DBFacade) WithContext(ctx context.Context) *gorm.DB {
	return GetGormDB().WithContext(ctx)
}

// Raw 执行原始SQL
func (db *DBFacade) Raw(sql string, values ...interface{}) *gorm.DB {
	return GetGormDB().Raw(sql, values...)
}

// Exec 执行SQL语句
func (db *DBFacade) Exec(sql string, values ...interface{}) *gorm.DB {
	return GetGormDB().Exec(sql, values...)
}

// Table 指定表名
func (db *DBFacade) Table(name string) *gorm.DB {
	return GetGormDB().Table(name)
}

// Instance 获取原始数据库实例
func (db *DBFacade) Instance() *gorm.DB {
	return GetGormDB()
}
