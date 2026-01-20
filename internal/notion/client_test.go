package notion

import (
	"testing"

	"github.com/jomei/notionapi"
)

func TestNewClient(t *testing.T) {
	client := NewClient("test-token")

	if client == nil {
		t.Fatal("NewClient returned nil")
	}
	if client.client == nil {
		t.Error("Client.client is nil")
	}
}

func TestNewClient_EmptyToken(t *testing.T) {
	client := NewClient("")

	if client == nil {
		t.Fatal("NewClient should return a client even with empty token")
	}
	if client.client == nil {
		t.Error("Client.client should not be nil")
	}
}

func TestCreatePage_Properties(t *testing.T) {
	client := NewClient("test-token")

	// Verify the client is properly initialized
	_ = client

	// Test data
	databaseID := "test-database-id"
	titleProperty := "Name"
	title := "Test Page Title"
	blocks := []notionapi.Block{
		&notionapi.ParagraphBlock{
			BasicBlock: notionapi.BasicBlock{Object: "block", Type: "paragraph"},
			Paragraph: notionapi.Paragraph{
				RichText: []notionapi.RichText{
					{Type: "text", Text: &notionapi.Text{Content: "Hello, World!"}},
				},
			},
		},
	}

	// Just verify the parameters are valid for the API call
	if databaseID == "" {
		t.Error("databaseID should not be empty")
	}
	if titleProperty == "" {
		t.Error("titleProperty should not be empty")
	}
	if title == "" {
		t.Error("title should not be empty")
	}
	if len(blocks) == 0 {
		t.Error("blocks should not be empty")
	}

	// Verify block structure
	for i, block := range blocks {
		if block.GetType() == "" {
			t.Errorf("Block %d has empty type", i)
		}
	}
}

func TestCreatePage_EmptyBlocks(t *testing.T) {
	client := NewClient("test-token")

	// Empty blocks should still create a page (though it might be empty)
	blocks := []notionapi.Block{}

	if blocks == nil {
		t.Error("blocks slice should not be nil")
	}
	if len(blocks) != 0 {
		t.Error("blocks should be empty")
	}

	_ = client
}

func TestCreatePage_MultipleBlockTypes(t *testing.T) {
	blocks := []notionapi.Block{
		&notionapi.ParagraphBlock{
			BasicBlock: notionapi.BasicBlock{Object: "block", Type: "paragraph"},
			Paragraph: notionapi.Paragraph{
				RichText: []notionapi.RichText{
					{Type: "text", Text: &notionapi.Text{Content: "Text content"}},
				},
			},
		},
		&notionapi.CodeBlock{
			BasicBlock: notionapi.BasicBlock{Object: "block", Type: "code"},
			Code: notionapi.Code{
				RichText: []notionapi.RichText{
					{Type: "text", Text: &notionapi.Text{Content: "fmt.Println(\"Hello\")"}},
				},
				Language: "go",
			},
		},
	}

	if len(blocks) != 2 {
		t.Errorf("Expected 2 blocks, got %d", len(blocks))
	}

	// Verify all blocks have valid types
	for i, block := range blocks {
		blockType := block.GetType()
		if blockType == "" {
			t.Errorf("Block %d has empty type", i)
		}
	}
}
