package handler

import (
	"errors"
	"fmt"
	"net"
	"net/http"
	"tenant-service/internal/model"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"shared.local/pkg/logger"
	"shared.local/pkg/trace"

	"tenant-service/internal/service"
)

type TenantHandler struct {
	service *service.TenantService
}

func NewTenantHandler(s *service.TenantService) *TenantHandler {
	return &TenantHandler{service: s}
}

// Register 主注册方法，协调所有路由注册
func (h *TenantHandler) Register(r *gin.Engine) {
	h.registerHealthRoutes(r)
	h.registerTenantRoutes(r)
}

// 健康检查相关路由
func (h *TenantHandler) registerHealthRoutes(r *gin.Engine) {
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

// 商户相关路由
func (h *TenantHandler) registerTenantRoutes(r *gin.Engine) {
	r.GET("/tenant/:id", h.getTenant)
}

// 商户查询
func (h *TenantHandler) getTenant(c *gin.Context) {
	ctx := c.Request.Context()

	logger.L().Info("get tenant in handler",
		zap.String("trace_id", trace.FromContext(ctx)),
		zap.String("user_id", c.Param("id")),
	)

	idStr := c.Param("id")
	var tenant *model.Tenant
	var err error

	if idStr != "" {
		var id uint
		if _, err := fmt.Sscanf(idStr, "%d", &id); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid tenant ID",
			})
			return
		}

		tenant, err = h.service.GetTenantByID(ctx, id)
	} else {
		host := c.Request.Host
		if h, _, err := net.SplitHostPort(host); err == nil {
			host = h
		}

		tenant, err = h.service.GetTenantByDomain(ctx, host)
	}

	if err != nil {
		if errors.Is(err, service.ErrTenantNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	c.JSON(http.StatusOK, tenant)
}
