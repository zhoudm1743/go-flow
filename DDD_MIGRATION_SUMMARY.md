# DDD 领域事件系统迁移总结

## 🎯 迁移概述

本次迁移将框架从**双事件系统并行**改为**专注DDD领域事件系统**，去掉了原有的`core/event`框架级事件系统，主要使用DDD的领域事件机制。

## 📊 迁移对比

### 迁移前
```
┌─ 框架级事件系统 (core/event/)
│  ├── EventBus (事件总线)
│  ├── EventStore (事件存储)
│  ├── DelayedEventScheduler (延时事件调度)
│  └── EventService (事件服务)
└─ DDD领域事件系统 (core/domain/)
   ├── DomainEvent (领域事件接口)
   ├── DomainEventDispatcher (事件分发器)
   └── BaseRepository (自动事件分发)
```

### 迁移后
```
└─ DDD领域事件系统 (core/domain/)
   ├── DomainEvent (领域事件接口)
   ├── DomainEventDispatcher (异步事件分发器)
   ├── BaseRepository (自动事件分发)
   ├── UnitOfWork (工作单元模式)
   ├── Specification (规约模式)
   └── DomainServiceRegistry (领域服务注册表)
```

## 🛠️ 执行的修改

### 1. 删除原有事件系统
```bash
# 删除的文件
core/event/module.go
core/event/service.go
core/event/store.go
core/event/bus.go
core/event/types.go
core/event/delayed_scheduler.go
core/event/cron_integration.go
core/event/README.md
```

### 2. 更新启动配置
```go
// boot/bootstrap.go - 修改前
import "github.com/zhoudm1743/go-flow/core/event"
var Module = fx.Options(
    // ...
    event.Module,
    // ...
)

// boot/bootstrap.go - 修改后
import "github.com/zhoudm1743/go-flow/core/domain"
var Module = fx.Options(
    // ...
    domain.Module,
    // ...
)
```

### 3. 创建DDD领域模块
```go
// core/domain/module.go - 新建
var Module = fx.Options(
    // 提供异步领域事件分发器
    fx.Provide(fx.Annotate(
        NewDefaultAsyncDomainEventDispatcher,
        fx.As(new(DomainEventDispatcher)),
    )),
    
    // 提供工作单元
    fx.Provide(func(db database.Database, dispatcher DomainEventDispatcher) UnitOfWork {
        return NewSimpleUnitOfWork(db.GetDB(), dispatcher)
    }),
    
    // 领域服务注册表
    fx.Provide(NewDomainServiceRegistry),
    
    // 系统事件处理器
    fx.Provide(NewSystemDomainEventHandlers),
    
    // 初始化
    fx.Invoke(InitDomainModule),
)
```

## 🏗️ DDD架构实现

### 1. 领域事件系统
```go
// 异步事件分发器 (10个工作协程，1000个事件缓冲)
dispatcher := NewAsyncDomainEventDispatcher(10, 1000)

// 自动注册系统事件处理器
dispatcher.Register(EventTypeEntityCreated, handlers.EntityCreatedHandler())
dispatcher.Register(EventTypeEntityUpdated, handlers.EntityUpdatedHandler())
dispatcher.Register(EventTypeEntityDeleted, handlers.EntityDeletedHandler())
```

### 2. 工作单元模式
```go
// 事务管理 + 自动事件分发
uow := NewSimpleUnitOfWork(db, dispatcher)
uow.RegisterNew(user)
uow.RegisterDirty(order)
uow.Commit() // 自动分发所有领域事件
```

### 3. 仓储模式
```go
// 泛型仓储，自动处理领域事件
type Repository[T AggregateRoot] interface {
    Save(aggregate T) error  // 自动添加创建/更新事件
    FindByID(id string) (T, error)
    FindBySpec(spec Specification[T]) ([]T, error)
    InTransaction(fn func(repo Repository[T]) error) error
}
```

### 4. 规约模式
```go
// 复合查询条件
activeSpec := NewSpecification(
    func(user *UserAggregate) bool { return user.Status == 1 },
    "status = ?", 1,
)

adminRoleSpec := NewSpecification(
    func(user *UserAggregate) bool { return contains(user.Roles, "admin") },
    "JSON_CONTAINS(roles, ?)", `"admin"`,
)

// 组合查询
users, err := repo.FindBySpec(activeSpec.And(adminRoleSpec))
```

## 📝 DDD使用示例

### 1. 聚合根定义
```go
// app/domain/user_aggregate.go
type UserAggregate struct {
    domain.BaseAggregateRoot
    Username    string      `gorm:"uniqueIndex;size:100;not null" json:"username"`
    Password    string      `gorm:"size:255;not null" json:"-"`
    Status      int8        `gorm:"default:1" json:"status"`
    Email       string      `gorm:"uniqueIndex;size:200;not null" json:"email"`
    Phone       types.Phone `gorm:"uniqueIndex;size:20;not null" json:"phone"`
    // ...
}

// 业务方法 - 自动触发领域事件
func (u *UserAggregate) ChangePassword(newPassword string) error {
    u.Password = newPassword
    
    // 添加密码修改事件
    u.AddDomainEvent(&PasswordChangedEvent{
        BaseDomainEvent: *domain.NewBaseDomainEvent(
            "user.password_changed",
            u.GetAggregateID(),
            u.GetAggregateType(),
            map[string]interface{}{
                "username": u.Username,
                "changed_at": time.Now(),
            },
        ),
        OldPasswordHash: oldPassword,
        NewPasswordHash: newPassword,
    })
    
    return nil
}
```

### 2. 仓储使用
```go
// app/domain/user_repository.go
type UserRepository interface {
    domain.Repository[*UserAggregate]
    FindByUsername(username string) (*UserAggregate, error)
    FindByEmail(email string) (*UserAggregate, error)
    FindActiveUsers() ([]*UserAggregate, error)
}

// 规约查询示例
func (r *userRepository) FindActiveUsers() ([]*UserAggregate, error) {
    spec := domain.NewSpecification[*UserAggregate](
        func(user *UserAggregate) bool { return user.Status == 1 },
        "status = ?", 1,
    )
    return r.FindBySpec(spec)
}
```

### 3. 事件处理器
```go
// 系统自动注册的事件处理器
type entityCreatedHandler struct {
    logger logger.Logger
}

func (h *entityCreatedHandler) Handle(event DomainEvent) error {
    h.logger.WithFields(map[string]interface{}{
        "event_id":       event.EventID(),
        "event_type":     event.EventType(),
        "aggregate_id":   event.AggregateID(),
        "aggregate_type": event.AggregateType(),
        "occurred_on":    event.OccurredOn(),
    }).Info("处理实体创建事件")
    return nil
}
```

## ⚡ 性能特性

### 1. 异步事件处理
- **工作协程**: 10个并发worker
- **缓冲队列**: 1000个事件缓冲
- **错误处理**: 异步处理，不阻塞主流程

### 2. 乐观锁并发控制
```go
// 自动版本控制
type BaseEntity struct {
    version int64 `json:"version" gorm:"column:version"`
    // ...
}

// 更新时自动检查版本冲突
result := db.Model(aggregate).
    Where("version = ?", originalVersion).
    Updates(aggregate)
    
if result.RowsAffected == 0 {
    return fmt.Errorf("optimistic lock failure")
}
```

### 3. 事务管理
```go
// 工作单元自动事务管理
func (uow *SimpleUnitOfWork) Commit() error {
    return uow.db.Transaction(func(tx *gorm.DB) error {
        // 批量处理实体
        // 批量分发事件
        return nil
    })
}
```

## 🎯 核心优势

### 1. **统一事件系统**
- 只保留DDD领域事件，避免系统复杂性
- 事件更贴近业务领域，语义更清晰

### 2. **自动化程度高**
- 仓储自动触发事件（创建、更新、删除）
- 工作单元自动事务管理和事件分发
- 乐观锁自动版本控制

### 3. **扩展性强**
- 支持自定义领域事件
- 支持复杂的规约组合查询
- 支持领域服务注册和管理

### 4. **性能优化**
- 异步事件处理，不阻塞业务流程
- 批量事件分发，提高效率
- 事务级别的事件一致性

## 🔄 如何使用新系统

### 1. 创建聚合根
```go
user := NewUserAggregate("testuser", "test@example.com", "13800138000")
```

### 2. 业务操作
```go
err := user.ChangePassword("newpassword")
err = user.Login("192.168.1.100")
```

### 3. 保存并触发事件
```go
repo := NewUserRepository(db)
err := repo.Save(user) // 自动分发领域事件
```

### 4. 复杂查询
```go
specs := UserSpecifications{}
activeAdmins, err := repo.FindBySpec(
    specs.ActiveUsers().And(specs.EmailDomain("@admin.com")),
)
```

## 📈 升级影响

### ✅ **保持兼容**
- 现有的传统模型（如 `SystemAdmin`）继续正常工作
- 可以渐进式地迁移到DDD模式
- 不影响现有API的功能

### 🚀 **新增能力**
- 完整的DDD基础设施
- 自动化的领域事件处理
- 规约模式的复杂查询
- 工作单元的事务管理

### 📚 **学习价值**
- 提供了完整的DDD实现参考
- 展示了如何设计领域驱动的系统
- 为团队学习DDD提供了实践基础

## 🎯 总结

通过这次迁移，框架成功地：

1. **简化了架构** - 从双事件系统统一为DDD领域事件
2. **提升了一致性** - 事件更贴近业务领域语义
3. **增强了可维护性** - 减少了系统复杂度
4. **保持了扩展性** - 为未来的DDD架构升级奠定基础

现在框架提供了完整的DDD基础设施，可以支持复杂的企业级业务场景，同时保持了简单易用的特性。 