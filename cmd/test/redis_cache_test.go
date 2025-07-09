package cmd

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/zhoudm1743/go-frame/pkg/cache"
	"github.com/zhoudm1743/go-frame/pkg/config"
	"github.com/zhoudm1743/go-frame/pkg/log"
)

// 创建测试用的Redis缓存实例，使用不同的前缀以便测试之间互不干扰
func newTestRedisCacheInstance(t *testing.T, prefix string) cache.Cache {
	// 创建配置
	cfg := &config.Config{}
	cfg.Cache.Type = "redis"
	cfg.Cache.Prefix = prefix    // 使用唯一前缀
	cfg.Cache.Host = "localhost" // 请确保Redis在此地址运行
	cfg.Cache.Port = 6379        // 默认Redis端口
	cfg.Cache.Password = ""      // 根据需要设置密码
	cfg.Cache.DB = 0             // 使用默认数据库
	cfg.Log.Level = "info"
	cfg.Log.Format = "text"
	cfg.Log.OutputPath = "stdout"

	// 创建日志
	logger, _ := log.NewLogger(log.LoggerParams{
		Config: cfg,
	})

	// 创建缓存
	redisCache, err := cache.NewRedisCache(cfg, logger)
	if err != nil {
		t.Fatalf("无法创建Redis缓存: %v", err)
	}

	return redisCache
}

// 测试前清理Redis数据库
func cleanupRedisDB(t *testing.T, c cache.Cache, prefix string) {
	// 删除所有指定前缀的键
	keys, err := c.Keys(prefix + "*")
	if err != nil {
		t.Fatalf("无法获取测试键: %v", err)
	}

	if len(keys) > 0 {
		_, err = c.Del(keys...)
		if err != nil {
			t.Fatalf("无法清理测试数据: %v", err)
		}
	}
}

// 测试Redis缓存的所有功能
func TestRedisCacheAll(t *testing.T) {
	// 跳过如果环境变量设置为跳过Redis测试
	// 可以通过设置环境变量 SKIP_REDIS_TESTS=true 跳过此测试
	// if os.Getenv("SKIP_REDIS_TESTS") == "true" {
	//     t.Skip("跳过Redis测试")
	// }

	// 基础操作测试
	t.Run("基础操作", func(t *testing.T) {
		prefix := "test:base:"
		c := newTestRedisCacheInstance(t, prefix)
		defer c.Close()

		// 测试连接
		err := c.Ping()
		if err != nil {
			t.Skipf("Redis服务不可用，跳过测试: %v", err)
			return
		}

		cleanupRedisDB(t, c, prefix)

		// 测试 Set/Get
		err = c.Set("key1", "value1", time.Minute)
		assert.NoError(t, err)

		val, err := c.Get("key1")
		assert.NoError(t, err)
		assert.Equal(t, "value1", val)

		// 测试不存在的键
		_, err = c.Get("non_existent_key")
		assert.ErrorIs(t, err, redis.Nil) // Redis使用redis.Nil作为键不存在的错误

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
		assert.ErrorIs(t, err, redis.Nil)

		cleanupRedisDB(t, c, prefix)
	})

	// 过期时间测试
	t.Run("过期时间", func(t *testing.T) {
		prefix := "test:expire:"
		c := newTestRedisCacheInstance(t, prefix)
		defer c.Close()

		// 测试连接
		err := c.Ping()
		if err != nil {
			t.Skipf("Redis服务不可用，跳过测试: %v", err)
			return
		}

		cleanupRedisDB(t, c, prefix)

		// 设置过期时间为1秒
		err = c.Set("expire_key", "value", time.Second)
		assert.NoError(t, err)

		// 立即获取应该成功
		val, err := c.Get("expire_key")
		assert.NoError(t, err)
		assert.Equal(t, "value", val)

		// 等待过期
		time.Sleep(1100 * time.Millisecond)

		// 获取应该失败
		_, err = c.Get("expire_key")
		assert.ErrorIs(t, err, redis.Nil)

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
		assert.ErrorIs(t, err, redis.Nil)

		cleanupRedisDB(t, c, prefix)
	})

	// 计数器操作测试
	t.Run("计数器操作", func(t *testing.T) {
		prefix := "test:counter:"
		c := newTestRedisCacheInstance(t, prefix)
		defer c.Close()

		// 测试连接
		err := c.Ping()
		if err != nil {
			t.Skipf("Redis服务不可用，跳过测试: %v", err)
			return
		}

		cleanupRedisDB(t, c, prefix)

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

		cleanupRedisDB(t, c, prefix)
	})

	// 哈希表操作测试
	t.Run("哈希表操作", func(t *testing.T) {
		prefix := "test:hash:"
		c := newTestRedisCacheInstance(t, prefix)
		defer c.Close()

		// 测试连接
		err := c.Ping()
		if err != nil {
			t.Skipf("Redis服务不可用，跳过测试: %v", err)
			return
		}

		cleanupRedisDB(t, c, prefix)

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
		assert.ErrorIs(t, err, redis.Nil)

		cleanupRedisDB(t, c, prefix)
	})

	// 列表操作测试
	t.Run("列表操作", func(t *testing.T) {
		prefix := "test:list:"
		c := newTestRedisCacheInstance(t, prefix)
		defer c.Close()

		// 测试连接
		err := c.Ping()
		if err != nil {
			t.Skipf("Redis服务不可用，跳过测试: %v", err)
			return
		}

		cleanupRedisDB(t, c, prefix)

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

		// 测试 LRange
		values, err := c.LRange("list", 0, -1)
		assert.NoError(t, err)
		assert.Equal(t, []string{"value2", "value1", "value3"}, values) // 注意Redis的LPush是头部插入

		// 测试 LPop/RPop
		val, err := c.LPop("list")
		assert.NoError(t, err)
		assert.Equal(t, "value2", val)

		val, err = c.RPop("list")
		assert.NoError(t, err)
		assert.Equal(t, "value3", val)

		length, err = c.LLen("list")
		assert.NoError(t, err)
		assert.Equal(t, int64(1), length)

		cleanupRedisDB(t, c, prefix)
	})

	// 集合操作测试
	t.Run("集合操作", func(t *testing.T) {
		prefix := "test:set:"
		c := newTestRedisCacheInstance(t, prefix)
		defer c.Close()

		// 测试连接
		err := c.Ping()
		if err != nil {
			t.Skipf("Redis服务不可用，跳过测试: %v", err)
			return
		}

		cleanupRedisDB(t, c, prefix)

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

		cleanupRedisDB(t, c, prefix)
	})

	// 有序集合操作测试
	t.Run("有序集合操作", func(t *testing.T) {
		prefix := "test:zset:"
		c := newTestRedisCacheInstance(t, prefix)
		defer c.Close()

		// 测试连接
		err := c.Ping()
		if err != nil {
			t.Skipf("Redis服务不可用，跳过测试: %v", err)
			return
		}

		cleanupRedisDB(t, c, prefix)

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

		card, err = c.ZCard("zset")
		assert.NoError(t, err)
		assert.Equal(t, int64(2), card)

		cleanupRedisDB(t, c, prefix)
	})

	// Context 相关测试
	t.Run("Context相关", func(t *testing.T) {
		prefix := "test:ctx:"
		c := newTestRedisCacheInstance(t, prefix)
		defer c.Close()

		// 测试连接
		err := c.Ping()
		if err != nil {
			t.Skipf("Redis服务不可用，跳过测试: %v", err)
			return
		}

		cleanupRedisDB(t, c, prefix)

		ctx := context.Background()

		err = c.SetCtx(ctx, "ctx_key", "ctx_value", time.Minute)
		assert.NoError(t, err)

		val, err := c.GetCtx(ctx, "ctx_key")
		assert.NoError(t, err)
		assert.Equal(t, "ctx_value", val)

		cleanupRedisDB(t, c, prefix)
	})

	// Keys 模式匹配测试
	t.Run("Keys模式匹配", func(t *testing.T) {
		prefix := "test:keys:"
		c := newTestRedisCacheInstance(t, prefix)
		defer c.Close()

		// 测试连接
		err := c.Ping()
		if err != nil {
			t.Skipf("Redis服务不可用，跳过测试: %v", err)
			return
		}

		cleanupRedisDB(t, c, prefix)

		// 设置多个键
		err = c.Set("pattern:key1", "value1", time.Minute)
		assert.NoError(t, err)
		err = c.Set("pattern:key2", "value2", time.Minute)
		assert.NoError(t, err)
		err = c.Set("pattern:other", "value3", time.Minute)
		assert.NoError(t, err)

		// 测试模式匹配
		keys, err := c.Keys("pattern:key*")
		assert.NoError(t, err)
		assert.Len(t, keys, 2)
		assert.Contains(t, keys, "pattern:key1")
		assert.Contains(t, keys, "pattern:key2")

		// 测试完全匹配
		keys, err = c.Keys("pattern:key1")
		assert.NoError(t, err)
		assert.Len(t, keys, 1)
		assert.Equal(t, "pattern:key1", keys[0])

		// 测试通配符
		keys, err = c.Keys("pattern:*")
		assert.NoError(t, err)
		assert.Len(t, keys, 3)

		cleanupRedisDB(t, c, prefix)
	})

	// 并发测试
	t.Run("并发操作", func(t *testing.T) {
		prefix := "test:concurrent:"
		c := newTestRedisCacheInstance(t, prefix)
		defer c.Close()

		// 测试连接
		err := c.Ping()
		if err != nil {
			t.Skipf("Redis服务不可用，跳过测试: %v", err)
			return
		}

		cleanupRedisDB(t, c, prefix)

		const concurrency = 10
		const iterations = 100

		// 并发设置和获取
		done := make(chan bool)
		for i := 0; i < concurrency; i++ {
			go func(id int) {
				for j := 0; j < iterations; j++ {
					key := fmt.Sprintf("concurrent:%d:%d", id, j)
					value := fmt.Sprintf("value:%d:%d", id, j)

					err := c.Set(key, value, time.Minute)
					assert.NoError(t, err)

					val, err := c.Get(key)
					assert.NoError(t, err)
					assert.Equal(t, value, val)
				}
				done <- true
			}(i)
		}

		// 等待所有goroutine完成
		for i := 0; i < concurrency; i++ {
			<-done
		}

		cleanupRedisDB(t, c, prefix)
	})
}
