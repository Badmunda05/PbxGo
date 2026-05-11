package user

import (
	"fmt"
	"os"
	"runtime"
	"time"

	"pbxgo/modules"

	"github.com/amarnathcjd/gogram/telegram"
)

var startTime = time.Now()

func uptime() string {
	d := time.Since(startTime).Round(time.Second)
	h, m, s := int(d.Hours()), int(d.Minutes())%60, int(d.Seconds())%60
	if h > 0 {
		return fmt.Sprintf("%dh %dm %ds", h, m, s)
	}
	if m > 0 {
		return fmt.Sprintf("%dm %ds", m, s)
	}
	return fmt.Sprintf("%ds", s)
}

func aliveHandler(m *telegram.NewMessage) error {
	me := m.Client.Me()
	mention := fmt.Sprintf(`<a href="tg://user?id=%d">%s</a>`, me.ID, me.FirstName)
	msg := fmt.Sprintf(
		"** 🥀 ᴘʙxɢᴏ 💞 **\n\n"+
			"❏ **ᴄʟɪᴇɴᴛ:** %s\n"+
			"├• **ᴠᴇʀsɪᴏɴ:** `v3.0.0`\n"+
			"├• **ᴜᴘᴛɪᴍᴇ:** `%s`\n"+
			"├• **ɢᴏ:** `%s`\n"+
			"└• **ᴏᴡɴᴇʀ:** [ʙᴀᴅ ᴍᴜɴᴅᴀ](https://t.me/Badmundaxd)",
		mention, uptime(), runtime.Version(),
	)
	modules.Reply(m, msg)
	return nil
}

func pingHandler(m *telegram.NewMessage) error {
	start := time.Now()
	sent, err := modules.Reply(m, "🏓 ᴘɪɴɢɪɴɢ...")
	if err != nil || sent == nil {
		return err
	}
	ms := time.Since(start).Milliseconds()
	me := m.Client.Me()
	mention := fmt.Sprintf(`<a href="tg://user?id=%d">%s</a>`, me.ID, me.FirstName)
	_, _ = sent.Edit(
		fmt.Sprintf("❏ **ᴘʙxɢᴏ**\n├• **sᴘᴇᴇᴅ** `%dms`\n├• **ᴜᴘᴛɪᴍᴇ** `%s`\n└• **ᴄʟɪᴇɴᴛ:** %s", ms, uptime(), mention),
		&telegram.SendOptions{ParseMode: telegram.HTML},
	)
	return nil
}

func restartHandler(m *telegram.NewMessage) error {
	modules.Reply(m, "♻️ <b>Restarting PbxGo...</b>")
	time.Sleep(time.Second)
	os.Exit(0)
	return nil
}

func init() {
	modules.RegisterUser(modules.Module{
		Name: "Core",
		Commands: []modules.CommandDef{
			{Pattern: "alive", Handler: aliveHandler},
			{Pattern: "ping", Handler: pingHandler},
			{Pattern: "restart", Handler: restartHandler},
			{Pattern: "rs", Handler: restartHandler},
		},
	})
}
