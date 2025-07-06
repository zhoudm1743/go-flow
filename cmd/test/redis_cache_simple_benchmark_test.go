package cmd

import (
	"fmt"
	"testing"

	"github.com/zhoudm1743/go-flow/pkg/cache"
	"github.com/zhoudm1743/go-flow/pkg/config"
	"github.com/zhoudm1743/go-flow/pkg/log"
)

// 创建基准测试用的Redis缓存实例
func newSimpleRedisCache(b *testing.B) cache.Cache {
	// 创建配置
	cfg := &config.Config{}
	cfg.Cache.Type = "redis"
	cfg.Cache.Prefix = "bench:"
	cfg.Cache.Host = "localhost"
	cfg.Cache.Port = 6379
	cfg.Cache.Password = ""
	cfg.Cache.DB = 0
	cfg.Log.Level = "error" // 降低日志级别，避免影响性能测试
	cfg.Log.Format = "text"
	cfg.Log.OutputPath = "none"

	// 创建日志
	logger, _ := log.NewLogger(log.LoggerParams{
		Config: cfg,
	})

	// 创建缓存
	redisCache, err := cache.NewRedisCache(cfg, logger)
	if err != nil {
		b.Skipf("无法创建Redis缓存: %v", err)
	}

	// 检查Redis连接
	err = redisCache.Ping()
	if err != nil {
		b.Skipf("Redis服务不可用，跳过基准测试: %v", err)
	}

	return redisCache
}

// 清理基准测试数据
func cleanupSimpleRedisDB(b *testing.B, c cache.Cache) {
	// 删除所有bench:前缀的键
	keys, err := c.Keys("bench:*")
	if err != nil {
		b.Fatalf("无法获取基准测试键: %v", err)
	}

	if len(keys) > 0 {
		_, err = c.Del(keys...)
		if err != nil {
			b.Fatalf("无法清理基准测试数据: %v", err)
		}
	}
}

// BenchmarkRedisSimpleGet 基准测试Get操作
func BenchmarkRedisSimpleGet(b *testing.B) {
	c := newSimpleRedisCache(b)
	defer c.Close()
	defer cleanupSimpleRedisDB(b, c)

	// 预先设置一个键
	err := c.Set("test-key", "test-value", 0)
	if err != nil {
		b.Fatalf("设置缓存失败: %v", err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := c.Get("test-key")
		if err != nil {
			b.Fatalf("获取缓存失败: %v", err)
		}
	}
}

// BenchmarkRedisSimpleGetParallel 并发基准测试Get操作
func BenchmarkRedisSimpleGetParallel(b *testing.B) {
	c := newSimpleRedisCache(b)
	defer c.Close()
	defer cleanupSimpleRedisDB(b, c)

	// 预先设置一个键
	err := c.Set("test-key", "test-value", 0)
	if err != nil {
		b.Fatalf("设置缓存失败: %v", err)
	}

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_, err := c.Get("test-key")
			if err != nil {
				b.Fatalf("获取缓存失败: %v", err)
			}
		}
	})
}

// BenchmarkRedisSimpleSet 基准测试Set操作
func BenchmarkRedisSimpleSet(b *testing.B) {
	c := newSimpleRedisCache(b)
	defer c.Close()
	defer cleanupSimpleRedisDB(b, c)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		err := c.Set("test-key", "test-value", 0)
		if err != nil {
			b.Fatalf("设置缓存失败: %v", err)
		}
	}
}

// BenchmarkRedisSimpleSetParallel 并发基准测试Set操作
func BenchmarkRedisSimpleSetParallel(b *testing.B) {
	c := newSimpleRedisCache(b)
	defer c.Close()
	defer cleanupSimpleRedisDB(b, c)

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		counter := 0
		for pb.Next() {
			key := fmt.Sprintf("test-key-%d", counter)
			err := c.Set(key, "test-value", 0)
			if err != nil {
				b.Fatalf("设置缓存失败: %v", err)
			}
			counter++
		}
	})
}
