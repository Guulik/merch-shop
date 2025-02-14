package model

type UserAuth struct {
	Id         int    `json:"id" db:"id"`
	Username   string `json:"username" db:"username"`
	PasswordDb string `json:"password" db:"password_hash"`
}
