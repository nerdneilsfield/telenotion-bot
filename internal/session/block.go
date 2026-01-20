package session

import "github.com/jomei/notionapi"

type Block interface {
	Kind() string
}

type TextBlock struct {
	RichText []notionapi.RichText
}

type CodeBlock struct {
	Content  string
	Language string
}

type ImageBlock struct {
	FileID   string
	FileURL  string
	Filename string
	Caption  string
}

func (t TextBlock) Kind() string {
	return "text"
}

func (c CodeBlock) Kind() string {
	return "code"
}

func (i ImageBlock) Kind() string {
	return "image"
}
