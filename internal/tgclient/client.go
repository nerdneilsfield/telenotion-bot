package tgclient

import (
	"io"
	"net"
	"net/http"
	"time"

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
		client: newHTTPClient(),
	}, nil
}

func newHTTPClient() *http.Client {
	return &http.Client{
		Timeout: 30 * time.Second,
		Transport: &http.Transport{
			DialContext: (&net.Dialer{
				Timeout:   5 * time.Second,
				KeepAlive: 30 * time.Second,
			}).DialContext,
			TLSHandshakeTimeout:   10 * time.Second,
			ResponseHeaderTimeout: 10 * time.Second,
			ExpectContinueTimeout: 1 * time.Second,
			MaxIdleConns:          100,
			IdleConnTimeout:       90 * time.Second,
		},
	}
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

func (c *Client) SendMessage(chatID int64, text string) error {
	msg := tgbotapi.NewMessage(chatID, text)
	_, err := c.bot.Send(msg)
	return err
}

func (c *Client) SetCommands(commands []tgbotapi.BotCommand) error {
	config := tgbotapi.NewSetMyCommands(commands...)
	_, err := c.bot.Request(config)
	return err
}
