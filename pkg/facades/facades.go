package facades

import (
	"sync"

	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"github.com/zhoudm1743/go-flow/pkg/cache"
	"github.com/zhoudm1743/go-flow/pkg/config"
	"github.com/zhoudm1743/go-flow/pkg/log"
	"go.uber.org/fx"
	"gorm.io/gorm"
)

// 全局命名空间
var (
	DB       *DBFacade
	Log      *LogFacade
	Config   *ConfigFacade
	Cache    *CacheFacade
	Validate *ValidateFacade
)

// 全局服务实例
var (
	dbInstance         *gorm.DB
	loggerInstance     log.Logger
	configInstance     *config.Config
	cacheInstance      cache.Cache
	validatorInstance  *validator.Validate
	translatorInstance ut.Translator

	// 确保线程安全
	mu sync.RWMutex
)

// DBFacade 数据库门面
type DBFacade struct{}

// LogFacade 日志门面
type LogFacade struct{}

// ConfigFacade 配置门面
type ConfigFacade struct{}

// CacheFacade 缓存门面
type CacheFacade struct{}

// ValidateFacade 验证器门面
type ValidateFacade struct{}

// FacadesParams 用于初始化门面的参数
type FacadesParams struct {
	fx.In
	DB     *gorm.DB
	Logger log.Logger
	Config *config.Config
	Cache  cache.Cache `optional:"true"`
}

func init() {
	// 初始化门面命名空间
	DB = &DBFacade{}
	Log = &LogFacade{}
	Config = &ConfigFacade{}
	Cache = &CacheFacade{}
	Validate = &ValidateFacade{}
}

// Initialize 初始化所有门面
func Initialize(p FacadesParams) {
	mu.Lock()
	defer mu.Unlock()

	dbInstance = p.DB
	loggerInstance = p.Logger
	configInstance = p.Config

	// 缓存实例可能为空，因为标记为可选
	if p.Cache != nil {
		cacheInstance = p.Cache
	}
}

// GetGormDB 获取数据库实例
func GetGormDB() *gorm.DB {
	mu.RLock()
	defer mu.RUnlock()

	if dbInstance == nil {
		panic("数据库实例未初始化，请先调用facades.Initialize")
	}

	return dbInstance
}

// GetLogger 获取日志实例
func GetLogger() log.Logger {
	mu.RLock()
	defer mu.RUnlock()

	if loggerInstance == nil {
		panic("日志实例未初始化，请先调用facades.Initialize")
	}

	return loggerInstance
}

// GetConfig 获取配置实例
func GetConfig() *config.Config {
	mu.RLock()
	defer mu.RUnlock()

	if configInstance == nil {
		panic("配置实例未初始化，请先调用facades.Initialize")
	}

	return configInstance
}

// GetCache 获取缓存实例
func GetCache() cache.Cache {
	mu.RLock()
	defer mu.RUnlock()

	if cacheInstance == nil {
		panic("缓存实例未初始化，请先调用facades.Initialize或确保缓存模块已注册")
	}

	return cacheInstance
}

// GetValidator 获取验证器实例
func GetValidator() *validator.Validate {
	mu.RLock()
	defer mu.RUnlock()
	return validatorInstance
}

// GetTranslator 获取翻译器实例
func GetTranslator() ut.Translator {
	mu.RLock()
	defer mu.RUnlock()
	return translatorInstance
}

// SetValidator 设置验证器实例
func SetValidator(validator *validator.Validate) {
	mu.Lock()
	defer mu.Unlock()
	validatorInstance = validator
}

// SetTranslator 设置翻译器实例
func SetTranslator(translator ut.Translator) {
	mu.Lock()
	defer mu.Unlock()
	translatorInstance = translator
}

// Module 提供门面模块
var Module = fx.Invoke(Initialize)
