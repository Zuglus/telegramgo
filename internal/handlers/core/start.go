package core

import (
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func HandleStart(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
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