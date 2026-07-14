package database

import (
	"context"
	"fmt"
	"log"

	"github.com/NickFinchD/chinese-learning-api/config"
	"github.com/jackc/pgx/v5"
)

func Connect(cfg *config.Config) *pgx.Conn {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.DB.Host,
		cfg.DB.Port,
		cfg.DB.User,
		cfg.DB.Password,
		cfg.DB.Name,
	)

	conn, err := pgx.Connect(context.Background(), dsn)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	log.Println("Connected to PostgreSQL")

	return conn
}
