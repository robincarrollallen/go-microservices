package handler

import (
	"fmt"
	"net/http"
	"tenant-service/internal/model/dto"
	"tenant-service/internal/service"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"shared.local/pkg/logger"
	"shared.local/pkg/pagination"
	"shared.local/pkg/response"
	"shared.local/pkg/trace"
)

type DomainHandler struct {
	service *service.DomainService
}

func NewDomainHandler(s *service.DomainService) *DomainHandler {
	return &DomainHandler{service: s}
}

func (h *DomainHandler) Register(r *gin.Engine) {
	h.registerDomainRoutes(r)
}

// 商户相关路由
func (h *DomainHandler) registerDomainRoutes(r *gin.Engine) {
	r.POST("/domain/create", h.createDomain)
	r.GET("/domains", h.getDomains)
}

func (h *DomainHandler) createDomain(c *gin.Context) {
	ctx := c.Request.Context()

	var req dto.CreateDomainRequest
	if err := c.ShouldBindJSON(&req); err != nil {

		logger.L().Warn("create domain request bind error",
			zap.String("trace_id", trace.FromContext(ctx)),
			zap.Any("body", err),
		)

		formattedErr := response.FormatValidationError(err)
		response.ErrorWithStatus(c, http.StatusBadRequest, 4000, formattedErr)
		return
	}

	logger.L().Info("create domain request",
		zap.String("trace_id", trace.FromContext(ctx)),
		zap.Any("body", req),
	)

	domain, err := h.service.CreateDomain(ctx, req)
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data":    domain,
	})
}

func (h *DomainHandler) getDomains(c *gin.Context) {
	ctx := c.Request.Context()

	// 解析分页参数，设置默认值
	page := 1
	pageSize := 10
	if p := c.DefaultQuery("page", "1"); p != "" {
		_, _ = parseScanInt(p, &page)
	}
	if ps := c.DefaultQuery("pageSize", "10"); ps != "" {
		_, _ = parseScanInt(ps, &pageSize)
	}

	// 解析可选的查询参数
	var tenantID *uint
	if tenantIDStr := c.Query("tenantID"); tenantIDStr != "" {
		var id uint
		if _, err := parseScanuint(tenantIDStr, &id); err == nil {
			tenantID = &id
		}
	}

	domain := c.Query("domain")

	// 构建请求
	req := dto.ListDomainsRequest{
		PaginationRequest: pagination.PaginationRequest{
			Page:     page,
			PageSize: pageSize,
		},
		TenantID: tenantID,
		Domain:   domain,
	}

	logger.L().Info("list domains request",
		zap.String("trace_id", trace.FromContext(ctx)),
		zap.Any("params", req),
	)

	// 调用 service
	result, err := h.service.ListDomains(ctx, req)
	if err != nil {
		c.Error(err)
		return
	}

	response.Success(c, result)
}

// 辅助函数：解析整数
func parseScanInt(s string, v *int) (int, error) {
	_, err := fmt.Sscanf(s, "%d", v)
	return *v, err
}

// 辅助函数：解析无符号整数
func parseScanuint(s string, v *uint) (uint, error) {
	_, err := fmt.Sscanf(s, "%d", v)
	return *v, err
}
