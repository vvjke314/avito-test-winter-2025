package ds

import "github.com/google/uuid"

type Merch struct {
	Id    uuid.UUID `db:"id" json:"id"`
	Name  string    `db:"name" json:"name"`
	Price int       `db:"price" json:"price"`
}
