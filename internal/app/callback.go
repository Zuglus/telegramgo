package app

import (
	"strconv"
	"strings"

	"telegramgo/internal/handlers/contributions"
	"telegramgo/internal/handlers/debts"
	"telegramgo/internal/states"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type CallbackRouter struct {
	bot *tgbotapi.BotAPI
}

func NewCallbackRouter(bot *tgbotapi.BotAPI) *CallbackRouter {
	return &CallbackRouter{bot: bot}
}

func (r *CallbackRouter) RouteCallback(update tgbotapi.Update) {
	callback := update.CallbackQuery
	state, exists := states.UserStates[callback.Message.Chat.ID]
	if !exists {
		// Если состояние пользователя не найдено, то это, скорее всего, команда из главного меню.
		// Обрабатываем её напрямую.
		switch callback.Data {
		case "add_contribution":
			contributions.HandleAddContribution(r.bot, callback)
		case "show_contributions":
			contributions.HandleShowContributions(r.bot, callback)
		case "show_debts":
			debts.HandleShowDebts(r.bot, callback)
		default:
			// Обработка неизвестных callback.Data для команд из главного меню
		}
		return // Прерываем выполнение, чтобы не обрабатывать callback как состояние
	}

	// Если состояние найдено, то продолжаем обработку в соответствии с состоянием
	switch state.Stage {
	case "awaiting_payment_month":
		r.handlePaymentMonthCallback(callback, state)
	case "awaiting_year", "awaiting_month_number":
		r.handleDebtCallback(callback, state)
	default:
		// Обработка неизвестных callback.Data для состояний
	}
}

func (r *CallbackRouter) handlePaymentMonthCallback(callback *tgbotapi.CallbackQuery, state states.UserState) {
	switch callback.Data {
	case "current_month", "previous_month":
		contributions.HandleSelectMonth(r.bot, callback, state)
	case "debt":
		debts.HandleDebt(r.bot, callback)
	case "confirm_contribution":
		contributions.HandleConfirmContribution(r.bot, callback)
	case "reject_contribution":
		contributions.HandleRejectContribution(r.bot, callback)
	}
}

func (r *CallbackRouter) handleDebtCallback(callback *tgbotapi.CallbackQuery, state states.UserState) {
	switch callback.Data {
	case "select_year":
		debts.HandleSelectYear(r.bot, callback, state)
	default:
		if strings.HasPrefix(callback.Data, "year_") {
			selectedYear, _ := strconv.Atoi(strings.TrimPrefix(callback.Data, "year_"))
			debts.HandleYearSelected(r.bot, callback, state, selectedYear)
		} else if strings.HasPrefix(callback.Data, "month_") {
			selectedMonth, _ := strconv.Atoi(strings.TrimPrefix(callback.Data, "month_"))
			debts.HandleMonthSelected(r.bot, callback, state, selectedMonth)
		}
	}
}