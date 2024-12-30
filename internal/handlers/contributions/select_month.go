package contributions

import (
	"fmt"
	"log"
	"time"

	"telegramgo/internal/states"
	"telegramgo/internal/telegram"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func HandleSelectMonth(bot *tgbotapi.BotAPI, callback *tgbotapi.CallbackQuery, state states.UserState) {
	userID := callback.Message.Chat.ID

	telegram.AnswerCallbackQuery(bot, callback)

	now := time.Now()
	var paymentMonth string

	switch callback.Data {
	case "current_month":
		paymentMonth = now.Format("2006-01")
	case "previous_month":
		paymentMonth = now.AddDate(0, -1, 0).Format("2006-01")
	}

	state.TempMember.Months = append(state.TempMember.Months, paymentMonth)
	// Устанавливаем состояние "ожидание даты взноса"
	state.Stage = "awaiting_contribution_date"
	userStates[userID] = state

	// Создаем клавиатуру с кнопками "Подтвердить" и "Отклонить"
	confirmationKeyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Подтвердить", "confirm_contribution"),
			tgbotapi.NewInlineKeyboardButtonData("Отклонить", "reject_contribution"),
		),
	)

	msg := tgbotapi.NewMessage(userID, fmt.Sprintf("Выбран месяц: %s\nТеперь введите дату взноса в формате ГГГГ-ММ-ДД:", paymentMonth))
	// msg.ReplyMarkup = confirmationKeyboard
	_, err := bot.Send(msg)
	if err != nil {
		log.Printf("Error sending message: %v", err)
	}
	// Удаляем сообщение с inline-кнопками
	telegram.DeleteInlineKeyboard(bot, callback)
}