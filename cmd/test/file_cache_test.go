package cmd

import (
	"os"
	"testing"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/zhoudm1743/go-frame/pkg/cache"
	"github.com/zhoudm1743/go-frame/pkg/config"
)

func TestFileCache(t *testing.T) {
	// 创建临时目录
	tempDir, err := os.MkdirTemp("", "file-cache-test")
	if err != nil {
		t.Fatalf("无法创建临时目录: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// 创建配置
	cfg := &config.Config{
		Cache: config.CacheConfig{
			Type:     "file",
			Prefix:   "test:",
			FilePath: tempDir,
		},
	}

	// 创建日志
	logger := logrus.New()
	logger.SetLevel(logrus.DebugLevel)

	// 创建文件缓存
	fileCache, err := cache.NewFileCache(cfg, logger)
	if err != nil {
		t.Fatalf("无法创建文件缓存: %v", err)
	}
	defer fileCache.Close()

	// 测试基本操作
	t.Run("基本操作", func(t *testing.T) {
		// 设置缓存
		err := fileCache.Set("key1", "value1", 0)
		assert.NoError(t, err)

		// 获取缓存
		value, err := fileCache.Get("key1")
		assert.NoError(t, err)
		assert.Equal(t, "value1", value)

		// 检查键是否存在
		count, err := fileCache.Exists("key1")
		assert.NoError(t, err)
		assert.Equal(t, int64(1), count)

		// 删除缓存
		count, err = fileCache.Del("key1")
		assert.NoError(t, err)
		assert.Equal(t, int64(1), count)

		// 检查键是否已删除
		_, err = fileCache.Get("key1")
		assert.Equal(t, cache.ErrKeyNotFound, err)
	})

	// 测试过期时间
	t.Run("过期时间", func(t *testing.T) {
		// 设置带过期时间的缓存
		err := fileCache.Set("key2", "value2", 1*time.Second)
		assert.NoError(t, err)

		// 获取缓存
		value, err := fileCache.Get("key2")
		assert.NoError(t, err)
		assert.Equal(t, "value2", value)

		// 等待过期
		time.Sleep(1500 * time.Millisecond)

		// 检查键是否已过期
		_, err = fileCache.Get("key2")
		assert.Equal(t, cache.ErrKeyNotFound, err)
	})

	// 测试计数器操作
	t.Run("计数器操作", func(t *testing.T) {
		// 自增
		count, err := fileCache.Incr("counter")
		assert.NoError(t, err)
		assert.Equal(t, int64(1), count)

		// 再次自增
		count, err = fileCache.Incr("counter")
		assert.NoError(t, err)
		assert.Equal(t, int64(2), count)

		// 自减
		count, err = fileCache.Decr("counter")
		assert.NoError(t, err)
		assert.Equal(t, int64(1), count)

		// 按指定值自增
		count, err = fileCache.IncrBy("counter", 10)
		assert.NoError(t, err)
		assert.Equal(t, int64(11), count)
	})

	// 测试哈希表操作
	t.Run("哈希表操作", func(t *testing.T) {
		// 设置哈希表字段
		count, err := fileCache.HSet("hash", "field1", "value1", "field2", "value2")
		assert.NoError(t, err)
		assert.Equal(t, int64(2), count)

		// 获取哈希表字段
		value, err := fileCache.HGet("hash", "field1")
		assert.NoError(t, err)
		assert.Equal(t, "value1", value)

		// 获取所有字段
		fields, err := fileCache.HGetAll("hash")
		assert.NoError(t, err)
		assert.Equal(t, map[string]string{
			"field1": "value1",
			"field2": "value2",
		}, fields)

		// 检查字段是否存在
		exists, err := fileCache.HExists("hash", "field1")
		assert.NoError(t, err)
		assert.True(t, exists)

		// 获取字段数量
		count, err = fileCache.HLen("hash")
		assert.NoError(t, err)
		assert.Equal(t, int64(2), count)

		// 删除字段
		count, err = fileCache.HDel("hash", "field1")
		assert.NoError(t, err)
		assert.Equal(t, int64(1), count)

		// 检查字段是否已删除
		exists, err = fileCache.HExists("hash", "field1")
		assert.NoError(t, err)
		assert.False(t, exists)
	})

	// 测试列表操作
	t.Run("列表操作", func(t *testing.T) {
		// 从左侧插入元素
		count, err := fileCache.LPush("list", "value1", "value2")
		assert.NoError(t, err)
		assert.Equal(t, int64(2), count)

		// 从右侧插入元素
		count, err = fileCache.RPush("list", "value3")
		assert.NoError(t, err)
		assert.Equal(t, int64(3), count)

		// 获取列表长度
		count, err = fileCache.LLen("list")
		assert.NoError(t, err)
		assert.Equal(t, int64(3), count)

		// 获取列表范围
		values, err := fileCache.LRange("list", 0, -1)
		assert.NoError(t, err)
		assert.Equal(t, []string{"value2", "value1", "value3"}, values)

		// 从左侧弹出元素
		value, err := fileCache.LPop("list")
		assert.NoError(t, err)
		assert.Equal(t, "value2", value)

		// 从右侧弹出元素
		value, err = fileCache.RPop("list")
		assert.NoError(t, err)
		assert.Equal(t, "value3", value)

		// 检查列表长度
		count, err = fileCache.LLen("list")
		assert.NoError(t, err)
		assert.Equal(t, int64(1), count)
	})

	// 测试集合操作
	t.Run("集合操作", func(t *testing.T) {
		// 添加集合成员
		count, err := fileCache.SAdd("set", "member1", "member2", "member3")
		assert.NoError(t, err)
		assert.Equal(t, int64(3), count)

		// 获取集合成员数
		count, err = fileCache.SCard("set")
		assert.NoError(t, err)
		assert.Equal(t, int64(3), count)

		// 检查成员是否存在
		exists, err := fileCache.SIsMember("set", "member1")
		assert.NoError(t, err)
		assert.True(t, exists)

		// 获取所有成员
		members, err := fileCache.SMembers("set")
		assert.NoError(t, err)
		assert.ElementsMatch(t, []string{"member1", "member2", "member3"}, members)

		// 删除成员
		count, err = fileCache.SRem("set", "member1")
		assert.NoError(t, err)
		assert.Equal(t, int64(1), count)

		// 检查成员是否已删除
		exists, err = fileCache.SIsMember("set", "member1")
		assert.NoError(t, err)
		assert.False(t, exists)
	})

	// 测试有序集合操作
	t.Run("有序集合操作", func(t *testing.T) {
		// 添加有序集合成员
		count, err := fileCache.ZAdd("zset",
			cache.Z{Score: 1.0, Member: "member1"},
			cache.Z{Score: 2.0, Member: "member2"},
			cache.Z{Score: 3.0, Member: "member3"},
		)
		assert.NoError(t, err)
		assert.Equal(t, int64(3), count)

		// 获取有序集合成员数
		count, err = fileCache.ZCard("zset")
		assert.NoError(t, err)
		assert.Equal(t, int64(3), count)

		// 获取成员分数
		score, err := fileCache.ZScore("zset", "member2")
		assert.NoError(t, err)
		assert.Equal(t, 2.0, score)

		// 获取范围内的成员
		members, err := fileCache.ZRange("zset", 0, -1)
		assert.NoError(t, err)
		assert.Equal(t, []string{"member1", "member2", "member3"}, members)

		// 获取带分数的范围内的成员
		membersWithScores, err := fileCache.ZRangeWithScores("zset", 0, -1)
		assert.NoError(t, err)
		assert.Equal(t, []cache.Z{
			{Score: 1.0, Member: "member1"},
			{Score: 2.0, Member: "member2"},
			{Score: 3.0, Member: "member3"},
		}, membersWithScores)

		// 删除成员
		count, err = fileCache.ZRem("zset", "member2")
		assert.NoError(t, err)
		assert.Equal(t, int64(1), count)

		// 检查成员是否已删除
		_, err = fileCache.ZScore("zset", "member2")
		assert.Equal(t, cache.ErrKeyNotFound, err)
	})

	// 测试持久化
	t.Run("持久化", func(t *testing.T) {
		// 设置缓存
		err := fileCache.Set("persistent", "value", 0)
		assert.NoError(t, err)

		// 关闭缓存
		err = fileCache.Close()
		assert.NoError(t, err)

		// 重新打开缓存
		fileCache, err = cache.NewFileCache(cfg, logger)
		assert.NoError(t, err)

		// 检查缓存是否仍然存在
		value, err := fileCache.Get("persistent")
		assert.NoError(t, err)
		assert.Equal(t, "value", value)
	})

	// 测试键操作
	t.Run("键操作", func(t *testing.T) {
		// 设置多个缓存
		err := fileCache.Set("key1", "value1", 0)
		assert.NoError(t, err)
		err = fileCache.Set("key2", "value2", 0)
		assert.NoError(t, err)
		err = fileCache.Set("otherkey", "value3", 0)
		assert.NoError(t, err)

		// 获取所有键
		keys, err := fileCache.Keys("*")
		assert.NoError(t, err)
		assert.Contains(t, keys, "key1")
		assert.Contains(t, keys, "key2")
		assert.Contains(t, keys, "otherkey")

		// 获取匹配的键
		keys, err = fileCache.Keys("key*")
		assert.NoError(t, err)
		assert.Contains(t, keys, "key1")
		assert.Contains(t, keys, "key2")
		assert.NotContains(t, keys, "otherkey")
	})

	// 测试 Ping
	t.Run("Ping", func(t *testing.T) {
		err := fileCache.Ping()
		assert.NoError(t, err)
	})
}
