package main

import (
	"encoding/base64"
	"io"
	"mime"
	"os"
	"path/filepath"
	"regexp"

	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
)

func encodeImage(filename string) (string, error) {
	// Open the image file
	file, err := os.Open(filename)
	if err != nil {
		return "", err
	}
	defer file.Close()
	fileContent, err := io.ReadAll(file)
	if err != nil {
		return "", err
	}
	// encode content to base64 URL encoded, with the mime type
	mimeType := mime.TypeByExtension(filepath.Ext(filename))
	encoded := base64.StdEncoding.Strict().EncodeToString(fileContent)
	return "data:" + mimeType + ";base64," + encoded, nil
}

func MDtoHTML(source string) []byte {
	// extensions := parser.NoIntraEmphasis | parser.Tables | parser.FencedCode | parser.Autolink | parser.Strikethrough | parser.SpaceHeadings | parser.HeadingIDs | parser.BackslashLineBreak | parser.DefinitionLists
	extensions := parser.CommonExtensions
	p := parser.NewWithExtensions(extensions)
	doc := p.Parse([]byte(source))

	htmlFlags := html.CommonFlags | html.UseXHTML
	opts := html.RendererOptions{Flags: htmlFlags}
	renderer := html.NewRenderer(opts)

	return markdown.Render(doc, renderer)
}

var (
	matchLatexBlock   = regexp.MustCompile(`(?ms)\\\[(.+?)\\\]`)
	matchLatexInline  = regexp.MustCompile(`(?ms)\\\((.+?)\\\)`)
	latexReplacements = map[*regexp.Regexp]string{
		matchLatexBlock:  "$$$$$1$$$$",
		matchLatexInline: "$$$1$",
	}
)

// fixKatex fixes the markdown katex by replacing \[\sand \s\]to $$, and \(\s and \s\) to $.
func fixKatex(chunk string) string {
	for re, replacement := range latexReplacements {
		chunk = re.ReplaceAllString(chunk, replacement)
	}
	return chunk
}
