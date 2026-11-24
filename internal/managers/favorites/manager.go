package favorites

import (
	"encoding/json"
	"io"
	"testing"

	"api-tests-template/internal/client/http/favorites"
	"api-tests-template/internal/managers/favorites/models"

	"github.com/stretchr/testify/require"
)

// GetFavoritesList получает список избранного пользователя и возвращает его модель
func GetFavoritesList(t *testing.T, token string) models.FavoritesList {
	raw := favorites.HttpGetFavorites(t, token)
	defer raw.Body.Close()

	body, err := io.ReadAll(raw.Body)
	require.NoError(t, err)

	var resp models.FavoritesList
	require.NoError(t, json.Unmarshal(body, &resp))
	return resp
}
