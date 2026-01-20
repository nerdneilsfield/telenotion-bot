package tgclient

import (
	"io"
	"net/http"

	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Client struct {
	bot    *tgbotapi.BotAPI
	client *http.Client
}

func NewClient(token string) (*Client, error) {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, err
	}

	return &Client{
		bot:    bot,
		client: http.DefaultClient,
	}, nil
}

func (c *Client) GetUpdates(offset int, timeout int) ([]tgbotapi.Update, error) {
	u := tgbotapi.NewUpdate(timeout)
	u.Offset = offset

	updates, err := c.bot.GetUpdates(u)
	if err != nil {
		return nil, err
	}

	return updates, nil
}

func (c *Client) IsAllowedChat(allowed []int64, chatID int64) bool {
	for _, id := range allowed {
		if id == chatID {
			return true
		}
	}
	return false
}

func (c *Client) GetFileURL(fileID string) (string, error) {
	file, err := c.bot.GetFile(tgbotapi.FileConfig{FileID: fileID})
	if err != nil {
		return "", err
	}

	return file.Link(c.bot.Token), nil
}

func (c *Client) DownloadFile(url string) ([]byte, error) {
	resp, err := c.client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return io.ReadAll(resp.Body)
}
