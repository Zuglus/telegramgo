package domain

import "time"

type Member struct {
	ID           int64
	Name         string
	Contribution float64
	Debt         float64
	Months       []string
}

type Contribution struct {
	ID          int64
	MemberID    int64
	Amount      float64
	Date        time.Time
	PaymentMonth string
}