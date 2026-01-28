package middleware

import (
	"errors"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

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

	logger.L().Warn("service error",
		zap.String("trace_id", traceID),
		zap.Error(err),
	)

	// 检查是否实现了 APIError 接口
	var apiErr response.APIError
	if errors.As(err, &apiErr) {
		logger.L().Warn("service error",
			zap.String("trace_id", traceID),
			zap.Int("error_code", apiErr.GetCode()),
			zap.String("message", apiErr.GetMessage()),
		)
		response.ErrorWithStatus(c, apiErr.GetStatus(), apiErr.GetCode(), apiErr.GetMessage())
		c.Errors = nil // 错误已处理，清除错误堆栈，防止被后续中间件重复处理
		return
	}

	// 不是业务错误，传给下一个中间件处理
	// 这里什么都不做，让错误继续传播
}
