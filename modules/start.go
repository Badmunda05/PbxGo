package modules

import (
	"fmt"
	"time"

	"github.com/amarnathcjd/gogram/telegram"
)

var startTime = time.Now()

func aliveHandler(m *telegram.NewMessage) error {
	msg := "🤖 <b>PbxGo is alive!</b> 🔥\n\n" +
		"📛 Bot: <code>PbxGo</code>\n" +
		"⚙️  Version: <code>v2.0.0</code>\n" +
		"📟 Status: <code>Online ✅</code>"
	Reply(m, msg)
	return nil
}

func pingHandler(m *telegram.NewMessage) error {
	uptime := time.Since(startTime).Round(time.Second)
	Reply(m, fmt.Sprintf("🏓 <b>Pong!</b>\n\n⏱ Uptime: <code>%s</code>", uptime))
	return nil
}

func init() {
	Register(ModuleInfo{
		Name:        "Core",
		Description: "Basic health-check commands.",
		Commands: []CommandInfo{
			{Pattern: "alive", Handler: aliveHandler, Sudo: true},
			{Pattern: "ping", Handler: pingHandler, Sudo: true},
		},
	})
}

// Ensure telegram import is used (SendOptions used in utils.go)
var _ = telegram.HTML
