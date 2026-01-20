package tgclient

import (
	"testing"

	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func TestEntitiesToRichText(t *testing.T) {
	mapper := NewMapper()

	text := "hello world"
	entities := []tgbotapi.MessageEntity{
		{Type: "bold", Offset: 0, Length: 5},
		{Type: "text_link", Offset: 6, Length: 5, URL: "https://example.com"},
	}

	result := mapper.EntitiesToRichText(text, entities)
	if len(result) != 3 {
		t.Fatalf("expected 3 segments, got %d", len(result))
	}

	if result[0].Text.Content != "hello" {
		t.Fatalf("expected first segment to be 'hello'")
	}
	if result[0].Annotations == nil || !result[0].Annotations.Bold {
		t.Fatalf("expected first segment to be bold")
	}

	if result[1].Text.Content != " " {
		t.Fatalf("expected middle segment to be a space")
	}

	if result[2].Text.Content != "world" {
		t.Fatalf("expected third segment to be 'world'")
	}
	if result[2].Text.Link == nil || result[2].Text.Link.Url != "https://example.com" {
		t.Fatalf("expected link on third segment")
	}
}

func TestEntitiesToRichTextPlain(t *testing.T) {
	mapper := NewMapper()

	text := "plain"
	result := mapper.EntitiesToRichText(text, nil)
	if len(result) != 1 {
		t.Fatalf("expected 1 segment, got %d", len(result))
	}

	if result[0].Text.Content != "plain" {
		t.Fatalf("expected plain text")
	}
}
