package cmd

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/zhoudm1743/go-flow/pkg/cache"
	"github.com/zhoudm1743/go-flow/pkg/config"
	"github.com/zhoudm1743/go-flow/pkg/log"
)

// TestFileCacheBenchmarkSetup 确保基准测试可以运行
func TestFileCacheBenchmarkSetup(t *testing.T) {
	// 这个测试函数什么都不做，只是确保基准测试可以运行
	t.Log("文件缓存基准测试设置完成")
}

// 创建基准测试用的文件缓存实例
func newBenchmarkFileCache(b *testing.B) cache.Cache {
	// 创建临时文件路径
	tempFile := fmt.Sprintf("benchmark_file_cache_%d.db", time.Now().UnixNano())

	// 创建配置
	cfg := &config.Config{}
	cfg.Cache.Type = "file"
	cfg.Cache.Prefix = "bench:"
	cfg.Cache.FilePath = tempFile
	cfg.Log.Level = "error" // 降低日志级别，避免影响性能测试
	cfg.Log.Format = "text"
	cfg.Log.OutputPath = "none"

	// 创建日志
	logger, _ := log.NewLogger(log.LoggerParams{
		Config: cfg,
	})

	// 创建缓存
	fileCache, err := cache.NewFileCache(cfg, logger)
	if err != nil {
		b.Fatalf("无法创建文件缓存: %v", err)
	}

	// 返回缓存实例，并在测试结束后清理文件
	b.Cleanup(func() {
		fileCache.Close()
		os.Remove(tempFile)
	})

	return fileCache
}

// 清理基准测试数据
func cleanupBenchmarkFileCache(b *testing.B, c cache.Cache) {
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

// BenchmarkFileCacheGet 基准测试Get操作
func BenchmarkFileCacheGet(b *testing.B) {
	c := newBenchmarkFileCache(b)
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
	for i := 0; i < b.N; i++ {
		key := fmt.Sprintf("key:%d", i%1000)
		_, err := c.Get(key)
		if err != nil {
			b.Fatalf("获取缓存失败: %v", err)
		}
	}
}

// BenchmarkFileCacheGetParallel 并发基准测试Get操作
func BenchmarkFileCacheGetParallel(b *testing.B) {
	c := newBenchmarkFileCache(b)
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

// BenchmarkFileCacheSet 基准测试Set操作
func BenchmarkFileCacheSet(b *testing.B) {
	c := newBenchmarkFileCache(b)
	defer c.Close()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		key := fmt.Sprintf("key:%d", i)
		value := fmt.Sprintf("value:%d", i)
		if err := c.Set(key, value, time.Hour); err != nil {
			b.Fatalf("设置缓存失败: %v", err)
		}
	}
}

// BenchmarkFileCacheSetParallel 并发基准测试Set操作
func BenchmarkFileCacheSetParallel(b *testing.B) {
	c := newBenchmarkFileCache(b)
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

// BenchmarkFileCacheHSet 基准测试HSet操作
func BenchmarkFileCacheHSet(b *testing.B) {
	c := newBenchmarkFileCache(b)
	defer c.Close()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		key := fmt.Sprintf("hash:%d", i/100)
		field := fmt.Sprintf("field:%d", i%100)
		value := fmt.Sprintf("value:%d", i)
		_, err := c.HSet(key, field, value)
		if err != nil {
			b.Fatalf("设置哈希缓存失败: %v", err)
		}
	}
}

// BenchmarkFileCacheHSetParallel 并发基准测试HSet操作
func BenchmarkFileCacheHSetParallel(b *testing.B) {
	c := newBenchmarkFileCache(b)
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

// BenchmarkFileCacheHGet 基准测试HGet操作
func BenchmarkFileCacheHGet(b *testing.B) {
	c := newBenchmarkFileCache(b)
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
	for i := 0; i < b.N; i++ {
		key := fmt.Sprintf("hash:%d", i/100%10)
		field := fmt.Sprintf("field:%d", i%100)
		_, err := c.HGet(key, field)
		if err != nil {
			b.Fatalf("获取哈希缓存失败: %v", err)
		}
	}
}

// BenchmarkFileCacheHGetParallel 并发基准测试HGet操作
func BenchmarkFileCacheHGetParallel(b *testing.B) {
	c := newBenchmarkFileCache(b)
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

// BenchmarkFileCacheLPush 基准测试LPush操作
func BenchmarkFileCacheLPush(b *testing.B) {
	c := newBenchmarkFileCache(b)
	defer c.Close()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		key := fmt.Sprintf("list:%d", i/100)
		value := fmt.Sprintf("value:%d", i)
		_, err := c.LPush(key, value)
		if err != nil {
			b.Fatalf("设置列表失败: %v", err)
		}
	}
}

// BenchmarkFileCacheLPushParallel 并发基准测试LPush操作
func BenchmarkFileCacheLPushParallel(b *testing.B) {
	c := newBenchmarkFileCache(b)
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

// BenchmarkFileCacheLRange 基准测试LRange操作
func BenchmarkFileCacheLRange(b *testing.B) {
	c := newBenchmarkFileCache(b)
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
	for i := 0; i < b.N; i++ {
		key := fmt.Sprintf("list:%d", i%10)
		_, err := c.LRange(key, 0, 10) // 只获取前10个
		if err != nil {
			b.Fatalf("获取列表失败: %v", err)
		}
	}
}

// BenchmarkFileCacheLRangeParallel 并发基准测试LRange操作
func BenchmarkFileCacheLRangeParallel(b *testing.B) {
	c := newBenchmarkFileCache(b)
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

// BenchmarkFileCacheSAdd 基准测试SAdd操作
func BenchmarkFileCacheSAdd(b *testing.B) {
	c := newBenchmarkFileCache(b)
	defer c.Close()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		key := fmt.Sprintf("set:%d", i/100)
		member := fmt.Sprintf("member:%d", i%100)
		_, err := c.SAdd(key, member)
		if err != nil {
			b.Fatalf("设置集合失败: %v", err)
		}
	}
}

// BenchmarkFileCacheSAddParallel 并发基准测试SAdd操作
func BenchmarkFileCacheSAddParallel(b *testing.B) {
	c := newBenchmarkFileCache(b)
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

// BenchmarkFileCacheSMembers 基准测试SMembers操作
func BenchmarkFileCacheSMembers(b *testing.B) {
	c := newBenchmarkFileCache(b)
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
	for i := 0; i < b.N; i++ {
		key := fmt.Sprintf("set:%d", i%10)
		_, err := c.SMembers(key)
		if err != nil {
			b.Fatalf("获取集合成员失败: %v", err)
		}
	}
}

// BenchmarkFileCacheSMembersParallel 并发基准测试SMembers操作
func BenchmarkFileCacheSMembersParallel(b *testing.B) {
	c := newBenchmarkFileCache(b)
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

// BenchmarkFileCacheZAdd 基准测试ZAdd操作
func BenchmarkFileCacheZAdd(b *testing.B) {
	c := newBenchmarkFileCache(b)
	defer c.Close()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		key := fmt.Sprintf("zset:%d", i/100)
		member := fmt.Sprintf("member:%d", i%100)
		_, err := c.ZAdd(key, cache.Z{Score: float64(i), Member: member})
		if err != nil {
			b.Fatalf("设置有序集合失败: %v", err)
		}
	}
}

// BenchmarkFileCacheZAddParallel 并发基准测试ZAdd操作
func BenchmarkFileCacheZAddParallel(b *testing.B) {
	c := newBenchmarkFileCache(b)
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

// BenchmarkFileCacheZRange 基准测试ZRange操作
func BenchmarkFileCacheZRange(b *testing.B) {
	c := newBenchmarkFileCache(b)
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
	for i := 0; i < b.N; i++ {
		key := fmt.Sprintf("zset:%d", i%10)
		_, err := c.ZRange(key, 0, 10) // 只获取前10个
		if err != nil {
			b.Fatalf("获取有序集合失败: %v", err)
		}
	}
}

// BenchmarkFileCacheZRangeParallel 并发基准测试ZRange操作
func BenchmarkFileCacheZRangeParallel(b *testing.B) {
	c := newBenchmarkFileCache(b)
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
