package modules

import (
	"github.com/amarnathcjd/gogram/telegram"
)

// .inviteall disabled
func inviteAllHandler(m *telegram.NewMessage) error {

	Reply(
		m,
		"❌ Invite module is not supported in this gogram version.",
	)

	return nil
}

func init() {

	Register(ModuleInfo{
		Name:        "Invite",
		Description: "Invite module disabled",
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
