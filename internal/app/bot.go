package app

import (
	"strconv"
	"strings"

	"telegramgo/internal/domain"
	"telegramgo/internal/handlers/contributions"
	"telegramgo/internal/handlers/core"
	"telegramgo/internal/handlers/debts"
	"telegramgo/internal/states"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func Run(bot *tgbotapi.BotAPI) {
	// Регистрация команд бота
	RegisterBotCommands(bot)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message != nil {
			processMessage(bot, update)
		} else if update.CallbackQuery != nil {
			processCallback(bot, update)
		}
	}
}

func processMessage(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	userID := update.Message.Chat.ID
	state, exists := states.UserStates[userID]

	if update.Message.IsCommand() {
		switch update.Message.Command() {
		case "start":
			core.HandleStart(bot, update)
		case "help":
			core.HandleHelp(bot, update)
		case "setstartdate":
			core.HandleSetStartDateCommand(bot, update)
		default:
			core.HandleUnknownCommand(bot, update)
		}
	} else if exists {
		// Если пользователь находится в каком-то состоянии, обрабатываем сообщение в соответствии с этим состоянием
		switch state.Stage {
		case "awaiting_name":
			contributions.HandleNameInput(bot, update)
		case "awaiting_amount":
			contributions.HandleAmountInput(bot, update)
		case "awaiting_payment_month":
			contributions.HandleMonthInput(bot, update)
		case "awaiting_year":
			debts.HandleYearInput(bot, update)
		case "awaiting_month_number":
			debts.HandleMonthNumberInput(bot, update)
		case "awaiting_member_name_for_start_date":
			core.HandleMemberNameForStartDateInput(bot, update)
		case "awaiting_start_date":
			core.HandleStartDateInput(bot, update)
		case "awaiting_contribution_date":
			contributions.HandleContributionDateInput(bot, update)
		}
	} else {
		// Обработка обычных сообщений
		core.HandleUnknownMessage(bot, update)
	}
}

func processCallback(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	callback := update.CallbackQuery
	userID := callback.Message.Chat.ID
	// Получаем текущее состояние пользователя из хранилища состояний
	state, exists := states.UserStates[userID]
	// Если состояние пользователя не найдено, устанавливаем состояние "idle"
	if !exists {
		state = states.UserState{Stage: "idle", TempMember: domain.Member{}}
	}

	switch callback.Data {
	case "add_contribution":
		contributions.HandleAddContribution(bot, callback)
	case "show_contributions":
		contributions.HandleShowContributions(bot, callback)
	case "show_debts":
		debts.HandleShowDebts(bot, callback)
	case "current_month", "previous_month":
		contributions.HandleSelectMonth(bot, callback, state)
	case "debt":
		debts.HandleDebt(bot, callback)
	case "select_year":
		debts.HandleSelectYear(bot, callback, state)
	case "confirm_contribution":
		contributions.HandleConfirmContribution(bot, callback)
	case "reject_contribution":
		contributions.HandleRejectContribution(bot, callback)
	default:
		if strings.HasPrefix(callback.Data, "year_") {
			selectedYear, _ := strconv.Atoi(strings.TrimPrefix(callback.Data, "year_"))
			debts.HandleYearSelected(bot, callback, state, selectedYear)
		} else if strings.HasPrefix(callback.Data, "month_") {
			selectedMonth, _ := strconv.Atoi(strings.TrimPrefix(callback.Data, "month_"))
			debts.HandleMonthSelected(bot, callback, state, selectedMonth)
		}
	}
}