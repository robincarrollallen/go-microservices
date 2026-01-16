package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"shared.local/pkg/trace"
)

const TraceIDHeader = "X-Trace-ID"

func TraceID() gin.HandlerFunc {
	return func(c *gin.Context) {
		traceID := c.GetHeader(TraceIDHeader)
		if traceID == "" {
			traceID = uuid.NewString()
		}

		// 写入标准 context.Context
		ctx := trace.WithTraceID(c.Request.Context(), traceID)
		c.Request = c.Request.WithContext(ctx)

		// 返回给客户端
		c.Writer.Header().Set(TraceIDHeader, traceID)

		c.Next()
	}
}
