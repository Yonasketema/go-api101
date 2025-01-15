package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

var db *sql.DB

func main() {
	err := godotenv.Load()

	var PORT = ":8080"
	log.Printf("server start %s", PORT)
	var DATABASE_URL = os.Getenv("DATABASE_URL")
	db, err = openDB(DATABASE_URL)

	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	fmt.Printf("database connection successfully")
	err = http.ListenAndServe(PORT, routes())
	log.Fatal(err)
	os.Exit(1)
}

func openDB(dbUrl string) (*sql.DB, error) {
	db, err := sql.Open("postgres", dbUrl)
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	defer cancel()
	err = db.PingContext(ctx)

	if err != nil {
		return nil, err
	}

	return db, nil

}
