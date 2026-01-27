package dto

import "time"

type TenantResponse struct {
	ID        uint        `json:"id"`
	Name      string      `json:"name"`
	Status    string      `json:"status"`
	Domains   []DomainDTO `json:"domains"`
	CreatedAt time.Time   `json:"created_at"`
	UpdatedAt time.Time   `json:"updated_at"`
}
