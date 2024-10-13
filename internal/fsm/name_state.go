package fsm

// import (
// 	"github.com/allieonpoppyfield/tg-bot/internal/services/bot"
// 	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
// )

// type NameState struct{}

// func (s *NameState) Execute(update tgbotapi.Update, bot *bot.BotService) (State, error) {
// 	chatID := update.Message.Chat.ID
// 	user := userStates[chatID]

// 	user.Name = update.Message.Text
// 	bot.SendMessage(chatID, "Введите свой возраст:")
// 	user.State = StateAge // Переход на следующий шаг
// 	return StateAge, nil
// }
