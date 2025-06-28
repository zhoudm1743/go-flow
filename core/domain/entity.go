package domain

import (
	"fmt"
	"reflect"
	"time"
)

// Entity 领域实体接口
type Entity interface {
	ID() string
	SetID(id string)
	GetVersion() int64
	SetVersion(version int64)
	IsNew() bool
	Validate() error
	GetDomainEvents() []DomainEvent
	AddDomainEvent(event DomainEvent)
	ClearDomainEvents()
}

// BaseEntity 基础实体
type BaseEntity struct {
	id           string        `json:"id" gorm:"primaryKey"`
	version      int64         `json:"version" gorm:"column:version"`
	domainEvents []DomainEvent `json:"-" gorm:"-"`
	CreatedAt    time.Time     `json:"created_at" gorm:"column:created_at"`
	UpdatedAt    time.Time     `json:"updated_at" gorm:"column:updated_at"`
}

// ID 获取实体ID
func (e *BaseEntity) ID() string {
	return e.id
}

// SetID 设置实体ID
func (e *BaseEntity) SetID(id string) {
	e.id = id
}

// GetVersion 获取版本号
func (e *BaseEntity) GetVersion() int64 {
	return e.version
}

// SetVersion 设置版本号
func (e *BaseEntity) SetVersion(version int64) {
	e.version = version
}

// IsNew 判断是否为新实体
func (e *BaseEntity) IsNew() bool {
	return e.id == ""
}

// Validate 验证实体
func (e *BaseEntity) Validate() error {
	if e.id == "" {
		return fmt.Errorf("entity ID cannot be empty")
	}
	return nil
}

// GetDomainEvents 获取领域事件
func (e *BaseEntity) GetDomainEvents() []DomainEvent {
	return e.domainEvents
}

// AddDomainEvent 添加领域事件
func (e *BaseEntity) AddDomainEvent(event DomainEvent) {
	e.domainEvents = append(e.domainEvents, event)
}

// ClearDomainEvents 清空领域事件
func (e *BaseEntity) ClearDomainEvents() {
	e.domainEvents = nil
}

// AggregateRoot 聚合根接口
type AggregateRoot interface {
	Entity
	GetAggregateID() string
	GetAggregateType() string
	CanBeDeleted() bool
}

// BaseAggregateRoot 基础聚合根
type BaseAggregateRoot struct {
	BaseEntity
	aggregateType string
}

// NewBaseAggregateRoot 创建基础聚合根
func NewBaseAggregateRoot(aggregateType string) *BaseAggregateRoot {
	return &BaseAggregateRoot{
		BaseEntity:    BaseEntity{},
		aggregateType: aggregateType,
	}
}

// GetAggregateID 获取聚合ID
func (ar *BaseAggregateRoot) GetAggregateID() string {
	return ar.ID()
}

// GetAggregateType 获取聚合类型
func (ar *BaseAggregateRoot) GetAggregateType() string {
	if ar.aggregateType != "" {
		return ar.aggregateType
	}
	// 默认使用类型名
	return reflect.TypeOf(ar).Elem().Name()
}

// CanBeDeleted 是否可以删除
func (ar *BaseAggregateRoot) CanBeDeleted() bool {
	// 默认可以删除，子类可以重写
	return true
}

// ValueObject 值对象接口
type ValueObject interface {
	Equals(other ValueObject) bool
	Validate() error
	String() string
}

// Repository 仓储接口
type Repository[T AggregateRoot] interface {
	// 基础CRUD操作
	Save(aggregate T) error
	FindByID(id string) (T, error)
	Delete(id string) error

	// 查询操作
	FindAll() ([]T, error)
	FindBySpec(spec Specification[T]) ([]T, error)
	Count() (int64, error)
	CountBySpec(spec Specification[T]) (int64, error)

	// 分页查询
	FindWithPagination(offset, limit int) ([]T, int64, error)
	FindBySpecWithPagination(spec Specification[T], offset, limit int) ([]T, int64, error)

	// 事务支持
	InTransaction(fn func(repo Repository[T]) error) error
}

// Specification 规约模式接口
type Specification[T AggregateRoot] interface {
	IsSatisfiedBy(candidate T) bool
	And(other Specification[T]) Specification[T]
	Or(other Specification[T]) Specification[T]
	Not() Specification[T]
	ToSQL() (string, []interface{})
}

// BaseSpecification 基础规约实现
type BaseSpecification[T AggregateRoot] struct {
	expression func(T) bool
	sqlClause  string
	sqlParams  []interface{}
}

// NewSpecification 创建新规约
func NewSpecification[T AggregateRoot](expr func(T) bool, sql string, params ...interface{}) *BaseSpecification[T] {
	return &BaseSpecification[T]{
		expression: expr,
		sqlClause:  sql,
		sqlParams:  params,
	}
}

// IsSatisfiedBy 检查是否满足规约
func (s *BaseSpecification[T]) IsSatisfiedBy(candidate T) bool {
	if s.expression == nil {
		return true
	}
	return s.expression(candidate)
}

// And 逻辑与操作
func (s *BaseSpecification[T]) And(other Specification[T]) Specification[T] {
	return &CompositeSpecification[T]{
		left:     s,
		right:    other,
		operator: "AND",
	}
}

// Or 逻辑或操作
func (s *BaseSpecification[T]) Or(other Specification[T]) Specification[T] {
	return &CompositeSpecification[T]{
		left:     s,
		right:    other,
		operator: "OR",
	}
}

// Not 逻辑非操作
func (s *BaseSpecification[T]) Not() Specification[T] {
	return &NotSpecification[T]{
		spec: s,
	}
}

// ToSQL 转换为SQL条件
func (s *BaseSpecification[T]) ToSQL() (string, []interface{}) {
	return s.sqlClause, s.sqlParams
}

// CompositeSpecification 组合规约
type CompositeSpecification[T AggregateRoot] struct {
	left     Specification[T]
	right    Specification[T]
	operator string
}

// IsSatisfiedBy 检查是否满足规约
func (c *CompositeSpecification[T]) IsSatisfiedBy(candidate T) bool {
	switch c.operator {
	case "AND":
		return c.left.IsSatisfiedBy(candidate) && c.right.IsSatisfiedBy(candidate)
	case "OR":
		return c.left.IsSatisfiedBy(candidate) || c.right.IsSatisfiedBy(candidate)
	default:
		return false
	}
}

// And 逻辑与操作
func (c *CompositeSpecification[T]) And(other Specification[T]) Specification[T] {
	return &CompositeSpecification[T]{
		left:     c,
		right:    other,
		operator: "AND",
	}
}

// Or 逻辑或操作
func (c *CompositeSpecification[T]) Or(other Specification[T]) Specification[T] {
	return &CompositeSpecification[T]{
		left:     c,
		right:    other,
		operator: "OR",
	}
}

// Not 逻辑非操作
func (c *CompositeSpecification[T]) Not() Specification[T] {
	return &NotSpecification[T]{
		spec: c,
	}
}

// ToSQL 转换为SQL条件
func (c *CompositeSpecification[T]) ToSQL() (string, []interface{}) {
	leftSQL, leftParams := c.left.ToSQL()
	rightSQL, rightParams := c.right.ToSQL()

	if leftSQL == "" {
		return rightSQL, rightParams
	}
	if rightSQL == "" {
		return leftSQL, leftParams
	}

	sql := fmt.Sprintf("(%s) %s (%s)", leftSQL, c.operator, rightSQL)
	params := append(leftParams, rightParams...)

	return sql, params
}

// NotSpecification 非规约
type NotSpecification[T AggregateRoot] struct {
	spec Specification[T]
}

// IsSatisfiedBy 检查是否满足规约
func (n *NotSpecification[T]) IsSatisfiedBy(candidate T) bool {
	return !n.spec.IsSatisfiedBy(candidate)
}

// And 逻辑与操作
func (n *NotSpecification[T]) And(other Specification[T]) Specification[T] {
	return &CompositeSpecification[T]{
		left:     n,
		right:    other,
		operator: "AND",
	}
}

// Or 逻辑或操作
func (n *NotSpecification[T]) Or(other Specification[T]) Specification[T] {
	return &CompositeSpecification[T]{
		left:     n,
		right:    other,
		operator: "OR",
	}
}

// Not 逻辑非操作
func (n *NotSpecification[T]) Not() Specification[T] {
	return n.spec
}

// ToSQL 转换为SQL条件
func (n *NotSpecification[T]) ToSQL() (string, []interface{}) {
	sql, params := n.spec.ToSQL()
	if sql == "" {
		return "", nil
	}
	return fmt.Sprintf("NOT (%s)", sql), params
}

// DomainService 领域服务接口
type DomainService interface {
	Name() string
}

// BaseDomainService 基础领域服务
type BaseDomainService struct {
	name string
}

// NewBaseDomainService 创建基础领域服务
func NewBaseDomainService(name string) *BaseDomainService {
	return &BaseDomainService{name: name}
}

// Name 获取服务名称
func (s *BaseDomainService) Name() string {
	return s.name
}
