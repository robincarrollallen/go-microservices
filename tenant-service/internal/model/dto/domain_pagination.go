package dto

import "shared.local/pkg/pagination"

type ListDomainsRequest struct {
	pagination.PaginationRequest
	TenantID *uint  `form:"tenantID"`
	Domain   string `form:"domain"`
}

type ListDomainsResponse struct {
	Items      []DomainResponse `json:"items"`
	Total      int64            `json:"total"`
	Page       int              `json:"page"`
	PageSize   int              `json:"pageSize"`
	TotalPages int              `json:"totalPages"`
}
