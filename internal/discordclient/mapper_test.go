package discordclient

import (
	"testing"

	"github.com/bwmarrin/discordgo"
)

func TestNormalizeContent(t *testing.T) {
	mentions := []*discordgo.User{{ID: "123", Username: "alice"}}
	content := "Hello <@123> <:wave:999>"
	result := NormalizeContent(content, mentions)

	if result != "Hello @alice :wave:" {
		t.Fatalf("NormalizeContent() = %q", result)
	}
}

func TestContentToRichTextBold(t *testing.T) {
	result := ContentToRichText("**bold**")
	if len(result) != 1 {
		t.Fatalf("expected 1 segment, got %d", len(result))
	}
	if result[0].Annotations == nil || !result[0].Annotations.Bold {
		t.Fatalf("expected bold annotation")
	}
}

func TestContentToRichTextLink(t *testing.T) {
	result := ContentToRichText("[link](https://example.com)")
	if len(result) != 1 {
		t.Fatalf("expected 1 segment, got %d", len(result))
	}
	if result[0].Text == nil || result[0].Text.Link == nil {
		t.Fatalf("expected link")
	}
	if result[0].Text.Link.Url != "https://example.com" {
		t.Fatalf("unexpected link url: %q", result[0].Text.Link.Url)
	}
}
