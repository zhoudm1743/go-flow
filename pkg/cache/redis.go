package cache

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/zhoudm1743/go-flow/pkg/config"
	"github.com/zhoudm1743/go-flow/pkg/log"
	"go.uber.org/fx"
)

// Cache 缓存接口
type Cache interface {
	// 默认方法（不带 Context，使用 unified.Background()）
	// 基础操作
	Get(key string) (string, error)
	Set(key string, value interface{}, expiration time.Duration) error
	Del(keys ...string) (int64, error)
	Exists(keys ...string) (int64, error)
	Expire(key string, expiration time.Duration) error
	TTL(key string) (time.Duration, error)

	// 字符串操作
	Incr(key string) (int64, error)
	Decr(key string) (int64, error)
	IncrBy(key string, value int64) (int64, error)

	// 哈希操作
	HGet(key, field string) (string, error)
	HSet(key string, values ...interface{}) (int64, error)
	HDel(key string, fields ...string) (int64, error)
	HGetAll(key string) (map[string]string, error)
	HExists(key, field string) (bool, error)
	HLen(key string) (int64, error)

	// 列表操作
	LPush(key string, values ...interface{}) (int64, error)
	RPush(key string, values ...interface{}) (int64, error)
	LPop(key string) (string, error)
	RPop(key string) (string, error)
	LLen(key string) (int64, error)
	LRange(key string, start, stop int64) ([]string, error)

	// 集合操作
	SAdd(key string, members ...interface{}) (int64, error)
	SRem(key string, members ...interface{}) (int64, error)
	SMembers(key string) ([]string, error)
	SIsMember(key string, member interface{}) (bool, error)
	SCard(key string) (int64, error)

	// 有序集合操作
	ZAdd(key string, members ...redis.Z) (int64, error)
	ZRem(key string, members ...interface{}) (int64, error)
	ZRange(key string, start, stop int64) ([]string, error)
	ZRangeWithScores(key string, start, stop int64) ([]redis.Z, error)
	ZCard(key string) (int64, error)
	ZScore(key, member string) (float64, error)

	// 其他操作
	Keys(pattern string) ([]string, error)
	Ping() error

	// 带 Context 的方法（精细控制）
	// 基础操作
	GetCtx(ctx context.Context, key string) (string, error)
	SetCtx(ctx context.Context, key string, value interface{}, expiration time.Duration) error
	DelCtx(ctx context.Context, keys ...string) (int64, error)
	ExistsCtx(ctx context.Context, keys ...string) (int64, error)
	ExpireCtx(ctx context.Context, key string, expiration time.Duration) error
	TTLCtx(ctx context.Context, key string) (time.Duration, error)

	// 字符串操作
	IncrCtx(ctx context.Context, key string) (int64, error)
	DecrCtx(ctx context.Context, key string) (int64, error)
	IncrByCtx(ctx context.Context, key string, value int64) (int64, error)

	// 哈希操作
	HGetCtx(ctx context.Context, key, field string) (string, error)
	HSetCtx(ctx context.Context, key string, values ...interface{}) (int64, error)
	HDelCtx(ctx context.Context, key string, fields ...string) (int64, error)
	HGetAllCtx(ctx context.Context, key string) (map[string]string, error)
	HExistsCtx(ctx context.Context, key, field string) (bool, error)
	HLenCtx(ctx context.Context, key string) (int64, error)

	// 列表操作
	LPushCtx(ctx context.Context, key string, values ...interface{}) (int64, error)
	RPushCtx(ctx context.Context, key string, values ...interface{}) (int64, error)
	LPopCtx(ctx context.Context, key string) (string, error)
	RPopCtx(ctx context.Context, key string) (string, error)
	LLenCtx(ctx context.Context, key string) (int64, error)
	LRangeCtx(ctx context.Context, key string, start, stop int64) ([]string, error)

	// 集合操作
	SAddCtx(ctx context.Context, key string, members ...interface{}) (int64, error)
	SRemCtx(ctx context.Context, key string, members ...interface{}) (int64, error)
	SMembersCtx(ctx context.Context, key string) ([]string, error)
	SIsMemberCtx(ctx context.Context, key string, member interface{}) (bool, error)
	SCardCtx(ctx context.Context, key string) (int64, error)

	// 有序集合操作
	ZAddCtx(ctx context.Context, key string, members ...redis.Z) (int64, error)
	ZRemCtx(ctx context.Context, key string, members ...interface{}) (int64, error)
	ZRangeCtx(ctx context.Context, key string, start, stop int64) ([]string, error)
	ZRangeWithScoresCtx(ctx context.Context, key string, start, stop int64) ([]redis.Z, error)
	ZCardCtx(ctx context.Context, key string) (int64, error)
	ZScoreCtx(ctx context.Context, key, member string) (float64, error)

	// 其他操作
	KeysCtx(ctx context.Context, pattern string) ([]string, error)
	PingCtx(ctx context.Context) error

	// 工具方法
	Close() error
	GetClient() *redis.Client
}

// RedisCache Redis 缓存实现
type RedisCache struct {
	client *redis.Client
	logger log.Logger
	prefix string
}

// RedisModule Redis缓存模块
var RedisModule = fx.Options(
	fx.Provide(NewRedisCache),
)

// NewRedisCache 创建新的 Redis 缓存实例
func NewRedisCache(cfg *config.Config, log log.Logger) (Cache, error) {
	// 创建 Redis 客户端
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.Redis.Host, cfg.Redis.Port),
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
	})

	// 测试连接
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := rdb.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("Redis 连接失败: %w", err)
	}

	log.WithFields(map[string]interface{}{
		"host":   cfg.Redis.Host,
		"port":   cfg.Redis.Port,
		"db":     cfg.Redis.DB,
		"prefix": cfg.Redis.Prefix,
	}).Info("Redis 连接成功")

	return &RedisCache{
		client: rdb,
		logger: log,
		prefix: cfg.Redis.Prefix,
	}, nil
}

// buildKey 构建带前缀的键
func (r *RedisCache) buildKey(key string) string {
	if r.prefix == "" {
		return key
	}
	return r.prefix + key
}

// GetClient 获取原始 Redis 客户端
func (r *RedisCache) GetClient() *redis.Client {
	return r.client
}

// Close 关闭连接
func (r *RedisCache) Close() error {
	return r.client.Close()
}

// ================== 默认方法（不带 Context） ==================

// 基础操作
func (r *RedisCache) Get(key string) (string, error) {
	return r.client.Get(context.Background(), r.buildKey(key)).Result()
}

func (r *RedisCache) Set(key string, value interface{}, expiration time.Duration) error {
	return r.client.Set(context.Background(), r.buildKey(key), value, expiration).Err()
}

func (r *RedisCache) Del(keys ...string) (int64, error) {
	// 转换所有键为带前缀的键
	prefixedKeys := make([]string, len(keys))
	for i, key := range keys {
		prefixedKeys[i] = r.buildKey(key)
	}
	return r.client.Del(context.Background(), prefixedKeys...).Result()
}

func (r *RedisCache) Exists(keys ...string) (int64, error) {
	// 转换所有键为带前缀的键
	prefixedKeys := make([]string, len(keys))
	for i, key := range keys {
		prefixedKeys[i] = r.buildKey(key)
	}
	return r.client.Exists(context.Background(), prefixedKeys...).Result()
}

func (r *RedisCache) Expire(key string, expiration time.Duration) error {
	return r.client.Expire(context.Background(), r.buildKey(key), expiration).Err()
}

func (r *RedisCache) TTL(key string) (time.Duration, error) {
	return r.client.TTL(context.Background(), r.buildKey(key)).Result()
}

// 字符串操作
func (r *RedisCache) Incr(key string) (int64, error) {
	return r.client.Incr(context.Background(), r.buildKey(key)).Result()
}

func (r *RedisCache) Decr(key string) (int64, error) {
	return r.client.Decr(context.Background(), r.buildKey(key)).Result()
}

func (r *RedisCache) IncrBy(key string, value int64) (int64, error) {
	return r.client.IncrBy(context.Background(), r.buildKey(key), value).Result()
}

// 哈希操作
func (r *RedisCache) HGet(key, field string) (string, error) {
	return r.client.HGet(context.Background(), r.buildKey(key), field).Result()
}

func (r *RedisCache) HSet(key string, values ...interface{}) (int64, error) {
	return r.client.HSet(context.Background(), r.buildKey(key), values...).Result()
}

func (r *RedisCache) HDel(key string, fields ...string) (int64, error) {
	return r.client.HDel(context.Background(), r.buildKey(key), fields...).Result()
}

func (r *RedisCache) HGetAll(key string) (map[string]string, error) {
	return r.client.HGetAll(context.Background(), r.buildKey(key)).Result()
}

func (r *RedisCache) HExists(key, field string) (bool, error) {
	return r.client.HExists(context.Background(), r.buildKey(key), field).Result()
}

func (r *RedisCache) HLen(key string) (int64, error) {
	return r.client.HLen(context.Background(), r.buildKey(key)).Result()
}

// 列表操作
func (r *RedisCache) LPush(key string, values ...interface{}) (int64, error) {
	return r.client.LPush(context.Background(), r.buildKey(key), values...).Result()
}

func (r *RedisCache) RPush(key string, values ...interface{}) (int64, error) {
	return r.client.RPush(context.Background(), r.buildKey(key), values...).Result()
}

func (r *RedisCache) LPop(key string) (string, error) {
	return r.client.LPop(context.Background(), r.buildKey(key)).Result()
}

func (r *RedisCache) RPop(key string) (string, error) {
	return r.client.RPop(context.Background(), r.buildKey(key)).Result()
}

func (r *RedisCache) LLen(key string) (int64, error) {
	return r.client.LLen(context.Background(), r.buildKey(key)).Result()
}

func (r *RedisCache) LRange(key string, start, stop int64) ([]string, error) {
	return r.client.LRange(context.Background(), r.buildKey(key), start, stop).Result()
}

// 集合操作
func (r *RedisCache) SAdd(key string, members ...interface{}) (int64, error) {
	return r.client.SAdd(context.Background(), r.buildKey(key), members...).Result()
}

func (r *RedisCache) SRem(key string, members ...interface{}) (int64, error) {
	return r.client.SRem(context.Background(), r.buildKey(key), members...).Result()
}

func (r *RedisCache) SMembers(key string) ([]string, error) {
	return r.client.SMembers(context.Background(), r.buildKey(key)).Result()
}

func (r *RedisCache) SIsMember(key string, member interface{}) (bool, error) {
	return r.client.SIsMember(context.Background(), r.buildKey(key), member).Result()
}

func (r *RedisCache) SCard(key string) (int64, error) {
	return r.client.SCard(context.Background(), r.buildKey(key)).Result()
}

// 有序集合操作
func (r *RedisCache) ZAdd(key string, members ...redis.Z) (int64, error) {
	return r.client.ZAdd(context.Background(), r.buildKey(key), members...).Result()
}

func (r *RedisCache) ZRem(key string, members ...interface{}) (int64, error) {
	return r.client.ZRem(context.Background(), r.buildKey(key), members...).Result()
}

func (r *RedisCache) ZRange(key string, start, stop int64) ([]string, error) {
	return r.client.ZRange(context.Background(), r.buildKey(key), start, stop).Result()
}

func (r *RedisCache) ZRangeWithScores(key string, start, stop int64) ([]redis.Z, error) {
	return r.client.ZRangeWithScores(context.Background(), r.buildKey(key), start, stop).Result()
}

func (r *RedisCache) ZCard(key string) (int64, error) {
	return r.client.ZCard(context.Background(), r.buildKey(key)).Result()
}

func (r *RedisCache) ZScore(key, member string) (float64, error) {
	return r.client.ZScore(context.Background(), r.buildKey(key), member).Result()
}

// 其他操作
func (r *RedisCache) Keys(pattern string) ([]string, error) {
	// 对于 Keys 操作，我们需要添加前缀到模式中
	prefixedPattern := r.buildKey(pattern)
	keys, err := r.client.Keys(context.Background(), prefixedPattern).Result()
	if err != nil {
		return nil, err
	}

	// 移除前缀
	if r.prefix != "" {
		for i, key := range keys {
			keys[i] = key[len(r.prefix):]
		}
	}

	return keys, nil
}

func (r *RedisCache) Ping() error {
	return r.client.Ping(context.Background()).Err()
}

// ================== 带 Context 的方法 ==================

// 基础操作
func (r *RedisCache) GetCtx(ctx context.Context, key string) (string, error) {
	return r.client.Get(ctx, r.buildKey(key)).Result()
}

func (r *RedisCache) SetCtx(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	return r.client.Set(ctx, r.buildKey(key), value, expiration).Err()
}

func (r *RedisCache) DelCtx(ctx context.Context, keys ...string) (int64, error) {
	// 转换所有键为带前缀的键
	prefixedKeys := make([]string, len(keys))
	for i, key := range keys {
		prefixedKeys[i] = r.buildKey(key)
	}
	return r.client.Del(ctx, prefixedKeys...).Result()
}

func (r *RedisCache) ExistsCtx(ctx context.Context, keys ...string) (int64, error) {
	// 转换所有键为带前缀的键
	prefixedKeys := make([]string, len(keys))
	for i, key := range keys {
		prefixedKeys[i] = r.buildKey(key)
	}
	return r.client.Exists(ctx, prefixedKeys...).Result()
}

func (r *RedisCache) ExpireCtx(ctx context.Context, key string, expiration time.Duration) error {
	return r.client.Expire(ctx, r.buildKey(key), expiration).Err()
}

func (r *RedisCache) TTLCtx(ctx context.Context, key string) (time.Duration, error) {
	return r.client.TTL(ctx, r.buildKey(key)).Result()
}

// 字符串操作
func (r *RedisCache) IncrCtx(ctx context.Context, key string) (int64, error) {
	return r.client.Incr(ctx, key).Result()
}

func (r *RedisCache) DecrCtx(ctx context.Context, key string) (int64, error) {
	return r.client.Decr(ctx, key).Result()
}

func (r *RedisCache) IncrByCtx(ctx context.Context, key string, value int64) (int64, error) {
	return r.client.IncrBy(ctx, key, value).Result()
}

// 哈希操作
func (r *RedisCache) HGetCtx(ctx context.Context, key, field string) (string, error) {
	return r.client.HGet(ctx, key, field).Result()
}

func (r *RedisCache) HSetCtx(ctx context.Context, key string, values ...interface{}) (int64, error) {
	return r.client.HSet(ctx, key, values...).Result()
}

func (r *RedisCache) HDelCtx(ctx context.Context, key string, fields ...string) (int64, error) {
	return r.client.HDel(ctx, key, fields...).Result()
}

func (r *RedisCache) HGetAllCtx(ctx context.Context, key string) (map[string]string, error) {
	return r.client.HGetAll(ctx, key).Result()
}

func (r *RedisCache) HExistsCtx(ctx context.Context, key, field string) (bool, error) {
	return r.client.HExists(ctx, key, field).Result()
}

func (r *RedisCache) HLenCtx(ctx context.Context, key string) (int64, error) {
	return r.client.HLen(ctx, key).Result()
}

// 列表操作
func (r *RedisCache) LPushCtx(ctx context.Context, key string, values ...interface{}) (int64, error) {
	return r.client.LPush(ctx, key, values...).Result()
}

func (r *RedisCache) RPushCtx(ctx context.Context, key string, values ...interface{}) (int64, error) {
	return r.client.RPush(ctx, key, values...).Result()
}

func (r *RedisCache) LPopCtx(ctx context.Context, key string) (string, error) {
	return r.client.LPop(ctx, key).Result()
}

func (r *RedisCache) RPopCtx(ctx context.Context, key string) (string, error) {
	return r.client.RPop(ctx, key).Result()
}

func (r *RedisCache) LLenCtx(ctx context.Context, key string) (int64, error) {
	return r.client.LLen(ctx, key).Result()
}

func (r *RedisCache) LRangeCtx(ctx context.Context, key string, start, stop int64) ([]string, error) {
	return r.client.LRange(ctx, key, start, stop).Result()
}

// 集合操作
func (r *RedisCache) SAddCtx(ctx context.Context, key string, members ...interface{}) (int64, error) {
	return r.client.SAdd(ctx, key, members...).Result()
}

func (r *RedisCache) SRemCtx(ctx context.Context, key string, members ...interface{}) (int64, error) {
	return r.client.SRem(ctx, key, members...).Result()
}

func (r *RedisCache) SMembersCtx(ctx context.Context, key string) ([]string, error) {
	return r.client.SMembers(ctx, key).Result()
}

func (r *RedisCache) SIsMemberCtx(ctx context.Context, key string, member interface{}) (bool, error) {
	return r.client.SIsMember(ctx, key, member).Result()
}

func (r *RedisCache) SCardCtx(ctx context.Context, key string) (int64, error) {
	return r.client.SCard(ctx, key).Result()
}

// 有序集合操作
func (r *RedisCache) ZAddCtx(ctx context.Context, key string, members ...redis.Z) (int64, error) {
	return r.client.ZAdd(ctx, key, members...).Result()
}

func (r *RedisCache) ZRemCtx(ctx context.Context, key string, members ...interface{}) (int64, error) {
	return r.client.ZRem(ctx, key, members...).Result()
}

func (r *RedisCache) ZRangeCtx(ctx context.Context, key string, start, stop int64) ([]string, error) {
	return r.client.ZRange(ctx, key, start, stop).Result()
}

func (r *RedisCache) ZRangeWithScoresCtx(ctx context.Context, key string, start, stop int64) ([]redis.Z, error) {
	return r.client.ZRangeWithScores(ctx, key, start, stop).Result()
}

func (r *RedisCache) ZCardCtx(ctx context.Context, key string) (int64, error) {
	return r.client.ZCard(ctx, key).Result()
}

func (r *RedisCache) ZScoreCtx(ctx context.Context, key, member string) (float64, error) {
	return r.client.ZScore(ctx, key, member).Result()
}

// 其他操作
func (r *RedisCache) KeysCtx(ctx context.Context, pattern string) ([]string, error) {
	return r.client.Keys(ctx, pattern).Result()
}

func (r *RedisCache) PingCtx(ctx context.Context) error {
	return r.client.Ping(ctx).Err()
}
