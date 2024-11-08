package main

import (
	"log"
	"os"

	"github.com/Jonathan1366/blockchain-money-transfer/db"
	"github.com/Jonathan1366/blockchain-money-transfer/handlers"
	"github.com/Jonathan1366/blockchain-money-transfer/routes"
	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()
	database:=db.Connect()
	authhandler:=handlers.InitialTransaction(database)
	routes.SetupRoutes(app, authhandler)
	// Setup Routes
	port:=os.Getenv("PORT")
	if port == ""{
		port ="3000"
	}
	log.Fatal(app.Listen(":"+port ))
}