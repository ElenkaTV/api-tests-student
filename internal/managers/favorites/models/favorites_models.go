package models

import advModels "api-tests-template/internal/managers/advertisements/models"

type FavoriteItem struct {
	FavoriteID    string                          `json:"favorite_id"`
	AddedAt       string                          `json:"added_at"`
	Advertisement advModels.AdvertisementResponse `json:"advertisement"`
}

type FavoritesList struct {
	Items      []FavoriteItem `json:"items"`
	Page       int            `json:"page"`
	PageSize   int            `json:"page_size"`
	Total      int            `json:"total"`
	TotalPages int            `json:"total_pages"`
}
