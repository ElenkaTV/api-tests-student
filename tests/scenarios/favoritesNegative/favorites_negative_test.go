package favoritesNegative

import (
	"net/http"
	"os"
	"testing"

	"api-tests-template/internal/client/http/advertisements"
	"api-tests-template/internal/client/http/favorites"
	advManager "api-tests-template/internal/managers/advertisements"
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

func (s *TestSuite) SetupSuite() {
	utils.LoadEnv()
	os.Setenv("API_URL", "https://testboard.avito.com/api/v1")

	loginMain := os.Getenv("TEST_LOGIN")
	passMain := os.Getenv("TEST_PASSWORD")
	loginCreator := os.Getenv("TEST_LOGIN_2")
	passCreator := os.Getenv("TEST_PASSWORD_2")

	require.NotEmpty(s.T(), loginMain, "TEST_LOGIN не задан в .env.override")
	require.NotEmpty(s.T(), loginCreator, "TEST_LOGIN_2 не задан в .env.override")

	s.T().Log("Логинимся под основным пользователем")
	authMain := auth.Login(s.T(), loginMain, passMain)
	s.tokenMain = authMain.Token

	s.T().Log("Логинимся под вторым пользователем (создатель)")
	authCreator := auth.Login(s.T(), loginCreator, passCreator)
	s.tokenCreator = authCreator.Token

	// Создаем объявление от второго пользователя
	s.T().Log("Создаем объявление для тестов")
	fakePhoto := []byte{0xFF, 0xD8, 0xFF, 0xE0, 0x00, 0x10, 0x4A, 0x46, 0x49, 0x46}
	resp := advManager.CreateAdvertisement(s.T(), s.tokenCreator, "Объявление для негативных тестов", "Описание", 500, 1, fakePhoto, "test.jpg")
	s.adID = resp.ID
	require.NotEmpty(s.T(), s.adID)
}

func (s *TestSuite) TearDownSuite() {
	if s.adID != "" {
		s.T().Log("Удаляем тестовое объявление")
		advertisements.HttpDeleteAdvertisement(s.T(), s.tokenCreator, s.adID)
	}
}

// 1. Добавление в избранное несуществующего ad_id (ожидаем 404)
func (s *TestSuite) Test_AddNonExistentFavorite() {
	// Используем валидный по формату, но несуществующий UUID
	nonExistentID := "11111111-1111-1111-1111-111111111111"
	raw := favorites.HttpAddToFavorites(s.T(), s.tokenMain, nonExistentID)
	defer raw.Body.Close()

	require.Equal(s.T(), http.StatusNotFound, raw.StatusCode, "Ожидался статус 404 для несуществующего ad_id")
}

// 2. Удаление из избранного несуществующего ad_id (ожидаем 404)
func (s *TestSuite) Test_RemoveNonExistentFavorite() {
	nonExistentID := "11111111-1111-1111-1111-111111111111"
	raw := favorites.HttpRemoveFromFavorites(s.T(), s.tokenMain, nonExistentID)
	defer raw.Body.Close()

	require.Equal(s.T(), http.StatusNotFound, raw.StatusCode, "Ожидался статус 404 при удалении несуществующего ad_id")
}

// 3. Запрос без авторизации на получение избранного (ожидаем 401)
func (s *TestSuite) Test_GetFavoritesWithoutAuth() {
	raw := favorites.HttpGetFavorites(s.T(), "")
	defer raw.Body.Close()

	require.Equal(s.T(), http.StatusUnauthorized, raw.StatusCode, "Ожидался статус 401 без авторизации")
}

// 4. Повторное добавление одного и того же объявления (ожидаем 409 )
func (s *TestSuite) Test_AddDuplicateFavorite() {
	// Сначала добавляем объявление в избранное
	s.T().Log("Добавляем объявление в избранное (первый раз)")
	rawFirst := favorites.HttpAddToFavorites(s.T(), s.tokenMain, s.adID)
	defer rawFirst.Body.Close()
	require.Equal(s.T(), http.StatusCreated, rawFirst.StatusCode, "Не удалось добавить в избранное (первый раз)")

	// Пытаемся добавить его же еще раз
	s.T().Log("Пытаемся добавить то же объявление в избранное (второй раз)")
	rawSecond := favorites.HttpAddToFavorites(s.T(), s.tokenMain, s.adID)
	defer rawSecond.Body.Close()

	// Ожидаем 409 
	require.Equal(s.T(), http.StatusConflict, rawSecond.StatusCode, "Ожидался статус 409 при повторном добавлении")

	// Убираем за собой 
	favorites.HttpRemoveFromFavorites(s.T(), s.tokenMain, s.adID)
}

func TestFavoritesNegativeSuite(t *testing.T) {
	suite.Run(t, new(TestSuite))
}
