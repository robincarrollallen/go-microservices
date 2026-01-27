package dto

import "time"

type TenantResponse struct {
	ID        uint      `json:"id"`
	Name      string    `json:"name"`
	Status    uint8     `json:"status"`
	Domains   []string  `json:"domains"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
