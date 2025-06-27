# Go-Flow

一个基于 Go 语言的现代化微服务框架，集成了日志、数据库、缓存等核心功能。

## 🚀 功能特性

### 🎯 核心服务

- **🔧 配置管理** - 基于 Viper 的配置服务，支持 YAML 配置文件
- **📝 日志服务** - 基于 Logrus 的结构化日志，支持文件轮转和颜色输出
- **💾 数据库服务** - 基于 GORM 的 ORM，支持 MySQL，包含自动迁移
- **🚄 缓存服务** - 基于 Redis 的缓存服务，封装常用操作（**全新 API 设计**）
- **🔄 依赖注入** - 基于 Fx 的依赖注入框架

### 📋 已实现功能

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
├── boot/                    # 启动模块
│   └── bootstrap.go         # 应用启动配置
├── core/                    # 核心服务
│   ├── config/              # 配置服务
│   │   └── config.go
│   ├── logger/              # 日志服务
│   │   ├── logger.go        # 主日志服务
│   │   └── fx_adapter.go    # Fx 日志适配器
│   ├── database/            # 数据库服务
│   │   ├── gorm.go          # GORM 配置
│   │   ├── models.go        # 数据模型
│   │   ├── migrator.go      # 数据库迁移
│   │   └── service.go       # 数据库服务示例
│   └── cache/               # 🔥 缓存服务 - 全新设计
│       ├── redis.go         # Redis 客户端（双API设计）
│       ├── helper.go        # 缓存助手（高级功能）
│       └── service.go       # 缓存服务示例
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
INFO   [2025-06-28 00:30:33] 开始数据库自动迁移
INFO   [2025-06-28 00:30:35] 数据库迁移完成
INFO   [2025-06-28 00:30:35] 应用启动
INFO   [2025-06-28 00:30:35] 数据库连接测试成功
INFO   [2025-06-28 00:30:35] Redis 连接测试成功
```

## 💡 使用示例

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

## 🌟 新缓存 API 设计亮点

### 双 API 设计理念

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

- ✅ **简洁性**: 大部分场景下无需关心 context
- ✅ **灵活性**: 需要超时控制时可使用 Ctx 方法
- ✅ **兼容性**: 满足不同使用场景的需求
- ✅ **一致性**: 所有方法都遵循相同的命名规范
- ✅ **功能性**: 支持 Redis 所有数据类型操作

## 📦 依赖

- **Fx** - 依赖注入框架
- **Viper** - 配置管理
- **Logrus** - 结构化日志
- **Lumberjack** - 日志轮转
- **GORM** - ORM 框架
- **Redis Go Client v9** - Redis 客户端

## 🔄 扩展

该框架设计为模块化，可以轻松添加新的服务：

1. 在 `core/` 下创建新服务包
2. 实现服务接口和 Fx 模块
3. 在 `boot/bootstrap.go` 中集成新模块
4. 在配置文件中添加相应配置

## �� 许可证

MIT License 