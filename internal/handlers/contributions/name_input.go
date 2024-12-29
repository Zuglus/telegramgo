package contributions

import (
	"log"

	"telegramgo/internal/states"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func HandleNameInput(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	userID := update.Message.Chat.ID
	state := states.UserStates[userID]

	// Сохраняем имя пользователя
	state.TempMember.Name = update.Message.Text
	state.Stage = "awaiting_amount"
	states.UserStates[userID] = state

	// Запрашиваем сумму взноса
	msg := tgbotapi.NewMessage(userID, "Введите сумму взноса:")
	_, err := bot.Send(msg)
	if err != nil {
		log.Printf("Error sending message: %v", err)
	}
}