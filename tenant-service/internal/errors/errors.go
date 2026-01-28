package apperror

import (
	"net/http"

	"shared.local/pkg/response"
)

// 预定义的常用错误
var (
	ErrNameExists     = response.NewAppError(4001, "Tenant name already exists", http.StatusConflict)
	ErrTenantNotFound = response.NewAppError(4004, "Tenant not found", http.StatusNotFound)
	ErrDomainExists   = response.NewAppError(4002, "Domain already exists", http.StatusConflict)
)
