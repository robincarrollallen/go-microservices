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
	tenantRepo := repo.NewTenantRepo(db)
	domainRepo := repo.NewDomainRepo(db)

	return &Container{
		TenantHandler: initTenantModule(tenantRepo),
		DomainHandler: initDomainModule(domainRepo, tenantRepo),
	}
}

// initTenantModule 初始化租户模块
func initTenantModule(tenantRepo *repo.TenantRepo) *handler.TenantHandler {
	tenantService := service.NewTenantService(tenantRepo)
	return handler.NewTenantHandler(tenantService)
}

// initDomainModule 初始化域名模块
func initDomainModule(domainRepo *repo.DomainRepo, tenantRepo *repo.TenantRepo) *handler.DomainHandler {
	domainService := service.NewDomainService(domainRepo, tenantRepo)
	return handler.NewDomainHandler(domainService)
}

// RegisterRoutes 注册所有模块的路由
func RegisterRoutes(c *Container, r *gin.Engine) {
	c.TenantHandler.Register(r)
	c.DomainHandler.Register(r)
}

