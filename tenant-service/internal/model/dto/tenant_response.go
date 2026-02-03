package dto

import (
	"time"

	"shared.local/pkg/pagination"
)

type TenantResponse struct {
	ID        uint      `json:"id"`
	Name      string    `json:"name"`
	Status    uint8     `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type ListTenantResponse = pagination.PaginationResponse[TenantResponse]
