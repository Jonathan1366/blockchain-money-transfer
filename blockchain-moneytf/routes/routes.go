package routes

import (
	"github.com/Jonathan1366/blockchain-money-transfer/handlers"
	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App, authhandler *handlers.AuthHandlers) {
	app.Get("/home", func(c *fiber.Ctx) error {
		return c.SendString(fiber.ErrForbidden.Message)
	})
	api:=app.Group("api")
	api.Post("/transaction", authhandler.CreateTransactionHandler)
	api.Get("/transactions", authhandler.GetAllTransactionHandler)
	api.Get("/users/:id", authhandler.GetUserByIDHandler)

	api.Put("/transaction/:id", authhandler.UpdateTransactionsHandler)
	api.Delete("/transaction/:id", authhandler.DeleteTransactionHandler)

	//endpoint block
	api.Get("/blocks", authhandler.GetAllBlocksHandler)
	
	api.Get("/mempool", authhandler.GetMempoolHandler)

	//endpoint mempool
	api.Get("/mine", authhandler.MinePendingTransactionsHandler)
	
}