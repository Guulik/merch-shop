package model

type UserInfo struct {
	Coins       int         `json:"coins"`
	Inventory   []Item      `json:"inventory"`
	CoinHistory CoinHistory `json:"coinHistory"`
}

type Item struct {
	Type     string `json:"type" db:"item"`
	Quantity int    `json:"quantity" db:"quantity"`
}
