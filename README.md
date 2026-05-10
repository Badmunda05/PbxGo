# 🤖 PbxGo

> Modular Telegram Userbot/Bot — **Go 1.24** | [gogram](https://github.com/amarnathcjd/gogram)

---

## ✨ Features

- ⚡ **Go 1.24** — modern stdlib, `log/slog`, structured logging
- 🔐 **OWNER_ID** — all `.` commands locked to owner only
- 🤖 **Bot Token** support via `BOT_TOKEN`
- 👤 **Userbot** support via `STRING_SESSION`
- 🗄️ **MongoDB v2** — optional persistent sudo storage
- 🧩 **Modular** — drop `.go` files in `modules/` to extend

---

## ⚙️ Setup

```bash
git clone https://github.com/youruser/PbxGo.git
cd PbxGo
cp sample.env .env
# Edit .env with your credentials
go mod tidy
go run .
```

---

## 📁 Project Structure

```
PbxGo/
├── client/client.go       # Client init, OWNER_ID filter
├── config/config.go       # .env loader
├── database/
│   ├── db.go              # MongoDB v2 connection
│   └── sudo.go            # Sudo CRUD
├── modules/
│   ├── module.go          # Module/Command types
│   ├── utils.go           # Reply() helper
│   ├── start.go           # .alive .ping /start
│   ├── help.go            # .help
│   ├── restart.go         # .restart .rs .reload
│   ├── purge.go           # .del .purge .purgeme
│   ├── spam.go            # .spam .delayspam .sspam
│   ├── broadcast.go       # .gcast .gucast
│   ├── clone.go           # .clone .revert
│   ├── emoji.go           # .emoji .cmoji
│   ├── hang.go            # .hang
│   ├── tagger.go          # .all .cancel
│   ├── invite.go          # .inviteall
│   └── sudo.go            # .addsudo .rmsudo .sudolist
├── main.go
├── go.mod
└── sample.env
```

---

## 💬 Commands

| Command | Description | Access |
|--------|-------------|--------|
| `/start` | Welcome message (bot mode) | Everyone |
| `.alive` | Bot status & uptime | Owner/Sudo |
| `.ping` | Speed check | Owner/Sudo |
| `.help` | Full command list | Owner/Sudo |
| `.restart` / `.rs` | Restart bot | Owner |
| `.del` | Delete replied message | Owner |
| `.purge` | Purge from replied msg | Owner |
| `.purgeme [n]` | Delete your last n msgs | Owner |
| `.spam [n] [text]` | Spam messages | Owner |
| `.delayspam [s] [n] [text]` | Delayed spam | Owner |
| `.sspam [n]` | Sticker spam (reply) | Owner |
| `.gcast [text]` | Broadcast to groups | Owner |
| `.gucast [text]` | Broadcast to private | Owner |
| `.clone [@user]` | Clone a user's profile | Owner |
| `.revert` | Revert your profile | Owner |
| `.emoji [text]` | Text to emoji art | Owner |
| `.cmoji [e] [text]` | Custom emoji art | Owner |
| `.hang [n]` | Send hang messages | Owner |
| `.all [text]` | Mention all members | Owner |
| `.cancel` | Stop tagger | Owner |
| `.inviteall [@grp]` | Invite members | Owner |
| `.addsudo [id]` | Add sudo user | Owner |
| `.rmsudo [id]` | Remove sudo user | Owner |
| `.sudolist` | List sudo users | Owner/Sudo |

---

## 📄 License

MIT © PbxGo | Support: @PBXCHATS | Updates: @PBX_UPDATE
