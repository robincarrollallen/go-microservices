package dto

import "shared.local/pkg/pagination"

type CreateDomainRequest struct {
	TenantID uint   `form:"tenantID" binding:"required"`             // 关联商户ID: 必传
	Domain   string `form:"domain" binding:"required,min=1,max=100"` // 域名: 必传, 长度1-100
	Status   *uint8 `form:"status"`                                  // 指针类型，用于判断是否传值
}

type ListDomainsRequest struct {
	pagination.PaginationRequest
	TenantID *uint  `form:"tenantID"`
	Domain   string `form:"domain"`
}
