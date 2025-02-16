package ds

type SendCoinRecord struct {
	ToUser string `json:"toUser"`
	Amount int    `json:"amount"`
}
