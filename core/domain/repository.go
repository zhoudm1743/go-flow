package domain

import (
	"fmt"
	"reflect"

	"gorm.io/gorm"
)

// BaseRepository GORM基础仓储实现
type BaseRepository[T AggregateRoot] struct {
	db    *gorm.DB
	table string
}

// NewBaseRepository 创建基础仓储
func NewBaseRepository[T AggregateRoot](db *gorm.DB) *BaseRepository[T] {
	repo := &BaseRepository[T]{
		db: db,
	}

	// 获取实体类型名作为表名
	var zero T
	entityType := reflect.TypeOf(zero).Elem()
	repo.table = entityType.Name()

	return repo
}

// Save 保存聚合根
func (r *BaseRepository[T]) Save(aggregate T) error {
	if aggregate.IsNew() {
		return r.create(aggregate)
	}
	return r.update(aggregate)
}

// create 创建新聚合根
func (r *BaseRepository[T]) create(aggregate T) error {
	// 添加创建事件
	aggregate.AddDomainEvent(NewEntityCreatedEvent(
		aggregate.GetAggregateID(),
		aggregate.GetAggregateType(),
		r.aggregateToMap(aggregate),
	))

	result := r.db.Create(aggregate)
	if result.Error != nil {
		return fmt.Errorf("failed to create %s: %w", r.table, result.Error)
	}

	return nil
}

// update 更新聚合根
func (r *BaseRepository[T]) update(aggregate T) error {
	// 获取原始版本进行乐观锁检查
	originalVersion := aggregate.GetVersion()

	// 增加版本号
	aggregate.SetVersion(originalVersion + 1)

	// 添加更新事件
	aggregate.AddDomainEvent(NewEntityUpdatedEvent(
		aggregate.GetAggregateID(),
		aggregate.GetAggregateType(),
		r.aggregateToMap(aggregate),
	))

	// 使用乐观锁更新
	result := r.db.Model(aggregate).
		Where("version = ?", originalVersion).
		Updates(aggregate)

	if result.Error != nil {
		return fmt.Errorf("failed to update %s: %w", r.table, result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("optimistic lock failure: %s has been modified by another transaction", r.table)
	}

	return nil
}

// FindByID 根据ID查找聚合根
func (r *BaseRepository[T]) FindByID(id string) (T, error) {
	var aggregate T
	result := r.db.Where("id = ?", id).First(&aggregate)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return aggregate, fmt.Errorf("%s with id %s not found", r.table, id)
		}
		return aggregate, fmt.Errorf("failed to find %s: %w", r.table, result.Error)
	}

	return aggregate, nil
}

// Delete 根据ID删除聚合根
func (r *BaseRepository[T]) Delete(id string) error {
	// 先查找聚合根以获取完整信息
	aggregate, err := r.FindByID(id)
	if err != nil {
		return err
	}

	// 检查是否可以删除
	if !aggregate.CanBeDeleted() {
		return fmt.Errorf("%s with id %s cannot be deleted", r.table, id)
	}

	// 添加删除事件
	aggregate.AddDomainEvent(NewEntityDeletedEvent(
		aggregate.GetAggregateID(),
		aggregate.GetAggregateType(),
	))

	result := r.db.Delete(&aggregate)
	if result.Error != nil {
		return fmt.Errorf("failed to delete %s: %w", r.table, result.Error)
	}

	return nil
}

// FindAll 查找所有聚合根
func (r *BaseRepository[T]) FindAll() ([]T, error) {
	var aggregates []T
	result := r.db.Find(&aggregates)

	if result.Error != nil {
		return nil, fmt.Errorf("failed to find all %s: %w", r.table, result.Error)
	}

	return aggregates, nil
}

// FindBySpec 根据规约查找聚合根
func (r *BaseRepository[T]) FindBySpec(spec Specification[T]) ([]T, error) {
	var aggregates []T

	query := r.db
	sql, params := spec.ToSQL()
	if sql != "" {
		query = query.Where(sql, params...)
	}

	result := query.Find(&aggregates)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to find %s by spec: %w", r.table, result.Error)
	}

	// 应用内存过滤（如果SQL过滤不够精确）
	var filteredAggregates []T
	for _, aggregate := range aggregates {
		if spec.IsSatisfiedBy(aggregate) {
			filteredAggregates = append(filteredAggregates, aggregate)
		}
	}

	return filteredAggregates, nil
}

// Count 统计总数
func (r *BaseRepository[T]) Count() (int64, error) {
	var count int64
	result := r.db.Model(new(T)).Count(&count)

	if result.Error != nil {
		return 0, fmt.Errorf("failed to count %s: %w", r.table, result.Error)
	}

	return count, nil
}

// CountBySpec 根据规约统计数量
func (r *BaseRepository[T]) CountBySpec(spec Specification[T]) (int64, error) {
	var count int64

	query := r.db.Model(new(T))
	sql, params := spec.ToSQL()
	if sql != "" {
		query = query.Where(sql, params...)
	}

	result := query.Count(&count)
	if result.Error != nil {
		return 0, fmt.Errorf("failed to count %s by spec: %w", r.table, result.Error)
	}

	return count, nil
}

// FindWithPagination 分页查找
func (r *BaseRepository[T]) FindWithPagination(offset, limit int) ([]T, int64, error) {
	var aggregates []T
	var total int64

	// 先获取总数
	if err := r.db.Model(new(T)).Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count %s: %w", r.table, err)
	}

	// 分页查询
	result := r.db.Offset(offset).Limit(limit).Find(&aggregates)
	if result.Error != nil {
		return nil, 0, fmt.Errorf("failed to find %s with pagination: %w", r.table, result.Error)
	}

	return aggregates, total, nil
}

// FindBySpecWithPagination 根据规约分页查找
func (r *BaseRepository[T]) FindBySpecWithPagination(spec Specification[T], offset, limit int) ([]T, int64, error) {
	var aggregates []T
	var total int64

	// 构建查询
	query := r.db
	sql, params := spec.ToSQL()
	if sql != "" {
		query = query.Where(sql, params...)
	}

	// 先获取总数
	if err := query.Model(new(T)).Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count %s by spec: %w", r.table, err)
	}

	// 分页查询
	result := query.Offset(offset).Limit(limit).Find(&aggregates)
	if result.Error != nil {
		return nil, 0, fmt.Errorf("failed to find %s by spec with pagination: %w", r.table, result.Error)
	}

	// 应用内存过滤
	var filteredAggregates []T
	for _, aggregate := range aggregates {
		if spec.IsSatisfiedBy(aggregate) {
			filteredAggregates = append(filteredAggregates, aggregate)
		}
	}

	return filteredAggregates, total, nil
}

// InTransaction 在事务中执行操作
func (r *BaseRepository[T]) InTransaction(fn func(repo Repository[T]) error) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// 创建事务版本的仓储
		txRepo := &BaseRepository[T]{
			db:    tx,
			table: r.table,
		}

		return fn(txRepo)
	})
}

// aggregateToMap 将聚合根转换为map用于事件数据
func (r *BaseRepository[T]) aggregateToMap(aggregate T) map[string]interface{} {
	result := make(map[string]interface{})

	// 使用反射获取聚合根的字段值
	val := reflect.ValueOf(aggregate)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	typ := val.Type()
	for i := 0; i < val.NumField(); i++ {
		field := typ.Field(i)
		value := val.Field(i)

		// 跳过私有字段和不可导出的字段
		if !field.IsExported() {
			continue
		}

		// 跳过方法和函数
		if value.Kind() == reflect.Func {
			continue
		}

		result[field.Name] = value.Interface()
	}

	return result
}

// UnitOfWork 工作单元接口
type UnitOfWork interface {
	RegisterNew(entity AggregateRoot)
	RegisterDirty(entity AggregateRoot)
	RegisterRemoved(entity AggregateRoot)
	Commit() error
	Rollback() error
}

// SimpleUnitOfWork 简单工作单元实现
type SimpleUnitOfWork struct {
	db              *gorm.DB
	newEntities     []AggregateRoot
	dirtyEntities   []AggregateRoot
	removedEntities []AggregateRoot
	eventDispatcher DomainEventDispatcher
}

// NewSimpleUnitOfWork 创建简单工作单元
func NewSimpleUnitOfWork(db *gorm.DB, eventDispatcher DomainEventDispatcher) *SimpleUnitOfWork {
	return &SimpleUnitOfWork{
		db:              db,
		newEntities:     make([]AggregateRoot, 0),
		dirtyEntities:   make([]AggregateRoot, 0),
		removedEntities: make([]AggregateRoot, 0),
		eventDispatcher: eventDispatcher,
	}
}

// RegisterNew 注册新实体
func (uow *SimpleUnitOfWork) RegisterNew(entity AggregateRoot) {
	uow.newEntities = append(uow.newEntities, entity)
}

// RegisterDirty 注册脏实体
func (uow *SimpleUnitOfWork) RegisterDirty(entity AggregateRoot) {
	uow.dirtyEntities = append(uow.dirtyEntities, entity)
}

// RegisterRemoved 注册删除实体
func (uow *SimpleUnitOfWork) RegisterRemoved(entity AggregateRoot) {
	uow.removedEntities = append(uow.removedEntities, entity)
}

// Commit 提交所有更改
func (uow *SimpleUnitOfWork) Commit() error {
	return uow.db.Transaction(func(tx *gorm.DB) error {
		var allEvents []DomainEvent

		// 处理新实体
		for _, entity := range uow.newEntities {
			if err := tx.Create(entity).Error; err != nil {
				return fmt.Errorf("failed to create entity: %w", err)
			}
			allEvents = append(allEvents, entity.GetDomainEvents()...)
			entity.ClearDomainEvents()
		}

		// 处理脏实体
		for _, entity := range uow.dirtyEntities {
			if err := tx.Save(entity).Error; err != nil {
				return fmt.Errorf("failed to update entity: %w", err)
			}
			allEvents = append(allEvents, entity.GetDomainEvents()...)
			entity.ClearDomainEvents()
		}

		// 处理删除实体
		for _, entity := range uow.removedEntities {
			if err := tx.Delete(entity).Error; err != nil {
				return fmt.Errorf("failed to delete entity: %w", err)
			}
			allEvents = append(allEvents, entity.GetDomainEvents()...)
			entity.ClearDomainEvents()
		}

		// 分发领域事件
		if uow.eventDispatcher != nil {
			if err := uow.eventDispatcher.DispatchAll(allEvents); err != nil {
				return fmt.Errorf("failed to dispatch domain events: %w", err)
			}
		}

		return nil
	})
}

// Rollback 回滚更改
func (uow *SimpleUnitOfWork) Rollback() error {
	// 清空所有注册的实体
	uow.newEntities = uow.newEntities[:0]
	uow.dirtyEntities = uow.dirtyEntities[:0]
	uow.removedEntities = uow.removedEntities[:0]
	return nil
}
