package advertisementCreateNegative

import (
	"net/http"
	"os"
	"testing"

	"api-tests-template/internal/client/http/advertisements"
	"api-tests-template/internal/managers/auth"
	"api-tests-template/internal/utils"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type TestSuite struct {
	suite.Suite
	token string
}

func (s *TestSuite) SetupSuite() {
	utils.LoadEnv()
	os.Setenv("API_URL", "https://testboard.avito.com/api/v1")

	login := os.Getenv("TEST_LOGIN")
	pass := os.Getenv("TEST_PASSWORD")
	require.NotEmpty(s.T(), login, "TEST_LOGIN не задан в .env.override")

	s.T().Log("Логинимся под основным пользователем")
	authData := auth.Login(s.T(), login, pass)
	s.token = authData.Token
	require.NotEmpty(s.T(), s.token, "Не удалось получить токен")
}

// 1. Создание без обязательного поля title (ожидаем 400)
func (s *TestSuite) Test_CreateWithoutTitle() {
	fakePhoto := []byte{0xFF, 0xD8, 0xFF, 0xE0, 0x00, 0x10, 0x4A, 0x46, 0x49, 0x46}
	raw := advertisements.HttpCreateAdvertisement(s.T(), s.token, "", "Описание", 1000, 1, fakePhoto, "test.jpg")
	defer raw.Body.Close()

	require.Equal(s.T(), http.StatusBadRequest, raw.StatusCode, "Ожидался статус 400 при пустом title")
}

// 2. Создание с пустым описанием (ожидаем 400, если description обязателен)
func (s *TestSuite) Test_CreateWithEmptyDescription() {
	fakePhoto := []byte{0xFF, 0xD8, 0xFF, 0xE0, 0x00, 0x10, 0x4A, 0x46, 0x49, 0x46}
	raw := advertisements.HttpCreateAdvertisement(s.T(), s.token, "Тест", "", 1000, 1, fakePhoto, "test.jpg")
	defer raw.Body.Close()

	require.Equal(s.T(), http.StatusBadRequest, raw.StatusCode, "Ожидался статус 400 при пустом description")
}

// 3. Создание без авторизации (без токена, ожидаем 401)
func (s *TestSuite) Test_CreateWithoutAuth() {
	fakePhoto := []byte{0xFF, 0xD8, 0xFF, 0xE0, 0x00, 0x10, 0x4A, 0x46, 0x49, 0x46}
	raw := advertisements.HttpCreateAdvertisement(s.T(), "", "Тест", "Описание", 1000, 1, fakePhoto, "test.jpg")
	defer raw.Body.Close()

	require.Equal(s.T(), http.StatusUnauthorized, raw.StatusCode, "Ожидался статус 401 без авторизации")
}

func TestCreateNegativeSuite(t *testing.T) {
	suite.Run(t, new(TestSuite))
}
