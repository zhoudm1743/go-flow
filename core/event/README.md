# Event 事件系统

Event事件系统是一个功能完整的事件驱动架构组件，支持即时事件发布/订阅和延时事件调度。该系统与Cron服务深度集成，可以实现复杂的延时任务调度。

## 功能特性

### 📡 事件发布/订阅
- **即时事件**: 支持同步和异步事件发布
- **事件路由**: 基于事件类型的智能路由
- **并发处理**: 多工作协程并发处理事件
- **重试机制**: 可配置的重试策略和超时控制

### ⏰ 延时事件调度
- **精确调度**: 支持秒级精度的延时事件
- **优先级控制**: 四级优先级（低、普通、高、关键）
- **重试策略**: 可配置的最大重试次数和重试间隔
- **状态跟踪**: 完整的事件状态生命周期管理

### 🔗 Cron集成
- **双重调度**: 短延时用轮询器，长延时用Cron任务
- **任务联动**: Cron任务状态变化自动发布事件
- **失败恢复**: 自动重试和错误处理

### 💾 持久化存储
- **Redis存储**: 基于Redis的高性能事件存储
- **事件历史**: 完整的事件处理历史记录
- **指标统计**: 实时的事件处理统计信息

## 核心组件

### 1. 事件总线 (EventBus)
负责事件的发布、订阅和分发：

```go
// 异步发布事件
err := eventBus.PublishAsync(ctx, event)

// 订阅事件
err := eventBus.Subscribe("user.login", handler, options)
```

### 2. 事件存储 (EventStore)
处理事件的持久化和查询：

```go
// 存储延时事件
err := eventStore.StoreDelayedEvent(ctx, delayedEvent)

// 获取准备就绪的事件
events, err := eventStore.GetReadyEvents(ctx, 50)
```

### 3. 延时事件调度器 (DelayedEventScheduler)
管理延时事件的调度和执行：

```go
// 调度延时事件
err := scheduler.ScheduleDelayedEvent(ctx, event)

// 取消延时事件
err := scheduler.CancelDelayedEvent(ctx, eventID)
```

### 4. 事件服务 (EventService)
提供高级的事件管理API：

```go
// 发布即时事件
err := eventService.PublishEventAsync(ctx, "user.login", payload, "auth-service")

// 调度延时事件
delayedEvent, err := eventService.ScheduleDelayedEvent(ctx, "user.notification", payload, "system", delayUntil, options)
```

## 事件类型

### 系统内置事件

#### Cron相关事件
- `cron.task.created` - Cron任务创建
- `cron.task.updated` - Cron任务更新  
- `cron.task.deleted` - Cron任务删除
- `cron.task.executed` - Cron任务执行
- `cron.task.failed` - Cron任务失败
- `cron.task.completed` - Cron任务完成

#### 系统事件
- `system.startup` - 系统启动
- `system.shutdown` - 系统关闭
- `system.error` - 系统错误
- `system.warning` - 系统警告

#### 用户事件
- `user.login` - 用户登录
- `user.logout` - 用户登出
- `user.registered` - 用户注册

#### 数据事件
- `data.created` - 数据创建
- `data.updated` - 数据更新
- `data.deleted` - 数据删除

## 使用示例

### 1. 发布即时事件

```go
func PublishUserLoginEvent(eventService event.EventService, userID string) error {
    ctx := context.Background()
    
    payload := map[string]interface{}{
        "user_id": userID,
        "ip": "192.168.1.100",
        "timestamp": time.Now(),
        "device": "web",
    }
    
    return eventService.PublishEventAsync(ctx, event.EventTypeUserLogin, payload, "auth-service")
}
```

### 2. 调度延时事件

```go
func ScheduleWelcomeNotification(eventService event.EventService, userID string) (*event.DelayedEvent, error) {
    ctx := context.Background()
    
    payload := map[string]interface{}{
        "user_id": userID,
        "title": "欢迎使用我们的服务",
        "content": "感谢您注册我们的服务，祝您使用愉快！",
        "type": "welcome",
    }
    
    options := &event.DelayedEventOptions{
        Priority:   event.EventPriorityHigh,
        MaxRetries: 3,
    }
    
    // 10分钟后发送欢迎通知
    delayUntil := time.Now().Add(10 * time.Minute)
    return eventService.ScheduleDelayedEvent(ctx, "user.notification", payload, "notification-service", delayUntil, options)
}
```

### 3. 创建事件处理器

```go
type NotificationHandler struct {
    logger *zap.Logger
}

func (h *NotificationHandler) Handle(ctx context.Context, event event.Event) error {
    payload := event.GetPayload().(map[string]interface{})
    
    userID := payload["user_id"].(string)
    title := payload["title"].(string)
    content := payload["content"].(string)
    
    h.logger.Info("Sending notification",
        zap.String("user_id", userID),
        zap.String("title", title))
    
    // 实际发送通知的逻辑
    return sendNotification(userID, title, content)
}

func (h *NotificationHandler) GetEventType() string {
    return "user.notification"
}

func (h *NotificationHandler) GetName() string {
    return "NotificationHandler"
}

// 注册处理器
func RegisterNotificationHandler(eventService event.EventService) error {
    handler := &NotificationHandler{logger: zap.NewNop()}
    return eventService.Subscribe("user.notification", handler)
}
```

### 4. 数据清理任务

```go
func ScheduleDataCleanup(eventService event.EventService) (*event.DelayedEvent, error) {
    ctx := context.Background()
    
    payload := map[string]interface{}{
        "table_name": "temp_uploads",
        "conditions": map[string]interface{}{
            "created_at": map[string]interface{}{
                "$lt": time.Now().Add(-24 * time.Hour),
            },
        },
        "action": "delete",
    }
    
    options := &event.DelayedEventOptions{
        Priority:   event.EventPriorityLow,
        MaxRetries: 2,
    }
    
    // 24小时后清理临时文件
    delayUntil := time.Now().Add(24 * time.Hour)
    return eventService.ScheduleDelayedEvent(ctx, event.EventTypeDataDeleted, payload, "cleanup-service", delayUntil, options)
}
```

### 5. 工作流步骤调度

```go
func ScheduleWorkflowStep(eventService event.EventService, workflowID, stepID string, delay time.Duration) (*event.DelayedEvent, error) {
    ctx := context.Background()
    
    payload := map[string]interface{}{
        "workflow_id": workflowID,
        "step_id": stepID,
        "action": "execute",
        "input_data": map[string]interface{}{
            "user_id": "12345",
            "order_id": "order_789",
        },
    }
    
    options := &event.DelayedEventOptions{
        Priority:   event.EventPriorityNormal,
        MaxRetries: 3,
    }
    
    delayUntil := time.Now().Add(delay)
    return eventService.ScheduleDelayedEvent(ctx, "workflow.step.execute", payload, "workflow-engine", delayUntil, options)
}
```

## 配置选项

### 事件总线配置

```go
options := &event.EventBusOptions{
    BufferSize:    1000,        // 事件缓冲区大小
    WorkerCount:   10,          // 工作协程数量
    BatchSize:     50,          // 批处理大小
    FlushInterval: time.Second, // 刷新间隔
}
```

### 调度器配置

```go
options := &event.SchedulerOptions{
    PollInterval:    30 * time.Second, // 轮询间隔
    BatchSize:       50,               // 批处理大小
    MaxRetries:      3,                // 最大重试次数
    RetryInterval:   5 * time.Minute,  // 重试间隔
    CleanupInterval: time.Hour,        // 清理间隔
    ExpireAfter:     24 * time.Hour,   // 事件过期时间
}
```

### 订阅选项

```go
options := &event.SubscriptionOptions{
    Async:      true,              // 是否异步处理
    MaxRetries: 3,                 // 最大重试次数
    RetryDelay: 5 * time.Second,   // 重试延迟
    Timeout:    30 * time.Second,  // 处理超时
    Persistent: false,             // 是否持久化订阅
}
```

## 事件状态

- `pending` - 待处理
- `delivered` - 已投递
- `failed` - 处理失败
- `expired` - 已过期

## 事件优先级

- `1` - 低优先级
- `5` - 普通优先级
- `8` - 高优先级
- `10` - 关键优先级

## 监控和指标

### 获取事件指标

```go
metrics, err := eventService.GetMetrics(ctx)
if err != nil {
    return err
}

fmt.Printf("总事件数: %d\n", metrics.TotalEvents)
fmt.Printf("待处理事件: %d\n", metrics.PendingEvents)
fmt.Printf("已投递事件: %d\n", metrics.DeliveredEvents)
fmt.Printf("失败事件: %d\n", metrics.FailedEvents)
```

### 查询事件历史

```go
filter := &event.EventFilter{
    Types:     []string{"user.login", "user.logout"},
    StartTime: &startTime,
    EndTime:   &endTime,
    Page:      1,
    PageSize:  20,
}

records, err := eventService.GetEventRecords(ctx, filter)
if err != nil {
    return err
}

for _, record := range records.Items {
    fmt.Printf("事件ID: %s, 类型: %s, 状态: %s\n", 
        record.ID, record.Type, record.Status)
}
```

### 查询延时事件

```go
filter := &event.EventFilter{
    Status:   event.EventStatusPending,
    Page:     1,
    PageSize: 50,
}

events, err := eventService.GetDelayedEvents(ctx, filter)
if err != nil {
    return err
}

for _, event := range events.Items {
    fmt.Printf("事件ID: %s, 延时到: %s, 优先级: %d\n", 
        event.ID, event.DelayUntil.Format("2006-01-02 15:04:05"), event.Priority)
}
```

## 最佳实践

### 1. 事件命名规范
- 使用点分隔的层次结构：`service.entity.action`
- 示例：`user.profile.updated`、`order.payment.completed`

### 2. 负载数据设计
- 保持负载数据简洁，避免过大的对象
- 使用结构化的数据格式
- 包含必要的上下文信息

### 3. 错误处理
- 实现幂等的事件处理器
- 合理设置重试次数和重试间隔
- 记录详细的错误日志

### 4. 性能优化
- 合理设置事件总线的工作协程数量
- 使用异步处理减少阻塞
- 定期清理过期的事件记录

### 5. 监控告警
- 监控事件处理延迟
- 设置失败事件的告警
- 跟踪事件处理的成功率

## 故障排查

### 1. 事件未被处理
- 检查事件类型是否正确
- 确认处理器已正确注册
- 查看事件总线是否正常运行

### 2. 延时事件不执行
- 检查延时时间设置是否正确
- 确认调度器是否正常运行
- 查看Redis连接状态

### 3. 性能问题
- 检查工作协程数量配置
- 监控Redis性能
- 优化事件处理器逻辑

### 4. 内存泄漏
- 检查事件订阅是否正确取消
- 确认事件缓冲区大小合理
- 定期清理过期事件

## 与Cron服务集成

事件系统与Cron服务深度集成，实现了以下功能：

### 1. Cron任务事件
- Cron任务的生命周期事件自动发布到事件系统
- 支持任务执行结果的事件通知

### 2. 延时任务调度
- 长延时事件通过Cron任务精确调度
- 短延时事件通过轮询器快速处理

### 3. 失败重试
- 失败的延时事件可以重新调度
- 与Cron任务的重试机制联动

这个事件系统为Go应用提供了完整的事件驱动架构支持，可以轻松实现复杂的业务逻辑解耦和异步处理。 