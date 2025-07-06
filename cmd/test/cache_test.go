package cmd

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/zhoudm1743/go-flow/pkg/cache"
	"github.com/zhoudm1743/go-flow/pkg/config"
	"github.com/zhoudm1743/go-flow/pkg/log"
)

// 创建测试用的内存缓存实例
func newTestMemoryCache(t *testing.T) cache.Cache {
	// 创建配置
	cfg := &config.Config{}
	cfg.Cache.Type = "memory"
	cfg.Cache.Prefix = "test:"
	cfg.Log.Level = "info"
	cfg.Log.Format = "text"
	cfg.Log.OutputPath = "stdout"

	// 创建日志
	logger, _ := log.NewLogger(log.LoggerParams{
		Config: cfg,
	})

	// 创建缓存
	memCache, err := cache.NewMemoryCache(cfg, logger)
	if err != nil {
		t.Fatalf("无法创建内存缓存: %v", err)
	}

	return memCache
}

// 测试基本的存取操作
func TestMemoryCacheBasicOps(t *testing.T) {
	c := newTestMemoryCache(t)
	defer c.Close()

	// 测试 Set/Get
	err := c.Set("key1", "value1", time.Minute)
	assert.NoError(t, err)

	val, err := c.Get("key1")
	assert.NoError(t, err)
	assert.Equal(t, "value1", val)

	// 测试不存在的键
	_, err = c.Get("non_existent_key")
	assert.Equal(t, cache.ErrKeyNotFound, err)

	// 测试 Exists
	count, err := c.Exists("key1")
	assert.NoError(t, err)
	assert.Equal(t, int64(1), count)

	count, err = c.Exists("key1", "non_existent_key")
	assert.NoError(t, err)
	assert.Equal(t, int64(1), count)

	// 测试 Del
	_, err = c.Del("key1")
	assert.NoError(t, err)

	_, err = c.Get("key1")
	assert.Equal(t, cache.ErrKeyNotFound, err)
}

// 测试过期时间
func TestMemoryCacheExpiration(t *testing.T) {
	c := newTestMemoryCache(t)
	defer c.Close()

	// 设置过期时间为1秒
	err := c.Set("expire_key", "value", time.Second)
	assert.NoError(t, err)

	// 立即获取应该成功
	val, err := c.Get("expire_key")
	assert.NoError(t, err)
	assert.Equal(t, "value", val)

	// 等待过期
	time.Sleep(1100 * time.Millisecond)

	// 获取应该失败
	_, err = c.Get("expire_key")
	assert.Equal(t, cache.ErrKeyNotFound, err)

	// 测试 TTL
	err = c.Set("ttl_key", "value", time.Minute)
	assert.NoError(t, err)

	ttl, err := c.TTL("ttl_key")
	assert.NoError(t, err)
	assert.True(t, ttl > 0)

	// 测试更新过期时间
	err = c.Expire("ttl_key", 2*time.Second)
	assert.NoError(t, err)

	ttl, err = c.TTL("ttl_key")
	assert.NoError(t, err)
	assert.True(t, ttl <= 2*time.Second)

	// 等待过期
	time.Sleep(2100 * time.Millisecond)

	_, err = c.Get("ttl_key")
	assert.Equal(t, cache.ErrKeyNotFound, err)
}

// 测试计数器操作
func TestMemoryCacheCounter(t *testing.T) {
	c := newTestMemoryCache(t)
	defer c.Close()

	// 测试 Incr
	val, err := c.Incr("counter")
	assert.NoError(t, err)
	assert.Equal(t, int64(1), val)

	val, err = c.Incr("counter")
	assert.NoError(t, err)
	assert.Equal(t, int64(2), val)

	// 测试 IncrBy
	val, err = c.IncrBy("counter", 3)
	assert.NoError(t, err)
	assert.Equal(t, int64(5), val)

	// 测试 Decr
	val, err = c.Decr("counter")
	assert.NoError(t, err)
	assert.Equal(t, int64(4), val)
}

// 测试哈希表操作
func TestMemoryCacheHash(t *testing.T) {
	c := newTestMemoryCache(t)
	defer c.Close()

	// 测试 HSet/HGet
	count, err := c.HSet("hash", "field1", "value1", "field2", "value2")
	assert.NoError(t, err)
	assert.Equal(t, int64(2), count)

	val, err := c.HGet("hash", "field1")
	assert.NoError(t, err)
	assert.Equal(t, "value1", val)

	// 测试 HGetAll
	all, err := c.HGetAll("hash")
	assert.NoError(t, err)
	assert.Equal(t, map[string]string{"field1": "value1", "field2": "value2"}, all)

	// 测试 HExists
	exists, err := c.HExists("hash", "field1")
	assert.NoError(t, err)
	assert.True(t, exists)

	exists, err = c.HExists("hash", "field3")
	assert.NoError(t, err)
	assert.False(t, exists)

	// 测试 HLen
	length, err := c.HLen("hash")
	assert.NoError(t, err)
	assert.Equal(t, int64(2), length)

	// 测试 HDel
	count, err = c.HDel("hash", "field1")
	assert.NoError(t, err)
	assert.Equal(t, int64(1), count)

	_, err = c.HGet("hash", "field1")
	assert.Equal(t, cache.ErrKeyNotFound, err)
}

// 测试列表操作
func TestMemoryCacheList(t *testing.T) {
	c := newTestMemoryCache(t)
	defer c.Close()

	// 测试 LPush/RPush
	count, err := c.LPush("list", "value1", "value2")
	assert.NoError(t, err)
	assert.Equal(t, int64(2), count)

	count, err = c.RPush("list", "value3")
	assert.NoError(t, err)
	assert.Equal(t, int64(3), count)

	// 测试 LLen
	length, err := c.LLen("list")
	assert.NoError(t, err)
	assert.Equal(t, int64(3), length)

	// 测试 LRange - Redis的LPush顺序是后加入的在前面，前加入的在后面
	values, err := c.LRange("list", 0, -1)
	assert.NoError(t, err)
	assert.Equal(t, []string{"value1", "value2", "value3"}, values)

	// 测试 LPop/RPop
	val, err := c.LPop("list")
	assert.NoError(t, err)
	assert.Equal(t, "value1", val)

	val, err = c.RPop("list")
	assert.NoError(t, err)
	assert.Equal(t, "value3", val)

	length, err = c.LLen("list")
	assert.NoError(t, err)
	assert.Equal(t, int64(1), length)
}

// 测试集合操作
func TestMemoryCacheSet(t *testing.T) {
	c := newTestMemoryCache(t)
	defer c.Close()

	// 测试 SAdd
	count, err := c.SAdd("set", "member1", "member2", "member3")
	assert.NoError(t, err)
	assert.Equal(t, int64(3), count)

	// 测试 SCard
	card, err := c.SCard("set")
	assert.NoError(t, err)
	assert.Equal(t, int64(3), card)

	// 测试 SIsMember
	isMember, err := c.SIsMember("set", "member1")
	assert.NoError(t, err)
	assert.True(t, isMember)

	isMember, err = c.SIsMember("set", "member4")
	assert.NoError(t, err)
	assert.False(t, isMember)

	// 测试 SMembers
	members, err := c.SMembers("set")
	assert.NoError(t, err)
	assert.Len(t, members, 3)
	assert.Contains(t, members, "member1")
	assert.Contains(t, members, "member2")
	assert.Contains(t, members, "member3")

	// 测试 SRem
	count, err = c.SRem("set", "member1")
	assert.NoError(t, err)
	assert.Equal(t, int64(1), count)

	isMember, err = c.SIsMember("set", "member1")
	assert.NoError(t, err)
	assert.False(t, isMember)
}

// 测试有序集合操作
func TestMemoryCacheZSet(t *testing.T) {
	c := newTestMemoryCache(t)
	defer c.Close()

	// 测试 ZAdd
	count, err := c.ZAdd("zset",
		cache.Z{Score: 1, Member: "member1"},
		cache.Z{Score: 2, Member: "member2"},
		cache.Z{Score: 3, Member: "member3"},
	)
	assert.NoError(t, err)
	assert.Equal(t, int64(3), count)

	// 测试 ZCard
	card, err := c.ZCard("zset")
	assert.NoError(t, err)
	assert.Equal(t, int64(3), card)

	// 测试 ZScore
	score, err := c.ZScore("zset", "member2")
	assert.NoError(t, err)
	assert.Equal(t, float64(2), score)

	// 测试 ZRange
	members, err := c.ZRange("zset", 0, -1)
	assert.NoError(t, err)
	assert.Equal(t, []string{"member1", "member2", "member3"}, members)

	// 测试 ZRangeWithScores
	zs, err := c.ZRangeWithScores("zset", 0, -1)
	assert.NoError(t, err)
	assert.Equal(t, []cache.Z{
		{Score: 1, Member: "member1"},
		{Score: 2, Member: "member2"},
		{Score: 3, Member: "member3"},
	}, zs)

	// 测试 ZRem
	count, err = c.ZRem("zset", "member1")
	assert.NoError(t, err)
	assert.Equal(t, int64(1), count)

	members, err = c.ZRange("zset", 0, -1)
	assert.NoError(t, err)
	assert.Equal(t, []string{"member2", "member3"}, members)
}

// 测试上下文相关操作
func TestMemoryCacheContext(t *testing.T) {
	c := newTestMemoryCache(t)
	defer c.Close()

	// 创建上下文
	ctx := context.Background()

	// 测试带上下文的操作
	err := c.SetCtx(ctx, "ctx_key", "value", time.Minute)
	assert.NoError(t, err)

	val, err := c.GetCtx(ctx, "ctx_key")
	assert.NoError(t, err)
	assert.Equal(t, "value", val)

	// 测试带取消的上下文
	ctxWithCancel, cancel := context.WithCancel(ctx)
	cancel() // 立即取消

	// 由于内存缓存实现忽略了上下文取消，这里应该仍然成功
	// 这个测试主要是确保接口兼容性
	val, err = c.GetCtx(ctxWithCancel, "ctx_key")
	assert.NoError(t, err)
	assert.Equal(t, "value", val)
}

// 测试 Keys 操作
func TestMemoryCacheKeys(t *testing.T) {
	c := newTestMemoryCache(t)
	defer c.Close()

	// 添加一些键
	err := c.Set("key1", "value1", time.Minute)
	assert.NoError(t, err)

	err = c.Set("key2", "value2", time.Minute)
	assert.NoError(t, err)

	err = c.Set("other", "value3", time.Minute)
	assert.NoError(t, err)

	// 测试 Keys - 通配符前缀匹配
	keys, err := c.Keys("key*")
	assert.NoError(t, err)
	assert.Len(t, keys, 2)
	assert.Contains(t, keys, "key1")
	assert.Contains(t, keys, "key2")

	// 测试 Keys 通配符 - 全匹配
	keys, err = c.Keys("*")
	assert.NoError(t, err)
	assert.Len(t, keys, 3)
	assert.Contains(t, keys, "key1")
	assert.Contains(t, keys, "key2")
	assert.Contains(t, keys, "other")

	// 测试精确匹配
	keys, err = c.Keys("key1")
	assert.NoError(t, err)
	assert.Len(t, keys, 1)
	assert.Contains(t, keys, "key1")

	// 测试 Ping
	err = c.Ping()
	assert.NoError(t, err)
}
