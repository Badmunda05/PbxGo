# 🤖 PbxGo

> Modular Telegram Userbot/Bot — **Go 1.24** | [gogram](https://github.com/amarnathcjd/gogram)

---

## ✨ Features

- ⚡ **Go 1.24** — `log/slog`, structured logging, modern stdlib
- 🔐 **OWNER_ID** — all commands locked to owner only (no spam risk)
- 🤖 **Bot Token** support via `BOT_TOKEN`
- 👤 **Userbot** support via `STRING_SESSION`
- 👥 **Sudo Management** — `.addsudo`, `.rmsudo`, `.listsudo`
- 🗄️ **MongoDB v2** — optional persistent sudo storage
- 🔒 **Thread-safe** — `sync.Map` for concurrent access
- 🧩 **Modular** — drop `.go` files in `modules/` to extend

---

## ⚙️ Setup

### 1. Requirements

- Go **1.24+** → https://go.dev/dl/
- Telegram API credentials → https://my.telegram.org/apps

### 2. Clone & configure

```bash
git clone https://github.com/youruser/PbxGo.git
cd PbxGo
cp sample.env .env
```

Edit `.env`:

```env
APP_ID=123456
APP_HASH=your_app_hash_here
OWNER_ID=123456789          # Your Telegram user ID (@userinfobot)
MONGO_URL=                  # Optional MongoDB URL
BOT_TOKEN=your_bot_token    # Bot mode (@BotFather)
STRING_SESSION=             # Userbot mode
```

> **Login priority:** `BOT_TOKEN` → `STRING_SESSION` → interactive prompt
> **Get your ID:** message @userinfobot on Telegram

### 3. Run

```bash
go mod tidy
go run .
```

---

## 📁 Project Structure

```
PbxGo/
├── client/
│   └── client.go        # Client init, OWNER_ID filter, handler registration
├── config/
│   └── config.go        # .env loader — AppID, AppHash, OwnerID, BotToken, etc.
├── database/
│   ├── db.go            # MongoDB v2 connection
│   └── sudo.go          # Sudo CRUD — sync.Map + MongoDB
├── modules/
│   ├── module.go        # Module/Command types + Register()
│   ├── utils.go         # Reply() helper
│   ├── start.go         # .alive, .ping, /start (bot mode)
│   └── sudo.go          # .addsudo, .rmsudo, .listsudo
├── main.go
├── go.mod               # Go 1.24
└── sample.env
```

---

## 💬 Commands

| Command | Description | Access |
|--------|-------------|--------|
| `/start` | Welcome message (bot mode) | Everyone |
| `.alive` | Bot status, uptime, version | Owner / Sudo |
| `.ping` | Ping with speed & uptime | Owner / Sudo |
| `.addsudo <id>` | Add sudo user | Owner only |
| `.rmsudo <id>` | Remove sudo user | Owner only |
| `.listsudo` | List all sudo users | Owner / Sudo |

> ⚠️ All `.` commands are **OWNER_ID locked** — random users cannot trigger them.

---

## 🧩 Adding a New Module

```go
// modules/hello.go
package modules

import "github.com/amarnathcjd/gogram/telegram"

func helloHandler(m *telegram.NewMessage) error {
    Reply(m, "👋 Hello from <b>MyModule</b>!")
    return nil
}

func init() {
    Register(ModuleInfo{
        Name:        "Hello",
        Description: "Says hello.",
        Commands: []CommandInfo{
            {Pattern: "hello", Handler: helloHandler, Sudo: false},
        },
    })
}
```

---

## 📄 License

MIT © PbxGo
