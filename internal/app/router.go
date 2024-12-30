package app

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Router struct {
	bot             *tgbotapi.BotAPI
	messageRouter   *MessageRouter
	callbackRouter  *CallbackRouter
}

func NewRouter(bot *tgbotapi.BotAPI) *Router {
	return &Router{
		bot:             bot,
		messageRouter:   NewMessageRouter(bot),
		callbackRouter:  NewCallbackRouter(bot),
	}
}

func (r *Router) RouteUpdate(update tgbotapi.Update) {
	if update.Message != nil {
		r.messageRouter.RouteMessage(update)
	} else if update.CallbackQuery != nil {
		r.callbackRouter.RouteCallback(update)
	}
}