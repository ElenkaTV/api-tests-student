package models

// PhotoResponse — структура фото в ответе
type PhotoResponse struct {
	ID        string `json:"id"`
	URL       string `json:"url"`
	SortOrder int    `json:"sort_order"`
}

// CreateAdvertisementResponse — ответ сервера после создания объявления
type CreateAdvertisementResponse struct {
	ID          string          `json:"id"`
	Title       string          `json:"title"`
	Description string          `json:"description"`
	Price       float64         `json:"price"`
	Quantity    int             `json:"quantity"`
	Photos      []PhotoResponse `json:"photos"`
	UserID      string          `json:"user_id"`
}

// AdvertisementResponse — ответ для GET /advertisement?id=
type AdvertisementResponse struct {
	ID          string          `json:"id"`
	Title       string          `json:"title"`
	Description string          `json:"description"`
	Price       float64         `json:"price"`
	Quantity    int             `json:"quantity"`
	Photos      []PhotoResponse `json:"photos"`
	UserID      string          `json:"user_id"`
}

// ErrorResponse — структура ошибки
type ErrorResponse struct {
	Error   string      `json:"error"`
	Message string      `json:"message"`
	Details interface{} `json:"details"`
}

// AdvertisementsList — структура ответа для поиска объявлений
type AdvertisementsList struct {
	Items []Advertisement `json:"items"`
}

// Advertisement — модель объявления для поиска
type Advertisement struct {
	ID          string          `json:"id"`
	Title       string          `json:"title"`
	Description string          `json:"description"`
	Price       float64         `json:"price"`
	Quantity    int             `json:"quantity"`
	Photos      []PhotoResponse `json:"photos"`
}
