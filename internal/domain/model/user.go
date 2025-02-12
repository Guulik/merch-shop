package model

type User struct {
	ID        int            `json:"id"`
	Username  string         `json:"username"`
	Password  string         `json:"-"`
	Coins     int            `json:"coins"`
	Inventory map[string]int `json:"inventory"`
}
