package config

import (
	"os"
	"sync"
)

// Config 程序配置结构
type Config struct {
	ServerPort string
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	JWTSecret  string
}

var (
	config     *Config
	configOnce sync.Once
)

// InitConfig 初始化配置
func InitConfig() {
	configOnce.Do(func() {
		config = &Config{
			ServerPort: getEnv("SERVER_PORT", "8080"),
			DBHost:     getEnv("DB_HOST", "localhost"),
			DBPort:     getEnv("DB_PORT", "5432"),
			DBUser:     getEnv("DB_USER", "admin"),
			DBPassword: getEnv("DB_PASSWORD", "123456"),
			DBName:     getEnv("DB_NAME", "cmdb"),
			JWTSecret:  getEnv("JWT_SECRET", "kk-ops-secure-jwt-token-2025"),
		}
	})
}

// GetConfig 获取配置实例
func GetConfig() *Config {
	if config == nil {
		InitConfig()
	}
	return config
}

// getEnv 获取环境变量，如果不存在则返回默认值
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
