package modules

import (
	"fmt"
	"pbxgo/database"
	"strconv"
	"strings"

	"github.com/amarnathcjd/gogram/telegram"
)

func addSudoHandler(m *telegram.NewMessage) error {
	userID, err := resolveUserID(m)
	if err != nil {
		Reply(m, "⚠️ Usage: <code>.addsudo &lt;user_id&gt;</code> or reply to a user")
		return nil
	}
	database.AddSudo(userID)
	Reply(m, fmt.Sprintf("✅ User <code>%d</code> added to sudo list.", userID))
	return nil
}

func rmSudoHandler(m *telegram.NewMessage) error {
	userID, err := resolveUserID(m)
	if err != nil {
		Reply(m, "⚠️ Usage: <code>.rmsudo &lt;user_id&gt;</code> or reply to a user")
		return nil
	}
	database.RemoveSudo(userID)
	Reply(m, fmt.Sprintf("✅ User <code>%d</code> removed from sudo list.", userID))
	return nil
}

func listSudoHandler(m *telegram.NewMessage) error {
	list := database.FetchSudoList()
	if len(list) == 0 {
		Reply(m, "📋 No sudo users found.")
		return nil
	}
	var sb strings.Builder
	sb.WriteString("👥 <b>Sudo Users:</b>\n")
	for _, id := range list {
		fmt.Fprintf(&sb, "• <code>%d</code>\n", id)
	}
	Reply(m, sb.String())
	return nil
}

func resolveUserID(m *telegram.NewMessage) (int64, error) {
	if m.IsReply() {
		r, err := m.GetReplyMessage()
		if err != nil {
			return 0, err
		}
		return r.SenderID(), nil
	}
	args := strings.Fields(m.Text())
	if len(args) < 2 {
		return 0, fmt.Errorf("no user id")
	}
	id, err := strconv.ParseInt(args[1], 10, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid id: %w", err)
	}
	return id, nil
}

func init() {
	Register(ModuleInfo{
		Name:        "Sudo",
		Description: "Manage sudo users (owner only).",
		Commands: []CommandInfo{
			{Pattern: "addsudo", Handler: addSudoHandler, Sudo: false},
			{Pattern: "rmsudo", Handler: rmSudoHandler, Sudo: false},
			{Pattern: "listsudo", Handler: listSudoHandler, Sudo: true},
		},
	})
}
