package contributions

import (
	"log"
	"strconv"

	"telegramgo/internal/states"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func HandleAmountInput(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	userID := update.Message.Chat.ID
	state := states.UserStates[userID]

	// Сохраняем сумму взноса
	amount, err := strconv.ParseFloat(update.Message.Text, 64)
	if err != nil {
		// Если не удалось преобразовать текст в число, отправляем сообщение об ошибке
		msg := tgbotapi.NewMessage(userID, "Пожалуйста, введите корректную сумму взноса (число).")
		_, err = bot.Send(msg)
		if err != nil {
			log.Printf("Error sending message: %v", err)
		}
		return
	}
	state.TempMember.Contribution = amount
	states.UserStates[userID] = state

	// Запрашиваем выбор месяца
	monthsKeyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Текущий месяц", "current_month"),
			tgbotapi.NewInlineKeyboardButtonData("Предыдущий месяц", "previous_month"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Ввести вручную", "debt"),
		),
	)

	msg := tgbotapi.NewMessage(userID, "Выберите месяц взноса:")
	msg.ReplyMarkup = monthsKeyboard
	_, err = bot.Send(msg)
	if err != nil {
		log.Printf("Error sending message: %v", err)
	}
}