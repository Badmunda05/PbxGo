package modules

import (
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/amarnathcjd/gogram/telegram"

	"yourbot/Config" // Change this if your config path is different
)

// .raid [count]
func raidHandler(m *telegram.NewMessage) error {
	args := GetArgs(m)
	parts := strings.SplitN(args, " ", 2)

	if len(parts) < 1 {
		Reply(m, "⚠️ Usage: <code>.raid 10</code>")
		return nil
	}

	var count int
	fmt.Sscanf(parts[0], "%d", &count)

	if count < 1 || count > 150 {
		Reply(m, "⚠️ Count must be between 1 and 150.")
		return nil
	}

	_, _ = m.Delete()

	// If replied to a message, forward spam it
	if m.IsReply() {
		replyMsg, err := m.GetReplyMessage()
		if err != nil {
			Reply(m, "❌ Could not get replied message.")
			return nil
		}

		for i := 0; i < count; i++ {
			_, _ = m.Client.ForwardMessage(m.ChatID(), replyMsg.ChatID(), replyMsg.ID())
			time.Sleep(80 * time.Millisecond)
		}
		return nil
	}

	// Otherwise use RAID texts from Config
	for i := 0; i < count; i++ {
		text := Config.RAID[rand.Intn(len(Config.RAID))]
		_, _ = m.Client.SendMessage(
			m.ChatID(),
			text,
			&telegram.SendOptions{
				ParseMode: telegram.MarkDown,
			},
		)
		time.Sleep(80 * time.Millisecond)
	}

	return nil
}

// .mraid [delay] [count]
func menuRaidHandler(m *telegram.NewMessage) error {
	args := GetArgs(m)
	parts := strings.SplitN(args, " ", 3)

	if len(parts) < 2 {
		Reply(m, "⚠️ Usage: <code>.mraid 0.8 15</code>")
		return nil
	}

	var delaySec float64
	var count int

	fmt.Sscanf(parts[0], "%f", &delaySec)
	fmt.Sscanf(parts[1], "%d", &count)

	if count < 1 || count > 100 {
		Reply(m, "⚠️ Count must be between 1 and 100.")
		return nil
	}

	delay := time.Duration(delaySec * float64(time.Second))

	_, _ = m.Delete()

	for i := 0; i < count; i++ {
		text := Config.RAID[rand.Intn(len(Config.RAID))]
		_, _ = m.Client.SendMessage(
			m.ChatID(),
			text,
			&telegram.SendOptions{ParseMode: telegram.MarkDown},
		)
		time.Sleep(delay)
	}

	return nil
}

func init() {
	rand.Seed(time.Now().UnixNano()) // For random selection

	Register(ModuleInfo{
		Name:        "Raid",
		Description: "Raid with texts from config",
		Commands: []CommandInfo{
			{
				Pattern: "raid",
				Handler: raidHandler,
				Sudo:    false,
			},
			{
				Pattern: "mraid",
				Handler: menuRaidHandler,
				Sudo:    false,
			},
			{
				Pattern: "menuraid",
				Handler: menuRaidHandler,
				Sudo:    false,
			},
		},
	})
}
