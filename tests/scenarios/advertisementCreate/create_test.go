package advertisementCreate

import (
	"os"
	"testing"

	"api-tests-template/internal/client/http/advertisements"
	advManager "api-tests-template/internal/managers/advertisements"
	"api-tests-template/internal/managers/auth"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type TestSuite struct {
	suite.Suite
	token string
	adID  string
}

func (s *TestSuite) SetupSuite() {
	// Принудительно задаем URL и токен
	os.Setenv("API_URL", "https://testboard.avito.com/api/v1")

	// ВСТАВЬТЕ ВАШИ РЕАЛЬНЫЕ ДАННЫЕ ОТ TESTBOARD
	loginData := auth.Login(s.T(), "sfme2013@ya.ru", "906307Zz")
	s.token = loginData.Token

	s.T().Log("Создаем тестовое объявление")

	// Имитация картинки (заглушка)
	fakePhoto := []byte{0xFF, 0xD8, 0xFF, 0xE0, 0x00, 0x10, 0x4A, 0x46, 0x49, 0x46}

	resp := advManager.CreateAdvertisement(s.T(), s.token, "Тестовое объявление для лабы", "Описание", 1000, 1, fakePhoto, "test.jpg")
	s.adID = resp.ID
	require.NotEmpty(s.T(), s.adID)
}

func (s *TestSuite) TearDownSuite() {
	if s.adID != "" {
		advertisements.HttpDeleteAdvertisement(s.T(), s.token, s.adID)
	}
}

func (s *TestSuite) Test_GetCreatedAdvertisement() {
	ad := advManager.GetAdvertisementByID(s.T(), s.adID)

	s.Run("Название совпадает", func() {
		require.Equal(s.T(), "Тестовое объявление для лабы", ad.Title)
	})
	s.Run("Цена совпадает", func() {
		require.Equal(s.T(), 1000.0, ad.Price)
	})
}

func TestCreateSuite(t *testing.T) {
	suite.Run(t, new(TestSuite))
}
