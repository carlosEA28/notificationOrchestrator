package main

import (
	"log"
	"os"

	"github.com/carlosEA28/notificationOrchestrator/internal/web/server"
	"github.com/joho/godotenv"
)

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Arquivo .env não encontrado")
	}

	port := getEnv("HTTP_PORT", "3003")
	srv := server.NewServer(port)

	if err := srv.Start(); err != nil {
		log.Fatal("Error starting server: ", err)
	}

	select {}
}
