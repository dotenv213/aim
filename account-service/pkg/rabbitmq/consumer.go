package rabbitmq

import (
	"encoding/json"
	"log"

	"github.com/dotenv213/aim/account-service/internal/domain"
	amqp "github.com/rabbitmq/amqp091-go"
)

type BalanceUpdateEvent struct {
	BankID uint    `json:"bank_id"`
	Amount float64 `json:"amount"`
	Type   string  `json:"type"`
}

type Consumer struct {
	conn    *amqp.Connection
	service domain.BankService
}

func NewConsumer(url string, service domain.BankService) (*Consumer, error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, err
	}
	return &Consumer{conn: conn, service: service}, nil
}

func (c *Consumer) Start() {
	ch, err := c.conn.Channel()
	if err != nil {
		log.Fatalf("Failed to open channel: %v", err)
	}
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"balance_updates", // queue name == sender
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("Failed to declare queue: %v", err)
	}

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer tag
		true,   // auto-ack delete the message when consumed
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	if err != nil {
		log.Fatalf("Failed to register consumer: %v", err)
	}

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			log.Printf("Received a message: %s", d.Body)

			var event BalanceUpdateEvent
			if err := json.Unmarshal(d.Body, &event); err != nil {
				log.Printf("Error parsing message: %v", err)
				continue
			}

			err := c.service.UpdateBalance(event.BankID, event.Amount, event.Type)
			if err != nil {
				log.Printf("Error updating balance: %v", err)
			} else {
				log.Printf("Balance updated for BankID: %d", event.BankID)
			}
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}
