package ds

type InfoResponse struct {
	Coins       int             `json:"coins"`
	Inventory   []ItemAmount    `json:"inventory"`
	CoinHistory TransferHistory `json:"coinHistory"`
}

type ItemAmount struct {
	Type     string `json:"type"`
	Quantity int    `json:"quantity"`
}

type ReceiveRecord struct {
	FromUser string `json:"fromUser"`
	Amount   int    `json:"amount"`
}

type SentRecord struct {
	ToUser string `json:"toUser"`
	Amount int    `json:"amount"`
}

type TransferHistory struct {
	Received []ReceiveRecord `json:"received"`
	Sent     []SentRecord    `json:"sent"`
}
