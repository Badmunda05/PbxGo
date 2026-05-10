package config

import (
	"log/slog"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	AppID         int32
	AppHash       string
	MongoURL      string
	StringSession string
	BotToken      string
	OwnerID       int64
}

var AppConfig Config

func Load() {
	if err := godotenv.Load(); err != nil {
		slog.Warn("No .env file found, using environment variables")
	}

	appID, err := strconv.ParseInt(mustGetEnv("APP_ID"), 10, 32)
	if err != nil {
		slog.Error("Invalid APP_ID", "error", err)
		os.Exit(1)
	}

	ownerID, err := strconv.ParseInt(mustGetEnv("OWNER_ID"), 10, 64)
	if err != nil {
		slog.Error("Invalid OWNER_ID — set your Telegram user ID in .env", "error", err)
		os.Exit(1)
	}

	AppConfig = Config{
		AppID:         int32(appID),
		AppHash:       mustGetEnv("APP_HASH"),
		MongoURL:      getEnv("MONGO_URL", ""),
		StringSession: getEnv("STRING_SESSION", ""),
		BotToken:      getEnv("BOT_TOKEN", ""),
		OwnerID:       ownerID,
	}

	slog.Info("✅ Config loaded", "owner_id", AppConfig.OwnerID)
}

func getEnv(key, fallback string) string {
	if v, ok := os.LookupEnv(key); ok && v != "" {
		return v
	}
	return fallback
}

func mustGetEnv(key string) string {
	v, ok := os.LookupEnv(key)
	if !ok || v == "" {
		slog.Error("Required env variable missing", "key", key)
		os.Exit(1)
	}
	return v
}
