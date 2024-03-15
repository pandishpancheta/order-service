package db

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"order-service/pkg/config"
	"strconv"
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
