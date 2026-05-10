package modules

import "github.com/amarnathcjd/gogram/telegram"

func helpHandler(m *telegram.NewMessage) error {
	msg := `<b>🤖 PbxGo — Command List</b>
<code>━━━━━━━━━━━━━━━━━━━</code>

<b>⚙️ CORE</b>
<code>.alive</code> — Bot status check
<code>.ping</code> — Speed &amp; uptime
<code>.help</code> — This menu
<code>.restart</code> — Restart the bot

<b>📢 BROADCAST</b>
<code>.gcast</code> [text/reply] — Broadcast to all groups
<code>.gucast</code> [text/reply] — Broadcast to all private chats

<b>👤 CLONE</b>
<code>.clone</code> [reply/username] — Clone a user's profile
<code>.revert</code> — Revert to original profile

<b>✏️ EMOJI</b>
<code>.emoji</code> [text] — Convert text to emoji art
<code>.cmoji</code> [emoji] [text] — Custom emoji art

<b>🔨 PURGE</b>
<code>.del</code> — Delete replied message
<code>.purge</code> — Purge from replied message
<code>.purgeme</code> [n] — Delete your last n messages

<b>💬 SPAM</b>
<code>.spam</code> [count] [text] — Spam messages
<code>.delayspam</code> [delay] [count] [text] — Delayed spam
<code>.sspam</code> [count] — Sticker spam (reply to sticker)

<b>🏓 HANG</b>
<code>.hang</code> [count] — Send hang messages

<b>👥 TAGGER</b>
<code>.all</code> [text] — Mention all members
<code>.cancel</code> — Stop tagger

<b>⚔️ RAID</b>
<code>.raid</code> [count] — Bold raid
<code>.hraid</code> [count] — Hindi raid
<code>.eraid</code> [count] — English raid
<code>.punraid</code> [count] — Punjabi raid
<code>.replyraid</code> [count] — Reply raid
<code>.hreplyraid</code> [count] — Hindi reply raid
<code>.ereplyraid</code> [count] — English reply raid
<code>.preplyraid</code> [count] — Punjabi reply raid

<b>📥 MEDIA</b>
<code>.save</code> — Download replied media

<b>🔐 SUDO</b>
<code>.addsudo</code> [id/reply] — Add sudo user
<code>.rmsudo</code> [id/reply] — Remove sudo user
<code>.sudolist</code> — List sudo users

<code>━━━━━━━━━━━━━━━━━━━</code>
<b>Support:</b> @PBXCHATS | <b>Updates:</b> @PBX_UPDATE`

	m.Delete()
	Reply(m, msg)
	return nil
}

func init() {
	Register(ModuleInfo{
		Name:        "Help",
		Description: "Help menu",
		Commands: []CommandInfo{
			{Pattern: "help", Handler: helpHandler, Sudo: true},
		},
	})
}

var _ = telegram.HTML
