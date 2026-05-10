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

	// Batches of 100
	count := 0
	for len(ids) > 0 {
		batch := ids
		if len(batch) > 100 {
			batch = ids[:100]
		}
		ids = ids[len(batch):]
		_, err := m.Client.DeleteMessages(m.ChatID(), batch, nil)
		if err == nil {
			count += len(batch)
		}
		time.Sleep(200 * time.Millisecond)
	}

	if ex != nil {
		done, _ := ex.Edit(fmt.Sprintf("✅ <b>Purge complete!</b> Deleted <code>%d</code> messages.", count),
			&telegram.SendOptions{ParseMode: telegram.HTML})
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
	if args == "" {
		m.Delete()
		return nil
	}
	var n int
	fmt.Sscanf(args, "%d", &n)
	if n < 1 {
		Reply(m, "⚠️ Usage: <code>.purgeme 10</code>")
		return nil
	}

	// Search own messages
	msgs, err := m.Client.SearchMessages(m.ChatID(), &telegram.SearchOptions{
		Query:  "",
		Limit:  int32(n + 1),
		FromID: m.Client.Me().ID,
	})
	if err != nil || len(msgs) == 0 {
		Reply(m, "❌ No messages found.")
		return nil
	}

	var ids []int32
	for _, msg := range msgs {
		ids = append(ids, int32(msg.ID))
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
