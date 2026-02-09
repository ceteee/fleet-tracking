package main

import (
	"encoding/json"
	"fleet-management-system/internal/config"
	"log"
	"math/rand/v2"
	"os"

	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/joho/godotenv"
)

type LocationPayload struct {
	VehicleID string  `json:"vehicle_id"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Timestamp int64   `json:"timestamp"`
}

func main() {
	if os.Getenv("APP_ENV") == "" {
		_ = godotenv.Load()
	}

	cfg := config.Load()

	options := mqtt.NewClientOptions()
	options.AddBroker(cfg.MqttURL())
	options.SetClientID("fleet-mqtt-publisher")

	client := mqtt.NewClient(options)

	if token := client.Connect(); token.Wait() && token.Error() != nil {
		log.Fatal("MQTT connection error", token.Error())
	}

	vehicleID := "B1234XYZ"
	topic := "/fleet/vehicle/" + vehicleID + "/location"

	log.Println("Publishing mock location to topic: ", topic)

	for {
		payload := LocationPayload{
			VehicleID: vehicleID,
			Latitude:  -6.2088 + rand.Float64()/1000,
			Longitude: 106.8456 + rand.Float64()/1000,
			Timestamp: time.Now().Unix(),
		}

		data, err := json.Marshal(payload)
		if err != nil {
			log.Println("marshal error:", err)
			continue
		}

		client.Publish(topic, 1, false, data)
		log.Println("Published:", string(data))

		time.Sleep(2 * time.Second)
	}
}
