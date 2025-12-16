package rabbitmq

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitMQProducer struct {
	conn    *amqp.Connection
	channel *amqp.Channel
}

func NewRabbitMQProducer(url string) (*RabbitMQProducer, error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to rabbitmq: %v", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("failed to open a channel: %v", err)
	}

	_, err = ch.QueueDeclare(
		"balance_updates", // queue name
		true,              // durable
		false,             // delete when unused
		false,             // exclusive
		false,             // no-wait
		nil,               // arguments
	)
	if err != nil {
		return nil, fmt.Errorf("failed to declare queue: %v", err)
	}

	return &RabbitMQProducer{
		conn:    conn,
		channel: ch,
	}, nil
}

type BalanceUpdateEvent struct {
	BankID uint    `json:"bank_id"`
	Amount float64 `json:"amount"`
	Type   string  `json:"type"` // "deposit" or "withdraw"
}

func (p *RabbitMQProducer) PublishBalanceUpdate(bankID uint, amount float64, trxType string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	event := BalanceUpdateEvent{
		BankID: bankID,
		Amount: amount,
		Type:   trxType,
	}

	body, err := json.Marshal(event)
	if err != nil {
		return err
	}

	err = p.channel.PublishWithContext(ctx,
		"",                // exchange
		"balance_updates", // routing key (queue name)
		false,             // mandatory
		false,             // immediate
		amqp.Publishing{
			ContentType:  "application/json",
			Body:         body,
			DeliveryMode: amqp.Persistent, // write message on disk
		})

	if err != nil {
		return fmt.Errorf("failed to publish message: %v", err)
	}

	log.Printf("Sent update event for BankID: %d", bankID)
	return nil
}

func (p *RabbitMQProducer) Close() {
	p.channel.Close()
	p.conn.Close()
}
