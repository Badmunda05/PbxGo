package bot

import (
	"fmt"
	"pbxgo/database"
	"pbxgo/modules"
	"sync"
	"time"

	"github.com/amarnathcjd/gogram/telegram"
)

// ─────────────────────────────────────────────────────────────
// SessionManager interface — injected by client to avoid cycle
// ─────────────────────────────────────────────────────────────

type SessionManager interface {
	AddSession(sessionStr string) (userID int64, firstName string, err error)
	RemoveSession(userID int64) error
	UserCount() int
	IsOnline(userID int64) bool
	IsAuth(id int64) bool
	IsOwner(id int64) bool
	AppID() int32
	AppHash() string
}

var SM SessionManager // injected by client.RegisterBotHandlers()

// ─────────────────────────────────────────────────────────────
// Pending OTP channels — one per active login flow
// ─────────────────────────────────────────────────────────────

type pendingFlow struct {
	codeCh chan string // receives OTP or password
	step   string     // "otp" or "2fa"
}

var flows sync.Map // chatID int64 → *pendingFlow

// ─────────────────────────────────────────────────────────────
// /start
// ─────────────────────────────────────────────────────────────

func startHandler(m *telegram.NewMessage) error {
	if SM == nil || !SM.IsAuth(m.SenderID()) {
		m.Reply("👋 Hey! This is a private bot.", nil)
		return nil
	}
	msg := fmt.Sprintf(
		"✨ <b>ᴘʙxɢᴏ v3.0.0</b> 🔥\n\n"+
			"<b>🤖 Active Sessions:</b> <code>%d</code>\n\n"+
			"<b>📋 Commands:</b>\n"+
			"• /newsession — Generate session (phone+OTP)\n"+
			"• /add <code>string</code> — Add manually\n"+
			"• /sessions — List all\n"+
			"• /remove <code>user_id</code> — Remove\n"+
			"• /help — All commands\n\n"+
			"<b>Support:</b> @PBXCHATS | @PBX_UPDATE",
		SM.UserCount(),
	)
	m.Reply(msg, &telegram.SendOptions{ParseMode: telegram.HTML})
	return nil
}

// ─────────────────────────────────────────────────────────────
// /newsession — interactive phone+OTP+2FA flow
// Uses gogram's Login() with CodeCallback & PasswordCallback
// that wait on channels fed by the text message handler
// ─────────────────────────────────────────────────────────────

func newSessionHandler(m *telegram.NewMessage) error {
	if SM == nil || !SM.IsAuth(m.SenderID()) {
		return nil
	}
	chatID := m.ChatID()

	// Cancel any existing flow for this chat
	if v, ok := flows.Load(chatID); ok {
		close(v.(*pendingFlow).codeCh)
		flows.Delete(chatID)
	}

	m.Reply(
		"👻 <b>ɴᴇᴡ sᴇssɪᴏɴ sᴇᴛᴜᴘ</b>\n\n"+
			"<b>Step 1</b> — Enter your phone number with country code:\n\n"+
			"<b>Example:</b> <code>+919876543210</code>\n\n"+
			"<i>Send /cancel to stop at any time.</i>",
		&telegram.SendOptions{ParseMode: telegram.HTML},
	)

	// Wait for phone number (next message)
	phone := waitMsg(m, chatID, "phone", 120)
	if phone == "" {
		return nil
	}
	if !isValidPhone(phone) {
		m.Reply("❌ Invalid phone number. Must start with <code>+</code> country code.\n\nTry /newsession again.",
			&telegram.SendOptions{ParseMode: telegram.HTML})
		return nil
	}

	prog, _ := m.Reply("`⏳ Sending OTP to Telegram...`", &telegram.SendOptions{ParseMode: telegram.MarkDown})

	// Create temp client
	c, err := telegram.NewClient(telegram.ClientConfig{
		AppID:         SM.AppID(),
		AppHash:       SM.AppHash(),
		MemorySession: true,
		LogLevel:      telegram.LogError,
	})
	if err != nil {
		editMsg(prog, "❌ Client error: <code>"+err.Error()+"</code>")
		return nil
	}
	if _, err = c.Conn(); err != nil {
		editMsg(prog, "❌ Connection failed: <code>"+err.Error()+"</code>")
		return nil
	}

	// Create channel for OTP/2FA
	codeCh := make(chan string, 1)
	pf := &pendingFlow{codeCh: codeCh, step: "otp"}
	flows.Store(chatID, pf)

	// Run login in goroutine (it blocks waiting for callbacks)
	doneCh := make(chan error, 1)
	go func() {
		loginErr := c.Login(phone, &telegram.LoginOptions{
			CodeCallback: func() (string, error) {
				editMsg(prog,
					"✅ <b>OTP Sent!</b>\n\n"+
						"<b>Step 2</b> — Enter the OTP from your Telegram app.\n"+
						"Separate digits with spaces:\n"+
						"<b>Example:</b> <code>2 4 1 7 4</code>\n\n"+
						"<i>Send /cancel to stop.</i>",
				)
				pf.step = "otp"
				code, ok := <-codeCh
				if !ok {
					return "", fmt.Errorf("cancelled")
				}
				return removeSpaces(code), nil
			},
			PasswordCallback: func() (string, error) {
				editMsg(prog,
					"🔐 <b>2FA Required</b>\n\n"+
						"<b>Step 3</b> — Enter your Two-Step Verification password:\n\n"+
						"<i>Send /cancel to stop.</i>",
				)
				pf.step = "2fa"
				// Reset channel for password
				codeCh2 := make(chan string, 1)
				pf.codeCh = codeCh2
				flows.Store(chatID, pf)
				pass, ok := <-codeCh2
				if !ok {
					return "", fmt.Errorf("cancelled")
				}
				return pass, nil
			},
		})
		doneCh <- loginErr
	}()

	// Wait for login to complete (max 5 min)
	select {
	case loginErr := <-doneCh:
		flows.Delete(chatID)
		if loginErr != nil {
			if loginErr.Error() == "cancelled" {
				editMsg(prog, "❌ <b>Cancelled.</b>")
			} else {
				editMsg(prog, "❌ Login failed: <code>"+loginErr.Error()+"</code>")
			}
			return nil
		}

		// Export & save
		sessionStr := c.ExportSession()
		_ = c.Stop()

		editMsg(prog, "✅ <b>Login successful!</b> Saving to database...")
		time.Sleep(500 * time.Millisecond)

		userID, firstName, err := SM.AddSession(sessionStr)
		if err != nil {
			editMsg(prog, "❌ Failed to save: <code>"+err.Error()+"</code>")
			return nil
		}

		// Send backup to Saved Messages
		go func() {
			time.Sleep(time.Second)
			m.Client.SendMessage(m.SenderID(),
				fmt.Sprintf("✅ <b>PbxGo Session Backup</b>\n\n<code>%s</code>\n\n⚠️ <b>Keep secret!</b>", sessionStr),
				&telegram.SendOptions{ParseMode: telegram.HTML},
			)
		}()

		editMsg(prog, fmt.Sprintf(
			"✅ <b>Session Added!</b>\n\n"+
				"👤 <b>Name:</b> %s\n"+
				"🆔 <b>ID:</b> <code>%d</code>\n"+
				"📊 <b>Active:</b> <code>%d</code>\n\n"+
				"<i>Userbot is online! 📨 Session sent to Saved Messages.</i>",
			firstName, userID, SM.UserCount(),
		))

	case <-time.After(5 * time.Minute):
		flows.Delete(chatID)
		editMsg(prog, "⏰ <b>Timeout!</b> Session setup cancelled.\n\nTry /newsession again.")
	}

	return nil
}

// waitMsg — block until next private message from user (or /cancel or timeout)
func waitMsg(m *telegram.NewMessage, chatID int64, label string, timeoutSec int) string {
	ch := make(chan string, 1)
	waiters.Store(chatID, ch)
	defer waiters.Delete(chatID)

	select {
	case text := <-ch:
		if text == "/cancel" {
			m.Reply("❌ <b>Cancelled!</b>", &telegram.SendOptions{ParseMode: telegram.HTML})
			return ""
		}
		return text
	case <-time.After(time.Duration(timeoutSec) * time.Second):
		m.Reply("⏰ <b>Timeout!</b> Try /newsession again.", &telegram.SendOptions{ParseMode: telegram.HTML})
		return ""
	}
}

// waiters — chats waiting for a phone number (step 1 only)
var waiters sync.Map // chatID → chan string

// ─────────────────────────────────────────────────────────────
// Text message handler — feeds OTP/pass into login goroutine
// ─────────────────────────────────────────────────────────────

func textHandler(m *telegram.NewMessage) error {
	if SM == nil || !SM.IsAuth(m.SenderID()) {
		return nil
	}
	chatID := m.ChatID()
	text := m.Text()

	// Skip commands (handled by command handlers)
	if len(text) > 0 && text[0] == '/' && !isCancelCmd(text) {
		return nil
	}

	// /cancel — stop any active flow
	if isCancelCmd(text) {
		cancelled := false
		if v, ok := flows.Load(chatID); ok {
			close(v.(*pendingFlow).codeCh)
			flows.Delete(chatID)
			cancelled = true
		}
		if v, ok := waiters.Load(chatID); ok {
			v.(chan string) <- "/cancel"
			cancelled = true
		}
		if cancelled {
			m.Reply("❌ <b>Cancelled!</b>", &telegram.SendOptions{ParseMode: telegram.HTML})
		}
		return nil
	}

	// Feed phone number waiter (step 1)
	if v, ok := waiters.Load(chatID); ok {
		v.(chan string) <- text
		return nil
	}

	// Feed OTP / 2FA into active login flow (step 2 / 3)
	if v, ok := flows.Load(chatID); ok {
		pf := v.(*pendingFlow)
		select {
		case pf.codeCh <- text:
		default:
		}
		return nil
	}

	return nil
}

func isCancelCmd(text string) bool {
	return text == "/cancel" || text == "/cancel@" || len(text) > 7 && text[:7] == "/cancel"
}

// ─────────────────────────────────────────────────────────────
// /add <session_string> — manual add
// ─────────────────────────────────────────────────────────────

func addSessionHandler(m *telegram.NewMessage) error {
	if SM == nil || !SM.IsAuth(m.SenderID()) {
		return nil
	}
	sessionStr := modules.GetArgs(m)
	if sessionStr == "" {
		m.Reply("⚠️ Usage: <code>/add YOUR_SESSION_STRING</code>", &telegram.SendOptions{ParseMode: telegram.HTML})
		return nil
	}
	prog, _ := m.Reply("`⏳ Validating...`", &telegram.SendOptions{ParseMode: telegram.MarkDown})
	userID, firstName, err := SM.AddSession(sessionStr)
	if err != nil {
		editMsg(prog, "❌ <b>Failed:</b> <code>"+err.Error()+"</code>")
		return nil
	}
	editMsg(prog, fmt.Sprintf(
		"✅ <b>Session Added!</b>\n\n👤 %s\n🆔 <code>%d</code>\n📊 Active: <code>%d</code>",
		firstName, userID, SM.UserCount(),
	))
	return nil
}

// ─────────────────────────────────────────────────────────────
// /sessions
// ─────────────────────────────────────────────────────────────

func listSessionsHandler(m *telegram.NewMessage) error {
	if SM == nil || !SM.IsAuth(m.SenderID()) {
		return nil
	}
	sessions, err := database.GetAllSessions()
	if err != nil || len(sessions) == 0 {
		m.Reply("📋 No sessions found.", nil)
		return nil
	}
	msg := fmt.Sprintf("👥 <b>Sessions (%d):</b>\n\n", len(sessions))
	for i, s := range sessions {
		icon := "🔴"
		if SM.IsOnline(s.UserID) {
			icon = "🟢"
		}
		msg += fmt.Sprintf("%s <b>%d.</b> %s — <code>%d</code>\n", icon, i+1, s.FirstName, s.UserID)
	}
	msg += "\n<code>🟢 online  🔴 offline</code>"
	m.Reply(msg, &telegram.SendOptions{ParseMode: telegram.HTML})
	return nil
}

// ─────────────────────────────────────────────────────────────
// /remove
// ─────────────────────────────────────────────────────────────

func removeSessionHandler(m *telegram.NewMessage) error {
	if SM == nil || !SM.IsOwner(m.SenderID()) {
		m.Reply("⚠️ Only owner can remove sessions.", nil)
		return nil
	}
	var userID int64
	fmt.Sscanf(modules.GetArgs(m), "%d", &userID)
	if userID == 0 {
		m.Reply("⚠️ Usage: <code>/remove USER_ID</code>", &telegram.SendOptions{ParseMode: telegram.HTML})
		return nil
	}
	if err := SM.RemoveSession(userID); err != nil {
		m.Reply("❌ <code>"+err.Error()+"</code>", &telegram.SendOptions{ParseMode: telegram.HTML})
		return nil
	}
	m.Reply(fmt.Sprintf("✅ Session <code>%d</code> removed.", userID), &telegram.SendOptions{ParseMode: telegram.HTML})
	return nil
}

// ─────────────────────────────────────────────────────────────
// /help
// ─────────────────────────────────────────────────────────────

func helpHandler(m *telegram.NewMessage) error {
	if SM == nil || !SM.IsAuth(m.SenderID()) {
		return nil
	}
	msg := `<b>🤖 PbxGo v3.0 — Commands</b>
<code>━━━━━━━━━━━━━━━━━━━</code>

<b>🔑 SESSION PANEL (Bot)</b>
<code>/newsession</code> — Add via phone+OTP (recommended)
<code>/add</code> [string] — Add manually
<code>/sessions</code> — List sessions
<code>/remove</code> [id] — Remove session
<code>/cancel</code> — Cancel active flow

<b>⚙️ USERBOT (. prefix)</b>
<code>.alive</code> / <code>.ping</code> / <code>.restart</code>

<b>💬 SPAM</b>
<code>.spam</code> [n] [text] / <code>.ds</code> [d] [n] [text] / <code>.sspam</code> [n]

<b>⚔️ RAID</b>
<code>.raid</code> / <code>.hraid</code> / <code>.eraid</code> / <code>.punraid</code> [n]
<code>.replyraid</code> / <code>.hreplyraid</code> / <code>.ereplyraid</code> / <code>.preplyraid</code> [n]

<b>🔨 PURGE</b>
<code>.del</code> / <code>.purge</code> / <code>.purgeme</code> [n]

<b>👥 TAGGER</b>
<code>.all</code> [text] / <code>.cancel</code>

<b>📢 OTHERS</b>
<code>.hang</code> [n] / <code>.gcast</code> / <code>.save</code>

<b>🔐 SUDO</b>
<code>.addsudo</code> / <code>.rmsudo</code> / <code>.sudolist</code>

<code>━━━━━━━━━━━━━━━━━━━</code>
@PBXCHATS | @PBX_UPDATE`
	m.Reply(msg, &telegram.SendOptions{ParseMode: telegram.HTML})
	return nil
}

// ─────────────────────────────────────────────────────────────
// Helpers
// ─────────────────────────────────────────────────────────────

func editMsg(msg *telegram.NewMessage, text string) {
	if msg == nil {
		return
	}
	_, _ = msg.Edit(text, &telegram.SendOptions{ParseMode: telegram.HTML})
}

func isValidPhone(p string) bool {
	if len(p) < 7 || p[0] != '+' {
		return false
	}
	for _, c := range p[1:] {
		if c < '0' || c > '9' {
			return false
		}
	}
	return true
}

func removeSpaces(s string) string {
	out := make([]byte, 0, len(s))
	for i := 0; i < len(s); i++ {
		if s[i] != ' ' {
			out = append(out, s[i])
		}
	}
	return string(out)
}

// ─────────────────────────────────────────────────────────────
// Registration
// ─────────────────────────────────────────────────────────────

func init() {
	modules.RegisterBot(modules.Module{
		Name: "Sessions",
		Commands: []modules.CommandDef{
			{Pattern: "start", Handler: startHandler},
			{Pattern: "newsession", Handler: newSessionHandler},
			{Pattern: "add", Handler: addSessionHandler},
			{Pattern: "sessions", Handler: listSessionsHandler},
			{Pattern: "remove", Handler: removeSessionHandler},
			{Pattern: "help", Handler: helpHandler},
		},
	})

	// Raw message handler — catches OTP/2FA inputs
	modules.RegisterBotRaw("message", textHandler)
}
