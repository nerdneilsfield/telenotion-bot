package notion

import (
	"context"

	"github.com/jomei/notionapi"
)

type Client struct {
	client *notionapi.Client
}

func NewClient(token string) *Client {
	return &Client{client: notionapi.NewClient(notionapi.Token(token))}
}

func (c *Client) CreatePage(ctx context.Context, databaseID string, titleProperty string, title string, children []notionapi.Block) (string, error) {
	page, err := c.client.Page.Create(ctx, &notionapi.PageCreateRequest{
		Parent: notionapi.Parent{Type: notionapi.ParentTypeDatabaseID, DatabaseID: notionapi.DatabaseID(databaseID)},
		Properties: notionapi.Properties{
			titleProperty: notionapi.TitleProperty{
				Title: []notionapi.RichText{{Text: &notionapi.Text{Content: title}}},
			},
		},
		Children: children,
	})
	if err != nil {
		return "", err
	}

	return string(page.ID), nil
}
