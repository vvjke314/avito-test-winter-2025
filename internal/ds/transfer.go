package ds

import (
	"time"

	"github.com/google/uuid"
)

type Transfer struct {
	Id        uuid.UUID `db:"id" json:"id"`
	From      uuid.UUID `db:"from_emp_id" json:"from_emp_id"`
	To        uuid.UUID `db:"to_emp_id" json:"to_emp_id"`
	Amount    int       `db:"amount" json:"amount"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
}
