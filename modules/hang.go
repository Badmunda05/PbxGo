package modules

import (
	"fmt"
	"time"

	"github.com/amarnathcjd/gogram/telegram"
)

const hangMsg = "⃟꙰⃟꙰꙰⃟꙰⃟꙰⃟꙰꙰⃟꙰⃟꙰⃟꙰꙰⃟꙰⃟꙰⃟꙰꙰⃟꙰⃟꙰⃟꙰꙰⃟꙰⃟꙰⃟꙰꙰⃟꙰⃟꙰⃟꙰꙰⃟꙰⃟꙰⃟꙰꙰⃟꙰⃟꙰⃟꙰꙰⃟꙰😈꙰꙰꙰꙰꙰꙰⃟꙰⃟꙰⃟꙰⃟"

// .hang [count]
func hangHandler(m *telegram.NewMessage) error {
	args := GetArgs(m)
	if args == "" {
		Reply(m, "⚠️ Usage: <code>.hang 3</code>")
		return nil
	}
	var count int
	fmt.Sscanf(args, "%d", &count)
	if count < 1 || count > 20 {
		Reply(m, "⚠️ Count must be between 1 and 20.")
		return nil
	}
	m.Delete()
	for i := 0; i < count; i++ {
		m.Client.SendMessage(m.ChatID(), hangMsg, nil)
		time.Sleep(300 * time.Millisecond)
	}
	return nil
}

func init() {
	Register(ModuleInfo{
		Name:        "Hang",
		Description: "Send hang messages",
		Commands: []CommandInfo{
			{Pattern: "hang", Handler: hangHandler, Sudo: false},
		},
	})
}

var _ = telegram.HTML
