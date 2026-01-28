package middleware

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"shared.local/pkg/logger"
	"shared.local/pkg/response"
	"shared.local/pkg/trace"
)

// ErrorHandler 通用错误处理中间件（处理 AppError 和通用错误）
func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		// 检查是否有错误
		if len(c.Errors) > 0 {
			err := c.Errors.Last().Err
			handleCommonError(c, err)
		}
	}
}

func handleCommonError(c *gin.Context, err error) {
	traceID := trace.FromContext(c.Request.Context())

	logger.L().Warn("system error",
		zap.String("trace_id", traceID),
		zap.Error(err),
	)

	// 检查是否是自定义 AppError
	var appErr *response.AppError
	if errors.As(err, &appErr) {
		logger.L().Warn("business error",
			zap.String("trace_id", traceID),
			zap.Int("error_code", appErr.Code),
			zap.String("message", appErr.Message),
		)
		response.ErrorWithStatus(c, appErr.Status, appErr.Code, appErr.Message)
		c.Errors = nil // 错误已处理，清除错误堆栈
		return
	}

	response.ErrorWithStatus(c, http.StatusInternalServerError, 5000, err.Error())
	c.Errors = nil // 错误已处理，清除错误堆栈
}

