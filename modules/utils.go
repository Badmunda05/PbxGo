package modules

import "github.com/amarnathcjd/gogram/telegram"

// Reply — userbot mode: edits own message, bot mode: replies
func Reply(m *telegram.NewMessage, text string) (*telegram.NewMessage, error) {
	if m.SenderID() == m.Client.Me().ID {
		return m.Edit(text, &telegram.SendOptions{ParseMode: telegram.HTML})
	}
	return m.Reply(text, &telegram.SendOptions{ParseMode: telegram.HTML})
}

// GetArgs — command ke baad ka text nikalta hai
func GetArgs(m *telegram.NewMessage) string {
	text := m.Text()
	if idx := findSpace(text); idx != -1 {
		return text[idx+1:]
	}
	return ""
}

func findSpace(s string) int {
	for i, c := range s {
		if c == ' ' {
			return i
		}
	}
	return -1
}
