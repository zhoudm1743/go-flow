# Go-Flow

一个基于 Go 语言的现代化微服务框架，集成了声明式路由、日志、数据库、缓存等核心功能。

## 🚀 功能特性

### 🎯 核心服务

- **🔧 配置管理** - 基于 Viper 的配置服务，支持 YAML 配置文件
- **📝 日志服务** - 基于 Logrus 的结构化日志，支持文件轮转和颜色输出
- **💾 数据库服务** - 基于 GORM 的 ORM，支持 MySQL，包含自动迁移
- **🚄 缓存服务** - 基于 Redis 的缓存服务，封装常用操作（**全新 API 设计**）
- **🌐 HTTP 路由** - 声明式路由系统，支持依赖注入和中间件链（**全新设计**）
- **🔄 依赖注入** - 基于 Fx 的依赖注入框架

### 📋 已实现功能

#### 🌐 HTTP 路由系统 (core/http) - 全新设计
- ✅ **声明式路由定义**: 简洁的路由组声明
- ✅ **依赖注入支持**: Fx 自动注入服务依赖
- ✅ **模块化设计**: 每个模块自管理路由
- ✅ **中间件链**: 支持可变参数处理器 `handlers ...RouteHandler`
- ✅ **函数式架构**: 零冗余的函数式设计
- ✅ **类型安全**: 编译时检查所有依赖关系

#### 配置服务 (core/config)
- ✅ YAML 配置文件支持
- ✅ 环境变量读取
- ✅ 配置结构化定义
- ✅ 默认配置自动生成

#### 日志服务 (core/logger)
- ✅ 多级别日志 (Debug, Info, Warn, Error, Fatal)
- ✅ 结构化日志 (带字段)
- ✅ 文件轮转 (Lumberjack)
- ✅ 彩色输出支持
- ✅ Fx 框架日志集成
- ✅ 全局日志函数

#### 数据库服务 (core/database)
- ✅ GORM 集成
- ✅ MySQL 驱动
- ✅ 连接池配置
- ✅ 自动迁移
- ✅ 示例模型 (User, Post)
- ✅ 基础 CRUD 服务

#### 🔥 缓存服务 (core/cache) - 全新设计
- ✅ **双 API 设计**：默认方法简洁，Ctx 后缀方法精细控制
- ✅ **基础操作**: Get, Set, Del, Exists, Expire, TTL
- ✅ **字符串操作**: Incr, Decr, IncrBy
- ✅ **哈希操作**: HGet, HSet, HGetAll, HDel, HExists, HLen
- ✅ **列表操作**: LPush, RPush, LPop, RPop, LLen, LRange
- ✅ **集合操作**: SAdd, SRem, SMembers, SIsMember, SCard
- ✅ **有序集合操作**: ZAdd, ZRem, ZRange, ZCard, ZScore
- ✅ **高级封装**: JSON 缓存, 分布式锁, 批量操作, 记忆模式

## 📁 项目结构

```
go-flow/
├── app/                     # 应用模块
│   └── admin/               # 管理模块
│       ├── routes/          # 路由定义
│       │   ├── enter.go     # 路由注册器
│       │   └── test/        # 测试路由模块
│       │       └── test.go  # 测试路由实现
│       ├── service/         # 业务服务
│       │   └── test/        # 测试服务
│       └── schemas/         # 数据结构
├── boot/                    # 启动模块
│   └── bootstrap.go         # 应用启动配置
├── core/                    # 核心服务
│   ├── config/              # 配置服务
│   │   └── config.go
│   ├── logger/              # 日志服务
│   │   ├── logger.go        # 主日志服务
│   │   └── fx_adapter.go    # Fx 日志适配器
│   ├── http/                # 🌐 HTTP 路由系统 - 全新设计
│   │   ├── router.go        # 声明式路由核心
│   │   ├── service.go       # HTTP 服务
│   │   ├── module.go        # Fx 模块定义
│   │   └── middleware.go    # 中间件支持
│   ├── database/            # 数据库服务
│   │   ├── gorm.go          # GORM 配置
│   │   ├── models.go        # 数据模型
│   │   ├── migrator.go      # 数据库迁移
│   │   ├── module.go        # Fx 模块定义
│   │   └── service.go       # 数据库服务示例
│   └── cache/               # 🔥 缓存服务 - 全新设计
│       ├── redis.go         # Redis 客户端（双API设计）
│       └── helper.go        # 缓存助手（高级功能）
├── config/                  # 配置文件
│   └── config.yaml          # 主配置文件
├── logs/                    # 日志文件目录
├── docker-compose.yml       # Docker 服务
├── go.mod                   # Go 模块定义
└── main.go                  # 主入口文件
```

## 🔧 配置说明

### config.yaml 配置文件结构

```yaml
app:
  name: "go-flow"
  version: "1.0.0"
  port: 8080
  env: "development"

database:
  host: "localhost"
  port: 3306
  username: "root"
  password: "123456"
  database: "go_flow"

redis:
  host: "localhost"
  port: 6379
  password: ""
  db: 0

log:
  level: "info"           # debug, info, warn, error
  format: "text"          # text, json
  output: "both"          # stdout, stderr, file, both
  lumberjack:
    filename: "logs/app.log"
    maxsize: 100          # MB
    maxage: 30            # 天
    maxbackups: 5         # 备份数量
    compress: true        # 是否压缩
```

## 🚦 快速开始

### 1. 启动依赖服务

```bash
# 使用 Docker Compose 启动 MySQL 和 Redis
docker-compose up -d
```

### 2. 运行应用

```bash
# 安装依赖
go mod download

# 运行应用
go run main.go
```

### 3. 验证服务

应用启动后会显示类似输出：

```
INFO   [2025-06-28 00:30:33] 数据库连接成功
INFO   [2025-06-28 00:30:33] Redis 连接成功
INFO   [2025-06-28 00:30:33] [GIN-debug] GET    /admin/test/test          
INFO   [2025-06-28 00:30:33] [GIN-debug] GET    /admin/test/test-multi    
INFO   [2025-06-28 00:30:33] [GIN-debug] POST   /admin/test/test          
INFO   [2025-06-28 00:30:33] HTTP 服务已启动，地址: http://localhost:8080
```

### 4. 测试API端点

```bash
# 测试基础路由
curl http://localhost:8080/admin/test/test

# 测试中间件链路由
curl http://localhost:8080/admin/test/test-multi

# 测试POST路由
curl -X POST http://localhost:8080/admin/test/test
```

## 💡 使用示例

### 🌐 声明式路由系统 - 全新设计

#### 🎯 路由模块定义（极简设计）

```go
// app/admin/routes/test/test.go
package test

import (
    "net/http"
    "github.com/gin-gonic/gin"
    "github.com/zhoudm1743/go-flow/app/admin/service/test"
    httpCore "github.com/zhoudm1743/go-flow/core/http"
)

type testRoutes struct {
    srv test.TestService
}

// NewTestGroup fx Provider函数，自动注入TestService并返回配置好的Group
func NewTestGroup(srv test.TestService) httpCore.Group {
    return httpCore.NewGroup("/test", 
        func() interface{} {
            return &testRoutes{srv: srv}
        }, 
        regTest,
    )
}

// regTest 注册测试路由（内部函数）
func regTest(rg *httpCore.BaseGroup, instance interface{}) error {
    r := instance.(*testRoutes)
    
    // 🔸 单个处理器
    rg.GET("/test", r.test)
    
    // 🔸 中间件链（多个处理器）
    rg.GET("/test-multi", 
        func(c *gin.Context) {
            c.Set("middleware_data", "from middleware")
            c.Next()
        },
        r.testMulti, // 主处理器
    )
    
    // 🔸 POST 路由
    rg.POST("/test", r.testPost)
    
    return nil
}

func (r *testRoutes) test(c *gin.Context) {
    res := r.srv.Test()
    c.JSON(http.StatusOK, res)
}
```

#### 🚀 路由注册器（函数式设计）

```go
// app/admin/routes/enter.go
package routes

import (
    "github.com/gin-gonic/gin"
    testRoutes "github.com/zhoudm1743/go-flow/app/admin/routes/test"
    httpCore "github.com/zhoudm1743/go-flow/core/http"
    "go.uber.org/fx"
)

// RouteRegistratorFunc 函数类型实现RouteRegistrator接口
type RouteRegistratorFunc func(*gin.Engine) error

func (f RouteRegistratorFunc) RegisterRoutes(engine *gin.Engine) error {
    return f(engine)
}

// NewAdminRouteRegistrator 创建admin路由注册器 - 简化为一个函数
func NewAdminRouteRegistrator(group httpCore.Group) RouteRegistratorResult {
    return RouteRegistratorResult{
        Registrator: RouteRegistratorFunc(func(engine *gin.Engine) error {
            return httpCore.RegisterModuleRoutes(engine, "admin", []httpCore.Group{
                group, // 🎉 终极简化！
            })
        }),
    }
}

// Module FX模块定义
var Module = fx.Options(
    fx.Provide(testRoutes.NewTestGroup),
    fx.Provide(NewAdminRouteRegistrator),
)
```

#### 🌟 路由系统特性

1. **🎯 零冗余设计**
   - 路由模块：一个函数 `NewTestGroup`
   - 路由注册器：一个函数类型 + 一个构造函数
   
2. **🔄 完全自动化**
   - Fx 自动注入 `TestService` 到路由模块
   - 路由注册器自动收集并注册所有路由

3. **⚡ 极致简洁**
   ```go
   // 添加新路由模块只需要：
   fx.Provide(newRoutes.NewUserGroup),  // 提供路由组
   ```

4. **🛡️ 类型安全**
   - 编译时检查所有依赖关系
   - 函数签名保证接口一致性

### 数据库操作

```go
// 在服务中注入数据库
func NewMyService(db database.Database, log logger.Logger) *MyService {
    return &MyService{db: db, logger: log}
}

// 使用 GORM
func (s *MyService) CreateUser(user *database.User) error {
    return s.db.GetDB().Create(user).Error
}
```

### 🔥 缓存操作 - 全新双 API 设计

```go
// 在服务中注入缓存
func NewMyService(cache cache.Cache, log logger.Logger) *MyService {
    return &MyService{cache: cache, logger: log}
}
```

#### 🌟 默认方法（简洁易用，无需 context）

```go
func (s *MyService) SimpleCacheData() error {
    // 🔸 基础操作
    err := s.cache.Set("user:name", "张三", time.Hour)
    value, err := s.cache.Get("user:name")
    
    // 🔸 哈希操作
    s.cache.HSet("user:1", "name", "李四", "age", "25", "city", "北京")
    name, err := s.cache.HGet("user:1", "name")
    allFields, err := s.cache.HGetAll("user:1")
    
    // 🔸 计数器操作
    count, err := s.cache.Incr("visit:count")
    count, err = s.cache.IncrBy("visit:count", 5)
    
    // 🔸 列表操作
    s.cache.RPush("messages", "消息1", "消息2", "消息3")
    message, err := s.cache.LPop("messages")
    length, err := s.cache.LLen("messages")
    
    // 🔸 集合操作
    s.cache.SAdd("tags", "golang", "redis", "cache")
    members, err := s.cache.SMembers("tags")
    exists, err := s.cache.SIsMember("tags", "golang")
    
    // 🔸 有序集合（排行榜）
    s.cache.ZAdd("leaderboard", 
        redis.Z{Score: 100, Member: "张三"},
        redis.Z{Score: 95, Member: "李四"},
    )
    topUsers, err := s.cache.ZRange("leaderboard", 0, 2)
    
    return err
}
```

#### 🎯 带 Context 方法（精细控制）

```go
func (s *MyService) CacheDataWithContext(ctx context.Context) error {
    // 使用超时控制的缓存操作
    err := s.cache.SetCtx(ctx, "key", "value", time.Hour)
    value, err := s.cache.GetCtx(ctx, "key")
    
    // 带 context 的哈希操作
    _, err = s.cache.HSetCtx(ctx, "hash", "field", "value")
    value, err = s.cache.HGetCtx(ctx, "hash", "field")
    
    // 批量操作
    keys, err := s.cache.KeysCtx(ctx, "user:*")
    deleted, err := s.cache.DelCtx(ctx, keys...)
    
    return err
}
```

### 🛠 高级缓存功能

#### JSON 缓存

```go
// 使用缓存助手
helper := cache.NewCacheHelper(cacheClient, logger, "myapp")

// 🔸 默认方法（简洁）
user := map[string]interface{}{"name": "张三", "age": 25}
err := helper.SetJSON("user:1", user, time.Hour)

var retrievedUser map[string]interface{}
err = helper.GetJSON("user:1", &retrievedUser)

// 🔸 带 Context 方法（精细控制）
err = helper.SetJSONCtx(ctx, "user:1", user, time.Hour)
err = helper.GetJSONCtx(ctx, "user:1", &retrievedUser)
```

#### 分布式锁

```go
// 🔸 简洁版
locked, err := helper.Lock("resource", time.Minute)
if locked {
    defer helper.Unlock("resource")
    // 执行需要锁保护的操作
}

// 🔸 带 Context 版
locked, err := helper.LockCtx(ctx, "resource", time.Minute)
if locked {
    defer helper.UnlockCtx(ctx, "resource")
    // 执行需要锁保护的操作
}

// 🔸 自动释放锁
err := helper.WithLock("resource", time.Minute, func() error {
    // 执行需要锁保护的操作
    return nil
})
```

#### 记忆模式（缓存函数结果）

```go
// 🔸 简洁版
result, err := helper.Remember("expensive_calc", time.Hour, func() (interface{}, error) {
    // 执行昂贵的计算
    return calculateSomething(), nil
})

// 🔸 JSON 版本
var result MyStruct
err = helper.RememberJSON("user_profile", time.Hour, &result, func() (interface{}, error) {
    return fetchUserFromDatabase(), nil
})
```

#### 批量操作

```go
// 🔸 批量设置
data := map[string]interface{}{
    "key1": "value1",
    "key2": "value2",
    "key3": "value3",
}
err := helper.BatchSet(data, time.Hour)

// 🔸 批量获取
results, err := helper.BatchGet([]string{"key1", "key2", "key3"})

// 🔸 模式删除
deletedCount, err := helper.FlushByPattern("temp:*")
```

### 日志使用

```go
// 在服务中使用
func (s *MyService) DoSomething() {
    s.logger.Info("开始执行操作")
    
    s.logger.WithFields(map[string]interface{}{
        "user_id": 123,
        "action": "create",
    }).Info("用户操作")
    
    s.logger.WithError(err).Error("操作失败")
}

// 全局使用
import "github.com/zhoudm1743/go-flow/core/logger"

logger.Info("全局日志信息")
logger.WithField("key", "value").Warn("全局警告")
```

## 🌟 设计亮点

### HTTP 路由系统

1. **🎯 声明式设计** - 路由定义简洁明了
2. **⚡ 零冗余架构** - 每行代码都有存在价值
3. **🔄 完全自动化** - Fx 框架处理所有依赖注入
4. **🛡️ 类型安全** - 编译时验证所有依赖关系
5. **📦 模块化** - 每个模块自管理路由和依赖

### 缓存双 API 设计

1. **🌟 默认方法** - 无需传递 `context`，使用简洁
   ```go
   cache.Set("key", "value", time.Hour)
   cache.Get("key")
   cache.HSet("hash", "field", "value")
   ```

2. **🎯 Ctx 后缀方法** - 需要精细控制时使用
   ```go
   cache.SetCtx(ctx, "key", "value", time.Hour)
   cache.GetCtx(ctx, "key")
   cache.HSetCtx(ctx, "hash", "field", "value")
   ```

### 设计优势

- ✅ **简洁性**: 大部分场景下无需关心复杂配置
- ✅ **灵活性**: 需要精细控制时提供完整功能
- ✅ **兼容性**: 满足不同使用场景的需求
- ✅ **一致性**: 所有API都遵循相同的设计理念
- ✅ **扩展性**: 模块化设计便于功能扩展

## 📦 依赖

- **Fx** - 依赖注入框架
- **Gin** - HTTP Web 框架
- **Viper** - 配置管理
- **Logrus** - 结构化日志
- **Lumberjack** - 日志轮转
- **GORM** - ORM 框架
- **Redis Go Client v9** - Redis 客户端

## 🔄 扩展

该框架设计为模块化，可以轻松添加新的功能：

### 添加新路由模块

1. 在 `app/admin/routes/` 下创建新模块目录
2. 实现 `NewXxxGroup` 函数
3. 在 `enter.go` 中添加 Provider

```go
// Module FX模块定义
var Module = fx.Options(
    fx.Provide(testRoutes.NewTestGroup),
    fx.Provide(userRoutes.NewUserGroup),  // 🆕 新路由模块
    fx.Provide(NewAdminRouteRegistrator),
)
```

### 添加新服务模块

1. 在 `core/` 下创建新服务包
2. 实现服务接口和 Fx 模块
3. 在 `boot/bootstrap.go` 中集成新模块
4. 在配置文件中添加相应配置

## 🎯 Casbin权限系统集成

### 概述
本项目已集成Casbin RBAC权限认证系统，提供完善的权限管理功能。

### 功能特性
- ✅ RBAC角色权限模型
- ✅ 用户角色分配管理
- ✅ 权限策略动态配置
- ✅ API权限中间件
- ✅ 数据库持久化存储
- ✅ 批量权限操作
- ✅ 默认权限策略初始化

### 权限模型
使用RBAC模型，支持：
- 用户(User) -> 角色(Role) -> 权限(Permission)
- 路径匹配权限控制（支持通配符）
- 角色继承机制

### API接口

#### 权限管理接口（需要admin角色）
```
POST   /api/admin/permissions/users/:userID/roles      # 为用户分配角色
DELETE /api/admin/permissions/users/:userID/roles/:role # 移除用户角色
GET    /api/admin/permissions/users/:userID/roles      # 获取用户角色

POST   /api/admin/permissions/policies                 # 添加权限策略
DELETE /api/admin/permissions/policies                 # 删除权限策略
GET    /api/admin/permissions/policies/:subject        # 获取主体权限

POST   /api/admin/permissions/roles                    # 创建角色
DELETE /api/admin/permissions/roles/:role              # 删除角色
GET    /api/admin/permissions/roles                    # 获取所有角色

POST   /api/admin/permissions/check                    # 检查权限
```

#### 受保护接口示例
```
GET    /api/admin/protected/profile                    # 获取用户信息
PUT    /api/admin/protected/profile                    # 更新用户信息
```

### 使用示例

#### 1. 分配用户角色
```bash
curl -X POST http://localhost:8080/api/admin/permissions/users/user123/roles \
  -H "Content-Type: application/json" \
  -H "X-User-ID: admin" \
  -d '{"role": "admin"}'
```

#### 2. 添加权限策略
```bash
curl -X POST http://localhost:8080/api/admin/permissions/policies \
  -H "Content-Type: application/json" \
  -H "X-User-ID: admin" \
  -d '{
    "subject": "user",
    "object": "/api/user/*",
    "action": "GET"
  }'
```

#### 3. 检查权限
```bash
curl -X POST http://localhost:8080/api/admin/permissions/check \
  -H "Content-Type: application/json" \
  -H "X-User-ID: admin" \
  -d '{
    "user_id": "user123",
    "resource": "/api/user/profile",
    "action": "GET"
  }'
```

#### 4. 访问受保护的API
```bash
curl -X GET http://localhost:8080/api/admin/protected/profile \
  -H "X-User-ID: user123"
```

### 默认角色和权限

#### 角色类型
- `admin`: 管理员，拥有所有权限
- `user`: 普通用户，拥有基础功能权限
- `guest`: 访客，只有公开接口权限

#### 默认权限策略
```
# 管理员权限
admin -> /api/* -> GET,POST,PUT,DELETE
admin -> /system/* -> GET,POST,PUT,DELETE

# 用户权限
user -> /api/user/* -> GET,POST,PUT
user -> /api/profile/* -> GET,PUT

# 访客权限
guest -> /api/public/* -> GET
guest -> /api/login -> POST
guest -> /api/register -> POST
```

### 中间件使用

#### 权限验证中间件
```go
// 使用Casbin权限验证
router.Use(http.CasbinAuthMiddleware(casbinService))

// 要求特定角色
router.Use(http.RequireRole("admin", casbinService))
```

### 配置说明

#### RBAC模型配置 (`config/rbac_model.conf`)
```conf
[request_definition]
r = sub, obj, act

[policy_definition]  
p = sub, obj, act

[role_definition]
g = _, _

[policy_effect]
e = some(where (p.eft == allow))

[matchers]
m = g(r.sub, p.sub) && keyMatch2(r.obj, p.obj) && r.act == p.act
```

### 开发指南

#### 在代码中使用
```go
// 注入CasbinService
func NewYourController(casbinService *casbin.CasbinService) *YourController {
    return &YourController{
        casbinService: casbinService,
    }
}

// 检查权限
func (c *YourController) SomeHandler(ctx *gin.Context) {
    userID := ctx.GetString("user_id")
    
    if !c.casbinService.CheckPermission(userID, "/api/resource", "GET") {
        ctx.JSON(403, gin.H{"error": "权限不足"})
        return
    }
    
    // 继续处理...
}
```

### 数据库表
权限数据存储在`casbin_rule`表中，包含以下字段：
- `ptype`: 策略类型 (p/g)
- `v0`: 主体 (用户/角色)
- `v1`: 对象 (资源路径)
- `v2`: 动作 (HTTP方法)

### 注意事项
1. 权限检查时使用`X-User-ID`请求头传递用户ID
2. 生产环境中应从JWT token中解析用户信息
3. 权限策略支持通配符匹配（使用`keyMatch2`）
4. 所有权限变更会自动持久化到数据库

## �� 许可证

MIT License 