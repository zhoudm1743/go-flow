package event

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/zhoudm1743/go-flow/core/logger"
	"github.com/zhoudm1743/go-flow/pkg/response"
	"go.uber.org/zap"
)

// EventStore 事件存储接口
type EventStore interface {
	// StoreDelayedEvent 存储延时事件
	StoreDelayedEvent(ctx context.Context, event *DelayedEvent) error
	// GetDelayedEvent 获取延时事件
	GetDelayedEvent(ctx context.Context, eventID string) (*DelayedEvent, error)
	// UpdateDelayedEvent 更新延时事件
	UpdateDelayedEvent(ctx context.Context, event *DelayedEvent) error
	// DeleteDelayedEvent 删除延时事件
	DeleteDelayedEvent(ctx context.Context, eventID string) error
	// GetReadyEvents 获取准备就绪的延时事件
	GetReadyEvents(ctx context.Context, limit int) ([]*DelayedEvent, error)
	// GetDelayedEvents 获取延时事件列表
	GetDelayedEvents(ctx context.Context, filter *EventFilter) (*response.PageResult[*DelayedEvent], error)
	// StoreEventRecord 存储事件记录
	StoreEventRecord(ctx context.Context, record *EventRecord) error
	// GetEventRecords 获取事件记录
	GetEventRecords(ctx context.Context, filter *EventFilter) (*response.PageResult[*EventRecord], error)
	// GetMetrics 获取事件指标
	GetMetrics(ctx context.Context) (*EventMetrics, error)
	// Cleanup 清理过期事件
	Cleanup(ctx context.Context, expireBefore time.Time) (int64, error)
}

// RedisEventStore Redis事件存储
type RedisEventStore struct {
	client redis.Cmdable
	logger logger.Logger
}

// Redis键前缀
const (
	keyDelayedEvents    = "event:delayed"      // 延时事件 Hash
	keyDelayedEventList = "event:delayed:list" // 延时事件列表 ZSet (按时间排序)
	keyEventRecords     = "event:records"      // 事件记录 Hash
	keyEventRecordList  = "event:records:list" // 事件记录列表 List
	keyEventMetrics     = "event:metrics"      // 事件指标 Hash
	keyEventSequence    = "event:sequence"     // 事件序列号
)

// NewRedisEventStore 创建Redis事件存储
func NewRedisEventStore(client redis.Cmdable, logger logger.Logger) *RedisEventStore {
	return &RedisEventStore{
		client: client,
		logger: logger,
	}
}

// StoreDelayedEvent 存储延时事件
func (s *RedisEventStore) StoreDelayedEvent(ctx context.Context, event *DelayedEvent) error {
	// 序列化事件
	data, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to marshal delayed event: %w", err)
	}

	pipe := s.client.Pipeline()

	// 存储事件数据
	eventKey := fmt.Sprintf("%s:%s", keyDelayedEvents, event.ID)
	pipe.HSet(ctx, eventKey, map[string]interface{}{
		"data":        string(data),
		"type":        event.Type,
		"source":      event.Source,
		"status":      string(event.Status),
		"priority":    int(event.Priority),
		"delay_until": event.DelayUntil.Unix(),
		"created_at":  event.CreatedAt.Unix(),
		"updated_at":  event.UpdatedAt.Unix(),
	})

	// 添加到延时事件列表 (按DelayUntil时间排序)
	score := float64(event.DelayUntil.Unix())
	pipe.ZAdd(ctx, keyDelayedEventList, redis.Z{
		Score:  score,
		Member: event.ID,
	})

	// 执行事务
	if _, err := pipe.Exec(ctx); err != nil {
		return fmt.Errorf("存储延时事件失败: %w", err)
	}

	s.logger.WithField("event_id", event.ID).
		WithField("event_type", event.Type).
		WithField("delay_until", event.DelayUntil).
		Debug("延时事件存储成功")

	return nil
}

// GetDelayedEvent 获取延时事件
func (s *RedisEventStore) GetDelayedEvent(ctx context.Context, eventID string) (*DelayedEvent, error) {
	data, err := s.client.HGet(ctx, fmt.Sprintf("%s:%s", keyDelayedEvents, eventID), "data").Result()
	if err != nil {
		if err == redis.Nil {
			return nil, fmt.Errorf("延时事件不存在: %s", eventID)
		}
		s.logger.WithField("event_id", eventID).
			WithError(err).Warn("获取延时事件失败")
		return nil, fmt.Errorf("获取延时事件失败: %w", err)
	}

	var event DelayedEvent
	if err := json.Unmarshal([]byte(data), &event); err != nil {
		return nil, fmt.Errorf("反序列化延时事件失败: %w", err)
	}

	return &event, nil
}

// UpdateDelayedEvent 更新延时事件
func (s *RedisEventStore) UpdateDelayedEvent(ctx context.Context, event *DelayedEvent) error {
	// 更新状态
	event.UpdatedAt = time.Now()

	// 检查事件是否存在
	eventKey := fmt.Sprintf("%s:%s", keyDelayedEvents, event.ID)
	exists, err := s.client.HExists(ctx, eventKey, "data").Result()
	if err != nil {
		return fmt.Errorf("检查事件存在性失败: %w", err)
	}
	if !exists {
		return fmt.Errorf("延时事件不存在: %s", event.ID)
	}

	// 序列化事件
	data, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to marshal delayed event: %w", err)
	}

	pipe := s.client.Pipeline()

	// 更新事件数据
	pipe.HSet(ctx, eventKey, map[string]interface{}{
		"data":        string(data),
		"status":      string(event.Status),
		"retry_count": event.RetryCount,
		"updated_at":  event.UpdatedAt.Unix(),
	})

	// 如果DelayUntil改变了，更新排序
	oldScore, err := s.client.ZScore(ctx, keyDelayedEventList, event.ID).Result()
	if err == nil {
		newScore := float64(event.DelayUntil.Unix())
		if oldScore != newScore {
			pipe.ZAdd(ctx, keyDelayedEventList, redis.Z{
				Score:  newScore,
				Member: event.ID,
			})
		}
	}

	// 执行事务
	if _, err := pipe.Exec(ctx); err != nil {
		return fmt.Errorf("更新延时事件失败: %w", err)
	}

	s.logger.WithField("event_id", event.ID).
		WithField("status", string(event.Status)).
		Debug("延时事件更新成功")

	return nil
}

// DeleteDelayedEvent 删除延时事件
func (s *RedisEventStore) DeleteDelayedEvent(ctx context.Context, eventID string) error {
	pipe := s.client.Pipeline()

	// 从Hash中删除事件数据
	pipe.HDel(ctx, fmt.Sprintf("%s:%s", keyDelayedEvents, eventID))

	// 从ZSet中删除时间索引
	pipe.ZRem(ctx, keyDelayedEventList, eventID)

	_, err := pipe.Exec(ctx)
	if err != nil {
		return fmt.Errorf("删除延时事件失败: %w", err)
	}

	s.logger.WithField("event_id", eventID).Debug("延时事件删除成功")

	return nil
}

// GetReadyEvents 获取准备就绪的事件
func (s *RedisEventStore) GetReadyEvents(ctx context.Context, limit int) ([]*DelayedEvent, error) {
	now := time.Now().UnixNano()

	// 获取到期的事件ID
	eventIDs, err := s.client.ZRangeByScore(ctx, keyDelayedEventList, &redis.ZRangeBy{
		Min:   "0",
		Max:   fmt.Sprintf("%d", now),
		Count: int64(limit),
	}).Result()
	if err != nil {
		return nil, fmt.Errorf("获取准备就绪的事件ID失败: %w", err)
	}

	if len(eventIDs) == 0 {
		return nil, nil
	}

	// 批量获取事件数据
	pipe := s.client.Pipeline()
	for _, eventID := range eventIDs {
		pipe.HGet(ctx, fmt.Sprintf("%s:%s", keyDelayedEvents, eventID), "data")
	}

	results, err := pipe.Exec(ctx)
	if err != nil {
		return nil, fmt.Errorf("获取准备就绪的事件失败: %w", err)
	}

	var events []*DelayedEvent
	for i, result := range results {
		cmd := result.(*redis.StringCmd)
		data, err := cmd.Result()
		if err != nil {
			if err == redis.Nil {
				continue // 事件可能已被删除
			}
			s.logger.WithField("event_id", eventIDs[i]).
				WithError(err).Warn("获取延时事件失败")
			continue
		}

		var event DelayedEvent
		if err := json.Unmarshal([]byte(data), &event); err != nil {
			s.logger.WithField("event_id", eventIDs[i]).
				WithError(err).Warn("反序列化延时事件失败")
			continue
		}

		// 只返回pending状态的事件
		if event.Status == EventStatusPending {
			events = append(events, &event)
		}
	}

	return events, nil
}

// GetDelayedEvents 获取延时事件列表
func (s *RedisEventStore) GetDelayedEvents(ctx context.Context, filter *EventFilter) (*response.PageResult[*DelayedEvent], error) {
	if filter.Page <= 0 {
		filter.Page = 1
	}
	if filter.PageSize <= 0 {
		filter.PageSize = 20
	}

	// 获取所有事件ID
	var eventIDs []string
	var err error

	if filter.StartTime != nil || filter.EndTime != nil {
		// 按时间范围查询
		min := "-inf"
		max := "+inf"
		if filter.StartTime != nil {
			min = strconv.FormatInt(filter.StartTime.Unix(), 10)
		}
		if filter.EndTime != nil {
			max = strconv.FormatInt(filter.EndTime.Unix(), 10)
		}

		eventIDs, err = s.client.ZRangeByScore(ctx, keyDelayedEventList, &redis.ZRangeBy{
			Min: min,
			Max: max,
		}).Result()
	} else {
		// 获取所有事件ID
		eventIDs, err = s.client.ZRange(ctx, keyDelayedEventList, 0, -1).Result()
	}

	if err != nil {
		return nil, fmt.Errorf("获取事件ID列表失败: %w", err)
	}

	// 过滤和分页
	var filteredEvents []*DelayedEvent
	for _, eventID := range eventIDs {
		event, err := s.GetDelayedEvent(ctx, eventID)
		if err != nil {
			s.logger.WithField("event_id", eventID).
				WithError(err).Warn("获取延时事件失败")
			continue
		}

		// 应用过滤器
		if s.matchesFilter(event, filter) {
			filteredEvents = append(filteredEvents, event)
		}
	}

	// 分页
	total := len(filteredEvents)
	start := (filter.Page - 1) * filter.PageSize
	end := start + filter.PageSize

	if start >= total {
		return &response.PageResult[*DelayedEvent]{
			Items:    []*DelayedEvent{},
			Total:    int64(total),
			Page:     int64(filter.Page),
			PageSize: int64(filter.PageSize),
		}, nil
	}

	if end > total {
		end = total
	}

	pagedEvents := filteredEvents[start:end]

	return &response.PageResult[*DelayedEvent]{
		Items:    pagedEvents,
		Total:    int64(total),
		Page:     int64(filter.Page),
		PageSize: int64(filter.PageSize),
	}, nil
}

// StoreEventRecord 存储事件记录
func (s *RedisEventStore) StoreEventRecord(ctx context.Context, record *EventRecord) error {
	// 生成记录ID
	if record.ID == "" {
		seq, err := s.client.Incr(ctx, keyEventSequence).Result()
		if err != nil {
			return fmt.Errorf("生成记录ID失败: %w", err)
		}
		record.ID = fmt.Sprintf("record_%d", seq)
	}

	// 序列化记录
	data, err := json.Marshal(record)
	if err != nil {
		return fmt.Errorf("序列化事件记录失败: %w", err)
	}

	pipe := s.client.Pipeline()

	// 存储记录数据
	recordKey := fmt.Sprintf("%s:%s", keyEventRecords, record.ID)
	pipe.HSet(ctx, recordKey, "data", string(data))

	// 添加到记录列表
	pipe.LPush(ctx, keyEventRecordList, record.ID)

	// 限制列表长度 (保留最近10000条记录)
	pipe.LTrim(ctx, keyEventRecordList, 0, 9999)

	// 执行事务
	if _, err := pipe.Exec(ctx); err != nil {
		return fmt.Errorf("存储事件记录失败: %w", err)
	}

	return nil
}

// GetEventRecords 获取事件记录
func (s *RedisEventStore) GetEventRecords(ctx context.Context, filter *EventFilter) (*response.PageResult[*EventRecord], error) {
	if filter.Page <= 0 {
		filter.Page = 1
	}
	if filter.PageSize <= 0 {
		filter.PageSize = 20
	}

	// 获取记录ID列表
	start := int64((filter.Page - 1) * filter.PageSize)
	end := start + int64(filter.PageSize) - 1

	recordIDs, err := s.client.LRange(ctx, keyEventRecordList, start, end).Result()
	if err != nil {
		return nil, fmt.Errorf("获取记录ID列表失败: %w", err)
	}

	// 获取总数
	total, err := s.client.LLen(ctx, keyEventRecordList).Result()
	if err != nil {
		return nil, fmt.Errorf("获取总数失败: %w", err)
	}

	// 批量获取记录数据
	var records []*EventRecord
	for _, recordID := range recordIDs {
		data, err := s.client.HGet(ctx, fmt.Sprintf("%s:%s", keyEventRecords, recordID), "data").Result()
		if err != nil {
			if err == redis.Nil {
				continue // 记录可能已被删除
			}
			s.logger.WithField("record_id", recordID).
				WithError(err).Warn("获取事件记录失败")
			continue
		}

		var record EventRecord
		if err := json.Unmarshal([]byte(data), &record); err != nil {
			s.logger.WithField("record_id", recordID).
				WithError(err).Warn("反序列化事件记录失败")
			continue
		}

		// 应用过滤器
		if s.matchesRecordFilter(&record, filter) {
			records = append(records, &record)
		}
	}

	return &response.PageResult[*EventRecord]{
		Items:    records,
		Total:    total,
		Page:     int64(filter.Page),
		PageSize: int64(filter.PageSize),
	}, nil
}

// GetMetrics 获取事件指标
func (s *RedisEventStore) GetMetrics(ctx context.Context) (*EventMetrics, error) {
	// 获取延时事件总数
	totalDelayed, err := s.client.ZCard(ctx, keyDelayedEventList).Result()
	if err != nil {
		return nil, fmt.Errorf("获取延时事件总数失败: %w", err)
	}

	// 获取准备就绪的事件数
	now := time.Now().Unix()
	readyCount, err := s.client.ZCount(ctx, keyDelayedEventList, "-inf", strconv.FormatInt(now, 10)).Result()
	if err != nil {
		return nil, fmt.Errorf("获取准备就绪事件数失败: %w", err)
	}

	// 获取记录总数
	totalRecords, err := s.client.LLen(ctx, keyEventRecordList).Result()
	if err != nil {
		return nil, fmt.Errorf("获取记录总数失败: %w", err)
	}

	return &EventMetrics{
		TotalEvents:     totalRecords,
		PendingEvents:   totalDelayed - readyCount,
		DeliveredEvents: readyCount,
		// TODO: 可以从缓存的指标中获取更详细的统计
	}, nil
}

// Cleanup 清理过期事件
func (s *RedisEventStore) Cleanup(ctx context.Context, expireBefore time.Time) (int64, error) {
	// 获取过期的事件ID
	expiredIDs, err := s.client.ZRangeByScore(ctx, keyDelayedEventList, &redis.ZRangeBy{
		Min: "-inf",
		Max: strconv.FormatInt(expireBefore.Unix(), 10),
	}).Result()
	if err != nil {
		return 0, fmt.Errorf("获取过期事件ID失败: %w", err)
	}

	if len(expiredIDs) == 0 {
		return 0, nil
	}

	// 批量删除过期事件
	var deletedCount int64
	for _, eventID := range expiredIDs {
		// 检查事件是否真的过期了
		event, err := s.GetDelayedEvent(ctx, eventID)
		if err != nil {
			s.logger.Warn("Failed to get event for cleanup",
				zap.String("event_id", eventID),
				zap.Error(err))
			continue
		}

		// 只删除失败状态且超过过期时间的事件
		if event.Status == EventStatusFailed && event.UpdatedAt.Before(expireBefore) {
			if err := s.DeleteDelayedEvent(ctx, eventID); err != nil {
				s.logger.WithField("event_id", eventID).
					WithError(err).Error("删除过期事件失败")
			} else {
				deletedCount++
			}
		}
	}

	s.logger.WithField("count", deletedCount).Info("清理过期事件完成")
	return deletedCount, nil
}

// matchesFilter 检查事件是否匹配过滤器
func (s *RedisEventStore) matchesFilter(event *DelayedEvent, filter *EventFilter) bool {
	// 检查事件类型
	if len(filter.Types) > 0 {
		found := false
		for _, t := range filter.Types {
			if t == event.Type {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}

	// 检查事件源
	if len(filter.Sources) > 0 {
		found := false
		for _, s := range filter.Sources {
			if s == event.Source {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}

	// 检查状态
	if filter.Status != "" && filter.Status != event.Status {
		return false
	}

	// 检查优先级
	if filter.Priority > 0 && filter.Priority != event.Priority {
		return false
	}

	return true
}

// matchesRecordFilter 检查记录是否匹配过滤器
func (s *RedisEventStore) matchesRecordFilter(record *EventRecord, filter *EventFilter) bool {
	// 检查事件类型
	if len(filter.Types) > 0 {
		found := false
		for _, t := range filter.Types {
			if t == record.Type {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}

	// 检查事件源
	if len(filter.Sources) > 0 {
		found := false
		for _, s := range filter.Sources {
			if s == record.Source {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}

	// 检查状态
	if filter.Status != "" && filter.Status != record.Status {
		return false
	}

	// 检查优先级
	if filter.Priority > 0 && filter.Priority != record.Priority {
		return false
	}

	return true
}

// 辅助方法
func (s *RedisEventStore) getDelayedEventKey() string {
	return keyDelayedEvents
}

func (s *RedisEventStore) getDelayedEventTimeKey() string {
	return keyDelayedEventList
}

func (s *RedisEventStore) getEventRecordsKey() string {
	return keyEventRecords
}
