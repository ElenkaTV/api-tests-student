package advertisementSearch

import (
	"os"
	"strings"
	"testing"
	"unicode"

	"api-tests-template/internal/assertions"
	advList "api-tests-template/internal/managers/advertisements"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

const (
	existingAdTitle = "Тестовый ждун"
	searchBeginning = "Тестовый"
	searchMiddle    = "вый жд"
	searchEnd       = "ждун"
)

type TestSuite struct {
	suite.Suite
}

func (s *TestSuite) SetupSuite() {
	os.Setenv("API_URL", "https://testboard.avito.com/api/v1")

	require.NotEmpty(s.T(), existingAdTitle, "задайте existingAdTitle")
	require.NotEmpty(s.T(), searchBeginning, "задайте searchBeginning")
	require.NotEmpty(s.T(), searchMiddle, "задайте searchMiddle")
	require.NotEmpty(s.T(), searchEnd, "задайте searchEnd")
}

func (s *TestSuite) TestSearchByTitle() {
	token := "fake-token-for-check"

	list := advList.GetAdvertisementsList(s.T(), token, searchMiddle)

	s.Run("В ответе есть хотя бы одно объявление", func() {
		require.NotEmpty(s.T(), list.Items)
	})

	s.Run("Среди найденных есть наше объявление", func() {
		assertions.ContainsAdvertisementByTitle(s.T(), list.Items, existingAdTitle)
	})
}

func (s *TestSuite) TestSearchByTitle_Substring() {
	cases := []struct {
		name   string
		search string
	}{
		{"начало", searchBeginning},
		{"середина", searchMiddle},
		{"конец", searchEnd},
		{"нижний регистр", strings.ToLower(searchMiddle)},
		{"верхний регистр", strings.ToUpper(searchMiddle)},
		{"смешанный регистр", mixCase(searchMiddle)},
	}

	for _, tc := range cases {
		s.Run(tc.name, func() {
			list := advList.GetAdvertisementsList(s.T(), "fake-token-for-check", tc.search)
			assertions.ContainsAdvertisementByTitle(s.T(), list.Items, existingAdTitle)
		})
	}
}

//  для проверки смешанного регистра
func mixCase(s string) string {
	runes := []rune(s)
	for i := 0; i < len(runes); i++ {
		if i%2 == 0 {
			runes[i] = unicode.ToUpper(runes[i])
		} else {
			runes[i] = unicode.ToLower(runes[i])
		}
	}
	return string(runes)
}

func TestSearchSuite(t *testing.T) {
	suite.Run(t, new(TestSuite))
}
