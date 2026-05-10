package modules

import (
	"fmt"
	"os"

	"github.com/amarnathcjd/gogram/telegram"
)

// .savemedia
func saveMediaHandler(m *telegram.NewMessage) error {

	// Must reply
	if !m.IsReply() {

		Reply(
			m,
			"↩️ Reply to a media message with <code>.savemedia</code>",
		)

		return nil
	}

	// Get replied message
	replyMsg, err := m.GetReplyMessage()

	if err != nil {

		Reply(
			m,
			"❌ Could not get replied message.",
		)

		return nil
	}

	// Check media
	media := replyMsg.Media()

	if media == nil {

		Reply(
			m,
			"❌ No media found in replied message.",
		)

		return nil
	}

	_, _ = m.Delete()

	// Create folder
	_ = os.MkdirAll("./downloads", 0755)

	// Progress message
	progress, _ := m.Client.SendMessage(
		m.ChatID(),
		"⬇️ **Downloading media...**",
		&telegram.SendOptions{
			ParseMode: telegram.MarkDown,
		},
	)

	// Download options
	opts := &telegram.DownloadOptions{
		FileName: "./downloads/",
	}

	// Download
	filePath, err := replyMsg.Download(opts)

	if err != nil {

		if progress != nil {

			_, _ = progress.Edit(
				"❌ Download failed:\n`"+err.Error()+"`",
				&telegram.SendOptions{
					ParseMode: telegram.MarkDown,
				},
			)
		}

		return nil
	}

	// Success
	if progress != nil {

		_, _ = progress.Edit(
			fmt.Sprintf(
				"✅ **Downloaded Successfully!**\n\n📁 Saved as:\n`%s`",
				filePath,
			),
			&telegram.SendOptions{
				ParseMode: telegram.MarkDown,
			},
		)
	}

	return nil
}

func init() {

	Register(ModuleInfo{
		Name:        "SaveMedia",
		Description: "Download replied media",
		Commands: []CommandInfo{
			{
				Pattern: "savemedia",
				Handler: saveMediaHandler,
				Sudo:    false,
			},
			{
				Pattern: "save",
				Handler: saveMediaHandler,
				Sudo:    false,
			},
		},
	})
}
