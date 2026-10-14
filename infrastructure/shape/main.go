package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/Space-Core-Dev/Mistaken.Recipe.Broker/declarations"
	"github.com/Space-Core-Dev/Mistaken.Recipe/application/services/helpers"
	"github.com/joho/godotenv"
	"github.com/rabbitmq/amqp091-go"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("environment variables loading failed %v", err)
	}

	brokerConnectionString := os.Getenv("RABBIT_CONNECTION_STRING")
	if brokerConnectionString == "" {
		log.Fatal("wrong environment variable RABBIT_CONNECTION_STRING")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	res := helpers.Retry(ctx, amqp091.Dial, brokerConnectionString, 5, 1, 1)
	if !res.IsSuccess() {
		log.Fatalf("event broker connection failed: %v", res.Error())
	}
	conn := res.Value()
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("broker channel failed: %v", err)
	}
	defer ch.Close()

	declaration := declarations.New(ch, "topic", queuePref, keyPref, names)
	err = declaration.Start()
	if err != nil {
		log.Fatalf("broker declarations failed: %v", err)
	}
}
