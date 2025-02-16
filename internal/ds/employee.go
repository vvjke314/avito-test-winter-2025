package ds

import "github.com/google/uuid"

type Employee struct {
	Id           uuid.UUID `db:"id" json:"id"`
	Username     string    `db:"username" json:"username"`
	PasswordHash string    `db:"password_hash" json:"-"`
	Coins        int       `db:"coins" json:"coins"`
}