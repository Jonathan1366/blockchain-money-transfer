package handlers

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"time"

	// "github.com/Jonathan1366/blockchain-money-transfer/models"
	"github.com/Jonathan1366/blockchain-money-transfer/models"
	"github.com/Jonathan1366/blockchain-money-transfer/repositories"
	"github.com/Jonathan1366/blockchain-money-transfer/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/jackc/pgx/v5"
)

type AuthHandlers struct{
	DB *pgxpool.Pool
	DefaultQueryExecMode pgx.QueryExecMode
}

func InitialTransaction(db *pgxpool.Pool) *AuthHandlers{
	return &AuthHandlers{DB: db}
}

func (h *AuthHandlers) CreateTransactionHandler(c *fiber.Ctx) error {
	transaction:= new(models.Transaction)
	if err:= c.BodyParser(transaction); err!=nil{
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error":"Cannot parse JSON",
		})
	}

		sender, err := repositories.GetUserByID(context.Background(), transaction.SenderID)
		if err != nil {
			log.Printf("Sender not found: ID %d", transaction.SenderID)
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Sender not found"})
		}
		receiver, err := repositories.GetUserByID(context.Background(), transaction.ReceiverID)
		if err != nil {
				log.Printf("Receiver not found: ID %d", transaction.ReceiverID)
				return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Receiver not found"})
		}

		if sender.Balance < transaction.Amount {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
					"error": "saldo tidak cukup",
			})
	}	
	
	sender.Balance -= transaction.Amount

	receiver.Balance += transaction.Amount

	if err := repositories.UpdateBalance(context.Background(), sender.ID, sender.Balance); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"error": "Failed to update sender balance",
			})
	}	

	if err := repositories.UpdateBalance(context.Background(), receiver.ID, receiver.Balance); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to update receiver balance",
		})
}
	transaction.Signature = utils.SignTransaction(sender.PrivateKey, fmt.Sprintf("%d%d%f", transaction.SenderID, transaction.ReceiverID, transaction.Amount))
	//buat hash dari transaksi misalnya sender_id + reciever_id+amount+ timestamp
	transaction.Waktu = time.Now().Format(time.RFC3339)
	transaction.TransactionHash= utils.GenerateHash(fmt.Sprintf("%d%d%f%s", transaction.SenderID, transaction.ReceiverID, transaction.Amount, transaction.Waktu) )

	//simpan transaksi ke db
	if err:= repositories.CreateTransaction(context.Background(), transaction); err!=nil{
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Failed to create transaction",
		})
	}	

	// if err := repositories.UpdateTransaction(context.Background(), transaction.SenderID, transaction ); err!=nil{
	// 	return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
	// 		"error":"failed to update balances"})
	// }//buat block baru


	// Mining a new block
	lastBlock, err := repositories.GetlastBlock(context.Background())
	if err != nil {
		if err == pgx.ErrNoRows {
			lastBlock = &models.Block{
				Id: 0,
				PreviousHash: "0",
				Hash: "genesis_hash",
				Timestamp: time.Now().Format(time.RFC3339),
			}
		}
		} else{
			log.Printf("Failed to retrieve new block: %v", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error":"fail to retrieve blocks",
			})
		}
		newblock:= models.Block{
			TransactionId: transaction.ID,
			PreviousHash: lastBlock.Hash,
			Timestamp: lastBlock.Timestamp,	
		}
		utils.MineBlock(&newblock, 4)
		
		if err := repositories.CreateBlock(context.Background(), &newblock); err != nil {
			log.Printf("Failed to create block: %v", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"error": "Failed to create block",
			})
	}

	return c.JSON(fiber.Map{"status":"success", "transaction":transaction, "block": newblock, "public_key":sender.PublicKey})
	}

	

func (h *AuthHandlers) GetTransactionHandler(c *fiber.Ctx) error  {
	idStr := c.Params("id")

	//str ke int
	id, err:=strconv.Atoi(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":"invalid id format",
		})
	}
	transaction, err:= repositories.GetTransactionByID(context.Background(), id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error":"Transaction not found",
		})
	}
	return c.JSON(transaction)
};

//handler utk mengambil semua transaksi

func (h* AuthHandlers) GetAllTransactionHandler(c *fiber.Ctx) error  {
	transactions, err := repositories.GetAllTransactions(context.Background())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve transaction",
		})
	}
	return c.JSON(transactions)
}

//handler utk memperbarui transaksi berdasarkan id
func (h *AuthHandlers) UpdateTransactionsHandler(c * fiber.Ctx) error {
	idStr:=c.Params("id")
	id, err:=strconv.Atoi(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid id format",
		})
	}
	transaction:= new(models.Transaction)
	if err:= c.BodyParser(transaction); err!=nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":"cannot parse JSON",
		})
	}
	if err:= repositories.UpdateTransaction(context.Background(), id, transaction); err !=nil{
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":"failed to update transaction",
		})
		
	}
	return c.JSON(fiber.Map{
		"status": "success",
		"transaction": transaction,
	})
}

//handler untuk menghapus transaksi berdasarkan ID
func (h *AuthHandlers) DeleteTransactionHandler(c *fiber.Ctx) error  {
	idStr:= c.Params("id")
	id, err:= strconv.Atoi(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":"invalid id format",
		})
	}
	if err:= repositories.DeleteTransaction(context.Background(), id); err!=nil{
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":"failed to delete transaction",
		})
	}
	return c.JSON(fiber.Map{
		"status":"deleted",
	})
}


func (h *AuthHandlers)GetAllBlocksHandler(c *fiber.Ctx) error  {
	blocks, err:= repositories.GetAllBlocks(context.Background())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":"Failed to retrieve blocks",
		})
	}
	return c.JSON(blocks)
}