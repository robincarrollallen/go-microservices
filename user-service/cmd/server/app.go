package main

import (
	"fmt"

	"user-service/internal/handler"
	"user-service/internal/repo"
	"user-service/internal/service"

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
	userRepo := repo.NewUserRepo()
	userService := service.NewUserService(userRepo)
	userHandler := handler.NewUserHandler(userService)
	userHandler.Register(r)

	addr := fmt.Sprintf(":%d", cfg.HTTP.Port)
	return r.Run(addr)
}
