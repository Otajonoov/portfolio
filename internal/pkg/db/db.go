package db

import (
	"context"
	"log"
	"portfolio/internal/pkg/config"

	"github.com/jackc/pgx/v5/pgxpool"
)

func ConnectToDatabase() *pgxpool.Pool {
	cfg := config.Load(".")

	connDB, err := pgxpool.New(context.Background(), cfg.Postgres.DatabaseURL)
	if err != nil {
		log.Fatal(err)
	}
	
	return connDB
}
