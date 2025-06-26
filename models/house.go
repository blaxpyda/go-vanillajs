package models

type House struct {
	ID          int     `json:"id"`
	Name		string  `json:"name"`
	Description string  `json:"description"`
	HouseTypeID int     `json:"house_type_id"`
	Price       float64 `json:"price"`
	Tags		[]string `json:"tags"`
	ImageURL	*string `json:"image_url"` // nullable
	CreatedAt   string  `json:"created_at"`
	UpdatedAt   string  `json:"updated_at"`
	AgentID     int     `json:"agent_id"`
}