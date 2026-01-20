# Telegram to Notion Bot

[![Go](https://img.shields.io/badge/go-1.23-blue.svg)](https://golang.org)
[![License](https://img.shields.io/github/license/nerdneilsfield/telenotion-bot)](LICENSE)
[![Release](https://img.shields.io/github/v/release/nerdneilsfield/telenotion-bot)](https://github.com/nerdneilsfield/telenotion-bot/releases)
[![GoReleaser](https://github.com/nerdneilsfield/telenotion-bot/actions/workflows/goreleaser.yml/badge.svg)](https://github.com/nerdneilsfield/telenotion-bot/actions/workflows/goreleaser.yml)

A Telegram bot that captures messages into a Notion database with MarkdownV2-aware formatting.

## Features

- Capture messages from configured Telegram chats between `/start` and `/end`
- Support `/clean` to clear current buffer
- Support `/discard` to abandon current session
- Convert Telegram MarkdownV2 entities to Notion rich text
- Upload images to GitHub and embed them as Notion image blocks
- Create Notion database pages with timestamped titles

## Configuration

Create a TOML configuration file (see `config.example.toml`):

```toml
[telegram]
token = "your-telegram-bot-token"
allowed_chat_ids = [123456789]

[notion]
token = "your-notion-integration-token"
database_id = "your-database-id"
title_property = "Name"

[github]
token = "your-github-pat"
repo = "owner/repo"
branch = "main"
path_prefix = "images/"

[title]
timezone = "Asia/Shanghai"
format = "2006-01-02 15:04"
```

## Usage

Run the bot:

```bash
telenotion-bot bot --config config.toml
```

### Commands

- `/start` - Start a new capture session
- `/clean` - Clear current buffer
- `/discard` - Abandon current session
- `/end` - Create Notion page and end session

### Supported MarkdownV2 Syntax

- `*bold text*` - Bold
- `_italic text_` - Italic
- `` `code` `` - Inline code
- ```code``` - Code block
- `[text](url)` - Link

### Image Handling

Images are automatically uploaded to GitHub when `/end` is called:
1. Image is downloaded from Telegram
2. Uploaded to configured GitHub repository
3. Notion embeds as GitHub raw URL

## Notion Database Requirements

Your Notion database must have:
- A `title` property (type: Title)
- No additional required properties (content is stored as child blocks)

## Development

Run tests:

```bash
go test ./...
```

Build:

```bash
go build ./...
```

Format code:

```bash
go fmt ./...
goimports -w .
```
