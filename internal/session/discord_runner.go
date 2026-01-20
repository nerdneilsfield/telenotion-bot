package session

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/jomei/notionapi"
	"github.com/nerdneilsfield/telenotion-bot/internal/config"
	"github.com/nerdneilsfield/telenotion-bot/internal/discordclient"
	"github.com/nerdneilsfield/telenotion-bot/internal/github"
	"github.com/nerdneilsfield/telenotion-bot/internal/notion"
	"go.uber.org/zap"
)

type DiscordRunner struct {
	cfg          *config.Config
	discord      *discordclient.Client
	notion       *notion.Client
	github       *github.Client
	stateMachine *StateMachine
	ctx          context.Context
	logger       *zap.Logger
}

const discordHelpText = "Commands:\n/start - start a new capture session\n/clean - clear the current buffer\n/discard - abandon the current session\n/end - create a Notion page and end session\n/help - show this help"

func NewDiscordRunner(cfg *config.Config, logger *zap.Logger) (*DiscordRunner, error) {
	client, err := discordclient.NewClient(cfg.Discord.Token)
	if err != nil {
		return nil, err
	}

	return &DiscordRunner{
		cfg:          cfg,
		discord:      client,
		notion:       notion.NewClient(cfg.Notion.Token),
		github:       github.NewClient(cfg.GitHub.Token, cfg.GitHub.Repo, cfg.GitHub.Branch, cfg.GitHub.PathPrefix),
		stateMachine: NewStateMachine(),
		logger:       logger,
	}, nil
}

func (r *DiscordRunner) Run(ctx context.Context) error {
	r.ctx = ctx
	session := r.discord.Session()
	session.AddHandler(r.handleInteraction)
	session.AddHandler(r.handleMessage)

	if err := r.discord.Open(); err != nil {
		return err
	}
	defer r.discord.Close()

	if err := r.discord.RegisterCommands(r.commands()); err != nil {
		if r.logger != nil {
			r.logger.Warn("failed to register discord commands", zap.Error(err))
		}
	}

	<-ctx.Done()
	return nil
}

func (r *DiscordRunner) commands() []*discordgo.ApplicationCommand {
	dmPermission := true

	return []*discordgo.ApplicationCommand{
		{Name: "start", Description: "Start a new capture session", DMPermission: &dmPermission},
		{Name: "clean", Description: "Clear the current buffer", DMPermission: &dmPermission},
		{Name: "discard", Description: "Discard the current session", DMPermission: &dmPermission},
		{Name: "end", Description: "Save to Notion and end session", DMPermission: &dmPermission},
		{Name: "help", Description: "Show available commands", DMPermission: &dmPermission},
	}
}

func (r *DiscordRunner) handleInteraction(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if i.Type != discordgo.InteractionApplicationCommand {
		return
	}

	userID := interactionUserID(i)
	if userID == "" {
		return
	}
	if !discordclient.IsAllowedUser(r.cfg.Discord.AllowedUserIDs, userID) {
		return
	}
	if i.GuildID != "" {
		r.respondInteraction(s, i, "Use these commands in a DM.", true)
		return
	}

	chatID, err := strconv.ParseInt(userID, 10, 64)
	if err != nil {
		return
	}

	command := i.ApplicationCommandData().Name
	response := ""

	switch command {
	case "help":
		response = discordHelpText
	case "start":
		if !r.stateMachine.IsActive(chatID) {
			r.stateMachine.StartSession(chatID)
			response = "Session started. Send messages or images, then /end to save."
		} else {
			response = "Session already active. Use /clean to clear or /end to save."
		}
	case "clean":
		if r.stateMachine.IsActive(chatID) {
			r.stateMachine.ClearSession(chatID)
			response = "Buffer cleared. Continue sending messages."
		} else {
			response = "No active session. Use /start first."
		}
	case "discard":
		if r.stateMachine.IsActive(chatID) {
			r.stateMachine.DiscardSession(chatID)
			response = "Session discarded."
		} else {
			response = "No active session to discard."
		}
	case "end":
		if r.stateMachine.IsActive(chatID) {
			session, ok := r.stateMachine.EndSession(chatID)
			if ok {
				if err := r.createNotionPage(r.ctx, session); err != nil {
					response = "Failed to save. Check logs for details."
					if r.logger != nil {
						r.logger.Error("failed to create notion page", zap.String("user_id", userID), zap.Error(err))
					}
				} else {
					response = "Saved to Notion."
				}
			}
		} else {
			response = "No active session. Use /start first."
		}
	default:
		return
	}

	if response != "" {
		r.respondInteraction(s, i, response, false)
	}
}

func (r *DiscordRunner) handleMessage(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author == nil || m.Author.Bot {
		return
	}
	if m.GuildID != "" {
		return
	}
	if !discordclient.IsAllowedUser(r.cfg.Discord.AllowedUserIDs, m.Author.ID) {
		return
	}

	chatID, err := strconv.ParseInt(m.Author.ID, 10, 64)
	if err != nil {
		return
	}
	if !r.stateMachine.IsActive(chatID) {
		r.stateMachine.StartSession(chatID)
	}

	r.collectMessage(m, chatID)
}

func (r *DiscordRunner) collectMessage(msg *discordgo.MessageCreate, chatID int64) {
	content := strings.TrimSpace(discordclient.NormalizeContent(msg.Content, msg.Mentions))
	if content != "" {
		if code := extractDiscordCodeBlock(content); code != nil {
			r.stateMachine.AppendBlock(chatID, *code)
		} else if richText := discordclient.ContentToRichText(content); len(richText) > 0 {
			r.stateMachine.AppendBlock(chatID, TextBlock{RichText: richText})
		}
	}

	for _, attachment := range msg.Attachments {
		if !isImageAttachment(attachment) {
			continue
		}
		r.stateMachine.AppendBlock(chatID, ImageBlock{FileURL: attachment.URL, Filename: attachment.Filename})
	}
}

func (r *DiscordRunner) respondInteraction(s *discordgo.Session, i *discordgo.InteractionCreate, content string, ephemeral bool) {
	flags := discordgo.MessageFlags(0)
	if ephemeral {
		flags = discordgo.MessageFlagsEphemeral
	}

	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: content,
			Flags:   flags,
		},
	})
	if err != nil && r.logger != nil {
		r.logger.Warn("failed to respond to interaction", zap.Error(err))
	}
}

func (r *DiscordRunner) createNotionPage(ctx context.Context, session *Session) error {
	loc, err := r.cfg.Title.Location()
	if err != nil {
		return err
	}

	blocks := make([]notionapi.Block, 0, len(session.Blocks))

	for _, block := range session.Blocks {
		switch b := block.(type) {
		case TextBlock:
			richTexts := splitRichTextEntries(b.RichText)
			for _, chunk := range chunkRichText(richTexts, notionRichTextBlockLimit) {
				blocks = append(blocks, &notionapi.ParagraphBlock{
					BasicBlock: notionapi.BasicBlock{Object: "block", Type: "paragraph"},
					Paragraph:  notionapi.Paragraph{RichText: chunk},
				})
			}
		case CodeBlock:
			language := b.Language
			if language == "" {
				language = "plain text"
			}
			richTexts := splitRichTextEntries([]notionapi.RichText{{Type: "text", Text: &notionapi.Text{Content: b.Content}}})
			for _, chunk := range chunkRichText(richTexts, notionRichTextBlockLimit) {
				blocks = append(blocks, &notionapi.CodeBlock{
					BasicBlock: notionapi.BasicBlock{Object: "block", Type: "code"},
					Code: notionapi.Code{
						RichText: chunk,
						Language: language,
					},
				})
			}
		case ImageBlock:
			if b.FileURL == "" {
				continue
			}
			data, err := downloadURL(b.FileURL)
			if err != nil {
				return err
			}

			filename := b.Filename
			if filename == "" {
				if b.FileID != "" {
					filename = b.FileID + ".jpg"
				} else {
					filename = "discord-image.jpg"
				}
			}
			rawURL, err := r.github.UploadImage(ctx, data, filename)
			if err != nil {
				return err
			}

			image := &notionapi.ImageBlock{
				BasicBlock: notionapi.BasicBlock{Object: "block", Type: "image"},
				Image:      notionapi.Image{External: &notionapi.FileObject{URL: rawURL}},
			}
			if b.Caption != "" {
				image.Image.Caption = []notionapi.RichText{{Type: "text", Text: &notionapi.Text{Content: b.Caption}}}
			}
			blocks = append(blocks, image)
		}
	}

	if len(blocks) == 0 {
		return fmt.Errorf("no blocks to save")
	}

	title := r.cfg.Title.FormatTime(loc)
	if r.logger != nil {
		r.logger.Info("creating notion page", zap.String("origin_property", r.cfg.Notion.OriginProperty), zap.String("origin", "Discord"))
	}
	_, err = r.notion.CreatePage(ctx, r.cfg.Notion.DatabaseID, r.cfg.Notion.TitleProperty, r.cfg.Notion.OriginProperty, title, "Discord", blocks)
	return err
}

func extractDiscordCodeBlock(content string) *CodeBlock {
	if !strings.HasPrefix(content, "```") || !strings.HasSuffix(content, "```") {
		return nil
	}

	trimmed := strings.TrimSuffix(strings.TrimPrefix(content, "```"), "```")
	trimmed = strings.TrimSpace(trimmed)
	if trimmed == "" {
		return nil
	}

	lines := strings.SplitN(trimmed, "\n", 2)
	language := ""
	body := trimmed
	if len(lines) == 2 && !strings.Contains(lines[0], " ") {
		language = strings.TrimSpace(lines[0])
		body = lines[1]
	}

	return &CodeBlock{Content: body, Language: language}
}

func isImageAttachment(attachment *discordgo.MessageAttachment) bool {
	if attachment == nil {
		return false
	}
	if attachment.ContentType != "" {
		return strings.HasPrefix(attachment.ContentType, "image/")
	}
	return attachment.Width > 0 && attachment.Height > 0
}

func interactionUserID(i *discordgo.InteractionCreate) string {
	if i.User != nil {
		return i.User.ID
	}
	if i.Member != nil && i.Member.User != nil {
		return i.Member.User.ID
	}
	return ""
}
