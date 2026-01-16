package main

import (
	"fmt"

	"shared.local/pkg/config"
	"shared.local/pkg/logger"
	"shared.local/pkg/middleware"
	"shared.local/pkg/response"

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

	// 健康检查
	r.GET("/health", func(c *gin.Context) {
		response.Success(c, gin.H{"status": "ok"})
	})

	// TODO: 添加用户相关的 handler

	addr := fmt.Sprintf(":%d", cfg.HTTP.Port)
	return r.Run(addr)
}
