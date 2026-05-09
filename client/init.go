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
	database.Init()
	database.LoadSudoUsers()

	cfg := telegram.ClientConfig{
		AppID:    config.AppConfig.AppID,
		AppHash:  config.AppConfig.AppHash,
		LogLevel: telegram.LogInfo,
		DeviceConfig: telegram.DeviceConfig{
			DeviceModel:   "PbxGo Bot",
			SystemVersion: "v2.0.0",
		},
	}

	// Priority: BOT_TOKEN > STRING_SESSION > Auth Prompt
	if config.AppConfig.BotToken != "" {
		cfg.BotToken = config.AppConfig.BotToken
	} else if config.AppConfig.StringSession != "" {
		cfg.StringSession = config.AppConfig.StringSession
	}

	var err error
	Client, err = telegram.NewClient(cfg)
	if err != nil {
		return err
	}

	Client.Conn()

	// If neither BOT_TOKEN nor STRING_SESSION provided, prompt login
	if config.AppConfig.BotToken == "" && config.AppConfig.StringSession == "" {
		if err := Client.AuthPrompt(); err == nil {
			fmt.Println("✅ Authentication successful!")
			fmt.Println("📋 Your String Session (save this):")
			fmt.Println(Client.ExportSession())
		}
	}

	OwnerID = Client.Me().ID
	fmt.Printf("✅ Logged in as: %s (ID: %d)\n", Client.Me().Username, OwnerID)
	return nil
}

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

func RegisterHandlers() {
	for _, module := range modules.RegisteredModules {
		for _, command := range module.Commands {
			var filter telegram.Filter
			if command.Sudo {
				filter = telegram.FilterFunc(IsSudoOrOwnerFilter)
			} else {
				filter = telegram.FilterFunc(IsOwnerFilter)
			}

			Client.On("cmd:"+command.Pattern, command.Func, filter)
		}
	}
}
