package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	bot, err := tgbotapi.NewBotAPI(os.Getenv("TELEGRAM_BOT_TOKEN"))
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true
	log.Printf("Authorized on account %s", bot.Self.UserName)

	initDB()

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message != nil {
			processMessage(bot, update)
		} else if update.CallbackQuery != nil {
			processCallback(bot, update)
		}
	}
}

func processMessage(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	userID := update.Message.Chat.ID
	state, exists := userStates[userID]

	if update.Message.IsCommand() {
		switch update.Message.Command() {
		case "start":
			handleStart(bot, update)
		case "help":
			handleHelp(bot, update)
		default:
			handleUnknownCommand(bot, update)
		}
	} else if exists {
		// Если пользователь находится в каком-то состоянии, обрабатываем сообщение в соответствии с этим состоянием
		switch state.Stage {
		case "awaiting_name":
			handleNameInput(bot, update)
		case "awaiting_amount":
			handleAmountInput(bot, update)
		case "awaiting_payment_month":
			handleMonthInput(bot, update)
		case "awaiting_year":
			handleYearInput(bot, update)
		case "awaiting_month_number":
			handleMonthNumberInput(bot, update)
		}
	} else {
		// Обработка обычных сообщений
		handleUnknownMessage(bot, update)
	}
}

func processCallback(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	callback := update.CallbackQuery
	userID := callback.Message.Chat.ID
	// Получаем текущее состояние пользователя из хранилища состояний
	state, exists := userStates[userID]
	// Если состояние пользователя не найдено, устанавливаем состояние "idle"
	if !exists {
		state = UserState{Stage: "idle", TempMember: Member{}}
	}

	switch callback.Data {
	case "add_contribution":
		handleAddContribution(bot, callback)
	case "show_contributions":
		handleShowContributions(bot, callback)
	case "show_debts":
		handleShowDebts(bot, callback)
	case "current_month", "previous_month":
		handleSelectMonth(bot, callback, state)
	case "debt":
		handleDebt(bot, callback)
	case "select_year":
		handleSelectYear(bot, callback, state)
	case "confirm_contribution":
		handleConfirmContribution(bot, callback)
	case "reject_contribution":
		handleRejectContribution(bot, callback)
	default:
		if strings.HasPrefix(callback.Data, "year_") {
			selectedYear, _ := strconv.Atoi(strings.TrimPrefix(callback.Data, "year_"))
			handleYearSelected(bot, callback, state, selectedYear)
		} else if strings.HasPrefix(callback.Data, "month_") {
			selectedMonth, _ := strconv.Atoi(strings.TrimPrefix(callback.Data, "month_"))
			handleMonthSelected(bot, callback, state, selectedMonth)
		}
	}
}

func handleYearInput(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	userID := update.Message.Chat.ID
	state := userStates[userID]

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
	userStates[userID] = state

	// Запрашиваем выбор месяца
	monthsKeyboard := generateMonthsKeyboard()

	msg := tgbotapi.NewMessage(userID, fmt.Sprintf("Выбран год: %d\nВыберите месяц:", selectedYear))
	msg.ReplyMarkup = monthsKeyboard
	_, err = bot.Send(msg)
	if err != nil {
		log.Printf("Error sending message: %v", err)
	}
}

func handleMonthNumberInput(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	userID := update.Message.Chat.ID
	state := userStates[userID]

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
	userStates[userID] = state

	// Отправляем сообщение об успешно выбранном месяце
	msg := tgbotapi.NewMessage(userID, fmt.Sprintf("Выбран месяц: %s", paymentMonth))
	_, err = bot.Send(msg)
	if err != nil {
		log.Printf("Error sending message: %v", err)
	}
}