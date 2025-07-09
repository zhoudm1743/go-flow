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
func newSimpleFileCache(b *testing.B) cache.Cache {
	// 创建临时文件路径
	tempFile := fmt.Sprintf("simple_benchmark_file_cache_%d.db", time.Now().UnixNano())

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

// BenchmarkSimpleGet 基准测试Get操作
func BenchmarkSimpleGet(b *testing.B) {
	c := newSimpleFileCache(b)
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

// BenchmarkSimpleGetParallel 并发基准测试Get操作
func BenchmarkSimpleGetParallel(b *testing.B) {
	c := newSimpleFileCache(b)
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

// BenchmarkSimpleSet 基准测试Set操作
func BenchmarkSimpleSet(b *testing.B) {
	c := newSimpleFileCache(b)
	defer c.Close()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		err := c.Set("test-key", "test-value", 0)
		if err != nil {
			b.Fatalf("设置缓存失败: %v", err)
		}
	}
}

// BenchmarkSimpleSetParallel 并发基准测试Set操作
func BenchmarkSimpleSetParallel(b *testing.B) {
	c := newSimpleFileCache(b)
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
