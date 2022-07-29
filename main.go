package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/google/uuid"
	kafka "github.com/segmentio/kafka-go"
)

type Account struct {
	Account  int    `json:"account"`
	Username string `json:"username"`
	Priority string `json:"priority"`
}

type Order struct {
	OrderID         string    `json:"orderId"`
	TransactionType int       `json:"transactionType"`
	From            Account   `json:"from"`
	To              Account   `json:"to"`
	Amount          int       `json:"amount"`
	Fee             int       `json:"fee"`
	Status          string    `json:"status"`
	CreatedAt       time.Time `json:"createdAt"`
	Currency        string    `json:"currency"`
	PutProceeds     string    `json:"putProceeds"`
	PayWith         string    `json:"payWith"`
}

func newOrder() Order {
	return Order{
		OrderID:         uuid.NewString(),
		TransactionType: 1,
		From: Account{
			Account:  1234,
			Username: "Billy Joel",
			Priority: "high",
		},
		To:          Account{},
		Amount:      1000000,
		Fee:         10000,
		Status:      "PENDING",
		CreatedAt:   time.Now(),
		Currency:    "USD",
		PutProceeds: "12345",
		PayWith:     "",
	}
}
func trigger(kafkaWriter *kafka.Writer) func(http.ResponseWriter, *http.Request) {
	return http.HandlerFunc(func(wrt http.ResponseWriter, req *http.Request) {
		order := newOrder()
		messageBody, err := json.Marshal(order)
		if err != nil {
			log.Fatalln(err)
		}
		msg := kafka.Message{
			Key:   []byte(fmt.Sprintf("address-%s", req.RemoteAddr)),
			Value: messageBody,
		}

		err = kafkaWriter.WriteMessages(req.Context(), msg)
		if err != nil {
			wrt.Write([]byte(err.Error()))
			log.Fatalln(err)
		}

		wrt.Write([]byte(fmt.Sprintf("New order posted to Kafka : %v", order)))
	})
}

func main() {
	kafkaWriter := &kafka.Writer{
		Addr:     kafka.TCP(os.Getenv("KAFKA_ADDRESS")),
		Topic:    os.Getenv("TOPIC_NAME"),
		Balancer: &kafka.LeastBytes{},
	}
	defer kafkaWriter.Close()

	http.HandleFunc("/", trigger(kafkaWriter))

	log.Println("Listening on http://localhost:8081")
	log.Fatal(http.ListenAndServe(":8081", nil))
}
