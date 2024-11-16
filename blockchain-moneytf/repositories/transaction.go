package repositories

import (
	"context"
	"log"

	"github.com/Jonathan1366/blockchain-money-transfer/db"
	"github.com/Jonathan1366/blockchain-money-transfer/models"
	"github.com/jackc/pgx/v4/pgxpool"
)


func CreateTransaction(ctx context.Context, db *pgxpool.Pool, transaction *models.Transaction) error  {
	err := db.QueryRow(ctx,"INSERT INTO transaction (sender_id, receiver_id, amount, signature, transaction_hash, waktu) VALUES($1, $2, $3, $4, $5, $6) RETURNING id", transaction.SenderID, transaction.ReceiverID, transaction.Amount, 
	transaction.Signature, transaction.TransactionHash, transaction.Waktu).Scan(&transaction.ID)
	if err != nil {
		log.Printf("Failed to create transaction: %v", err)
	}	
	return err
}

// func CreateBlock(ctx context.Context, db *pgxpool.Pool , block *models.Block) error {
// 	transactionJSON, err:= json.Marshal(block.Transactions)
// 	if err != nil {
// 		log.Printf("failed to marshal transactions: %v", err)
// 		return err
// 	}
// 	err = db.QueryRow(ctx, "INSERT INTO blocks (hash, previous_hash, nonce, timestamp, transactions) VALUES ($1, $2, $3, $4, $5) RETURNING id", block.Hash, block.PreviousHash, block.Nonce, block.Timestamp, transactionJSON).Scan(&block.Id)
// 	if err != nil {
// 		log.Printf("Failed to create new block: %v", err)
// 	}
// 	return err
	
// }

func CreateBlock(ctx context.Context, db *pgxpool.Pool, block *models.Block) error {
	query := `INSERT INTO blocks (id, previous_hash, hash, timestamp, nonce) VALUES ($1, $2, $3, $4, $5)`
	_, err := db.Exec(ctx, query, block.Id, block.PreviousHash, block.Hash, block.Timestamp, block.Nonce)
	return err
}


func GetUserByID(ctx context.Context, db*pgxpool.Pool, id int) (*models.User, error) {
	user := &models.User{}
	err := db.QueryRow(ctx, `SELECT id, name, balance, public_key, private_key FROM users WHERE id = $1`,
		id).Scan(&user.ID, &user.Name, &user.Balance, &user.PublicKey, &user.PrivateKey)
	if err != nil {
		return nil, err
	}
	return user, nil
}

// func GetlastBlock(ctx context.Context, db*pgxpool.Pool) (*models.Block, error)  {
// 	block:=&models.Block{}
// 	var transactionsJSON []byte

// 	err:= db.QueryRow(ctx, "SELECT id, previous_hash, hash, nonce, timestamp, transactions from blocks  WHERE id = $1", block.Id).Scan(
// 		&block.Id,
// 		&block.PreviousHash,
// 		&block.Hash,
// 		&block.Nonce,
// 		&block.Timestamp,
// 		&block.Transactions,
// 	)
// 		if err != nil {
// 			return nil, err
// 		}
// 		err = json.Unmarshal(transactionsJSON, &block.Transactions)
// 		if err != nil {
// 			return nil, err
// 		}
// 		return block, nil

// }

func GetlastBlock(ctx context.Context, db *pgxpool.Pool) (*models.Block, error) {
	query := `SELECT id, previous_hash, hash, timestamp, nonce FROM blocks ORDER BY id DESC LIMIT 1`
	row := db.QueryRow(ctx, query)

	var block models.Block
	if err := row.Scan(&block.Id, &block.PreviousHash, &block.Hash, &block.Timestamp, &block.Nonce); err != nil {
			return nil, err
	}
	return &block, nil
}


func UpdateBalance(ctx context.Context, db*pgxpool.Pool, userID int, balance float64) error {
	_, err := db.Exec(ctx, `UPDATE users SET balance = $1 WHERE id = $2`, balance, userID)
	return err
}


func GetTransactionByID(ctx context.Context, id int)(*models.Transaction, error)  {
	conn:=db.Connect()
	transaction:= &models.Transaction{}

	err:= conn.QueryRow(ctx, "SELECT id, sender_id, receiver_id, amount, transaction_hash, waktu FROM transaction where id=$1", id).Scan(
		&transaction.ID,
		&transaction.SenderID,
		&transaction.ReceiverID,
		&transaction.Amount,
		&transaction.TransactionHash,
		&transaction.Waktu,
	)
	if err != nil {
		return nil, err
	}
	return transaction, nil
}


func GetAllTransactions(ctx context.Context)([]models.Transaction, error)  {
	conn:= db.Connect()
	rows, err:= conn.Query(ctx,"SELECT id, sender_id, receiver_id, amount, transaction_hash, waktu from transaction")
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
			&transaction.Waktu,
		)
		if err != nil {
			return nil, err
		}
		transactions=append(transactions, transaction )
	}
	return transactions, nil
}

func UpdateTransaction(ctx context.Context, SenderID int, transaction *models.Transaction) error  {
	conn:= db.Connect()
	tx, err:= conn.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	// Update saldo sender
	_, err = tx.Exec(ctx, `UPDATE users SET balance = balance - $1 WHERE id = $2`, transaction.Amount, SenderID)
	if err != nil {
			return err
	}

	// Update saldo receiver
	_, err = tx.Exec(ctx, `UPDATE users SET balance = balance + $1 WHERE id = $2`, transaction.Amount, transaction.ReceiverID)
	if err != nil {
			return err
	}
	// Commit transaksi jika semua berhasil
	if err := tx.Commit(ctx); err != nil {
		return err
}

return nil
}

//fungsi utk menghapus transaksi berdasarkan id 
func DeleteTransaction(ctx context.Context, id int) error  {
	conn:=db.Connect()
	_, err:= conn.Exec(ctx, "DELETE from transaction where id=$1", id)
	return err
}



// func GetAllBlocks(ctx context.Context) ([]models.Block, error)  {
// 	conn:=db.Connect()
// 	rows, err:= conn.Query(ctx, "SELECT id, transaction_id, previous_hash, hash, timestamp from blocks")
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer rows.Close()
	
// 	blocks := []models.Block{}
// 	for rows.Next() {
// 		var block models.Block
// 		err:=rows.Scan(
// 			&block.Id,
// 			&block.Transaction,
// 			&block.PreviousHash,
// 			&block.Hash,
// 			&block.Timestamp,
// 		)
// 		if err != nil {
// 			return nil, err
// 		}
// 		blocks = append(blocks, block)
// 	}
// 	return blocks, nil
// }