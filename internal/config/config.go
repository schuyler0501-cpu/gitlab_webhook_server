package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

// Config 应用配置
type Config struct {
	Port          string
	Environment   string
	LogLevel      string
	WebhookSecret string
	Database      DatabaseConfig
	WorkerPool    WorkerPoolConfig
	RateLimit     RateLimitConfig
	GitLab        GitLabConfig
}

// WorkerPoolConfig 工作池配置
type WorkerPoolConfig struct {
	Workers  int
	QueueSize int
}

// RateLimitConfig 限流配置
type RateLimitConfig struct {
	Limit  int
	Window string // 时间窗口，如 "1m", "1h"
}

// GitLabConfig GitLab API 配置
type GitLabConfig struct {
	BaseURL string
	Token   string
}

// DatabaseConfig 数据库配置
type DatabaseConfig struct {
	Type     string // 数据库类型: mysql, postgresql
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string // PostgreSQL 使用
	TimeZone string
	Charset  string // MySQL 使用
}

// Load 加载配置
func Load() (*Config, error) {
	// 尝试加载 .env 文件（如果存在）
	_ = godotenv.Load()

	cfg := &Config{
		Port:          getEnv("PORT", "3000"),
		Environment:   getEnv("NODE_ENV", "development"),
		LogLevel:      getEnv("LOG_LEVEL", "info"),
		WebhookSecret: getEnv("GITLAB_WEBHOOK_SECRET", ""),
		Database: DatabaseConfig{
			Type:     getEnv("DB_TYPE", "mysql"), // 默认 MySQL
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "3306"), // MySQL 默认端口
			User:     getEnv("DB_USER", "root"),
			Password: getEnv("DB_PASSWORD", ""),
			DBName:   getEnv("DB_NAME", "gitlab_webhook"),
			SSLMode:  getEnv("DB_SSLMODE", "disable"),
			TimeZone: getEnv("DB_TIMEZONE", "Asia/Shanghai"),
			Charset:  getEnv("DB_CHARSET", "utf8mb4"),
		},
		WorkerPool: WorkerPoolConfig{
			Workers:   getEnvInt("WORKER_POOL_WORKERS", 10),
			QueueSize: getEnvInt("WORKER_POOL_QUEUE_SIZE", 100),
		},
		RateLimit: RateLimitConfig{
			Limit:  getEnvInt("RATE_LIMIT_LIMIT", 100),
			Window: getEnv("RATE_LIMIT_WINDOW", "1m"),
		},
		GitLab: GitLabConfig{
			BaseURL: getEnv("GITLAB_BASE_URL", ""),
			Token:   getEnv("GITLAB_TOKEN", ""),
		},
	}

	return cfg, nil
}

// GetDSN 获取数据库连接字符串
func (c *Config) GetDSN() string {
	switch c.Database.Type {
	case "postgresql", "postgres":
		return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s TimeZone=%s",
			c.Database.Host,
			c.Database.Port,
			c.Database.User,
			c.Database.Password,
			c.Database.DBName,
			c.Database.SSLMode,
			c.Database.TimeZone,
		)
	case "mysql":
		fallthrough
	default:
		// MySQL DSN 格式: user:password@tcp(host:port)/dbname?charset=utf8mb4&parseTime=True&loc=Local
		dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=True&loc=%s",
			c.Database.User,
			c.Database.Password,
			c.Database.Host,
			c.Database.Port,
			c.Database.DBName,
			c.Database.Charset,
			c.Database.TimeZone,
		)
		return dsn
	}
}

// GetDatabaseType 获取数据库类型
func (c *Config) GetDatabaseType() string {
	return c.Database.Type
}

// getEnv 获取环境变量，如果不存在则返回默认值
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// getEnvInt 获取整数环境变量
func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		var result int
		if _, err := fmt.Sscanf(value, "%d", &result); err == nil {
			return result
		}
	}
	return defaultValue
}

