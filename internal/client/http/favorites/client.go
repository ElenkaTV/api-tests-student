package favorites

import (
	"net/http"
	"testing"

	"api-tests-template/internal/constants/path"
	apiRunner "api-tests-template/internal/helpers/api-runner"
)

// HttpAddToFavorites добавляет объявление в избранное
func HttpAddToFavorites(t *testing.T, token, adID string) *http.Response {
	return apiRunner.GetRunner().Auth(token).Create().
		Post(path.FavoritesPath).
		JSON(map[string]string{"advertisement_id": adID}).
		Expect(t).
		End().Response
}

// HttpGetFavorites получает список избранного
func HttpGetFavorites(t *testing.T, token string) *http.Response {
	return apiRunner.GetRunner().Auth(token).Create().
		Get(path.FavoritesPath).
		Expect(t).
		End().Response // Убрали Status(http.StatusOK)
}

// HttpRemoveFromFavorites удаляет из избранного по ID объявления
func HttpRemoveFromFavorites(t *testing.T, token, adID string) *http.Response {
	return apiRunner.GetRunner().Auth(token).Create().
		Delete(path.FavoritesPath + "/" + adID).
		Expect(t).
		End().Response
}
