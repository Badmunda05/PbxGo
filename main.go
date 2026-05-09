package main

import (
	"context"
	"log/slog"
	"os"
	"pbxgo/client"
	"pbxgo/config"
	"pbxgo/database"

	_ "pbxgo/modules" // register all modules via init()
)

func main() {
	// Structured JSON logging (Go 1.24 standard)
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
	slog.SetDefault(logger)

	slog.Info("🚀 PbxGo starting", "version", "v2.0.0")

	// Load config from .env
	config.Load()

	// Init database
	database.Init()
	database.LoadSudoUsers()

	// Init Telegram client
	ctx := context.Background()
	if err := client.Init(); err != nil {
		slog.Error("Failed to start client", "error", err)
		os.Exit(1)
	}

	// Register all module handlers
	client.RegisterHandlers()

	slog.Info("✅ PbxGo is running — press Ctrl+C to stop")

	// Block until signal
	client.Run(ctx)
}
