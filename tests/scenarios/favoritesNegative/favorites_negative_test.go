package favoritesNegative

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"testing"

	"api-tests-template/internal/client/http/advertisements"
	"api-tests-template/internal/client/http/favorites"
	advManager "api-tests-template/internal/managers/advertisements"
	"api-tests-template/internal/managers/advertisements/models"
	"api-tests-template/internal/managers/auth"
	"api-tests-template/internal/utils"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type TestSuite struct {
	suite.Suite
	tokenMain    string
	tokenCreator string
	adID         string
}

// SetupSuite подготавливает окружение: логин двух пользователей и создание объявления
func (s *TestSuite) SetupSuite() {
	utils.LoadEnv()
	os.Setenv("API_URL", "https://testboard.avito.com/api/v1")

	loginMain := os.Getenv("TEST_LOGIN")
	passMain := os.Getenv("TEST_PASSWORD")
	loginCreator := os.Getenv("TEST_LOGIN_2")
	passCreator := os.Getenv("TEST_PASSWORD_2")

	require.NotEmpty(s.T(), loginMain, "TEST_LOGIN не задан")
	require.NotEmpty(s.T(), loginCreator, "TEST_LOGIN_2 не задан")

	authMain := auth.Login(s.T(), loginMain, passMain)
	s.tokenMain = authMain.Token

	authCreator := auth.Login(s.T(), loginCreator, passCreator)
	s.tokenCreator = authCreator.Token

	s.T().Log("Читаем тестовую картинку из файла")
	photoPath := filepath.Join("..", "..", "testdata", "test.jpg")
	photoBytes, err := os.ReadFile(photoPath)
	require.NoError(s.T(), err, "Не удалось прочитать файл картинки: %s", photoPath)

	s.T().Log("Создаем объявление для тестов")
	resp := advManager.CreateAdvertisement(s.T(), s.tokenCreator, "Объявление для негативных тестов", "Описание", 500, 1, photoBytes, "test.jpg")
	s.adID = resp.ID
	require.NotEmpty(s.T(), s.adID)
}

// TearDownSuite удаляет тестовое объявление после завершения тестов
func (s *TestSuite) TearDownSuite() {
	if s.adID != "" {
		advertisements.HttpDeleteAdvertisement(s.T(), s.tokenCreator, s.adID)
	}
}

// Test_AddNonExistentFavorite проверяет ошибку 404 при добавлении несуществующего объявления
func (s *TestSuite) Test_AddNonExistentFavorite() {
	nonExistentID := "11111111-1111-1111-1111-111111111111"
	raw := favorites.HttpAddToFavorites(s.T(), s.tokenMain, nonExistentID)
	defer raw.Body.Close()

	require.Equal(s.T(), http.StatusNotFound, raw.StatusCode)
	body, _ := io.ReadAll(raw.Body)
	var errResp models.ErrorResponse
	_ = json.Unmarshal(body, &errResp)
	require.Equal(s.T(), "not_found", errResp.Error)
}

// Test_RemoveNonExistentFavorite проверяет ошибку 404 при удалении несуществующего объявления
func (s *TestSuite) Test_RemoveNonExistentFavorite() {
	nonExistentID := "11111111-1111-1111-1111-111111111111"
	raw := favorites.HttpRemoveFromFavorites(s.T(), s.tokenMain, nonExistentID)
	defer raw.Body.Close()

	require.Equal(s.T(), http.StatusNotFound, raw.StatusCode)
	body, _ := io.ReadAll(raw.Body)
	var errResp models.ErrorResponse
	_ = json.Unmarshal(body, &errResp)
	require.Equal(s.T(), "not_found", errResp.Error)
}

// Test_GetFavoritesWithoutAuth проверяет ошибку 401 без авторизации
func (s *TestSuite) Test_GetFavoritesWithoutAuth() {
	raw := favorites.HttpGetFavorites(s.T(), "")
	defer raw.Body.Close()
	require.Equal(s.T(), http.StatusUnauthorized, raw.StatusCode)
}

// Test_AddFavoriteWithoutAuth проверяет ошибку 401 при добавлении в избранное без авторизации
func (s *TestSuite) Test_AddFavoriteWithoutAuth() {
	raw := favorites.HttpAddToFavorites(s.T(), "", s.adID)
	defer raw.Body.Close()
	require.Equal(s.T(), http.StatusUnauthorized, raw.StatusCode)
}

// Test_RemoveFavoriteWithoutAuth проверяет ошибку 401 при удалении из избранного без авторизации
func (s *TestSuite) Test_RemoveFavoriteWithoutAuth() {
	raw := favorites.HttpRemoveFromFavorites(s.T(), "", s.adID)
	defer raw.Body.Close()
	require.Equal(s.T(), http.StatusUnauthorized, raw.StatusCode)
}

// Test_AddDuplicateFavorite проверяет ошибку 409 при повторном добавлении в избранное
func (s *TestSuite) Test_AddDuplicateFavorite() {
	favorites.HttpAddToFavorites(s.T(), s.tokenMain, s.adID)

	raw := favorites.HttpAddToFavorites(s.T(), s.tokenMain, s.adID)
	defer raw.Body.Close()
	require.Equal(s.T(), http.StatusConflict, raw.StatusCode)

	body, _ := io.ReadAll(raw.Body)
	var errResp models.ErrorResponse
	_ = json.Unmarshal(body, &errResp)
	require.Equal(s.T(), "already_exists", errResp.Error)

	favorites.HttpRemoveFromFavorites(s.T(), s.tokenMain, s.adID)
}

// TestFavoritesNegativeSuite запускает набор негативных тестов
func TestFavoritesNegativeSuite(t *testing.T) {
	suite.Run(t, new(TestSuite))
}
