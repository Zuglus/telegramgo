package states

import "telegramgo/internal/domain"

// Структура для хранения текущего состояния пользователя
type UserState struct {
	Stage       string         // "idle", "awaiting_name", "awaiting_amount", "awaiting_payment_month", "awaiting_year", "awaiting_month_number"
	TempMember  domain.Member // Временные данные пользователя
	SelectedYear int          // Выбранный год
	Months      []Month      // Доступные месяцы для выбора
	TempAmount  float64      // Временное хранение суммы взноса
}

// Определение структуры Month
type Month struct {
	Name  string
	Value int
}

// Хранилище состояний пользователей (для простоты используем map, но в будущем лучше заменить на БД)
var UserStates = make(map[int64]UserState)