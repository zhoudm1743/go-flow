package captcha

import (
	"fmt"
	"log"
	"time"

	"github.com/zhoudm1743/go-flow/core/cache"
)

// MathCaptchaExample 算数验证码使用示例
func MathCaptchaExample(cacheClient cache.Cache) {
	// 方法1: 使用默认的算数验证码服务
	mathService := NewMathService(cacheClient)

	// 生成算数验证码
	id, b64Image, answer, err := mathService.Generate()
	if err != nil {
		log.Printf("生成算数验证码失败: %v", err)
		return
	}

	fmt.Printf("算数验证码生成成功:\n")
	fmt.Printf("ID: %s\n", id)
	fmt.Printf("答案: %s\n", answer)                 // 这个在实际应用中不应该显示给用户
	fmt.Printf("Base64图片: %s...\n", b64Image[:50]) // 只显示前50个字符

	// 验证算数验证码
	userAnswer := answer // 模拟用户输入正确答案
	isValid := mathService.Verify(id, userAnswer, true)
	fmt.Printf("验证结果: %v\n", isValid)

	// 方法2: 使用自定义配置
	customConfig := &CaptchaConfig{
		Type:       TypeMath,
		Width:      300,
		Height:     100,
		NoiseCount: 3,
		ShowLine:   true,
		UseRedis:   true,
		Expiration: 10 * time.Minute,
		KeyPrefix:  "math_captcha",
	}

	customMathService := NewService(cacheClient, customConfig)

	// 生成自定义配置的算数验证码
	id2, b64Image2, answer2, err := customMathService.Generate()
	if err != nil {
		log.Printf("生成自定义算数验证码失败: %v", err)
		return
	}

	fmt.Printf("\n自定义算数验证码生成成功:\n")
	fmt.Printf("ID: %s\n", id2)
	fmt.Printf("答案: %s\n", answer2)
	fmt.Printf("Base64图片: %s...\n", b64Image2[:50])
}

// GenerateMathCaptcha 生成算数验证码的便捷函数
func GenerateMathCaptcha(cacheClient cache.Cache) (id, b64Image, answer string, err error) {
	service := NewMathService(cacheClient)
	return service.Generate()
}

// VerifyMathCaptcha 验证算数验证码的便捷函数
func VerifyMathCaptcha(cacheClient cache.Cache, id, userAnswer string) bool {
	service := NewMathService(cacheClient)
	return service.Verify(id, userAnswer, true)
}

// CreateMathCaptchaWithSize 创建指定尺寸的算数验证码
func CreateMathCaptchaWithSize(cacheClient cache.Cache, width, height int) (id, b64Image, answer string, err error) {
	config := &CaptchaConfig{
		Type:       TypeMath,
		Width:      width,
		Height:     height,
		NoiseCount: 5,
		ShowLine:   true,
		UseRedis:   true,
		Expiration: 5 * time.Minute,
		KeyPrefix:  "captcha",
	}

	service := NewService(cacheClient, config)
	return service.Generate()
}
