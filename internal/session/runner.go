package session

import (
	"context"
	"fmt"

	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/jomei/notionapi"
	"github.com/nerdneilsfield/telenotion-bot/internal/config"
	"github.com/nerdneilsfield/telenotion-bot/internal/github"
	"github.com/nerdneilsfield/telenotion-bot/internal/notion"
	"github.com/nerdneilsfield/telenotion-bot/internal/tgclient"
	"go.uber.org/zap"
)

type Runner struct {
	cfg          *config.Config
	telegram     *tgclient.Client
	notion       *notion.Client
	github       *github.Client
	stateMachine *StateMachine
	mapper       *tgclient.Mapper
	logger       *zap.Logger
}

const helpText = "Commands:\n/start - start a new capture session\n/clean - clear the current buffer\n/discard - abandon the current session\n/end - create a Notion page and end session\n/help - show this help"

func NewRunner(cfg *config.Config, logger *zap.Logger) (*Runner, error) {
	client, err := tgclient.NewClient(cfg.Telegram.Token)
	if err != nil {
		return nil, err
	}

	return &Runner{
		cfg:          cfg,
		telegram:     client,
		notion:       notion.NewClient(cfg.Notion.Token),
		github:       github.NewClient(cfg.GitHub.Token, cfg.GitHub.Repo, cfg.GitHub.Branch, cfg.GitHub.PathPrefix),
		stateMachine: NewStateMachine(),
		mapper:       tgclient.NewMapper(),
		logger:       logger,
	}, nil
}

func (r *Runner) Run(ctx context.Context) error {
	if err := r.registerCommands(); err != nil {
		if r.logger != nil {
			r.logger.Warn("failed to register bot commands", zap.Error(err))
		}
	}

	offset := 0

	for {
		select {
		case <-ctx.Done():
			return nil
		default:
		}

		updates, err := r.telegram.GetUpdates(offset, 60)
		if err != nil {
			return err
		}

		for _, update := range updates {
			if update.UpdateID >= offset {
				offset = update.UpdateID + 1
			}

			if err := r.handleUpdate(ctx, update); err != nil {
				return err
			}
		}
	}
}

func (r *Runner) handleUpdate(ctx context.Context, update tgbotapi.Update) error {
	if update.Message == nil {
		return nil
	}

	chatID := update.Message.Chat.ID
	if !r.telegram.IsAllowedChat(r.cfg.Telegram.AllowedChatIDs, chatID) {
		return nil
	}

	switch update.Message.Text {
	case "/help":
		r.reply(chatID, helpText)
	case "/start":
		if !r.stateMachine.IsActive(chatID) {
			r.stateMachine.StartSession(chatID)
			r.reply(chatID, "Session started. Send messages or images, then /end to save.")
			if r.logger != nil {
				r.logger.Info("session started", zap.Int64("chat_id", chatID))
			}
		} else {
			r.reply(chatID, "Session already active. Use /clean to clear or /end to save.")
		}
	case "/clean":
		if r.stateMachine.IsActive(chatID) {
			r.stateMachine.ClearSession(chatID)
			r.reply(chatID, "Buffer cleared. Continue sending messages.")
			if r.logger != nil {
				r.logger.Info("session cleared", zap.Int64("chat_id", chatID))
			}
		} else {
			r.reply(chatID, "No active session. Use /start first.")
		}
	case "/discard":
		if r.stateMachine.IsActive(chatID) {
			r.stateMachine.DiscardSession(chatID)
			r.reply(chatID, "Session discarded.")
			if r.logger != nil {
				r.logger.Info("session discarded", zap.Int64("chat_id", chatID))
			}
		} else {
			r.reply(chatID, "No active session to discard.")
		}
	case "/end":
		if r.stateMachine.IsActive(chatID) {
			session, ok := r.stateMachine.EndSession(chatID)
			if ok {
				if err := r.createNotionPage(ctx, session); err != nil {
					r.reply(chatID, "Failed to save. Check logs for details. Session remains active for retry.")
					if r.logger != nil {
						r.logger.Error("failed to create notion page", zap.Int64("chat_id", chatID), zap.Error(err))
					}
					// Keep session active for retry instead of exiting
					return nil
				}
				r.reply(chatID, "Saved to Notion.")
				if r.logger != nil {
					r.logger.Info("session ended", zap.Int64("chat_id", chatID))
				}
			}
		} else {
			r.reply(chatID, "No active session. Use /start first.")
		}
	default:
		if r.stateMachine.IsActive(chatID) {
			r.collectMessage(update.Message)
		}
	}

	return nil
}

func (r *Runner) collectMessage(msg *tgbotapi.Message) {
	if msg.Text != "" {
		if code := extractCodeBlock(msg); code != nil {
			r.stateMachine.AppendBlock(msg.Chat.ID, code)
			return
		}

		richText := r.mapper.EntitiesToRichText(msg.Text, msg.Entities)
		if len(richText) > 0 {
			r.stateMachine.AppendBlock(msg.Chat.ID, TextBlock{RichText: richText})
		}
	}

	if len(msg.Photo) > 0 {
		photo := msg.Photo[len(msg.Photo)-1]
		r.stateMachine.AppendBlock(msg.Chat.ID, ImageBlock{FileID: photo.FileID, Caption: msg.Caption})
	}
}

func (r *Runner) registerCommands() error {
	commands := []tgbotapi.BotCommand{
		{Command: "start", Description: "Start a new capture session"},
		{Command: "clean", Description: "Clear the current buffer"},
		{Command: "discard", Description: "Discard the current session"},
		{Command: "end", Description: "Save to Notion and end session"},
		{Command: "help", Description: "Show available commands"},
	}

	return r.telegram.SetCommands(commands)
}

func (r *Runner) reply(chatID int64, text string) {
	if err := r.telegram.SendMessage(chatID, text); err != nil {
		if r.logger != nil {
			r.logger.Warn("failed to send reply", zap.Int64("chat_id", chatID), zap.Error(err))
		}
	}
}

func extractCodeBlock(msg *tgbotapi.Message) *CodeBlock {
	if msg.Text == "" {
		return nil
	}

	if len(msg.Entities) != 1 {
		return nil
	}

	entity := msg.Entities[0]
	if entity.Type != "pre" || entity.Offset != 0 || entity.Length != len([]rune(msg.Text)) {
		return nil
	}

	return &CodeBlock{Content: msg.Text, Language: entity.Language}
}

func (r *Runner) createNotionPage(ctx context.Context, session *Session) error {
	loc, err := r.cfg.Title.Location()
	if err != nil {
		return err
	}

	blocks := make([]notionapi.Block, 0, len(session.Blocks))

	for _, block := range session.Blocks {
		switch b := block.(type) {
		case TextBlock:
			blocks = append(blocks, &notionapi.ParagraphBlock{
				BasicBlock: notionapi.BasicBlock{Object: "block", Type: "paragraph"},
				Paragraph:  notionapi.Paragraph{RichText: b.RichText},
			})
		case CodeBlock:
			language := b.Language
			if language == "" {
				language = "plain text"
			}
			blocks = append(blocks, &notionapi.CodeBlock{
				BasicBlock: notionapi.BasicBlock{Object: "block", Type: "code"},
				Code: notionapi.Code{
					RichText: []notionapi.RichText{
						{Type: "text", Text: &notionapi.Text{Content: b.Content}},
					},
					Language: language,
				},
			})
		case ImageBlock:
			if b.FileID == "" {
				continue
			}

			fileURL, err := r.telegram.GetFileURL(b.FileID)
			if err != nil {
				return err
			}

			data, err := r.telegram.DownloadFile(fileURL)
			if err != nil {
				return err
			}

			filename := fmt.Sprintf("%s.jpg", b.FileID)
			rawURL, err := r.github.UploadImage(ctx, data, filename)
			if err != nil {
				return err
			}

			image := &notionapi.ImageBlock{
				BasicBlock: notionapi.BasicBlock{Object: "block", Type: "image"},
				Image: notionapi.Image{
					External: &notionapi.FileObject{URL: rawURL},
				},
			}

			if b.Caption != "" {
				image.Image.Caption = []notionapi.RichText{
					{Type: "text", Text: &notionapi.Text{Content: b.Caption}},
				}
			}

			blocks = append(blocks, image)
		}
	}

	if len(blocks) == 0 {
		return fmt.Errorf("no blocks to save")
	}

	title := r.cfg.Title.FormatTime(loc)
	_, err = r.notion.CreatePage(ctx, r.cfg.Notion.DatabaseID, r.cfg.Notion.TitleProperty, title, blocks)
	return err
}
