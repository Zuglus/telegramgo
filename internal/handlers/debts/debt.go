package debts

import (
	"log"
	"strconv"
	"time"

	"telegramgo/internal/telegram"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func HandleDebt(bot *tgbotapi.BotAPI, callback *tgbotapi.CallbackQuery) {
	userID := callback.Message.Chat.ID

	// Запрашиваем выбор года
	currentYear := time.Now().Year()
	yearsKeyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(strconv.Itoa(currentYear), "year_"+strconv.Itoa(currentYear)),
			tgbotapi.NewInlineKeyboardButtonData(strconv.Itoa(currentYear-1), "year_"+strconv.Itoa(currentYear-1)),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Другой год", "select_year"),
		),
	)

	msg := tgbotapi.NewMessage(userID, "Выберите год:")
	msg.ReplyMarkup = yearsKeyboard
	_, err := bot.Send(msg)
	if err != nil {
		log.Printf("Error sending message: %v", err)
	}
	// Удаляем сообщение с inline-кнопками
	telegram.DeleteInlineKeyboard(bot, callback)
}