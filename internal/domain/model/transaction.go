package model

import "time"

type Transaction struct {
	ID         int
	FromUserId int
	ToUserId   int
	Amount     int
	Time       time.Time
}
