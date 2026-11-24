package favorites

import (
	"os"
	"path/filepath"
	"testing"

	"api-tests-template/internal/client/http/advertisements"
	"api-tests-template/internal/client/http/favorites"
	advManager "api-tests-template/internal/managers/advertisements"
	"api-tests-template/internal/managers/auth"
	favManager "api-tests-template/internal/managers/favorites"
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

// SetupSuite подготавливает окружение: логин двух пользователей
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
}

// SetupTest создает новое объявление перед каждым тестом
func (s *TestSuite) SetupTest() {
	s.T().Log("Создаем объявление для текущего теста")

	photoPath := filepath.Join("..", "..", "testdata", "test.jpg")
	photoBytes, err := os.ReadFile(photoPath)
	require.NoError(s.T(), err, "Не удалось прочитать файл картинки: %s", photoPath)

	resp := advManager.CreateAdvertisement(s.T(), s.tokenCreator, "Объявление для избранного", "Описание", 500, 1, photoBytes, "test.jpg")
	s.adID = resp.ID
	require.NotEmpty(s.T(), s.adID)
}

// TearDownTest удаляет объявление после каждого теста
func (s *TestSuite) TearDownTest() {
	if s.adID != "" {
		// ВРЕМЕННО ЗАКОММЕНТИРОВАНО, ЧТОБЫ ОБЪЯВЛЕНИЕ ОСТАЛОСЬ НА САЙТЕ
		advertisements.HttpDeleteAdvertisement(s.T(), s.tokenCreator, s.adID)
		s.adID = ""
	}
}

// Test_FavoritesCycle проверяет полный цикл работы с избранным
func (s *TestSuite) Test_FavoritesCycle() {
	s.T().Log("Добавляем в избранное")
	favorites.HttpAddToFavorites(s.T(), s.tokenMain, s.adID)

	s.Run("Объявление в списке избранного", func() {
		list := favManager.GetFavoritesList(s.T(), s.tokenMain)
		require.NotEmpty(s.T(), list.Items)
	})

	s.T().Log("Удаляем из избранного")
	favorites.HttpRemoveFromFavorites(s.T(), s.tokenMain, s.adID)

	s.Run("Объявление пропало из списка", func() {
		list := favManager.GetFavoritesList(s.T(), s.tokenMain)
		found := false
		for _, item := range list.Items {
			if item.Advertisement.ID == s.adID {
				found = true
				break
			}
		}
		require.False(s.T(), found, "Объявление не должно быть в избранном")
	})
}

// Test_DeleteAdvertisementRemovesFromFavorites проверяет, что при удалении объявления оно пропадает из избранного
func (s *TestSuite) Test_DeleteAdvertisementRemovesFromFavorites() {
	s.T().Log("Добавляем в избранное")
	favorites.HttpAddToFavorites(s.T(), s.tokenMain, s.adID)

	s.T().Log("Удаляем само объявление от лица владельца")
	// ВРЕМЕННО ЗАКОММЕНТИРОВАНО, ЧТОБЫ ОБЪЯВЛЕНИЕ ОСТАЛОСЬ НА САЙТЕ
	advertisements.HttpDeleteAdvertisement(s.T(), s.tokenCreator, s.adID)

	s.T().Log("Проверяем, что объявление пропало из избранного")
	list := favManager.GetFavoritesList(s.T(), s.tokenMain)
	found := false
	for _, item := range list.Items {
		if item.Advertisement.ID == s.adID {
			found = true
			break
		}
	}
	require.False(s.T(), found, "Объявление должно быть удалено из избранного после удаления самого объявления")

	// Обнуляем ID, чтобы TearDownTest не пытался удалить его снова
	s.adID = ""
}

// TestFavoritesSuite запускает тестовый набор
func TestFavoritesSuite(t *testing.T) {
	suite.Run(t, new(TestSuite))
}
