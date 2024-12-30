package message

import (
	"telegramgo/internal/handlers/core"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type CommandRouter struct {
	bot *tgbotapi.BotAPI
}

func NewCommandRouter(bot *tgbotapi.BotAPI) *CommandRouter {
	return &CommandRouter{bot: bot}
}

func (cr *CommandRouter) HandleCommand(update tgbotapi.Update) {
	switch update.Message.Command() {
	case "start":
		core.HandleStart(cr.bot, update)
	case "help":
		core.HandleHelp(cr.bot, update)
	case "setstartdate":
		core.HandleSetStartDateCommand(cr.bot, update)
	default:
		core.HandleUnknownCommand(cr.bot, update)
	}
}