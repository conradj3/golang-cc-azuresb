package main

import (
	"context"
	"log"
	"os"

	servicebus "github.com/Azure/azure-service-bus-go"
	"github.com/joho/godotenv"
)

var (
	connStr   string
	queueName string
)

func main() {
	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	connStr = os.Getenv("SERVICEBUS_CONNECTION_STRING")
	queueName = os.Getenv("QUEUE_NAME")

	// Create a new service bus namespace
	ns, err := servicebus.NewNamespace(servicebus.NamespaceWithConnectionString(connStr))
	if err != nil {
		log.Fatal("Error creating service bus namespace:", err)
	}

	// Get a handle to the queue
	q, err := ns.NewQueue(queueName)
	if err != nil {
		log.Fatal("Error creating new queue:", err)
	}

	// Receive a message from the queue
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err = q.ReceiveOne(ctx, servicebus.HandlerFunc(func(ctx context.Context, msg *servicebus.Message) error {
		log.Println("Received message:", string(msg.Data))
		return msg.Complete(ctx)
	}))

	if err != nil {
		log.Fatal("Error receiving message:", err)
	}

	log.Println("Successfully received and processed message.")
}
