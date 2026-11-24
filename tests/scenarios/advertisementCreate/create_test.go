package advertisementCreate

import (
	"os"
	"testing"

	"api-tests-template/internal/client/http/advertisements"
	advManager "api-tests-template/internal/managers/advertisements"
	"api-tests-template/internal/managers/auth"
	"api-tests-template/internal/utils"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type TestSuite struct {
	suite.Suite
	token string
	adID  string
}

// SetupSuite подготавливает окружение: логин и создание тестового объявления
func (s *TestSuite) SetupSuite() {
	utils.LoadEnv()
	os.Setenv("API_URL", "https://testboard.avito.com/api/v1")

	login := os.Getenv("TEST_LOGIN")
	pass := os.Getenv("TEST_PASSWORD")
	require.NotEmpty(s.T(), login, "TEST_LOGIN не задан в .env.override")

	authData := auth.Login(s.T(), login, pass)
	s.token = authData.Token

	s.T().Log("Создаем тестовое объявление")
	fakePhoto := []byte{0xFF, 0xD8, 0xFF, 0xE0, 0x00, 0x10, 0x4A, 0x46, 0x49, 0x46}
	resp := advManager.CreateAdvertisement(s.T(), s.token, "Тестовое объявление для лабы", "Полное описание", 1000, 2, fakePhoto, "test.jpg")
	s.adID = resp.ID
	require.NotEmpty(s.T(), s.adID)
}

// TearDownSuite удаляет тестовое объявление после завершения тестов
func (s *TestSuite) TearDownSuite() {
	if s.adID != "" {
		s.T().Log("Удаляем тестовое объявление")
		advertisements.HttpDeleteAdvertisement(s.T(), s.token, s.adID)
	}
}

// Test_GetCreatedAdvertisement проверяет, что созданное объявление доступно по ID и поля совпадают
func (s *TestSuite) Test_GetCreatedAdvertisement() {
	s.T().Log("Получаем объявление по ID")
	ad := advManager.GetAdvertisementByID(s.T(), s.adID)

	s.Run("Название совпадает", func() {
		require.Equal(s.T(), "Тестовое объявление для лабы", ad.Title)
	})
	s.Run("Описание совпадает", func() {
		require.Equal(s.T(), "Полное описание", ad.Description)
	})
	s.Run("Цена совпадает", func() {
		require.Equal(s.T(), 1000.0, ad.Price)
	})
	s.Run("Количество совпадает", func() {
		require.Equal(s.T(), 2, ad.Quantity)
	})
	s.Run("Количество фото равно 1", func() {
		require.Len(s.T(), ad.Photos, 1)
	})
}

// Test_SearchCreatedAdvertisement проверяет, что созданное объявление находится в поиске
func (s *TestSuite) Test_SearchCreatedAdvertisement() {
	s.T().Log("Ищем созданное объявление")
	list := advManager.GetAdvertisementsList(s.T(), s.token, "Тестовое объявление для лабы")

	s.Run("Объявление найдено в поиске", func() {
		found := false
		for _, item := range list.Items {
			if item.ID == s.adID {
				found = true
				break
			}
		}
		require.True(s.T(), found, "Объявление с ID %s не найдено в поиске", s.adID)
	})
}

// TestCreateSuite запускает тестовый набор
func TestCreateSuite(t *testing.T) {
	suite.Run(t, new(TestSuite))
}
