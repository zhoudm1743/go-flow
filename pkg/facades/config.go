package facades

import (
	"time"

	"github.com/zhoudm1743/go-flow/pkg/config"
)

// GetString 获取字符串配置
func (c *ConfigFacade) GetString(key string, defaultValue ...string) string {
	switch key {
	case "app.name":
		if GetConfig().App.Name != "" {
			return GetConfig().App.Name
		}
	case "app.version":
		if GetConfig().App.Version != "" {
			return GetConfig().App.Version
		}
	case "app.mode":
		if GetConfig().App.Mode != "" {
			return GetConfig().App.Mode
		}
	case "log.level":
		if GetConfig().Log.Level != "" {
			return GetConfig().Log.Level
		}
	case "log.format":
		if GetConfig().Log.Format != "" {
			return GetConfig().Log.Format
		}
	case "log.output_path":
		if GetConfig().Log.OutputPath != "" {
			return GetConfig().Log.OutputPath
		}
	case "database.driver":
		if GetConfig().Database.Driver != "" {
			return GetConfig().Database.Driver
		}
	case "database.dsn":
		if GetConfig().Database.DSN != "" {
			return GetConfig().Database.DSN
		}
	case "database.log_level":
		if GetConfig().Database.LogLevel != "" {
			return GetConfig().Database.LogLevel
		}
	case "http.host":
		if GetConfig().HTTP.Host != "" {
			return GetConfig().HTTP.Host
		}
	}

	if len(defaultValue) > 0 {
		return defaultValue[0]
	}
	return ""
}

// GetInt 获取整数配置
func (c *ConfigFacade) GetInt(key string, defaultValue ...int) int {
	switch key {
	case "http.port":
		if GetConfig().HTTP.Port != 0 {
			return GetConfig().HTTP.Port
		}
	case "http.max_header_bytes":
		if GetConfig().HTTP.MaxHeaderBytes != 0 {
			return GetConfig().HTTP.MaxHeaderBytes
		}
	case "database.max_open_conns":
		if GetConfig().Database.MaxOpenConns != 0 {
			return GetConfig().Database.MaxOpenConns
		}
	case "database.max_idle_conns":
		if GetConfig().Database.MaxIdleConns != 0 {
			return GetConfig().Database.MaxIdleConns
		}
	}

	if len(defaultValue) > 0 {
		return defaultValue[0]
	}
	return 0
}

// GetDuration 获取时间间隔配置
func (c *ConfigFacade) GetDuration(key string, defaultValue ...time.Duration) time.Duration {
	switch key {
	case "http.read_timeout":
		if GetConfig().HTTP.ReadTimeout != 0 {
			return GetConfig().HTTP.ReadTimeout
		}
	case "http.write_timeout":
		if GetConfig().HTTP.WriteTimeout != 0 {
			return GetConfig().HTTP.WriteTimeout
		}
	case "database.conn_max_lifetime":
		if GetConfig().Database.ConnMaxLifetime != 0 {
			return GetConfig().Database.ConnMaxLifetime
		}
	}

	if len(defaultValue) > 0 {
		return defaultValue[0]
	}
	return 0
}

// GetBool 获取布尔值配置
func (c *ConfigFacade) GetBool(key string, defaultValue ...bool) bool {
	// 当前配置中没有布尔值，这里预留接口
	if len(defaultValue) > 0 {
		return defaultValue[0]
	}
	return false
}

// Instance 获取原始配置实例
func (c *ConfigFacade) Instance() *config.Config {
	return GetConfig()
}
