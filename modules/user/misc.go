package user

import (
	"fmt"
	"pbxgo/modules"
	"time"

	"github.com/amarnathcjd/gogram/telegram"
)

const hangMsg = "⃟꙰⃟꙰꙰⃟꙰⃟꙰⃟꙰꙰⃟꙰⃟꙰⃟꙰꙰⃟꙰⃟꙰⃟꙰꙰⃟꙰⃟꙰⃟꙰꙰⃟꙰⃟꙰⃟꙰꙰⃟꙰⃟꙰⃟꙰꙰⃟꙰⃟꙰⃟꙰꙰⃟꙰⃟꙰⃟꙰꙰⃟꙰😈꙰꙰꙰꙰꙰꙰⃟꙰⃟꙰⃟꙰⃟"

func hangHandler(m *telegram.NewMessage) error {
	var count int
	fmt.Sscanf(modules.GetArgs(m), "%d", &count)
	if count < 1 || count > 20 {
		modules.Reply(m, "⚠️ Count: 1–20")
		return nil
	}
	_, _ = m.Delete()
	for i := 0; i < count; i++ {
		_, _ = m.Client.SendMessage(m.ChatID(), hangMsg, nil)
		time.Sleep(300 * time.Millisecond)
	}
	return nil
}

func gcastHandler(m *telegram.NewMessage) error {
	text := modules.GetArgs(m)
	if text == "" && !m.IsReply() {
		modules.Reply(m, "⚠️ Usage: <code>.gcast text</code> or reply")
		return nil
	}
	status, _ := modules.Reply(m, "`📢 Broadcasting...`")
	dialogs, err := m.Client.GetDialogs(&telegram.DialogOptions{Limit: 500})
	if err != nil {
		if status != nil {
			status.Edit("`❌ Failed to get dialogs.`", nil)
		}
		return nil
	}
	var replyMsg *telegram.NewMessage
	if m.IsReply() {
		replyMsg, _ = m.GetReplyMessage()
	}
	done, failed := 0, 0
	for _, d := range dialogs {
		chatID := d.GetID()
		if chatID == 0 {
			continue
		}
		var err error
		if replyMsg != nil {
			_, err = m.Client.Forward(chatID, m.ChatID(), []int32{int32(replyMsg.ID)})
		} else {
			_, err = m.Client.SendMessage(chatID, text, &telegram.SendOptions{ParseMode: telegram.HTML})
		}
		if err != nil {
			failed++
		} else {
			done++
		}
		time.Sleep(300 * time.Millisecond)
	}
	if status != nil {
		status.Edit(fmt.Sprintf("✅ <b>Done!</b> Sent: <code>%d</code> Failed: <code>%d</code>", done, failed),
			&telegram.SendOptions{ParseMode: telegram.HTML})
	}
	return nil
}

func saveMediaHandler(m *telegram.NewMessage) error {
	if !m.IsReply() {
		modules.Reply(m, "↩️ Reply to media with <code>.save</code>")
		return nil
	}
	r, err := m.GetReplyMessage()
	if err != nil || r.Media() == nil {
		modules.Reply(m, "❌ No media found.")
		return nil
	}
	_, _ = m.Delete()
	prog, _ := modules.Reply(m, "`⬇️ Downloading...`")
	fp, err := r.Download(&telegram.DownloadOptions{FileName: "./downloads/"})
	if err != nil {
		if prog != nil {
			prog.Edit("`❌ Download failed: "+err.Error()+"`", nil)
		}
		return nil
	}
	if prog != nil {
		prog.Edit(fmt.Sprintf("✅ **Downloaded:** `%s`", fp),
			&telegram.SendOptions{ParseMode: telegram.MarkDown})
	}
	return nil
}

func init() {
	modules.RegisterUser(modules.Module{
		Name: "Hang",
		Commands: []modules.CommandDef{
			{Pattern: "hang", Handler: hangHandler},
		},
	})
	modules.RegisterUser(modules.Module{
		Name: "Broadcast",
		Commands: []modules.CommandDef{
			{Pattern: "gcast", Handler: gcastHandler},
			{Pattern: "gucast", Handler: gcastHandler},
		},
	})
	modules.RegisterUser(modules.Module{
		Name: "SaveMedia",
		Commands: []modules.CommandDef{
			{Pattern: "save", Handler: saveMediaHandler},
			{Pattern: "savemedia", Handler: saveMediaHandler},
		},
	})
}
