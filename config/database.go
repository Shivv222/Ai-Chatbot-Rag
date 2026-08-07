package config

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func ConnectDB() {

	host := "localhost"
	port := 5432
	user := "postgres"
	password := "MySecure@123"
	dbname := "ai-chatbot"

	psqlInfo := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname,
	)

	var err error

	DB, err = sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatal("Error opening database:", err)
	}

	err = DB.Ping()
	if err != nil {
		log.Fatal("Error connecting to database:", err)
	}

	log.Println("✅ Database Connected Successfully")
}