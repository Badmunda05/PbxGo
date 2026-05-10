package modules

import (
	"fmt"
	"time"

	"github.com/amarnathcjd/gogram/telegram"
)

// .del
func delHandler(m *telegram.NewMessage) error {

	if m.IsReply() {

		r, err := m.GetReplyMessage()

		if err == nil {

			_, _ = m.Client.DeleteMessages(
				m.ChatID(),
				[]int32{int32(r.ID)},
				false,
			)
		}
	}

	_, _ = m.Delete()

	return nil
}

// .purge
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

		_, err := m.Client.DeleteMessages(
			m.ChatID(),
			batch,
			false,
		)

		if err == nil {
			count += len(batch)
		}

		time.Sleep(200 * time.Millisecond)
	}

	if ex != nil {

		done, _ := ex.Edit(
			fmt.Sprintf(
				"✅ <b>Purge done!</b> Deleted <code>%d</code> messages.",
				count,
			),
			&telegram.SendOptions{
				ParseMode: telegram.HTML,
			},
		)

		if done != nil {

			time.Sleep(2 * time.Second)

			_, _ = done.Delete()
		}
	}

	return nil
}

// .purgeme [n]
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

	// Simplified purgeMe for old gogram
	var ids []int32

	start := int32(m.ID) - n

	if start < 1 {
		start = 1
	}

	for i := start; i <= int32(m.ID); i++ {
		ids = append(ids, i)
	}

	if len(ids) == 0 {
		Reply(m, "❌ No messages found.")
		return nil
	}

	_, _ = m.Client.DeleteMessages(
		m.ChatID(),
		ids,
		false,
	)

	return nil
}

func init() {

	Register(ModuleInfo{
		Name:        "Purge",
		Description: "Delete/purge messages",
		Commands: []CommandInfo{
			{
				Pattern: "del",
				Handler: delHandler,
				Sudo:    false,
			},
			{
				Pattern: "purge",
				Handler: purgeHandler,
				Sudo:    false,
			},
			{
				Pattern: "purgeme",
				Handler: purgeMeHandler,
				Sudo:    false,
			},
		},
	})
}
