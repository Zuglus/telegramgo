package app

import (
	"strconv"
	"strings"

	"telegramgo/internal/handlers/contributions"
	"telegramgo/internal/handlers/core"
	"telegramgo/internal/handlers/debts"
	"telegramgo/internal/states"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (r *Router) routeMessage(update tgbotapi.Update) {
	userID := update.Message.Chat.ID
	state, exists := states.UserStates[userID]

	if update.Message.IsCommand() {
		r.handleCommand(update)
		return // Важно! Добавляем return, чтобы не обрабатывать команду как состояние
	}

	if exists {
		r.handleStatefulMessage(update, state)
	} else {
		core.HandleUnknownMessage(r.bot, update)
	}
}

func (r *Router) handleCommand(update tgbotapi.Update) {
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

func (r *Router) handleStatefulMessage(update tgbotapi.Update, state states.UserState) {
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

func (r *Router) routeCallback(update tgbotapi.Update) {
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

func (r *Router) handlePaymentMonthCallback(callback *tgbotapi.CallbackQuery, state states.UserState) {
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

func (r *Router) handleDebtCallback(callback *tgbotapi.CallbackQuery, state states.UserState) {
	if callback.Data == "select_year" {
		debts.HandleSelectYear(r.bot, callback, state)
		return
	}

	if strings.HasPrefix(callback.Data, "year_") {
		selectedYear, _ := strconv.Atoi(strings.TrimPrefix(callback.Data, "year_"))
		debts.HandleYearSelected(r.bot, callback, state, selectedYear)
		return
	}

	if strings.HasPrefix(callback.Data, "month_") {
		selectedMonth, _ := strconv.Atoi(strings.TrimPrefix(callback.Data, "month_"))
		debts.HandleMonthSelected(r.bot, callback, state, selectedMonth)
		return
	}
}