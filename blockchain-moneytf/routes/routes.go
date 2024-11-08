package routes

import (
	"github.com/Jonathan1366/blockchain-money-transfer/handlers"
	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App, authhandler *handlers.AuthHandlers) {
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString(fiber.ErrForbidden.Message)
	})
	api:=app.Group("api")
	api.Post("/transaction", authhandler.CreateTransactionHandler)
	api.Get("/transaction/:id", authhandler.GetTransactionHandler)
	api.Get("/transactions", authhandler.GetAllTransactionHandler)
	api.Put("/transaction/:id", authhandler.UpdateTransactionsHandler)
	api.Delete("/transaction/:id", authhandler.DeleteTransactionHandler)
	api.Get("/blocks", authhandler.GetAllBlocksHandler)
}