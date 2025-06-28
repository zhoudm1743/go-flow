package captcha

import (
	"context"
	"time"

	"github.com/mojocn/base64Captcha"
	"github.com/zhoudm1743/go-flow/core/cache"
)

// 内存存储（修复参数错误）
var store = base64Captcha.NewMemoryStore(200, 10*time.Minute)

// RedisStore Redis验证码存储实现
type RedisStore struct {
	cache      cache.Cache
	keyPrefix  string
	expiration time.Duration
}

// NewRedisStore 创建Redis验证码存储
func NewRedisStore(cacheClient cache.Cache, expiration time.Duration, keyPrefix string) base64Captcha.Store {
	if keyPrefix == "" {
		keyPrefix = "captcha"
	}
	return &RedisStore{
		cache:      cacheClient,
		keyPrefix:  keyPrefix,
		expiration: expiration,
	}
}

// buildKey 构建Redis键
func (rs *RedisStore) buildKey(id string) string {
	return rs.keyPrefix + ":" + id
}

// Set 设置验证码
func (rs *RedisStore) Set(id string, value string) error {
	key := rs.buildKey(id)
	return rs.cache.SetCtx(context.Background(), key, value, rs.expiration)
}

// Get 获取验证码
func (rs *RedisStore) Get(id string, clear bool) string {
	key := rs.buildKey(id)

	value, err := rs.cache.GetCtx(context.Background(), key)
	if err != nil {
		return ""
	}

	// 如果需要清除，删除key
	if clear {
		rs.cache.DelCtx(context.Background(), key)
	}

	return value
}

// Verify 验证验证码
func (rs *RedisStore) Verify(id, answer string, clear bool) bool {
	stored := rs.Get(id, clear)
	return stored == answer && stored != ""
}

// GetStore 获取默认的内存存储
func GetStore() base64Captcha.Store {
	return store
}

// GetRedisStore 获取Redis存储实例
func GetRedisStore(cacheClient cache.Cache) base64Captcha.Store {
	return NewRedisStore(cacheClient, 10*time.Minute, "captcha")
}
