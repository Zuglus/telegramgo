package main

import (
	"log"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func main() {
	bot, err := tgbotapi.NewBotAPI(os.Getenv("TELEGRAM_BOT_TOKEN"))
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true
	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message != nil {
			if update.Message.IsCommand() {
                switch update.Message.Command() {
                case "start":
                    handleStart(bot, update)
                case "help":
                    handleHelp(bot, update)
                default:
                    handleUnknownCommand(bot, update)
                }
			} else if update.CallbackQuery != nil {
                //обработка нажатия на кнопку
            } else {
				// Обработка обычных сообщений
			}
		}
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