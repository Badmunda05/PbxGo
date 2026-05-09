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

	// ─── USERBOT SESSION ─────────────────────

	if config.AppConfig.StringSession != "" {

		slog.Info(
			"Starting in Userbot mode",
		)

		cfg.StringSession =
			config.AppConfig.StringSession

	} else {

		slog.Info(
			"No session found — login required",
		)
	}

	var err error

	// ─── Create Client ───────────────────────

	Client, err = telegram.NewClient(cfg)

	if err != nil {
		return fmt.Errorf(
			"client creation failed: %w",
			err,
		)
	}

	// ─── Connect ─────────────────────────────

	_, err = Client.Conn()

	if err != nil {
		return fmt.Errorf(
			"connection failed: %w",
			err,
		)
	}

	// ─── Login Prompt ────────────────────────

	if config.AppConfig.StringSession == "" {

		if err = Client.AuthPrompt(); err != nil {

			return fmt.Errorf(
				"auth failed: %w",
				err,
			)
		}

		session := Client.ExportSession()

		slog.Info(
			"Login successful",
		)

		fmt.Println(
			"\n--- STRING SESSION ---",
		)

		fmt.Println(session)

		fmt.Println(
			"----------------------\n",
		)
	}

	me := Client.Me()

	OwnerID = me.ID

	slog.Info(
		"Logged in",
		"name", me.FirstName,
		"username", me.Username,
		"id", me.ID,
	)

	return nil
}

func RegisterHandlers() {

	Client.SetCommandPrefixes(".")

	for _, mod := range modules.RegisteredModules {

		for _, cmd := range mod.Commands {

			Client.On(
				"cmd:"+cmd.Pattern,
				cmd.Handler,
			)
		}

		slog.Info(
			"Module registered",
			"name", mod.Name,
			"commands", len(mod.Commands),
		)
	}
}

func Run(ctx context.Context) {

	_ = ctx

	quit := make(chan os.Signal, 1)

	signal.Notify(
		quit,
		syscall.SIGINT,
		syscall.SIGTERM,
	)

	go func() {

		<-quit

		slog.Info(
			"Shutting down PbxGo...",
		)

		database.Disconnect()

		os.Exit(0)
	}()

	Client.Idle()
}

// ─────────────────────────────────────────────
// Filters
// ─────────────────────────────────────────────

func isOwner(
	m *telegram.NewMessage,
) bool {

	return m.SenderID() ==
		Client.Me().ID
}

func isSudo(
	m *telegram.NewMessage,
) bool {

	return database.IsSudo(
		m.SenderID(),
	)
}

func isSudoOrOwner(
	m *telegram.NewMessage,
) bool {

	return isOwner(m) ||
		isSudo(m)
}
