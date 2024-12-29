package main

import (
	"log"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"google.golang.org/api/sheets/v4"
)

var srv *sheets.Service
var bot *tgbotapi.BotAPI
var spreadsheetId = "1RArmgBABcOKmn3lCM2z4xKW3KdDcrD6nJN2SosD3luQ" // Идентификатор твоей таблицы

func main() {
	// Инициализация бота
	var err error
	bot, err = tgbotapi.NewBotAPI(os.Getenv("TELEGRAM_BOT_TOKEN"))
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true
	log.Printf("Authorized on account %s", bot.Self.UserName)

	// Инициализация Google Sheets API
	srv, err = initSheetsService() // Реализуй эту функцию в sheets.go
	if err != nil {
		log.Fatalf("Unable to retrieve Sheets client: %v", err)
	}

	// Обработка обновлений
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message != nil {
			if update.Message.IsCommand() {
				handleCommand(update)
			} else {
				// Обработка обычных сообщений (например, ввод данных для нового взноса)
				handleMessage(update)
			}
		}
	}
}
