package modules

import (
	"strings"

	"github.com/amarnathcjd/gogram/telegram"
)

var rtextList = []string{
	"a","b","c","d","e","f","g","h","i","j","k","l","m",
	"n","o","p","q","r","s","t","u","v","w","x","y","z",
	"0","1","2","3","4","5","6","7","8","9",
}

var remojiList = []string{
	"⁭\n      💖\n    💖💖\n  💖💖💖\n💖💖💖💖\n  💖💖💖\n    💖💖\n      💖\n",
	"⁭\n💗💗💗💗\n💗      💗\n💗💗💗💗\n💗💗💗💗\n💗      💗\n💗      💗\n",
	"⁭\n  💛💛💛\n💛      💛\n💛\n  💛💛\n      💛\n💛      💛\n  💛💛💛\n",
	"⁭\n💙💙💙💙\n💙      💙\n💙      💙\n💙      💙\n💙      💙\n💙💙💙💙\n",
	"⁭\n💟💟💟💟\n💟\n💟💟💟\n💟\n💟💟💟💟\n",
	"⁭\n💚💚💚💚\n💚\n💚💚💚\n💚\n💚\n",
	"⁭\n  💜💜💜\n💜      💜\n💜  💜💜\n💜      💜\n  💜💜💜\n",
	"⁭\n💖      💖\n💖      💖\n💖💖💖💖\n💖      💖\n💖      💖\n",
	"⁭\n💗💗💗\n  💗\n  💗\n  💗\n💗💗💗\n",
	"⁭\n  💛💛\n    💛\n    💛\n💛  💛\n  💛💛\n",
	"⁭\n💙  💙\n💙💙\n💙💙\n💙  💙\n💙  💙\n",
	"⁭\n💟\n💟\n💟\n💟\n💟💟💟\n",
	"⁭\n💚  💚\n💚💚💚\n💚💚💚\n💚  💚\n",
	"⁭\n💜  💜\n💜💜\n💜  💜\n💜  💜\n",
	"⁭\n  💖💖\n💖    💖\n💖💖💖💖\n💖    💖\n💖    💖\n",
	"⁭\n💗💗💗\n💗    💗\n💗💗💗\n💗\n💗\n",
	"⁭\n  💛💛\n💛    💛\n💛💛💛💛\n💛    💛\n💛    💛\n",
	"⁭\n💙💙💙\n💙    💙\n💙💙💙\n💙    💙\n💙    💙\n",
	"⁭\n  💟💟\n💟\n  💟💟\n      💟\n  💟💟\n",
	"⁭\n💚💚💚\n  💚\n  💚\n  💚\n",
	"⁭\n💜    💜\n💜    💜\n  💜💜\n    💜\n",
	"⁭\n💖  💖\n  💖\n  💖\n  💖\n",
	"⁭\n💗  💗\n💗💗💗\n💗  💗\n",
	"⁭\n💛  💛\n  💛\n  💛\n  💛\n",
	"⁭\n💙  💙\n  💙\n  💙\n  💙\n",
	"⁭\n💟💟\n    💟\n  💟\n💟💟\n",
	"⁭\n  💗💗\n💗    💗\n💗    💗\n  💗💗\n",
	"⁭\n  💙\n💙💙\n  💙\n  💙\n💙💙💙\n",
	"⁭\n  💟\n💟\n  💟\n      💟\n  💟\n",
	"⁭\n  💛\n💛\n  💛💛\n      💛\n  💛\n",
	"⁭\n💖💖💖\n    💖\n  💖💖\n💖\n💖💖💖\n",
	"⁭\n💚💚💚\n💚\n💚💚\n💚\n💚💚💚\n",
	"⁭\n  💜💜\n💜    💜\n    💜💜\n💜    💜\n  💜💜\n",
	"⁭\n💗💗💗\n💗\n💗\n💗\n💗💗💗\n",
	"⁭\n  💙💙\n💙    💙\n  💙💙\n💙    💙\n  💙💙\n",
	"⁭\n  💟💟\n💟    💟\n  💟💟\n      💟\n  💟💟\n",
}

// .emoji [text]
func emojiHandler(m *telegram.NewMessage) error {
	op, _ := Reply(m, "`✨ Emojifying...`")

	args := GetArgs(m)
	if args == "" && m.IsReply() {
		r, _ := m.GetReplyMessage()
		if r != nil {
			args = r.Text()
		}
	}
	if args == "" {
		op.Edit("`⚠️ Usage: .emoji hello`", nil)
		return nil
	}

	var result strings.Builder
	for _, ch := range strings.ToLower(args) {
		c := string(ch)
		found := false
		for i, t := range rtextList {
			if c == t && i < len(remojiList) {
				result.WriteString(remojiList[i])
				found = true
				break
			}
		}
		if !found {
			result.WriteString(c)
		}
	}

	out := result.String()
	if strings.TrimSpace(out) == "" {
		op.Edit("`❌ Could not convert text.`", nil)
		return nil
	}
	op.Edit(out, nil)
	return nil
}

// .cmoji [emoji] [text]
func cmojiHandler(m *telegram.NewMessage) error {
	op, _ := Reply(m, "`✨ Emojifying with custom emoji...`")

	args := GetArgs(m)
	if args == "" && m.IsReply() {
		r, _ := m.GetReplyMessage()
		if r != nil {
			args = r.Text()
		}
	}
	if args == "" {
		op.Edit("`⚠️ Usage: .cmoji 😎 hello`", nil)
		return nil
	}

	parts := strings.SplitN(args, " ", 2)
	var emoji, text string
	if len(parts) == 2 {
		emoji = parts[0]
		text = parts[1]
	} else {
		emoji = "😎"
		text = args
	}

	templates := []string{
		"⁭\n      {e}\n    {e}{e}\n  {e}{e}{e}\n{e}{e}{e}{e}\n  {e}{e}{e}\n    {e}{e}\n      {e}\n",
		"⁭\n{e}{e}{e}{e}\n{e}      {e}\n{e}{e}{e}{e}\n{e}{e}{e}{e}\n{e}      {e}\n{e}      {e}\n",
		"⁭\n  {e}{e}{e}\n{e}      {e}\n{e}\n  {e}{e}\n      {e}\n{e}      {e}\n  {e}{e}{e}\n",
		"⁭\n{e}{e}{e}{e}\n{e}      {e}\n{e}      {e}\n{e}      {e}\n{e}      {e}\n{e}{e}{e}{e}\n",
		"⁭\n{e}{e}{e}{e}\n{e}\n{e}{e}{e}\n{e}\n{e}{e}{e}{e}\n",
		"⁭\n{e}{e}{e}{e}\n{e}\n{e}{e}{e}\n{e}\n{e}\n",
		"⁭\n  {e}{e}{e}\n{e}      {e}\n{e}  {e}{e}\n{e}      {e}\n  {e}{e}{e}\n",
		"⁭\n{e}      {e}\n{e}      {e}\n{e}{e}{e}{e}\n{e}      {e}\n{e}      {e}\n",
		"⁭\n{e}{e}{e}\n  {e}\n  {e}\n  {e}\n{e}{e}{e}\n",
		"⁭\n  {e}{e}\n    {e}\n    {e}\n{e}  {e}\n  {e}{e}\n",
	}

	var result strings.Builder
	for _, ch := range strings.ToLower(text) {
		c := string(ch)
		found := false
		for i, t := range rtextList {
			if c == t && i < len(templates) {
				tmpl := strings.ReplaceAll(templates[i], "{e}", emoji)
				result.WriteString(tmpl)
				found = true
				break
			}
		}
		if !found {
			result.WriteString(c)
		}
	}

	out := result.String()
	if strings.TrimSpace(out) == "" {
		op.Edit("`❌ Could not convert text.`", nil)
		return nil
	}
	op.Edit(out, nil)
	return nil
}

func init() {
	Register(ModuleInfo{
		Name:        "Emoji",
		Description: "Convert text to emoji art",
		Commands: []CommandInfo{
			{Pattern: "emoji", Handler: emojiHandler, Sudo: false},
			{Pattern: "ej", Handler: emojiHandler, Sudo: false},
			{Pattern: "cmoji", Handler: cmojiHandler, Sudo: false},
			{Pattern: "cj", Handler: cmojiHandler, Sudo: false},
		},
	})
}

var _ = telegram.HTML
