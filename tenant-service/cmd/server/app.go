package main

import (
	"fmt"
	"tenant-service/internal/bootstrap"
	"tenant-service/internal/middleware"

	"shared.local/pkg/config"
	"shared.local/pkg/database"
	"shared.local/pkg/logger"
	pkgMiddleware "shared.local/pkg/middleware"

	"github.com/gin-gonic/gin"
)

func run() error {
	cfg, err := config.LoadBase() // 一次性加载所有配置
	if err != nil {
		return fmt.Errorf("load config failed: %w", err)
	}

	logger.Init(cfg.Log.Level) // 初始化日志
	defer logger.L().Sync()

	dbCfg := database.PostgresConfig{
		Host:         cfg.DB.Host,
		Port:         cfg.DB.Port,
		User:         cfg.DB.User,
		Password:     cfg.DB.Password,
		DBName:       cfg.DB.DBName,
		SSLMode:      cfg.DB.SSLMode,
		MaxIdleConns: cfg.DB.MaxIdleConns,
		MaxOpenConns: cfg.DB.MaxOpenConns,
		MaxLifetime:  cfg.DB.MaxLifetime,
	}

	db, err := database.NewPostgres(dbCfg) // 将配置转换为 *gorm.DB 连接对象
	if err != nil {
		return fmt.Errorf("init database failed: %w", err)
	}

	r := gin.New()
	r.Use(
		pkgMiddleware.TraceID(),
		pkgMiddleware.Logger(),
		pkgMiddleware.ErrorHandler(),
		middleware.ServiceErrorHandler(),
		pkgMiddleware.Recovery(),
	)

	// 初始化依赖并注册路由
	container := bootstrap.SetupDependencies(db)
	bootstrap.RegisterRoutes(container, r)

	addr := fmt.Sprintf(":%d", cfg.HTTP.Port)
	return r.Run(addr)
}
