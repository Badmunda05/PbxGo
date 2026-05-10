package modules

import (
	"os"
	"time"

	"github.com/amarnathcjd/gogram/telegram"
)

func restartHandler(m *telegram.NewMessage) error {
	Reply(m, "♻️ <b>PbxGo Restarting...</b>\n\n⏳ Wait 10-15 seconds...")
	time.Sleep(1 * time.Second)
	os.Exit(0)
	return nil
}

func init() {
	Register(ModuleInfo{
		Name:        "Restart",
		Description: "Restart the bot",
		Commands: []CommandInfo{
			{Pattern: "restart", Handler: restartHandler, Sudo: false},
			{Pattern: "rs", Handler: restartHandler, Sudo: false},
			{Pattern: "reload", Handler: restartHandler, Sudo: false},
		},
	})
}

var _ = telegram.HTML
