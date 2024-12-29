package main

import (
	"log"
	"os"
	"strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
)

// Структура для хранения текущего состояния пользователя
type UserState struct {
	Stage      string // "idle", "awaiting_name", "awaiting_amount"
	TempMember Member // Временные данные пользователя
}

// Хранилище состояний пользователей (для простоты используем map, но в будущем лучше заменить на БД)
var userStates = make(map[int64]UserState)

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
		}
	} else {
		// Обработка обычных сообщений
		handleUnknownMessage(bot, update)
	}
}

func processCallback(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	callback := update.CallbackQuery
	switch callback.Data {
	case "add_contribution":
		handleAddContribution(bot, callback)
	case "show_contributions":
		handleShowContributions(bot, callback)
	case "show_debts":
		handleShowDebts(bot, callback)
	}
}

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
            messageText += member.Name + ": " + strconv.FormatFloat(member.Contribution, 'f', 2, 64) + "\n"
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
	msg := tgbotapi.NewMessage(callback.Message.Chat.ID, "Функция показа долгов пока не реализована.")
	_, err := bot.Send(msg)
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

	// Записываем данные в БД
	memberID, err := addOrUpdateMember(state.TempMember.Name)
	if err != nil {
		log.Printf("Failed to add or update member: %v", err)
		msg := tgbotapi.NewMessage(userID, "Ошибка при сохранении данных в БД.")
		_, _ = bot.Send(msg)
		return
	}

	err = addContribution(memberID, state.TempMember.Contribution, "")
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