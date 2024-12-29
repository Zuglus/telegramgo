package telegram

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"telegramgo/internal/states"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// AnswerCallbackQuery отвечает на callback query, убирая "часики" на кнопке.
func AnswerCallbackQuery(bot *tgbotapi.BotAPI, callbackQuery *tgbotapi.CallbackQuery) {
	callbackConfig := tgbotapi.NewCallback(callbackQuery.ID, "")
	if _, err := bot.Request(callbackConfig); err != nil {
		log.Printf("Error responding to callback: %v", err)
	}
}

// DeleteInlineKeyboard удаляет inline-клавиатуру из сообщения.
func DeleteInlineKeyboard(bot *tgbotapi.BotAPI, callbackQuery *tgbotapi.CallbackQuery) {
	deleteMsg := tgbotapi.NewDeleteMessage(callbackQuery.Message.Chat.ID, callbackQuery.Message.MessageID)
	_, err := bot.Send(deleteMsg)
	if err != nil {
		log.Printf("Failed to delete message: %v", err)
	}
}

// HandleSelectMonth обрабатывает выбор месяца из inline-клавиатуры.
func HandleSelectMonth(bot *tgbotapi.BotAPI, callback *tgbotapi.CallbackQuery, state states.UserState) {
	userID := callback.Message.Chat.ID

	// answerCallbackQuery отвечает на callback query, убирая "часики" на кнопке.
	AnswerCallbackQuery(bot, callback)

	now := time.Now()
	var paymentMonth string

	switch callback.Data {
	case "current_month":
		paymentMonth = now.Format("2006-01")
	case "previous_month":
		paymentMonth = now.AddDate(0, -1, 0).Format("2006-01")
	}

	state.TempMember.Months = append(state.TempMember.Months, paymentMonth)
	state.Stage = "awaiting_payment_month"
	states.UserStates[userID] = state

	// Создаем клавиатуру с кнопками "Подтвердить" и "Отклонить"
	confirmationKeyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Подтвердить", "confirm_contribution"),
			tgbotapi.NewInlineKeyboardButtonData("Отклонить", "reject_contribution"),
		),
	)

	msg := tgbotapi.NewMessage(userID, fmt.Sprintf("Выбран месяц: %s\nПодтвердите добавление взноса.", paymentMonth))
	msg.ReplyMarkup = confirmationKeyboard
	_, err := bot.Send(msg)
	if err != nil {
		log.Printf("Error sending message: %v", err)
	}
	// Удаляем сообщение с inline-кнопками
	DeleteInlineKeyboard(bot, callback)
}

func HandleDebt(bot *tgbotapi.BotAPI, callback *tgbotapi.CallbackQuery) {
	userID := callback.Message.Chat.ID

	// Запрашиваем выбор года
	currentYear := time.Now().Year()
	yearsKeyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(strconv.Itoa(currentYear), "year_"+strconv.Itoa(currentYear)),
			tgbotapi.NewInlineKeyboardButtonData(strconv.Itoa(currentYear-1), "year_"+strconv.Itoa(currentYear-1)),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Другой год", "select_year"),
		),
	)

	msg := tgbotapi.NewMessage(userID, "Выберите год:")
	msg.ReplyMarkup = yearsKeyboard
	_, err := bot.Send(msg)
	if err != nil {
		log.Printf("Error sending message: %v", err)
	}
	// Удаляем сообщение с inline-кнопками
	DeleteInlineKeyboard(bot, callback)
}

func HandleSelectYear(bot *tgbotapi.BotAPI, callback *tgbotapi.CallbackQuery, state states.UserState) {
	userID := callback.Message.Chat.ID

	state.Stage = "awaiting_year"
	states.UserStates[userID] = state

	// Запрашиваем ввод года
	msg := tgbotapi.NewMessage(userID, "Введите год в формате ГГГГ:")
	_, err := bot.Send(msg)
	if err != nil {
		log.Printf("Error sending message: %v", err)
	}
	// Удаляем сообщение с inline-кнопками
	DeleteInlineKeyboard(bot, callback)
}

func HandleYearSelected(bot *tgbotapi.BotAPI, callback *tgbotapi.CallbackQuery, state states.UserState, selectedYear int) {
	userID := callback.Message.Chat.ID

	// Обновляем состояние пользователя с выбранным годом
	state.SelectedYear = selectedYear
	state.Stage = "awaiting_month_number"
	states.UserStates[userID] = state

	// Запрашиваем выбор месяца
	monthsKeyboard := GenerateMonthsKeyboard()

	msg := tgbotapi.NewMessage(userID, fmt.Sprintf("Выбран год: %d\nВыберите месяц:", selectedYear))
	msg.ReplyMarkup = monthsKeyboard
	_, err := bot.Send(msg)
	if err != nil {
		log.Printf("Error sending message: %v", err)
	}
	// Удаляем сообщение с inline-кнопками
	DeleteInlineKeyboard(bot, callback)
}

// Вспомогательная функция для генерации клавиатуры с месяцами
func GenerateMonthsKeyboard() tgbotapi.InlineKeyboardMarkup {
	var rows [][]tgbotapi.InlineKeyboardButton
	var currentRow []tgbotapi.InlineKeyboardButton

	for i := 1; i <= 12; i++ {
		monthName := time.Month(i).String() // Получаем название месяца
		button := tgbotapi.NewInlineKeyboardButtonData(monthName, fmt.Sprintf("month_%d", i))
		currentRow = append(currentRow, button)

		if len(currentRow) == 3 {
			rows = append(rows, currentRow)
			currentRow = []tgbotapi.InlineKeyboardButton{}
		}
	}

	// Добавляем оставшиеся кнопки, если они есть
	if len(currentRow) > 0 {
		rows = append(rows, currentRow)
	}

	return tgbotapi.NewInlineKeyboardMarkup(rows...)
}

func HandleMonthSelected(bot *tgbotapi.BotAPI, callback *tgbotapi.CallbackQuery, state states.UserState, selectedMonth int) {
	userID := callback.Message.Chat.ID

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
	DeleteInlineKeyboard(bot, callback)
}