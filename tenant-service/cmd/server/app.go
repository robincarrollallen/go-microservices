package main

import (
	"fmt"
	"tenant-service/internal/handler"
	"tenant-service/internal/middleware"
	"tenant-service/internal/repo"
	"tenant-service/internal/service"

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

	// 依赖初始化
	tenantRepo := repo.NewTenantRepo(db)
	tenantService := service.NewTenantService(tenantRepo)
	tenantHandler := handler.NewTenantHandler(tenantService)

	r := gin.New()
	r.Use(
		pkgMiddleware.TraceID(),
		pkgMiddleware.Logger(),
		pkgMiddleware.ErrorHandler(),
		middleware.ServiceErrorHandler(),
		pkgMiddleware.Recovery(),
	)
	tenantHandler.Register(r) // 注册所有路由

	addr := fmt.Sprintf(":%d", cfg.HTTP.Port)
	return r.Run(addr)
}
