package contributions

import (
	"fmt"
	"log"

	"telegramgo/internal/repository"
	"telegramgo/internal/telegram"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func HandleShowContributions(bot *tgbotapi.BotAPI, callback *tgbotapi.CallbackQuery) {

	// Получаем список взносов из БД
	members, err := repository.GetContributions()
	if err != nil {
		log.Printf("Error getting contributions: %v", err)
		msg := tgbotapi.NewMessage(callback.Message.Chat.ID, "Ошибка при получении списка взносов.")
		_, _ = bot.Send(msg)
		return
	}

	// Формируем сообщение со списком взносов
	var messageText string
	if len(members) == 0 {
		messageText = "Взносов пока нет."
	} else {
		for _, member := range members {
			messageText += fmt.Sprintf("%s: \nМесяцы: %v\n\n", member.Name, member.Months)
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