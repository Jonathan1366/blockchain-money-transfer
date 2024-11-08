package repositories

import (
	"context"

	"github.com/Jonathan1366/blockchain-money-transfer/db"
	"github.com/Jonathan1366/blockchain-money-transfer/models"
)

func CreateTransaction(ctx context.Context, transaction *models.Transaction) error  {
	conn:=db.Connect()
	_, err:=conn.Exec(ctx,"INSERT INTO transaction (sender_id, receiver_id, amount, transaction_hash, timestamp) VALUES($1,$2,$3,$4, $5)", transaction.SenderID, transaction.ReceiverID, transaction.Amount, transaction.TransactionHash, transaction.Timestamp)
	return err
}

func GetTransactionByID(ctx context.Context, id int)(*models.Transaction, error)  {
	conn:=db.Connect()
	transaction:= &models.Transaction{}

	err:= conn.QueryRow(ctx, "SELECT id, sender_id, reciever_id, amount, transaction_hash, timestamp FROM transaction where id=$1", id).Scan(
		&transaction.ID,
		&transaction.SenderID,
		&transaction.ReceiverID,
		&transaction.Amount,
		&transaction.TransactionHash,
		&transaction.Timestamp,
	)
	if err != nil {
		return nil, err
	}
	return transaction, nil
}


func GetAllTransactions(ctx context.Context)([]models.Transaction, error)  {
	conn:= db.Connect()
	rows, err:= conn.Query(ctx,"SELECT id, sender_id, reciever_id, amount, transaction_hash, timestamp from transaction")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	transactions := []models.Transaction{}
	
	for rows.Next(){
		var transaction models.Transaction
		err:= rows.Scan(
			&transaction.ID,
			&transaction.SenderID,
			&transaction.ReceiverID,
			&transaction.Amount,
			&transaction.TransactionHash,
			&transaction.Timestamp,
		)
		if err != nil {
			return nil, err
		}
		transactions=append(transactions, transaction )
	}
	return transactions, nil
}

func UpdateTransaction(ctx context.Context, id int, transaction *models.Transaction) error  {
	conn:= db.Connect()
	_, err:= conn.Exec(ctx, "UPDATE transaction SET sender_id=$1, receiver_id=$2, amount=$3, transaction_hash=$4, timestamp =$5 WHERE id=$6",
		transaction.SenderID,
		transaction.ReceiverID,
		transaction.Amount,
		transaction.TransactionHash,
		transaction.Timestamp,
		id,
	)
	return err
}

//fungsi utk menghapus transaksi berdasarkan id 
func DeleteTransaction(ctx context.Context, id int) error  {
	conn:=db.Connect()
	_, err:= conn.Exec(ctx, "DELETE from transaction where id=$1", id)
	return err
}

func GetlastBlock(ctx context.Context) (*models.Block, error)  {
	conn:=db.Connect()
	block:=&models.Block{}

	err:= conn.QueryRow(ctx, "SELECT id, transaction_id, previous_hash, hash, timestamp from blocks order by id desc limit 1").Scan(
		&block.Id,
		&block.TransactionId,
		&block.PreviousHash,
		&block.Hash,
		&block.Timestamp,)
		if err != nil {
			return nil, err
		}
		return block, nil
}

func CreateBlock(ctx context.Context, block *models.Block)error  {
	conn:= db.Connect()
	_, err:=conn.Exec(ctx, "INSERT INTO blocks (transaction_id, previous_hash, timestamp) VALUES ($1, $2, $3, $4)", block.TransactionId, block.PreviousHash, block.Hash, block.Timestamp)
	return err
}

func GetAllBlocks(ctx context.Context) ([]models.Block, error)  {
	conn:=db.Connect()
	rows, err:= conn.Query(ctx, "SELECT id, transaction_id, previous_hash, hash, timestamp from blocks")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	blocks := []models.Block{}
	for rows.Next() {
		var block models.Block
		err:=rows.Scan(
			&block.Id,
			&block.TransactionId,
			&block.PreviousHash,
			&block.Hash,
			&block.Timestamp,
		)
		if err != nil {
			return nil, err
		}
		blocks = append(blocks, block)
	}
	return blocks, nil
}