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

	db, err := sql.Open("postgres", os.Getenv("DB_URL"))
	if err != nil {
		log.Fatalf("Error connecting to chirpy: %v", err)
	}

	if err := db.Ping(); err != nil {
		log.Fatal("Chirpy DB appears to be down. Please start it and try again...")
	}

	mux := http.NewServeMux()
	cnfg := &config.APIConfig{
		DB:        database.New(db),
		JWTSecret: os.Getenv("JWT_SECRET"),
	}

	s := server.NewAPIServer(address, mux, cnfg)
	s.Run()
}
