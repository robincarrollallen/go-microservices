package bootstrap

import (
	"tenant-service/internal/handler"
	"tenant-service/internal/repo"
	"tenant-service/internal/service"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// SetupDependencies 初始化所有依赖并返回容器
func SetupDependencies(db *gorm.DB) *Container {
	return &Container{
		TenantHandler: initTenantModule(db),
		DomainHandler: initDomainModule(db),
	}
}

// initTenantModule 初始化租户模块
func initTenantModule(db *gorm.DB) *handler.TenantHandler {
	tenantRepo := repo.NewTenantRepo(db)
	tenantService := service.NewTenantService(tenantRepo)
	return handler.NewTenantHandler(tenantService)
}

// initDomainModule 初始化域名模块
func initDomainModule(db *gorm.DB) *handler.DomainHandler {
	domainRepo := repo.NewDomainRepo(db)
	domainService := service.NewDomainService(domainRepo)
	return handler.NewDomainHandler(domainService)
}

// RegisterRoutes 注册所有模块的路由
func RegisterRoutes(c *Container, r *gin.Engine) {
	c.TenantHandler.Register(r)
	c.DomainHandler.Register(r)
}

