package handler

import (
	"net/http"
	"tenant-service/internal/service"

	"github.com/gin-gonic/gin"
)

type TenantHandler struct {
	service *service.TenantService
}

func NewCommonHandler(s *service.TenantService) *TenantHandler {
	return &TenantHandler{service: s}
}

func (h *TenantHandler) Register(r *gin.Engine) {
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	r.GET("/common/:id", h.getTenant)
}

func (h *TenantHandler) getTenant(c *gin.Context) {
	id := c.Param("id")
	user := h.service.GetCommon(id)

	c.JSON(http.StatusOK, gin.H{
		"id":   id,
		"name": user,
	})
}
