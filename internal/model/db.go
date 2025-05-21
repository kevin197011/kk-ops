package model

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/kevin197011/kk-ops/internal/config"
)

var (
	dbPool *pgxpool.Pool
)

// InitDB 初始化数据库连接
func InitDB() error {
	cfg := config.GetConfig()

	// 确保配置被正确加载
	log.Printf("Database configuration: Host=%s, Port=%s, User=%s, DBName=%s",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBName)

	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName)

	log.Printf("Connecting to database: %s:%s/%s (user: %s)", cfg.DBHost, cfg.DBPort, cfg.DBName, cfg.DBUser)

	var err error

	// 尝试最多3次连接
	for i := 0; i < 3; i++ {
		dbPool, err = pgxpool.New(context.Background(), connStr)
		if err == nil {
			break
		}
		log.Printf("Failed to create database pool (attempt %d/3): %v", i+1, err)
		time.Sleep(1 * time.Second)
	}

	if err != nil {
		log.Printf("All attempts to connect to database failed: %v", err)
		// 不要在这里退出程序，而是让调用者决定如何处理错误
		return err
	}

	// 测试连接
	if err := dbPool.Ping(context.Background()); err != nil {
		log.Printf("Failed to ping database: %v", err)
		return err
	}

	log.Println("Connected to database successfully")

	// 验证数据库表是否存在
	if err := verifyTables(); err != nil {
		log.Printf("Database tables verification failed: %v", err)
		// 这个错误可以不那么致命
	}

	return nil
}

// verifyTables 验证必要的表是否存在
func verifyTables() error {
	ctx := context.Background()

	tables := []string{"users", "servers", "pipelines", "pipeline_runs"}

	for _, table := range tables {
		var exists bool
		query := `
			SELECT EXISTS (
				SELECT FROM information_schema.tables
				WHERE table_schema = 'public'
				AND table_name = $1
			)
		`

		err := dbPool.QueryRow(ctx, query, table).Scan(&exists)
		if err != nil {
			return fmt.Errorf("failed to check if table %s exists: %w", table, err)
		}

		if !exists {
			log.Printf("WARNING: Table '%s' does not exist", table)
		} else {
			log.Printf("Table '%s' exists", table)
		}
	}

	return nil
}

// GetDB 获取数据库连接池
func GetDB() *pgxpool.Pool {
	// 如果连接池为空，尝试重新初始化
	if dbPool == nil {
		log.Println("Database connection pool is nil, trying to reinitialize")
		if err := InitDB(); err != nil {
			log.Printf("Failed to reinitialize database: %v", err)
			return nil
		}
	}
	return dbPool
}

// CloseDB 关闭数据库连接
func CloseDB() {
	if dbPool != nil {
		dbPool.Close()
		log.Println("Database connection closed")
	}
}
