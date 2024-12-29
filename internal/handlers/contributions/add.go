package contributions

import (
	"log"

	"telegramgo/internal/domain"
	"telegramgo/internal/states"
	"telegramgo/internal/telegram"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func HandleAddContribution(bot *tgbotapi.BotAPI, callback *tgbotapi.CallbackQuery) {
	// Начинаем процесс добавления взноса
	userID := callback.Message.Chat.ID
	// Добавляем иницилизацию в states
	states.UserStates[userID] = states.UserState{Stage: "awaiting_name", TempMember: domain.Member{Months: []string{}}}

	// Отправляем сообщение пользователю с запросом имени
	msg := tgbotapi.NewMessage(callback.Message.Chat.ID, "Введите имя:")
	_, err := bot.Send(msg)
	if err != nil {
		log.Printf("Error sending message: %v", err)
	}
	// Удаляем сообщение с inline-кнопками
	telegram.DeleteInlineKeyboard(bot, callback)
}