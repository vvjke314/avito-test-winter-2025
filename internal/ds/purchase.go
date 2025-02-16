package ds

import (
	"time"

	"github.com/google/uuid"
)

type Purchase struct {
	Id         uuid.UUID `db:"id" json:"id"`
	EmployeeId uuid.UUID `db:"employee_id" json:"employee_id"`
	MerchId    uuid.UUID `db:"merch_id" json:"merch_id"`
	Amount     int       `db:"amount" json:"amount"`
	CreatedAt  time.Time `db:"created_at" json:"created_at"`
}
