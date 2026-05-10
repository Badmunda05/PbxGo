package modules

import (
	"fmt"
	"time"

	"github.com/amarnathcjd/gogram/telegram"
)

// .gcast — sab groups/channels vich broadcast
func gcastHandler(m *telegram.NewMessage) error {
	text := GetArgs(m)
	hasReply := m.IsReply()

	if text == "" && !hasReply {
		Reply(m, "⚠️ Usage: <code>.gcast message</code> or reply to a message")
		return nil
	}

	status, _ := Reply(m, "`📢 Starting broadcast to groups...`")

	dialogs, err := m.Client.GetDialogs(&telegram.DialogOptions{
		Limit: 500,
	})

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

		dtype := dialog.DialogClassName()

		if dtype != "group" && dtype != "channel" && dtype != "supergroup" {
			continue
		}

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
				&telegram.SendOptions{
					ParseMode: telegram.HTML,
				},
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
		done,
		failed,
	)

	status.Edit(msg, &telegram.SendOptions{
		ParseMode: telegram.HTML,
	})

	return nil
}

// .gucast — sab private users vich broadcast
func gucastHandler(m *telegram.NewMessage) error {

	text := GetArgs(m)
	hasReply := m.IsReply()

	if text == "" && !hasReply {
		Reply(m, "⚠️ Usage: <code>.gucast message</code> or reply to a message")
		return nil
	}

	status, _ := Reply(m, "`📢 Starting broadcast to private chats...`")

	dialogs, err := m.Client.GetDialogs(&telegram.DialogOptions{
		Limit: 500,
	})

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

		if dialog.DialogClassName() != "user" {
			continue
		}

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
				&telegram.SendOptions{
					ParseMode: telegram.HTML,
				},
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
		done,
		failed,
	)

	status.Edit(msg, &telegram.SendOptions{
		ParseMode: telegram.HTML,
	})

	return nil
}

func init() {
	Register(ModuleInfo{
		Name:        "Broadcast",
		Description: "Broadcast to groups or private chats",
		Commands: []CommandInfo{
			{
				Pattern: "gcast",
				Handler: gcastHandler,
				Sudo:    false,
			},
			{
				Pattern: "gucast",
				Handler: gucastHandler,
				Sudo:    false,
			},
		},
	})
}
