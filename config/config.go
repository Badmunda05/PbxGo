package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	AppID         int32
	AppHash       string
	MongoURL      string
	StringSession string
	BotToken      string
}

var AppConfig Config

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, using environment variables")
	}

	AppConfig = Config{
		AppID:         toInt32(getEnv("APP_ID", "")),
		AppHash:       getEnv("APP_HASH", ""),
		MongoURL:      getEnv("MONGO_URL", ""),
		StringSession: getEnv("STRING_SESSION", ""),
		BotToken:      getEnv("BOT_TOKEN", ""),
	}
}

func getEnv(key, defaultVal string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultVal
}

func toInt32(s string) int32 {
	var i int32
	_, err := fmt.Sscanf(s, "%d", &i)
	if err != nil {
		log.Fatalf("Invalid APP_ID: %v", err)
	}
	return i
}
