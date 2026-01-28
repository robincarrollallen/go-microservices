package handler

import (
	"net/http"
	"tenant-service/internal/model/dto"
	"tenant-service/internal/service"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"shared.local/pkg/logger"
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
