package modules

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/amarnathcjd/gogram/telegram"
)

const originalDataFile = "original_profile.json"

type profileBackup struct {
	FirstName  string `json:"first_name"`
	LastName   string `json:"last_name"`
	Bio        string `json:"bio"`
	Username   string `json:"username"`
	PhotoPaths []string `json:"photos"`
}

func loadBackup() map[string]profileBackup {
	data := make(map[string]profileBackup)
	if b, err := os.ReadFile(originalDataFile); err == nil {
		_ = json.Unmarshal(b, &data)
	}
	return data
}

func saveBackup(data map[string]profileBackup) {
	b, _ := json.MarshalIndent(data, "", "  ")
	_ = os.WriteFile(originalDataFile, b, 0644)
}

// .clone [username/reply]
func cloneHandler(m *telegram.NewMessage) error {
	op, _ := Reply(m, "`🔍 Processing...`")

	args := GetArgs(m)
	var targetID int64

	if m.IsReply() {
		r, err := m.GetReplyMessage()
		if err != nil {
			op.Edit("`❌ Could not get replied message.`", nil)
			return nil
		}
		targetID = r.SenderID()
	} else if args != "" {
		user, err := m.Client.GetUser(args)
		if err != nil {
			op.Edit("`❌ User not found.`", nil)
			return nil
		}
		targetID = user.ID
	} else {
		op.Edit("`⚠️ Usage: .clone @username or reply to a user`", nil)
		return nil
	}

	me := m.Client.Me()
	myID := fmt.Sprintf("%d", me.ID)

	// Backup if not done yet
	backups := loadBackup()
	if _, exists := backups[myID]; !exists {
		op.Edit("`📦 Backing up your profile...`", nil)
		backup := profileBackup{
			FirstName: me.FirstName,
			LastName:  me.LastName,
			Username:  me.Username,
		}
		backups[myID] = backup
		saveBackup(backups)
	}

	// Get target user info
	targetUser, err := m.Client.GetUser(targetID)
	if err != nil {
		op.Edit("`❌ Could not fetch target user.`", nil)
		return nil
	}

	op.Edit("`✏️ Updating profile...`", nil)

	// Update name & bio
	_ = m.Client.UpdateProfile(&telegram.UpdateProfileParams{
		FirstName: targetUser.FirstName,
		LastName:  targetUser.LastName,
	})

	op.Edit(
		fmt.Sprintf(
			"**✅ Cloned Successfully!**\n\n"+
				"**👤 Name:** `%s %s`\n"+
				"**🔗 Username:** `@%s`\n\n"+
				"_Your original data is safely backed up!_",
			targetUser.FirstName, targetUser.LastName,
			targetUser.Username,
		),
		&telegram.SendOptions{ParseMode: telegram.Markdown},
	)
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
		op.Edit("`⚠️ No backup found!`", nil)
		return nil
	}

	_ = m.Client.UpdateProfile(&telegram.UpdateProfileParams{
		FirstName: orig.FirstName,
		LastName:  orig.LastName,
	})

	delete(backups, myID)
	saveBackup(backups)

	op.Edit(
		fmt.Sprintf(
			"**✅ Reverted Successfully!**\n\n"+
				"**👤 Name:** `%s %s`\n"+
				"**🔗 Username:** `@%s`\n\n"+
				"_✨ Backup cleaned!_",
			orig.FirstName, orig.LastName, orig.Username,
		),
		&telegram.SendOptions{ParseMode: telegram.Markdown},
	)
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
