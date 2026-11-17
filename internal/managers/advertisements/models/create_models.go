package models

type PhotoResponse struct {
	ID        string `json:"id"`
	URL       string `json:"url"`
	SortOrder int    `json:"sort_order"`
}

type CreateAdvertisementResponse struct {
	ID          string          `json:"id"`
	Title       string          `json:"title"`
	Description string          `json:"description"`
	Price       float64         `json:"price"`
	Quantity    int             `json:"quantity"`
	Photos      []PhotoResponse `json:"photos"`
	UserID      string          `json:"user_id"`
}

type AdvertisementResponse struct {
	ID          string          `json:"id"`
	Title       string          `json:"title"`
	Description string          `json:"description"`
	Price       float64         `json:"price"`
	Quantity    int             `json:"quantity"`
	Photos      []PhotoResponse `json:"photos"`
	UserID      string          `json:"user_id"`
}

type ErrorResponse struct {
	Error   string      `json:"error"`
	Message string      `json:"message"`
	Details interface{} `json:"details"`
}
