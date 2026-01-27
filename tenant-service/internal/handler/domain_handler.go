package handler

import (
	"tenant-service/internal/service"

	"github.com/gin-gonic/gin"
)

type DomainHandler struct {
	service *service.DomainService
}

func NewDomainHandler(s *service.DomainService) *DomainHandler {
	return &DomainHandler{service: s}
}

func (h *DomainHandler) Register(r *gin.Engine) {

}
