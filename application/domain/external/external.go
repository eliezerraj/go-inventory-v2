package external

type ProductRequest struct {
	ID			int				`json:"id,omitempty"`
	Sku			string			`json:"sku,omitempty"`
	Type		string 			`json:"type,omitempty"`
	Name		string 			`json:"name,omitempty"`
	Status		string 			`json:"status,omitempty"`
	LeadTime	int					`json:"lead_time" validate:"required,gte=0"`
	Inventory	*InventoryRequest	`json:"inventory,omitempty"`
	Price		*PriceRequest		`json:"price,omitempty"`
}

type InventoryRequest struct {
	Available		int		`json:"available,omitempty"`
	Pending			int		`json:"pending,omitempty"`
	Sold			int		`json:"sold,omitempty"`
}

type PriceRequest struct {
    Currency string  `json:"currency" validate:"required"`
    Amount   float64 `json:"amount" validate:"required,gte=0"`
}

type ProductResponse struct {
	Response    string	`json:"response"`
	Product		any		`json:"product"`
}