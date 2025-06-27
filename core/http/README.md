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