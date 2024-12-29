package contributions

import (
	"log"
	"time"

	"telegramgo/internal/repository"
	"telegramgo/internal/states"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func HandleMonthInput(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	userID := update.Message.Chat.ID
	state := states.UserStates[userID]

	// Получаем текущую дату
	currentDate := time.Now().Format("2006-01-02")

	// Записываем данные в БД
	memberID, err := repository.AddOrUpdateMember(state.TempMember.Name)
	if err != nil {
		log.Printf("Failed to add or update member: %v", err)
		msg := tgbotapi.NewMessage(userID, "Ошибка при сохранении данных в БД.")
		_, _ = bot.Send(msg)
		return
	}

	err = repository.AddContribution(memberID, state.TempMember.Contribution, currentDate, state.TempMember.Months[0])
	if err != nil {
		log.Printf("Failed to add contribution: %v", err)
		msg := tgbotapi.NewMessage(userID, "Ошибка при сохранении данных в БД.")
		_, _ = bot.Send(msg)
		return
	}

	// Очищаем состояние пользователя
	delete(states.UserStates, userID)

	// Отправляем подтверждение
	msg := tgbotapi.NewMessage(userID, "Взнос успешно добавлен!")
	_, err = bot.Send(msg)
	if err != nil {
		log.Printf("Error sending message: %v", err)
	}
}