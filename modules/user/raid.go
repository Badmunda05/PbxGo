package user

import (
	"fmt"
	"math/rand"
	"pbxgo/config"
	"pbxgo/modules"
	"time"

	"github.com/amarnathcjd/gogram/telegram"
)

func doRaid(m *telegram.NewMessage, list []string, max int) error {
	var count int
	fmt.Sscanf(modules.GetArgs(m), "%d", &count)
	if count < 1 || count > max {
		modules.Reply(m, fmt.Sprintf("⚠️ Count: 1–%d", max))
		return nil
	}
	_, _ = m.Delete()
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < count; i++ {
		_, _ = m.Client.SendMessage(m.ChatID(), list[rng.Intn(len(list))], &telegram.SendOptions{ParseMode: telegram.MarkDown})
		time.Sleep(130 * time.Millisecond)
	}
	return nil
}

func doReplyRaid(m *telegram.NewMessage, list []string, max int) error {
	if !m.IsReply() {
		modules.Reply(m, "↩️ Reply to a message.")
		return nil
	}
	var count int
	fmt.Sscanf(modules.GetArgs(m), "%d", &count)
	if count < 1 || count > max {
		modules.Reply(m, fmt.Sprintf("⚠️ Count: 1–%d", max))
		return nil
	}
	r, err := m.GetReplyMessage()
	if err != nil {
		return nil
	}
	_, _ = m.Delete()
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < count; i++ {
		_, _ = m.Client.SendMessage(m.ChatID(), list[rng.Intn(len(list))],
			&telegram.SendOptions{ReplyID: int32(r.ID), ParseMode: telegram.MarkDown})
		time.Sleep(130 * time.Millisecond)
	}
	return nil
}

func init() {
	modules.RegisterUser(modules.Module{
		Name: "Raid",
		Commands: []modules.CommandDef{
			{Pattern: "raid", Handler: func(m *telegram.NewMessage) error { return doRaid(m, config.RAID, 50) }},
			{Pattern: "hraid", Handler: func(m *telegram.NewMessage) error { return doRaid(m, config.HRAID, 50) }},
			{Pattern: "eraid", Handler: func(m *telegram.NewMessage) error { return doRaid(m, config.ERAID, 50) }},
			{Pattern: "punraid", Handler: func(m *telegram.NewMessage) error { return doRaid(m, config.PUNRAID, 50) }},
			{Pattern: "replyraid", Handler: func(m *telegram.NewMessage) error { return doReplyRaid(m, config.RAID, 50) }},
			{Pattern: "hreplyraid", Handler: func(m *telegram.NewMessage) error { return doReplyRaid(m, config.HRAID, 50) }},
			{Pattern: "ereplyraid", Handler: func(m *telegram.NewMessage) error { return doReplyRaid(m, config.ERAID, 50) }},
			{Pattern: "preplyraid", Handler: func(m *telegram.NewMessage) error { return doReplyRaid(m, config.PUNRAID, 50) }},
		},
	})
}
