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
	Price		*Price		`json:"price,omitempty"`
	Inventory	*Inventory	`json:"inventory,omitempty"`
	CreatedAt	*time.Time 	`json:"created_at,omitempty"`
	ExpiresAt	*time.Time 	`json:"expires_at,omitempty"`
	UpdatedAt	*time.Time 	`json:"updated_at,omitempty"`	
}

type Inventory struct {
	ID				int			`json:"id,omitempty"`
	ProductId		int			`json:"fk_product_id,omitempty"`
	Available		int			`json:"available,omitempty"`
	Pending			int			`json:"pending,omitempty"`
	Sold			int			`json:"sold,omitempty"`
	CreatedAt		*time.Time 	`json:"created_at,omitempty"`
	UpdatedAt		*time.Time 	`json:"updated_at,omitempty"`	
}

type Price struct {
	ID				int			`json:"id,omitempty"`
	ProductId 		int		 	`json:"product_id,omitempty"`
	Currency		string		`json:"currency,omitempty"`
	Amount			float64		`json:"amount,omitempty"`
	StartedAt		*time.Time 	`json:"started_at,omitempty"`
	EndedAt			*time.Time 	`json:"ended_at,omitempty"`
	CreatedAt		*time.Time 	`json:"created_at,omitempty"`
	UpdatedAt		*time.Time 	`json:"updated_at,omitempty"`	
}