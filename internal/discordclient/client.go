package discordclient

import "github.com/bwmarrin/discordgo"

type Client struct {
	session *discordgo.Session
}

func NewClient(token string) (*Client, error) {
	session, err := discordgo.New("Bot " + token)
	if err != nil {
		return nil, err
	}

	session.Identify.Intents = discordgo.IntentsGuilds | discordgo.IntentsDirectMessages | discordgo.IntentsMessageContent

	return &Client{session: session}, nil
}

func (c *Client) Session() *discordgo.Session {
	return c.session
}

func (c *Client) Open() error {
	return c.session.Open()
}

func (c *Client) Close() error {
	return c.session.Close()
}

func (c *Client) RegisterCommands(commands []*discordgo.ApplicationCommand) error {
	if c.session.State.User == nil {
		user, err := c.session.User("@me")
		if err != nil {
			return err
		}
		c.session.State.User = user
	}

	appID := c.session.State.User.ID
	for _, command := range commands {
		if _, err := c.session.ApplicationCommandCreate(appID, "", command); err != nil {
			return err
		}
	}

	return nil
}
