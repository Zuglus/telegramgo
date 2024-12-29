package main

import (
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// answerCallbackQuery отвечает на callback query, убирая "часики" на кнопке.
func answerCallbackQuery(bot *tgbotapi.BotAPI, callbackQuery *tgbotapi.CallbackQuery) {
	callbackConfig := tgbotapi.NewCallback(callbackQuery.ID, "")
	if _, err := bot.Request(callbackConfig); err != nil {
		log.Printf("Error responding to callback: %v", err)
	}
}