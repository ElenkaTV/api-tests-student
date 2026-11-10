package assertions

import (
	"testing"

	"api-tests-template/internal/managers/advertisements/models"
	"github.com/stretchr/testify/assert"
)

func ContainsAdvertisementByTitle(t *testing.T, items []models.Advertisement, title string) {
	for _, item := range items {
		if item.Title == title {
			return
		}
	}
	assert.Fail(t, "Объявление с заголовком '%s' не найдено в списке", title)
}