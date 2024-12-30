package app

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func Run(bot *tgbotapi.BotAPI) {
    // Инициализируем роутер
    router := NewRouter(bot)

    // Регистрация команд бота
    RegisterBotCommands(bot)

    u := tgbotapi.NewUpdate(0)
    u.Timeout = 60
    updates := bot.GetUpdatesChan(u)

    for update := range updates {
        router.RouteUpdate(update) // Передаем обновления в роутер
    }
}