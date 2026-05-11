package user

import (
	"fmt"
	"pbxgo/modules"
	"time"

	"github.com/amarnathcjd/gogram/telegram"
)

func delHandler(m *telegram.NewMessage) error {
	if m.IsReply() {
		if r, err := m.GetReplyMessage(); err == nil {
			_, _ = m.Client.DeleteMessages(m.ChatID(), []int32{int32(r.ID)}, false)
		}
	}
	_, _ = m.Delete()
	return nil
}

func purgeHandler(m *telegram.NewMessage) error {
	if !m.IsReply() {
		modules.Reply(m, "↩️ Reply to a message to purge from.")
		return nil
	}
	r, err := m.GetReplyMessage()
	if err != nil {
		return nil
	}
	ex, _ := modules.Reply(m, "`⏳ Purging...`")
	var ids []int32
	for i := int32(r.ID); i <= int32(m.ID); i++ {
		ids = append(ids, i)
	}
	count := 0
	for len(ids) > 0 {
		b := 100
		if len(ids) < b {
			b = len(ids)
		}
		if _, err := m.Client.DeleteMessages(m.ChatID(), ids[:b], false); err == nil {
			count += b
		}
		ids = ids[b:]
		time.Sleep(150 * time.Millisecond)
	}
	if ex != nil {
		done, _ := ex.Edit(fmt.Sprintf("✅ <b>Purged</b> <code>%d</code> messages.", count),
			&telegram.SendOptions{ParseMode: telegram.HTML})
		if done != nil {
			time.Sleep(3 * time.Second)
			_, _ = done.Delete()
		}
	}
	return nil
}

func purgeMeHandler(m *telegram.NewMessage) error {
	var n int32 = 10
	fmt.Sscanf(modules.GetArgs(m), "%d", &n)
	if n < 1 || n > 100 {
		modules.Reply(m, "⚠️ Count: 1–100")
		return nil
	}
	var ids []int32
	start := int32(m.ID) - n
	if start < 1 {
		start = 1
	}
	for i := start; i <= int32(m.ID); i++ {
		ids = append(ids, i)
	}
	_, _ = m.Client.DeleteMessages(m.ChatID(), ids, false)
	return nil
}

func init() {
	modules.RegisterUser(modules.Module{
		Name: "Purge",
		Commands: []modules.CommandDef{
			{Pattern: "del", Handler: delHandler},
			{Pattern: "purge", Handler: purgeHandler},
			{Pattern: "purgeme", Handler: purgeMeHandler},
		},
	})
}
