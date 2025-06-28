package captcha

import (
	"context"
	"time"

	"github.com/mojocn/base64Captcha"
	"github.com/zhoudm1743/go-flow/core/cache"
)

// CaptchaType 验证码类型
type CaptchaType int

const (
	TypeDigit  CaptchaType = iota // 数字验证码
	TypeMath                      // 算数验证码
	TypeString                    // 字符串验证码
	TypeAudio                     // 音频验证码
)

// CaptchaConfig 验证码配置
type CaptchaConfig struct {
	Type       CaptchaType // 验证码类型
	Width      int
	Height     int
	Length     int
	NoiseCount int
	ShowLine   bool
	UseRedis   bool
	Expiration time.Duration
	KeyPrefix  string
	// 算数验证码专用配置
	MathFont string // 数学验证码字体
}

// DefaultConfig 默认配置
var DefaultConfig = &CaptchaConfig{
	Type:       TypeDigit,
	Width:      240,
	Height:     80,
	Length:     4,
	NoiseCount: 5,
	ShowLine:   true,
	UseRedis:   true,
	Expiration: 5 * time.Minute,
	KeyPrefix:  "captcha",
}

// DefaultMathConfig 默认算数验证码配置
var DefaultMathConfig = &CaptchaConfig{
	Type:       TypeMath,
	Width:      240,
	Height:     80,
	Length:     2, // 算数验证码长度指的是操作数个数
	NoiseCount: 5,
	ShowLine:   true,
	UseRedis:   true,
	Expiration: 5 * time.Minute,
	KeyPrefix:  "captcha",
	MathFont:   "", // 使用默认字体
}

// Service 验证码服务
type Service struct {
	driver base64Captcha.Driver
	store  base64Captcha.Store
	config *CaptchaConfig
}

// NewService 创建验证码服务
func NewService(cacheClient cache.Cache, config *CaptchaConfig) *Service {
	if config == nil {
		config = DefaultConfig
	}

	// 根据类型创建不同的驱动
	var driver base64Captcha.Driver
	switch config.Type {
	case TypeMath:
		// ShowLine转换为int：true=1, false=0
		showLineInt := 0
		if config.ShowLine {
			showLineInt = 1
		}
		driver = base64Captcha.NewDriverMath(
			config.Height,
			config.Width,
			config.NoiseCount,
			showLineInt,
			nil,        // 使用默认背景色
			nil,        // 使用默认字体
			[]string{}, // 字体文件路径数组，空数组使用默认字体
		)
	case TypeDigit:
		driver = base64Captcha.NewDriverDigit(
			config.Height,
			config.Width,
			config.Length,
			0.7,
			config.NoiseCount,
		)
	case TypeString:
		// ShowLine转换为int：true=1, false=0
		showLineInt := 0
		if config.ShowLine {
			showLineInt = 1
		}
		driver = base64Captcha.NewDriverString(
			config.Height,
			config.Width,
			config.NoiseCount,
			showLineInt,
			config.Length,
			"abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789",
			nil,
			nil,
			nil,
		)
	default:
		// 默认使用数字验证码
		driver = base64Captcha.NewDriverDigit(
			config.Height,
			config.Width,
			config.Length,
			0.7,
			config.NoiseCount,
		)
	}

	// 选择存储方式
	var captchaStore base64Captcha.Store
	if config.UseRedis && cacheClient != nil {
		captchaStore = NewRedisStore(cacheClient, config.Expiration, config.KeyPrefix)
	} else {
		captchaStore = GetStore()
	}

	return &Service{
		driver: driver,
		store:  captchaStore,
		config: config,
	}
}

// NewMathService 创建算数验证码服务的便捷方法
func NewMathService(cacheClient cache.Cache) *Service {
	return NewService(cacheClient, DefaultMathConfig)
}

// Generate 生成验证码
func (s *Service) Generate() (id, b64s, answer string, err error) {
	captcha := base64Captcha.NewCaptcha(s.driver, s.store)
	id, b64s, answer, err = captcha.Generate()
	return id, b64s, answer, err
}

// Verify 验证验证码
func (s *Service) Verify(id, answer string, clear bool) bool {
	captcha := base64Captcha.NewCaptcha(s.driver, s.store)
	return captcha.Verify(id, answer, clear)
}

// GenerateWithContext 带context的生成验证码
func (s *Service) GenerateWithContext(ctx context.Context) (id, b64s, answer string, err error) {
	// 检查context是否已取消
	select {
	case <-ctx.Done():
		return "", "", "", ctx.Err()
	default:
	}

	return s.Generate()
}

// VerifyWithContext 带context的验证验证码
func (s *Service) VerifyWithContext(ctx context.Context, id, answer string, clear bool) bool {
	// 检查context是否已取消
	select {
	case <-ctx.Done():
		return false
	default:
	}

	return s.Verify(id, answer, clear)
}
