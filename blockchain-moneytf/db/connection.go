package db

import (
	"context"
	"os"

	"github.com/gofiber/fiber/v2/log"
	"github.com/jackc/pgx/v4"
	"github.com/joho/godotenv"
)

func Connect() *pgx.Conn {
	err:=godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading.env file")
	}
	databaseURL:=os.Getenv("DATABASE_URL")
	if databaseURL==""{
		log.Fatal("DATABASE_URL is not set in the environment")
	}
	conn, err:=pgx.Connect(context.Background(), databaseURL)
	if err != nil {
		log.Fatalf("Unable to connect database")
	}
	return conn
}


