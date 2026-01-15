package middleware

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"tenant-service/internal/trace"
)

const TraceIDKey = "trace_id"

func TraceID() gin.HandlerFunc {
	return func(c *gin.Context) {
		traceID := c.GetHeader("X-Trace-ID")
		if traceID == "" {
			traceID = uuid.NewString()
		}

		// 写入标准 context.Context
		ctx := context.WithValue(c.Request.Context(), trace.TraceIDKey(), traceID)
		c.Request = c.Request.WithContext(ctx)

		// 返回给客户端
		c.Writer.Header().Set("X-Trace-ID", traceID)

		c.Next()
	}
}
