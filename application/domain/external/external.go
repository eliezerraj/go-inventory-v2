package external

type ProductRequest struct {
	ID			int			`json:"id,omitempty"`
	Sku			string		`json:"sku,omitempty"`
	Type		string 		`json:"type,omitempty"`
	Name		string 		`json:"name,omitempty"`
	Status		string 		`json:"status,omitempty"`
	LeadTime	int			`json:"lead_time,omitempty"`
}

type ProductResponse struct {
	Response    string	`json:"response"`
	Product		any	`json:"product"`
}