package advertisements

import (
	"bytes"
	"mime/multipart"
	"net/http"
	"strconv"
	"testing"

	"api-tests-template/internal/constants/path"
	apiRunner "api-tests-template/internal/helpers/api-runner"
)

// HttpCreateAdvertisement создает объявление с фотографиями
func HttpCreateAdvertisement(t *testing.T, token, title, description string, price float64, quantity int, photoBytes []byte, photoName string) *http.Response {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	_ = writer.WriteField("title", title)
	_ = writer.WriteField("description", description)
	if price > 0 {
		_ = writer.WriteField("price", strconv.FormatFloat(price, 'f', -1, 64))
	}
	if quantity > 0 {
		_ = writer.WriteField("quantity", strconv.Itoa(quantity))
	}

	if len(photoBytes) > 0 {
		part, _ := writer.CreateFormFile("photos", photoName)
		_, _ = part.Write(photoBytes)
	}
	writer.Close()

	return apiRunner.GetRunner().Auth(token).Create().
		Post(path.AdvertisementPath).
		Header("Content-Type", writer.FormDataContentType()).
		Body(body.String()).
		Expect(t).
		End().Response
}

// HttpGetAdvertisementByID получает объявление по его ID
func HttpGetAdvertisementByID(t *testing.T, id string) *http.Response {
	return apiRunner.GetRunner().Create().
		Get(path.AdvertisementPath).
		Query("id", id).
		Expect(t).
		Status(http.StatusOK).
		End().Response
}

// HttpDeleteAdvertisement удаляет объявление по его ID
func HttpDeleteAdvertisement(t *testing.T, token, id string) *http.Response {
	return apiRunner.GetRunner().Auth(token).Create().
		Delete(path.AdvertisementPath).
		Query("id", id).
		Expect(t).
		End().Response
}
