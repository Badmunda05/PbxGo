package bot

import (
	"pbxgo/modules"

	"github.com/amarnathcjd/gogram/telegram"
)

// Simple test command
func testHandler(m *telegram.NewMessage) error {

	_, _ = m.Reply(
		"✅ Bot Working Successfully!",
		nil,
	)

	return nil
}

func init() {

	modules.RegisterBot(modules.Module{
		Name: "Test",
		Commands: []modules.CommandDef{
			{
				Pattern: "test",
				Handler: testHandler,
			},
		},
	})
}
