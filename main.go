package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"

	"github.com/MrR0b0t1001/Chirpy/config"
	"github.com/MrR0b0t1001/Chirpy/internal/database"
	"github.com/MrR0b0t1001/Chirpy/server"
)

const address = ":8080"

func main() {
	if err := godotenv.Load(".env"); err != nil {
		log.Fatalf("Error loading .env file: %v\n", err)
	}

	dbURL := os.Getenv("DB_URL")
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Println("Error connecting to chirpy")
	}

	mux := http.NewServeMux()
	cnfg := &config.APIConfig{
		DB: database.New(db),
	}

	s := server.NewAPIServer(address, mux, cnfg)
	s.Run()
}
