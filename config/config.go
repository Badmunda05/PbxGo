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

var AppConfig *Config

func init() {

	// ─── Load .env ───────────────────────────

	err := godotenv.Load()

	if err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	// ─── Config Init ─────────────────────────

	AppConfig = &Config{
		AppID:         toInt32(getEnv("APP_ID", "0")),
		AppHash:       getEnv("APP_HASH", ""),
		MongoURL:      getEnv("MONGO_URL", ""),
		StringSession: getEnv("STRING_SESSION", ""),
		BotToken:      getEnv("BOT_TOKEN", ""),
	}

	// ─── Validation ──────────────────────────

	if AppConfig.AppID == 0 {
		log.Fatal("APP_ID is missing")
	}

	if AppConfig.AppHash == "" {
		log.Fatal("APP_HASH is missing")
	}

	if AppConfig.BotToken == "" {
		log.Fatal("BOT_TOKEN is missing")
	}
}

// ─────────────────────────────────────────────
// Get ENV
// ─────────────────────────────────────────────

func getEnv(key string, defaultVal string) string {

	value := os.Getenv(key)

	if value == "" {
		return defaultVal
	}

	return value
}

// ─────────────────────────────────────────────
// String → int32
// ─────────────────────────────────────────────

func toInt32(s string) int32 {

	var i int32

	_, err := fmt.Sscanf(s, "%d", &i)

	if err != nil {
		return 0
	}

	return i
}
