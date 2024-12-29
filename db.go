package main

import (
	"database/sql"
	"log"

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
func addContribution(memberID int64, amount float64, date string) error {
	_, err := db.Exec("INSERT INTO contributions (member_id, amount, date) VALUES (?, ?, ?)", memberID, amount, date)
	return err
}