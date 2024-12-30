package app

import (
	"telegramgo/internal/handlers/contributions"
	"telegramgo/internal/handlers/core"
	"telegramgo/internal/handlers/debts"
	"telegramgo/internal/states"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type MessageRouter struct {
	bot *tgbotapi.BotAPI
}

func NewMessageRouter(bot *tgbotapi.BotAPI) *MessageRouter {
	return &MessageRouter{bot: bot}
}

func (r *MessageRouter) RouteMessage(update tgbotapi.Update) {
	userID := update.Message.Chat.ID
	state, exists := states.UserStates[userID]

	if update.Message.IsCommand() {
		r.handleCommand(update)
		return
	}

	if exists {
		r.handleStatefulMessage(update, state)
	} else {
		core.HandleUnknownMessage(r.bot, update)
	}
}

func (r *MessageRouter) handleCommand(update tgbotapi.Update) {
	switch update.Message.Command() {
	case "start":
		core.HandleStart(r.bot, update)
	case "help":
		core.HandleHelp(r.bot, update)
	case "setstartdate":
		core.HandleSetStartDateCommand(r.bot, update)
	default:
		core.HandleUnknownCommand(r.bot, update)
	}
}

func (r *MessageRouter) handleStatefulMessage(update tgbotapi.Update, state states.UserState) {
	userID := update.Message.Chat.ID
	switch state.Stage {
	case "awaiting_name":
		contributions.HandleNameInput(r.bot, update)
	case "awaiting_amount":
		contributions.HandleAmountInput(r.bot, update)
	case "awaiting_payment_month":
		contributions.HandleMonthInput(r.bot, update)
	case "awaiting_year":
		debts.HandleYearInput(r.bot, update)
	case "awaiting_month_number":
		debts.HandleMonthNumberInput(r.bot, update)
	case "awaiting_member_name_for_start_date":
		core.HandleMemberNameForStartDateInput(r.bot, update)
	case "awaiting_start_date":
		core.HandleStartDateInput(r.bot, update)
	case "awaiting_contribution_date":
		contributions.HandleContributionDateInput(r.bot, update)
	default:
		core.HandleUnknownMessage(r.bot, update) // На случай, если попадём в неизвестное состояние
		delete(states.UserStates, userID)       // Сбрасываем состояние
	}
}