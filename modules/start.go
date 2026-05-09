package modules

import (
	"fmt"
	"time"

	"github.com/amarnathcjd/gogram/telegram"
)

var startTime = time.Now()

// ─────────────────────────────────────────────
// Alive
// ─────────────────────────────────────────────

func aliveHandler(
	m *telegram.NewMessage,
) error {

	msg :=
		"✨ <b>ᴘʙxɢᴏ ɪs ᴏɴʟɪɴᴇ</b>\n\n" +
			"⚡ <b>sᴛᴀᴛᴜs:</b> <code>ᴏɴʟɪɴᴇ ✅</code>\n" +
			"🧠 <b>ʟɪʙʀᴀʀʏ:</b> <code>ɢᴏɢʀᴀᴍ</code>\n" +
			"⚙️ <b>ᴠᴇʀsɪᴏɴ:</b> <code>v2.0.0</code>\n" +
			"🔥 <b>ᴍᴏᴅᴇ:</b> <code>ᴜsᴇʀʙᴏᴛ</code>"

	_, err := Reply(m, msg)

	return err
}

// ─────────────────────────────────────────────
// Ping
// ─────────────────────────────────────────────

func pingHandler(
	m *telegram.NewMessage,
) error {

	uptime := time.Since(
		startTime,
	).Round(time.Second)

	msg := fmt.Sprintf(
		"🏓 <b>ᴘᴏɴɢ!</b>\n\n"+
			"⚡ <b>sᴘᴇᴇᴅ:</b> <code>ғᴀsᴛ ᴀғ</code>\n"+
			"⏱ <b>ᴜᴘᴛɪᴍᴇ:</b> <code>%s</code>",
		uptime,
	)

	_, err := Reply(m, msg)

	return err
}

// ─────────────────────────────────────────────
// Start
// ─────────────────────────────────────────────

func startHandler(
	m *telegram.NewMessage,
) error {

	msg :=
		"👋 <b>ᴡᴇʟᴄᴏᴍᴇ ᴛᴏ ᴘʙxɢᴏ</b>\n\n" +
			"⚡ ғᴀsᴛ & ʟɪɢʜᴛᴡᴇɪɢʜᴛ ɢᴏ ᴜsᴇʀʙᴏᴛ\n" +
			"🧠 ᴘᴏᴡᴇʀᴇᴅ ʙʏ <b>ɢᴏɢʀᴀᴍ</b>\n" +
			"🔥 ʙᴜɪʟᴛ ꜰᴏʀ sᴘᴇᴇᴅ\n\n" +
			"📌 <b>ᴄᴏᴍᴍᴀɴᴅs:</b>\n" +
			"• <code>.alive</code>\n" +
			"• <code>.ping</code>\n" +
			"• <code>.help</code>"

	_, err := Reply(m, msg)

	return err
}

// ─────────────────────────────────────────────
// Register
// ─────────────────────────────────────────────

func init() {

	Register(ModuleInfo{
		Name:        "Core",
		Description: "Basic core commands.",

		Commands: []CommandInfo{

			{
				Pattern: "start",
				Handler: startHandler,
				Sudo:    true,
			},

			{
				Pattern: "alive",
				Handler: aliveHandler,
				Sudo:    true,
			},

			{
				Pattern: "ping",
				Handler: pingHandler,
				Sudo:    true,
			},
		},
	})
}

var _ = telegram.HTML
