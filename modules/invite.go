package modules

import (
	"fmt"
	"time"

	"github.com/amarnathcjd/gogram/telegram"
)

// .inviteall @group
func inviteAllHandler(m *telegram.NewMessage) error {

	args := GetArgs(m)

	if args == "" {
		Reply(m, "⚠️ Usage: <code>.inviteall @groupusername</code>")
		return nil
	}

	status, _ := Reply(m, "`👥 Fetching members...`")

	members, _, err := m.Client.GetChatMembers(
		args,
		&telegram.ParticipantOptions{
			Limit: 100,
		},
	)

	if err != nil {
		status.Edit("`❌ Failed to fetch members.`", nil)
		return nil
	}

	done := 0
	failed := 0

	for _, member := range members {

		if member.User == nil {
			continue
		}

		if member.User.Bot {
			continue
		}

		peer, err := m.Client.ResolvePeer(m.ChatID())

		if err != nil {
			failed++
			continue
		}

		userPeer, err := m.Client.ResolvePeer(member.User.ID)

		if err != nil {
			failed++
			continue
		}

		_, err = m.Client.API().ChannelsInviteToChannel(
			peer.InputChannel(),
			[]telegram.InputUser{
				userPeer.InputUser(),
			},
		)

		if err != nil {
			failed++
		} else {
			done++
		}

		time.Sleep(700 * time.Millisecond)
	}

	status.Edit(
		fmt.Sprintf(
			"✅ Invite Finished\n\n📨 Success: %d\n❌ Failed: %d",
			done,
			failed,
		),
		nil,
	)

	return nil
}

func init() {

	Register(ModuleInfo{
		Name:        "Invite",
		Description: "Invite members from another group",
		Commands: []CommandInfo{
			{
				Pattern: "inviteall",
				Handler: inviteAllHandler,
				Sudo:    false,
			},
			{
				Pattern: "invitesall",
				Handler: inviteAllHandler,
				Sudo:    false,
			},
		},
	})
}
