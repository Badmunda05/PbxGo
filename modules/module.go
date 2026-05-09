package modules

import "github.com/amarnathcjd/gogram/telegram"

// HandlerFunc is the standard handler signature for all PbxGo commands.
type HandlerFunc = func(*telegram.NewMessage) error

// CommandInfo defines a single bot command.
type CommandInfo struct {
	Pattern string
	Handler HandlerFunc
	Sudo    bool // true = owner OR sudo users can run; false = owner only
}

// ModuleInfo groups related commands under a named module.
type ModuleInfo struct {
	Name        string
	Description string
	Commands    []CommandInfo
}

// RegisteredModules holds all modules loaded via init().
var RegisteredModules []ModuleInfo

// Register adds a module to the global list.
func Register(m ModuleInfo) {
	RegisteredModules = append(RegisteredModules, m)
}
