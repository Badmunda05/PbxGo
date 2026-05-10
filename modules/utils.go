package modules

import (
	"pbxgo/client"

	"github.com/amarnathcjd/gogram/telegram"
)

// Reply sends a response:
//   - Userbot mode: if the sender is the bot itself, it edits the message (cleaner UX)
//   - Bot mode: always replies to the user's message
//   - Userbot mode (other sender): replies to the message
func Reply(m *telegram.NewMessage, text string) (*telegram.NewMessage, error) {
	if !client.IsBotMode && m.SenderID() == m.Client.Me().ID {
		// Userbot: edit own message instead of replying to self
		return m.Edit(text, &telegram.SendOptions{ParseMode: telegram.HTML})
	}
	return m.Reply(text, &telegram.SendOptions{ParseMode: telegram.HTML})
}

// GetArgs returns the text after the command word.
// Example: ".spam 5 hello" → "5 hello"
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
