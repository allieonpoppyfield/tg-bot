package fsm

// import (
// 	"github.com/allieonpoppyfield/tg-bot/internal/services/bot"
// 	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
// )

// type MainState struct{}

// func (s *MainState) Execute(update tgbotapi.Update, bot *bot.BotService) (State, error) {
// 	chatID := update.Message.Chat.ID
// 	// Создаем кнопки
// 	keyboard := tgbotapi.NewReplyKeyboard(
// 		tgbotapi.NewKeyboardButtonRow(
// 			tgbotapi.NewKeyboardButton("Оставить заявку"),
// 			tgbotapi.NewKeyboardButton("Поиск заявок"),
// 		),
// 		tgbotapi.NewKeyboardButtonRow(
// 			tgbotapi.NewKeyboardButton("Мой профиль"),
// 			tgbotapi.NewKeyboardButton("Мои заявки"),
// 		),
// 	)

// 	// Отправляем сообщение с клавиатурой
// 	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Выберите действие:")
// 	msg.ReplyMarkup = keyboard

// 	bot.SendMessage(msg)
// 	return StateAge, nil
// }
