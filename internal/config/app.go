package config

import (
	"fmt"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	Mode  string
	Host  string
	Port  string
	Store StoreConfig
	MQTT  MQTTConfig
}

func New() *Config {
	err := godotenv.Load(".env." + getEnv("MODE", "development"))
	if err != nil {
		fmt.Println(err.Error())
	}
	return &Config{
		Host: getEnv("SERVER_HOST", "localhost"),
		Port: getEnv("SERVER_PORT", "8000"),
		Mode: getEnv("MODE", "development"),
		Store: StoreConfig{
			Driver: getEnv("DB_DRIVER", "postgres"),
			DSN:    buildDSN(),
		},
		MQTT: MQTTConfig{
			Host:     getEnv("MQTT_HOST", "localhost"),
			Port:     getEnv("MQTT_PORT", "1883"),
			ClientID: getEnv("MQTT_CLIENT_ID", ""),
			Username: getEnv("MQTT_USER", "safe-homie"),
			Password: getEnv("MQTT_PASSWORD", "safe-homie"),
			Topics:   strings.Split(getEnv("MQTT_TOPICS", ""), ","),
		},
	}
}

func getEnv(key string, fallback string) string {
	if value, exist := os.LookupEnv(key); exist {
		return value
	}
	return fallback
}

func buildDSN() string {
	user := getEnv("DB_USER", "postgres")
	password := getEnv("DB_PASSWORD", "password")
	host := getEnv("DB_HOST", "localhost")
	port := getEnv("DB_PORT", "5432")
	dbname := getEnv("DB_NAME", "postgres")
	return "postgres://" + user + ":" + password + "@" + host + ":" + port + "/" + dbname + "?sslmode=disable"
}

func (cfg *Config) isDev() bool {
	return cfg.Mode == "development"
}

func (cfg *Config) AppAddress() string {
	return cfg.Host + ":" + cfg.Port
}
