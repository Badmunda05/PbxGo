package modules

import (
	"fmt"
	"strings"
	"time"

	"github.com/amarnathcjd/gogram/telegram"
)

// .spam [count] [text]
func spamHandler(m *telegram.NewMessage) error {
	args := GetArgs(m)
	parts := strings.SplitN(args, " ", 2)
	if len(parts) < 2 {
		Reply(m, "⚠️ Usage: <code>.spam 5 hello</code>")
		return nil
	}
	var count int
	fmt.Sscanf(parts[0], "%d", &count)
	if count < 1 || count > 200 {
		Reply(m, "⚠️ Count must be between 1 and 200.")
		return nil
	}
	text := parts[1]
	m.Delete()
	for i := 0; i < count; i++ {
		m.Client.SendMessage(m.ChatID(), text, &telegram.SendOptions{ParseMode: telegram.Markdown})
		time.Sleep(100 * time.Millisecond)
	}
	return nil
}

// .delayspam [delay_sec] [count] [text]
func delaySpamHandler(m *telegram.NewMessage) error {
	args := GetArgs(m)
	parts := strings.SplitN(args, " ", 3)
	if len(parts) < 3 {
		Reply(m, "⚠️ Usage: <code>.delayspam 1.5 5 hello</code>")
		return nil
	}
	var delaySec float64
	var count int
	fmt.Sscanf(parts[0], "%f", &delaySec)
	fmt.Sscanf(parts[1], "%d", &count)
	text := parts[2]
	if count < 1 || count > 100 {
		Reply(m, "⚠️ Count must be between 1 and 100.")
		return nil
	}
	m.Delete()
	delay := time.Duration(delaySec * float64(time.Second))
	for i := 0; i < count; i++ {
		m.Client.SendMessage(m.ChatID(), text, &telegram.SendOptions{ParseMode: telegram.Markdown})
		time.Sleep(delay)
	}
	return nil
}

// .sspam [count] — reply to sticker
func stickerSpamHandler(m *telegram.NewMessage) error {
	if !m.IsReply() {
		Reply(m, "↩️ Reply to a sticker with <code>.sspam 5</code>")
		return nil
	}
	args := GetArgs(m)
	var count int
	fmt.Sscanf(args, "%d", &count)
	if count < 1 || count > 100 {
		Reply(m, "⚠️ Count must be between 1 and 100.")
		return nil
	}

	r, err := m.GetReplyMessage()
	if err != nil || r.Sticker() == nil {
		Reply(m, "↩️ Reply to a sticker.")
		return nil
	}

	sticker := r.Sticker()
	m.Delete()
	for i := 0; i < count; i++ {
		m.Client.SendMedia(m.ChatID(), sticker.FileID, &telegram.MediaOptions{})
		time.Sleep(100 * time.Millisecond)
	}
	return nil
}

func init() {
	Register(ModuleInfo{
		Name:        "Spam",
		Description: "Spam messages/stickers",
		Commands: []CommandInfo{
			{Pattern: "spam", Handler: spamHandler, Sudo: false},
			{Pattern: "ds", Handler: delaySpamHandler, Sudo: false},
			{Pattern: "delayspam", Handler: delaySpamHandler, Sudo: false},
			{Pattern: "sspam", Handler: stickerSpamHandler, Sudo: false},
		},
	})
}
