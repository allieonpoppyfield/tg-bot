package fsm

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/allieonpoppyfield/tg-bot/internal/services/bot"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

type FillingSurveyState struct {
	state              State
	previousState      State
	availableStateList []State
}

type User struct {
	Name        string `json:"name"`
	Age         int    `json:"age"`
	Gender      int16  `json:"gender"`
	Photo       string `json:"photo_url"`
	Description string `json:"description"`
}

func (s *FillingSurveyState) Execute(update tgbotapi.Update, api *tgbotapi.BotAPI, service *bot.BotService) (State, error) {
	ctx := context.Background()
	chatId := update.Message.Chat.ID

	// Получаем текущее состояние из кеша
	cacheVal, err := service.Cache.Client.Get(ctx, fmt.Sprintf(SESSION_STATE_CACHE_KEY, chatId)).Result()
	if err != nil {
		// Если нет состояния, начинаем с имени
		return s.askName(ctx, chatId, api, service)
	}

	// Обрабатываем состояние на основе текущего значения
	return s.handleState(ctx, chatId, update, cacheVal, api, service)
}

func (s *FillingSurveyState) askName(ctx context.Context, chatId int64, api *tgbotapi.BotAPI, service *bot.BotService) (State, error) {
	// Сохраняем текущее состояние
	service.Cache.Client.Set(ctx, fmt.Sprintf(SESSION_STATE_CACHE_KEY, chatId), string(sName), 5*time.Minute)
	msg := tgbotapi.NewMessage(chatId, "Как вас зовут?")
	api.Send(msg)
	return StateFllingSurvey, nil
}

func (s *FillingSurveyState) handleState(ctx context.Context, chatId int64, update tgbotapi.Update, stateCache string, api *tgbotapi.BotAPI, service *bot.BotService) (State, error) {
	switch SessionState(stateCache) {
	case sName:
		u := User{
			Name: update.Message.Text,
		}
		service.Cache.Client.Set(ctx, fmt.Sprintf(SESSION_STATE_CACHE_KEY, chatId), string(sAge), 1*time.Minute)
		service.Cache.Client.Set(ctx, fmt.Sprintf(PROFILE_CACHE_KEY, chatId), marshalUser(u), 1*time.Minute)
		api.Send(tgbotapi.NewMessage(chatId, "Сколько вам лет?"))
	case sAge:
		var u User
		uStr := service.Cache.Client.Get(ctx, fmt.Sprintf(PROFILE_CACHE_KEY, chatId)).Val()
		json.Unmarshal([]byte(uStr), &u) // Десериализуем пользователя

		age, err := strconv.Atoi(update.Message.Text)
		if err != nil || age < 14 || age > 99 {
			api.Send(tgbotapi.NewMessage(chatId, "Введите корректный возраст"))
			return StateFllingSurvey, nil
		}
		u.Age = age
		service.Cache.Client.Set(ctx, fmt.Sprintf(SESSION_STATE_CACHE_KEY, chatId), string(sGender), 5*time.Minute)
		service.Cache.Client.Set(ctx, fmt.Sprintf(PROFILE_CACHE_KEY, chatId), marshalUser(u), 1*time.Minute)
		api.Send(tgbotapi.NewMessage(chatId, "Какого вы пола? (м / ж)"))
	case sGender:
		var u User
		uStr := service.Cache.Client.Get(ctx, fmt.Sprintf(PROFILE_CACHE_KEY, chatId)).Val()
		json.Unmarshal([]byte(uStr), &u) // Десериализуем пользователя

		if !strings.EqualFold(update.Message.Text, "м") && !strings.EqualFold(update.Message.Text, "ж") {
			api.Send(tgbotapi.NewMessage(chatId, "Нужно указать ваш пол (м / ж)"))
			return StateFllingSurvey, nil
		}
		u.Gender = 1
		if strings.EqualFold(update.Message.Text, "ж") {
			u.Gender = 2
		}
		service.Cache.Client.Set(ctx, fmt.Sprintf(SESSION_STATE_CACHE_KEY, chatId), string(sPhoto), 5*time.Minute)
		service.Cache.Client.Set(ctx, fmt.Sprintf(PROFILE_CACHE_KEY, chatId), marshalUser(u), 1*time.Minute)
		api.Send(tgbotapi.NewMessage(chatId, "Пришлите свою лучшую фотографию"))
	case sPhoto:
		log.Println("внизу мессага")
		log.Println(update.Message)
		if update.Message.Photo == nil {
			api.Send(tgbotapi.NewMessage(chatId, "Пришлите одну вашу фотографию"))
			return StateFllingSurvey, nil
		}
		log.Print("1488 я в фото 100")
		photos := *update.Message.Photo
		log.Print("1488 я в фото 102")
		largestPhoto := photos[len(photos)-1]
		log.Print("1488 я в фото 105")
		fileConfig := tgbotapi.FileConfig{
			FileID: largestPhoto.FileID,
		}

		file, err := api.GetFile(fileConfig)
		if err != nil {
			log.Panic(err)
		}
		downloadURL := file.Link(api.Token)
		log.Printf("Скачайте фото по URL: %s", downloadURL)

		// Скачиваем фото локально
		localFileName := fmt.Sprintf("%d.jpg", chatId)
		err = downloadPhoto(downloadURL, localFileName)
		if err != nil {
			log.Fatalf("Ошибка при скачивании фото: %v", err)
		}

		// Загрузка на Яндекс.Диск
		link, err := service.Ya.UploadToYandexDisk(localFileName)
		if err != nil {
			log.Fatalf("Ошибка при загрузке на Яндекс.Диск: %v", err)
		}
		err = deletePhoto(localFileName)
		if err != nil {
			return s.state, err
		}

		var u User
		uStr := service.Cache.Client.Get(ctx, fmt.Sprintf(PROFILE_CACHE_KEY, chatId)).Val()
		json.Unmarshal([]byte(uStr), &u) // Десериализуем пользователя
		u.Photo = link
		if err != nil {
			return s.state, err
		}
		service.Cache.Client.Set(ctx, fmt.Sprintf(PROFILE_CACHE_KEY, chatId), marshalUser(u), 1*time.Minute)

		service.Cache.Client.Set(ctx, fmt.Sprintf(SESSION_STATE_CACHE_KEY, chatId), string(sDescription), 5*time.Minute)
		api.Send(tgbotapi.NewMessage(chatId, "Расскажите немного о себе"))

	case sDescription:
		var u User
		uStr := service.Cache.Client.Get(ctx, fmt.Sprintf(PROFILE_CACHE_KEY, chatId)).Val()
		json.Unmarshal([]byte(uStr), &u) // Десериализуем пользователя
		u.Description = update.Message.Text
		if err := service.Repo.InsertUser(ctx, chatId, u.Name, u.Age, u.Gender, u.Description, u.Photo); err != nil {
			api.Send(tgbotapi.NewMessage(chatId, err.Error()))
			return StateMain, nil
		}
		service.Cache.Client.Del(ctx, fmt.Sprintf(SESSION_STATE_CACHE_KEY, chatId))
		service.Cache.Client.Del(ctx, fmt.Sprintf(PROFILE_CACHE_KEY, chatId))
		api.Send(tgbotapi.NewMessage(chatId, "Профиль успешно сохранен!"))
		st := MainState{
			state:         StateMain,
			previousState: StateFllingSurvey,
		}
		st.Execute(update, api, service)
	}

	return StateMain, nil
}

// Упрощаем маршалинг объекта User
func marshalUser(u User) []byte {
	um, _ := json.Marshal(u)
	return um
}

func (s *FillingSurveyState) GetState() State {
	return s.state
}

func (s *FillingSurveyState) SetPreviousState(st State) {
	s.previousState = st
}

func (s *FillingSurveyState) GetPreviousState() State {
	return s.previousState
}

func (s *FillingSurveyState) GetAvailableStateList() []State {
	return s.availableStateList
}

// Функция для скачивания фото локально
func downloadPhoto(url string, fileName string) error {
	response, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("не удалось скачать фото: %v", err)
	}
	defer response.Body.Close()

	// Создаем файл на локальной системе
	file, err := os.Create(fileName)
	if err != nil {
		return fmt.Errorf("не удалось создать файл: %v", err)
	}
	defer file.Close()

	_, err = io.Copy(file, response.Body)
	if err != nil {
		return fmt.Errorf("не удалось сохранить фото: %v", err)
	}

	return nil
}

// Функция для скачивания фото локально
func deletePhoto(fileName string) error {
	err := os.Remove(fileName)
	if err != nil {
		return fmt.Errorf("не удалось удалить файл: %v", err)
	}
	return nil
}
