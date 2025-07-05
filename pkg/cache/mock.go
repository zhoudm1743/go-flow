package cache

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/zhoudm1743/go-flow/pkg/config"
	"github.com/zhoudm1743/go-flow/pkg/log"
	"go.uber.org/fx"
)

// 错误定义
var (
	ErrKeyNotFound = errors.New("键不存在")
)

// MockCache 内存缓存实现，用于测试
type MockCache struct {
	data   map[string]interface{}
	expiry map[string]time.Time
	prefix string
	mu     sync.RWMutex
	logger log.Logger
}

// NewMockCache 创建内存缓存
func NewMockCache(cfg *config.Config, log log.Logger) (Cache, error) {
	log.Info("使用内存缓存（模拟Redis）")

	return &MockCache{
		data:   make(map[string]interface{}),
		expiry: make(map[string]time.Time),
		prefix: cfg.Redis.Prefix,
		logger: log,
	}, nil
}

// buildKey 构建带前缀的键
func (m *MockCache) buildKey(key string) string {
	if m.prefix == "" {
		return key
	}
	return m.prefix + key
}

// isExpired 检查键是否过期
func (m *MockCache) isExpired(key string) bool {
	if exp, ok := m.expiry[key]; ok {
		return exp.Before(time.Now())
	}
	return false
}

// cleanExpired 清理过期的键
func (m *MockCache) cleanExpired(key string) {
	if m.isExpired(key) {
		delete(m.data, key)
		delete(m.expiry, key)
	}
}

// GetClient 获取Redis客户端（模拟版本返回nil）
func (m *MockCache) GetClient() *redis.Client {
	return nil
}

// Close 关闭连接
func (m *MockCache) Close() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.data = make(map[string]interface{})
	m.expiry = make(map[string]time.Time)
	return nil
}

// Get 获取缓存
func (m *MockCache) Get(key string) (string, error) {
	return m.GetCtx(context.Background(), key)
}

// Set 设置缓存
func (m *MockCache) Set(key string, value interface{}, expiration time.Duration) error {
	return m.SetCtx(context.Background(), key, value, expiration)
}

// Del 删除缓存
func (m *MockCache) Del(keys ...string) (int64, error) {
	return m.DelCtx(context.Background(), keys...)
}

// Exists 检查键是否存在
func (m *MockCache) Exists(keys ...string) (int64, error) {
	return m.ExistsCtx(context.Background(), keys...)
}

// Expire 设置过期时间
func (m *MockCache) Expire(key string, expiration time.Duration) error {
	return m.ExpireCtx(context.Background(), key, expiration)
}

// TTL 获取过期时间
func (m *MockCache) TTL(key string) (time.Duration, error) {
	return m.TTLCtx(context.Background(), key)
}

// GetCtx 获取缓存
func (m *MockCache) GetCtx(ctx context.Context, key string) (string, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	fullKey := m.buildKey(key)
	m.cleanExpired(fullKey)

	if val, ok := m.data[fullKey]; ok {
		if str, ok := val.(string); ok {
			return str, nil
		}
		return "", errors.New("值类型不是字符串")
	}

	return "", ErrKeyNotFound
}

// SetCtx 设置缓存
func (m *MockCache) SetCtx(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	fullKey := m.buildKey(key)
	m.data[fullKey] = value

	if expiration > 0 {
		m.expiry[fullKey] = time.Now().Add(expiration)
	} else {
		delete(m.expiry, fullKey)
	}

	return nil
}

// DelCtx 删除缓存
func (m *MockCache) DelCtx(ctx context.Context, keys ...string) (int64, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	var count int64
	for _, key := range keys {
		fullKey := m.buildKey(key)
		if _, ok := m.data[fullKey]; ok {
			delete(m.data, fullKey)
			delete(m.expiry, fullKey)
			count++
		}
	}

	return count, nil
}

// ExistsCtx 检查键是否存在
func (m *MockCache) ExistsCtx(ctx context.Context, keys ...string) (int64, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	var count int64
	for _, key := range keys {
		fullKey := m.buildKey(key)
		m.cleanExpired(fullKey)
		if _, ok := m.data[fullKey]; ok {
			count++
		}
	}

	return count, nil
}

// ExpireCtx 设置过期时间
func (m *MockCache) ExpireCtx(ctx context.Context, key string, expiration time.Duration) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	fullKey := m.buildKey(key)
	if _, ok := m.data[fullKey]; !ok {
		return ErrKeyNotFound
	}

	if expiration > 0 {
		m.expiry[fullKey] = time.Now().Add(expiration)
	} else {
		delete(m.expiry, fullKey)
	}

	return nil
}

// TTLCtx 获取过期时间
func (m *MockCache) TTLCtx(ctx context.Context, key string) (time.Duration, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	fullKey := m.buildKey(key)
	m.cleanExpired(fullKey)

	if _, ok := m.data[fullKey]; !ok {
		return -2 * time.Second, nil // Redis返回-2表示键不存在
	}

	if exp, ok := m.expiry[fullKey]; ok {
		ttl := exp.Sub(time.Now())
		if ttl < 0 {
			return -1 * time.Second, nil // Redis返回-1表示键已过期
		}
		return ttl, nil
	}

	return -1 * time.Second, nil // 没有设置过期时间
}

// Incr 自增
func (m *MockCache) Incr(key string) (int64, error) {
	return m.IncrCtx(context.Background(), key)
}

// Decr 自减
func (m *MockCache) Decr(key string) (int64, error) {
	return m.DecrCtx(context.Background(), key)
}

// IncrBy 增加指定值
func (m *MockCache) IncrBy(key string, value int64) (int64, error) {
	return m.IncrByCtx(context.Background(), key, value)
}

// IncrCtx 自增
func (m *MockCache) IncrCtx(ctx context.Context, key string) (int64, error) {
	return m.IncrByCtx(ctx, key, 1)
}

// DecrCtx 自减
func (m *MockCache) DecrCtx(ctx context.Context, key string) (int64, error) {
	return m.IncrByCtx(ctx, key, -1)
}

// IncrByCtx 增加指定值
func (m *MockCache) IncrByCtx(ctx context.Context, key string, value int64) (int64, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	fullKey := m.buildKey(key)
	m.cleanExpired(fullKey)

	var current int64
	if val, ok := m.data[fullKey]; ok {
		switch v := val.(type) {
		case int64:
			current = v
		case string:
			// 尝试转换字符串为整数
			var err error
			current, err = parseInt64(v)
			if err != nil {
				return 0, err
			}
		default:
			return 0, errors.New("值类型不是整数")
		}
	}

	current += value
	m.data[fullKey] = current

	return current, nil
}

// parseInt64 尝试将字符串转换为int64
func parseInt64(s string) (int64, error) {
	var val int64
	_, err := fmt.Sscanf(s, "%d", &val)
	return val, err
}

// HGet 获取哈希字段
func (m *MockCache) HGet(key, field string) (string, error) {
	return m.HGetCtx(context.Background(), key, field)
}

// HSet 设置哈希字段
func (m *MockCache) HSet(key string, values ...interface{}) (int64, error) {
	return m.HSetCtx(context.Background(), key, values...)
}

// HDel 删除哈希字段
func (m *MockCache) HDel(key string, fields ...string) (int64, error) {
	return m.HDelCtx(context.Background(), key, fields...)
}

// HGetAll 获取所有哈希字段
func (m *MockCache) HGetAll(key string) (map[string]string, error) {
	return m.HGetAllCtx(context.Background(), key)
}

// HExists 检查哈希字段是否存在
func (m *MockCache) HExists(key, field string) (bool, error) {
	return m.HExistsCtx(context.Background(), key, field)
}

// HLen 获取哈希字段数量
func (m *MockCache) HLen(key string) (int64, error) {
	return m.HLenCtx(context.Background(), key)
}

// HGetCtx 获取哈希字段
func (m *MockCache) HGetCtx(ctx context.Context, key, field string) (string, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	fullKey := m.buildKey(key)
	m.cleanExpired(fullKey)

	if val, ok := m.data[fullKey]; ok {
		if hash, ok := val.(map[string]string); ok {
			if value, exists := hash[field]; exists {
				return value, nil
			}
			return "", ErrKeyNotFound
		}
		return "", errors.New("值类型不是哈希")
	}

	return "", ErrKeyNotFound
}

// HSetCtx 设置哈希字段
func (m *MockCache) HSetCtx(ctx context.Context, key string, values ...interface{}) (int64, error) {
	if len(values)%2 != 0 {
		return 0, errors.New("参数数量必须是偶数")
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	fullKey := m.buildKey(key)
	m.cleanExpired(fullKey)

	var hash map[string]string
	if val, ok := m.data[fullKey]; ok {
		if h, ok := val.(map[string]string); ok {
			hash = h
		} else {
			hash = make(map[string]string)
			m.data[fullKey] = hash
		}
	} else {
		hash = make(map[string]string)
		m.data[fullKey] = hash
	}

	var count int64
	for i := 0; i < len(values); i += 2 {
		field := fmt.Sprint(values[i])
		value := fmt.Sprint(values[i+1])

		if _, ok := hash[field]; !ok {
			count++
		}
		hash[field] = value
	}

	return count, nil
}

// HDelCtx 删除哈希字段
func (m *MockCache) HDelCtx(ctx context.Context, key string, fields ...string) (int64, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	fullKey := m.buildKey(key)
	m.cleanExpired(fullKey)

	if val, ok := m.data[fullKey]; ok {
		if hash, ok := val.(map[string]string); ok {
			var count int64
			for _, field := range fields {
				if _, exists := hash[field]; exists {
					delete(hash, field)
					count++
				}
			}
			return count, nil
		}
		return 0, errors.New("值类型不是哈希")
	}

	return 0, nil
}

// HGetAllCtx 获取所有哈希字段
func (m *MockCache) HGetAllCtx(ctx context.Context, key string) (map[string]string, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	fullKey := m.buildKey(key)
	m.cleanExpired(fullKey)

	if val, ok := m.data[fullKey]; ok {
		if hash, ok := val.(map[string]string); ok {
			result := make(map[string]string)
			for k, v := range hash {
				result[k] = v
			}
			return result, nil
		}
		return nil, errors.New("值类型不是哈希")
	}

	return make(map[string]string), nil
}

// HExistsCtx 检查哈希字段是否存在
func (m *MockCache) HExistsCtx(ctx context.Context, key, field string) (bool, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	fullKey := m.buildKey(key)
	m.cleanExpired(fullKey)

	if val, ok := m.data[fullKey]; ok {
		if hash, ok := val.(map[string]string); ok {
			_, exists := hash[field]
			return exists, nil
		}
		return false, errors.New("值类型不是哈希")
	}

	return false, nil
}

// HLenCtx 获取哈希字段数量
func (m *MockCache) HLenCtx(ctx context.Context, key string) (int64, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	fullKey := m.buildKey(key)
	m.cleanExpired(fullKey)

	if val, ok := m.data[fullKey]; ok {
		if hash, ok := val.(map[string]string); ok {
			return int64(len(hash)), nil
		}
		return 0, errors.New("值类型不是哈希")
	}

	return 0, nil
}

// LRange 获取列表范围
func (m *MockCache) LRange(key string, start, stop int64) ([]string, error) {
	return m.LRangeCtx(context.Background(), key, start, stop)
}

// LRangeCtx 获取列表范围
func (m *MockCache) LRangeCtx(ctx context.Context, key string, start, stop int64) ([]string, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	fullKey := m.buildKey(key)
	m.cleanExpired(fullKey)

	if val, ok := m.data[fullKey]; ok {
		if list, ok := val.([]string); ok {
			length := int64(len(list))

			// 处理负索引
			if start < 0 {
				start = length + start
				if start < 0 {
					start = 0
				}
			}
			if stop < 0 {
				stop = length + stop
				if stop < 0 {
					stop = -1
				}
			}

			// 处理超出范围的索引
			if start >= length || start > stop {
				return []string{}, nil
			}
			if stop >= length {
				stop = length - 1
			}

			result := make([]string, 0, stop-start+1)
			for i := start; i <= stop; i++ {
				result = append(result, list[i])
			}
			return result, nil
		}
		return nil, errors.New("值类型不是列表")
	}

	return []string{}, nil
}

// LLen 获取列表长度
func (m *MockCache) LLen(key string) (int64, error) {
	return m.LLenCtx(context.Background(), key)
}

// LLenCtx 获取列表长度
func (m *MockCache) LLenCtx(ctx context.Context, key string) (int64, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	fullKey := m.buildKey(key)
	m.cleanExpired(fullKey)

	if val, ok := m.data[fullKey]; ok {
		if list, ok := val.([]string); ok {
			return int64(len(list)), nil
		}
		return 0, errors.New("值类型不是列表")
	}

	return 0, nil
}

// LPush 向列表头部添加元素
func (m *MockCache) LPush(key string, values ...interface{}) (int64, error) {
	return m.LPushCtx(context.Background(), key, values...)
}

// LPushCtx 向列表头部添加元素
func (m *MockCache) LPushCtx(ctx context.Context, key string, values ...interface{}) (int64, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	fullKey := m.buildKey(key)
	m.cleanExpired(fullKey)

	var list []string
	if val, ok := m.data[fullKey]; ok {
		if l, ok := val.([]string); ok {
			list = l
		} else {
			list = make([]string, 0)
		}
	} else {
		list = make([]string, 0)
	}

	// 将值转换为字符串并添加到列表头部
	for i := len(values) - 1; i >= 0; i-- {
		list = append([]string{fmt.Sprint(values[i])}, list...)
	}

	m.data[fullKey] = list
	return int64(len(list)), nil
}

// RPush 向列表尾部添加元素
func (m *MockCache) RPush(key string, values ...interface{}) (int64, error) {
	return m.RPushCtx(context.Background(), key, values...)
}

// RPushCtx 向列表尾部添加元素
func (m *MockCache) RPushCtx(ctx context.Context, key string, values ...interface{}) (int64, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	fullKey := m.buildKey(key)
	m.cleanExpired(fullKey)

	var list []string
	if val, ok := m.data[fullKey]; ok {
		if l, ok := val.([]string); ok {
			list = l
		} else {
			list = make([]string, 0)
		}
	} else {
		list = make([]string, 0)
	}

	// 将值转换为字符串并添加到列表尾部
	for _, v := range values {
		list = append(list, fmt.Sprint(v))
	}

	m.data[fullKey] = list
	return int64(len(list)), nil
}

// LPop 从列表头部弹出元素
func (m *MockCache) LPop(key string) (string, error) {
	return m.LPopCtx(context.Background(), key)
}

// LPopCtx 从列表头部弹出元素
func (m *MockCache) LPopCtx(ctx context.Context, key string) (string, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	fullKey := m.buildKey(key)
	m.cleanExpired(fullKey)

	if val, ok := m.data[fullKey]; ok {
		if list, ok := val.([]string); ok {
			if len(list) == 0 {
				return "", ErrKeyNotFound
			}

			result := list[0]
			m.data[fullKey] = list[1:]
			return result, nil
		}
		return "", errors.New("值类型不是列表")
	}

	return "", ErrKeyNotFound
}

// RPop 从列表尾部弹出元素
func (m *MockCache) RPop(key string) (string, error) {
	return m.RPopCtx(context.Background(), key)
}

// RPopCtx 从列表尾部弹出元素
func (m *MockCache) RPopCtx(ctx context.Context, key string) (string, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	fullKey := m.buildKey(key)
	m.cleanExpired(fullKey)

	if val, ok := m.data[fullKey]; ok {
		if list, ok := val.([]string); ok {
			if len(list) == 0 {
				return "", ErrKeyNotFound
			}

			lastIndex := len(list) - 1
			result := list[lastIndex]
			m.data[fullKey] = list[:lastIndex]
			return result, nil
		}
		return "", errors.New("值类型不是列表")
	}

	return "", ErrKeyNotFound
}

// 简单实现其他必要方法，返回空值或错误
func (m *MockCache) SAdd(key string, members ...interface{}) (int64, error) {
	return 0, errors.New("未实现")
}

func (m *MockCache) SRem(key string, members ...interface{}) (int64, error) {
	return 0, errors.New("未实现")
}

func (m *MockCache) SMembers(key string) ([]string, error) {
	return nil, errors.New("未实现")
}

func (m *MockCache) SIsMember(key string, member interface{}) (bool, error) {
	return false, errors.New("未实现")
}

func (m *MockCache) SCard(key string) (int64, error) {
	return 0, errors.New("未实现")
}

func (m *MockCache) ZAdd(key string, members ...redis.Z) (int64, error) {
	return 0, errors.New("未实现")
}

func (m *MockCache) ZRem(key string, members ...interface{}) (int64, error) {
	return 0, errors.New("未实现")
}

func (m *MockCache) ZRange(key string, start, stop int64) ([]string, error) {
	return nil, errors.New("未实现")
}

func (m *MockCache) ZRangeWithScores(key string, start, stop int64) ([]redis.Z, error) {
	return nil, errors.New("未实现")
}

func (m *MockCache) ZCard(key string) (int64, error) {
	return 0, errors.New("未实现")
}

func (m *MockCache) ZScore(key, member string) (float64, error) {
	return 0, errors.New("未实现")
}

func (m *MockCache) Keys(pattern string) ([]string, error) {
	return m.KeysCtx(context.Background(), pattern)
}

func (m *MockCache) KeysCtx(ctx context.Context, pattern string) ([]string, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	// 简单实现，不支持通配符
	var keys []string
	prefix := m.buildKey("")
	for k := range m.data {
		if strings.HasPrefix(k, prefix) {
			keys = append(keys, k[len(prefix):])
		}
	}
	return keys, nil
}

func (m *MockCache) Ping() error {
	return m.PingCtx(context.Background())
}

func (m *MockCache) PingCtx(ctx context.Context) error {
	return nil
}

// 其他未实现的Ctx方法
func (m *MockCache) SAddCtx(ctx context.Context, key string, members ...interface{}) (int64, error) {
	return 0, errors.New("未实现")
}

func (m *MockCache) SRemCtx(ctx context.Context, key string, members ...interface{}) (int64, error) {
	return 0, errors.New("未实现")
}

func (m *MockCache) SMembersCtx(ctx context.Context, key string) ([]string, error) {
	return nil, errors.New("未实现")
}

func (m *MockCache) SIsMemberCtx(ctx context.Context, key string, member interface{}) (bool, error) {
	return false, errors.New("未实现")
}

func (m *MockCache) SCardCtx(ctx context.Context, key string) (int64, error) {
	return 0, errors.New("未实现")
}

func (m *MockCache) ZAddCtx(ctx context.Context, key string, members ...redis.Z) (int64, error) {
	return 0, errors.New("未实现")
}

func (m *MockCache) ZRemCtx(ctx context.Context, key string, members ...interface{}) (int64, error) {
	return 0, errors.New("未实现")
}

func (m *MockCache) ZRangeCtx(ctx context.Context, key string, start, stop int64) ([]string, error) {
	return nil, errors.New("未实现")
}

func (m *MockCache) ZRangeWithScoresCtx(ctx context.Context, key string, start, stop int64) ([]redis.Z, error) {
	return nil, errors.New("未实现")
}

func (m *MockCache) ZCardCtx(ctx context.Context, key string) (int64, error) {
	return 0, errors.New("未实现")
}

func (m *MockCache) ZScoreCtx(ctx context.Context, key, member string) (float64, error) {
	return 0, errors.New("未实现")
}

// MockCacheModule 提供内存缓存模块
var MockCacheModule = fx.Options(
	fx.Provide(
		fx.Annotated{
			Name:   "mock-cache",
			Target: NewMockCache,
		},
	),
)
