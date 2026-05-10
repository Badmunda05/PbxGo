package modules

import (
	"fmt"
	"time"

	"github.com/amarnathcjd/gogram/telegram"
)

// .inviteall [@group] — dusre group de members nu current group vich invite karo
func inviteAllHandler(m *telegram.NewMessage) error {
	args := GetArgs(m)
	if args == "" {
		Reply(m, "⚠️ Usage: <code>.inviteall @groupusername</code>")
		return nil
	}

	status, _ := Reply(m, "`👥 Fetching members...`")

	members, err := m.Client.GetChatMembers(args, &telegram.GetChatMembersParams{
		Limit: 200,
	})
	if err != nil {
		status.Edit("`❌ Failed to get members from that group.`", nil)
		return nil
	}

	status.Edit("`📨 Inviting members...`", nil)
	done, failed := 0, 0

	for _, member := range members {
		user := member.GetUser()
		if user == nil || user.Bot {
			continue
		}
		err := m.Client.AddChatMembers(m.ChatID(), []int64{user.ID})
		if err != nil {
			failed++
		} else {
			done++
		}
		time.Sleep(500 * time.Millisecond)
	}

	msg := fmt.Sprintf(
		"✅ <b>Invite done!</b>\n📨 Invited: <code>%d</code>\n❌ Failed: <code>%d</code>",
		done, failed,
	)
	status.Edit(msg, &telegram.SendOptions{ParseMode: telegram.HTML})
	return nil
}

func init() {
	Register(ModuleInfo{
		Name:        "Invite",
		Description: "Invite members from another group",
		Commands: []CommandInfo{
			{Pattern: "inviteall", Handler: inviteAllHandler, Sudo: false},
			{Pattern: "invitesall", Handler: inviteAllHandler, Sudo: false},
		},
	})
}
