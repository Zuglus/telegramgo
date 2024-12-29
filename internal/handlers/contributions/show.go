package contributions

import (
	"fmt"
	"log"

	"telegramgo/internal/repository"

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
            messageText += fmt.Sprintf("%s: %.2f\nМесяцы: %v\n\n", member.Name, member.Contribution, member.Months)
        }
    }

    // Отправляем сообщение пользователю
    msg := tgbotapi.NewMessage(callback.Message.Chat.ID, messageText)
    _, err = bot.Send(msg)
    if err != nil {
        log.Printf("Error sending message: %v", err)
    }
    // Удаляем сообщение с inline-кнопками
    deleteMsg := tgbotapi.NewDeleteMessage(callback.Message.Chat.ID, callback.Message.MessageID)
    _, err = bot.Send(deleteMsg)
    if err != nil {
        log.Printf("Failed to delete message: %v", err)
    }
}