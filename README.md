# 🤖 PbxGo

> Modular Telegram Userbot/Bot — written in **Go 1.24** using [gogram](https://github.com/amarnathcjd/gogram)

---

## ✨ Features

- ⚡ **Go 1.24** — latest language features, `log/slog`, structured logging
- 🤖 **Bot Token** support (run as a normal bot via `BOT_TOKEN`)
- 👤 **Userbot** support (run as your account via `STRING_SESSION`)
- 🔐 **Sudo Management** — `.addsudo`, `.rmsudo`, `.listsudo`
- 🗄️ **MongoDB v2** — optional persistent sudo storage (`mongo-driver/v2`)
- 🔒 **Thread-safe** — `sync.Map` for concurrent sudo access
- 🧩 **Modular** — drop new `.go` files in `modules/` to extend

---

## ⚙️ Setup

### 1. Requirements

- Go **1.24+** → https://go.dev/dl/
- Telegram API credentials → https://my.telegram.org/apps
- (Optional) MongoDB instance

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
MONGO_URL=mongodb://localhost:27017   # optional

# Pick ONE:
BOT_TOKEN=your_bot_token_here         # Bot mode (@BotFather)
STRING_SESSION=                       # Userbot mode
```

**Login priority:** `BOT_TOKEN` → `STRING_SESSION` → interactive prompt

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
│   └── client.go        # Telegram client, filters, handler registration
├── config/
│   └── config.go        # .env loader (uses log/slog, os.LookupEnv)
├── database/
│   ├── db.go            # MongoDB v2 connection
│   └── sudo.go          # Sudo CRUD (sync.Map + MongoDB)
├── modules/
│   ├── module.go        # Module/Command types + Register()
│   ├── utils.go         # Reply() helper
│   ├── start.go         # .alive, .ping
│   └── sudo.go          # .addsudo, .rmsudo, .listsudo
├── main.go              # Entry point
├── go.mod               # Go 1.24
└── sample.env
```

---

## 💬 Commands

| Command | Description | Access |
|--------|-------------|--------|
| `.alive` | Show bot status & version | Owner / Sudo |
| `.ping` | Ping with uptime | Owner / Sudo |
| `.addsudo <id>` | Add sudo user (reply or ID) | Owner only |
| `.rmsudo <id>` | Remove sudo user | Owner only |
| `.listsudo` | List all sudo users | Owner / Sudo |

---

## 🧩 Adding a New Module

Create `modules/mymodule.go`:

```go
package modules

import "github.com/amarnathcjd/gogram/telegram"

func helloHandler(m *telegram.NewMessage) error {
    Reply(m, "👋 Hello from <b>MyModule</b>!")
    return nil
}

func init() {
    Register(ModuleInfo{
        Name:        "MyModule",
        Description: "Says hello.",
        Commands: []CommandInfo{
            {Pattern: "hello", Handler: helloHandler, Sudo: true},
        },
    })
}
```

That's it — no wiring needed. `init()` auto-registers it.

---

## 📄 License

MIT © PbxGo
