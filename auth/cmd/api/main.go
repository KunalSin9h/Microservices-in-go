package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"auth/data"

	_ "github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
)

const PORT = "5002"

var count int64

type Config struct {
	DB     *sql.DB
	Models data.Models
}

func main() {

	conn := connectDB()

	if conn == nil {
		panic("Can't connect to postgres database")
	}

	app := Config{
		DB:     conn,
		Models: data.New(conn),
	}

	server := &http.Server{
		Addr:    fmt.Sprintf(":%s", PORT),
		Handler: app.routes(),
	}

	log.Println("Starting Authentication server at port:", PORT)
	log.Fatal(server.ListenAndServe())
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("pgx", dsn)
	if err != nil || db.Ping() != nil {
		return nil, err
	}

	return db, nil
}

func connectDB() *sql.DB {
	dsn := os.Getenv("DSN")
	for {
		connection, err := openDB(dsn)
		if err != nil {
			log.Println("Trying to connect...")
			count++
		} else {
			log.Println("Connected to postgres database!")
			return connection
		}

		if count > 10 {
			log.Println("Can't connect to postgres db even after 10 tries")
			return nil
		}

		time.Sleep(500 * time.Millisecond)
		continue
	}
}
