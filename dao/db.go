package dao

import (
	"database/sql"
	"fmt"
	"log"
	"sync"

	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "shannee_ahirwar_ftc"
	password = "postgres"
	dbname   = "fulfillment"
)

var (
	db     *sql.DB
	dbOnce sync.Once
)

func InitDB() (*sql.DB, error) {
	var err error
	dbOnce.Do(func() {
		psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
			"password=%s dbname=%s sslmode=disable",
			host, port, user, password, dbname)

		db, err = sql.Open("postgres", psqlInfo)
		if err != nil {
			log.Fatalf("Failed to connect to database: %v", err)
			return
		}

		if err = db.Ping(); err != nil {
			log.Fatalf("Failed to ping database: %v", err)
			return
		}

		log.Println("Database connected successfully!")
	})

	return db, err
}

func GetDB() *sql.DB {
	return db
}
