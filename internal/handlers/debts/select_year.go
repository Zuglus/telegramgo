package debts

import (
	"log"

	"telegramgo/internal/states"
	"telegramgo/internal/telegram"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func HandleSelectYear(bot *tgbotapi.BotAPI, callback *tgbotapi.CallbackQuery, state states.UserState) {
	userID := callback.Message.Chat.ID

	telegram.AnswerCallbackQuery(bot, callback)

	state.Stage = "awaiting_year"
	states.UserStates[userID] = state

	// Запрашиваем ввод года
	msg := tgbotapi.NewMessage(userID, "Введите год в формате ГГГГ:")
	_, err := bot.Send(msg)
	if err != nil {
		log.Printf("Error sending message: %v", err)
	}
	// Удаляем сообщение с inline-кнопками
	telegram.DeleteInlineKeyboard(bot, callback)
}