package models

type User struct {
	ID         int     `json:"id"`
	Name       string  `json:"name"`
	Balance    float64 `json:"balance"`
	PublicKey  string  `json:"public_key"`
	PrivateKey string  `json:"-"`
}

type Transaction struct {
	ID              int     `json:"id"`
	SenderID        int     `json:"sender_id"`
	ReceiverID      int     `json:"receiver_id"`
	Amount          float64 `json:"amount"`
	Signature       string  `json:"signature"`
	TransactionHash string  `json:"transaction_hash"`
	Waktu           string  `json:"waktu"`
}

type Block struct {
	Id           int           `json:"id"`
	PreviousHash string        `json:"previous_hash"`
	Hash         string        `json:"hash"`
	Nonce        int           `json:"nonce"`
	Timestamp    string        `json:"timestamp"`
	Transactions  []Transaction `json:"transactions"`
}
