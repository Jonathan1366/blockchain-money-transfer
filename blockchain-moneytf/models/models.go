package models

type User struct {
	ID      int    `json:"id"`
	Name    string `json:"name"`
	Balance string `json:"balance"`
}

type Transaction struct {
	ID              int     `json:"id"`
	SenderID        int     `json:"sender_id"`
	ReceiverID      int     `json:"receiver_id"`
	Amount          float64 `json:"amount"`
	TransactionHash string  `json:"transaction_hash"`
	Timestamp       string  `json:"timestamp"`
}

type Block struct {
	Id            int    `json:"id"`
	TransactionId int    `json:"transaction_id"`
	PreviousHash  string `json:"previous_hash"`
	Hash          string `json:"hash"`
	Timestamp     string `json:"timestamp"`
}
