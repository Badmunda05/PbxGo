package main

import (
	"context"
	"log/slog"
	"os"
	"pbxgo/client"
	"pbxgo/config"
	"pbxgo/database"
)

func main() {
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo})))
	slog.Info("🚀 PbxGo v3.0 starting — Multi-Session Edition")

	config.Load()
	database.Init()
	database.LoadSudoUsers()

	// 1. Start the bot
	if err := client.StartBot(); err != nil {
		slog.Error("Bot start failed", "error", err)
		os.Exit(1)
	}

	// 2. Load all saved userbot sessions from MongoDB
	client.LoadUserSessions()

	// 3. Register bot command handlers
	client.RegisterBotHandlers()

	slog.Info("✅ PbxGo is running",
		"bot", client.Bot.Me().Username,
		"user_sessions", len(client.Users),
	)

	slog.Info("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	slog.Info("Add userbots via the bot:")
	slog.Info("  /add YOUR_SESSION_STRING")
	slog.Info("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

	client.Run(context.Background())
}
