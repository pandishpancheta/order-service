package db

import (
	"database/sql"
	"fmt"
	"order-service/pkg/config"
	"strconv"

	_ "github.com/lib/pq"
)

func Init(cfg *config.Config) *sql.DB {
	port, err := strconv.Atoi(cfg.DB_PORT)
	if err != nil {
		panic(err)
	}

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable", cfg.DB_HOST,
		port, cfg.DB_USER, cfg.DB_PASS, cfg.DB_NAME)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	fmt.Println("Successfully connected!")
	return db
}

func InitTable(db *sql.DB) {
	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS orders (
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		id UUID PRIMARY KEY,
		user_id UUID NOT NULL,
		listing_id UUID NOT NULL,
		status VARCHAR(255) NOT NULL DEFAULT 'pending',
		token_uri VARCHAR(255) NOT NULL,
		FOREIGN KEY (user_id) REFERENCES users(id),
		FOREIGN KEY (listing_id) REFERENCES listings(id)
	)`)
	if err != nil {
		panic(err)
	}

	fmt.Println("Successfully created table!")
}
