package core

import (
	"log"
	"time"

	"telegramgo/internal/repository"
	"telegramgo/internal/states"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func HandleSetStartDateCommand(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	userID := update.Message.Chat.ID

	// Проверяем, есть ли у пользователя состояние
	_, exists := states.UserStates[userID]
	if exists {
		// Сбрасываем состояние пользователя, если оно было установлено
		delete(states.UserStates, userID)
	}

	// Устанавливаем состояние пользователя на ожидание имени участника
	states.UserStates[userID] = states.UserState{Stage: "awaiting_member_name_for_start_date"}

	// Запрашиваем у пользователя имя участника
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Введите имя участника:")
	_, err := bot.Send(msg)
	if err != nil {
		log.Printf("Error sending message: %v", err)
	}
}

func HandleMemberNameForStartDateInput(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	userID := update.Message.Chat.ID
	memberName := update.Message.Text

	// Проверяем существование участника
	_, err := repository.GetMember(memberName)
	if err != nil {
		log.Printf("Error getting member: %v", err)
		msg := tgbotapi.NewMessage(userID, "Участник с таким именем не найден.")
		_, _ = bot.Send(msg)
		// Сбрасываем состояние пользователя, так как участник не найден
		delete(states.UserStates, userID)
		return
	}

	// Устанавливаем состояние пользователя на ожидание начальной даты
	state := states.UserStates[userID]
	state.TempMember.Name = memberName // Сохраняем имя участника
	state.Stage = "awaiting_start_date"
	states.UserStates[userID] = state

	// Запрашиваем у пользователя начальную дату
	msg := tgbotapi.NewMessage(userID, "Введите начальную дату в формате ГГГГ-ММ-ДД:")
	_, err = bot.Send(msg)
	if err != nil {
		log.Printf("Error sending message: %v", err)
	}
}

func HandleStartDateInput(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	userID := update.Message.Chat.ID
	startDate := update.Message.Text
	state := states.UserStates[userID]

	// Проверяем формат даты
	_, err := time.Parse("2006-01-02", startDate)
	if err != nil {
		msg := tgbotapi.NewMessage(userID, "Неверный формат даты. Пожалуйста, используйте ГГГГ-ММ-ДД.")
		_, _ = bot.Send(msg)
		return
	}

	// Сохраняем начальную дату в БД
	err = repository.SetMemberStartDate(state.TempMember.Name, startDate)
	if err != nil {
		log.Printf("Error setting start date: %v", err)
		msg := tgbotapi.NewMessage(userID, "Ошибка при сохранении начальной даты.")
		_, _ = bot.Send(msg)
		return
	}

	// Очищаем состояние пользователя
	delete(states.UserStates, userID)

	// Отправляем подтверждение
	msg := tgbotapi.NewMessage(userID, "Начальная дата успешно установлена.")
	_, err = bot.Send(msg)
	if err != nil {
		log.Printf("Error sending message: %v", err)
	}
}