package client

import (
	"fmt"
	"main/config"
	"main/database"
	"main/modules"

	"github.com/amarnathcjd/gogram/telegram"
)

var (
	Client  *telegram.Client
	OwnerID int64 = 0
)

func InitClient() error {

	// ─── Init Database ───────────────────────

	database.Init()
	database.LoadSudoUsers()

	// ─── Client Config ──────────────────────

	cfg := telegram.ClientConfig{
		AppID:    config.AppConfig.AppID,
		AppHash:  config.AppConfig.AppHash,
		LogLevel: telegram.LogInfo,

		DeviceConfig: telegram.DeviceConfig{
			DeviceModel:   "PbxGo Bot",
			SystemVersion: "v2.0.0",
		},
	}

	// ─── String Session ─────────────────────

	if config.AppConfig.StringSession != "" {
		cfg.StringSession = config.AppConfig.StringSession
	}

	var (
		err error
	)

	// ─── BOT MODE ───────────────────────────

	if config.AppConfig.BotToken != "" {

		Client, err = telegram.NewClient(
			config.AppConfig.BotToken,
			cfg,
		)

	} else {

		// ─── USERBOT MODE ───────────────────

		Client, err = telegram.NewClient(cfg)
	}

	if err != nil {
		return err
	}

	// ─── Connect ────────────────────────────

	Client.Conn()

	// ─── Auth Prompt (Userbot only) ─────────

	if config.AppConfig.BotToken == "" &&
		config.AppConfig.StringSession == "" {

		if err := Client.AuthPrompt(); err == nil {

			fmt.Println("✅ Authentication successful!")

			fmt.Println("📋 Your String Session:")

			fmt.Println(Client.ExportSession())
		}
	}

	// ─── Logged User ────────────────────────

	OwnerID = Client.Me().ID

	fmt.Printf(
		"✅ Logged in as: %s (ID: %d)\n",
		Client.Me().Username,
		OwnerID,
	)

	return nil
}

// ─────────────────────────────────────────────
// Filters
// ─────────────────────────────────────────────

func IsOwnerFilter(m *telegram.NewMessage) bool {
	return m.SenderID() == m.Client.Me().ID
}

func IsSudoFilter(m *telegram.NewMessage) bool {

	sender := m.SenderID()

	return database.IsSudo(sender)
}

func IsSudoOrOwnerFilter(m *telegram.NewMessage) bool {
	return IsOwnerFilter(m) || IsSudoFilter(m)
}

// ─────────────────────────────────────────────
// Register Handlers
// ─────────────────────────────────────────────

func RegisterHandlers() {

	for _, module := range modules.RegisteredModules {

		for _, command := range module.Commands {

			var filter telegram.Filter

			if command.Sudo {

				filter = telegram.FilterFunc(
					IsSudoOrOwnerFilter,
				)

			} else {

				filter = telegram.FilterFunc(
					IsOwnerFilter,
				)
			}

			Client.On(
				"cmd:"+command.Pattern,
				command.Func,
				filter,
			)
		}
	}
}
