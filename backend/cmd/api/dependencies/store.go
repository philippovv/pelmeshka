package dependencies

import (
	"gitlab.praktikum-services.ru/Stasyan/momo-store/internal/store/dumplings"
	"gitlab.praktikum-services.ru/Stasyan/momo-store/internal/store/dumplings/fake"
)

// NewFakeDumplingsStore returns new fake store for app
func NewFakeDumplingsStore() (dumplings.Store, error) {
	packs := []dumplings.Product{
		{
			ID:          1,
			Name:        "Пельмени",
			Description: "ВКУСНЕНЬКИЕ",
			Price:       5.00,
			Image:       "https://storage.yandexcloud.net/picture-for-site/%D0%BF%D0%B5%D0%BB%D1%8C%D0%BC%D0%B5%D1%88%D0%BA%D0%B8.jpg",
		},
		{
			ID:          2,
			Name:        "Хинкали",
			Description: "С чернилами девопса",
			Price:       3.50,
			Image:       "https://storage.yandexcloud.net/picture-for-site/%D1%85%D0%B8%D0%BD%D0%BA%D0%B0%D0%BB%D0%B8.jpg",
		},
		{
			ID:          3,
			Name:        "Манты",
			Description: "С мясом молодых разработчиков",
			Price:       2.75,
			Image:       "https://storage.yandexcloud.net/picture-for-site/%D0%BC%D0%B0%D0%BD%D1%82%D1%8B.jpg",
		},
		{
			ID:          4,
			Name:        "Буузы",
			Description: "с любовью бабушки",
			Price:       4.00,
			Image:       "https://storage.yandexcloud.net/picture-for-site/%D0%91%D1%83%D1%83%D0%B7%D1%8B.jpg",
		},
		{
			ID:          5,
			Name:        "Цзяоцзы",
			Description: "С говядиной и свининой",
			Price:       7.25,
			Image:       "https://storage.yandexcloud.net/picture-for-site/%D0%A6%D0%B7%D1%8F%D0%BE%D1%86%D0%B7%D1%8B.jpg",
		},
		{
			ID:          6,
			Name:        "Гедза",
			Description: "С соевым мясом",
			Price:       3.50,
			Image:       "https://storage.yandexcloud.net/picture-for-site/%D0%93%D0%B5%D0%B4%D0%B7%D0%B0.jpg",
		},
		{
			ID:          7,
			Name:        "Дим-самы",
			Description: "С какой-то девушкой",
			Price:       2.65,
			Image:       "https://storage.yandexcloud.net/picture-for-site/%D0%94%D0%B8%D0%BC-%D1%81%D0%B0%D0%BC%D1%8B.webp",
		},
		{
			ID:          8,
			Name:        "Момо",
			Description: "в подарок Аватар",
			Price:       5.00,
			Image:       "https://storage.yandexcloud.net/picture-for-site/%D0%9C%D0%BE%D0%BC%D0%BE.jpg",
		},
		{
			ID:          9,
			Name:        "Вонтоны",
			Description: "С крепким алкогольным",
			Price:       4.10,
			Image:       "https://storage.yandexcloud.net/picture-for-site/%D0%92%D0%BE%D0%BD%D1%82%D0%BE%D0%BD%D1%8B.jpg",
		},
		{
			ID:          10,
			Name:        "Баоцзы",
			Description: "Оптом",
			Price:       4.20,
			Image:       "https://storage.yandexcloud.net/picture-for-site/%D0%91%D0%B0%D0%BE%D1%86%D0%B7%D1%8B.jpg",
		},
		{
			ID:          11,
			Name:        "Кундюмы",
			Description: "По старорусски",
			Price:       5.45,
			Image:       "https://storage.yandexcloud.net/picture-for-site/%D0%9A%D1%83%D0%BD%D0%B4%D1%8E%D0%BC%D1%8B.jpg",
		},
		{
			ID:          12,
			Name:        "Курзе",
			Description: "Цветные",
			Price:       3.25,
			Image:       "https://storage.yandexcloud.net/picture-for-site/%D0%9A%D1%83%D1%80%D0%B7%D0%B5.jpg",
		},
		{
			ID:          13,
			Name:        "Бораки",
			Description: "От души!",
			Price:       4.00,
			Image:       "https://storage.yandexcloud.net/picture-for-site/%D0%91%D0%BE%D1%80%D0%B0%D0%BA%D0%B8.jpg",
		},
		{
			ID:          14,
			Name:        "Равиоли",
			Description: "Я из Италии прибыль",
			Price:       2.90,
			Image:       "https://storage.yandexcloud.net/picture-for-site/%D0%A0%D0%B0%D0%B2%D0%B8%D0%BE%D0%BB%D0%B8.jpg",
		},
	}

	store := fake.NewStore()
	store.SetAvailablePacks(packs...)

	return store, nil
}
