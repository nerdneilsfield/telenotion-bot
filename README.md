# 🚀 Telegram to Notion Bot — Your Personal Capture Assistant ✨

[![Go 1.23](https://img.shields.io/badge/Go-1.23-blue?logo=go)](https://golang.org)
[![MIT License](https://img.shields.io/badge/License-MIT-green)](LICENSE)
[![CI Status](https://img.shields.io/github/actions/workflow/status/nerdneilsfield/telenotion-bot/build.yml?label=CI&branch=master)](https://github.com/nerdneilsfield/telenotion-bot/actions)
[![Latest Release](https://img.shields.io/github/v/release/nerdneilsfield/telenotion-bot?color=orange)](https://github.com/nerdneilsfield/telenotion-bot/releases)
[![Downloads](https://img.shields.io/github/downloads/nerdneilsfield/telenotion-bot/total)](https://github.com/nerdneilsfield/telenotion-bot/releases)
[![Docker Image](https://img.shields.io/badge/Docker-ghcr.io%2Fnerdneilsfield%2Ftelenotion-bot-blue?logo=docker)](https://github.com/nerdneilsfield/telenotion-bot/pkgs/container/telenotion-bot)
[![Go Report Card](https://goreportcard.com/badge/github.com/nerdneilsfield/telenotion-bot)](https://goreportcard.com/report/github.com/nerdneilsfield/telenotion-bot)

[中文说明](README_ZH.md) | [English](README.md)

---

## 👋 Hey everyone! Gotta share something amazing with you!

Listen... Every time you see something useful on Telegram, do you find yourself:
1. Copy-pasting to Notion
2. Manually fixing the formatting
3. Uploading images
4. ... Just to feel exhausted afterwards 😭

**WELL! That's all in the PAST now!** 🎉

---

## ✨ What can this thing do?

**Telegram to Notion Bot** is your personal capture assistant:

| Superpower | Description |
|------------|-------------|
| 📝 **Message Capture** | Every message during a session, all remembered for you! |
| 🎨 **Format Preservation** | Bold, italic, code blocks, links — copied exactly! |
| 🖼️ **Image Hosting** | Auto-upload to GitHub, displayed directly in Notion! |
| ⏰ **Timestamp Titles** | Auto-generate page titles with current time |
| 💾 **Block Storage** | All content saved as Notion child blocks |

**TL;DR: Copy-paste? Nah! Just send and go, Notion organizes itself!** 😎

---

## 🎯 Core Features at a Glance

| Feature | Command | Effect |
|---------|---------|--------|
| `/start` ✨ | Start new session | Enable capture mode |
| `/clean` 🧹 | Clear buffer | Clear current content, session continues |
| `/discard` 🔄 | Discard session | Start fresh |
| `/end` 💾 | Save to Notion | Generate page, end session |
| `/help` 📖 | Show help | Display all commands |

**Markdown Support**: `*bold*` → ✅ | `_italic_` → ✅ | `` `code` `` → ✅ | ```code block``` → ✅ | `[link](url)` → ✅

**Image Handling**: Just send! Automatically download Telegram images → upload to GitHub → embed in Notion 🖼️

---

## 🚀 Get Started in 5 Minutes

### Option 1: Download & Run (Recommended for Beginners)

```bash
# 1️⃣ Download the latest version
# Visit https://github.com/nerdneilsfield/telenotion-bot/releases
# Download and extract the archive for your system

# 2️⃣ Create config file
cp config.example.toml config.toml

# 3️⃣ Edit config (see detailed config below)
vim config.toml

# 4️⃣ Launch!
./telenotion-bot bot -c config.toml
```

### Option 2: Build from Source (For Developers)

```bash
# 1️⃣ Clone the repo
git clone https://github.com/nerdneilsfield/telenotion-bot.git
cd telenotion-bot

# 2️⃣ Build
go build -o telenotion-bot .

# 3️⃣ Launch (same as above)
./telenotion-bot bot -c config.toml
```

### Option 3: Docker One-Click Deployment 🐳

```bash
# The easiest deployment method!
docker-compose up -d
```

---

## ⚙️ Config Details (3 Services, 5 Minutes)

> 💡 **Pro Tip**: Both environment variables AND TOML config file are supported! Environment variables recommended for Docker deployment~

### Telegram Config

```toml
[telegram]
token = "YOUR_TELEGRAM_BOT_TOKEN"  # Create via @BotFather
allowed_chat_ids = [123456789, 987654321]  # Allowed chat/group IDs
```

### Notion Config

```toml
[notion]
token = "YOUR_NOTION_INTEGRATION_TOKEN"  # https://www.notion.so/my-integrations
database_id = "YOUR_DATABASE_ID"  # Long string in database URL
title_property = "Name"  # Title field name in your database
```

### GitHub Config (Image Hosting)

```toml
[github]
token = "YOUR_GITHUB_PAT"  # https://github.com/settings/tokens
repo = "username/repo"     # e.g., "nerdneilsfield/my-images"
branch = "main"            # Branch name
path_prefix = "images/"    # Image storage directory
```

### Title Format Config

```toml
[title]
timezone = "Asia/Shanghai"  # Timezone
format = "2006-01-02 15:04" # Page title format
```

### Log Config

```toml
[log]
level = "info"   # debug | info | warn | error
file = ""        # Log file path, empty = stdout only
```

### 🔐 Environment Variable Support (Essential for Docker!)

```bash
# Don't want config.toml? No problem!
export TELEGRAM_TOKEN="xxx"
export TELEGRAM_ALLOWED_CHAT_IDS="123,456,789"
export NOTION_TOKEN="xxx"
export NOTION_DATABASE_ID="xxx"
export NOTION_TITLE_PROPERTY="Name"
export GITHUB_TOKEN="xxx"
export GITHUB_REPO="owner/repo"
export GITHUB_BRANCH="main"
export GITHUB_PATH_PREFIX="images/"
export TITLE_TIMEZONE="Asia/Shanghai"
export TITLE_FORMAT="2006-01-02 15:04"
export LOG_LEVEL="info"
export LOG_FILE=""

# Then run (no -c flag needed)
./telenotion-bot bot
```

---

## 📱 How to Use

### Step 1: Start Capturing ✨

```
/start
```

Bot replies:
> *"Session started. Send messages or images, then /end to save."* ✨

### Step 2: Start Sending Messages 📝

Send whatever you want!

**Supported formats:**
- `*bold text*` → **bold text**
- `_italic text_` → *italic text*
- `` `console.log('hi')` `` → `console.log('hi')`
- ```javascript\nconsole.log('code block')\n``` → Code block
- `[Visit Google](https://google.com)` → [Visit Google](https://google.com)

**Images**: Just send! Bot handles everything automatically~ 🖼️

### Step 3: Save to Notion 💾

```
/end
```

Done! Check Notion for your new page! 🎉

---

## 🛠️ Developer Friendly

### Build & Test

```bash
# Build
go build ./...

# Tests (90%+ coverage! 🎉)
go test ./... -cover

# Detailed test output
go test ./... -v
```

### Code Style

```bash
# Format
go fmt ./...

# Organize imports
goimports -w .
```

### Tech Stack

| Tech | Purpose |
|------|---------|
| Go 1.23+ | Development language |
| Telegram Bot API | Message receiving |
| Notion API | Page creation |
| GitHub Contents API | Image hosting |
| Zap | Structured logging |
| TOML | Config format |

---

## 📋 Notion Database Requirements

Your database only needs:

| Requirement | Description |
|-------------|-------------|
| ✅ Title property | Field of type Title |
| ✅ No other required fields | Content stored as child blocks |

---

## 🐛 Having Issues?

1. **Check logs**: Use `-v` flag for detailed output
2. **Common problems**:
   - `telegram.token is required` → Check your Token
   - `notion.database_id is required` → Check database ID
   - Image upload failed → Check GitHub Token permissions
3. **Still stuck?** → [Open an Issue](https://github.com/nerdneilsfield/telenotion-bot/issues)

---

## 🤝 Contributing

Found a bug? Have an idea?

**We welcome all contributions!** 🌟

- 🐛 Report bugs
- 💡 Suggest features
- 🔧 Submit code
- 📖 Improve docs

---

## 📝 License

MIT License — **Free! Open source! Use it however you want!** 🎊

---

## 💬 Final Words

Hope this tool helps you save time on all that repetitive work!

**Got ideas? Let's chat on GitHub!**

[🐙 GitHub](https://github.com/nerdneilsfield/telenotion-bot) | [🐛 Report Issues](https://github.com/nerdneilsfield/telenotion-bot/issues)

---

**Made with ❤️ and a lot of ☕**

*Your Telegram → Notion bridge, serving you~* 🚀
