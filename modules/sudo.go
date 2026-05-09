package modules

import (
	"fmt"
	"main/database"
	"strconv"
	"strings"

	"github.com/amarnathcjd/gogram/telegram"
)

func AddSudoHandler(m *telegram.NewMessage) error {
	args := strings.Fields(m.Text())
	if len(args) < 2 && !m.IsReply() {
		EditOrReply(m, "⚠️ Usage: <code>.addsudo &lt;user_id&gt;</code>")
		return nil
	}

	var userID int64
	if m.IsReply() {
		r, err := m.GetReplyMessage()
		if err != nil {
			EditOrReply(m, "❌ Failed to get replied message.")
			return nil
		}
		userID = r.SenderID()
	} else {
		userIDStr := args[1]
		userIDx, err := strconv.ParseInt(userIDStr, 10, 64)
		if err != nil {
			EditOrReply(m, "❌ Invalid user ID.")
			return nil
		}
		userID = userIDx
	}

	database.AddSudo(userID)
	EditOrReply(m, fmt.Sprintf("✅ User <code>%d</code> added to sudo list.", userID))
	return nil
}

func RemoveSudoHandler(m *telegram.NewMessage) error {
	args := strings.Fields(m.Text())
	if len(args) < 2 && !m.IsReply() {
		EditOrReply(m, "⚠️ Usage: <code>.rmsudo &lt;user_id&gt;</code>")
		return nil
	}

	var userID int64
	if m.IsReply() {
		r, err := m.GetReplyMessage()
		if err != nil {
			EditOrReply(m, "❌ Failed to get replied message.")
			return nil
		}
		userID = r.SenderID()
	} else {
		userIDStr := args[1]
		userIDx, err := strconv.ParseInt(userIDStr, 10, 64)
		if err != nil {
			EditOrReply(m, "❌ Invalid user ID.")
			return nil
		}
		userID = userIDx
	}

	database.RemoveSudo(userID)
	EditOrReply(m, fmt.Sprintf("✅ User <code>%d</code> removed from sudo list.", userID))
	return nil
}

func ListSudoHandler(m *telegram.NewMessage) error {
	sudos := database.FetchSudoList()
	if len(sudos) == 0 {
		EditOrReply(m, "📋 No sudo users found.")
		return nil
	}

	list := "👥 <b>Sudo Users:</b>\n"
	for _, id := range sudos {
		list += fmt.Sprintf("• <code>%d</code>\n", id)
	}
	EditOrReply(m, list)
	return nil
}

func init() {
	RegisterModule(ModuleInfo{
		Name:        "Sudo Module",
		Description: "Manages sudo users for PbxGo.",
		Commands: []CommandInfo{
			{
				Pattern: "addsudo",
				Func:    AddSudoHandler,
				Sudo:    false,
			},
			{
				Pattern: "rmsudo",
				Func:    RemoveSudoHandler,
				Sudo:    false,
			},
			{
				Pattern: "listsudo",
				Func:    ListSudoHandler,
				Sudo:    true,
			},
		},
	})
}
