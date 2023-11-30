package parsing

import (
	"bytes"

	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/text"
)

type TailwindTransformer struct {
	headlineMapping map[int]string
}

func getHeadlineMappings() map[int]string {
	return map[int]string{
		1: "text-4xl my-4",
		2: "text-3xl font-semibold my-2",
		3: "text-2xl font-semibold",
		4: "text-xl font-semibold",
	}
}

func NewTailwindTransformer() *TailwindTransformer {
	return &TailwindTransformer{
		headlineMapping: getHeadlineMappings(),
	}
}

func (t *TailwindTransformer) Transform(node *ast.Document, reader text.Reader, pc parser.Context) {
	ast.Walk(node, func(n ast.Node, entering bool) (ast.WalkStatus, error) {
		switch nt := n.(type) {
		case *ast.Heading:
			var buf bytes.Buffer
			buf.WriteString(t.headlineMapping[nt.Level])
			nt.SetAttributeString("class", buf.Bytes())
		case *ast.Paragraph:
			var buf bytes.Buffer
			buf.WriteString("text-lg leading-relaxed mb-4")
			nt.SetAttributeString("class", buf.Bytes())
		case *ast.Link:
			var buf bytes.Buffer
			buf.WriteString("underline text-primary")
			nt.SetAttributeString("class", buf.Bytes())
		case *ast.List:
			var buf bytes.Buffer
			buf.WriteString("list-disc list-inside mb-4 leading-relaxed text-lg")
			nt.SetAttributeString("class", buf.Bytes())
		}
		return ast.WalkContinue, nil
	})
}
