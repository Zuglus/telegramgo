package domain

import "time"

type Member struct {
	ID        int64
	Name      string
	StartDate string
	Months    []string
	Debt      float64
}

// Contribution представляет собой отдельный взнос.
type Contribution struct {
	ID           int64
	MemberID     int64
	Amount       float64
	Date         time.Time
	PaymentMonth string
}