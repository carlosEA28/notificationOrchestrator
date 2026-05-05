package main

import (
	"context"
	"log"
	"os"

	"github.com/carlosEA28/channelWorkers/internal/lib"
	"github.com/carlosEA28/channelWorkers/internal/web/server"
	"github.com/joho/godotenv"
)

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func main() {
	log.Println("Starting channel workers service")

	ctx := context.Background()

	err := godotenv.Load()
	if err != nil {
		log.Println(".env file not found")
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

	err = rabbitmq.Consume(ctx, "notification.delivery.queue", "notification.delivery.v1", "notificationDelivery")
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
