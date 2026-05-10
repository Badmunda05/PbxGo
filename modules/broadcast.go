package modules

import (
	"fmt"
	"time"

	"github.com/amarnathcjd/gogram/telegram"
)

// .gcast — broadcast to all dialogs (groups + private)
func gcastHandler(m *telegram.NewMessage) error {
	text := GetArgs(m)
	hasReply := m.IsReply()

	if text == "" && !hasReply {
		Reply(m, "⚠️ Usage: <code>.gcast message</code> or reply to a message")
		return nil
	}

	status, _ := Reply(m, "`📢 Starting broadcast...`")

	dialogs, err := m.Client.GetDialogs(&telegram.DialogOptions{Limit: 500})
	if err != nil {
		status.Edit("`❌ Failed to get dialogs.`", nil)
		return nil
	}

	var replyMsg *telegram.NewMessage
	if hasReply {
		replyMsg, _ = m.GetReplyMessage()
	}

	done := 0
	failed := 0

	for _, dialog := range dialogs {
		chatID := dialog.GetID()
		if chatID == 0 {
			continue
		}

		var err error
		if replyMsg != nil {
			_, err = m.Client.Forward(
				chatID,
				m.ChatID(),
				[]int32{int32(replyMsg.ID)},
			)
		} else {
			_, err = m.Client.SendMessage(
				chatID,
				text,
				&telegram.SendOptions{ParseMode: telegram.HTML},
			)
		}

		if err != nil {
			failed++
		} else {
			done++
		}
		time.Sleep(300 * time.Millisecond)
	}

	msg := fmt.Sprintf(
		"✅ <b>Broadcast done!</b>\n📤 Sent: <code>%d</code>\n❌ Failed: <code>%d</code>",
		done, failed,
	)
	status.Edit(msg, &telegram.SendOptions{ParseMode: telegram.HTML})
	return nil
}

// .gucast — broadcast to private chats only
func gucastHandler(m *telegram.NewMessage) error {
	text := GetArgs(m)
	hasReply := m.IsReply()

	if text == "" && !hasReply {
		Reply(m, "⚠️ Usage: <code>.gucast message</code> or reply to a message")
		return nil
	}

	status, _ := Reply(m, "`📢 Starting private broadcast...`")

	dialogs, err := m.Client.GetDialogs(&telegram.DialogOptions{Limit: 500})
	if err != nil {
		status.Edit("`❌ Failed to get dialogs.`", nil)
		return nil
	}

	var replyMsg *telegram.NewMessage
	if hasReply {
		replyMsg, _ = m.GetReplyMessage()
	}

	done := 0
	failed := 0

	for _, dialog := range dialogs {
		chatID := dialog.GetID()
		if chatID == 0 {
			continue
		}

		var err error
		if replyMsg != nil {
			_, err = m.Client.Forward(
				chatID,
				m.ChatID(),
				[]int32{int32(replyMsg.ID)},
			)
		} else {
			_, err = m.Client.SendMessage(
				chatID,
				text,
				&telegram.SendOptions{ParseMode: telegram.HTML},
			)
		}

		if err != nil {
			failed++
		} else {
			done++
		}
		time.Sleep(300 * time.Millisecond)
	}

	msg := fmt.Sprintf(
		"✅ <b>Broadcast done!</b>\n📤 Sent: <code>%d</code>\n❌ Failed: <code>%d</code>",
		done, failed,
	)
	status.Edit(msg, &telegram.SendOptions{ParseMode: telegram.HTML})
	return nil
}

func init() {
	Register(ModuleInfo{
		Name:        "Broadcast",
		Description: "Broadcast messages",
		Commands: []CommandInfo{
			{Pattern: "gcast", Handler: gcastHandler, Sudo: false},
			{Pattern: "gucast", Handler: gucastHandler, Sudo: false},
		},
	})
}
