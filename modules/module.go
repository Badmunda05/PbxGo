package modules

import "github.com/amarnathcjd/gogram/telegram"

type HandlerFunc = func(*telegram.NewMessage) error

type CommandInfo struct {
	Pattern string
	Handler HandlerFunc
	Sudo    bool
}

type ModuleInfo struct {
	Name        string
	Description string
	Commands    []CommandInfo
}

var RegisteredModules []ModuleInfo

func Register(m ModuleInfo) {
	RegisteredModules = append(RegisteredModules, m)
}
