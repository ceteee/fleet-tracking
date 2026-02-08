package config

import (
	"fmt"
	"log"
	"os"
)

type Config struct {
	AppEnv  string
	AppHost string
	AppPort string

	// PostgreSQL
	DBHost string
	DBPort string
	DBUser string
	DBPass string
	DBName string

	// MQTT
	MQTTHost  string
	MQTTPort  string
	MQTTTopic string

	// RabbitMQ
	RabbitMQHost  string
	RabbitMQPort  string
	RabbitMQUser  string
	RabbitMQPass  string
	RabbitMQVHost string
}

func Load() *Config {
	cfg := &Config{
		AppEnv:  getEnv("APP_ENV", "development"),
		AppHost: getEnv("APP_HOST", "localhost"),
		AppPort: getEnv("APP_PORT", "8080"),

		// PostgreSQL
		DBHost: getEnv("DB_HOST", "localhost"),
		DBPort: getEnv("DB_POST", "5432"),
		DBUser: getEnv("DB_USER", "postgres"),
		DBPass: getEnv("DB_PASS", "chrstmbn"),
		DBName: getEnv("DB_NAME", "fleet_app"),

		// MQTT
		MQTTHost:  getEnv("MQTT_HOST", "localhost"),
		MQTTPort:  getEnv("MQTT_PORT", "1883"),
		MQTTTopic: getEnv("MQTT_TOPIC", "/fleet/vehicle/+/location"),

		// RabbitMQ
		RabbitMQHost:  getEnv("RABBITMQ_HOST", "localhost"),
		RabbitMQPort:  getEnv("RABBITMQ_PORT", "5672"),
		RabbitMQUser:  getEnv("RABBITMQ_USER", "guest"),
		RabbitMQPass:  getEnv("RABBITMQ_PASS", "guest"),
		RabbitMQVHost: getEnv("RABBITMQ_VHOST", "/"),
	}

	log.Println("Config loaded:", cfg.AppEnv)

	return cfg
}

func getEnv(key, defaultVal string) string {
	if val, status := os.LookupEnv(key); status {
		return val
	}

	return defaultVal
}

func (c *Config) PostgresDSN() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		c.DBUser,
		c.DBPass,
		c.DBHost,
		c.DBPort,
		c.DBName,
	)
}

func (c *Config) MqttURL() string {
	return fmt.Sprintf("tcp://%s:%s", c.MQTTHost, c.MQTTPort)
}

func (c *Config) RabbitMQURL() string {
	return fmt.Sprintf("amqp://%s:%s@%s:%s", c.RabbitMQUser, c.RabbitMQPass, c.RabbitMQHost, c.RabbitMQPort)
}
