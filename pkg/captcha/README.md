# 验证码模块使用说明

本模块提供了多种类型的验证码生成和验证功能，支持Redis存储和内存存储。

## 支持的验证码类型

- **数字验证码** (`TypeDigit`): 生成数字字符组合
- **算数验证码** (`TypeMath`): 生成数学计算题，如 "3+5=?" 
- **字符串验证码** (`TypeString`): 生成字母数字混合字符
- **音频验证码** (`TypeAudio`): 音频形式的验证码（可扩展）

## 快速开始

### 1. 生成算数验证码（推荐）

```go
package main

import (
    "log"
    "github.com/zhoudm1743/go-flow/pkg/captcha"
    "github.com/zhoudm1743/go-flow/core/cache"
)

func main() {
    // 假设您已经有了cache客户端
    var cacheClient cache.Cache
    
    // 方法1: 使用便捷函数
    id, b64Image, answer, err := captcha.GenerateMathCaptcha(cacheClient)
    if err != nil {
        log.Fatal(err)
    }
    
    log.Printf("验证码ID: %s", id)
    log.Printf("验证码图片(Base64): %s", b64Image)
    log.Printf("正确答案: %s", answer) // 仅用于测试，实际不应暴露
    
    // 验证用户输入
    userInput := "8" // 假设用户输入的答案
    isValid := captcha.VerifyMathCaptcha(cacheClient, id, userInput)
    log.Printf("验证结果: %v", isValid)
}
```

### 2. 使用自定义配置

```go
// 创建自定义算数验证码配置
config := &captcha.CaptchaConfig{
    Type:       captcha.TypeMath,
    Width:      300,
    Height:     100,
    NoiseCount: 3,
    ShowLine:   true,
    UseRedis:   true,
    Expiration: 10 * time.Minute,
    KeyPrefix:  "math_captcha",
}

// 创建服务
service := captcha.NewService(cacheClient, config)

// 生成验证码
id, b64Image, answer, err := service.Generate()
```

### 3. 使用其他类型的验证码

```go
// 数字验证码
digitConfig := &captcha.CaptchaConfig{
    Type:       captcha.TypeDigit,
    Width:      240,
    Height:     80,
    Length:     4,
    NoiseCount: 5,
    UseRedis:   true,
}

digitService := captcha.NewService(cacheClient, digitConfig)
id, b64Image, answer, err := digitService.Generate()
```

## API 参考

### 主要类型

```go
// 验证码类型
type CaptchaType int

const (
    TypeDigit  CaptchaType = iota // 数字验证码
    TypeMath                      // 算数验证码  
    TypeString                    // 字符串验证码
    TypeAudio                     // 音频验证码
)

// 验证码配置
type CaptchaConfig struct {
    Type       CaptchaType   // 验证码类型
    Width      int           // 图片宽度
    Height     int           // 图片高度
    Length     int           // 验证码长度
    NoiseCount int           // 噪点数量
    ShowLine   bool          // 是否显示干扰线
    UseRedis   bool          // 是否使用Redis存储
    Expiration time.Duration // 过期时间
    KeyPrefix  string        // Redis键前缀
    MathFont   string        // 数学验证码字体
}
```

### 主要方法

```go
// 创建服务
func NewService(cacheClient cache.Cache, config *CaptchaConfig) *Service
func NewMathService(cacheClient cache.Cache) *Service

// 生成验证码
func (s *Service) Generate() (id, b64s, answer string, err error)
func (s *Service) GenerateWithContext(ctx context.Context) (id, b64s, answer string, err error)

// 验证验证码
func (s *Service) Verify(id, answer string, clear bool) bool
func (s *Service) VerifyWithContext(ctx context.Context, id, answer string, clear bool) bool

// 便捷函数
func GenerateMathCaptcha(cacheClient cache.Cache) (id, b64Image, answer string, err error)
func VerifyMathCaptcha(cacheClient cache.Cache, id, userAnswer string) bool
func CreateMathCaptchaWithSize(cacheClient cache.Cache, width, height int) (id, b64Image, answer string, err error)
```

## 存储选项

### Redis存储（推荐）
- 支持分布式部署
- 自动过期清理
- 高性能

### 内存存储
- 适用于单机部署
- 简单快速
- 重启后数据丢失

## 注意事项

1. **安全性**: 永远不要将验证码答案暴露给前端
2. **过期时间**: 建议设置合理的过期时间（5-10分钟）
3. **清理策略**: 验证成功后建议立即清除验证码
4. **Redis配置**: 确保Redis连接正常，否则会自动降级到内存存储

## 前端集成示例

```html
<!DOCTYPE html>
<html>
<head>
    <title>算数验证码示例</title>
</head>
<body>
    <div>
        <img id="captcha" src="" alt="验证码" onclick="refreshCaptcha()">
        <input type="text" id="answer" placeholder="请输入计算结果">
        <button onclick="verifyCaptcha()">验证</button>
    </div>

    <script>
        let captchaId = '';
        
        // 生成新验证码
        function refreshCaptcha() {
            fetch('/api/captcha/generate', {method: 'POST'})
                .then(res => res.json())
                .then(data => {
                    captchaId = data.id;
                    document.getElementById('captcha').src = 'data:image/png;base64,' + data.image;
                });
        }
        
        // 验证验证码
        function verifyCaptcha() {
            const answer = document.getElementById('answer').value;
            fetch('/api/captcha/verify', {
                method: 'POST',
                headers: {'Content-Type': 'application/json'},
                body: JSON.stringify({id: captchaId, answer: answer})
            })
            .then(res => res.json())
            .then(data => {
                alert(data.valid ? '验证成功' : '验证失败');
                if (!data.valid) refreshCaptcha();
            });
        }
        
        // 页面加载时生成验证码
        refreshCaptcha();
    </script>
</body>
</html>
```

## 常见问题

**Q: 为什么选择算数验证码？**
A: 算数验证码具有以下优势：
- 对用户更友好，容易理解
- 具有一定的防机器识别能力
- 答案简短，用户输入方便
- 支持无障碍访问

**Q: 如何自定义验证码样式？**
A: 可以通过修改CaptchaConfig中的Width、Height、NoiseCount、ShowLine等参数来调整样式。

**Q: Redis连接失败怎么办？**
A: 系统会自动降级到内存存储，确保功能正常运行。 