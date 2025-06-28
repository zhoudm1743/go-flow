# HTTP 服务模块

本模块提供了基于 Gin 框架的 HTTP 服务，集成了项目的 logger 作为 Gin 的日志输出。

## 特性

- ✅ Gin 框架集成
- ✅ 自定义 Logger 集成
- ✅ CORS 支持
- ✅ 请求日志记录
- ✅ 错误恢复中间件
- ✅ 健康检查端点
- ✅ 配置化的超时设置
- ✅ 优雅关闭

## 配置

在 `config/config.yaml` 中添加 HTTP 配置：

```yaml
app:
  name: "go-flow"
  version: "1.0.0"
  port: 8080
  env: "development"
  http:
    read_timeout: 30      # 读取超时时间（秒）
    write_timeout: 30     # 写入超时时间（秒）
    idle_timeout: 120     # 空闲超时时间（秒）
    enable_cors: true     # 是否启用 CORS
```

## 默认路由

- `GET /health` - 健康检查
- `GET /ping` - 简单的 ping 测试
- `GET /api/v1/status` - 应用状态信息

## 使用方式

HTTP 服务会自动通过 fx 依赖注入启动，无需手动配置。

### 添加自定义路由

可以通过获取 Gin 引擎来添加自定义路由：

```go
func registerCustomRoutes(service *http.Service) {
    engine := service.GetEngine()
    
    // 添加自定义路由
    engine.GET("/custom", func(c *gin.Context) {
        c.JSON(200, gin.H{"message": "custom route"})
    })
}
```

## 日志集成

所有的 HTTP 请求都会通过项目的统一 logger 进行记录，包括：

- 请求信息（IP、方法、路径、状态码、延迟等）
- 错误信息
- 恐慌恢复信息

日志格式和输出方式遵循项目的全局日志配置。

# HTTP 路由系统

这是一个基于 Gin 的模块化路由注册系统，支持按照声明式的方式定义和注册路由。

## 新的路由注册模式

### 1. 定义路由组

使用 `httpCore.NewGroup` 创建路由组：

```go
// 定义路由组
var UserGroup = httpCore.NewGroup("/users", newUserRoutes, regUserRoutes, httpCore.AuthMiddleware)
```

参数说明：
- `prefix`: 路由前缀，如 "/users"
- `constructor`: 构造函数（可为 nil）
- `register`: 注册函数，类型为 `func(*httpCore.BaseGroup) error`
- `middleware...`: 中间件列表

### 2. 创建注册函数

注册函数负责定义具体的路由，支持多个处理器：

```go
func regUserRoutes(rg *httpCore.BaseGroup) error {
    return httpCore.Reg(func() {
        // 单个处理器
        rg.GET("", getUserList)
        
        // 多个处理器：中间件 + 主处理器
        rg.GET("/:id", 
            validateIDMiddleware,    // 中间件1
            logRequestMiddleware,    // 中间件2  
            getUserByID,            // 主处理器
        )
        
        rg.POST("", 
            validateInputMiddleware, // 验证中间件
            createUser,             // 主处理器
        )
    })
}
```

### 3. 注册模块路由

使用路由注册器注册模块路由：

```go
func (r *AdminRouteRegistrator) RegisterRoutes(engine *gin.Engine) error {
    return httpCore.RegisterModuleRoutes(engine, "admin", []httpCore.Group{
        UserGroup,
        SystemGroup,
        TestGroup,
    })
}
```

### 4. 完整示例

```go
package routes

import (
    "net/http"
    "github.com/gin-gonic/gin"
    httpCore "github.com/zhoudm1743/go-flow/core/http"
)

// 定义路由组
var UserGroup = httpCore.NewGroup("/users", nil, regUserRoutes, httpCore.AuthMiddleware)

// 注册函数
func regUserRoutes(rg *httpCore.BaseGroup) error {
    return httpCore.Reg(func() {
        // 单个处理器
        rg.GET("", func(c *gin.Context) {
            c.JSON(http.StatusOK, gin.H{"message": "用户列表"})
        })
        
        // 多个处理器：中间件 + 主处理器
        rg.GET("/:id", 
            func(c *gin.Context) {
                // 中间件：验证ID
                id := c.Param("id")
                if id == "" {
                    c.JSON(http.StatusBadRequest, gin.H{"error": "ID不能为空"})
                    c.Abort()
                    return
                }
                c.Set("user_id", id)
                c.Next()
            },
            func(c *gin.Context) {
                // 主处理器
                userID := c.GetString("user_id")
                c.JSON(http.StatusOK, gin.H{
                    "message": "获取用户",
                    "user_id": userID,
                })
            },
        )
        
        rg.POST("", 
            func(c *gin.Context) {
                // 验证中间件
                var data map[string]interface{}
                if err := c.ShouldBindJSON(&data); err != nil {
                    c.JSON(http.StatusBadRequest, gin.H{"error": "数据格式错误"})
                    c.Abort()
                    return
                }
                c.Set("validated_data", data)
                c.Next()
            },
            func(c *gin.Context) {
                // 主处理器
                data := c.MustGet("validated_data")
                c.JSON(http.StatusOK, gin.H{
                    "message": "创建用户成功",
                    "data":    data,
                })
            },
        )
    })
}
```

### 5. 支持的 HTTP 方法

BaseGroup 支持以下 HTTP 方法，都支持可变参数：
- `GET(path, handlers...)`
- `POST(path, handlers...)`
- `PUT(path, handlers...)`
- `DELETE(path, handlers...)`
- `PATCH(path, handlers...)`
- `HEAD(path, handlers...)`
- `OPTIONS(path, handlers...)`

### 6. 中间件链

每个路由方法都支持传入多个处理器，形成中间件链：

```go
rg.GET("/protected", 
    authMiddleware,      // 认证中间件
    rateLimitMiddleware, // 限流中间件
    logMiddleware,       // 日志中间件
    mainHandler,         // 主处理器
)
```

处理器执行顺序：
1. 认证中间件 -> 调用 `c.Next()`
2. 限流中间件 -> 调用 `c.Next()`
3. 日志中间件 -> 调用 `c.Next()`
4. 主处理器 -> 返回响应

### 7. 路由组级别的中间件

除了路由级别的中间件链，还可以在路由组级别添加中间件：

```go
var AdminGroup = httpCore.NewGroup("/admin", nil, regAdminRoutes, 
    httpCore.AuthMiddleware,    // 组级别认证中间件
    httpCore.AdminMiddleware,   // 组级别管理员权限中间件
)
```

最终的中间件执行顺序：
1. 组级别中间件
2. 路由级别中间件链
3. 主处理器

### 8. 路由结构

使用这个系统，您的路由将按以下结构组织：

```
/admin/test/test          - 单个处理器
/admin/test/test-multi    - 多个处理器（中间件+主处理器）
/admin/users              - 用户管理（支持中间件链）
/admin/users/:id          - 单个用户操作（支持验证中间件）
```

## 新特性

✅ **可变参数处理器**: 每个路由方法支持 `handlers ...RouteHandler`  
✅ **中间件链**: 灵活的中间件组合  
✅ **路由级别中间件**: 精细化的权限控制  
✅ **组级别中间件**: 统一的模块级别中间件  
✅ **声明式路由定义**: 使用全局变量声明路由组  
✅ **模块化**: 每个模块独立管理路由  
✅ **类型安全**: 编译时路由检查  
✅ **简洁API**: 最少的样板代码  
✅ **FX集成**: 完美集成依赖注入 