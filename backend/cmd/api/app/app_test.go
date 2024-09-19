package app

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"

	"gitlab.praktikum-services.ru/Stasyan/momo-store/cmd/api/dependencies"
)

func TestFakeAppIntegrational(t *testing.T) {
	store, err := dependencies.NewFakeDumplingsStore()
	assert.NoError(t, err)
	app, err := NewInstance(store)
	assert.NoError(t, err)

	t.Run("create_order", func(t *testing.T) {
		for i := 1; i <= 10; i++ {
			t.Run("id"+strconv.Itoa(i), func(t *testing.T) {
				r := httptest.NewRequest("POST", "/orders", nil)
				w := httptest.NewRecorder()
				app.CreateOrderController(w, r)

				assert.Equal(t, http.StatusOK, w.Code)
				assert.Equal(t, "application/json", w.Header().Get("Content-Type"))
				fmt.Fprintln(os.Stdout, "_____")
				fmt.Fprintln(os.Stdout, w.Body.String())
				fmt.Fprintln(os.Stdout, "_____")

				expectedJSON, err := json.Marshal(map[string]interface{}{"id": i})
				assert.NoError(t, err)
				assert.JSONEq(t, string(expectedJSON), w.Body.String())
			})
		}
	})

	t.Run("list_dumplings", func(t *testing.T) {
		r := httptest.NewRequest("GET", "/packs", nil)
		w := httptest.NewRecorder()
		app.ListDumplingsController(w, r)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, "application/json", w.Header().Get("Content-Type"))

		fmt.Fprintln(os.Stdout, "_____")
		fmt.Fprintln(os.Stdout, w.Body.String())
		fmt.Fprintln(os.Stdout, "_____")

		expectedJSON := "{\"results\":[{\"id\":1,\"name\":\"Пельмени\",\"price\":5.00,\"description\":\"ВКУСНЕНЬКИЕ\",\"image\":\"https://storage.yandexcloud.net/picture-for-site/%D0%BF%D0%B5%D0%BB%D1%8C%D0%BC%D0%B5%D1%88%D0%BA%D0%B8.jpg\"},{\"id\":2,\"name\":\"Хинкали\",\"price\":3.50,\"description\":\"С чернилами девопса\",\"image\":\"https://storage.yandexcloud.net/picture-for-site/%D1%85%D0%B8%D0%BD%D0%BA%D0%B0%D0%BB%D0%B8.jpg\"},{\"id\":3,\"name\":\"Манты\",\"price\":2.75,\"description\":\"С мясом молодых разработчиков\",\"image\":\"https://storage.yandexcloud.net/picture-for-site/%D0%BC%D0%B0%D0%BD%D1%82%D1%8B.jpg\"},{\"id\":4,\"name\":\"Буузы\",\"price\":4.00,\"description\":\"с любовью бабушки\",\"image\":\"https://storage.yandexcloud.net/picture-for-site/%D0%91%D1%83%D1%83%D0%B7%D1%8B.jpg\"},{\"id\":5,\"name\":\"Цзяоцзы\",\"price\":7.25,\"description\":\"С говядиной и свининой\",\"image\":\"https://storage.yandexcloud.net/picture-for-site/%D0%A6%D0%B7%D1%8F%D0%BE%D1%86%D0%B7%D1%8B.jpg\"},{\"id\":6,\"name\":\"Гедза\",\"price\":3.50,\"description\":\"С соевым мясом\",\"image\":\"https://storage.yandexcloud.net/picture-for-site/%D0%93%D0%B5%D0%B4%D0%B7%D0%B0.jpg\"},{\"id\":7,\"name\":\"Дим-самы\",\"price\":2.65,\"description\":\"С какой-то девушкой\",\"image\":\"https://storage.yandexcloud.net/picture-for-site/%D0%94%D0%B8%D0%BC-%D1%81%D0%B0%D0%BC%D1%8B.webp\"},{\"id\":8,\"name\":\"Момо\",\"price\":5.00,\"description\":\"в подарок Аватар\",\"image\":\"https://storage.yandexcloud.net/picture-for-site/%D0%9C%D0%BE%D0%BC%D0%BE.jpg\"},{\"id\":9,\"name\":\"Вонтоны\",\"price\":4.10,\"description\":\"С крепким алкогольным\",\"image\":\"https://storage.yandexcloud.net/picture-for-site/%D0%92%D0%BE%D0%BD%D1%82%D0%BE%D0%BD%D1%8B.jpg\"},{\"id\":10,\"name\":\"Баоцзы\",\"price\":4.20,\"description\":\"Оптом\",\"image\":\"https://storage.yandexcloud.net/picture-for-site/%D0%91%D0%B0%D0%BE%D1%86%D0%B7%D1%8B.jpg\"},{\"id\":11,\"name\":\"Кундюмы\",\"price\":5.45,\"description\":\"По старорусски\",\"image\":\"https://storage.yandexcloud.net/picture-for-site/%D0%9A%D1%83%D0%BD%D0%B4%D1%8E%D0%BC%D1%8B.jpg\"},{\"id\":12,\"name\":\"Курзе\",\"price\":3.25,\"description\":\"Цветные\",\"image\":\"https://storage.yandexcloud.net/picture-for-site/%D0%9A%D1%83%D1%80%D0%B7%D0%B5.jpg\"},{\"id\":13,\"name\":\"Бораки\",\"price\":4.00,\"description\":\"От души!\",\"image\":\"https://storage.yandexcloud.net/picture-for-site/%D0%91%D0%BE%D1%80%D0%B0%D0%BA%D0%B8.jpg\"},{\"id\":14,\"name\":\"Равиоли\",\"price\":2.90,\"description\":\"Я из Италии прибыль\",\"image\":\"https://storage.yandexcloud.net/picture-for-site/%D0%A0%D0%B0%D0%B2%D0%B8%D0%BE%D0%BB%D0%B8.jpg\"}]}\n"

		assert.NoError(t, err)
		assert.JSONEq(t, string(expectedJSON), w.Body.String())
	})

	t.Run("healthcheck", func(t *testing.T) {
		r := httptest.NewRequest("GET", "/health", nil)
		w := httptest.NewRecorder()
		app.HealthcheckController(w, r)

		assert.Equal(t, http.StatusOK, w.Code)
	})
}
