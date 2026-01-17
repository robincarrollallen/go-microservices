package handler

import (
	"net"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"shared.local/pkg/logger"
	"shared.local/pkg/trace"

	"user-service/internal/service"
)

type UserHandler struct {
	service *service.UserService
}

func NewUserHandler(s *service.UserService) *UserHandler {
	return &UserHandler{service: s}
}

// Register 主注册方法，协调所有路由注册
func (h *UserHandler) Register(r *gin.Engine) {
	h.registerHealthRoutes(r)
	h.registerUserRoutes(r)
}

// 健康检查相关路由
func (h *UserHandler) registerHealthRoutes(r *gin.Engine) {
	r.GET("/health", func(c *gin.Context) {
		host := c.Request.Host
		if h, _, err := net.SplitHostPort(host); err == nil {
			host = h
		}

		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
			"host":   host,
		})
	})
}

// 用户相关路由
func (h *UserHandler) registerUserRoutes(r *gin.Engine) {
	r.GET("/user/:id", h.getUser)
}

// 用户查询
func (h *UserHandler) getUser(c *gin.Context) {
	ctx := c.Request.Context()

	logger.L().Info("get user in handler",
		zap.String("trace_id", trace.FromContext(ctx)),
		zap.String("user_id", c.Param("id")),
	)

	id := c.Param("id")

	user := h.service.GetUser(ctx, id)

	c.JSON(http.StatusOK, gin.H{
		"id":   id,
		"name": user,
	})
}
