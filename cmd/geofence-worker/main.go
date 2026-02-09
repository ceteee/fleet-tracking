package main

import (
	"encoding/json"
	"fleet-management-system/internal/config"
	"log"
	"os"
	"time"

	rabbitmqConn "fleet-management-system/internal/infra/rabbitmq"

	"github.com/joho/godotenv"
	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {

	if os.Getenv("APP_ENV") == "" {
		_ = godotenv.Load()
	}

	cfg := config.Load()

	conn, ch := rabbitmqConn.MustConnect(cfg)
	defer conn.Close()
	defer ch.Close()

	err := ch.ExchangeDeclare(
		"fleet.events",
		"topic",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatal("exchange declare failed:", err)
	}

	ch.QueueDeclare(
		"geofence_alerts",
		true,
		false,
		false,
		false,
		nil,
	)

	ch.QueueBind(
		"geofence_alerts",
		"",
		"fleet.events",
		false,
		nil,
	)

	msgs, _ := ch.Consume(
		"geofence_alerts",
		"",
		true,
		false,
		false,
		false,
		nil,
	)

	log.Println("Geofence worker started")

	for msg := range msgs {
		var payload map[string]any
		json.Unmarshal(msg.Body, &payload)
		log.Println("GEOFENCE EVENT DISPATCHED:", payload)
	}
}

func ConnectRabbitMQ(cfg *config.Config) (*amqp.Connection, *amqp.Channel, error) {
	var conn *amqp.Connection
	var ch *amqp.Channel
	var err error

	for i := 1; i <= 10; i++ {
		conn, err = amqp.Dial(cfg.RabbitMQURL())
		if err == nil {
			ch, err = conn.Channel()
			if err == nil {
				log.Println("RabbitMQ connected")
				return conn, ch, nil
			}
		}

		log.Printf("RabbitMQ not ready, retry %d/10...", i)
		time.Sleep(2 * time.Second)
	}

	return nil, nil, err
}
