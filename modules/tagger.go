package modules

import (
	"fmt"
	"sync"
	"time"

	"github.com/amarnathcjd/gogram/telegram"
)

var (
	activeTags   = make(map[int64]bool)
	activeTagsMu sync.Mutex
)

// .all [text] — mention all members
func mentionAllHandler(m *telegram.NewMessage) error {
	chatID := m.ChatID()
	args := GetArgs(m)
	hasReply := m.IsReply()

	if args == "" && !hasReply {
		Reply(m, "⚠️ Usage: <code>.all Hello everyone!</code> or reply to a message")
		return nil
	}

	m.Delete()

	activeTagsMu.Lock()
	activeTags[chatID] = true
	activeTagsMu.Unlock()

	members, _, err := m.Client.GetChatMembers(chatID, &telegram.ParticipantOptions{
		Filter: &telegram.ChannelParticipantsRecent{},
		Limit:  200,
	})
	if err != nil {
		Reply(m, "❌ Failed to get members.")
		activeTagsMu.Lock()
		delete(activeTags, chatID)
		activeTagsMu.Unlock()
		return nil
	}

	var replyMsg *telegram.NewMessage
	if hasReply {
		replyMsg, _ = m.GetReplyMessage()
	}

	batch := ""
	count := 0
	for _, member := range members {
		activeTagsMu.Lock()
		active := activeTags[chatID]
		activeTagsMu.Unlock()
		if !active {
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
			if replyMsg != nil {
				replyMsg.Reply(batch, &telegram.SendOptions{ParseMode: telegram.HTML})
			} else {
				m.Client.SendMessage(chatID, args+"\n\n"+batch, &telegram.SendOptions{ParseMode: telegram.HTML})
			}
			time.Sleep(2 * time.Second)
			batch = ""
			count = 0
		}
	}

	if batch != "" {
		if replyMsg != nil {
			replyMsg.Reply(batch, &telegram.SendOptions{ParseMode: telegram.HTML})
		} else {
			m.Client.SendMessage(chatID, args+"\n\n"+batch, &telegram.SendOptions{ParseMode: telegram.HTML})
		}
	}

	activeTagsMu.Lock()
	delete(activeTags, chatID)
	activeTagsMu.Unlock()
	return nil
}

// .cancel — tagger band karo
func cancelHandler(m *telegram.NewMessage) error {
	chatID := m.ChatID()
	activeTagsMu.Lock()
	active := activeTags[chatID]
	activeTagsMu.Unlock()

	if !active {
		Reply(m, "⚠️ No active tagger in this chat.")
		return nil
	}

	activeTagsMu.Lock()
	delete(activeTags, chatID)
	activeTagsMu.Unlock()
	Reply(m, "✅ <b>Tagger cancelled.</b>")
	return nil
}

func init() {
	Register(ModuleInfo{
		Name:        "Tagger",
		Description: "Mention all members in a group",
		Commands: []CommandInfo{
			{Pattern: "all", Handler: mentionAllHandler, Sudo: false},
			{Pattern: "cancel", Handler: cancelHandler, Sudo: false},
		},
	})
}
