package middleware

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"tenant-service/internal/errors"

	"shared.local/pkg/logger"
	"shared.local/pkg/response"
	"shared.local/pkg/trace"
)

// ServiceErrorHandler tenant-service 业务错误处理中间件
func ServiceErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		// 检查是否有错误
		if len(c.Errors) > 0 {
			err := c.Errors.Last().Err
			handleServiceError(c, err)
		}
	}
}

func handleServiceError(c *gin.Context, err error) {
	traceID := trace.FromContext(c.Request.Context())

	// 检查特定的业务错误
	if errors.Is(err, apperror.ErrNameExists) {
		logger.L().Warn("duplicate tenant name",
			zap.String("trace_id", traceID),
		)
		response.ErrorWithStatus(c, http.StatusConflict, 4001, "Tenant name already exists")
		c.Errors = nil // 错误已处理，清除错误堆栈，防止被后续中间件重复处理
		return
	}

	if errors.Is(err, apperror.ErrTenantNotFound) {
		logger.L().Warn("tenant not found",
			zap.String("trace_id", traceID),
		)
		response.ErrorWithStatus(c, http.StatusNotFound, 4004, "Tenant not found")
		c.Errors = nil // 错误已处理，清除错误堆栈，防止被后续中间件重复处理
		return
	}

	if errors.Is(err, apperror.ErrDomainExists) {
		logger.L().Warn("duplicate domain",
			zap.String("trace_id", traceID),
		)
		response.ErrorWithStatus(c, http.StatusConflict, 4002, "Domain already exists")
		c.Errors = nil // 错误已处理，清除错误堆栈，防止被后续中间件重复处理
		return
	}

	// 如果不是本服务的业务错误，传给下一个中间件处理
	// 这里什么都不做，让错误继续传播
}
