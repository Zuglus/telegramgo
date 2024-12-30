package callback

import (
	"strconv"
	"strings"

	"telegramgo/internal/handlers/debts"
	"telegramgo/internal/states"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type DebtsCallbackRouter struct {
	bot *tgbotapi.BotAPI
}

func NewDebtsCallbackRouter(bot *tgbotapi.BotAPI) *DebtsCallbackRouter {
	return &DebtsCallbackRouter{bot: bot}
}

func (dr *DebtsCallbackRouter) HandleCallback(update tgbotapi.Update) {
	callback := update.CallbackQuery
	state := states.UserStates[callback.Message.Chat.ID]

	switch {
	case callback.Data == "show_debts":
		debts.HandleShowDebts(dr.bot, callback)
	case callback.Data == "debt":
		debts.HandleDebt(dr.bot, callback)
	case callback.Data == "select_year":
		debts.HandleSelectYear(dr.bot, callback, state)
	case strings.HasPrefix(callback.Data, "year_"):
		dr.handleYearCallback(callback, state)
	case strings.HasPrefix(callback.Data, "month_"):
		dr.handleMonthCallback(callback, state)
	}
}

func (dr *DebtsCallbackRouter) handleYearCallback(callback *tgbotapi.CallbackQuery, state states.UserState) {
	selectedYear, _ := strconv.Atoi(strings.TrimPrefix(callback.Data, "year_"))
	debts.HandleYearSelected(dr.bot, callback, state, selectedYear)
}

func (dr *DebtsCallbackRouter) handleMonthCallback(callback *tgbotapi.CallbackQuery, state states.UserState) {
	selectedMonth, _ := strconv.Atoi(strings.TrimPrefix(callback.Data, "month_"))
	debts.HandleMonthSelected(dr.bot, callback, state, selectedMonth)
}