package message

import (
	"telegramgo/internal/handlers/contributions"
	"telegramgo/internal/handlers/core"
	"telegramgo/internal/handlers/debts"
	"telegramgo/internal/states"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type StateRouter struct {
	bot *tgbotapi.BotAPI
}

func NewStateRouter(bot *tgbotapi.BotAPI) *StateRouter {
	return &StateRouter{bot: bot}
}

func (sr *StateRouter) HandleStatefulMessage(update tgbotapi.Update, state states.UserState) {
	userID := update.Message.Chat.ID
	switch state.Stage {
	case "awaiting_name":
		contributions.HandleNameInput(sr.bot, update)
	case "awaiting_amount":
		contributions.HandleAmountInput(sr.bot, update)
	case "awaiting_payment_month":
		contributions.HandleMonthInput(sr.bot, update)
	case "awaiting_year":
		debts.HandleYearInput(sr.bot, update)
	case "awaiting_month_number":
		debts.HandleMonthNumberInput(sr.bot, update)
	case "awaiting_member_name_for_start_date":
		core.HandleMemberNameForStartDateInput(sr.bot, update)
	case "awaiting_start_date":
		core.HandleStartDateInput(sr.bot, update)
	case "awaiting_contribution_date":
		contributions.HandleContributionDateInput(sr.bot, update)
	default:
		core.HandleUnknownMessage(sr.bot, update) // На случай, если попадём в неизвестное состояние
		delete(states.UserStates, userID)       // Сбрасываем состояние
	}
}