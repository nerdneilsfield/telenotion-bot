package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/BurntSushi/toml"
)

type Config struct {
	Telegram Telegram `toml:"telegram"`
	Discord  Discord  `toml:"discord"`
	Notion   Notion   `toml:"notion"`
	GitHub   GitHub   `toml:"github"`
	Title    Title    `toml:"title"`
	Log      Log      `toml:"log"`
}

type Telegram struct {
	Token          string  `toml:"token"`
	AllowedChatIDs []int64 `toml:"allowed_chat_ids"`
}

type Discord struct {
	Token          string   `toml:"token"`
	AllowedUserIDs []string `toml:"allowed_user_ids"`
}

type Notion struct {
	Token          string `toml:"token"`
	DatabaseID     string `toml:"database_id"`
	TitleProperty  string `toml:"title_property"`
	OriginProperty string `toml:"origin_property"`
}

type GitHub struct {
	Token      string `toml:"token"`
	Repo       string `toml:"repo"`
	Branch     string `toml:"branch"`
	PathPrefix string `toml:"path_prefix"`
}

type Title struct {
	Timezone string `toml:"timezone"`
	Format   string `toml:"format"`
}

type Log struct {
	Level string `toml:"level"`
	File  string `toml:"file"`
}

// Environment variable names for config overrides
const (
	EnvTelegramToken      = "TELEGRAM_TOKEN"
	EnvTelegramAllowedIDs = "TELEGRAM_ALLOWED_CHAT_IDS"
	EnvDiscordToken       = "DISCORD_TOKEN"
	EnvDiscordAllowedIDs  = "DISCORD_ALLOWED_USER_IDS"
	EnvNotionToken        = "NOTION_TOKEN"
	EnvNotionDatabaseID   = "NOTION_DATABASE_ID"
	EnvNotionTitleProp    = "NOTION_TITLE_PROPERTY"
	EnvNotionOriginProp   = "NOTION_ORIGIN_PROPERTY"
	EnvGitHubToken        = "GITHUB_TOKEN"
	EnvGitHubRepo         = "GITHUB_REPO"
	EnvGitHubBranch       = "GITHUB_BRANCH"
	EnvGitHubPathPrefix   = "GITHUB_PATH_PREFIX"
	EnvTitleTimezone      = "TITLE_TIMEZONE"
	EnvTitleFormat        = "TITLE_FORMAT"
	EnvLogLevel           = "LOG_LEVEL"
	EnvLogFile            = "LOG_FILE"
)

func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var cfg Config
	if err := toml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	// Apply environment variable overrides
	cfg.applyEnvOverrides()

	cfg.Normalize()
	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	return &cfg, nil
}

// LoadFromEnv loads configuration from environment variables only.
// Useful for container deployments where config file is not available.
func LoadFromEnv() (*Config, error) {
	cfg := &Config{}

	// Apply environment variable overrides
	cfg.applyEnvOverrides()

	cfg.Normalize()
	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	return cfg, nil
}

// applyEnvOverrides applies environment variable overrides to config values.
func (c *Config) applyEnvOverrides() {
	// Telegram
	if v := os.Getenv(EnvTelegramToken); v != "" {
		c.Telegram.Token = v
	}
	if v := os.Getenv(EnvTelegramAllowedIDs); v != "" {
		c.Telegram.AllowedChatIDs = parseInt64List(v)
	}

	if v := os.Getenv(EnvDiscordToken); v != "" {
		c.Discord.Token = v
	}
	if v := os.Getenv(EnvDiscordAllowedIDs); v != "" {
		c.Discord.AllowedUserIDs = parseStringList(v)
	}

	// Notion
	if v := os.Getenv(EnvNotionToken); v != "" {
		c.Notion.Token = v
	}
	if v := os.Getenv(EnvNotionDatabaseID); v != "" {
		c.Notion.DatabaseID = v
	}
	if v := os.Getenv(EnvNotionTitleProp); v != "" {
		c.Notion.TitleProperty = v
	}
	if v := os.Getenv(EnvNotionOriginProp); v != "" {
		c.Notion.OriginProperty = v
	}

	// GitHub
	if v := os.Getenv(EnvGitHubToken); v != "" {
		c.GitHub.Token = v
	}
	if v := os.Getenv(EnvGitHubRepo); v != "" {
		c.GitHub.Repo = v
	}
	if v := os.Getenv(EnvGitHubBranch); v != "" {
		c.GitHub.Branch = v
	}
	if v := os.Getenv(EnvGitHubPathPrefix); v != "" {
		c.GitHub.PathPrefix = v
	}

	// Title
	if v := os.Getenv(EnvTitleTimezone); v != "" {
		c.Title.Timezone = v
	}
	if v := os.Getenv(EnvTitleFormat); v != "" {
		c.Title.Format = v
	}

	// Log
	if v := os.Getenv(EnvLogLevel); v != "" {
		c.Log.Level = v
	}
	if v := os.Getenv(EnvLogFile); v != "" {
		c.Log.File = v
	}
}

// parseInt64List parses a comma-separated list of int64 values.
func parseInt64List(s string) []int64 {
	if s == "" {
		return nil
	}
	parts := strings.Split(s, ",")
	result := make([]int64, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if v, err := strconv.ParseInt(p, 10, 64); err == nil {
			result = append(result, v)
		}
	}
	return result
}

func parseStringList(s string) []string {
	if s == "" {
		return nil
	}
	parts := strings.Split(s, ",")
	result := make([]string, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			result = append(result, p)
		}
	}
	return result
}

func (c *Config) Normalize() {
	if c.Title.Timezone == "" {
		c.Title.Timezone = "Asia/Shanghai"
	}
	if c.Title.Format == "" {
		c.Title.Format = "2006-01-02 15:04"
	}
	if c.Notion.TitleProperty == "" {
		c.Notion.TitleProperty = "Name"
	}
	if c.Notion.OriginProperty == "" {
		c.Notion.OriginProperty = "Origin"
	}
	if c.Log.Level == "" {
		c.Log.Level = "info"
	}
}

func (c *Config) Validate() error {
	telegramEnabled := c.Telegram.Token != "" || len(c.Telegram.AllowedChatIDs) > 0
	discordEnabled := c.Discord.Token != "" || len(c.Discord.AllowedUserIDs) > 0

	if !telegramEnabled && !discordEnabled {
		return fmt.Errorf("telegram or discord configuration is required")
	}

	if telegramEnabled {
		if c.Telegram.Token == "" {
			return fmt.Errorf("telegram.token is required")
		}
		if len(c.Telegram.AllowedChatIDs) == 0 {
			return fmt.Errorf("telegram.allowed_chat_ids is required")
		}
	}

	if discordEnabled {
		if c.Discord.Token == "" {
			return fmt.Errorf("discord.token is required")
		}
		if len(c.Discord.AllowedUserIDs) == 0 {
			return fmt.Errorf("discord.allowed_user_ids is required")
		}
	}

	if c.Notion.Token == "" {
		return fmt.Errorf("notion.token is required")
	}
	if c.Notion.DatabaseID == "" {
		return fmt.Errorf("notion.database_id is required")
	}
	if c.GitHub.Token == "" {
		return fmt.Errorf("github.token is required")
	}
	if c.GitHub.Repo == "" {
		return fmt.Errorf("github.repo is required")
	}
	if c.GitHub.Branch == "" {
		return fmt.Errorf("github.branch is required")
	}
	if _, err := time.LoadLocation(c.Title.Timezone); err != nil {
		return fmt.Errorf("title.timezone is invalid: %w", err)
	}

	return nil
}

func (t Title) Location() (*time.Location, error) {
	return time.LoadLocation(t.Timezone)
}

func (t Title) FormatTime(loc *time.Location) string {
	return time.Now().In(loc).Format(t.Format)
}
