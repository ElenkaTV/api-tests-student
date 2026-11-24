package advertisementCreateNegative

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
	"strings"
	"testing"

	"api-tests-template/internal/client/http/advertisements"
	"api-tests-template/internal/managers/advertisements/models"
	"api-tests-template/internal/managers/auth"
	"api-tests-template/internal/utils"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type TestSuite struct {
	suite.Suite
	token string
}

// SetupSuite подготавливает окружение: логин пользователя
func (s *TestSuite) SetupSuite() {
	utils.LoadEnv()
	os.Setenv("API_URL", "https://testboard.avito.com/api/v1")
	login := os.Getenv("TEST_LOGIN")
	pass := os.Getenv("TEST_PASSWORD")
	require.NotEmpty(s.T(), login, "TEST_LOGIN не задан в .env.override")

	authData := auth.Login(s.T(), login, pass)
	s.token = authData.Token
}

// Test_CreateWithoutTitle проверяет ошибку 400 при создании объявления без title
func (s *TestSuite) Test_CreateWithoutTitle() {
	fakePhoto := []byte{0xFF, 0xD8}
	raw := advertisements.HttpCreateAdvertisement(s.T(), s.token, "", "Описание", 1000, 1, fakePhoto, "test.jpg")
	defer raw.Body.Close()

	require.Equal(s.T(), http.StatusBadRequest, raw.StatusCode)
	body, _ := io.ReadAll(raw.Body)
	var errResp models.ErrorResponse
	_ = json.Unmarshal(body, &errResp)
	require.Equal(s.T(), "invalid_request", errResp.Error)
	require.Equal(s.T(), "title is required", errResp.Message)
}

// Test_CreateWithEmptyDescription проверяет ошибку 400 при создании объявления с пустым description
func (s *TestSuite) Test_CreateWithEmptyDescription() {
	fakePhoto := []byte{0xFF, 0xD8}
	raw := advertisements.HttpCreateAdvertisement(s.T(), s.token, "Тест", "", 1000, 1, fakePhoto, "test.jpg")
	defer raw.Body.Close()

	require.Equal(s.T(), http.StatusBadRequest, raw.StatusCode)
	body, _ := io.ReadAll(raw.Body)
	var errResp models.ErrorResponse
	_ = json.Unmarshal(body, &errResp)
	require.Equal(s.T(), "invalid_request", errResp.Error)
}

// Test_CreateWithLongTitle проверяет ошибку 400 при слишком длинном title
func (s *TestSuite) Test_CreateWithLongTitle() {
	longTitle := strings.Repeat("A", 300)
	fakePhoto := []byte{0xFF, 0xD8}
	raw := advertisements.HttpCreateAdvertisement(s.T(), s.token, longTitle, "Описание", 1000, 1, fakePhoto, "test.jpg")
	defer raw.Body.Close()

	require.Equal(s.T(), http.StatusBadRequest, raw.StatusCode)
	body, _ := io.ReadAll(raw.Body)
	var errResp models.ErrorResponse
	_ = json.Unmarshal(body, &errResp)
	require.Contains(s.T(), errResp.Message, "title")
}

// Test_CreateWithoutAuth проверяет ошибку 401 при создании без авторизации
func (s *TestSuite) Test_CreateWithoutAuth() {
	fakePhoto := []byte{0xFF, 0xD8}
	raw := advertisements.HttpCreateAdvertisement(s.T(), "", "Тест", "Описание", 1000, 1, fakePhoto, "test.jpg")
	defer raw.Body.Close()

	require.Equal(s.T(), http.StatusUnauthorized, raw.StatusCode)
	body, _ := io.ReadAll(raw.Body)
	var errResp models.ErrorResponse
	_ = json.Unmarshal(body, &errResp)
	require.Equal(s.T(), "unauthorized", errResp.Error)
}

// TestCreateNegativeSuite запускает набор негативных тестов
func TestCreateNegativeSuite(t *testing.T) {
	suite.Run(t, new(TestSuite))
}
