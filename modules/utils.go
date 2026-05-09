package modules

import "github.com/amarnathcjd/gogram/telegram"

// Reply sends a message; if the sender is the bot itself it edits, otherwise replies.
func Reply(m *telegram.NewMessage, text string) (*telegram.NewMessage, error) {
	if m.SenderID() == m.Client.Me().ID {
		return m.Edit(text, telegram.SendOptions{ParseMode: telegram.HTML})
	}
	return m.Reply(text, &telegram.SendOptions{ParseMode: telegram.HTML})
}
