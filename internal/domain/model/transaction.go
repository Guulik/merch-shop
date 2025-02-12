package model

import "time"

type Transaction struct {
	ID           int
	FromUserId   int    `db:"fromUserId"`
	FromUsername string `db:"fromUsername"`
	ToUserId     int    `db:"toUserId"`
	ToUsername   string `db:"toUsername"`
	Amount       int    `db:"amount"`
	Time         time.Time
}

type CoinHistory struct {
	Received []Received `json:"received"`
	Sent     []Sent     `json:"sent"`
}

type Received struct {
	FromUser string `json:"fromUser"`
	Amount   int    `json:"amount"`
}
type Sent struct {
	ToUser string `json:"toUser"`
	Amount int    `json:"amount"`
}
