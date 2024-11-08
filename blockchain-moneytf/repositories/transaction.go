package repositories

import (
	"context"

	"github.com/Jonathan1366/blockchain-money-transfer/db"
	"github.com/Jonathan1366/blockchain-money-transfer/models"
)

func CreateTransaction(ctx context.Context, transaction *models.Transaction) error  {
	conn:=db.Connect()
	_, err:=conn.Exec(ctx,"INSERT INTO transactions (sender_id, receiver_id, amount, transaction_hash, timestamp) VALUES($1,$2,$3,$4, $5)", transaction.SenderID, transaction.ReceiverID, transaction.Amount, transaction.TransactionHash, transaction.Timestamp)
	return err
}