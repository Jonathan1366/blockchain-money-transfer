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

	// Mining a new block
	lastBlock, err := repositories.GetlastBlock(context.Background())
	if err != nil {
		log.Printf("Failed to retrieve last block: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve last block",
		})
	}
	//buat block baru
	newBlock := models.Block{
		TransactionId: transaction.ID,
		PreviousHash:  lastBlock.Hash,
		Timestamp:     transaction.Waktu,
	}

	utils.MineBlock(&newBlock, 4) // Mining with difficulty 4
	if err := repositories.CreateBlock(context.Background(), &newBlock); err != nil {
		log.Printf("Failed to create block: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create block",
		})
	}
	return c.JSON(fiber.Map{
		"status":"success",
		"transaction":transaction,
		"block": newBlock,
		"public_key":sender.PublicKey,
	})
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