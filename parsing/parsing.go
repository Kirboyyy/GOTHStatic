package parsing

import (
	"blog/model"
	"bytes"
	"fmt"

	chromahtml "github.com/alecthomas/chroma/v2/formatters/html"
	"github.com/yuin/goldmark"

	highlighting "github.com/yuin/goldmark-highlighting/v2"
	meta "github.com/yuin/goldmark-meta"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer/html"
	"github.com/yuin/goldmark/util"
)

func ConvertMarkdownToHTML(source string) (string, model.Post) {
	markdown := goldmark.New(
		goldmark.WithExtensions(
			extension.GFM,
			meta.Meta,
			highlighting.NewHighlighting(
				highlighting.WithStyle("monokai"),
				highlighting.WithFormatOptions(
					chromahtml.WithLineNumbers(true),
				),
			),
		),
		goldmark.WithParserOptions(
			parser.WithAutoHeadingID(),
			parser.WithASTTransformers(
				util.PrioritizedValue{Value: NewTailwindTransformer(), Priority: 0},
			)),
		goldmark.WithRendererOptions(
			html.WithHardWraps(),
		),
	)

	var buf bytes.Buffer
	context := parser.NewContext()
	if err := markdown.Convert([]byte(source), &buf, parser.WithContext(context)); err != nil {
		panic(err)
	}
	metaData := meta.Get(context)

	id := -1337 // negative value for insert behaviour
	if metaData["id"] != nil {
		id = metaData["id"].(int)
	}

	tagsInterface := metaData["tags"].([]interface{})
	var tags []string
	for _, tag := range tagsInterface {
		tags = append(tags, tag.(string))
	}

	return buf.String(), model.Post{
		ID:          id,
		Title:       fmt.Sprintf("%v", metaData["title"]),
		Subtitle:    fmt.Sprintf("%v", metaData["subtitle"]),
		Description: fmt.Sprintf("%v", metaData["description"]),
		Image:       fmt.Sprintf("%v", metaData["image"]),
		Slug:        fmt.Sprintf("%v", metaData["slug"]),
		Tags:        tags,
	}
}
