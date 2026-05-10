package modules

import (
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/amarnathcjd/gogram/telegram"

	"pbxgo/config"
)

// ====================== Normal Raids ======================

// .raid [count]
func raidHandler(m *telegram.NewMessage) error {
	return genericTextRaid(m, config.RAID, 50, 150*time.Millisecond)
}

// .hraid [count]
func hraidHandler(m *telegram.NewMessage) error {
	return genericTextRaid(m, config.HRAID, 50, 130*time.Millisecond)
}

// .eraid [count]
func eraidHandler(m *telegram.NewMessage) error {
	return genericTextRaid(m, config.ERAID, 50, 130*time.Millisecond)
}

// .punraid [count]
func punraidHandler(m *telegram.NewMessage) error {
	return genericTextRaid(m, config.PUNRAID, 50, 140*time.Millisecond)
}

// ====================== Reply Raids ======================

// .replyraid [count]
func replyRaidHandler(m *telegram.NewMessage) error {
	return genericReplyRaid(m, config.RAID, 50, 150*time.Millisecond)
}

// .preplyraid [count]
func preplyRaidHandler(m *telegram.NewMessage) error {
	return genericReplyRaid(m, config.PUNRAID, 50, 140*time.Millisecond)
}

// .hreplyraid [count]
func hreplyRaidHandler(m *telegram.NewMessage) error {
	return genericReplyRaid(m, config.HRAID, 50, 130*time.Millisecond)
}

// .ereplyraid [count]
func ereplyRaidHandler(m *telegram.NewMessage) error {
	return genericReplyRaid(m, config.ERAID, 50, 130*time.Millisecond)
}

// ====================== Generic Functions ======================

// Normal Text Raid
func genericTextRaid(
	m *telegram.NewMessage,
	raidList []string,
	maxCount int,
	delay time.Duration,
) error {

	args := GetArgs(m)
	parts := strings.SplitN(args, " ", 2)

	if len(parts) < 1 {
		Reply(m, "⚠️ Usage: <code>.raid 10</code>")
		return nil
	}

	var count int

	fmt.Sscanf(parts[0], "%d", &count)

	if count < 1 || count > maxCount {

		Reply(
			m,
			fmt.Sprintf(
				"⚠️ Count must be between 1 and %d.",
				maxCount,
			),
		)

		return nil
	}

	_, _ = m.Delete()

	for i := 0; i < count; i++ {

		text := raidList[rand.Intn(len(raidList))]

		_, _ = m.Client.SendMessage(
			m.ChatID(),
			text,
			&telegram.SendOptions{
				ParseMode: telegram.MarkDown,
			},
		)

		time.Sleep(delay)
	}

	return nil
}

// Reply Raid
func genericReplyRaid(
	m *telegram.NewMessage,
	raidList []string,
	maxCount int,
	delay time.Duration,
) error {

	if !m.IsReply() {
		Reply(m, "↩️ Reply to a message.")
		return nil
	}

	args := GetArgs(m)
	parts := strings.SplitN(args, " ", 2)

	if len(parts) < 1 {
		Reply(m, "⚠️ Usage: <code>.replyraid 10</code>")
		return nil
	}

	var count int

	fmt.Sscanf(parts[0], "%d", &count)

	if count < 1 || count > maxCount {

		Reply(
			m,
			fmt.Sprintf(
				"⚠️ Count must be between 1 and %d.",
				maxCount,
			),
		)

		return nil
	}

	replyMsg, err := m.GetReplyMessage()

	if err != nil {
		Reply(m, "❌ Could not get replied message.")
		return nil
	}

	_, _ = m.Delete()

	for i := 0; i < count; i++ {

		text := raidList[rand.Intn(len(raidList))]

		_, _ = m.Client.SendMessage(
			m.ChatID(),
			text,
			&telegram.SendOptions{
				ReplyID:   int32(replyMsg.ID),
				ParseMode: telegram.MarkDown,
			},
		)

		time.Sleep(delay)
	}

	return nil
}

func init() {

	rand.Seed(time.Now().UnixNano())

	Register(ModuleInfo{
		Name:        "Raid",
		Description: "Multi-language raid commands",
		Commands: []CommandInfo{
			{
				Pattern: "raid",
				Handler: raidHandler,
				Sudo:    false,
			},
			{
				Pattern: "hraid",
				Handler: hraidHandler,
				Sudo:    false,
			},
			{
				Pattern: "eraid",
				Handler: eraidHandler,
				Sudo:    false,
			},
			{
				Pattern: "punraid",
				Handler: punraidHandler,
				Sudo:    false,
			},
			{
				Pattern: "replyraid",
				Handler: replyRaidHandler,
				Sudo:    false,
			},
			{
				Pattern: "preplyraid",
				Handler: preplyRaidHandler,
				Sudo:    false,
			},
			{
				Pattern: "hreplyraid",
				Handler: hreplyRaidHandler,
				Sudo:    false,
			},
			{
				Pattern: "ereplyraid",
				Handler: ereplyRaidHandler,
				Sudo:    false,
			},
		},
	})
}
