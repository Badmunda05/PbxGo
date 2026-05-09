package modules

import (
	"fmt"
	"runtime"
	"time"

	"github.com/amarnathcjd/gogram/telegram"
)

var startTime = time.Now()

const version = "v2.0.0"

func goVersion() string {
	return runtime.Version()
}

func uptime() string {
	d := time.Since(startTime).Round(time.Second)
	h := int(d.Hours())
	m := int(d.Minutes()) % 60
	s := int(d.Seconds()) % 60
	if h > 0 {
		return fmt.Sprintf("%dh %dm %ds", h, m, s)
	}
	if m > 0 {
		return fmt.Sprintf("%dm %ds", m, s)
	}
	return fmt.Sprintf("%ds", s)
}

// aliveHandler — .alive command
func aliveHandler(m *telegram.NewMessage) error {
	msg := fmt.Sprintf(
		"** 🥀 ᴘʙxɢᴏ 💞 **\n\n"+
			"❏ **ᴠᴇʀsɪᴏɴ**: `%s`\n"+
			"├• **ᴜᴘᴛɪᴍᴇ**: `%s`\n"+
			"├• **ɢᴏ**: `%s`\n"+
			"├• **sᴜᴘᴘᴏʀᴛ-𝐂ʜᴀᴛ**: [|| ˹ᴘʙx_sᴜᴘᴘᴏʀᴛ˼ ||](https://t.me/PBXCHATS)\n"+
			"├• **ᴜᴘᴅᴀᴛᴇs**: [ᴘʙx_ᴜᴘᴅᴀᴛᴇ](https://t.me/PBX_UPDATE)\n"+
			"└• **ᴏᴡɴᴇʀ**: [ʙᴀᴅ ᴍᴜɴᴅᴀ](https://t.me/Badmundaxd)",
		version,
		uptime(),
		goVersion(),
	)
	Reply(m, msg)
	return nil
}

// pingHandler — .ping command
func pingHandler(m *telegram.NewMessage) error {
	start := time.Now()
	sent, err := Reply(m, "🏓 ᴘɪɴɢɪɴɢ...")
	if err != nil || sent == nil {
		return err
	}
	duration := time.Since(start).Milliseconds()

	me := m.Client.Me()
	mention := fmt.Sprintf("<a href=\"tg://user?id=%d\">%s</a>", me.ID, me.FirstName)

	msg := fmt.Sprintf(
		"❏ **╰☞ ᴘʙxɢᴏ**\n"+
			"├• **╰☞ 𝐒ᴘᴇᴇᴅ** `%dms`\n"+
			"├• **╰☞ 𝐔ᴘᴛɪᴍᴇ** `%s`\n"+
			"└• **╰☞ 𝐍ᴀᴍᴇ:** %s",
		duration,
		uptime(),
		mention,
	)

	sent.Edit(msg, &telegram.SendOptions{ParseMode: telegram.HTML})
	return nil
}

// StartBotHandler — /start command (bot mode, public — no owner filter)
func StartBotHandler(m *telegram.NewMessage) error {
	msg := "✨ 𝗛ᴇʏ 𝗧ʜᴇʀᴇ..! 👋\n\n" +
		"ɪ'ᴍ <b>ᴘʙxɢᴏ</b> — ᴀ ꜰᴀsᴛ &amp; ᴘᴏᴡᴇʀꜰᴜʟ ɢᴏ ᴜsᴇʀʙᴏᴛ ⚡\n\n" +
		"❍ ᴍᴀɴᴀɢᴇ ʏᴏᴜʀ ᴜsᴇʀʙᴏᴛ ᴇᴀsɪʟʏ\n" +
		"❍ ꜰᴀsᴛ ᴘᴇʀꜰᴏʀᴍᴀɴᴄᴇ\n" +
		"❍ ʟɪɢʜᴛᴡᴇɪɢʜᴛ &amp; sᴍᴏᴏᴛʜ\n" +
		"❍ ᴘᴏᴡᴇʀᴇᴅ ʙʏ ɢᴏɢʀᴀᴍ 🔥\n\n" +
		"📖 <b>ᴄᴏᴍᴍᴀɴᴅs:</b>\n" +
		"• <code>.alive</code>\n" +
		"• <code>.ping</code>\n" +
		"• <code>.help</code>"

	m.Reply(msg, &telegram.SendOptions{ParseMode: telegram.HTML})
	return nil
}

func init() {
	Register(ModuleInfo{
		Name:        "Core",
		Description: "Basic health-check commands.",
		Commands: []CommandInfo{
			{Pattern: "alive", Handler: aliveHandler, Sudo: true},
			{Pattern: "ping", Handler: pingHandler, Sudo: true},
		},
	})
}
