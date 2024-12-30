package app

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func Run(bot *tgbotapi.BotAPI) {
	// Регистрация команд бота
	RegisterBotCommands(bot)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := bot.GetUpdatesChan(u)

	router := NewRouter(bot) // Инициализируем роутер

	for update := range updates {
		router.RouteUpdate(update) // Передаем обновления в роутер
	}
}