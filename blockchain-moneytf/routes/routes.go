package routes

import (
	"github.com/Jonathan1366/blockchain-money-transfer/handlers"
	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App, authhandler *handlers.AuthHandlers) {
	api:=app.Group("/api")
	api.Post("/transaction", authhandler.CreateTransactionHandler)
	
}