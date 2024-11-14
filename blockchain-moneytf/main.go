package main

import (
	"log"
	"os"

	"github.com/Jonathan1366/blockchain-money-transfer/db"
	"github.com/Jonathan1366/blockchain-money-transfer/handlers"
	"github.com/Jonathan1366/blockchain-money-transfer/routes"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func main() {
	app := fiber.New()
	app.Use(cors.New(cors.Config{
		AllowOrigins: "http://localhost:5173", // Specify the allowed origin
		AllowHeaders: "Origin, Content-Type, Accept",
	}))
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