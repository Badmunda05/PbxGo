package modules

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/amarnathcjd/gogram/telegram"
)

const backupFile = "pbxgo_profile_backup.json"

type profileBackup struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

func loadBackup() map[string]profileBackup {
	data := make(map[string]profileBackup)
	if b, err := os.ReadFile(backupFile); err == nil {
		_ = json.Unmarshal(b, &data)
	}
	return data
}

func saveBackup(data map[string]profileBackup) {
	b, _ := json.MarshalIndent(data, "", "  ")
	_ = os.WriteFile(backupFile, b, 0644)
}

// .clone [@username / reply]
func cloneHandler(m *telegram.NewMessage) error {
	op, _ := Reply(m, "`🔍 Fetching user info...`")

	var targetUser *telegram.UserObj

	if m.IsReply() {
		r, err := m.GetReplyMessage()
		if err != nil {
			op.Edit("`❌ Could not get replied message.`", nil)
			return nil
		}
		u, err := m.Client.GetUser(r.SenderID())
		if err != nil {
			op.Edit("`❌ Could not fetch sender info.`", nil)
			return nil
		}
		targetUser = u
	} else {
		args := GetArgs(m)
		if args == "" {
			op.Edit("`⚠️ Usage: .clone @username or reply to a user`", nil)
			return nil
		}
		// Strip @ if present
		args = strings.TrimPrefix(args, "@")
		u, err := m.Client.GetUser(args)
		if err != nil {
			op.Edit("`❌ User not found.`", nil)
			return nil
		}
		targetUser = u
	}

	me := m.Client.Me()
	myID := fmt.Sprintf("%d", me.ID)

	// Backup original profile once
	backups := loadBackup()
	if _, exists := backups[myID]; !exists {
		op.Edit("`📦 Backing up your original profile...`", nil)
		backups[myID] = profileBackup{
			FirstName: me.FirstName,
			LastName:  me.LastName,
		}
		saveBackup(backups)
	}

	op.Edit("`✏️ Cloning profile...`", nil)

	// Update name via accounts.UpdateProfile
	_, err := m.Client.UpdateProfile(targetUser.FirstName, targetUser.LastName, "")
	if err != nil {
		op.Edit(fmt.Sprintf("`❌ Failed to update profile: %s`", err.Error()), nil)
		return nil
	}

	op.Edit(fmt.Sprintf(
		"**✅ Clone Successful!**\n\n"+
			"**👤 Name:** `%s %s`\n"+
			"**🔗 Username:** `@%s`\n\n"+
			"_Your original profile is backed up. Use_ `.revert` _to restore._",
		targetUser.FirstName, targetUser.LastName, targetUser.Username,
	), &telegram.SendOptions{ParseMode: telegram.MarkDown})
	return nil
}

// .revert
func revertHandler(m *telegram.NewMessage) error {
	op, _ := Reply(m, "`🔄 Reverting to original profile...`")

	me := m.Client.Me()
	myID := fmt.Sprintf("%d", me.ID)
	backups := loadBackup()

	orig, exists := backups[myID]
	if !exists {
		op.Edit("`⚠️ No backup found! Clone first with .clone`", nil)
		return nil
	}

	_, err := m.Client.UpdateProfile(orig.FirstName, orig.LastName, "")
	if err != nil {
		op.Edit(fmt.Sprintf("`❌ Failed to revert: %s`", err.Error()), nil)
		return nil
	}

	delete(backups, myID)
	saveBackup(backups)

	op.Edit(fmt.Sprintf(
		"**✅ Reverted Successfully!**\n\n"+
			"**👤 Name:** `%s %s`\n\n"+
			"_Backup cleared. You're back to original!_",
		orig.FirstName, orig.LastName,
	), &telegram.SendOptions{ParseMode: telegram.MarkDown})
	return nil
}

func init() {
	Register(ModuleInfo{
		Name:        "Clone",
		Description: "Clone/revert Telegram profile",
		Commands: []CommandInfo{
			{Pattern: "clone", Handler: cloneHandler, Sudo: false},
			{Pattern: "revert", Handler: revertHandler, Sudo: false},
		},
	})
}
