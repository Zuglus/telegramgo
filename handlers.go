package main

import (
	"fmt"
	"log"
	"strconv"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func handleStart(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	var keyboard = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Добавить взнос", "add_contribution"),
			tgbotapi.NewInlineKeyboardButtonData("Показать взносы", "show_contributions"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Показать долги", "show_debts"),
		),
	)

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Выберите действие:")
	msg.ReplyMarkup = keyboard
	_, err := bot.Send(msg)
	if err != nil {
		log.Printf("Error sending message: %v", err)
	}
}

func handleHelp(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "/start - начать\n/help - помощь")
	_, err := bot.Send(msg)
	if err != nil {
		log.Printf("Error sending message: %v", err)
	}
}

func handleUnknownCommand(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Неизвестная команда")
	_, err := bot.Send(msg)
	if err != nil {
		log.Printf("Error sending message: %v", err)
	}
}

func handleUnknownMessage(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Нажмите /start")
	_, err := bot.Send(msg)
	if err != nil {
		log.Printf("Error sending message: %v", err)
	}
}

func handleAddContribution(bot *tgbotapi.BotAPI, callback *tgbotapi.CallbackQuery) {
	// Отвечаем на callback, чтобы убрать "часики" на кнопке
	callbackConfig := tgbotapi.NewCallback(callback.ID, "")
	if _, err := bot.Request(callbackConfig); err != nil {
		log.Printf("Error responding to callback: %v", err)
	}

	// Начинаем процесс добавления взноса
	userID := callback.Message.Chat.ID
	userStates[userID] = UserState{Stage: "awaiting_name", TempMember: Member{}}

	// Отправляем сообщение пользователю с запросом имени
	msg := tgbotapi.NewMessage(callback.Message.Chat.ID, "Введите имя:")
	_, err := bot.Send(msg)
	if err != nil {
		log.Printf("Error sending message: %v", err)
	}
}

func handleShowContributions(bot *tgbotapi.BotAPI, callback *tgbotapi.CallbackQuery) {
	// Отвечаем на callback, чтобы убрать "часики" на кнопке
	callbackConfig := tgbotapi.NewCallback(callback.ID, "")
	if _, err := bot.Request(callbackConfig); err != nil {
		log.Printf("Error responding to callback: %v", err)
	}

	// Получаем список взносов из БД
	members, err := getContributions()
	if err != nil {
		log.Printf("Error getting contributions: %v", err)
		msg := tgbotapi.NewMessage(callback.Message.Chat.ID, "Ошибка при получении списка взносов.")
		_, _ = bot.Send(msg)
		return
	}

	// Формируем сообщение со списком взносов
	var messageText string
	if len(members) == 0 {
		messageText = "Взносов пока нет."
	} else {
		for _, member := range members {
			messageText += fmt.Sprintf("%s: %.2f\nМесяцы: %v\n\n", member.Name, member.Contribution, member.Months)
		}
	}

	// Отправляем сообщение пользователю
	msg := tgbotapi.NewMessage(callback.Message.Chat.ID, messageText)
	_, err = bot.Send(msg)
	if err != nil {
		log.Printf("Error sending message: %v", err)
	}
}

func handleShowDebts(bot *tgbotapi.BotAPI, callback *tgbotapi.CallbackQuery) {
	// Отвечаем на callback, чтобы убрать "часики" на кнопке
	callbackConfig := tgbotapi.NewCallback(callback.ID, "")
	if _, err := bot.Request(callbackConfig); err != nil {
		log.Printf("Error responding to callback: %v", err)
	}

	// Получаем список долгов из БД
	members, err := getDebts()
	if err != nil {
		log.Printf("Error getting debts: %v", err)
		msg := tgbotapi.NewMessage(callback.Message.Chat.ID, "Ошибка при получении списка долгов.")
		_, _ = bot.Send(msg)
		return
	}

	// Формируем сообщение со списком долгов
	var messageText string
	if len(members) == 0 {
		messageText = "Долгов пока нет."
	} else {
		for _, member := range members {
			messageText += fmt.Sprintf("%s: %.2f\nОплаченные месяцы: %v\n\n", member.Name, member.Debt, member.Months)
		}
	}

	// Отправляем сообщение пользователю
	msg := tgbotapi.NewMessage(callback.Message.Chat.ID, messageText)
	_, err = bot.Send(msg)
	if err != nil {
		log.Printf("Error sending message: %v", err)
	}
}

func handleNameInput(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	userID := update.Message.Chat.ID
	state := userStates[userID]

	// Сохраняем имя пользователя
	state.TempMember.Name = update.Message.Text
	state.Stage = "awaiting_amount"
	userStates[userID] = state

	// Запрашиваем сумму взноса
	msg := tgbotapi.NewMessage(userID, "Введите сумму взноса:")
	_, err := bot.Send(msg)
	if err != nil {
		log.Printf("Error sending message: %v", err)
	}
}

func handleAmountInput(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	userID := update.Message.Chat.ID
	state := userStates[userID]

	// Сохраняем сумму взноса
	amount, err := strconv.ParseFloat(update.Message.Text, 64)
	if err != nil {
		// Если не удалось преобразовать текст в число, отправляем сообщение об ошибке
		msg := tgbotapi.NewMessage(userID, "Пожалуйста, введите корректную сумму взноса (число).")
		_, err = bot.Send(msg)
		if err != nil {
			log.Printf("Error sending message: %v", err)
		}
		return
	}
	state.TempMember.Contribution = amount
	userStates[userID] = state

	// Запрашиваем выбор месяца
	monthsKeyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Текущий месяц", "current_month"),
			tgbotapi.NewInlineKeyboardButtonData("Предыдущий месяц", "previous_month"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Ввести вручную", "debt"),
		),
	)

	msg := tgbotapi.NewMessage(userID, "Выберите месяц взноса:")
	msg.ReplyMarkup = monthsKeyboard
	_, err = bot.Send(msg)
	if err != nil {
		log.Printf("Error sending message: %v", err)
	}
}

// ... (предыдущий код) ...

func handleMonthInput(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	userID := update.Message.Chat.ID
	state := userStates[userID]

	// Получаем текущую дату
	currentDate := time.Now().Format("2006-01-02")

	// Записываем данные в БД
	memberID, err := addOrUpdateMember(state.TempMember.Name)
	if err != nil {
		log.Printf("Failed to add or update member: %v", err)
		msg := tgbotapi.NewMessage(userID, "Ошибка при сохранении данных в БД.")
		_, _ = bot.Send(msg)
		return
	}

	// Добавляем взнос в базу данных
	err = addContribution(memberID, state.TempMember.Contribution, currentDate, state.TempMember.Months[0])
	if err != nil {
		log.Printf("Failed to add contribution: %v", err)
		msg := tgbotapi.NewMessage(userID, "Ошибка при сохранении данных в БД.")
		_, _ = bot.Send(msg)
		return
	}

	// Очищаем состояние пользователя
	delete(userStates, userID)

	// Отправляем подтверждение
	msg := tgbotapi.NewMessage(userID, "Взнос успешно добавлен!")
	_, err = bot.Send(msg)
	if err != nil {
		log.Printf("Error sending message: %v", err)
	}
}

func handleSelectMonth(bot *tgbotapi.BotAPI, callback *tgbotapi.CallbackQuery, state UserState) {
	userID := callback.Message.Chat.ID

	// Отвечаем на callback, чтобы убрать "часики" на кнопке
	callbackConfig := tgbotapi.NewCallback(callback.ID, "")
	if _, err := bot.Request(callbackConfig); err != nil {
		log.Printf("Error responding to callback: %v", err)
	}

	now := time.Now()
	var paymentMonth string

	switch callback.Data {
	case "current_month":
		paymentMonth = now.Format("2006-01")
	case "previous_month":
		paymentMonth = now.AddDate(0, -1, 0).Format("2006-01")
	}

	state.TempMember.Months = append(state.TempMember.Months, paymentMonth)
	state.Stage = "awaiting_confirmation" // Изменяем состояние на "ожидание подтверждения"
	userStates[userID] = state

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
}

func handleDebt(bot *tgbotapi.BotAPI, callback *tgbotapi.CallbackQuery) {
	userID := callback.Message.Chat.ID
	// Отвечаем на callback, чтобы убрать "часики" на кнопке
	callbackConfig := tgbotapi.NewCallback(callback.ID, "")
	if _, err := bot.Request(callbackConfig); err != nil {
		log.Printf("Error responding to callback: %v", err)
	}

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
}

func handleSelectYear(bot *tgbotapi.BotAPI, callback *tgbotapi.CallbackQuery, state UserState) {
	userID := callback.Message.Chat.ID

	// Отвечаем на callback, чтобы убрать "часики" на кнопке
	callbackConfig := tgbotapi.NewCallback(callback.ID, "")
	if _, err := bot.Request(callbackConfig); err != nil {
		log.Printf("Error responding to callback: %v", err)
	}

	state.Stage = "awaiting_year"
	userStates[userID] = state

	// Запрашиваем ввод года
	msg := tgbotapi.NewMessage(userID, "Введите год в формате ГГГГ:")
	_, err := bot.Send(msg)
	if err != nil {
		log.Printf("Error sending message: %v", err)
	}
}

func handleYearSelected(bot *tgbotapi.BotAPI, callback *tgbotapi.CallbackQuery, state UserState, selectedYear int) {
	userID := callback.Message.Chat.ID

	// Отвечаем на callback, чтобы убрать "часики" на кнопке
	callbackConfig := tgbotapi.NewCallback(callback.ID, "")
	if _, err := bot.Request(callbackConfig); err != nil {
		log.Printf("Error responding to callback: %v", err)
	}

	// Обновляем состояние пользователя с выбранным годом
	state.SelectedYear = selectedYear
	state.Stage = "awaiting_month_number"
	userStates[userID] = state

	// Запрашиваем выбор месяца
	monthsKeyboard := generateMonthsKeyboard()

	msg := tgbotapi.NewMessage(userID, fmt.Sprintf("Выбран год: %d\nВыберите месяц:", selectedYear))
	msg.ReplyMarkup = monthsKeyboard
	_, err := bot.Send(msg)
	if err != nil {
		log.Printf("Error sending message: %v", err)
	}
}

// Вспомогательная функция для генерации клавиатуры с месяцами
func generateMonthsKeyboard() tgbotapi.InlineKeyboardMarkup {
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

func handleMonthSelected(bot *tgbotapi.BotAPI, callback *tgbotapi.CallbackQuery, state UserState, selectedMonth int) {
	userID := callback.Message.Chat.ID

	// Отвечаем на callback, чтобы убрать "часики" на кнопке
	callbackConfig := tgbotapi.NewCallback(callback.ID, "")
	if _, err := bot.Request(callbackConfig); err != nil {
		log.Printf("Error responding to callback: %v", err)
	}

	// Формируем месяц в формате ГГГГ-ММ
	paymentMonth := fmt.Sprintf("%d-%02d", state.SelectedYear, selectedMonth)

	// Обновляем состояние пользователя
	state.TempMember.Months = append(state.TempMember.Months, paymentMonth)
	state.Stage = "awaiting_payment_month"
	userStates[userID] = state

	// Отправляем сообщение об успешно выбранном месяце
	msg := tgbotapi.NewMessage(userID, fmt.Sprintf("Выбран месяц: %s", paymentMonth))
	_, err := bot.Send(msg)
	if err != nil {
		log.Printf("Error sending message: %v", err)
	}
}

// Добавляем новые функции для обработки подтверждения и отклонения
func handleConfirmContribution(bot *tgbotapi.BotAPI, callback *tgbotapi.CallbackQuery) {
    userID := callback.Message.Chat.ID
    state := userStates[userID]

    // Отвечаем на callback, чтобы убрать "часики" на кнопке
    callbackConfig := tgbotapi.NewCallback(callback.ID, "")
    if _, err := bot.Request(callbackConfig); err != nil {
        log.Printf("Error responding to callback: %v", err)
    }

    // Получаем текущую дату
    currentDate := time.Now().Format("2006-01-02")

    // Записываем данные в БД
    memberID, err := addOrUpdateMember(state.TempMember.Name)
    if err != nil {
        log.Printf("Failed to add or update member: %v", err)
        msg := tgbotapi.NewMessage(userID, "Ошибка при сохранении данных в БД.")
        _, _ = bot.Send(msg)
        return
    }

    err = addContribution(memberID, state.TempMember.Contribution, currentDate, state.TempMember.Months[0])
    if err != nil {
        log.Printf("Failed to add contribution: %v", err)
        msg := tgbotapi.NewMessage(userID, "Ошибка при сохранении данных в БД.")
        _, _ = bot.Send(msg)
        return
    }

    // Очищаем состояние пользователя
    delete(userStates, userID)

    // Отправляем подтверждение
    msg := tgbotapi.NewMessage(userID, "Взнос успешно добавлен!")
    _, err = bot.Send(msg)
    if err != nil {
        log.Printf("Error sending message: %v", err)
    }
}

func handleRejectContribution(bot *tgbotapi.BotAPI, callback *tgbotapi.CallbackQuery) {
	userID := callback.Message.Chat.ID

	// Отвечаем на callback, чтобы убрать "часики" на кнопке
	callbackConfig := tgbotapi.NewCallback(callback.ID, "")
	if _, err := bot.Request(callbackConfig); err != nil {
		log.Printf("Error responding to callback: %v", err)
	}

	// Очищаем состояние пользователя
	delete(userStates, userID)

	// Отправляем сообщение об отмене
	msg := tgbotapi.NewMessage(userID, "Ввод взноса отменен.")
	_, err := bot.Send(msg)
	if err != nil {
		log.Printf("Error sending message: %v", err)
	}
}