package yandex

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

const apiKey = "y0_AgAAAAAzf0BIAAybvwAAAAEUjPvsAACAPb89XvxP1KLPL7UkaNEbjfnrOA"

type YandexService struct {
}

func New() *YandexService {
	return &YandexService{}
}

// Функция для загрузки файла на Яндекс.Диск
func (s *YandexService) UploadToYandexDisk(fileName string) (string, error) {
	// Получаем URL для загрузки файла
	uploadURL, err := s.GetYandexUploadURL(fileName)
	if err != nil {
		return "", fmt.Errorf("не удалось получить URL для загрузки на Яндекс.Диск: %v", err)
	}

	// Загружаем файл
	file, err := os.Open(fileName)
	if err != nil {
		return "", fmt.Errorf("не удалось открыть файл: %v", err)
	}
	defer file.Close()

	request, err := http.NewRequest("PUT", uploadURL, file)
	if err != nil {
		return "", fmt.Errorf("не удалось создать запрос для загрузки: %v", err)
	}

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return "", fmt.Errorf("не удалось загрузить файл на Яндекс.Диск: %v", err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusCreated {
		return "", fmt.Errorf("не удалось загрузить файл, статус: %s", response.Status)
	}

	return "https://disk.yandex.ru/d/" + fileName, nil
}

// Функция для получения URL для загрузки на Яндекс.Диск
func (s *YandexService) GetYandexUploadURL(fileName string) (string, error) {
	// URL API Яндекс.Диска для загрузки файлов
	url := "https://cloud-api.yandex.net/v1/disk/resources/upload?path=" + fileName + "&overwrite=true"
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", fmt.Errorf("не удалось создать запрос для получения URL загрузки: %v", err)
	}

	// Устанавливаем заголовок с OAuth токеном
	request.Header.Set("Authorization", "OAuth "+apiKey)

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return "", fmt.Errorf("не удалось выполнить запрос: %v", err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return "", fmt.Errorf("ошибка при получении URL загрузки, статус: %s", response.Status)
	}

	var result struct {
		Href string `json:"href"`
	}
	err = json.NewDecoder(response.Body).Decode(&result)
	if err != nil {
		return "", fmt.Errorf("не удалось декодировать ответ: %v", err)
	}

	return result.Href, nil
}

// Метод для получения фото по ссылке
func (s *YandexService) DownloadPhoto(link string, savePath string) error {
	// Отправляем GET-запрос на URL файла
	response, err := http.Get(link)
	if err != nil {
		return fmt.Errorf("не удалось получить фото по ссылке: %v", err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("не удалось получить фото, статус: %s", response.Status)
	}

	// Создаем файл для сохранения
	out, err := os.Create(savePath)
	if err != nil {
		return fmt.Errorf("не удалось создать файл для сохранения: %v", err)
	}
	defer out.Close()

	// Копируем содержимое ответа в файл
	_, err = io.Copy(out, response.Body)
	if err != nil {
		return fmt.Errorf("не удалось сохранить фото: %v", err)
	}

	return nil
}
