package advertisements

import (
	"encoding/json"
	"io"
	"testing"

	"api-tests-template/internal/client/http/advertisements"
	"api-tests-template/internal/managers/advertisements/models"

	"github.com/stretchr/testify/require"
)

// CreateAdvertisement создает объявление и возвращает ответ с его ID
func CreateAdvertisement(t *testing.T, token, title, desc string, price float64, quantity int, photoBytes []byte, photoName string) models.CreateAdvertisementResponse {
	raw := advertisements.HttpCreateAdvertisement(t, token, title, desc, price, quantity, photoBytes, photoName)
	defer raw.Body.Close()

	body, err := io.ReadAll(raw.Body)
	require.NoError(t, err)

	var resp models.CreateAdvertisementResponse
	require.NoError(t, json.Unmarshal(body, &resp))
	return resp
}

// GetAdvertisementByID получает объявление по ID и возвращает его модель
func GetAdvertisementByID(t *testing.T, id string) models.AdvertisementResponse {
	raw := advertisements.HttpGetAdvertisementByID(t, id)
	defer raw.Body.Close()

	body, err := io.ReadAll(raw.Body)
	require.NoError(t, err)

	var resp models.AdvertisementResponse
	require.NoError(t, json.Unmarshal(body, &resp))
	return resp
}
