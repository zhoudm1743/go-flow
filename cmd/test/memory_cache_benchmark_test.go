package cmd

import (
	"fmt"
	"testing"
	"time"

	"github.com/zhoudm1743/go-frame/pkg/cache"
	"github.com/zhoudm1743/go-frame/pkg/config"
	"github.com/zhoudm1743/go-frame/pkg/log"
)

// 创建基准测试用的内存缓存实例
func newBenchmarkMemoryCache(b *testing.B) cache.Cache {
	// 创建配置
	cfg := &config.Config{}
	cfg.Cache.Type = "memory"
	cfg.Cache.Prefix = "bench:"
	cfg.Log.Level = "error" // 降低日志级别，避免影响性能测试
	cfg.Log.Format = "text"
	cfg.Log.OutputPath = "stdout"

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

// BenchmarkMemoryCacheGet 基准测试Get操作
func BenchmarkMemoryCacheGet(b *testing.B) {
	c := newBenchmarkMemoryCache(b)
	defer c.Close()

	// 预先设置一些键
	for i := 0; i < 1000; i++ {
		key := fmt.Sprintf("key:%d", i)
		value := fmt.Sprintf("value:%d", i)
		if err := c.Set(key, value, time.Hour); err != nil {
			b.Fatalf("设置缓存失败: %v", err)
		}
	}

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		counter := 0
		for pb.Next() {
			key := fmt.Sprintf("key:%d", counter%1000)
			_, err := c.Get(key)
			if err != nil {
				b.Fatalf("获取缓存失败: %v", err)
			}
			counter++
		}
	})
}

// BenchmarkMemoryCacheSet 基准测试Set操作
func BenchmarkMemoryCacheSet(b *testing.B) {
	c := newBenchmarkMemoryCache(b)
	defer c.Close()

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		counter := 0
		for pb.Next() {
			key := fmt.Sprintf("key:%d", counter)
			value := fmt.Sprintf("value:%d", counter)
			if err := c.Set(key, value, time.Hour); err != nil {
				b.Fatalf("设置缓存失败: %v", err)
			}
			counter++
		}
	})
}

// BenchmarkMemoryCacheHSet 基准测试HSet操作
func BenchmarkMemoryCacheHSet(b *testing.B) {
	c := newBenchmarkMemoryCache(b)
	defer c.Close()

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		counter := 0
		for pb.Next() {
			key := fmt.Sprintf("hash:%d", counter/100)
			field := fmt.Sprintf("field:%d", counter%100)
			value := fmt.Sprintf("value:%d", counter)
			_, err := c.HSet(key, field, value)
			if err != nil {
				b.Fatalf("设置哈希缓存失败: %v", err)
			}
			counter++
		}
	})
}

// BenchmarkMemoryCacheHGet 基准测试HGet操作
func BenchmarkMemoryCacheHGet(b *testing.B) {
	c := newBenchmarkMemoryCache(b)
	defer c.Close()

	// 预先设置一些哈希键
	for i := 0; i < 10; i++ {
		key := fmt.Sprintf("hash:%d", i)
		for j := 0; j < 100; j++ {
			field := fmt.Sprintf("field:%d", j)
			value := fmt.Sprintf("value:%d", i*100+j)
			_, err := c.HSet(key, field, value)
			if err != nil {
				b.Fatalf("设置哈希缓存失败: %v", err)
			}
		}
	}

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		counter := 0
		for pb.Next() {
			key := fmt.Sprintf("hash:%d", counter/100%10)
			field := fmt.Sprintf("field:%d", counter%100)
			_, err := c.HGet(key, field)
			if err != nil {
				b.Fatalf("获取哈希缓存失败: %v", err)
			}
			counter++
		}
	})
}

// BenchmarkMemoryCacheZAdd 基准测试ZAdd操作
func BenchmarkMemoryCacheZAdd(b *testing.B) {
	c := newBenchmarkMemoryCache(b)
	defer c.Close()

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		counter := 0
		for pb.Next() {
			key := fmt.Sprintf("zset:%d", counter/100)
			member := fmt.Sprintf("member:%d", counter%100)
			_, err := c.ZAdd(key, cache.Z{Score: float64(counter), Member: member})
			if err != nil {
				b.Fatalf("设置有序集合失败: %v", err)
			}
			counter++
		}
	})
}

// BenchmarkMemoryCacheZRange 基准测试ZRange操作
func BenchmarkMemoryCacheZRange(b *testing.B) {
	c := newBenchmarkMemoryCache(b)
	defer c.Close()

	// 预先设置一些有序集合
	for i := 0; i < 10; i++ {
		key := fmt.Sprintf("zset:%d", i)
		for j := 0; j < 100; j++ {
			member := fmt.Sprintf("member:%d", j)
			_, err := c.ZAdd(key, cache.Z{Score: float64(j), Member: member})
			if err != nil {
				b.Fatalf("设置有序集合失败: %v", err)
			}
		}
	}

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		counter := 0
		for pb.Next() {
			key := fmt.Sprintf("zset:%d", counter%10)
			_, err := c.ZRange(key, 0, 10) // 只获取前10个
			if err != nil {
				b.Fatalf("获取有序集合失败: %v", err)
			}
			counter++
		}
	})
}

// BenchmarkMemoryCacheLPush 基准测试LPush操作
func BenchmarkMemoryCacheLPush(b *testing.B) {
	c := newBenchmarkMemoryCache(b)
	defer c.Close()

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		counter := 0
		for pb.Next() {
			key := fmt.Sprintf("list:%d", counter/100)
			value := fmt.Sprintf("value:%d", counter)
			_, err := c.LPush(key, value)
			if err != nil {
				b.Fatalf("设置列表失败: %v", err)
			}
			counter++
		}
	})
}

// BenchmarkMemoryCacheLRange 基准测试LRange操作
func BenchmarkMemoryCacheLRange(b *testing.B) {
	c := newBenchmarkMemoryCache(b)
	defer c.Close()

	// 预先设置一些列表
	for i := 0; i < 10; i++ {
		key := fmt.Sprintf("list:%d", i)
		for j := 0; j < 100; j++ {
			value := fmt.Sprintf("value:%d", j)
			_, err := c.LPush(key, value)
			if err != nil {
				b.Fatalf("设置列表失败: %v", err)
			}
		}
	}

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		counter := 0
		for pb.Next() {
			key := fmt.Sprintf("list:%d", counter%10)
			_, err := c.LRange(key, 0, 10) // 只获取前10个
			if err != nil {
				b.Fatalf("获取列表失败: %v", err)
			}
			counter++
		}
	})
}

// BenchmarkMemoryCacheSAdd 基准测试SAdd操作
func BenchmarkMemoryCacheSAdd(b *testing.B) {
	c := newBenchmarkMemoryCache(b)
	defer c.Close()

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		counter := 0
		for pb.Next() {
			key := fmt.Sprintf("set:%d", counter/100)
			member := fmt.Sprintf("member:%d", counter%100)
			_, err := c.SAdd(key, member)
			if err != nil {
				b.Fatalf("设置集合失败: %v", err)
			}
			counter++
		}
	})
}

// BenchmarkMemoryCacheSMembers 基准测试SMembers操作
func BenchmarkMemoryCacheSMembers(b *testing.B) {
	c := newBenchmarkMemoryCache(b)
	defer c.Close()

	// 预先设置一些集合
	for i := 0; i < 10; i++ {
		key := fmt.Sprintf("set:%d", i)
		for j := 0; j < 100; j++ {
			member := fmt.Sprintf("member:%d", j)
			_, err := c.SAdd(key, member)
			if err != nil {
				b.Fatalf("设置集合失败: %v", err)
			}
		}
	}

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		counter := 0
		for pb.Next() {
			key := fmt.Sprintf("set:%d", counter%10)
			_, err := c.SMembers(key)
			if err != nil {
				b.Fatalf("获取集合成员失败: %v", err)
			}
			counter++
		}
	})
}
