package session

import "github.com/jomei/notionapi"

const (
	notionRichTextLimit      = 2000
	notionRichTextBlockLimit = 100
)

func splitRichTextEntries(richTexts []notionapi.RichText) []notionapi.RichText {
	result := make([]notionapi.RichText, 0, len(richTexts))
	for _, rt := range richTexts {
		if rt.Text == nil || rt.Text.Content == "" {
			result = append(result, rt)
			continue
		}
		content := rt.Text.Content
		for len(content) > 0 {
			chunk, rest := splitRunes(content, notionRichTextLimit)
			copyRt := rt
			copyRt.Text = &notionapi.Text{Content: chunk, Link: rt.Text.Link}
			result = append(result, copyRt)
			content = rest
		}
	}
	return result
}

func splitRunes(text string, limit int) (string, string) {
	runes := []rune(text)
	if len(runes) <= limit {
		return text, ""
	}
	return string(runes[:limit]), string(runes[limit:])
}

func chunkRichText(richTexts []notionapi.RichText, limit int) [][]notionapi.RichText {
	if len(richTexts) == 0 {
		return nil
	}
	if limit <= 0 || len(richTexts) <= limit {
		return [][]notionapi.RichText{richTexts}
	}

	chunks := make([][]notionapi.RichText, 0, (len(richTexts)+limit-1)/limit)
	for len(richTexts) > 0 {
		if len(richTexts) <= limit {
			chunks = append(chunks, richTexts)
			break
		}
		chunks = append(chunks, richTexts[:limit])
		richTexts = richTexts[limit:]
	}
	return chunks
}
