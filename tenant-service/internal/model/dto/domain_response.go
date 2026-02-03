package dto

import "time"

type DomainResponse struct {
	ID        uint      `json:"id"`
	Domain    string    `json:"domain"`
	Status    uint8     `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type ListDomainsResponse struct {
	Items      []DomainResponse `json:"items"`
	Total      int64            `json:"total"`
	Page       int              `json:"page"`
	PageSize   int              `json:"pageSize"`
	TotalPages int              `json:"totalPages"`
}
