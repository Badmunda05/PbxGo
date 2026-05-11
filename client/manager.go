package client

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"pbxgo/config"
	"pbxgo/database"
	"sync"
	"syscall"

	"github.com/amarnathcjd/gogram/telegram"
)

// ─────────────────────────────────────────────
//  Global state
// ─────────────────────────────────────────────

var (
	Bot   *telegram.Client      // bot account (BotFather token)
	Users []*telegram.Client    // userbot accounts (session strings from DB)
	mu    sync.RWMutex
)

// ─────────────────────────────────────────────
//  Start bot
// ─────────────────────────────────────────────

func StartBot() error {
	var err error
	Bot, err = telegram.NewClient(telegram.ClientConfig{
		AppID:    config.App.AppID,
		AppHash:  config.App.AppHash,
		LogLevel: telegram.LogInfo,
		DeviceConfig: telegram.DeviceConfig{
			DeviceModel:   "PbxGo",
			SystemVersion: "Go 1.26",
			AppVersion:    "v3.0.0",
		},
	})
	if err != nil {
		return fmt.Errorf("bot client create: %w", err)
	}
	if _, err = Bot.Conn(); err != nil {
		return fmt.Errorf("bot connect: %w", err)
	}
	if err = Bot.LoginBot(config.App.BotToken); err != nil {
		return fmt.Errorf("bot login: %w", err)
	}
	me := Bot.Me()
	slog.Info("✅ Bot started", "username", me.Username, "id", me.ID)
	return nil
}

// ─────────────────────────────────────────────
//  Load all saved userbot sessions from MongoDB
// ─────────────────────────────────────────────

func LoadUserSessions() {
	sessions, err := database.GetAllSessions()
	if err != nil {
		slog.Error("Failed to load sessions", "error", err)
		return
	}
	for i, s := range sessions {
		if err := addUserClient(fmt.Sprintf("PbxUser#%d", i+1), s.Session); err != nil {
			slog.Error("Session start failed", "user_id", s.UserID, "error", err)
		}
	}
	slog.Info("✅ User sessions loaded", "count", len(Users))
}

// addUserClient — create + start one userbot client
func addUserClient(name, sessionString string) error {
	c, err := telegram.NewClient(telegram.ClientConfig{
		AppID:         config.App.AppID,
		AppHash:       config.App.AppHash,
		StringSession: sessionString,
		LogLevel:      telegram.LogWarn,
		DeviceConfig: telegram.DeviceConfig{
			DeviceModel:   "PbxGo",
			SystemVersion: "Go 1.26",
			AppVersion:    "v3.0.0",
		},
	})
	if err != nil {
		return fmt.Errorf("create: %w", err)
	}
	if _, err = c.Conn(); err != nil {
		return fmt.Errorf("connect: %w", err)
	}
	me := c.Me()
	slog.Info("User session started", "name", name, "user", me.FirstName, "id", me.ID)

	mu.Lock()
	Users = append(Users, c)
	mu.Unlock()

	// Register all userbot command handlers for this client
	RegisterUserHandlers(c)
	return nil
}

// ─────────────────────────────────────────────
//  Add a new session at runtime (called from bot panel)
// ─────────────────────────────────────────────

func AddSessionRuntime(sessionString string) (int64, string, error) {
	// Validate session by connecting temporarily
	c, err := telegram.NewClient(telegram.ClientConfig{
		AppID:         config.App.AppID,
		AppHash:       config.App.AppHash,
		StringSession: sessionString,
		LogLevel:      telegram.LogWarn,
	})
	if err != nil {
		return 0, "", fmt.Errorf("invalid session: %w", err)
	}
	if _, err = c.Conn(); err != nil {
		return 0, "", fmt.Errorf("session connect failed: %w", err)
	}
	me := c.Me()
	userID := me.ID
	firstName := me.FirstName

	// Save to DB
	if err = database.AddSession(userID, sessionString, firstName); err != nil {
		return 0, "", fmt.Errorf("db save failed: %w", err)
	}

	name := fmt.Sprintf("PbxUser#%d", len(Users)+1)
	// Register handlers and keep client alive
	mu.Lock()
	Users = append(Users, c)
	mu.Unlock()
	RegisterUserHandlers(c)

	slog.Info("New session added at runtime", "name", name, "user", firstName, "id", userID)
	return userID, firstName, nil
}

// ─────────────────────────────────────────────
//  Remove a session at runtime
// ─────────────────────────────────────────────

func RemoveSessionRuntime(userID int64) error {
	mu.Lock()
	defer mu.Unlock()
	for i, c := range Users {
		if c.Me().ID == userID {
			Users = append(Users[:i], Users[i+1:]...)
			slog.Info("Session removed", "user_id", userID)
			break
		}
	}
	return database.RemoveSession(userID)
}

// ─────────────────────────────────────────────
//  Auth check helpers
// ─────────────────────────────────────────────

func IsOwner(id int64) bool {
	return id == config.App.OwnerID
}

func IsAuth(id int64) bool {
	if IsOwner(id) {
		return true
	}
	for _, s := range config.App.SudoUsers {
		if s == id {
			return true
		}
	}
	return database.IsSudo(id)
}

// ─────────────────────────────────────────────
//  Run — block until SIGINT/SIGTERM
// ─────────────────────────────────────────────

func Run(_ context.Context) {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	slog.Info("Shutting down PbxGo...")
	database.Disconnect()
	os.Exit(0)
}
