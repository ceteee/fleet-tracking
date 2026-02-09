package main

import (
	"fleet-management-system/internal/config"
	"fleet-management-system/internal/fleet/vehicle"
	"log"
	"os"

	vehicleHTTP "fleet-management-system/internal/transport/http/vehicle"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"

	mqttConsumer "fleet-management-system/internal/transport/mqtt"

	database "fleet-management-system/internal/infra/postres"
	rabbitmqConn "fleet-management-system/internal/infra/rabbitmq"
	rabbitmq "fleet-management-system/internal/transport/rabbitmq"
)

func main() {
	if os.Getenv("APP_ENV") == "" {
		_ = godotenv.Load()
	}

	cfg := config.Load()

	db := database.ConnectPostgres(cfg)
	defer db.Close()

	conn, ch := rabbitmqConn.MustConnect(cfg)
	defer conn.Close()
	defer ch.Close()

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
