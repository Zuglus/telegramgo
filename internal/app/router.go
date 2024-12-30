package app

import (
	"telegramgo/internal/app/callback"
	"telegramgo/internal/app/message"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Router struct {
	bot                 *tgbotapi.BotAPI
	commandRouter       *message.CommandRouter
	stateRouter         *message.StateRouter
	contributionsRouter *callback.ContributionsCallbackRouter
	debtsRouter         *callback.DebtsCallbackRouter
}

func NewRouter(bot *tgbotapi.BotAPI) *Router {
	return &Router{
		bot:                 bot,
		commandRouter:       message.NewCommandRouter(bot),
		stateRouter:         message.NewStateRouter(bot),
		contributionsRouter: callback.NewContributionsCallbackRouter(bot),
		debtsRouter:         callback.NewDebtsCallbackRouter(bot),
	}
}

func (r *Router) RouteUpdate(update tgbotapi.Update) {
	if update.Message != nil {
		r.routeMessage(update)
	} else if update.CallbackQuery != nil {
		r.routeCallback(update)
	}
}