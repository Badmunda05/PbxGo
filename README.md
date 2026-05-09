# 🤖 PbxGo

> A powerful, modular Telegram Userbot/Bot built with [gogram](https://github.com/amarnathcjd/gogram) — written in Go.

---

## ✨ Features

- **Bot Token Support** — Run as a regular Telegram bot using `BOT_TOKEN`
- **Userbot Support** — Run as a userbot using `STRING_SESSION`
- **Sudo Management** — Add, remove, and list sudo users with `.addsudo`, `.rmsudo`, `.listsudo`
- **MongoDB Integration** — Optional persistent storage for sudo users
- **Modular Design** — Easily add new modules in the `modules/` folder
- **Environment Config** — All credentials loaded from `.env` file

---

## ⚙️ Setup

### 1. Clone the repository

```bash
git clone https://github.com/youruser/PbxGo.git
cd PbxGo
```

### 2. Configure `.env`

Copy `sample.env` to `.env` and fill in your credentials:

```bash
cp sample.env .env
```

Edit `.env`:

```env
APP_ID=123456                         # From https://my.telegram.org
APP_HASH=your_app_hash_here           # From https://my.telegram.org
MONGO_URL=mongodb://localhost:27017   # Optional: MongoDB connection URL
STRING_SESSION=                       # Your Telegram string session (for userbot)
BOT_TOKEN=your_bot_token_here         # Your bot token from @BotFather
```

> **Priority:** `BOT_TOKEN` > `STRING_SESSION` > Interactive login prompt

### 3. Install dependencies

```bash
go mod tidy
```

### 4. Run PbxGo

```bash
go run .
```

---

## 📁 Project Structure

```
PbxGo/
├── client/
│   └── init.go          # Telegram client initialization & handler registration
├── config/
│   └── config.go        # Environment config loader
├── database/
│   ├── init.go          # MongoDB connection
│   └── sudo.go          # Sudo user management
├── modules/
│   ├── commands.go      # Module registration system
│   ├── start.go         # Alive & Ping commands
│   ├── sudo.go          # Sudo management commands
│   └── utils.go         # Helper utilities
├── main.go              # Entry point
├── go.mod
├── sample.env
└── README.md
```

---

## 💬 Commands

| Command | Description | Access |
|--------|-------------|--------|
| `.alive` | Check if PbxGo is running | Owner / Sudo |
| `.ping` | Ping with uptime | Owner / Sudo |
| `.addsudo <id>` | Add a sudo user | Owner only |
| `.rmsudo <id>` | Remove a sudo user | Owner only |
| `.listsudo` | List all sudo users | Owner / Sudo |

---

## 🧩 Adding New Modules

Create a new `.go` file in `modules/` and register your commands:

```go
package modules

import "github.com/amarnathcjd/gogram/telegram"

func MyHandler(m *telegram.NewMessage) error {
    EditOrReply(m, "Hello from my module!")
    return nil
}

func init() {
    RegisterModule(ModuleInfo{
        Name:        "My Module",
        Description: "Does something cool.",
        Commands: []CommandInfo{
            {
                Pattern: "hello",
                Func:    MyHandler,
                Sudo:    true,
            },
        },
    })
}
```

---

## 📄 License

MIT © PbxGo
