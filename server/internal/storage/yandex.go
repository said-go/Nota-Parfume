package storage

import (
	"fmt"
	"mime/multipart"
)

type YandexStorage struct {
	Token string
}

func NewYandexStorage(token string) *YandexStorage {
	return &YandexStorage{
		Token: token,
	}
}

func (s *YandexStorage) Upload(file *multipart.FileHeader) (string, error) {
	// пока заглушка
	filename := file.Filename

	// временно просто возвращаем "фейковую ссылку"
	url := fmt.Sprintf("https://yandex.disk/%s", filename)

	return url, nil
}