package main

import (
	"database/sql"
	"log"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

func initDB() {
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
            name TEXT UNIQUE
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

// Member представляет собой запись в таблице members
type Member struct {
	ID           int64
	Name         string
	Contribution float64
	Debt         float64
	Months       []string
}

// addOrUpdateMember добавляет нового участника или обновляет существующего
func addOrUpdateMember(name string) (int64, error) {
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

// addContribution добавляет взнос в таблицу contributions
func addContribution(memberID int64, amount float64, date string, paymentMonth string) error {
	_, err := db.Exec("INSERT INTO contributions (member_id, amount, date, payment_month) VALUES (?, ?, ?, ?)", memberID, amount, date, paymentMonth)
	return err
}

// getContributions возвращает список всех взносов
func getContributions() ([]Member, error) {
	rows, err := db.Query(`
		SELECT m.name, c.payment_month, c.amount
		FROM members m
		LEFT JOIN contributions c ON m.id = c.member_id
		ORDER BY m.name, c.payment_month
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var membersMap = make(map[string]*Member)
	for rows.Next() {
		var name, month string
		var amount float64
		if err := rows.Scan(&name, &month, &amount); err != nil {
			return nil, err
		}

		if _, ok := membersMap[name]; !ok {
			membersMap[name] = &Member{Name: name, Months: []string{}, Contribution: 0}
		}
		membersMap[name].Months = append(membersMap[name].Months, month)
		membersMap[name].Contribution += amount
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	var members []Member
	for _, member := range membersMap {
		members = append(members, *member)
	}

	return members, nil
}

// getDebts возвращает список долгов
func getDebts() ([]Member, error) {
	rows, err := db.Query(`
        SELECT m.name, GROUP_CONCAT(DISTINCT c.payment_month), COUNT(DISTINCT c.payment_month)
        FROM members m
        LEFT JOIN contributions c ON m.id = c.member_id
		WHERE c.payment_month IS NOT NULL
        GROUP BY m.name
        ORDER BY m.name
    `)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var members []Member
	for rows.Next() {
		var member Member
		var monthsPaidStr *string
		var monthsPaidCount int

		if err := rows.Scan(&member.Name, &monthsPaidStr, &monthsPaidCount); err != nil {
			return nil, err
		}

		// Заполняем срез Months оплаченными месяцами
		member.Months = []string{}
		if monthsPaidStr != nil {
			member.Months = splitMonths(*monthsPaidStr)
		}
		member.Debt = float64(monthsPaidCount)

		members = append(members, member)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return members, nil
}

// Вспомогательная функция для разделения строки с месяцами на срез строк
func splitMonths(monthsStr string) []string {
	if monthsStr == "" {
		return []string{}
	}
	return strings.Split(monthsStr, ",")
}