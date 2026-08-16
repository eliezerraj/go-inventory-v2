package entity

import (
	"time"
)

type Product struct {
	ID			int			`json:"id,omitempty"`
	Sku			string		`json:"sku,omitempty"`
	Type		string 		`json:"type,omitempty"`
	Name		string 		`json:"name,omitempty"`
	Status		string 		`json:"status,omitempty"`
	LeadTime	int			`json:"lead_time,omitempty"`
	CreatedAt	*time.Time 	`json:"created_at,omitempty"`
	UpdatedAt	*time.Time 	`json:"updated_at,omitempty"`	
}