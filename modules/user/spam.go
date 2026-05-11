package user

import (
	"fmt"
	"pbxgo/modules"
	"strings"
	"time"

	"github.com/amarnathcjd/gogram/telegram"
)

func spamHandler(m *telegram.NewMessage) error {
	args := modules.GetArgs(m)
	parts := strings.SplitN(args, " ", 2)
	if len(parts) < 2 {
		modules.Reply(m, "⚠️ Usage: <code>.spam 5 hello</code>")
		return nil
	}
	var count int
	fmt.Sscanf(parts[0], "%d", &count)
	if count < 1 || count > 200 {
		modules.Reply(m, "⚠️ Count: 1–200")
		return nil
	}
	_, _ = m.Delete()
	for i := 0; i < count; i++ {
		_, _ = m.Client.SendMessage(m.ChatID(), parts[1], &telegram.SendOptions{ParseMode: telegram.MarkDown})
		time.Sleep(100 * time.Millisecond)
	}
	return nil
}

func delaySpamHandler(m *telegram.NewMessage) error {
	args := modules.GetArgs(m)
	parts := strings.SplitN(args, " ", 3)
	if len(parts) < 3 {
		modules.Reply(m, "⚠️ Usage: <code>.ds 1.5 5 hello</code>")
		return nil
	}
	var delaySec float64
	var count int
	fmt.Sscanf(parts[0], "%f", &delaySec)
	fmt.Sscanf(parts[1], "%d", &count)
	if count < 1 || count > 100 {
		modules.Reply(m, "⚠️ Count: 1–100")
		return nil
	}
	_, _ = m.Delete()
	delay := time.Duration(delaySec * float64(time.Second))
	for i := 0; i < count; i++ {
		_, _ = m.Client.SendMessage(m.ChatID(), parts[2], &telegram.SendOptions{ParseMode: telegram.MarkDown})
		time.Sleep(delay)
	}
	return nil
}

func stickerSpamHandler(m *telegram.NewMessage) error {
	if !m.IsReply() {
		modules.Reply(m, "↩️ Reply to a sticker with <code>.sspam 5</code>")
		return nil
	}
	var count int
	fmt.Sscanf(modules.GetArgs(m), "%d", &count)
	if count < 1 || count > 100 {
		modules.Reply(m, "⚠️ Count: 1–100")
		return nil
	}
	r, err := m.GetReplyMessage()
	if err != nil || r.Media() == nil {
		modules.Reply(m, "↩️ Reply to a sticker/media.")
		return nil
	}
	_, _ = m.Delete()
	for i := 0; i < count; i++ {
		_, _ = m.Client.SendMedia(m.ChatID(), r.Media(), &telegram.MediaOptions{})
		time.Sleep(100 * time.Millisecond)
	}
	return nil
}

func init() {
	modules.RegisterUser(modules.Module{
		Name: "Spam",
		Commands: []modules.CommandDef{
			{Pattern: "spam", Handler: spamHandler},
			{Pattern: "ds", Handler: delaySpamHandler},
			{Pattern: "delayspam", Handler: delaySpamHandler},
			{Pattern: "sspam", Handler: stickerSpamHandler},
		},
	})
}
