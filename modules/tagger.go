package modules

import (
	"fmt"
	"sync"
	"time"

	"github.com/amarnathcjd/gogram/telegram"
)

var (
	spamChats   = make(map[int64]bool)
	spamChatsMu sync.Mutex
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

	spamChatsMu.Lock()
	spamChats[chatID] = true
	spamChatsMu.Unlock()

	members, err := m.Client.GetChatMembers(chatID, &telegram.GetChatMembersParams{
		Limit: 200,
	})
	if err != nil {
		Reply(m, "❌ Failed to get members.")
		spamChatsMu.Lock()
		delete(spamChats, chatID)
		spamChatsMu.Unlock()
		return nil
	}

	var replyMsg *telegram.NewMessage
	if hasReply {
		replyMsg, _ = m.GetReplyMessage()
	}

	batch := ""
	count := 0
	for _, member := range members {
		spamChatsMu.Lock()
		active := spamChats[chatID]
		spamChatsMu.Unlock()
		if !active {
			break
		}

		user := member.GetUser()
		if user == nil || user.Bot {
			continue
		}

		batch += fmt.Sprintf(`<a href="tg://user?id=%d">%s</a>, `, user.ID, user.FirstName)
		count++

		if count == 5 {
			if replyMsg != nil {
				replyMsg.Reply(batch, &telegram.SendOptions{ParseMode: telegram.HTML})
			} else {
				txt := args + "\n\n" + batch
				m.Client.SendMessage(chatID, txt, &telegram.SendOptions{ParseMode: telegram.HTML})
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

	spamChatsMu.Lock()
	delete(spamChats, chatID)
	spamChatsMu.Unlock()
	return nil
}

// .cancel — tagger stop karo
func cancelHandler(m *telegram.NewMessage) error {
	chatID := m.ChatID()
	spamChatsMu.Lock()
	active := spamChats[chatID]
	spamChatsMu.Unlock()

	if !active {
		Reply(m, "⚠️ No active tagger in this chat.")
		return nil
	}

	spamChatsMu.Lock()
	delete(spamChats, chatID)
	spamChatsMu.Unlock()
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
