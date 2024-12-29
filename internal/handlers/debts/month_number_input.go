package debts

import (
	"fmt"
	"log"
	"strconv"
	"telegramgo/internal/states"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func HandleMonthNumberInput(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	userID := update.Message.Chat.ID
	state := states.UserStates[userID]

	// Пытаемся преобразовать введенный текст в число (месяц)
	selectedMonth, err := strconv.Atoi(update.Message.Text)
	if err != nil {
		// Если не удалось преобразовать в число, отправляем сообщение об ошибке
		msg := tgbotapi.NewMessage(userID, "Пожалуйста, введите корректный номер месяца (1-12).")
		_, err = bot.Send(msg)
		if err != nil {
			log.Printf("Error sending message: %v", err)
		}
		return
	}

	// Проверяем, что месяц находится в допустимом диапазоне
	if selectedMonth < 1 || selectedMonth > 12 {
		msg := tgbotapi.NewMessage(userID, "Пожалуйста, введите номер месяца в диапазоне от 1 до 12.")
		_, err = bot.Send(msg)
		if err != nil {
			log.Printf("Error sending message: %v", err)
		}
		return
	}

	// Формируем месяц в формате ГГГГ-ММ
	paymentMonth := fmt.Sprintf("%d-%02d", state.SelectedYear, selectedMonth)

	// Обновляем состояние пользователя
	state.TempMember.Months = append(state.TempMember.Months, paymentMonth)
	state.Stage = "awaiting_payment_month"
	states.UserStates[userID] = state

	// Отправляем сообщение об успешно выбранном месяце
	msg := tgbotapi.NewMessage(userID, fmt.Sprintf("Выбран месяц: %s", paymentMonth))
	_, err = bot.Send(msg)
	if err != nil {
		log.Printf("Error sending message: %v", err)
	}
}