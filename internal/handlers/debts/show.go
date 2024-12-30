package debts

import (
	"fmt"
	"log"

	"telegramgo/internal/repository"
	"telegramgo/internal/telegram"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func HandleShowDebts(bot *tgbotapi.BotAPI, callback *tgbotapi.CallbackQuery) {
	// Получаем список долгов из БД
	members, err := repository.GetDebts()
	if err != nil {
		log.Printf("Error getting debts: %v", err)
		msg := tgbotapi.NewMessage(callback.Message.Chat.ID, "Ошибка при получении списка долгов.")
		_, _ = bot.Send(msg)
		return
	}

	// Формируем сообщение со списком долгов
	var messageText string
	if len(members) == 0 {
		messageText = "Долгов пока нет."
	} else {
		for _, member := range members {
			messageText += fmt.Sprintf("%s: %.0f\nОплаченные месяцы: %v\n\n", member.Name, member.Debt, member.Months)
		}
	}

	// Отправляем сообщение пользователю
	msg := tgbotapi.NewMessage(callback.Message.Chat.ID, messageText)
	_, err = bot.Send(msg)
	if err != nil {
		log.Printf("Error sending message: %v", err)
	}
	// Удаляем сообщение с inline-кнопками
	telegram.DeleteInlineKeyboard(bot, callback)
}