package modules

import (
	"fmt"
	"os"

	"github.com/amarnathcjd/gogram/telegram"
)

// .save / .savemedia — download replied media
func saveMediaHandler(m *telegram.NewMessage) error {
	if !m.IsReply() {
		Reply(m, "↩️ Reply to a media message with <code>.savemedia</code>")
		return nil
	}

	replyMsg, err := m.GetReplyMessage()
	if err != nil {
		Reply(m, "❌ Could not get replied message.")
		return nil
	}

	media := replyMsg.Media()
	if media == nil {
		Reply(m, "❌ No media found in replied message.")
		return nil
	}

	_, _ = m.Delete()

	_ = os.MkdirAll("./downloads", 0755)

	progress, _ := m.Client.SendMessage(
		m.ChatID(),
		"⬇️ **Downloading media...**",
		&telegram.SendOptions{ParseMode: telegram.MarkDown},
	)

	opts := &telegram.DownloadOptions{
		FileName: "./downloads/",
	}

	filePath, err := replyMsg.Download(opts)
	if err != nil {
		if progress != nil {
			_, _ = progress.Edit(
				"❌ Download failed:\n`"+err.Error()+"`",
				&telegram.SendOptions{ParseMode: telegram.MarkDown},
			)
		}
		return nil
	}

	if progress != nil {
		_, _ = progress.Edit(
			fmt.Sprintf("✅ **Downloaded Successfully!**\n\n📁 Saved as:\n`%s`", filePath),
			&telegram.SendOptions{ParseMode: telegram.MarkDown},
		)
	}

	return nil
}

func init() {
	Register(ModuleInfo{
		Name:        "SaveMedia",
		Description: "Download replied media",
		Commands: []CommandInfo{
			{Pattern: "savemedia", Handler: saveMediaHandler, Sudo: false},
			{Pattern: "save", Handler: saveMediaHandler, Sudo: false},
		},
	})
}
