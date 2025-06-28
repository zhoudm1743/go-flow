package concurrency

import (
	"context"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/zhoudm1743/go-flow/core/errors"
	"github.com/zhoudm1743/go-flow/core/logger"
)

// RateLimiter 限流器接口
type RateLimiter interface {
	Allow() bool
	AllowN(n int) bool
	Wait(ctx context.Context) error
	Burst() int
	Limit() float64
}

// TokenBucketLimiter 令牌桶限流器
type TokenBucketLimiter struct {
	rps        float64    // 每秒生成的令牌数
	burst      int        // 桶的容量
	tokens     float64    // 当前令牌数
	lastUpdate time.Time  // 最后更新时间
	mutex      sync.Mutex // 保护并发访问
	logger     logger.Logger
}

// NewTokenBucketLimiter 创建令牌桶限流器
func NewTokenBucketLimiter(rps float64, burst int, log logger.Logger) *TokenBucketLimiter {
	now := time.Now()
	return &TokenBucketLimiter{
		rps:        rps,
		burst:      burst,
		tokens:     float64(burst), // 初始时令牌桶是满的
		lastUpdate: now,
		logger:     log,
	}
}

// updateTokens 更新令牌数量
func (t *TokenBucketLimiter) updateTokens(now time.Time) {
	elapsed := now.Sub(t.lastUpdate).Seconds()
	t.tokens += elapsed * t.rps
	if t.tokens > float64(t.burst) {
		t.tokens = float64(t.burst)
	}
	t.lastUpdate = now
}

func (t *TokenBucketLimiter) Allow() bool {
	return t.AllowN(1)
}

func (t *TokenBucketLimiter) AllowN(n int) bool {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	now := time.Now()
	t.updateTokens(now)

	if t.tokens >= float64(n) {
		t.tokens -= float64(n)
		return true
	}
	return false
}

func (t *TokenBucketLimiter) Wait(ctx context.Context) error {
	for {
		if t.Allow() {
			return nil
		}

		// 计算等待时间
		sleepTime := time.Second / time.Duration(t.rps)

		select {
		case <-time.After(sleepTime):
			continue
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

func (t *TokenBucketLimiter) Burst() int {
	return t.burst
}

func (t *TokenBucketLimiter) Limit() float64 {
	return t.rps
}

// SlidingWindowLimiter 滑动窗口限流器
type SlidingWindowLimiter struct {
	requests   []time.Time
	maxRequest int
	window     time.Duration
	mutex      sync.Mutex
	logger     logger.Logger
}

// NewSlidingWindowLimiter 创建滑动窗口限流器
func NewSlidingWindowLimiter(maxRequest int, window time.Duration, log logger.Logger) *SlidingWindowLimiter {
	return &SlidingWindowLimiter{
		requests:   make([]time.Time, 0),
		maxRequest: maxRequest,
		window:     window,
		logger:     log,
	}
}

func (s *SlidingWindowLimiter) Allow() bool {
	return s.AllowN(1)
}

func (s *SlidingWindowLimiter) AllowN(n int) bool {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	now := time.Now()
	// 清理过期请求
	cutoff := now.Add(-s.window)
	for len(s.requests) > 0 && s.requests[0].Before(cutoff) {
		s.requests = s.requests[1:]
	}

	// 检查是否超过限制
	if len(s.requests)+n > s.maxRequest {
		return false
	}

	// 记录当前请求
	for i := 0; i < n; i++ {
		s.requests = append(s.requests, now)
	}
	return true
}

func (s *SlidingWindowLimiter) Wait(ctx context.Context) error {
	for {
		if s.Allow() {
			return nil
		}

		select {
		case <-time.After(100 * time.Millisecond):
			continue
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

func (s *SlidingWindowLimiter) Burst() int {
	return s.maxRequest
}

func (s *SlidingWindowLimiter) Limit() float64 {
	return float64(s.maxRequest) / s.window.Seconds()
}

// IPRateLimiter IP级别的限流器
type IPRateLimiter struct {
	limiters     map[string]*TokenBucketLimiter
	lastAccess   map[string]time.Time
	mutex        sync.RWMutex
	rps          float64
	burst        int
	cleanupTimer *time.Timer
	logger       logger.Logger
}

// NewIPRateLimiter 创建IP限流器
func NewIPRateLimiter(rps float64, burst int, log logger.Logger) *IPRateLimiter {
	limiter := &IPRateLimiter{
		limiters:   make(map[string]*TokenBucketLimiter),
		lastAccess: make(map[string]time.Time),
		rps:        rps,
		burst:      burst,
		logger:     log,
	}

	// 启动清理goroutine
	limiter.startCleanup()
	return limiter
}

// GetLimiter 获取IP对应的限流器
func (i *IPRateLimiter) GetLimiter(ip string) *TokenBucketLimiter {
	i.mutex.RLock()
	limiter, exists := i.limiters[ip]
	if exists {
		i.lastAccess[ip] = time.Now()
		i.mutex.RUnlock()
		return limiter
	}
	i.mutex.RUnlock()

	i.mutex.Lock()
	defer i.mutex.Unlock()

	// 双重检查
	if limiter, exists := i.limiters[ip]; exists {
		i.lastAccess[ip] = time.Now()
		return limiter
	}

	// 创建新的限流器
	limiter = NewTokenBucketLimiter(i.rps, i.burst, i.logger)
	i.limiters[ip] = limiter
	i.lastAccess[ip] = time.Now()

	return limiter
}

// Allow 检查IP是否允许请求
func (i *IPRateLimiter) Allow(ip string) bool {
	limiter := i.GetLimiter(ip)
	return limiter.Allow()
}

// startCleanup 启动清理过期限流器的goroutine
func (i *IPRateLimiter) startCleanup() {
	go func() {
		ticker := time.NewTicker(5 * time.Minute)
		defer ticker.Stop()

		for range ticker.C {
			i.CleanupExpired(10 * time.Minute)
		}
	}()
}

// CleanupExpired 清理过期的限流器
func (i *IPRateLimiter) CleanupExpired(maxIdle time.Duration) {
	i.mutex.Lock()
	defer i.mutex.Unlock()

	now := time.Now()
	cutoff := now.Add(-maxIdle)

	for ip, lastTime := range i.lastAccess {
		if lastTime.Before(cutoff) {
			delete(i.limiters, ip)
			delete(i.lastAccess, ip)
			i.logger.Debugf("Cleaned up expired limiter for IP: %s", ip)
		}
	}
}

// ConcurrencyLimiter 并发限制器
type ConcurrencyLimiter struct {
	semaphore      chan struct{}
	maxConcurrency int
	logger         logger.Logger
}

// NewConcurrencyLimiter 创建并发限制器
func NewConcurrencyLimiter(maxConcurrency int, log logger.Logger) *ConcurrencyLimiter {
	return &ConcurrencyLimiter{
		semaphore:      make(chan struct{}, maxConcurrency),
		maxConcurrency: maxConcurrency,
		logger:         log,
	}
}

// Acquire 获取并发许可
func (c *ConcurrencyLimiter) Acquire(ctx context.Context) error {
	select {
	case c.semaphore <- struct{}{}:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

// Release 释放并发许可
func (c *ConcurrencyLimiter) Release() {
	select {
	case <-c.semaphore:
	default:
		c.logger.Warn("Release called on empty semaphore")
	}
}

// TryAcquire 尝试获取并发许可（非阻塞）
func (c *ConcurrencyLimiter) TryAcquire() bool {
	select {
	case c.semaphore <- struct{}{}:
		return true
	default:
		return false
	}
}

// Available 获取可用的并发许可数
func (c *ConcurrencyLimiter) Available() int {
	return c.maxConcurrency - len(c.semaphore)
}

// ============ Gin中间件 ============

// RateLimitMiddleware 限流中间件
func RateLimitMiddleware(limiter RateLimiter) gin.HandlerFunc {
	return func(c *gin.Context) {
		if !limiter.Allow() {
			errors.SendError(c, errors.ErrTooManyRequests)
			return
		}
		c.Next()
	}
}

// IPRateLimitMiddleware IP限流中间件
func IPRateLimitMiddleware(ipLimiter *IPRateLimiter) gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()
		if !ipLimiter.Allow(ip) {
			errors.SendError(c, errors.ErrTooManyRequests.WithDetails(map[string]interface{}{
				"client_ip": ip,
				"message":   "Rate limit exceeded for this IP",
			}))
			return
		}
		c.Next()
	}
}

// ConcurrencyLimitMiddleware 并发限制中间件
func ConcurrencyLimitMiddleware(concLimiter *ConcurrencyLimiter) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 尝试获取并发许可
		if !concLimiter.TryAcquire() {
			errors.SendError(c, errors.NewError(
				errors.ServiceUnavailable,
				errors.ErrorTypeSystem,
				"Too many concurrent requests",
			))
			return
		}

		// 确保在请求结束时释放许可
		defer concLimiter.Release()

		c.Next()
	}
}

// AdaptiveRateLimitMiddleware 自适应限流中间件
func AdaptiveRateLimitMiddleware(baseRPS float64, maxRPS float64, log logger.Logger) gin.HandlerFunc {
	limiter := NewTokenBucketLimiter(baseRPS, int(baseRPS), log)
	currentRPS := baseRPS

	var successCount, totalCount int64
	var mutex sync.Mutex

	// 定期调整限流速率
	go func() {
		ticker := time.NewTicker(10 * time.Second)
		defer ticker.Stop()

		for range ticker.C {
			mutex.Lock()
			if totalCount > 0 {
				successRate := float64(successCount) / float64(totalCount)

				// 根据成功率调整RPS
				if successRate > 0.95 && currentRPS < maxRPS {
					currentRPS = min(currentRPS*1.1, maxRPS)
					limiter = NewTokenBucketLimiter(currentRPS, int(currentRPS), log)
					log.Infof("Increased rate limit to %.2f RPS (success rate: %.2f)", currentRPS, successRate)
				} else if successRate < 0.8 && currentRPS > baseRPS {
					currentRPS = max(currentRPS*0.9, baseRPS)
					limiter = NewTokenBucketLimiter(currentRPS, int(currentRPS), log)
					log.Infof("Decreased rate limit to %.2f RPS (success rate: %.2f)", currentRPS, successRate)
				}
			}

			// 重置计数器
			successCount = 0
			totalCount = 0
			mutex.Unlock()
		}
	}()

	return func(c *gin.Context) {
		if !limiter.Allow() {
			mutex.Lock()
			totalCount++
			mutex.Unlock()

			errors.SendError(c, errors.ErrTooManyRequests)
			return
		}

		c.Next()

		// 记录请求结果
		mutex.Lock()
		totalCount++
		if c.Writer.Status() < 500 {
			successCount++
		}
		mutex.Unlock()
	}
}

// CircuitBreakerMiddleware 熔断器中间件
func CircuitBreakerMiddleware(failureThreshold int, timeout time.Duration, log logger.Logger) gin.HandlerFunc {
	var failureCount int
	var lastFailureTime time.Time
	var state = "closed" // closed, open, half-open
	var mutex sync.Mutex

	return func(c *gin.Context) {
		mutex.Lock()
		currentState := state
		currentLastFailure := lastFailureTime
		mutex.Unlock()

		// 检查熔断器状态
		switch currentState {
		case "open":
			if time.Since(currentLastFailure) > timeout {
				mutex.Lock()
				state = "half-open"
				mutex.Unlock()
				log.Info("Circuit breaker state changed to half-open")
			} else {
				errors.SendError(c, errors.NewError(
					errors.ServiceUnavailable,
					errors.ErrorTypeSystem,
					"Service temporarily unavailable (circuit breaker open)",
				))
				return
			}
		}

		c.Next()

		// 检查请求结果
		mutex.Lock()
		defer mutex.Unlock()

		if c.Writer.Status() >= 500 {
			failureCount++
			lastFailureTime = time.Now()

			if failureCount >= failureThreshold {
				state = "open"
				log.Warnf("Circuit breaker opened due to %d failures", failureCount)
			}
		} else {
			if state == "half-open" {
				state = "closed"
				failureCount = 0
				log.Info("Circuit breaker closed after successful request")
			}
		}
	}
}

// 辅助函数
func min(a, b float64) float64 {
	if a < b {
		return a
	}
	return b
}

func max(a, b float64) float64 {
	if a > b {
		return a
	}
	return b
}
