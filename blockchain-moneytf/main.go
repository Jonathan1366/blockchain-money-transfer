package main

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()
	// Setup Routes
	port:=os.Getenv("PORT")
	if port == ""{
		port ="3000"
	}
	log.Fatal(app.Listen(":"+port ))
}