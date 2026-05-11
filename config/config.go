package config

import (
	"log/slog"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

type Cfg struct {
	AppID     int32
	AppHash   string
	BotToken  string
	MongoURL  string
	OwnerID   int64
	LoggerID  int64
	SudoUsers []int64
	Handlers  []string
}

var App Cfg

func Load() {
	_ = godotenv.Load()

	appID, _ := strconv.ParseInt(mustEnv("API_ID"), 10, 32)
	ownerID, _ := strconv.ParseInt(mustEnv("OWNER_ID"), 10, 64)
	loggerID, _ := strconv.ParseInt(getEnv("LOGGER_ID", "0"), 10, 64)

	var sudos []int64
	for _, s := range strings.Fields(getEnv("SUDO_USERS", "")) {
		if id, err := strconv.ParseInt(s, 10, 64); err == nil {
			sudos = append(sudos, id)
		}
	}

	handlers := strings.Fields(getEnv("HANDLERS", ". ! ?"))

	App = Cfg{
		AppID:     int32(appID),
		AppHash:   mustEnv("API_HASH"),
		BotToken:  mustEnv("BOT_TOKEN"),
		MongoURL:  mustEnv("DATABASE_URL"),
		OwnerID:   ownerID,
		LoggerID:  loggerID,
		SudoUsers: sudos,
		Handlers:  handlers,
	}

	slog.Info("✅ Config loaded", "owner_id", App.OwnerID, "handlers", App.Handlers)
}

func getEnv(k, fallback string) string {
	if v, ok := os.LookupEnv(k); ok && v != "" {
		return v
	}
	return fallback
}

func mustEnv(k string) string {
	v, ok := os.LookupEnv(k)
	if !ok || v == "" {
		slog.Error("Required env missing", "key", k)
		os.Exit(1)
	}
	return v
}
