package dto

import "time"

type DomainDTO struct {
	ID        uint      `json:"id"`
	Domain    string    `json:"domain"`
	CreatedAt time.Time `json:"created_at"`
}
