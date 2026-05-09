package modules

import (
	"fmt"
	"time"

	"github.com/amarnathcjd/gogram/telegram"
)

var (
	StartTime = time.Now()
)

func AliveHandler(m *telegram.NewMessage) error {
	aliveMsg := "🤖 <b>PbxGo is alive and kicking!</b> 🔥\n\n"
	aliveMsg += "Bot: <code>PbxGo</code>\n"
	aliveMsg += "Status: <code>Active ✅</code>\n"
	aliveMsg += "Framework: <code>gogram</code>"

	EditOrReply(m, aliveMsg)
	return nil
}

func PingHandler(m *telegram.NewMessage) error {
	uptime := time.Since(StartTime)
	EditOrReply(m, fmt.Sprintf("🏓 <b>Pong!</b>\n\n⏱ Uptime: <code>%v</code>", uptime))
	return nil
}

func init() {
	RegisterModule(ModuleInfo{
		Name:        "Start Module",
		Description: "Provides basic commands like alive and ping.",
		Commands: []CommandInfo{
			{
				Pattern: "alive",
				Func:    AliveHandler,
				Sudo:    true,
			},
			{
				Pattern: "ping",
				Func:    PingHandler,
				Sudo:    true,
			},
		},
	})
}
