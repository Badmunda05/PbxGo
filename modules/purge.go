package modules

import (
	"fmt"
	"time"

	"github.com/amarnathcjd/gogram/telegram"
)

// .del — replied message delete karo
func delHandler(m *telegram.NewMessage) error {
	if m.IsReply() {
		r, err := m.GetReplyMessage()
		if err == nil {
			m.Client.DeleteMessages(m.ChatID(), []int32{int32(r.ID)}, nil)
		}
	}
	m.Delete()
	return nil
}

// .purge — replied message se lekar current tak sab delete
func purgeHandler(m *telegram.NewMessage) error {
	if !m.IsReply() {
		Reply(m, "↩️ Reply karo kisi message pe purge karne ke liye.")
		return nil
	}
	r, err := m.GetReplyMessage()
	if err != nil {
		return nil
	}

	ex, _ := Reply(m, "`⏳ Purging messages...`")

	var ids []int32
	for i := int32(r.ID); i <= int32(m.ID); i++ {
		ids = append(ids, i)
	}

	count := 0
	for len(ids) > 0 {
		batchSize := 100
		if len(ids) < batchSize {
			batchSize = len(ids)
		}
		batch := ids[:batchSize]
		ids = ids[batchSize:]
		_, err := m.Client.DeleteMessages(m.ChatID(), batch, nil)
		if err == nil {
			count += len(batch)
		}
		time.Sleep(200 * time.Millisecond)
	}

	if ex != nil {
		done, _ := ex.Edit(
			fmt.Sprintf("✅ <b>Purge done!</b> Deleted <code>%d</code> messages.", count),
			&telegram.SendOptions{ParseMode: telegram.HTML},
		)
		if done != nil {
			time.Sleep(2 * time.Second)
			done.Delete()
		}
	}
	return nil
}

// .purgeme [n] — apne last n messages delete karo
func purgeMeHandler(m *telegram.NewMessage) error {
	args := GetArgs(m)
	var n int32 = 10
	if args != "" {
		fmt.Sscanf(args, "%d", &n)
	}
	if n < 1 || n > 100 {
		Reply(m, "⚠️ Count must be between 1 and 100.")
		return nil
	}

	// Get message history and filter own messages
	history, err := m.Client.GetMessages(m.ChatID(), &telegram.MessagesOptions{
		Limit: n * 3, // fetch more to find own msgs
	})
	if err != nil || len(history) == 0 {
		Reply(m, "❌ No messages found.")
		return nil
	}

	myID := m.Client.Me().ID
	var ids []int32
	for _, msg := range history {
		if int32(len(ids)) >= n {
			break
		}
		if msg.SenderID() == myID {
			ids = append(ids, int32(msg.ID))
		}
	}

	if len(ids) == 0 {
		Reply(m, "❌ No own messages found.")
		return nil
	}

	m.Client.DeleteMessages(m.ChatID(), ids, nil)
	return nil
}

func init() {
	Register(ModuleInfo{
		Name:        "Purge",
		Description: "Delete/purge messages",
		Commands: []CommandInfo{
			{Pattern: "del", Handler: delHandler, Sudo: false},
			{Pattern: "purge", Handler: purgeHandler, Sudo: false},
			{Pattern: "purgeme", Handler: purgeMeHandler, Sudo: false},
		},
	})
}
