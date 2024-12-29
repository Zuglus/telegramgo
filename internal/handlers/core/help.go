package core

import (
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func HandleHelp(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "/start - начать\n/help - помощь")
	_, err := bot.Send(msg)
	if err != nil {
		log.Printf("Error sending message: %v", err)
	}
}