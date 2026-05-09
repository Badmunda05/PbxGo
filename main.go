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
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
	slog.SetDefault(logger)

	slog.Info("🚀 PbxGo starting", "version", "v2.0.0")

	config.Load()
	database.Init()
	database.LoadSudoUsers()

	ctx := context.Background()
	if err := client.Init(); err != nil {
		slog.Error("Failed to start client", "error", err)
		os.Exit(1)
	}

	client.RegisterHandlers()
	slog.Info("✅ PbxGo is running — press Ctrl+C to stop")
	client.Run(ctx)
}
