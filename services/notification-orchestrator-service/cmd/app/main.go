package main

import (
	"context"
	"log"
	"os"

	"github.com/carlosEA28/notificationOrchestrator/internal/lib"
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
	log.Println("Starting notification orchestrator service")

	ctx := context.Background()

	err := godotenv.Load()
	if err != nil {
		log.Println("Arquivo .env não encontrado")
	}

	rabbitmqUrl := os.Getenv("RABBITMQ_URL")
	if rabbitmqUrl == "" {
		log.Fatal("RABBITMQ_URL is required")
	}
	log.Printf("RabbitMQ URL: %s", rabbitmqUrl)

	rabbitmq, err := lib.NewRabbitmqProvider()
	if err != nil {
		log.Fatal("Error connecting to RabbitMQ: ", err)
	}
	defer rabbitmq.Close()

	err = rabbitmq.Consume(ctx, "notification.requested", "notificationRequested", "notificationRequested")
	if err != nil {
		log.Fatal("Error starting consumer: ", err)
	}
	log.Println("RabbitMQ consumer started")

	port := getEnv("HTTP_PORT", "3003")
	srv := server.NewServer(port)

	go func() {
		if err := srv.Start(); err != nil {
			log.Fatal("Error starting server: ", err)
		}
	}()

	select {}
}
