package cmd

import (
	"fmt"
	"testing"

	"github.com/zhoudm1743/go-frame/pkg/cache"
	"github.com/zhoudm1743/go-frame/pkg/config"
	"github.com/zhoudm1743/go-frame/pkg/log"
)

// 创建基准测试用的内存缓存实例
func newSimpleMemoryCache(b *testing.B) cache.Cache {
	// 创建配置
	cfg := &config.Config{}
	cfg.Cache.Type = "memory"
	cfg.Cache.Prefix = "bench:"
	cfg.Log.Level = "error" // 降低日志级别，避免影响性能测试
	cfg.Log.Format = "text"
	cfg.Log.OutputPath = "none"

	// 创建日志
	logger, _ := log.NewLogger(log.LoggerParams{
		Config: cfg,
	})

	// 创建缓存
	memCache, err := cache.NewMemoryCache(cfg, logger)
	if err != nil {
		b.Fatalf("无法创建内存缓存: %v", err)
	}

	return memCache
}

// BenchmarkMemorySimpleGet 基准测试Get操作
func BenchmarkMemorySimpleGet(b *testing.B) {
	c := newSimpleMemoryCache(b)
	defer c.Close()

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

// BenchmarkMemorySimpleGetParallel 并发基准测试Get操作
func BenchmarkMemorySimpleGetParallel(b *testing.B) {
	c := newSimpleMemoryCache(b)
	defer c.Close()

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

// BenchmarkMemorySimpleSet 基准测试Set操作
func BenchmarkMemorySimpleSet(b *testing.B) {
	c := newSimpleMemoryCache(b)
	defer c.Close()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		err := c.Set("test-key", "test-value", 0)
		if err != nil {
			b.Fatalf("设置缓存失败: %v", err)
		}
	}
}

// BenchmarkMemorySimpleSetParallel 并发基准测试Set操作
func BenchmarkMemorySimpleSetParallel(b *testing.B) {
	c := newSimpleMemoryCache(b)
	defer c.Close()

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
