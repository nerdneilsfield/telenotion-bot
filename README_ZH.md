# 🚀 Telegram + Discord to Notion Bot — 你的私人捕获小助手 ✨

[![Go 1.23](https://img.shields.io/badge/Go-1.23-blue?logo=go)](https://golang.org)
[![MIT License](https://img.shields.io/badge/License-MIT-green)](LICENSE)
[![CI Status](https://img.shields.io/github/actions/workflow/status/nerdneilsfield/telenotion-bot/build.yml?label=CI&branch=master)](https://github.com/nerdneilsfield/telenotion-bot/actions)
[![Latest Release](https://img.shields.io/github/v/release/nerdneilsfield/telenotion-bot?color=orange)](https://github.com/nerdneilsfield/telenotion-bot/releases)
[![Downloads](https://img.shields.io/github/downloads/nerdneilsfield/telenotion-bot/total)](https://github.com/nerdneilsfield/telenotion-bot/releases)
[![Docker Image](https://img.shields.io/badge/Docker-ghcr.io%2Fnerdneilsfield%2Ftelenotion-bot-blue?logo=docker)](https://github.com/nerdneilsfield/telenotion-bot/pkgs/container/telenotion-bot)
[![Go Report Card](https://goreportcard.com/badge/github.com/nerdneilsfield/telenotion-bot)](https://goreportcard.com/report/github.com/nerdneilsfield/telenotion-bot)

---

## 👋 姐妹们！有个超好用的东西必须安利给你们！

讲真！是不是每次在 Telegram 或 Discord 看到什么有用的信息，都得：
1. 复制粘贴到 Notion
2. 手动调整格式
3. 上传图片
4. ... 一套流程下来累觉不爱 😭

**现在！这一切都！不！存！在！了！** 🎉

---

## ✨ 这玩意儿能干嘛？

**Telegram + Discord to Notion Bot** 就是你的私人捕获小助手：

| 超能力 | 说明 |
|--------|------|
| 📝 **消息捕获** | 会话期间的所有消息，全部帮你记住！ |
| 🎨 **格式保留** | 粗体、斜体、代码块、链接，原样搬运！ |
| 🖼️ **图片托管** | 自动上传到 GitHub，Notion 里直接显示！ |
| ⏰ **时间戳标题** | 自动用当前时间生成页面标题 |
| 💾 **块级存储** | 所有内容以 Notion 子块形式保存 |

**一句话：复制粘贴？不存在！发完就走，Notion 自动帮你整理好！** 😎

---

## 🎯 核心功能一览

| 功能 | 命令 | 效果 |
|------|------|------|
| `/start` ✨ | 开始新会话 | 开启捕获模式 |
| `/clean` 🧹 | 清空缓存 | 清除当前内容，会话继续 |
| `/discard` 🔄 | 放弃会话 | 重新开始 |
| `/end` 💾 | 保存到 Notion | 生成页面，结束会话 |
| `/help` 📖 | 查看帮助 | 显示所有命令 |

**Markdown 支持**：`*粗体*` → ✅ | `_斜体_` → ✅ | `` `代码` `` → ✅ | ```代码块``` → ✅ | `[链接](url)` → ✅

**图片处理**：直接发！自动下载 Telegram/Discord 图片 → 上传到 GitHub → 嵌入 Notion 🖼️

---

## 🤖 机器人申请与使用（Telegram + Discord）

### Telegram
1. 用 `@BotFather` 创建机器人并获取 Token。
2. 用 `@wczj_userinfo_bot` 获取你的 chat ID。
3. 将 chat ID 写入 `allowed_chat_ids`。

### Discord
1. 在 Discord Developer Portal 创建应用并添加 Bot。
2. 复制 Bot Token，并勾选 **Message Content Intent**。
3. 邀请机器人时包含 `bot` + `applications.commands` 权限。
4. 将你的用户 ID 写入 `allowed_user_ids`。

### 使用方法（两端通用）
- 私聊机器人输入 `/start`，发送消息或图片，最后 `/end`。
- `/end` 后可直接发消息自动开启新会话。

---

## 通过 Homebrew / Scoop 安装

```bash
# macOS / Linux
brew install --cask nerdneilsfield/apps/telenotion-bot
```

```powershell
# Windows
scoop bucket add nerdneilsfield https://github.com/nerdneilsfield/homebrew-apps
scoop install nerdneilsfield/telenotion-bot
```

## 🚀 5 分钟快速上手

### 方式一：下载运行（推荐新手）

```bash
# 1️⃣ 下载最新版本
# 访问 https://github.com/nerdneilsfield/telenotion-bot/releases
# 下载对应系统的压缩包并解压

# 2️⃣ 创建配置文件
cp config.example.toml config.toml

# 3️⃣ 编辑配置（见下方配置详解）
vim config.toml

# 4️⃣ 启动！
./telenotion-bot bot -c config.toml
```

### 方式二：源码编译（适合开发者）

```bash
# 1️⃣ 克隆项目
git clone https://github.com/nerdneilsfield/telenotion-bot.git
cd telenotion-bot

# 2️⃣ 构建
go build -o telenotion-bot .

# 3️⃣ 启动（同上）
./telenotion-bot bot -c config.toml
```

### 方式三：Docker 一键部署 🐳

```bash
# 最简单的部署方式！
docker-compose up -d
```

---

## ⚙️ 配置详解（3 个服务，5 分钟搞定）

> 💡 **小贴士**：环境变量和 TOML 配置文件都支持！Docker 部署首选环境变量~

### Telegram 配置

```toml
[telegram]
token = "你的Telegram Bot Token"  # @BotFather 创建
allowed_chat_ids = [123456789, 987654321]  # 允许使用机器人的群组/用户ID
```

可以用 `@wczj_userinfo_bot` 获取 chat ID。

### Discord 配置

```toml
[discord]
token = "你的Discord Bot Token"  # Discord Developer Portal
allowed_user_ids = ["123456789012345678"]  # 允许使用机器人的用户ID
```

Discord 设置说明：
- 在 Discord Developer Portal 勾选 Message Content Intent。
- 邀请机器人时包含 `applications.commands` 权限。

### Notion 配置

```toml
[notion]
token = "你的Notion Integration Token"  # https://www.notion.so/my-integrations
database_id = "你的数据库ID"  # 数据库 URL 中的一大串字符
title_property = "Name"  # 数据库的标题字段名
origin_property = "Origin"  # Select 字段，选项 Discord/Telegram
```

### GitHub 配置（图片托管）

```toml
[github]
token = "你的GitHub PAT"  # https://github.com/settings/tokens
repo = "用户名/仓库名"    # 比如 "nerdneilsfield/my-images"
branch = "main"          # 默认分支
telegram_branch = "telegram" # Telegram 专用分支（可选）
discord_branch = "discord"   # Discord 专用分支（可选）
path_prefix = "images/"  # 图片存放目录
```

设置 `telegram_branch`/`discord_branch` 后会覆盖 `branch`。

### 媒体配置

```toml
[media]
max_image_size_mb = 20
allowed_image_types = ["image/jpeg", "image/png", "image/gif", "image/webp"]
```

### 标题格式配置

```toml
[title]
timezone = "Asia/Shanghai"  # 时区
format = "2006-01-02 15:04" # 页面标题格式
```

### 日志配置

```toml
[log]
level = "info"   # debug | info | warn | error
file = ""        # 日志文件路径，留空则只输出到 stdout
```

### 🔐 环境变量支持（Docker 必备！）

```bash
# 不需要 config.toml？没问题！
export TELEGRAM_TOKEN="xxx"
export TELEGRAM_ALLOWED_CHAT_IDS="123,456,789"
export DISCORD_TOKEN="xxx"
export DISCORD_ALLOWED_USER_IDS="123456789012345678"
export NOTION_TOKEN="xxx"
export NOTION_DATABASE_ID="xxx"
export NOTION_TITLE_PROPERTY="Name"
export NOTION_ORIGIN_PROPERTY="Origin"
export GITHUB_TOKEN="xxx"
export GITHUB_REPO="owner/repo"
export GITHUB_BRANCH="main"
export GITHUB_TELEGRAM_BRANCH="telegram"
export GITHUB_DISCORD_BRANCH="discord"
export GITHUB_PATH_PREFIX="images/"
export MEDIA_MAX_IMAGE_SIZE_MB="20"
export MEDIA_ALLOWED_IMAGE_TYPES="image/jpeg,image/png,image/gif,image/webp"
export TITLE_TIMEZONE="Asia/Shanghai"
export TITLE_FORMAT="2006-01-02 15:04"
export LOG_LEVEL="info"
export LOG_FILE=""

# 然后运行（不需要 -c 参数）
./telenotion-bot bot
```

---

## 📱 使用教程

### Step 1：开始捕获 ✨

Telegram 或 Discord 私聊：
```
/start
```

机器人回复：
> *"Session started. Send messages or images, then /end to save."* ✨

在 `/end` 之后也可以不输入 `/start`，直接发消息会自动开启新会话。

### Step 2：开始发消息 📝

想发啥发啥！

**支持的格式：**
- `*这是粗体*` → **这是粗体**
- `_这是斜体_` → *这是斜体*
- `` `console.log('hi')` `` → `console.log('hi')`
- ```javascript\nconsole.log('code block')\n``` → 代码块
- `[点击访问 Google](https://google.com)` → [点击访问 Google](https://google.com)

**图片**：直接发！机器人自动帮你处理~ 🖼️

### Step 3：保存到 Notion 💾

```
/end
```

搞定了！去 Notion 看你的新页面吧！🎉

---

## 🛠️ 开发者友好

### 构建 & 测试

```bash
# 构建
go build ./...

# 测试（覆盖率超 90%！🎉）
go test ./... -cover

# 单元测试详情
go test ./... -v
```

### 代码规范

```bash
# 格式化
go fmt ./...

# 整理 import
goimports -w .
```

### 技术栈

| 技术 | 用途 |
|------|------|
| Go 1.23+ | 开发语言 |
| Telegram Bot API | 消息接收 |
| Discord API | Slash 命令 + 私聊捕获 |
| Notion API | 页面创建 |
| GitHub Contents API | 图片托管 |
| Zap | 结构化日志 |
| TOML | 配置格式 |

---

## 📋 Notion 数据库要求

你的数据库只需要：

| 要求 | 说明 |
|------|------|
| ✅ 标题属性 | 类型为 Title 的字段 |
| ✅ 无其他必填项 | 内容以子块形式存储 |

---

## 🐛 遇到问题？

1. **检查日志**：用 `-v` 参数看详细输出
2. **常见问题**：
   - `telegram.token is required` → 检查 Token
   - `notion.database_id is required` → 检查数据库 ID
   - 图片上传失败 → 检查 GitHub Token 权限
3. **还没解决？** → [提 Issue](https://github.com/nerdneilsfield/telenotion-bot/issues)

---

## 🤝 贡献指南

发现 Bug？有新想法？

**我们欢迎各种贡献！** 🌟

- 🐛 报告 Bug
- 💡 提出建议
- 🔧 提交代码
- 📖 完善文档

---

## 📝 License

MIT License — **免费！开源！随便用！** 🎊

---

## 💬 最后的最后

希望这个工具能帮你省下那些重复劳动的时间！

**有想法？来 GitHub 找我聊！**

[🐙 GitHub](https://github.com/nerdneilsfield/telenotion-bot) | [🐛 Report Issues](https://github.com/nerdneilsfield/telenotion-bot/issues)

---

**Made with ❤️ and a lot of ☕**

*你的 Telegram/Discord → Notion 桥梁，正在为你服务~* 🚀
