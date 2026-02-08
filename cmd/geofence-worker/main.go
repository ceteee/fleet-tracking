package main

import (
	"encoding/json"
	"fleet-management-system/internal/config"
	"log"

	"github.com/joho/godotenv"
	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {

	if err := godotenv.Load(); err != nil {
		log.Fatal("Error .env failed to load")
	}

	cfg := config.Load()

	conn, err := amqp.Dial(cfg.RabbitMQURL())
	if err != nil {
		log.Fatal(err)
	}

	ch, _ := conn.Channel()

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

	log.Println("🚀 Geofence worker started")

	for msg := range msgs {
		var payload map[string]any
		json.Unmarshal(msg.Body, &payload)
		log.Println("GEOFENCE EVENT DISPATCHED:", payload)
	}
}
