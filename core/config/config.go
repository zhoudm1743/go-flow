package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
	"go.uber.org/fx"
)

// Config 配置结构体
type Config struct {
	App      AppConfig      `mapstructure:"app"`
	Database DatabaseConfig `mapstructure:"database"`
	Redis    RedisConfig    `mapstructure:"redis"`
	Log      LogConfig      `mapstructure:"log"`
	Jwt      JwtConfig      `mapstructure:"jwt"`
}

// AppConfig 应用配置
type AppConfig struct {
	Name    string     `mapstructure:"name"`
	Version string     `mapstructure:"version"`
	Port    int        `mapstructure:"port"`
	Env     string     `mapstructure:"env"`
	HTTP    HTTPConfig `mapstructure:"http"`
}

// HTTPConfig HTTP 服务配置
type HTTPConfig struct {
	ReadTimeout  int  `mapstructure:"read_timeout"`  // 读取超时时间（秒）
	WriteTimeout int  `mapstructure:"write_timeout"` // 写入超时时间（秒）
	IdleTimeout  int  `mapstructure:"idle_timeout"`  // 空闲超时时间（秒）
	EnableCORS   bool `mapstructure:"enable_cors"`   // 是否启用 CORS
}

// DatabaseConfig 数据库配置
type DatabaseConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	Database string `mapstructure:"database"`
}

// RedisConfig Redis配置
type RedisConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
}

// LogConfig 日志配置
type LogConfig struct {
	Level      string           `mapstructure:"level"`
	Format     string           `mapstructure:"format"`
	Output     string           `mapstructure:"output"`
	Lumberjack LumberjackConfig `mapstructure:"lumberjack"`
}

// LumberjackConfig lumberjack 日志轮转配置
type LumberjackConfig struct {
	Filename   string `mapstructure:"filename"`   // 日志文件路径
	MaxSize    int    `mapstructure:"maxsize"`    // 单个日志文件最大大小（MB）
	MaxAge     int    `mapstructure:"maxage"`     // 日志文件保留天数
	MaxBackups int    `mapstructure:"maxbackups"` // 最多保留多少个备份文件
	Compress   bool   `mapstructure:"compress"`   // 是否压缩旧日志文件
}

// JwtConfig JWT配置
type JwtConfig struct {
	Secret         string `mapstructure:"secret"`          // JWT 签名密钥
	ExpiresIn      string `mapstructure:"expires_in"`      // 访问token过期时间（小时）
	RefreshExpires string `mapstructure:"refresh_expires"` // 刷新token过期时间（小时）
	Issuer         string `mapstructure:"issuer"`          // 签发者
}

// Module fx模块
var Module = fx.Options(
	fx.Provide(NewConfig),
)

// NewConfig 创建新的配置实例
func NewConfig() (*Config, error) {
	v := viper.New()

	// 设置配置文件名和类型
	v.SetConfigName("config")
	v.SetConfigType("yaml")

	// 添加配置文件搜索路径
	v.AddConfigPath(".")
	v.AddConfigPath("./config")
	v.AddConfigPath("../config")

	// 读取环境变量
	v.AutomaticEnv()

	// 读取配置文件
	if err := v.ReadInConfig(); err != nil {
		// 如果配置文件不存在，创建默认配置文件
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			fmt.Println("配置文件未找到，创建默认配置文件...")
			if err := createDefaultConfig(); err != nil {
				return nil, fmt.Errorf("创建默认配置文件失败: %w", err)
			}
			// 重新读取配置文件
			if err := v.ReadInConfig(); err != nil {
				return nil, fmt.Errorf("读取配置文件失败: %w", err)
			}
		} else {
			return nil, fmt.Errorf("读取配置文件失败: %w", err)
		}
	}

	// 解析配置到结构体
	var config Config
	if err := v.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("解析配置文件失败: %w", err)
	}

	fmt.Printf("成功加载配置文件: %s\n", v.ConfigFileUsed())
	return &config, nil
}

// createDefaultConfig 创建默认配置文件
func createDefaultConfig() error {
	defaultConfig := `app:
  name: "go-flow"
  version: "1.0.0"
  port: 8080
  env: "development"
  http:
    read_timeout: 30
    write_timeout: 30
    idle_timeout: 120
    enable_cors: true

database:
  host: "localhost"
  port: 3306
  username: "root"
  password: ""
  database: "go_flow"

redis:
  host: "localhost"
  port: 6379
  password: ""
  db: 0

log:
  level: "info"
  format: "json"
  output: "stdout"
  lumberjack:
    filename: "logs/app.log"
    maxsize: 100
    maxage: 7
    maxbackups: 3
    compress: true

jwt:
  secret: "your-secret-key-change-in-production"
  expires_in: 24
  refresh_expires: 168
  issuer: "go-flow"
`

	// 确保配置目录存在
	configDir := "config"
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return fmt.Errorf("创建配置目录失败: %w", err)
	}

	// 写入默认配置文件
	configPath := filepath.Join(configDir, "config.yaml")
	if err := os.WriteFile(configPath, []byte(defaultConfig), 0644); err != nil {
		return fmt.Errorf("写入配置文件失败: %w", err)
	}

	return nil
}
