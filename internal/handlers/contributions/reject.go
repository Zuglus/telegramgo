package contributions

import (
	"log"

	"telegramgo/internal/states"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func HandleRejectContribution(bot *tgbotapi.BotAPI, callback *tgbotapi.CallbackQuery) {
	userID := callback.Message.Chat.ID

	// Отвечаем на callback, чтобы убрать "часики" на кнопке
	callbackConfig := tgbotapi.NewCallback(callback.ID, "")
	if _, err := bot.Request(callbackConfig); err != nil {
		log.Printf("Error responding to callback: %v", err)
	}

	// Очищаем состояние пользователя
	delete(states.UserStates, userID)

	// Отправляем сообщение об отмене
	msg := tgbotapi.NewMessage(userID, "Ввод взноса отменен.")
	_, err := bot.Send(msg)
	if err != nil {
		log.Printf("Error sending message: %v", err)
	}
}