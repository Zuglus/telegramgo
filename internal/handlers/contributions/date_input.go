package contributions

import (
	"log"
	"time"

	"telegramgo/internal/repository"
	"telegramgo/internal/states"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func HandleContributionDateInput(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	userID := update.Message.Chat.ID
	state := states.UserStates[userID]

	// Пытаемся распарсить введенную дату
	contributionDate, err := time.Parse("2006-01-02", update.Message.Text)
	if err != nil {
		msg := tgbotapi.NewMessage(userID, "Неверный формат даты. Пожалуйста, используйте ГГГГ-ММ-ДД.")
		_, _ = bot.Send(msg)
		return
	}

	// Записываем данные в БД
	memberID, err := repository.AddOrUpdateMember(state.TempMember.Name)
	if err != nil {
		log.Printf("Failed to add or update member: %v", err)
		msg := tgbotapi.NewMessage(userID, "Ошибка при сохранении данных в БД.")
		_, _ = bot.Send(msg)
		return
	}

	// Добавляем взнос в базу данных с указанной датой
	err = repository.AddContribution(memberID, state.TempAmount, contributionDate.Format("2006-01-02"), state.TempMember.Months[0])
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