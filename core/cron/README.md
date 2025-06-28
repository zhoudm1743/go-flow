# Cron调度服务

## 概述

这是一个基于Redis持久化的完整cron调度服务，支持系统任务和自定义HTTP任务的调度执行。

## 功能特性

- 🕒 **标准Cron表达式支持** - 支持秒级精度的cron表达式
- 🔄 **任务持久化** - 基于Redis的任务数据持久化存储
- 📊 **执行历史** - 完整的任务执行历史记录和统计
- 🔧 **系统任务** - 内置系统任务处理器，支持扩展
- 🌐 **HTTP任务** - 支持自定义HTTP请求任务
- ⚡ **高性能** - 基于goroutine的并发执行
- 🛡️ **容错机制** - 支持重试、超时控制
- 📈 **监控统计** - 提供详细的执行统计和状态监控

## 架构设计

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   Service层     │    │   Scheduler层   │    │  Repository层   │
│  业务逻辑接口   │───▶│   任务调度器    │───▶│   数据持久化    │
└─────────────────┘    └─────────────────┘    └─────────────────┘
          │                       │                       │
          ▼                       ▼                       ▼
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   Task实现      │    │   Cron引擎      │    │   Redis存储     │
│ HTTP/System任务 │    │ robfig/cron/v3  │    │  任务+历史数据  │
└─────────────────┘    └─────────────────┘    └─────────────────┘
```

## 任务类型

### 1. HTTP任务

支持所有标准HTTP方法的自定义任务：

```json
{
    "name": "API健康检查",
    "description": "定期检查API服务状态",
    "cron": "0 */5 * * * *",
    "type": "http",
    "config": {
        "url": "https://api.example.com/health",
        "method": "GET",
        "headers": {
            "Authorization": "Bearer your-token",
            "Content-Type": "application/json"
        },
        "body": "",
        "timeout": 30000000000,
        "retry_count": 3,
        "retry_delay": 5000000000,
        "expected_code": 200
    }
}
```

**HTTP任务特性：**
- 支持GET、POST、PUT、DELETE等所有HTTP方法
- 自动Content-Type检测（JSON/表单）
- 灵活的请求头配置
- 期望状态码验证
- 重试机制和超时控制

### 2. 系统任务

内置系统任务处理器，支持扩展：

#### 日志处理器 (`log`)
```json
{
    "handler_name": "log",
    "parameters": {
        "level": "info",
        "message": "系统状态检查完成"
    }
}
```

#### 缓存清理处理器 (`cache_clean`)
```json
{
    "handler_name": "cache_clean",
    "parameters": {
        "pattern": "temp:*"
    }
}
```

#### 数据库维护处理器 (`db_maintenance`)
```json
{
    "handler_name": "db_maintenance",
    "parameters": {
        "action": "analyze"
    }
}
```

## Cron表达式格式

支持6位或7位cron表达式（包含秒）：

```
┌───────────── 秒 (0-59)
│ ┌─────────── 分 (0-59)
│ │ ┌───────── 时 (0-23)
│ │ │ ┌─────── 日 (1-31)
│ │ │ │ ┌───── 月 (1-12)
│ │ │ │ │ ┌─── 周 (0-6, 周日=0)
│ │ │ │ │ │
* * * * * *
```

**常用示例：**
- `0 */5 * * * *` - 每5分钟执行
- `0 0 2 * * *` - 每天凌晨2点执行
- `0 0 9 * * 1-5` - 工作日上午9点执行
- `0 30 8 1 * *` - 每月1日上午8点30分执行

## 任务状态

- `active` - 激活状态，正在调度执行
- `paused` - 暂停状态，暂时停止调度
- `stopped` - 停止状态，永久停止

## 执行状态

- `running` - 执行中
- `success` - 执行成功
- `failed` - 执行失败
- `timeout` - 执行超时

## Redis存储结构

```
cron:task:{task_id}        # 任务详细信息
cron:tasks                 # 任务ID集合
cron:stats:{task_id}       # 任务统计信息
cron:history:{result_id}   # 执行结果详情
cron:history_list          # 执行历史有序集合
```

## API使用示例

### 创建HTTP任务

```go
// 创建HTTP任务
req := &CreateTaskRequest{
    Name:        "定期数据同步",
    Description: "每小时同步一次用户数据",
    Cron:        "0 0 * * * *",
    Type:        TaskTypeHTTP,
    Config: HTTPConfig{
        URL:          "https://api.example.com/sync",
        Method:       "POST",
        Headers:      map[string]string{"Authorization": "Bearer token"},
        Body:         `{"sync_type": "users"}`,
        Timeout:      30 * time.Second,
        RetryCount:   3,
        RetryDelay:   5 * time.Second,
        ExpectedCode: 200,
    },
}

taskInfo, err := cronService.CreateTask(ctx, req)
```

### 控制任务

```go
// 启动任务
err := cronService.StartTask(ctx, taskID)

// 暂停任务
err := cronService.PauseTask(ctx, taskID)

// 停止任务
err := cronService.StopTask(ctx, taskID)

// 立即执行
err := cronService.ExecuteTaskNow(ctx, taskID)
```

### 查询任务

```go
// 获取任务列表
req := &TaskListRequest{
    Page:     1,
    PageSize: 10,
    Status:   TaskStatusActive,
}
result, err := cronService.ListTasks(ctx, req)

// 获取执行历史
historyReq := &TaskExecHistoryRequest{
    TaskID:   taskID,
    Page:     1,
    PageSize: 20,
}
history, err := cronService.GetTaskExecHistory(ctx, historyReq)
```

## 自定义系统任务处理器

```go
type CustomHandler struct {
    logger logger.Logger
}

func (h *CustomHandler) Handle(ctx context.Context, params map[string]interface{}) (string, error) {
    // 实现自定义逻辑
    return "执行完成", nil
}

func (h *CustomHandler) GetName() string {
    return "custom_handler"
}

func (h *CustomHandler) GetDescription() string {
    return "自定义任务处理器"
}

// 注册处理器
RegisterSystemTaskHandler(&CustomHandler{logger: log})
```

## 性能监控

服务提供详细的执行统计：

```go
// 获取调度器状态
status, err := cronService.GetSchedulerStatus(ctx)
// 输出：运行状态、启动时间、任务数量、活跃任务数

// 获取任务统计
stats, err := repository.GetTaskStats(ctx, taskID)
// 输出：执行次数、成功次数、失败次数、平均执行时间
```

## 最佳实践

### 1. 任务设计
- 保持任务幂等性
- 避免长时间运行的任务
- 设置合理的超时时间
- 使用重试机制处理临时失败

### 2. 性能优化
- 合理设置任务执行频率
- 避免同时执行大量任务
- 定期清理执行历史
- 监控Redis内存使用

### 3. 错误处理
- 设置合适的重试次数和延迟
- 记录详细的错误信息
- 监控任务失败率
- 设置告警机制

### 4. 安全考虑
- HTTP任务使用HTTPS
- 敏感信息不要记录在日志中
- 限制任务执行权限
- 定期审计任务配置

## 故障排查

### 常见问题

1. **任务不执行**
   - 检查cron表达式是否正确
   - 确认任务状态为`active`
   - 查看调度器是否正常运行

2. **HTTP任务失败**
   - 检查URL是否可访问
   - 验证请求头和认证信息
   - 确认期望状态码设置

3. **Redis连接问题**
   - 检查Redis服务状态
   - 验证连接配置
   - 查看网络连通性

### 日志级别

- `DEBUG` - 详细的调试信息
- `INFO` - 正常的运行信息
- `WARN` - 警告信息（如重试）
- `ERROR` - 错误信息

## 扩展开发

### 添加新的任务类型

1. 实现`Task`接口
2. 创建对应的配置结构
3. 在`Repository`中添加序列化/反序列化逻辑
4. 在`Service`中添加创建逻辑

### 添加新的系统任务处理器

1. 实现`SystemTaskHandler`接口
2. 在模块初始化时注册
3. 提供配置文档和使用示例

## 配置示例

```yaml
# config.yaml
redis:
  host: localhost
  port: 6379
  password: ""
  db: 0
```

## 依赖项

- `github.com/robfig/cron/v3` - Cron表达式解析和调度
- `github.com/redis/go-redis/v9` - Redis客户端
- `github.com/google/uuid` - UUID生成
- `go.uber.org/fx` - 依赖注入框架 