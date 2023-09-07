package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	_ "github.com/lib/pq"

	"time"

	"github.com/joho/godotenv"
)

type DBInfo struct {
	address  string
	user     string
	password string
	name     string
}

func main() {
	fmt.Println("Helloes")

	dbInfo, err := handleEnvVariables()
	if err != nil {
		log.Fatal(err.Error())
	}

	db, err := openDB(dbInfo)
	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	var version string
	if err := db.QueryRow("select version()").Scan(&version); err != nil {
		panic(err)
	}

	log.Printf("version - %s", version)

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Heelo :D"))
	})
	server := http.Server{
		Addr:    ":4000",
		Handler: mux,
	}

	if err := server.ListenAndServe(); err != nil {
		fmt.Printf("ooopps, %v", err)
	}
}

func openDB(dbInfo DBInfo) (*sql.DB, error) {

	dsn := fmt.Sprintf("postgres://%s:%s@%s/%s", dbInfo.user, dbInfo.password, dbInfo.address, dbInfo.name)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = db.PingContext(ctx)
	if err != nil {
		return nil, err
	}
	log.Print("Connected to DB :)")

	return db, nil
}

func handleEnvVariables() (DBInfo, error) {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	dbAddress := os.Getenv("DB_ADDRESS")
	if dbAddress == "" {
		log.Fatal("Missing DB_ADDRESS from .env")
	}
	dbName := os.Getenv("DB_NAME")
	if dbName == "" {
		log.Fatal("Missing DB_NAME from .env")
	}
	dbUser := os.Getenv("DB_USER")
	if dbUser == "" {
		log.Fatal("Missing DB_USER from .env")
	}
	dbPassword := os.Getenv("DB_PASSWORD")
	if dbPassword == "" {
		log.Fatal("Missing DB_PASSWORD from .env")
	}

	return DBInfo{address: dbAddress, user: dbUser, password: dbPassword, name: dbName}, err
}
