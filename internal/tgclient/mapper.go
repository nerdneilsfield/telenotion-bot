package tgclient

import (
	"sort"

	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/jomei/notionapi"
)

type Mapper struct{}

func NewMapper() *Mapper {
	return &Mapper{}
}

func (m *Mapper) EntitiesToRichText(text string, entities []tgbotapi.MessageEntity) []notionapi.RichText {
	if text == "" {
		return nil
	}

	if len(entities) == 0 {
		return []notionapi.RichText{plainRichText(text)}
	}

	sorted := append([]tgbotapi.MessageEntity(nil), entities...)
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].Offset < sorted[j].Offset
	})

	runes := []rune(text)
	cursor := 0
	result := make([]notionapi.RichText, 0, len(sorted)+1)

	for _, entity := range sorted {
		start := entity.Offset
		end := entity.Offset + entity.Length

		if start < 0 || start >= len(runes) {
			continue
		}
		if end > len(runes) {
			end = len(runes)
		}

		if start > cursor {
			result = append(result, plainRichText(string(runes[cursor:start])))
		}

		segment := string(runes[start:end])
		result = append(result, entityToRichText(segment, entity))
		cursor = end
	}

	if cursor < len(runes) {
		result = append(result, plainRichText(string(runes[cursor:])))
	}

	return result
}

func plainRichText(text string) notionapi.RichText {
	return notionapi.RichText{
		Type: "text",
		Text: &notionapi.Text{Content: text},
	}
}

func entityToRichText(text string, entity tgbotapi.MessageEntity) notionapi.RichText {
	richText := notionapi.RichText{
		Type: "text",
		Text: &notionapi.Text{Content: text},
	}

	annotations := notionapi.Annotations{}

	switch entity.Type {
	case "bold":
		annotations.Bold = true
	case "italic":
		annotations.Italic = true
	case "code":
		annotations.Code = true
	case "text_link":
		if entity.URL != "" {
			richText.Text.Link = &notionapi.Link{Url: entity.URL}
		}
	}

	if annotations.Bold || annotations.Italic || annotations.Code {
		richText.Annotations = &annotations
	}

	return richText
}
