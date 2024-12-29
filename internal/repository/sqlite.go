package repository

import (
	"database/sql"
	"log"
	"strings"
	"time"

	"telegramgo/internal/domain"

	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

func InitDB() {
	var err error
	db, err = sql.Open("sqlite3", "./data/contributions.db")
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}

	createTables()
}

func createTables() {
	_, err := db.Exec(`
        CREATE TABLE IF NOT EXISTS members (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            name TEXT UNIQUE,
            start_date TEXT DEFAULT ''
        );
        CREATE TABLE IF NOT EXISTS contributions (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            member_id INTEGER,
            amount REAL,
            date TEXT,
            payment_month TEXT,
            FOREIGN KEY (member_id) REFERENCES members(id)
        );
    `)
	if err != nil {
		log.Fatalf("Failed to create tables: %v", err)
	}
}

// AddOrUpdateMember добавляет нового участника или обновляет существующего
func AddOrUpdateMember(name string) (int64, error) {
	// Проверяем, существует ли участник с таким именем
	var memberID int64
	err := db.QueryRow("SELECT id FROM members WHERE name = ?", name).Scan(&memberID)
	if err != nil {
		if err == sql.ErrNoRows {
			// Участник не найден, добавляем нового
			result, err := db.Exec("INSERT INTO members (name) VALUES (?)", name)
			if err != nil {
				return 0, err
			}
			memberID, err = result.LastInsertId()
			if err != nil {
				return 0, err
			}
		} else {
			// Произошла другая ошибка при запросе
			return 0, err
		}
	}
	// Если участник уже существует, memberID будет содержать его ID,
	// и нового участника добавлять не нужно
	return memberID, nil
}

// AddContribution добавляет взнос в таблицу contributions
func AddContribution(memberID int64, amount float64, date string, paymentMonth string) error {
	_, err := db.Exec("INSERT INTO contributions (member_id, amount, date, payment_month) VALUES (?, ?, ?, ?)", memberID, amount, date, paymentMonth)
	return err
}

// GetContributions возвращает список всех взносов, сгруппированных по участникам
func GetContributions() ([]domain.Member, error) {
	rows, err := db.Query(`
        SELECT m.name, c.payment_month, c.amount, m.start_date
        FROM members m
        LEFT JOIN contributions c ON m.id = c.member_id
        ORDER BY m.name, c.payment_month
    `)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	membersMap := make(map[string]*domain.Member)
	for rows.Next() {
		var name, month string
		var amount float64
		var startDate *string
		if err := rows.Scan(&name, &month, &amount, &startDate); err != nil {
			return nil, err
		}

		if _, ok := membersMap[name]; !ok {
			membersMap[name] = &domain.Member{Name: name, Months: []string{}, StartDate: *startDate}
		}
		membersMap[name].Months = append(membersMap[name].Months, month)
		// Сумму взносов больше не храним в структуре Member, так как она вычисляется
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	var members []domain.Member
	for _, member := range membersMap {
		members = append(members, *member)
	}

	return members, nil
}

// GetMember ищет участника по имени и возвращает его данные
func GetMember(name string) (*domain.Member, error) {
	member := &domain.Member{}
	err := db.QueryRow("SELECT id, name, start_date FROM members WHERE name = ?", name).Scan(&member.ID, &member.Name, &member.StartDate)
	if err != nil {
		return nil, err
	}
	return member, nil
}

// GetDebts возвращает список долгов
func GetDebts() ([]domain.Member, error) {
	rows, err := db.Query(`
        SELECT m.name, GROUP_CONCAT(DISTINCT c.payment_month), m.start_date
        FROM members m
        LEFT JOIN contributions c ON m.id = c.member_id
        GROUP BY m.name
        ORDER BY m.name
    `)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var members []domain.Member
	for rows.Next() {
		var member domain.Member
		var monthsPaidStr, startDate *string

		if err := rows.Scan(&member.Name, &monthsPaidStr, &startDate); err != nil {
			return nil, err
		}

		// Устанавливаем начальную дату, если она есть
		if startDate != nil {
			member.StartDate = *startDate
		}

		// Заполняем срез Months оплаченными месяцами
		member.Months = []string{}
		if monthsPaidStr != nil {
			member.Months = splitMonths(*monthsPaidStr)
		}

		// Рассчитываем количество неоплаченных месяцев
		if member.StartDate == "" {
			member.Debt = 0 // Если начальная дата не установлена, считаем, что долга нет
		} else {
			member.Debt = float64(calculateUnpaidMonths(member.StartDate, member.Months))
		}

		members = append(members, member)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return members, nil
}

// SetMemberStartDate устанавливает начальную дату для участника
func SetMemberStartDate(name string, startDate string) error {
	_, err := db.Exec("UPDATE members SET start_date = ? WHERE name = ?", startDate, name)
	return err
}

// GetMemberStartDate получает начальную дату для участника
func GetMemberStartDate(name string) (string, error) {
	var startDate string
	err := db.QueryRow("SELECT start_date FROM members WHERE name = ?", name).Scan(&startDate)
	if err != nil {
		return "", err
	}
	return startDate, nil
}

// Вспомогательная функция для расчета количества неоплаченных месяцев
func calculateUnpaidMonths(startDateStr string, paidMonths []string) int {
	startDate, err := time.Parse("2006-01-02", startDateStr)
	if err != nil {
		log.Printf("Error parsing start date: %v", err)
		return 0
	}

	// Приводим оплаченные месяцы к формату time.Time для упрощения сравнения
	paidMonthsTime := make(map[time.Time]bool)
	for _, monthStr := range paidMonths {
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

	return unpaidMonths
}

// Вспомогательная функция для разделения строки с месяцами на срез строк
func splitMonths(monthsStr string) []string {
	if monthsStr == "" {
		return []string{}
	}
	return strings.Split(monthsStr, ",")
}