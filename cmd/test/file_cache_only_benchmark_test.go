package cmd

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/zhoudm1743/go-frame/pkg/cache"
	"github.com/zhoudm1743/go-frame/pkg/config"
	"github.com/zhoudm1743/go-frame/pkg/log"
)

// 创建基准测试用的文件缓存实例
func newFileCacheBenchmark(b *testing.B) cache.Cache {
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

// BenchmarkFileCacheBasicGet 基准测试基本Get操作
func BenchmarkFileCacheBasicGet(b *testing.B) {
	c := newFileCacheBenchmark(b)
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

// BenchmarkFileCacheBasicGetParallel 并发基准测试基本Get操作
func BenchmarkFileCacheBasicGetParallel(b *testing.B) {
	c := newFileCacheBenchmark(b)
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

// BenchmarkFileCacheBasicSet 基准测试基本Set操作
func BenchmarkFileCacheBasicSet(b *testing.B) {
	c := newFileCacheBenchmark(b)
	defer c.Close()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		err := c.Set("test-key", "test-value", 0)
		if err != nil {
			b.Fatalf("设置缓存失败: %v", err)
		}
	}
}

// BenchmarkFileCacheBasicSetParallel 并发基准测试基本Set操作
func BenchmarkFileCacheBasicSetParallel(b *testing.B) {
	c := newFileCacheBenchmark(b)
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

// BenchmarkFileCacheMultiGet 基准测试多键Get操作
func BenchmarkFileCacheMultiGet(b *testing.B) {
	c := newFileCacheBenchmark(b)
	defer c.Close()

	// 预先设置多个键
	for i := 0; i < 1000; i++ {
		key := fmt.Sprintf("key:%d", i)
		value := fmt.Sprintf("value:%d", i)
		if err := c.Set(key, value, 0); err != nil {
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

// BenchmarkFileCacheMultiGetParallel 并发基准测试多键Get操作
func BenchmarkFileCacheMultiGetParallel(b *testing.B) {
	c := newFileCacheBenchmark(b)
	defer c.Close()

	// 预先设置多个键
	for i := 0; i < 1000; i++ {
		key := fmt.Sprintf("key:%d", i)
		value := fmt.Sprintf("value:%d", i)
		if err := c.Set(key, value, 0); err != nil {
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

// BenchmarkFileCacheHSetGet 基准测试哈希表操作
func BenchmarkFileCacheHSetGet(b *testing.B) {
	c := newFileCacheBenchmark(b)
	defer c.Close()

	// 预先设置哈希表
	for i := 0; i < 100; i++ {
		field := fmt.Sprintf("field:%d", i)
		value := fmt.Sprintf("value:%d", i)
		_, err := c.HSet("hash:test", field, value)
		if err != nil {
			b.Fatalf("设置哈希缓存失败: %v", err)
		}
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		field := fmt.Sprintf("field:%d", i%100)
		_, err := c.HGet("hash:test", field)
		if err != nil {
			b.Fatalf("获取哈希缓存失败: %v", err)
		}
	}
}

// BenchmarkFileCacheHSetGetParallel 并发基准测试哈希表操作
func BenchmarkFileCacheHSetGetParallel(b *testing.B) {
	c := newFileCacheBenchmark(b)
	defer c.Close()

	// 预先设置哈希表
	for i := 0; i < 100; i++ {
		field := fmt.Sprintf("field:%d", i)
		value := fmt.Sprintf("value:%d", i)
		_, err := c.HSet("hash:test", field, value)
		if err != nil {
			b.Fatalf("设置哈希缓存失败: %v", err)
		}
	}

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		counter := 0
		for pb.Next() {
			field := fmt.Sprintf("field:%d", counter%100)
			_, err := c.HGet("hash:test", field)
			if err != nil {
				b.Fatalf("获取哈希缓存失败: %v", err)
			}
			counter++
		}
	})
}

// BenchmarkFileCacheConcurrentMixed 并发基准测试混合操作
func BenchmarkFileCacheConcurrentMixed(b *testing.B) {
	c := newFileCacheBenchmark(b)
	defer c.Close()

	// 预先设置一些数据
	for i := 0; i < 100; i++ {
		// 普通键值对
		key := fmt.Sprintf("key:%d", i)
		value := fmt.Sprintf("value:%d", i)
		if err := c.Set(key, value, 0); err != nil {
			b.Fatalf("设置缓存失败: %v", err)
		}

		// 哈希表
		field := fmt.Sprintf("field:%d", i)
		if _, err := c.HSet("hash:test", field, value); err != nil {
			b.Fatalf("设置哈希缓存失败: %v", err)
		}

		// 列表
		if _, err := c.LPush("list:test", value); err != nil {
			b.Fatalf("设置列表失败: %v", err)
		}

		// 集合
		if _, err := c.SAdd("set:test", value); err != nil {
			b.Fatalf("设置集合失败: %v", err)
		}

		// 有序集合
		if _, err := c.ZAdd("zset:test", cache.Z{Score: float64(i), Member: value}); err != nil {
			b.Fatalf("设置有序集合失败: %v", err)
		}
	}

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		counter := 0
		for pb.Next() {
			// 每次迭代执行不同的操作
			switch counter % 5 {
			case 0: // Get
				key := fmt.Sprintf("key:%d", counter%100)
				_, err := c.Get(key)
				if err != nil {
					b.Fatalf("获取缓存失败: %v", err)
				}
			case 1: // HGet
				field := fmt.Sprintf("field:%d", counter%100)
				_, err := c.HGet("hash:test", field)
				if err != nil {
					b.Fatalf("获取哈希缓存失败: %v", err)
				}
			case 2: // LRange
				_, err := c.LRange("list:test", 0, 10)
				if err != nil {
					b.Fatalf("获取列表失败: %v", err)
				}
			case 3: // SMembers
				_, err := c.SMembers("set:test")
				if err != nil {
					b.Fatalf("获取集合成员失败: %v", err)
				}
			case 4: // ZRange
				_, err := c.ZRange("zset:test", 0, 10)
				if err != nil {
					b.Fatalf("获取有序集合失败: %v", err)
				}
			}
			counter++
		}
	})
}

// BenchmarkFileCacheConcurrentReadWrite 并发基准测试读写操作
func BenchmarkFileCacheConcurrentReadWrite(b *testing.B) {
	c := newFileCacheBenchmark(b)
	defer c.Close()

	// 预先设置一些数据
	for i := 0; i < 100; i++ {
		key := fmt.Sprintf("key:%d", i)
		value := fmt.Sprintf("value:%d", i)
		if err := c.Set(key, value, 0); err != nil {
			b.Fatalf("设置缓存失败: %v", err)
		}
	}

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		counter := 0
		for pb.Next() {
			if counter%2 == 0 { // 读操作
				key := fmt.Sprintf("key:%d", counter%100)
				_, err := c.Get(key)
				if err != nil {
					b.Fatalf("获取缓存失败: %v", err)
				}
			} else { // 写操作
				key := fmt.Sprintf("key:%d", counter%100)
				value := fmt.Sprintf("new-value:%d", counter)
				err := c.Set(key, value, 0)
				if err != nil {
					b.Fatalf("设置缓存失败: %v", err)
				}
			}
			counter++
		}
	})
}
