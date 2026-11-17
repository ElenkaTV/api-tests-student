package favorites

import (
	"os"
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
	tokenMain    string // Основной пользователь (добавляет в избранное)
	tokenCreator string // Второй пользователь (создает объявление)
	adID         string // ID созданного объявления
}

func (s *TestSuite) SetupSuite() {
	// 1. Принудительно загружаем .env.override
	utils.LoadEnv()
	os.Setenv("API_URL", "https://testboard.avito.com/api/v1")

	// 2. Читаем данные из окружения
	loginMain := os.Getenv("TEST_LOGIN")
	passMain := os.Getenv("TEST_PASSWORD")

	loginCreator := os.Getenv("TEST_LOGIN_2")
	passCreator := os.Getenv("TEST_PASSWORD_2")

	// Проверяем, что данные заданы (если нет — тест упадет с понятной ошибкой)
	require.NotEmpty(s.T(), loginMain, "TEST_LOGIN не задан в .env.override")
	require.NotEmpty(s.T(), loginCreator, "TEST_LOGIN_2 не задан в .env.override")

	// 3. Логинимся под ОСНОВНЫМ пользователем (кто будет добавлять в избранное)
	s.T().Log("Логинимся под основным пользователем")
	authMain := auth.Login(s.T(), loginMain, passMain)
	s.tokenMain = authMain.Token

	// 4. Логинимся под ВТОРЫМ пользователем (создатель объявления)
	s.T().Log("Логинимся под вторым пользователем (создатель)")
	authCreator := auth.Login(s.T(), loginCreator, passCreator)
	s.tokenCreator = authCreator.Token

	// 5. Создаем объявление от лица ВТОРОГО пользователя
	s.T().Log("Создаем объявление от имени второго пользователя")
	// Имитация картинки (заглушка)
	fakePhoto := []byte{0xFF, 0xD8, 0xFF, 0xE0, 0x00, 0x10, 0x4A, 0x46, 0x49, 0x46}

	resp := advManager.CreateAdvertisement(s.T(), s.tokenCreator, "Объявление для избранного", "Описание", 500, 1, fakePhoto, "test.jpg")
	s.adID = resp.ID
	require.NotEmpty(s.T(), s.adID, "Не удалось создать объявление")
}

func (s *TestSuite) TearDownSuite() {
	if s.adID != "" {
		s.T().Log("Удаляем тестовое объявление (под вторым пользователем)")
		// Удаляем под токеном создателя
		advertisements.HttpDeleteAdvertisement(s.T(), s.tokenCreator, s.adID)
	}
}

func (s *TestSuite) Test_FavoritesCycle() {
	// 1. Добавляем в избранное (под ОСНОВНЫМ пользователем)
	s.T().Log("Добавляем чужое объявление в избранное (под основным пользователем)")
	favorites.HttpAddToFavorites(s.T(), s.tokenMain, s.adID)

	// 2. Проверяем, что оно в списке избранного (под ОСНОВНЫМ)
	s.Run("Объявление в списке избранного", func() {
		list := favManager.GetFavoritesList(s.T(), s.tokenMain)
		require.NotEmpty(s.T(), list.Items, "Список избранного пуст")

		found := false
		for _, item := range list.Items {
			if item.Advertisement.ID == s.adID {
				found = true
				break
			}
		}
		require.True(s.T(), found, "Наше объявление не найдено в списке избранного")
	})

	// 3. Удаляем из избранного (под ОСНОВНЫМ)
	s.T().Log("Удаляем из избранного (под основным пользователем)")
	favorites.HttpRemoveFromFavorites(s.T(), s.tokenMain, s.adID)

	// 4. Проверяем, что его больше нет (под ОСНОВНЫМ)
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

func TestFavoritesSuite(t *testing.T) {
	suite.Run(t, new(TestSuite))
}
