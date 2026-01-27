package bootstrap

import (
	"tenant-service/internal/handler"
)

// Container 依赖容器，存储所有模块的处理器
type Container struct {
	TenantHandler *handler.TenantHandler
	DomainHandler *handler.DomainHandler
}

