package debts

import (
	"fmt"
	"log"

	"telegramgo/internal/states"
	"telegramgo/internal/telegram"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func HandleYearSelected(bot *tgbotapi.BotAPI, callback *tgbotapi.CallbackQuery, state states.UserState, selectedYear int) {
	userID := callback.Message.Chat.ID

	// Обновляем состояние пользователя с выбранным годом
	state.SelectedYear = selectedYear
	state.Stage = "awaiting_month_number"
	states.UserStates[userID] = state

	// Запрашиваем выбор месяца
	monthsKeyboard := telegram.GenerateMonthsKeyboard()

	msg := tgbotapi.NewMessage(userID, fmt.Sprintf("Выбран год: %d\nВыберите месяц:", selectedYear))
	msg.ReplyMarkup = monthsKeyboard
	_, err := bot.Send(msg)
	if err != nil {
		log.Printf("Error sending message: %v", err)
	}
	// Удаляем сообщение с inline-кнопками
	telegram.DeleteInlineKeyboard(bot, callback)
}