package handlers

import (
	"context"
	"fmt"
	"time"

	modelss "github.com/Jonathan1366/blockchain-money-transfer/models"
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

func (h*AuthHandlers) CreateTransactionHandler(c *fiber.Ctx) error {
	transaction:= new(modelss.Transaction)
	if err:= c.BodyParser(transaction); err!=nil{
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error":"Cannot parse JSON",
		})
	}

	//buat hash dari transaksi misalnya sender_id + reciever_id+amount+ timestamp
	transaction.Timestamp = time.Now().Format(time.RFC3339)
	transaction.TransactionHash= utils.GenerateHash(fmt.Sprintf("%d%d%f%s", transaction.SenderID, transaction.ReceiverID, transaction.Amount, transaction.Timestamp) )

	//simpan transaksi ke db
	if err:= repositories.CreateTransaction(context.Background(), transaction); err!=nil{
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Failed to create transaction",
		})
	}
	
	return c.JSON(fiber.Map{"status":"success", "transaction":transaction})
}