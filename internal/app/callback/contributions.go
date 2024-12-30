package callback

import (
	"telegramgo/internal/handlers/contributions"
	"telegramgo/internal/states"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type ContributionsCallbackRouter struct {
	bot *tgbotapi.BotAPI
}

func NewContributionsCallbackRouter(bot *tgbotapi.BotAPI) *ContributionsCallbackRouter {
	return &ContributionsCallbackRouter{bot: bot}
}

func (cr *ContributionsCallbackRouter) HandleCallback(update tgbotapi.Update) {
	callback := update.CallbackQuery
	state := states.UserStates[callback.Message.Chat.ID]

	switch {
	case callback.Data == "add_contribution":
		contributions.HandleAddContribution(cr.bot, callback)
	case callback.Data == "show_contributions":
		contributions.HandleShowContributions(cr.bot, callback)
	case callback.Data == "current_month" || callback.Data == "previous_month":
		contributions.HandleSelectMonth(cr.bot, callback, state)
	case callback.Data == "confirm_contribution":
		contributions.HandleConfirmContribution(cr.bot, callback)
	case callback.Data == "reject_contribution":
		contributions.HandleRejectContribution(cr.bot, callback)
	}
}