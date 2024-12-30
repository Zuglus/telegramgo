package domain

import (
	"fmt"
	"log"
	"time"
)

// CalculateDebt рассчитывает долг участника.
func CalculateDebt(member *Member) (int, error) {
	if member.StartDate == "" {
		return 0, nil
	}

	startDate, err := time.Parse("2006-01-02", member.StartDate)
	if err != nil {
		return 0, fmt.Errorf("error parsing start date: %w", err)
	}

	// Приводим оплаченные месяцы к формату time.Time для упрощения сравнения
	paidMonthsTime := make(map[time.Time]bool)
	for _, monthStr := range member.Months {
		monthTime, err := time.Parse("2006-01", monthStr)
		if err != nil {
			log.Printf("Error parsing paid month: %v", err)
			continue
		}
		paidMonthsTime[monthTime] = true
	}

	now := time.Now()
	unpaidMonths := 0
	for startDate.Before(now) {
		// Проверяем, был ли оплачен текущий месяц
		if _, ok := paidMonthsTime[time.Date(startDate.Year(), startDate.Month(), 1, 0, 0, 0, 0, time.UTC)]; !ok {
			unpaidMonths++
		}
		startDate = startDate.AddDate(0, 1, 0) // Переходим к следующему месяцу
	}

	return unpaidMonths, nil
}