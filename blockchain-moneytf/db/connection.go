package db

import (
	"context"
	"os"

	"github.com/gofiber/fiber/v2/log"
	"github.com/jackc/pgx/v4"
	"github.com/joho/godotenv"
)

func Connect() *pgx.Conn {
	if os.Getenv("RAILWAY_ENVIRONMENT")==""{
		err:=godotenv.Load()
		if err != nil {
			log.Fatalf("No .env file found, proceeding with environment variables")
		}
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


