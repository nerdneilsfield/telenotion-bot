package config

import (
	"os"
	"testing"
	"time"
)

func TestLoad_ValidConfig(t *testing.T) {
	// Create a valid config file
	content := `
[telegram]
token = "test-token"
allowed_chat_ids = [123456789]

[notion]
token = "notion-token"
database_id = "database-id"
title_property = "Name"
origin_property = "Origin"

[github]
token = "github-token"
repo = "owner/repo"
branch = "main"
telegram_branch = "tg"
discord_branch = "discord"
path_prefix = "images/"

[title]
timezone = "Asia/Shanghai"
format = "2006-01-02 15:04"

[log]
level = "info"
file = ""
`
	tmpFile, err := os.CreateTemp("", "config-*.toml")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	if _, err := tmpFile.WriteString(content); err != nil {
		t.Fatalf("Failed to write temp file: %v", err)
	}
	tmpFile.Close()

	cfg, err := Load(tmpFile.Name())
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	// Verify values
	if cfg.Telegram.Token != "test-token" {
		t.Errorf("Telegram.Token = %q, want %q", cfg.Telegram.Token, "test-token")
	}
	if cfg.Telegram.AllowedChatIDs[0] != 123456789 {
		t.Errorf("AllowedChatIDs = %v, want %v", cfg.Telegram.AllowedChatIDs, []int64{123456789})
	}
	if cfg.Notion.DatabaseID != "database-id" {
		t.Errorf("DatabaseID = %q, want %q", cfg.Notion.DatabaseID, "database-id")
	}
	if cfg.Notion.OriginProperty != "Origin" {
		t.Errorf("OriginProperty = %q, want %q", cfg.Notion.OriginProperty, "Origin")
	}
	if cfg.GitHub.Repo != "owner/repo" {
		t.Errorf("GitHub.Repo = %q, want %q", cfg.GitHub.Repo, "owner/repo")
	}
	if cfg.GitHub.TelegramBranch != "tg" {
		t.Errorf("GitHub.TelegramBranch = %q, want %q", cfg.GitHub.TelegramBranch, "tg")
	}
	if cfg.GitHub.DiscordBranch != "discord" {
		t.Errorf("GitHub.DiscordBranch = %q, want %q", cfg.GitHub.DiscordBranch, "discord")
	}
	if cfg.Title.Timezone != "Asia/Shanghai" {
		t.Errorf("Timezone = %q, want %q", cfg.Title.Timezone, "Asia/Shanghai")
	}
}

func TestLoad_FileNotFound(t *testing.T) {
	_, err := Load("/nonexistent/path/config.toml")
	if err == nil {
		t.Error("Load() should return error for missing file")
	}
}

func TestLoad_InvalidTOML(t *testing.T) {
	tmpFile, err := os.CreateTemp("", "config-*.toml")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	// Write invalid TOML
	if _, err := tmpFile.WriteString("invalid toml content ["); err != nil {
		t.Fatalf("Failed to write temp file: %v", err)
	}
	tmpFile.Close()

	_, err = Load(tmpFile.Name())
	if err == nil {
		t.Error("Load() should return error for invalid TOML")
	}
}

func TestNormalize(t *testing.T) {
	cfg := &Config{
		Telegram: Telegram{Token: "token", AllowedChatIDs: []int64{1}},
		Notion:   Notion{Token: "token", DatabaseID: "id"},
		GitHub:   GitHub{Token: "token", Repo: "repo", Branch: "main"},
		Title:    Title{},
		Log:      Log{},
	}

	cfg.Normalize()

	if cfg.Title.Timezone != "Asia/Shanghai" {
		t.Errorf("Timezone = %q, want %q", cfg.Title.Timezone, "Asia/Shanghai")
	}
	if cfg.Title.Format != "2006-01-02 15:04" {
		t.Errorf("Format = %q, want %q", cfg.Title.Format, "2006-01-02 15:04")
	}
	if cfg.Notion.TitleProperty != "Name" {
		t.Errorf("TitleProperty = %q, want %q", cfg.Notion.TitleProperty, "Name")
	}
	if cfg.Notion.OriginProperty != "Origin" {
		t.Errorf("OriginProperty = %q, want %q", cfg.Notion.OriginProperty, "Origin")
	}
	if cfg.Log.Level != "info" {
		t.Errorf("Log.Level = %q, want %q", cfg.Log.Level, "info")
	}
	if cfg.Media.MaxImageSizeMB != 20 {
		t.Errorf("Media.MaxImageSizeMB = %d, want %d", cfg.Media.MaxImageSizeMB, 20)
	}
	if len(cfg.Media.AllowedImageTypes) == 0 {
		t.Error("Media.AllowedImageTypes should not be empty")
	}
}

func TestGitHubBranchFallback(t *testing.T) {
	g := GitHub{Branch: "main"}
	if g.BranchForTelegram() != "main" {
		t.Errorf("BranchForTelegram = %q, want %q", g.BranchForTelegram(), "main")
	}
	if g.BranchForDiscord() != "main" {
		t.Errorf("BranchForDiscord = %q, want %q", g.BranchForDiscord(), "main")
	}

	g.TelegramBranch = "tg"
	g.DiscordBranch = "dc"
	if g.BranchForTelegram() != "tg" {
		t.Errorf("BranchForTelegram = %q, want %q", g.BranchForTelegram(), "tg")
	}
	if g.BranchForDiscord() != "dc" {
		t.Errorf("BranchForDiscord = %q, want %q", g.BranchForDiscord(), "dc")
	}
}

func TestValidate_MissingToken(t *testing.T) {
	tests := []struct {
		name    string
		cfg     Config
		wantErr string
	}{
		{
			name: "no platform configured",
			cfg: Config{
				Notion: Notion{Token: "token", DatabaseID: "id"},
				GitHub: GitHub{Token: "token", Repo: "repo", Branch: "main"},
				Title:  Title{Timezone: "UTC"},
			},
			wantErr: "telegram or discord configuration is required",
		},
		{
			name: "missing telegram token",
			cfg: Config{
				Telegram: Telegram{AllowedChatIDs: []int64{1}},
				Notion:   Notion{Token: "token", DatabaseID: "id"},
				GitHub:   GitHub{Token: "token", Repo: "repo", Branch: "main"},
				Title:    Title{Timezone: "UTC"},
			},
			wantErr: "telegram.token is required",
		},
		{
			name: "missing allowed chat ids",
			cfg: Config{
				Telegram: Telegram{Token: "token"},
				Notion:   Notion{Token: "token", DatabaseID: "id"},
				GitHub:   GitHub{Token: "token", Repo: "repo", Branch: "main"},
				Title:    Title{Timezone: "UTC"},
			},
			wantErr: "telegram.allowed_chat_ids is required",
		},
		{
			name: "missing discord token",
			cfg: Config{
				Discord: Discord{AllowedUserIDs: []string{"user-1"}},
				Notion:  Notion{Token: "token", DatabaseID: "id"},
				GitHub:  GitHub{Token: "token", Repo: "repo", Branch: "main"},
				Title:   Title{Timezone: "UTC"},
			},
			wantErr: "discord.token is required",
		},
		{
			name: "missing discord allowlist",
			cfg: Config{
				Discord: Discord{Token: "token"},
				Notion:  Notion{Token: "token", DatabaseID: "id"},
				GitHub:  GitHub{Token: "token", Repo: "repo", Branch: "main"},
				Title:   Title{Timezone: "UTC"},
			},
			wantErr: "discord.allowed_user_ids is required",
		},
		{
			name: "missing notion token",
			cfg: Config{
				Telegram: Telegram{Token: "token", AllowedChatIDs: []int64{1}},
				Notion:   Notion{DatabaseID: "id"},
				GitHub:   GitHub{Token: "token", Repo: "repo", Branch: "main"},
				Title:    Title{Timezone: "UTC"},
			},
			wantErr: "notion.token is required",
		},
		{
			name: "missing database id",
			cfg: Config{
				Telegram: Telegram{Token: "token", AllowedChatIDs: []int64{1}},
				Notion:   Notion{Token: "token"},
				GitHub:   GitHub{Token: "token", Repo: "repo", Branch: "main"},
				Title:    Title{Timezone: "UTC"},
			},
			wantErr: "notion.database_id is required",
		},
		{
			name: "missing github token",
			cfg: Config{
				Telegram: Telegram{Token: "token", AllowedChatIDs: []int64{1}},
				Notion:   Notion{Token: "token", DatabaseID: "id"},
				GitHub:   GitHub{Repo: "repo", Branch: "main"},
				Title:    Title{Timezone: "UTC"},
			},
			wantErr: "github.token is required",
		},
		{
			name: "missing github repo",
			cfg: Config{
				Telegram: Telegram{Token: "token", AllowedChatIDs: []int64{1}},
				Notion:   Notion{Token: "token", DatabaseID: "id"},
				GitHub:   GitHub{Token: "token", Branch: "main"},
				Title:    Title{Timezone: "UTC"},
			},
			wantErr: "github.repo is required",
		},
		{
			name: "missing github branch for telegram",
			cfg: Config{
				Telegram: Telegram{Token: "token", AllowedChatIDs: []int64{1}},
				Notion:   Notion{Token: "token", DatabaseID: "id"},
				GitHub:   GitHub{Token: "token", Repo: "repo"},
				Title:    Title{Timezone: "UTC"},
			},
			wantErr: "github.telegram_branch or github.branch is required",
		},
		{
			name: "missing github branch for discord",
			cfg: Config{
				Discord: Discord{Token: "token", AllowedUserIDs: []string{"user-1"}},
				Notion:  Notion{Token: "token", DatabaseID: "id"},
				GitHub:  GitHub{Token: "token", Repo: "repo"},
				Title:   Title{Timezone: "UTC"},
			},
			wantErr: "github.discord_branch or github.branch is required",
		},
		{
			name: "invalid timezone",
			cfg: Config{
				Telegram: Telegram{Token: "token", AllowedChatIDs: []int64{1}},
				Notion:   Notion{Token: "token", DatabaseID: "id"},
				GitHub:   GitHub{Token: "token", Repo: "repo", Branch: "main"},
				Title:    Title{Timezone: "Invalid/Timezone"},
			},
			wantErr: "title.timezone is invalid: unknown time zone Invalid/Timezone",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.cfg.Validate()
			if err == nil {
				t.Errorf("Validate() should return error")
				return
			}
			if err.Error() != tt.wantErr {
				t.Errorf("Validate() error = %q, want %q", err.Error(), tt.wantErr)
			}
		})
	}
}

func TestValidate_Valid(t *testing.T) {
	cfg := Config{
		Telegram: Telegram{Token: "token", AllowedChatIDs: []int64{1}},
		Notion:   Notion{Token: "token", DatabaseID: "id", OriginProperty: "Origin"},
		GitHub:   GitHub{Token: "token", Repo: "repo", Branch: "main"},
		Title:    Title{Timezone: "UTC"},
	}

	err := cfg.Validate()
	if err != nil {
		t.Errorf("Validate() unexpected error = %v", err)
	}
}

func TestValidate_ValidTelegramBranchOverride(t *testing.T) {
	cfg := Config{
		Telegram: Telegram{Token: "token", AllowedChatIDs: []int64{1}},
		Notion:   Notion{Token: "token", DatabaseID: "id"},
		GitHub:   GitHub{Token: "token", Repo: "repo", TelegramBranch: "telegram"},
		Title:    Title{Timezone: "UTC"},
	}

	err := cfg.Validate()
	if err != nil {
		t.Errorf("Validate() unexpected error = %v", err)
	}
}

func TestValidate_ValidDiscordOnly(t *testing.T) {
	cfg := Config{
		Discord: Discord{Token: "token", AllowedUserIDs: []string{"user-1"}},
		Notion:  Notion{Token: "token", DatabaseID: "id"},
		GitHub:  GitHub{Token: "token", Repo: "repo", Branch: "main"},
		Title:   Title{Timezone: "UTC"},
	}

	err := cfg.Validate()
	if err != nil {
		t.Errorf("Validate() unexpected error = %v", err)
	}
}

func TestValidate_ValidDiscordBranchOverride(t *testing.T) {
	cfg := Config{
		Discord: Discord{Token: "token", AllowedUserIDs: []string{"user-1"}},
		Notion:  Notion{Token: "token", DatabaseID: "id"},
		GitHub:  GitHub{Token: "token", Repo: "repo", DiscordBranch: "discord"},
		Title:   Title{Timezone: "UTC"},
	}

	err := cfg.Validate()
	if err != nil {
		t.Errorf("Validate() unexpected error = %v", err)
	}
}

func TestTitleLocation(t *testing.T) {
	title := Title{Timezone: "America/New_York"}

	loc, err := title.Location()
	if err != nil {
		t.Fatalf("Location() error = %v", err)
	}

	if loc == nil {
		t.Error("Location() returned nil")
	}
}

func TestTitleLocation_Invalid(t *testing.T) {
	title := Title{Timezone: "Invalid/Zone"}

	_, err := title.Location()
	if err == nil {
		t.Error("Location() should return error for invalid timezone")
	}
}

func TestTitleFormatTime(t *testing.T) {
	title := Title{Timezone: "UTC", Format: "2006-01-02"}

	loc, _ := time.LoadLocation("UTC")
	result := title.FormatTime(loc)

	// Result should be in YYYY-MM-DD format
	if len(result) != 10 {
		t.Errorf("FormatTime() result = %q, length = %d, want 10", result, len(result))
	}
}

func TestApplyEnvOverrides(t *testing.T) {
	// Set environment variables
	os.Setenv(EnvTelegramToken, "env-telegram-token")
	os.Setenv(EnvTelegramAllowedIDs, "111,222,333")
	os.Setenv(EnvDiscordToken, "env-discord-token")
	os.Setenv(EnvDiscordAllowedIDs, "user-1,user-2")
	os.Setenv(EnvNotionToken, "env-notion-token")
	os.Setenv(EnvNotionDatabaseID, "env-database-id")
	os.Setenv(EnvNotionTitleProp, "Title")
	os.Setenv(EnvNotionOriginProp, "Origin")
	os.Setenv(EnvGitHubToken, "env-github-token")
	os.Setenv(EnvGitHubRepo, "env-owner/env-repo")
	os.Setenv(EnvGitHubBranch, "develop")
	os.Setenv(EnvGitHubTelegramBranch, "env-tg")
	os.Setenv(EnvGitHubDiscordBranch, "env-discord")
	os.Setenv(EnvGitHubPathPrefix, "env-images/")
	os.Setenv(EnvMediaMaxImageSizeMB, "25")
	os.Setenv(EnvMediaAllowedTypes, "image/jpeg,image/png")
	os.Setenv(EnvTitleTimezone, "America/New_York")
	os.Setenv(EnvTitleFormat, "2006-01-02")
	os.Setenv(EnvLogLevel, "debug")
	os.Setenv(EnvLogFile, "env-log.log")

	defer func() {
		// Clean up environment variables
		os.Unsetenv(EnvTelegramToken)
		os.Unsetenv(EnvTelegramAllowedIDs)
		os.Unsetenv(EnvDiscordToken)
		os.Unsetenv(EnvDiscordAllowedIDs)
		os.Unsetenv(EnvNotionToken)
		os.Unsetenv(EnvNotionDatabaseID)
		os.Unsetenv(EnvNotionTitleProp)
		os.Unsetenv(EnvNotionOriginProp)
		os.Unsetenv(EnvGitHubToken)
		os.Unsetenv(EnvGitHubRepo)
		os.Unsetenv(EnvGitHubBranch)
		os.Unsetenv(EnvGitHubTelegramBranch)
		os.Unsetenv(EnvGitHubDiscordBranch)
		os.Unsetenv(EnvGitHubPathPrefix)
		os.Unsetenv(EnvMediaMaxImageSizeMB)
		os.Unsetenv(EnvMediaAllowedTypes)
		os.Unsetenv(EnvTitleTimezone)
		os.Unsetenv(EnvTitleFormat)
		os.Unsetenv(EnvLogLevel)
		os.Unsetenv(EnvLogFile)
	}()

	cfg := &Config{
		Telegram: Telegram{Token: "file-token", AllowedChatIDs: []int64{123}},
		Notion:   Notion{Token: "file-notion", DatabaseID: "file-db"},
		GitHub:   GitHub{Token: "file-github", Repo: "file/repo", Branch: "main"},
		Title:    Title{Timezone: "UTC", Format: "2006"},
		Log:      Log{Level: "info", File: ""},
	}

	cfg.applyEnvOverrides()

	// Verify environment overrides were applied
	if cfg.Telegram.Token != "env-telegram-token" {
		t.Errorf("Telegram.Token = %q, want %q", cfg.Telegram.Token, "env-telegram-token")
	}
	if len(cfg.Telegram.AllowedChatIDs) != 3 {
		t.Errorf("AllowedChatIDs length = %d, want 3", len(cfg.Telegram.AllowedChatIDs))
	}
	if cfg.Discord.Token != "env-discord-token" {
		t.Errorf("Discord.Token = %q, want %q", cfg.Discord.Token, "env-discord-token")
	}
	if len(cfg.Discord.AllowedUserIDs) != 2 {
		t.Errorf("AllowedUserIDs length = %d, want 2", len(cfg.Discord.AllowedUserIDs))
	}
	if cfg.Notion.Token != "env-notion-token" {
		t.Errorf("Notion.Token = %q, want %q", cfg.Notion.Token, "env-notion-token")
	}
	if cfg.Notion.DatabaseID != "env-database-id" {
		t.Errorf("DatabaseID = %q, want %q", cfg.Notion.DatabaseID, "env-database-id")
	}
	if cfg.Notion.OriginProperty != "Origin" {
		t.Errorf("OriginProperty = %q, want %q", cfg.Notion.OriginProperty, "Origin")
	}
	if cfg.GitHub.Token != "env-github-token" {
		t.Errorf("GitHub.Token = %q, want %q", cfg.GitHub.Token, "env-github-token")
	}
	if cfg.GitHub.Repo != "env-owner/env-repo" {
		t.Errorf("GitHub.Repo = %q, want %q", cfg.GitHub.Repo, "env-owner/env-repo")
	}
	if cfg.GitHub.Branch != "develop" {
		t.Errorf("GitHub.Branch = %q, want %q", cfg.GitHub.Branch, "develop")
	}
	if cfg.GitHub.TelegramBranch != "env-tg" {
		t.Errorf("GitHub.TelegramBranch = %q, want %q", cfg.GitHub.TelegramBranch, "env-tg")
	}
	if cfg.GitHub.DiscordBranch != "env-discord" {
		t.Errorf("GitHub.DiscordBranch = %q, want %q", cfg.GitHub.DiscordBranch, "env-discord")
	}
	if cfg.Media.MaxImageSizeMB != 25 {
		t.Errorf("Media.MaxImageSizeMB = %d, want %d", cfg.Media.MaxImageSizeMB, 25)
	}
	if len(cfg.Media.AllowedImageTypes) != 2 {
		t.Errorf("Media.AllowedImageTypes length = %d, want %d", len(cfg.Media.AllowedImageTypes), 2)
	}
	if cfg.Title.Timezone != "America/New_York" {
		t.Errorf("Title.Timezone = %q, want %q", cfg.Title.Timezone, "America/New_York")
	}
	if cfg.Log.Level != "debug" {
		t.Errorf("Log.Level = %q, want %q", cfg.Log.Level, "debug")
	}
}

func TestParseInt64List(t *testing.T) {
	tests := []struct {
		input    string
		expected []int64
	}{
		{"123,456,789", []int64{123, 456, 789}},
		{"  123 , 456  ", []int64{123, 456}},
		{"", nil},
		{"invalid", nil},
		{"123,abc,456", []int64{123, 456}},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := parseInt64List(tt.input)
			if len(result) != len(tt.expected) {
				t.Errorf("parseInt64List(%q) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}

func TestParseStringList(t *testing.T) {
	tests := []struct {
		input    string
		expected []string
	}{
		{"a,b,c", []string{"a", "b", "c"}},
		{"  a ,  b  ", []string{"a", "b"}},
		{"", nil},
		{",,", nil},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := parseStringList(tt.input)
			if len(result) != len(tt.expected) {
				t.Errorf("parseStringList(%q) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}

func TestLoadFromEnv(t *testing.T) {
	// Set required environment variables
	os.Setenv(EnvTelegramToken, "env-token")
	os.Setenv(EnvTelegramAllowedIDs, "123")
	os.Setenv(EnvNotionToken, "notion-token")
	os.Setenv(EnvNotionDatabaseID, "database-id")
	os.Setenv(EnvGitHubToken, "github-token")
	os.Setenv(EnvGitHubRepo, "owner/repo")
	os.Setenv(EnvGitHubBranch, "main")

	defer func() {
		os.Unsetenv(EnvTelegramToken)
		os.Unsetenv(EnvTelegramAllowedIDs)
		os.Unsetenv(EnvNotionToken)
		os.Unsetenv(EnvNotionDatabaseID)
		os.Unsetenv(EnvGitHubToken)
		os.Unsetenv(EnvGitHubRepo)
		os.Unsetenv(EnvGitHubBranch)
	}()

	cfg, err := LoadFromEnv()
	if err != nil {
		t.Fatalf("LoadFromEnv() error = %v", err)
	}

	if cfg.Telegram.Token != "env-token" {
		t.Errorf("Telegram.Token = %q, want %q", cfg.Telegram.Token, "env-token")
	}
	if cfg.Notion.DatabaseID != "database-id" {
		t.Errorf("DatabaseID = %q, want %q", cfg.Notion.DatabaseID, "database-id")
	}
	if cfg.Notion.OriginProperty != "Origin" {
		t.Errorf("OriginProperty = %q, want %q", cfg.Notion.OriginProperty, "Origin")
	}
}

func TestLoadFromEnvDiscordOnly(t *testing.T) {
	os.Setenv(EnvDiscordToken, "discord-token")
	os.Setenv(EnvDiscordAllowedIDs, "user-1")
	os.Setenv(EnvNotionToken, "notion-token")
	os.Setenv(EnvNotionDatabaseID, "database-id")
	os.Setenv(EnvGitHubToken, "github-token")
	os.Setenv(EnvGitHubRepo, "owner/repo")
	os.Setenv(EnvGitHubBranch, "main")

	defer func() {
		os.Unsetenv(EnvDiscordToken)
		os.Unsetenv(EnvDiscordAllowedIDs)
		os.Unsetenv(EnvNotionToken)
		os.Unsetenv(EnvNotionDatabaseID)
		os.Unsetenv(EnvGitHubToken)
		os.Unsetenv(EnvGitHubRepo)
		os.Unsetenv(EnvGitHubBranch)
	}()

	cfg, err := LoadFromEnv()
	if err != nil {
		t.Fatalf("LoadFromEnv() error = %v", err)
	}

	if cfg.Discord.Token != "discord-token" {
		t.Errorf("Discord.Token = %q, want %q", cfg.Discord.Token, "discord-token")
	}
	if cfg.Notion.DatabaseID != "database-id" {
		t.Errorf("DatabaseID = %q, want %q", cfg.Notion.DatabaseID, "database-id")
	}
	if cfg.Notion.OriginProperty != "Origin" {
		t.Errorf("OriginProperty = %q, want %q", cfg.Notion.OriginProperty, "Origin")
	}
}
