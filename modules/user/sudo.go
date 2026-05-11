package user

import (
	"fmt"
	"pbxgo/database"
	"pbxgo/modules"
	"strconv"
	"strings"

	"github.com/amarnathcjd/gogram/telegram"
)

func addSudoHandler(m *telegram.NewMessage) error {
	uid, err := resolveUID(m)
	if err != nil {
		modules.Reply(m, "⚠️ Usage: <code>.addsudo id</code> or reply")
		return nil
	}
	database.AddSudo(uid)
	modules.Reply(m, fmt.Sprintf("✅ <code>%d</code> added to sudo.", uid))
	return nil
}

func rmSudoHandler(m *telegram.NewMessage) error {
	uid, err := resolveUID(m)
	if err != nil {
		modules.Reply(m, "⚠️ Usage: <code>.rmsudo id</code> or reply")
		return nil
	}
	database.RemoveSudo(uid)
	modules.Reply(m, fmt.Sprintf("✅ <code>%d</code> removed from sudo.", uid))
	return nil
}

func sudoListHandler(m *telegram.NewMessage) error {
	list := database.FetchSudoList()
	if len(list) == 0 {
		modules.Reply(m, "📋 No sudo users.")
		return nil
	}
	var sb strings.Builder
	sb.WriteString("👥 <b>Sudo Users:</b>\n")
	for i, id := range list {
		fmt.Fprintf(&sb, "%d. <code>%d</code>\n", i+1, id)
	}
	modules.Reply(m, sb.String())
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
	args := modules.GetArgs(m)
	if args == "" {
		return 0, fmt.Errorf("no id")
	}
	return strconv.ParseInt(strings.TrimSpace(args), 10, 64)
}

func init() {
	modules.RegisterUser(modules.Module{
		Name: "Sudo",
		Commands: []modules.CommandDef{
			{Pattern: "addsudo", Handler: addSudoHandler, OwnerOnly: true},
			{Pattern: "asd", Handler: addSudoHandler, OwnerOnly: true},
			{Pattern: "rmsudo", Handler: rmSudoHandler, OwnerOnly: true},
			{Pattern: "delsudo", Handler: rmSudoHandler, OwnerOnly: true},
			{Pattern: "sudolist", Handler: sudoListHandler},
		},
	})
}
