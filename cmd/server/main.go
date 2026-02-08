package main

import (
	"database/sql"
	"fleet-management-system/internal/config"
	"fleet-management-system/internal/fleet/vehicle"
	"log"

	vehicleHTTP "fleet-management-system/internal/transport/http/vehicle"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"

	mqttConsumer "fleet-management-system/internal/transport/mqtt"

	rabbitmq "fleet-management-system/internal/transport/rabbitmq"

	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error .env failed to load")
	}

	cfg := config.Load()

	db, err := sql.Open("postgres", cfg.PostgresDSN())

	if err != nil {
		log.Fatal(err)
	}

	if err := db.Ping(); err != nil {
		log.Fatal("failed to connect database:", err)
	}

	conn, err := amqp.Dial(cfg.RabbitMQURL())

	if err != nil {
		log.Fatal("failed to connect rabbitmq:", err)
	}

	defer conn.Close()

	publisher, err := rabbitmq.NewPublisher(conn)
	if err != nil {
		log.Fatal("failed to create rabbitmq publisher:", err)
	}

	repo := vehicle.NewRepository(db)
	service := vehicle.NewService(repo, publisher)
	handler := vehicleHTTP.NewHandler(service)

	server := gin.Default()
	vehicleGroup := server.Group("/vehicles")
	{
		vehicleGroup.POST("/locations", handler.CreateLocation)
		vehicleGroup.GET("/:vehicle_id/location", handler.GetLatestLocation)
		vehicleGroup.GET("/:vehicle_id/history", handler.GetLocationHistory)
	}

	consumer := mqttConsumer.NewConsumer(
		cfg,
		"fleet-mqtt-consumer",
		service,
	)

	go consumer.Start()

	log.Println("Welcome to fleet service")
	if err := server.Run(":" + cfg.AppPort); err != nil {
		log.Fatal(err)
	}
}
