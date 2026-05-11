package client

import (
	"fmt"
	"log/slog"
	"pbxgo/config"
	"pbxgo/modules"
	"pbxgo/modules/bot"

	_ "pbxgo/modules/user" // register user modules via init()

	"github.com/amarnathcjd/gogram/telegram"
)

// ─────────────────────────────────────────────────────────────
// sessionManagerImpl — implements bot.SessionManager
// ─────────────────────────────────────────────────────────────

type sessionManagerImpl struct{}

func (s *sessionManagerImpl) AddSession(sessionStr string) (int64, string, error) {
	return AddSessionRuntime(sessionStr)
}
func (s *sessionManagerImpl) RemoveSession(userID int64) error {
	return RemoveSessionRuntime(userID)
}
func (s *sessionManagerImpl) UserCount() int {
	mu.RLock()
	defer mu.RUnlock()
	return len(Users)
}
func (s *sessionManagerImpl) IsOnline(userID int64) bool {
	mu.RLock()
	defer mu.RUnlock()
	for _, u := range Users {
		if u.Me().ID == userID {
			return true
		}
	}
	return false
}
func (s *sessionManagerImpl) IsAuth(id int64) bool  { return IsAuth(id) }
func (s *sessionManagerImpl) IsOwner(id int64) bool { return IsOwner(id) }
func (s *sessionManagerImpl) AppID() int32          { return config.App.AppID }
func (s *sessionManagerImpl) AppHash() string       { return config.App.AppHash }

// ─────────────────────────────────────────────────────────────
// RegisterUserHandlers — attach user modules to a userbot client
// ─────────────────────────────────────────────────────────────

func RegisterUserHandlers(c *telegram.Client) {
	c.SetCommandPrefixes(". ! ?")
	for _, mod := range modules.UserModules {
		for _, cmd := range mod.Commands {
			if cmd.OwnerOnly {
				c.On("cmd:"+cmd.Pattern, cmd.Handler, ownerFilter(c))
			} else {
				c.On("cmd:"+cmd.Pattern, cmd.Handler, authFilter)
			}
		}
		slog.Info("User module registered", "module", mod.Name, "client", c.Me().FirstName)
	}
}

// ─────────────────────────────────────────────────────────────
// RegisterBotHandlers — attach bot modules + raw handlers
// ─────────────────────────────────────────────────────────────

func RegisterBotHandlers() {
	// Inject session manager (breaks import cycle)
	bot.SM = &sessionManagerImpl{}

	Bot.SetCommandPrefixes("/")

	// Command handlers
	for _, mod := range modules.BotModules {
		for _, cmd := range mod.Commands {
			Bot.On("cmd:"+cmd.Pattern, cmd.Handler, botAuthFilter)
		}
		slog.Info("Bot module registered", "module", mod.Name)
	}

	// Raw handlers (e.g. all private messages for conversation flow)
	for _, rh := range modules.BotRawHandlers {
		Bot.On(rh.Event, rh.Handler, privateFilter)
	}

	slog.Info("✅ All handlers registered",
		"user_modules", len(modules.UserModules),
		"bot_modules", len(modules.BotModules),
		"raw_handlers", len(modules.BotRawHandlers),
		"user_clients", len(Users),
	)
}

// ─────────────────────────────────────────────────────────────
// Filters
// ─────────────────────────────────────────────────────────────

func ownerFilter(c *telegram.Client) func(*telegram.NewMessage) error {
	return func(m *telegram.NewMessage) error {
		if m.SenderID() == c.Me().ID {
			return nil
		}
		return fmt.Errorf("owner only")
	}
}

func authFilter(m *telegram.NewMessage) error {
	if IsAuth(m.SenderID()) || m.SenderID() == m.Client.Me().ID {
		return nil
	}
	return fmt.Errorf("unauthorized")
}

func botAuthFilter(m *telegram.NewMessage) error {
	// /start is always accessible
	if m.Text() == "/start" {
		return nil
	}
	if IsAuth(m.SenderID()) {
		return nil
	}
	return fmt.Errorf("unauthorized")
}

// privateFilter — only allow private chats (not groups) for conversation handler
func privateFilter(m *telegram.NewMessage) error {
	if m.ChatID() == m.SenderID() { // private chat: chatID == userID
		return nil
	}
	return fmt.Errorf("not private")
}
