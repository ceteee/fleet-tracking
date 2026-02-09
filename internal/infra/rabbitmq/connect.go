package rabbitmq

import (
	"log"
	"time"

	"fleet-management-system/internal/config"

	amqp "github.com/rabbitmq/amqp091-go"
)

func MustConnect(cfg *config.Config) (*amqp.Connection, *amqp.Channel) {
	var (
		conn *amqp.Connection
		ch   *amqp.Channel
		err  error
	)

	for i := 1; i <= 10; i++ {
		conn, err = amqp.Dial(cfg.RabbitMQURL())
		if err == nil {
			ch, err = conn.Channel()
			if err == nil {
				log.Println("RabbitMQ connected")
				return conn, ch
			}
		}

		log.Printf("RabbitMQ not ready, retry %d/10...", i)
		time.Sleep(2 * time.Second)
	}

	log.Fatal("RabbitMQ failed after retries:", err)
	return nil, nil
}
