# Telegram to Notion Bot

[![Go](https://img.shields.io/badge/go-1.23-blue.svg)](https://golang.org)
[![License](https://img.shields.io/github/license/nerdneilsfield/telenotion-bot)](LICENSE)
[![Release](https://img.shields.io/github/v/release/nerdneilsfield/telenotion-bot)](https://github.com/nerdneilsfield/telenotion-bot/releases)
[![GoReleaser](https://github.com/nerdneilsfield/telenotion-bot/actions/workflows/goreleaser.yml/badge.svg)](https://github.com/nerdneilsfield/telenotion-bot/actions/workflows/goreleaser.yml)

[English](README.md)

这是一个将 Telegram 消息捕获并写入 Notion 数据库的机器人，支持 MarkdownV2 格式。

## 特性

- 从配置的 Telegram 聊天中捕获消息并存入 Notion
- 支持 `/start`、`/clean`、`/discard`、`/end` 命令
- 将 MarkdownV2 格式转换为 Notion 富文本
- 自动上传图片到 GitHub 仓库并嵌入到 Notion
- 使用时间格式作为 Notion 页面标题

## 配置

创建 TOML 配置文件（参考 `config.example.toml`）：

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

[log]
level = "info"
file = ""
```

## 使用

运行 bot：

```bash
telenotion-bot bot --config config.toml
```

### 命令

- `/start` - 开始新的捕获会话
- `/clean` - 清空当前缓存
- `/discard` - 放弃当前会话
- `/end` - 创建 Notion 页面并结束会话
- `/help` - 查看可用命令

### 支持的 Markdown 语法

- `*bold text*` - 粗体
- `_italic text_` - 斜体
- ` `code` `` - 行内代码
- `[text](url)` - 链接

### Notion 数据库要求

你的 Notion 数据库必须具有：
- `title` 属性（类型：Title）
- 不需要额外的内容属性（内容存储为子块）

## 开发

```bash
# 构建
go build ./...

# 测试
go test ./...

# 格式化代码
go fmt ./...
goimports -w .
```

## 许可证

MIT License
