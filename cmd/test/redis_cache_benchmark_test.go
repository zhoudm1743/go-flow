package cmd

import (
	"fmt"
	"testing"
	"time"

	"github.com/zhoudm1743/go-frame/pkg/cache"
	"github.com/zhoudm1743/go-frame/pkg/config"
	"github.com/zhoudm1743/go-frame/pkg/log"
)

// 创建基准测试用的Redis缓存实例
func newBenchmarkRedisCacheInstance(b *testing.B) cache.Cache {
	// 创建配置
	cfg := &config.Config{}
	cfg.Cache.Type = "redis"
	cfg.Cache.Prefix = "benchmark:"
	cfg.Cache.Host = "localhost"
	cfg.Cache.Port = 6379
	cfg.Cache.Password = ""
	cfg.Cache.DB = 0
	cfg.Log.Level = "error" // 基准测试时禁止日志输出
	cfg.Log.Format = "text"
	cfg.Log.OutputPath = "none"

	// 创建日志
	logger, _ := log.NewLogger(log.LoggerParams{
		Config: cfg,
	})

	// 创建缓存
	redisCache, err := cache.NewRedisCache(cfg, logger)
	if err != nil {
		b.Fatalf("无法创建Redis缓存: %v", err)
	}

	return redisCache
}

// 清理基准测试数据
func cleanupBenchmarkRedisDB(b *testing.B, c cache.Cache) {
	// 删除所有benchmark:前缀的键
	keys, err := c.Keys("benchmark:*")
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

// 跳过基准测试如果Redis不可用
func skipIfRedisUnavailable(b *testing.B, c cache.Cache) {
	err := c.Ping()
	if err != nil {
		b.Skipf("Redis服务不可用，跳过基准测试: %v", err)
	}
}

// 基准测试: Get 方法
func BenchmarkRedisCache_Get(b *testing.B) {
	c := newBenchmarkRedisCacheInstance(b)
	defer c.Close()

	// 检查Redis连接并跳过测试如果不可用
	skipIfRedisUnavailable(b, c)

	// 清理旧数据并准备测试数据
	cleanupBenchmarkRedisDB(b, c)
	err := c.Set("bench_get", "value", 0)
	if err != nil {
		b.Fatalf("无法设置测试数据: %v", err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := c.Get("bench_get")
		if err != nil {
			b.Fatalf("Get失败: %v", err)
		}
	}
	b.StopTimer()

	cleanupBenchmarkRedisDB(b, c)
}

// 基准测试: Set 方法
func BenchmarkRedisCache_Set(b *testing.B) {
	c := newBenchmarkRedisCacheInstance(b)
	defer c.Close()

	// 检查Redis连接并跳过测试如果不可用
	skipIfRedisUnavailable(b, c)

	// 清理旧数据
	cleanupBenchmarkRedisDB(b, c)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		key := fmt.Sprintf("bench_set_%d", i)
		err := c.Set(key, "value", time.Minute)
		if err != nil {
			b.Fatalf("Set失败: %v", err)
		}
	}
	b.StopTimer()

	cleanupBenchmarkRedisDB(b, c)
}

// 基准测试: HGet 方法
func BenchmarkRedisCache_HGet(b *testing.B) {
	c := newBenchmarkRedisCacheInstance(b)
	defer c.Close()

	// 检查Redis连接并跳过测试如果不可用
	skipIfRedisUnavailable(b, c)

	// 清理旧数据并准备测试数据
	cleanupBenchmarkRedisDB(b, c)
	_, err := c.HSet("bench_hash", "field1", "value1", "field2", "value2")
	if err != nil {
		b.Fatalf("无法设置测试数据: %v", err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := c.HGet("bench_hash", "field1")
		if err != nil {
			b.Fatalf("HGet失败: %v", err)
		}
	}
	b.StopTimer()

	cleanupBenchmarkRedisDB(b, c)
}

// 基准测试: ZRange 方法
func BenchmarkRedisCache_ZRange(b *testing.B) {
	c := newBenchmarkRedisCacheInstance(b)
	defer c.Close()

	// 检查Redis连接并跳过测试如果不可用
	skipIfRedisUnavailable(b, c)

	// 清理旧数据并准备测试数据 - 添加100个成员
	cleanupBenchmarkRedisDB(b, c)
	members := make([]cache.Z, 100)
	for i := 0; i < 100; i++ {
		members[i] = cache.Z{
			Score:  float64(i),
			Member: fmt.Sprintf("member%d", i),
		}
	}
	_, err := c.ZAdd("bench_zset", members...)
	if err != nil {
		b.Fatalf("无法设置测试数据: %v", err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := c.ZRange("bench_zset", 0, -1)
		if err != nil {
			b.Fatalf("ZRange失败: %v", err)
		}
	}
	b.StopTimer()

	cleanupBenchmarkRedisDB(b, c)
}
