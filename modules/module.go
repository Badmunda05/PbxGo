package modules

import "github.com/amarnathcjd/gogram/telegram"

// IsBotMode — set by client after login
var IsBotMode = true

type HandlerFunc = func(*telegram.NewMessage) error

type CommandDef struct {
	Pattern   string
	Handler   HandlerFunc
	OwnerOnly bool
}

type Module struct {
	Name     string
	Commands []CommandDef
}

// UserModules — handlers for userbot clients
var UserModules []Module

// BotModules — handlers for bot client (commands)
var BotModules []Module

// BotRawHandlers — raw event handlers for bot (e.g. all messages for conversation)
type RawHandler struct {
	Event   string
	Handler HandlerFunc
}

var BotRawHandlers []RawHandler

func RegisterUser(m Module)  { UserModules = append(UserModules, m) }
func RegisterBot(m Module)   { BotModules = append(BotModules, m) }

// RegisterBotRaw — register a raw event handler on the bot (e.g. "message" for all texts)
func RegisterBotRaw(event string, handler HandlerFunc) {
	BotRawHandlers = append(BotRawHandlers, RawHandler{Event: event, Handler: handler})
}
