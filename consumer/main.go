package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

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
		log.Fatal("error loading .env file")
	}
	log.Println("loaded environment variables from .env file")

	connStr = os.Getenv("SERVICEBUS_CONNECTION_STRING")
	queueName = os.Getenv("QUEUE_NAME")

	// Connecting to the namespace with a connection string
	ns, err := servicebus.NewNamespace(servicebus.NamespaceWithConnectionString(connStr))
	if err != nil {
		log.Fatal("Error connecting service bus namespace:", err)
	}
	log.Println("conected service bus namespace")

	// Get a handle to the queue
	q, err := ns.NewQueue(queueName)
	if err != nil {
		log.Fatal("Error establishing queue:", err)
	}
	log.Println("recieved handle to queue")

	// Create a context that can be cancelled when a signal is received
	ctx, cancel := context.WithCancel(context.Background())

	// Handle SIGTERM signal
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGTERM)
	go func() {
		<-sigCh
		log.Println("received SIGTERM signal, cancelling context")
		cancel()
	}()

	// Receive message from the queue
	err = q.Receive(ctx, servicebus.HandlerFunc(func(ctx context.Context, msg *servicebus.Message) error {
		log.Printf("received message: %s\n", string(msg.Data))
		return msg.Complete(ctx)
	}))
	if err != nil {
		log.Fatal("Error receiving message from queue:", err)
	}
}
