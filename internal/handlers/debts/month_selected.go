package debts

import (
	"fmt"
	"log"

	"telegramgo/internal/states"
	"telegramgo/internal/telegram"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func HandleMonthSelected(bot *tgbotapi.BotAPI, callback *tgbotapi.CallbackQuery, state states.UserState, selectedMonth int) {
	userID := callback.Message.Chat.ID

	// Отвечаем на callback, чтобы убрать "часики" на кнопке
	telegram.AnswerCallbackQuery(bot, callback)

	// Формируем месяц в формате ГГГГ-ММ
	paymentMonth := fmt.Sprintf("%d-%02d", state.SelectedYear, selectedMonth)

	// Обновляем состояние пользователя
	state.TempMember.Months = append(state.TempMember.Months, paymentMonth)
	state.Stage = "awaiting_payment_month"
	states.UserStates[userID] = state

	// Отправляем сообщение об успешно выбранном месяце
	msg := tgbotapi.NewMessage(userID, fmt.Sprintf("Выбран месяц: %s", paymentMonth))
	_, err := bot.Send(msg)
	if err != nil {
		log.Printf("Error sending message: %v", err)
	}
	// Удаляем сообщение с inline-кнопками
	telegram.DeleteInlineKeyboard(bot, callback)
}