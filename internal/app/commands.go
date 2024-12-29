package app

import (
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// RegisterBotCommands регистрирует команды для бота в Telegram.
func RegisterBotCommands(bot *tgbotapi.BotAPI) {
	commands := tgbotapi.NewSetMyCommands(
		tgbotapi.BotCommand{
			Command:     "start",
			Description: "Начать работу с ботом",
		},
		tgbotapi.BotCommand{
			Command:     "help",
			Description: "Помощь",
		},
		tgbotapi.BotCommand{
			Command:     "setstartdate",
			Description: "Установить начальную дату",
		},
	)

	if _, err := bot.Request(commands); err != nil {
		log.Panic(err)
	}
}