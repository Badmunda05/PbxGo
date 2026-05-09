package modules

import "github.com/amarnathcjd/gogram/telegram"

func EditOrReply(m *telegram.NewMessage, text string) (*telegram.NewMessage, error) {
	if m.SenderID() == m.Client.Me().ID {
		return m.Edit(text)
	}
	return m.Reply(text)
}
