package db

import (
	"context"
	"log"

	"github.com/jackc/pgx/v4/pgxpool"
)

func Connect() *pgxpool.Pool {

	databaseURL := "postgresql://postgres.ufzhvxdimzrqjuxlvdvx:[YOUR-PASSWORD]@aws-0-ap-southeast-1.pooler.supabase.com:6543/postgres"
	if databaseURL == "" {
		log.Fatal("DATABASE_URL is not set in the environment")
	}

	// Konfigurasi pool koneksi
	config, err := pgxpool.ParseConfig(databaseURL)
	if err != nil {
		log.Fatalf("Unable to parse DATABASE_URL: %v", err)
	}

	// Buat pool koneksi
	conn, err := pgxpool.ConnectConfig(context.Background(), config)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v", err)
	}

	return conn
}


