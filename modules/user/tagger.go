package user

import (
	"fmt"
	"pbxgo/modules"
	"sync"
	"time"

	"github.com/amarnathcjd/gogram/telegram"
)

var (
	activeTags sync.Map
)

func tagallHandler(m *telegram.NewMessage) error {
	chatID := m.ChatID()
	args := modules.GetArgs(m)
	_, _ = m.Delete()
	activeTags.Store(chatID, true)

	members, _, err := m.Client.GetChatMembers(chatID, &telegram.ParticipantOptions{
		Filter: &telegram.ChannelParticipantsRecent{},
		Limit:  200,
	})
	if err != nil {
		activeTags.Delete(chatID)
		return nil
	}

	batch := ""
	count := 0
	for _, member := range members {
		if v, _ := activeTags.Load(chatID); v == false {
			break
		}
		if member.User == nil || member.User.Bot {
			continue
		}
		u := member.User
		name := u.FirstName
		if name == "" {
			name = "User"
		}
		batch += fmt.Sprintf(`<a href="tg://user?id=%d">%s</a> `, u.ID, name)
		count++
		if count == 5 {
			_, _ = m.Client.SendMessage(chatID, args+"\n\n"+batch, &telegram.SendOptions{ParseMode: telegram.HTML})
			time.Sleep(2 * time.Second)
			batch, count = "", 0
		}
	}
	if batch != "" {
		_, _ = m.Client.SendMessage(chatID, args+"\n\n"+batch, &telegram.SendOptions{ParseMode: telegram.HTML})
	}
	activeTags.Delete(chatID)
	return nil
}

func cancelHandler(m *telegram.NewMessage) error {
	chatID := m.ChatID()
	if v, ok := activeTags.Load(chatID); !ok || v == false {
		modules.Reply(m, "⚠️ No active tagger.")
		return nil
	}
	activeTags.Store(chatID, false)
	modules.Reply(m, "✅ <b>Tagger cancelled.</b>")
	return nil
}

func init() {
	modules.RegisterUser(modules.Module{
		Name: "Tagger",
		Commands: []modules.CommandDef{
			{Pattern: "all", Handler: tagallHandler},
			{Pattern: "cancel", Handler: cancelHandler},
		},
	})
}
