package handler

import (
	"fmt"
	"net"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"shared.local/pkg/logger"
	"shared.local/pkg/response"
	"shared.local/pkg/trace"

	"tenant-service/internal/model/dto"
	"tenant-service/internal/model/entity"
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
	r.GET("/tenant/info", h.getTenant)
	r.POST("/tenant/create", h.createTenant)
}

// 商户创建
func (h *TenantHandler) createTenant(c *gin.Context) {
	ctx := c.Request.Context()

	var req dto.CreateTenantRequest
	if err := c.ShouldBindJSON(&req); err != nil {

		logger.L().Warn("raw request body",
			zap.String("trace_id", trace.FromContext(ctx)),
			zap.Any("body", err),
		)

		formattedErr := response.FormatValidationError(err)
		response.ErrorWithStatus(c, http.StatusBadRequest, 4000, formattedErr)
		return
	}

	logger.L().Info("create tenant request",
		zap.String("trace_id", trace.FromContext(ctx)),
		zap.Any("body", req),
	)

	tenant, err := h.service.CreateTenant(ctx, req)
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data":    tenant,
	})
}

// 商户查询
func (h *TenantHandler) getTenant(c *gin.Context) {
	ctx := c.Request.Context()

	logger.L().Info("get tenant in handler",
		zap.String("trace_id", trace.FromContext(ctx)),
		zap.String("user_id", c.Param("id")),
	)

	idStr := c.Param("id")
	var tenant *entity.Tenant
	var err error

	if idStr != "" {
		var id uint
		if _, err := fmt.Sscanf(idStr, "%d", &id); err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
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
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data":    tenant,
	})
}
