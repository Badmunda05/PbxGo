package modules

import (
	"fmt"
	"pbxgo/database"
	"strconv"
	"strings"

	"github.com/amarnathcjd/gogram/telegram"
)

func addSudoHandler(m *telegram.NewMessage) error {
	userID, err := resolveUID(m)
	if err != nil {
		Reply(m, "⚠️ Usage: <code>.addsudo &lt;user_id&gt;</code> or reply to a user")
		return nil
	}
	database.AddSudo(userID)
	Reply(m, fmt.Sprintf("✅ User <code>%d</code> added to sudo list.", userID))
	return nil
}

func rmSudoHandler(m *telegram.NewMessage) error {
	userID, err := resolveUID(m)
	if err != nil {
		Reply(m, "⚠️ Usage: <code>.rmsudo &lt;user_id&gt;</code> or reply to a user")
		return nil
	}
	database.RemoveSudo(userID)
	Reply(m, fmt.Sprintf("✅ User <code>%d</code> removed from sudo list.", userID))
	return nil
}

func sudoListHandler(m *telegram.NewMessage) error {
	ex, _ := Reply(m, "`⏳ Processing...`")
	list := database.FetchSudoList()
	if len(list) == 0 {
		ex.Edit("📋 No sudo users found.", nil)
		return nil
	}
	var sb strings.Builder
	sb.WriteString("👥 <b>Sudo Users:</b>\n")
	for i, id := range list {
		fmt.Fprintf(&sb, "%d. <code>%d</code>\n", i+1, id)
	}
	ex.Edit(sb.String(), &telegram.SendOptions{ParseMode: telegram.HTML})
	return nil
}

func resolveUID(m *telegram.NewMessage) (int64, error) {
	if m.IsReply() {
		r, err := m.GetReplyMessage()
		if err != nil {
			return 0, err
		}
		return r.SenderID(), nil
	}
	args := GetArgs(m)
	if args == "" {
		return 0, fmt.Errorf("no id provided")
	}
	id, err := strconv.ParseInt(strings.TrimSpace(args), 10, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid id")
	}
	return id, nil
}

func init() {
	Register(ModuleInfo{
		Name:        "Sudo",
		Description: "Manage sudo users",
		Commands: []CommandInfo{
			// Sudo: false = owner only (add/remove is sensitive, owner-only)
			{Pattern: "addsudo", Handler: addSudoHandler, Sudo: false},
			{Pattern: "asd", Handler: addSudoHandler, Sudo: false},
			{Pattern: "rmsudo", Handler: rmSudoHandler, Sudo: false},
			{Pattern: "delsudo", Handler: rmSudoHandler, Sudo: false},
			// Sudo: true = owner + sudo users can view the list
			{Pattern: "sudolist", Handler: sudoListHandler, Sudo: true},
			{Pattern: "sdl", Handler: sudoListHandler, Sudo: true},
		},
	})
}
