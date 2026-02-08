package mqtt

import (
	"context"
	"encoding/json"
	"fleet-management-system/internal/config"
	"fleet-management-system/internal/fleet/vehicle"
	"log"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type Consumer struct {
	client  mqtt.Client
	cfg     *config.Config
	service *vehicle.Service
}

func NewConsumer(c *config.Config, clientID string, service *vehicle.Service) *Consumer {
	options := mqtt.NewClientOptions()
	options.AddBroker(c.MqttURL())
	options.SetClientID(clientID)

	client := mqtt.NewClient(options)

	if token := client.Connect(); token.Wait() && token.Error() != nil {
		log.Fatal("MQTT connection Error: ", token.Error())
	}

	log.Println("MQTT Consumer connected")

	return &Consumer{
		client:  client,
		cfg:     c,
		service: service,
	}
}

func (c *Consumer) handleMessage(
	client mqtt.Client,
	msg mqtt.Message,
) {
	log.Println("MQTT message received")
	log.Println("Topic   :", msg.Topic())

	var req vehicle.LocationRequest

	if err := json.Unmarshal(msg.Payload(), &req); err != nil {
		log.Println("invalid JSON payload:", err)
		return
	}

	if req.VehicleID == "" {
		log.Println("validation error: vehicle_id is required")
		return
	}
	if req.Timestamp <= 0 {
		log.Println("validation error: invalid timestamp")
		return
	}

	ctx, cancel := context.WithTimeout(
		context.Background(),
		5*time.Second,
	)
	defer cancel()

	if err := c.service.RecordLocation(
		ctx,
		req.VehicleID,
		req.Latitude,
		req.Longitude,
		time.Unix(req.Timestamp, 0),
	); err != nil {
		log.Println("failed to record location:", err)
	}
}

func (c *Consumer) Start() {
	topic := c.cfg.MQTTTopic

	if token := c.client.Subscribe(topic, 1, c.handleMessage); token.Wait() && token.Error() != nil {
		log.Fatal("MQTT subscribe error:", token.Error())
	}

	log.Println("Subscribed to topic:", topic)
}
