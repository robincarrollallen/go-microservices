package main

import (
	"fmt"

	"tenant-service/internal/handler"
	"tenant-service/internal/repo"
	"tenant-service/internal/service"

	"shared.local/pkg/config"
	"shared.local/pkg/logger"
	"shared.local/pkg/middleware"

	"github.com/gin-gonic/gin"
)

func run() error {
	cfg, err := config.LoadBase()
	if err != nil {
		return fmt.Errorf("load config failed: %w", err)
	}

	logger.Init(cfg.Log.Level)
	defer logger.L().Sync()

	r := gin.New()
	r.Use(
		middleware.TraceID(),
		middleware.Logger(),
		middleware.Recovery(),
	)

	// 依赖初始化
	tenantRepo := repo.NewTenantRepo()
	tenantService := service.NewTenantService(tenantRepo)
	tenantHandler := handler.NewTenantHandler(tenantService)
	tenantHandler.Register(r)

	addr := fmt.Sprintf(":%d", cfg.HTTP.Port)
	return r.Run(addr)
}
