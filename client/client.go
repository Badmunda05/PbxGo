package client

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"pbxgo/config"
	"pbxgo/database"
	"pbxgo/modules"
	"syscall"

	"github.com/amarnathcjd/gogram/telegram"
)

var (
	Client  *telegram.Client
	OwnerID int64
)

func Init() error {
	cfg := telegram.ClientConfig{
		AppID:    config.AppConfig.AppID,
		AppHash:  config.AppConfig.AppHash,
		LogLevel: telegram.LogInfo,
		DeviceConfig: telegram.DeviceConfig{
			DeviceModel:    "PbxGo",
			SystemVersion:  "Go 1.24",
			AppVersion:     "v2.0.0",
			LangCode:       "en",
			SystemLangCode: "en-US",
		},
	}

	if config.AppConfig.StringSession != "" {
		cfg.StringSession = config.AppConfig.StringSession
	}

	var err error
	Client, err = telegram.NewClient(cfg)
	if err != nil {
		return fmt.Errorf("client creation failed: %w", err)
	}

	if _, err = Client.Conn(); err != nil {
		return fmt.Errorf("connection failed: %w", err)
	}

	switch {
	case config.AppConfig.BotToken != "":
		slog.Info("Starting in Bot mode")
		if err = Client.LoginBot(config.AppConfig.BotToken); err != nil {
			return fmt.Errorf("bot login failed: %w", err)
		}
	case config.AppConfig.StringSession != "":
		slog.Info("Starting in Userbot mode")
	default:
		slog.Info("No session — prompting interactive login")
		if err = Client.AuthPrompt(); err != nil {
			return fmt.Errorf("auth failed: %w", err)
		}
		fmt.Println("\n--- STRING SESSION ---")
		fmt.Println(Client.ExportSession())
		fmt.Println("----------------------\n")
	}

	me := Client.Me()
	OwnerID = config.AppConfig.OwnerID
	slog.Info("Logged in", "name", me.FirstName, "id", me.ID, "owner_id", OwnerID)
	return nil
}

func ownerFilter(m *telegram.NewMessage) error {
	if m.SenderID() != config.AppConfig.OwnerID {
		return fmt.Errorf("unauthorized")
	}
	return nil
}

func sudoOrOwnerFilter(m *telegram.NewMessage) error {
	sid := m.SenderID()
	if sid == config.AppConfig.OwnerID || database.IsSudo(sid) {
		return nil
	}
	return fmt.Errorf("unauthorized")
}

func RegisterHandlers() {
	Client.SetCommandPrefixes(".")

	for _, mod := range modules.RegisteredModules {
		for _, cmd := range mod.Commands {
			if cmd.Sudo {
				Client.On("cmd:"+cmd.Pattern, cmd.Handler, sudoOrOwnerFilter)
			} else {
				Client.On("cmd:"+cmd.Pattern, cmd.Handler, ownerFilter)
			}
		}
		slog.Info("Module registered", "name", mod.Name, "commands", len(mod.Commands))
	}

	// /start — bot mode vich public
	if config.AppConfig.BotToken != "" {
		Client.On("message:/start", modules.StartBotHandler)
	}
}

func Run(ctx context.Context) {
	_ = ctx
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-quit
		slog.Info("Shutting down PbxGo...")
		database.Disconnect()
		os.Exit(0)
	}()
	Client.Idle()
}
