package facades

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

// ================== CacheFacade 方法 ==================

// Get 获取缓存
func (f *CacheFacade) Get(key string) (string, error) {
	return GetCache().Get(key)
}

// Set 设置缓存
func (f *CacheFacade) Set(key string, value interface{}, expiration time.Duration) error {
	return GetCache().Set(key, value, expiration)
}

// Del 删除缓存
func (f *CacheFacade) Del(keys ...string) (int64, error) {
	return GetCache().Del(keys...)
}

// Exists 检查键是否存在
func (f *CacheFacade) Exists(keys ...string) (int64, error) {
	return GetCache().Exists(keys...)
}

// Expire 设置过期时间
func (f *CacheFacade) Expire(key string, expiration time.Duration) error {
	return GetCache().Expire(key, expiration)
}

// TTL 获取过期时间
func (f *CacheFacade) TTL(key string) (time.Duration, error) {
	return GetCache().TTL(key)
}

// Incr 自增
func (f *CacheFacade) Incr(key string) (int64, error) {
	return GetCache().Incr(key)
}

// Decr 自减
func (f *CacheFacade) Decr(key string) (int64, error) {
	return GetCache().Decr(key)
}

// HGet 获取哈希字段
func (f *CacheFacade) HGet(key, field string) (string, error) {
	return GetCache().HGet(key, field)
}

// HSet 设置哈希字段
func (f *CacheFacade) HSet(key string, values ...interface{}) (int64, error) {
	return GetCache().HSet(key, values...)
}

// HGetAll 获取所有哈希字段
func (f *CacheFacade) HGetAll(key string) (map[string]string, error) {
	return GetCache().HGetAll(key)
}

// Keys 获取匹配的键
func (f *CacheFacade) Keys(pattern string) ([]string, error) {
	return GetCache().Keys(pattern)
}

// Client 获取原始Redis客户端
func (f *CacheFacade) Client() *redis.Client {
	return GetCache().GetClient()
}

// WithContext 使用上下文执行操作
func (f *CacheFacade) WithContext(ctx context.Context) CacheContextFacade {
	return CacheContextFacade{ctx: ctx}
}

// CacheContextFacade 带上下文的缓存门面
type CacheContextFacade struct {
	ctx context.Context
}

// Get 获取缓存
func (f CacheContextFacade) Get(key string) (string, error) {
	return GetCache().GetCtx(f.ctx, key)
}

// Set 设置缓存
func (f CacheContextFacade) Set(key string, value interface{}, expiration time.Duration) error {
	return GetCache().SetCtx(f.ctx, key, value, expiration)
}

// Del 删除缓存
func (f CacheContextFacade) Del(keys ...string) (int64, error) {
	return GetCache().DelCtx(f.ctx, keys...)
}

// Exists 检查键是否存在
func (f CacheContextFacade) Exists(keys ...string) (int64, error) {
	return GetCache().ExistsCtx(f.ctx, keys...)
}

// Expire 设置过期时间
func (f CacheContextFacade) Expire(key string, expiration time.Duration) error {
	return GetCache().ExpireCtx(f.ctx, key, expiration)
}

// TTL 获取过期时间
func (f CacheContextFacade) TTL(key string) (time.Duration, error) {
	return GetCache().TTLCtx(f.ctx, key)
}

// LPush 向列表头部添加元素
func (f *CacheFacade) LPush(key string, values ...interface{}) (int64, error) {
	return GetCache().LPush(key, values...)
}

// RPush 向列表尾部添加元素
func (f *CacheFacade) RPush(key string, values ...interface{}) (int64, error) {
	return GetCache().RPush(key, values...)
}

// LPop 从列表头部弹出元素
func (f *CacheFacade) LPop(key string) (string, error) {
	return GetCache().LPop(key)
}

// RPop 从列表尾部弹出元素
func (f *CacheFacade) RPop(key string) (string, error) {
	return GetCache().RPop(key)
}

// LLen 获取列表长度
func (f *CacheFacade) LLen(key string) (int64, error) {
	return GetCache().LLen(key)
}

// LRange 获取列表范围
func (f *CacheFacade) LRange(key string, start, stop int64) ([]string, error) {
	return GetCache().LRange(key, start, stop)
}
