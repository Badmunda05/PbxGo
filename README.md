# 🤖 PbxGo

> Modular Telegram Userbot/Bot — **Go 1.24** | [gogram](https://github.com/amarnathcjd/gogram)

---

## ✨ Features

- ⚡ **Go 1.24** — modern stdlib, `log/slog`, structured logging
- 🔐 **Owner Lock** — all `.` commands locked to `OWNER_ID` only (both bot and userbot mode)
- 👥 **Sudo System** — owner can grant trusted users sudo access
- 🤖 **Bot Mode** — via `BOT_TOKEN`, `/start` is public for everyone
- 👤 **Userbot Mode** — via `STRING_SESSION`
- 🗄️ **MongoDB v2** — optional persistent sudo storage (works without MongoDB too)
- 🧩 **Modular** — add a new `.go` file in `modules/` to extend

---

## ⚙️ Setup

```bash
git clone https://github.com/youruser/PbxGo.git
cd PbxGo
cp sample.env .env
# Fill in your credentials in .env
go mod tidy
go run .
```

---

## 🔧 .env Configuration

```env
APP_ID=12345678
APP_HASH=your_app_hash_here
OWNER_ID=your_telegram_user_id

# Choose one mode:
STRING_SESSION=your_string_session   # Userbot mode
BOT_TOKEN=your_bot_token             # Bot mode

# Optional — for persistent sudo storage
MONGO_URL=mongodb://localhost:27017
```

> ⚠️ `OWNER_ID` is required — bot will not start without it.

---

## 📁 Project Structure

```
PbxGo/
├── client/client.go       # Client init, owner/sudo filters, handler registration
├── config/
│   ├── config.go          # .env loader
│   └── data.go            # Raid text data
├── database/
│   ├── db.go              # MongoDB v2 connection
│   └── sudo.go            # Sudo CRUD (add/remove/list/check)
├── modules/
│   ├── module.go          # Module and Command type definitions
│   ├── utils.go           # Reply() helper, GetArgs()
│   ├── start.go           # .alive  .ping  /start
│   ├── help.go            # .help
│   ├── restart.go         # .restart  .rs  .reload
│   ├── purge.go           # .del  .purge  .purgeme
│   ├── spam.go            # .spam  .delayspam  .sspam
│   ├── broadcast.go       # .gcast  .gucast
│   ├── clone.go           # .clone  .revert
│   ├── emoji.go           # .emoji  .cmoji
│   ├── hang.go            # .hang
│   ├── tagger.go          # .all  .cancel
│   ├── raid.go            # .raid  .hraid  .eraid  .punraid
│   └── sudo.go            # .addsudo  .rmsudo  .sudolist
├── main.go
├── go.mod
└── sample.env
```

---

## 💬 Commands

### 🔑 Access Levels
| Level | Who |
|-------|-----|
| **Owner** | Only the user whose ID is set as `OWNER_ID` in `.env` |
| **Owner/Sudo** | Owner + any user added via `.addsudo` |
| **Everyone** | All users (bot mode `/start` only) |

---

### ⚙️ Core
| Command | Description | Access |
|---------|-------------|--------|
| `/start` | Welcome message | Everyone (bot mode only) |
| `.alive` | Bot status and uptime | Owner/Sudo |
| `.ping` | Speed and uptime check | Owner/Sudo |
| `.help` | Full command list | Owner/Sudo |
| `.restart` / `.rs` / `.reload` | Restart the bot | Owner |

---

### 🗑️ Purge
| Command | Description | Access |
|---------|-------------|--------|
| `.del` | Delete the replied message | Owner |
| `.purge` | Purge all messages from replied message onward | Owner |
| `.purgeme [n]` | Delete your own last n messages | Owner |

---

### 💬 Spam
| Command | Description | Access |
|---------|-------------|--------|
| `.spam [n] [text]` | Send a message n times | Owner |
| `.delayspam [s] [n] [text]` | Send a message n times with s seconds delay | Owner |
| `.sspam [n]` | Spam a sticker n times (reply to sticker) | Owner |

---

### 📢 Broadcast
| Command | Description | Access |
|---------|-------------|--------|
| `.gcast [text/reply]` | Broadcast to all groups | Owner |
| `.gucast [text/reply]` | Broadcast to all private chats | Owner |

---

### 👤 Clone
| Command | Description | Access |
|---------|-------------|--------|
| `.clone [@user/reply]` | Clone a user's profile (name, bio, photo) | Owner |
| `.revert` | Revert back to your original profile | Owner |

---

### ✏️ Emoji Art
| Command | Description | Access |
|---------|-------------|--------|
| `.emoji [text]` | Convert text to emoji art | Owner |
| `.cmoji [emoji] [text]` | Emoji art with a custom emoji | Owner |

---

### 🏓 Hang
| Command | Description | Access |
|---------|-------------|--------|
| `.hang [n]` | Send hang messages | Owner |

---

### 👥 Tagger
| Command | Description | Access |
|---------|-------------|--------|
| `.all [text]` | Mention all members in the group | Owner |
| `.cancel` | Stop the tagger | Owner |

---

### ⚔️ Raid
| Command | Description | Access |
|---------|-------------|--------|
| `.raid` | Default raid messages | Owner |
| `.hraid` | Hindi raid messages | Owner |
| `.eraid` | English raid messages | Owner |
| `.punraid` | Punjabi raid messages | Owner |

---

### 🔐 Sudo Management
| Command | Description | Access |
|---------|-------------|--------|
| `.addsudo [id/reply]` / `.asd` | Add a sudo user | Owner |
| `.rmsudo [id/reply]` / `.delsudo` | Remove a sudo user | Owner |
| `.sudolist` / `.sdl` | View the sudo user list | Owner/Sudo |

---

## 🔐 Security

- All `.` commands are **owner-locked** — verified by Telegram `SenderID`, cannot be spoofed
- Works the same in **both bot mode and userbot mode**
- No other user, regardless of group or chat, can trigger any command
- Sudo users can only run commands marked `Sudo: true` (e.g. `.alive`, `.ping`, `.sudolist`)
- Sensitive commands (`.addsudo`, `.rmsudo`, `.restart`, `.gcast`, etc.) are **owner only**

---

## 📄 License

MIT © PbxGo | Support: @PBXCHATS | Updates: @PBX_UPDATE
