package modules

import (
	"fmt"
	"time"

	"github.com/amarnathcjd/gogram/telegram"
)

// .gcast — sab groups vich broadcast
func gcastHandler(m *telegram.NewMessage) error {
	text := GetArgs(m)
	hasReply := m.IsReply()
	if text == "" && !hasReply {
		Reply(m, "⚠️ Usage: <code>.gcast message</code> or reply to a message")
		return nil
	}

	status, _ := Reply(m, "`📢 Starting global broadcast to groups...`")
	done, failed := 0, 0

	dialogs, err := m.Client.GetDialogs(nil)
	if err != nil {
		Reply(m, "❌ Failed to get dialogs.")
		return nil
	}

	var replyMsg *telegram.NewMessage
	if hasReply {
		replyMsg, _ = m.GetReplyMessage()
	}

	for _, d := range dialogs {
		chat := d.GetChat()
		if chat == nil {
			continue
		}
		peer := chat.GetPeer()
		if peer == nil {
			continue
		}
		_, isGroup := peer.(*telegram.InputPeerChannel)
		_, isChat := peer.(*telegram.InputPeerChat)
		if !isGroup && !isChat {
			continue
		}
		var err error
		if replyMsg != nil {
			_, err = replyMsg.Copy(chat.GetID(), nil)
		} else {
			_, err = m.Client.SendMessage(chat.GetID(), text, nil)
		}
		if err != nil {
			failed++
		} else {
			done++
		}
		time.Sleep(300 * time.Millisecond)
	}

	if status != nil {
		status.Edit(
			fmt.Sprintf("✅ <b>Broadcast done!</b>\n📤 Sent: <code>%d</code>\n❌ Failed: <code>%d</code>", done, failed),
			&telegram.SendOptions{ParseMode: telegram.HTML},
		)
	}
	return nil
}

// .gucast — sab private chats vich broadcast
func gucastHandler(m *telegram.NewMessage) error {
	text := GetArgs(m)
	hasReply := m.IsReply()
	if text == "" && !hasReply {
		Reply(m, "⚠️ Usage: <code>.gucast message</code> or reply to a message")
		return nil
	}

	status, _ := Reply(m, "`📢 Starting global broadcast to private chats...`")
	done, failed := 0, 0

	dialogs, err := m.Client.GetDialogs(nil)
	if err != nil {
		Reply(m, "❌ Failed to get dialogs.")
		return nil
	}

	var replyMsg *telegram.NewMessage
	if hasReply {
		replyMsg, _ = m.GetReplyMessage()
	}

	for _, d := range dialogs {
		chat := d.GetChat()
		if chat == nil {
			continue
		}
		peer := chat.GetPeer()
		if peer == nil {
			continue
		}
		_, isUser := peer.(*telegram.InputPeerUser)
		if !isUser {
			continue
		}
		var err error
		if replyMsg != nil {
			_, err = replyMsg.Copy(chat.GetID(), nil)
		} else {
			_, err = m.Client.SendMessage(chat.GetID(), text, nil)
		}
		if err != nil {
			failed++
		} else {
			done++
		}
		time.Sleep(300 * time.Millisecond)
	}

	if status != nil {
		status.Edit(
			fmt.Sprintf("✅ <b>Broadcast done!</b>\n📤 Sent: <code>%d</code>\n❌ Failed: <code>%d</code>", done, failed),
			&telegram.SendOptions{ParseMode: telegram.HTML},
		)
	}
	return nil
}

func init() {
	Register(ModuleInfo{
		Name:        "Broadcast",
		Description: "Broadcast to groups or private chats",
		Commands: []CommandInfo{
			{Pattern: "gcast", Handler: gcastHandler, Sudo: false},
			{Pattern: "gucast", Handler: gucastHandler, Sudo: false},
		},
	})
}
