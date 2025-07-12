package config

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/spf13/viper"
	"go.uber.org/fx"
)

// Config 配置结构体
type Config struct {
	App      AppConfig      `mapstructure:"app"`
	HTTP     HTTPConfig     `mapstructure:"http"`
	Database DatabaseConfig `mapstructure:"database"`
	Log      LogConfig      `mapstructure:"log"`
	Cache    CacheConfig    `mapstructure:"cache"`
}

// AppConfig 应用配置
type AppConfig struct {
	Name    string
	Version string
	Mode    string // dev, test, prod
}

// HTTPConfig HTTP服务配置
type HTTPConfig struct {
	Host           string
	Port           int
	Engine         string // 引擎类型："gin" 或 "fiber"
	ReadTimeout    time.Duration
	WriteTimeout   time.Duration
	MaxHeaderBytes int
	MaxBodySize    int // 请求体大小限制
}

// DatabaseConfig 数据库配置
type DatabaseConfig struct {
	Driver          string
	DSN             string
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
	LogLevel        string
}

// LogConfig 日志配置
type LogConfig struct {
	Level      string // debug, info, warn, error
	Format     string // json, text
	OutputPath string
}

// CacheConfig Cache缓存配置
type CacheConfig struct {
	Type     string // 缓存类型：memory、redis 或 file
	Host     string
	Port     int
	Password string
	DB       int
	Prefix   string // 键前缀
	FilePath string // 文件缓存路径，仅当 Type 为 file 时使用
}

// NewConfig 创建配置
// 修改 pkg/config/config.go 中的 NewConfig 函数
func NewConfig() (*Config, error) {
	// 创建默认配置
	config := &Config{}

	// 设置默认值（无论配置文件是否存在）
	setDefaultConfig(config)

	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		configPath = "config"
	}

	// 确保配置目录存在
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		fmt.Printf("配置目录 %s 不存在，使用默认配置\n", configPath)
		return config, nil
	}

	v := viper.New()
	v.AddConfigPath(configPath)
	configName := os.Getenv("CONFIG_NAME")
	v.SetConfigName(configName)
	if configName == "" {
		v.SetConfigName("config")
	}
	v.SetConfigType("yaml")

	// 读取环境变量
	v.AutomaticEnv()

	// 读取配置文件
	if err := v.ReadInConfig(); err != nil {
		fmt.Printf("读取配置文件错误: %v，使用默认配置\n", err)
		return config, nil
	}

	// 将文件配置合并到默认配置
	if err := v.Unmarshal(config); err != nil {
		return nil, fmt.Errorf("解析配置错误: %w", err)
	}

	return config, nil
}

// 设置默认值
func setDefaultConfig(config *Config) {
	if config.App.Name == "" {
		config.App.Name = "goflow-app"
	}
	if config.App.Version == "" {
		config.App.Version = "0.1.0"
	}
	if config.App.Mode == "" {
		config.App.Mode = "dev"
	}

	// HTTP默认配置
	if config.HTTP.Host == "" {
		config.HTTP.Host = "0.0.0.0"
	}
	if config.HTTP.Port == 0 {
		config.HTTP.Port = 8080
	}
	if config.HTTP.Engine == "" {
		config.HTTP.Engine = "fiber" // 默认使用Gin引擎
	}
	if config.HTTP.ReadTimeout == 0 {
		config.HTTP.ReadTimeout = 10 * time.Second
	}
	if config.HTTP.WriteTimeout == 0 {
		config.HTTP.WriteTimeout = 10 * time.Second
	}
	if config.HTTP.MaxHeaderBytes == 0 {
		config.HTTP.MaxHeaderBytes = 1 << 20 // 1MB
	}
	if config.HTTP.MaxBodySize == 0 {
		config.HTTP.MaxBodySize = 4 << 20 // 4MB
	}

	// 日志默认配置
	if config.Log.Level == "" {
		config.Log.Level = "info"
	}
	if config.Log.Format == "" {
		config.Log.Format = "json"
	}
	if config.Log.OutputPath == "" {
		config.Log.OutputPath = "stdout"
	}

	// 数据库默认配置
	if config.Database.MaxOpenConns == 0 {
		config.Database.MaxOpenConns = 100
	}
	if config.Database.MaxIdleConns == 0 {
		config.Database.MaxIdleConns = 10
	}
	if config.Database.ConnMaxLifetime == 0 {
		config.Database.ConnMaxLifetime = time.Hour
	}
	if config.Database.LogLevel == "" {
		config.Database.LogLevel = "error"
	}

	// Cache默认配置
	if config.Cache.Type == "" {
		config.Cache.Type = "memory" // 默认使用内存缓存
	}
	if config.Cache.Host == "" {
		config.Cache.Host = "localhost"
	}
	if config.Cache.Port == 0 {
		config.Cache.Port = 6379
	}
	if config.Cache.DB < 0 {
		config.Cache.DB = 0
	}
	if config.Cache.Prefix == "" {
		config.Cache.Prefix = "goflow:"
	}
	if config.Cache.FilePath == "" {
		// 默认使用当前工作目录下的 storage/cache 目录
		cwd, _ := os.Getwd()
		config.Cache.FilePath = filepath.Join(cwd, "storage", "cache")
	}
}

// Module 提供配置模块
var Module = fx.Provide(NewConfig)
