package modules

type CommandInfo struct {
	Pattern string
	Func    any
	Sudo    bool
}

type ModuleInfo struct {
	Name        string
	Description string
	Commands    []CommandInfo
}

var RegisteredModules = []ModuleInfo{}

func RegisterModule(module ModuleInfo) {
	RegisteredModules = append(RegisteredModules, module)
}
