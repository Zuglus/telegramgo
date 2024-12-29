package debts

import (
	"fmt"
	"log"
	"strconv"

	"telegramgo/internal/states"
	"telegramgo/internal/telegram"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func HandleYearInput(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	userID := update.Message.Chat.ID
	state := states.UserStates[userID]

	// Пытаемся преобразовать введенный текст в число (год)
	selectedYear, err := strconv.Atoi(update.Message.Text)
	if err != nil {
		// Если не удалось преобразовать в число, отправляем сообщение об ошибке
		msg := tgbotapi.NewMessage(userID, "Пожалуйста, введите корректный год в формате ГГГГ.")
		_, err = bot.Send(msg)
		if err != nil {
			log.Printf("Error sending message: %v", err)
		}
		return
	}

	// Обновляем состояние пользователя с выбранным годом
	state.SelectedYear = selectedYear
	state.Stage = "awaiting_month_number"
	states.UserStates[userID] = state

	// Запрашиваем выбор месяца
	monthsKeyboard := telegram.GenerateMonthsKeyboard()

	msg := tgbotapi.NewMessage(userID, fmt.Sprintf("Выбран год: %d\nВыберите месяц:", selectedYear))
	msg.ReplyMarkup = monthsKeyboard
	_, err = bot.Send(msg)
	if err != nil {
		log.Printf("Error sending message: %v", err)
	}
}