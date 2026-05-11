package modules

import (
	"strings"

	"github.com/amarnathcjd/gogram/telegram"
)

// Reply — edit if own message (userbot), reply otherwise
func Reply(m *telegram.NewMessage, text string) (*telegram.NewMessage, error) {
	if m.SenderID() == m.Client.Me().ID {
		return m.Edit(text, &telegram.SendOptions{ParseMode: telegram.HTML})
	}
	return m.Reply(text, &telegram.SendOptions{ParseMode: telegram.HTML})
}

// GetArgs — text after command word
func GetArgs(m *telegram.NewMessage) string {
	t := m.Text()
	if i := strings.IndexByte(t, ' '); i != -1 {
		return strings.TrimSpace(t[i+1:])
	}
	return ""
}
