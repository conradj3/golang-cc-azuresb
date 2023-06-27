package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	servicebus "github.com/Azure/azure-service-bus-go"
	"github.com/joho/godotenv"
)

type MessageRequest struct {
	Count int `json:"count"`
}

var (
	connStr   string
	queueName string
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	connStr = os.Getenv("SERVICEBUS_CONNECTION_STRING")
	queueName = os.Getenv("QUEUE_NAME")

	http.HandleFunc("/createMessages", createMessagesHandler)
	http.HandleFunc("/clearMessages", clearMessagesHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func createMessagesHandler(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var req MessageRequest
	err := decoder.Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	client, err := servicebus.NewNamespace(servicebus.NamespaceWithConnectionString(connStr))
	if err != nil {
		http.Error(w, "Error creating service bus client", http.StatusInternalServerError)
		return
	}

	sender, err := client.NewQueue(queueName)
	if err != nil {
		http.Error(w, "Error creating queue sender", http.StatusInternalServerError)
		return
	}

	ctx := context.Background()

	for i := 0; i < req.Count; i++ {
		message := "Message " + strconv.Itoa(i+1)
		err = sender.Send(ctx, &servicebus.Message{
			Data: []byte(message),
		})
		if err != nil {
			http.Error(w, "Error sending message", http.StatusInternalServerError)
			return
		}
	}

	fmt.Fprintf(w, "Successfully sent %d messages to queue.\n", req.Count)
}

func clearMessagesHandler(w http.ResponseWriter, r *http.Request) {
	ns, err := servicebus.NewNamespace(servicebus.NamespaceWithConnectionString(connStr))
	if err != nil {
		http.Error(w, "Error creating service bus namespace", http.StatusInternalServerError)
		return
	}

	q, err := ns.NewQueue(queueName)
	if err != nil {
		http.Error(w, "Error creating new queue", http.StatusInternalServerError)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel() // Important to avoid leaks

	receiver, err := q.NewReceiver(ctx)
	if err != nil {
		http.Error(w, "Error creating new receiver", http.StatusInternalServerError)
		return
	}
	defer receiver.Close(ctx)

	for {
		// Try to receive a message
		if err := receiver.ReceiveOne(ctx, servicebus.HandlerFunc(func(ctx context.Context, msg *servicebus.Message) error {
			return msg.Complete(ctx)
		})); err != nil {
			// If no message is received within the timeout, break the loop
			if err == context.DeadlineExceeded {
				break
			}
			// If an error other than DeadlineExceeded occurs, return a 500 error
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	fmt.Fprintln(w, "Successfully cleared all messages from queue.")
}
