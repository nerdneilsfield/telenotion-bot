package discordclient

import (
	"regexp"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/jomei/notionapi"
)

var (
	codeRe   = regexp.MustCompile("`([^`]+)`")
	boldRe   = regexp.MustCompile(`\*\*([^*]+)\*\*`)
	italicRe = regexp.MustCompile(`\*([^*]+)\*|_([^_]+)_`)
	linkRe   = regexp.MustCompile(`\[([^\]]+)\]\(([^)]+)\)`)
	emojiRe  = regexp.MustCompile(`<a?:([a-zA-Z0-9_]+):\d+>`)
)

func NormalizeContent(content string, mentions []*discordgo.User) string {
	for _, user := range mentions {
		if user == nil {
			continue
		}
		content = strings.ReplaceAll(content, "<@"+user.ID+">", "@"+user.Username)
		content = strings.ReplaceAll(content, "<@!"+user.ID+">", "@"+user.Username)
	}

	return emojiRe.ReplaceAllString(content, ":$1:")
}

func ContentToRichText(content string) []notionapi.RichText {
	if content == "" {
		return nil
	}

	remaining := content
	result := make([]notionapi.RichText, 0, 4)

	for len(remaining) > 0 {
		match := nextMatch(remaining)
		if match.kind == "" {
			result = append(result, plainRichText(remaining))
			break
		}
		if match.start > 0 {
			result = append(result, plainRichText(remaining[:match.start]))
		}

		result = append(result, match.toRichText())
		remaining = remaining[match.end:]
	}

	return result
}

type tokenMatch struct {
	kind  string
	start int
	end   int
	text  string
	url   string
}

func nextMatch(content string) tokenMatch {
	matches := []struct {
		kind string
		re   *regexp.Regexp
	}{
		{kind: "code", re: codeRe},
		{kind: "bold", re: boldRe},
		{kind: "italic", re: italicRe},
		{kind: "link", re: linkRe},
	}

	best := tokenMatch{}
	best.start = -1

	for _, candidate := range matches {
		idx := candidate.re.FindStringSubmatchIndex(content)
		if idx == nil {
			continue
		}
		start, end := idx[0], idx[1]
		if best.start != -1 && start >= best.start {
			continue
		}
		best = tokenMatch{kind: candidate.kind, start: start, end: end}

		switch candidate.kind {
		case "code":
			best.text = content[idx[2]:idx[3]]
		case "bold":
			best.text = content[idx[2]:idx[3]]
		case "italic":
			if idx[2] != -1 {
				best.text = content[idx[2]:idx[3]]
			} else {
				best.text = content[idx[4]:idx[5]]
			}
		case "link":
			best.text = content[idx[2]:idx[3]]
			best.url = content[idx[4]:idx[5]]
		}
	}

	if best.start == -1 {
		return tokenMatch{}
	}

	return best
}

func (t tokenMatch) toRichText() notionapi.RichText {
	richText := notionapi.RichText{
		Type: "text",
		Text: &notionapi.Text{Content: t.text},
	}

	annotations := notionapi.Annotations{}
	switch t.kind {
	case "code":
		annotations.Code = true
	case "bold":
		annotations.Bold = true
	case "italic":
		annotations.Italic = true
	case "link":
		richText.Text.Link = &notionapi.Link{Url: t.url}
	}

	if annotations.Bold || annotations.Italic || annotations.Code {
		richText.Annotations = &annotations
	}

	return richText
}

func plainRichText(text string) notionapi.RichText {
	return notionapi.RichText{
		Type: "text",
		Text: &notionapi.Text{Content: text},
	}
}
