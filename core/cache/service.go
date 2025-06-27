package cache

import (
	"context"
	"fmt"
	"time"

	"github.com/zhoudm1743/go-flow/core/logger"
	"go.uber.org/fx"
)

// UserCacheService 用户缓存服务示例
type UserCacheService struct {
	helper *CacheHelper
	logger logger.Logger
}

// NewUserCacheService 创建用户缓存服务
func NewUserCacheService(cache Cache, log logger.Logger) *UserCacheService {
	helper := NewCacheHelper(cache, log, "user")
	return &UserCacheService{
		helper: helper,
		logger: log,
	}
}

// CacheUser 缓存用户信息
func (s *UserCacheService) CacheUser(ctx context.Context, userID uint, user interface{}) error {
	key := fmt.Sprintf("info:%d", userID)
	return s.helper.SetJSONCtx(ctx, key, user, time.Hour*24) // 缓存24小时
}

// GetCachedUser 获取缓存的用户信息
func (s *UserCacheService) GetCachedUser(ctx context.Context, userID uint, dest interface{}) error {
	key := fmt.Sprintf("info:%d", userID)
	return s.helper.GetJSONCtx(ctx, key, dest)
}

// CacheUserSession 缓存用户会话
func (s *UserCacheService) CacheUserSession(ctx context.Context, sessionID string, userID uint) error {
	key := fmt.Sprintf("session:%s", sessionID)
	return s.helper.cache.SetCtx(ctx, s.helper.buildKey(key), userID, time.Hour*2) // 会话2小时过期
}

// GetUserFromSession 从会话获取用户ID
func (s *UserCacheService) GetUserFromSession(ctx context.Context, sessionID string) (string, error) {
	key := fmt.Sprintf("session:%s", sessionID)
	return s.helper.cache.GetCtx(ctx, s.helper.buildKey(key))
}

// DeleteUserSession 删除用户会话
func (s *UserCacheService) DeleteUserSession(ctx context.Context, sessionID string) error {
	key := fmt.Sprintf("session:%s", sessionID)
	_, err := s.helper.cache.DelCtx(ctx, s.helper.buildKey(key))
	return err
}

// IncrementLoginCount 增加登录次数
func (s *UserCacheService) IncrementLoginCount(ctx context.Context, userID uint) (int64, error) {
	key := fmt.Sprintf("login_count:%d", userID)
	return s.helper.cache.IncrCtx(ctx, s.helper.buildKey(key))
}

// GetLoginCount 获取登录次数
func (s *UserCacheService) GetLoginCount(ctx context.Context, userID uint) (string, error) {
	key := fmt.Sprintf("login_count:%d", userID)
	return s.helper.cache.GetCtx(ctx, s.helper.buildKey(key))
}

// CacheUserPreferences 缓存用户偏好设置
func (s *UserCacheService) CacheUserPreferences(ctx context.Context, userID uint, preferences map[string]string) error {
	key := fmt.Sprintf("preferences:%d", userID)
	fullKey := s.helper.buildKey(key)

	// 使用哈希存储偏好设置
	values := make([]interface{}, 0, len(preferences)*2)
	for k, v := range preferences {
		values = append(values, k, v)
	}

	_, err := s.helper.cache.HSetCtx(ctx, fullKey, values...)
	if err == nil {
		// 设置过期时间
		s.helper.cache.ExpireCtx(ctx, fullKey, time.Hour*24*7) // 7天过期
	}
	return err
}

// GetUserPreferences 获取用户偏好设置
func (s *UserCacheService) GetUserPreferences(ctx context.Context, userID uint) (map[string]string, error) {
	key := fmt.Sprintf("preferences:%d", userID)
	return s.helper.cache.HGetAllCtx(ctx, s.helper.buildKey(key))
}

// SetUserPreference 设置单个用户偏好
func (s *UserCacheService) SetUserPreference(ctx context.Context, userID uint, key, value string) error {
	hashKey := fmt.Sprintf("preferences:%d", userID)
	fullKey := s.helper.buildKey(hashKey)
	_, err := s.helper.cache.HSetCtx(ctx, fullKey, key, value)
	return err
}

// GetUserPreference 获取单个用户偏好
func (s *UserCacheService) GetUserPreference(ctx context.Context, userID uint, key string) (string, error) {
	hashKey := fmt.Sprintf("preferences:%d", userID)
	return s.helper.cache.HGetCtx(ctx, s.helper.buildKey(hashKey), key)
}

// AddToUserQueue 添加到用户队列
func (s *UserCacheService) AddToUserQueue(ctx context.Context, userID uint, message string) error {
	key := fmt.Sprintf("queue:%d", userID)
	_, err := s.helper.cache.RPushCtx(ctx, s.helper.buildKey(key), message)
	return err
}

// PopFromUserQueue 从用户队列弹出消息
func (s *UserCacheService) PopFromUserQueue(ctx context.Context, userID uint) (string, error) {
	key := fmt.Sprintf("queue:%d", userID)
	return s.helper.cache.LPopCtx(ctx, s.helper.buildKey(key))
}

// GetUserQueueLength 获取用户队列长度
func (s *UserCacheService) GetUserQueueLength(ctx context.Context, userID uint) (int64, error) {
	key := fmt.Sprintf("queue:%d", userID)
	return s.helper.cache.LLenCtx(ctx, s.helper.buildKey(key))
}

// AddUserTag 添加用户标签
func (s *UserCacheService) AddUserTag(ctx context.Context, userID uint, tags ...string) error {
	key := fmt.Sprintf("tags:%d", userID)
	fullKey := s.helper.buildKey(key)

	// 转换为 interface{} 切片
	values := make([]interface{}, len(tags))
	for i, tag := range tags {
		values[i] = tag
	}

	_, err := s.helper.cache.SAddCtx(ctx, fullKey, values...)
	return err
}

// RemoveUserTag 移除用户标签
func (s *UserCacheService) RemoveUserTag(ctx context.Context, userID uint, tags ...string) error {
	key := fmt.Sprintf("tags:%d", userID)
	fullKey := s.helper.buildKey(key)

	// 转换为 interface{} 切片
	values := make([]interface{}, len(tags))
	for i, tag := range tags {
		values[i] = tag
	}

	_, err := s.helper.cache.SRemCtx(ctx, fullKey, values...)
	return err
}

// GetUserTags 获取用户所有标签
func (s *UserCacheService) GetUserTags(ctx context.Context, userID uint) ([]string, error) {
	key := fmt.Sprintf("tags:%d", userID)
	return s.helper.cache.SMembersCtx(ctx, s.helper.buildKey(key))
}

// HasUserTag 检查用户是否有某个标签
func (s *UserCacheService) HasUserTag(ctx context.Context, userID uint, tag string) (bool, error) {
	key := fmt.Sprintf("tags:%d", userID)
	return s.helper.cache.SIsMemberCtx(ctx, s.helper.buildKey(key), tag)
}

// 简化版方法（不需要 Context）

// SimpleCacheUser 缓存用户信息（简化版）
func (s *UserCacheService) SimpleCacheUser(userID uint, user interface{}) error {
	key := fmt.Sprintf("info:%d", userID)
	return s.helper.SetJSON(key, user, time.Hour*24) // 缓存24小时
}

// SimpleGetCachedUser 获取缓存的用户信息（简化版）
func (s *UserCacheService) SimpleGetCachedUser(userID uint, dest interface{}) error {
	key := fmt.Sprintf("info:%d", userID)
	return s.helper.GetJSON(key, dest)
}

// SimpleCacheUserSession 缓存用户会话（简化版）
func (s *UserCacheService) SimpleCacheUserSession(sessionID string, userID uint) error {
	key := fmt.Sprintf("session:%s", sessionID)
	return s.helper.cache.Set(s.helper.buildKey(key), userID, time.Hour*2) // 会话2小时过期
}

// SimpleGetUserFromSession 从会话获取用户ID（简化版）
func (s *UserCacheService) SimpleGetUserFromSession(sessionID string) (string, error) {
	key := fmt.Sprintf("session:%s", sessionID)
	return s.helper.cache.Get(s.helper.buildKey(key))
}

// SimpleDeleteUserSession 删除用户会话（简化版）
func (s *UserCacheService) SimpleDeleteUserSession(sessionID string) error {
	key := fmt.Sprintf("session:%s", sessionID)
	_, err := s.helper.cache.Del(s.helper.buildKey(key))
	return err
}

// SimpleIncrementLoginCount 增加登录次数（简化版）
func (s *UserCacheService) SimpleIncrementLoginCount(userID uint) (int64, error) {
	key := fmt.Sprintf("login_count:%d", userID)
	return s.helper.cache.Incr(s.helper.buildKey(key))
}

// SimpleGetLoginCount 获取登录次数（简化版）
func (s *UserCacheService) SimpleGetLoginCount(userID uint) (string, error) {
	key := fmt.Sprintf("login_count:%d", userID)
	return s.helper.cache.Get(s.helper.buildKey(key))
}

// CacheServiceModule 缓存服务模块
var CacheServiceModule = fx.Options(
	fx.Provide(NewUserCacheService),
)
