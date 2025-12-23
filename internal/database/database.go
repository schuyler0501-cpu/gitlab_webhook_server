package database

import (
	"fmt"
	"strings"

	"gitlab-webhook-server/internal/config"
	"gitlab-webhook-server/internal/model"

	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// DB 全局数据库实例
var DB *gorm.DB

// Init 初始化数据库连接
func Init(cfg *config.Config, zapLogger *zap.Logger) error {
	dsn := cfg.GetDSN()
	dbType := strings.ToLower(cfg.GetDatabaseType())

	// 配置 GORM 日志
	var gormLogger logger.Interface
	if cfg.Environment == "development" {
		gormLogger = logger.Default.LogMode(logger.Info)
	} else {
		gormLogger = logger.Default.LogMode(logger.Error)
	}

	// 根据数据库类型选择驱动
	var db *gorm.DB
	var err error

	switch dbType {
	case "postgresql", "postgres":
		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
			Logger: gormLogger,
		})
	case "mysql":
		fallthrough
	default:
		// 默认使用 MySQL
		db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
			Logger: gormLogger,
		})
	}

	if err != nil {
		return fmt.Errorf("连接数据库失败: %w", err)
	}

	// 测试连接
	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("获取数据库实例失败: %w", err)
	}

	if err := sqlDB.Ping(); err != nil {
		return fmt.Errorf("数据库连接测试失败: %w", err)
	}

	// 配置连接池
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)

	DB = db
	zapLogger.Info("数据库连接成功",
		zap.String("type", dbType),
		zap.String("host", cfg.Database.Host),
		zap.String("database", cfg.Database.DBName),
	)

	return nil
}

// Migrate 执行数据库迁移
func Migrate() error {
	if DB == nil {
		return fmt.Errorf("数据库未初始化")
	}

	// 自动迁移表结构
	err := DB.AutoMigrate(
		&model.Commit{},
		&model.CommitFile{},
		&model.CommitLanguage{},
		&model.MemberContribution{},
		&model.MemberLanguageStat{},
	)
	if err != nil {
		return fmt.Errorf("数据库迁移失败: %w", err)
	}

	return nil
}

// Close 关闭数据库连接
func Close() error {
	if DB == nil {
		return nil
	}

	sqlDB, err := DB.DB()
	if err != nil {
		return err
	}

	return sqlDB.Close()
}

