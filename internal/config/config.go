package config

import (
	"fmt"
	"os"
	"time"

	"github.com/BurntSushi/toml"
)

type Config struct {
	Telegram Telegram `toml:"telegram"`
	Notion   Notion   `toml:"notion"`
	GitHub   GitHub   `toml:"github"`
	Title    Title    `toml:"title"`
	Log      Log      `toml:"log"`
}

type Telegram struct {
	Token          string  `toml:"token"`
	AllowedChatIDs []int64 `toml:"allowed_chat_ids"`
}

type Notion struct {
	Token         string `toml:"token"`
	DatabaseID    string `toml:"database_id"`
	TitleProperty string `toml:"title_property"`
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

func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var cfg Config
	if err := toml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	cfg.Normalize()
	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	return &cfg, nil
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
	if c.Log.Level == "" {
		c.Log.Level = "info"
	}
}

func (c *Config) Validate() error {
	if c.Telegram.Token == "" {
		return fmt.Errorf("telegram.token is required")
	}
	if len(c.Telegram.AllowedChatIDs) == 0 {
		return fmt.Errorf("telegram.allowed_chat_ids is required")
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
