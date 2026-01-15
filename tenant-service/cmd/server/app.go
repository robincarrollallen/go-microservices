package main

import (
	"fmt"
	
	"tenant-service/internal/config"
	"tenant-service/internal/handler"
	"tenant-service/internal/logger"
	"tenant-service/internal/middleware"
	"tenant-service/internal/repo"
	"tenant-service/internal/service"

	"github.com/gin-gonic/gin"
)

func run() error {
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("load config failed: %w", err)
	}

	logger.Init(cfg.Log.Level)
	defer logger.L().Sync()

	r := gin.New()
	r.Use(
		middleware.TraceID(),
		middleware.Logger(),
		gin.Recovery(),
	)

	// 依赖初始化
	tenantRepo := repo.NewTenantRepo()
	tenantService := service.NewTenantService(tenantRepo)
	tenantHandler := handler.NewTenantHandler(tenantService)
	tenantHandler.Register(r)

	addr := fmt.Sprintf(":%d", cfg.HTTP.Port)
	return r.Run(addr)
}
