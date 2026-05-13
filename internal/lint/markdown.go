package lint

import (
	"path/filepath"
	"strings"

	"github.com/yuin/goldmark"
	gast "github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/extension"
	extast "github.com/yuin/goldmark/extension/ast"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/text"
)

type MarkdownDocument struct {
	Headings []Heading
	Tables   []Table
}

type Heading struct {
	Level     int
	Text      string
	StartLine int
}

type Table struct {
	Section   string
	StartLine int
	Headers   []string
	Rows      []TableRow
}

type TableRow struct {
	Cells []TableCell
}

type TableCell struct {
	Text  string
	Links []Link
}

type Link struct {
	Text        string
	Destination string
}

func ParseMarkdown(source []byte) (MarkdownDocument, error) {
	markdown := goldmark.New(goldmark.WithExtensions(extension.Table))
	reader := text.NewReader(source)
	root := markdown.Parser().Parse(reader, parser.WithContext(parser.NewContext()))
	doc := MarkdownDocument{}
	lineIndex := newLineIndex(source)
	currentSection := ""
	for child := root.FirstChild(); child != nil; child = child.NextSibling() {
		switch node := child.(type) {
		case *gast.Heading:
			heading := Heading{Level: node.Level, Text: inlineText(node, source), StartLine: lineIndex.lineForNode(node)}
			doc.Headings = append(doc.Headings, heading)
			currentSection = heading.Text
		case *extast.Table:
			doc.Tables = append(doc.Tables, parseTable(node, currentSection, source, lineIndex))
		}
	}
	return doc, nil
}

func (d MarkdownDocument) HasHeading(text string) bool {
	for _, heading := range d.Headings {
		if heading.Text == text {
			return true
		}
	}
	return false
}

func (d MarkdownDocument) TableBySection(section string) *Table {
	for i := range d.Tables {
		if d.Tables[i].Section == section {
			return &d.Tables[i]
		}
	}
	return nil
}

func parseTable(node *extast.Table, section string, source []byte, index lineIndex) Table {
	table := Table{Section: section, StartLine: index.lineForNodeDeep(node)}
	for child := node.FirstChild(); child != nil; child = child.NextSibling() {
		switch row := child.(type) {
		case *extast.TableHeader:
			table.Headers = cellTexts(cellsFromRow(row, source))
		case *extast.TableRow:
			table.Rows = append(table.Rows, TableRow{Cells: cellsFromRow(row, source)})
		}
	}
	return table
}

func cellsFromRow(row gast.Node, source []byte) []TableCell {
	var cells []TableCell
	for child := row.FirstChild(); child != nil; child = child.NextSibling() {
		if _, ok := child.(*extast.TableCell); !ok {
			continue
		}
		cells = append(cells, TableCell{Text: inlineText(child, source), Links: linksInNode(child, source)})
	}
	return cells
}

func cellTexts(cells []TableCell) []string {
	out := make([]string, 0, len(cells))
	for _, cell := range cells {
		out = append(out, cell.Text)
	}
	return out
}

func inlineText(node gast.Node, source []byte) string {
	var parts []string
	_ = gast.Walk(node, func(child gast.Node, entering bool) (gast.WalkStatus, error) {
		if !entering || child == node {
			return gast.WalkContinue, nil
		}
		switch typed := child.(type) {
		case *gast.Text:
			parts = append(parts, string(typed.Value(source)))
		case *gast.String:
			parts = append(parts, string(typed.Value))
		}
		return gast.WalkContinue, nil
	})
	return strings.Join(strings.Fields(strings.Join(parts, " ")), " ")
}

func linksInNode(node gast.Node, source []byte) []Link {
	var links []Link
	_ = gast.Walk(node, func(child gast.Node, entering bool) (gast.WalkStatus, error) {
		if !entering {
			return gast.WalkContinue, nil
		}
		link, ok := child.(*gast.Link)
		if !ok {
			return gast.WalkContinue, nil
		}
		links = append(links, Link{Text: inlineText(link, source), Destination: normalizeMarkdownLink(string(link.Destination))})
		return gast.WalkContinue, nil
	})
	return links
}

func normalizeMarkdownLink(destination string) string {
	destination = filepath.ToSlash(strings.TrimSpace(destination))
	destination = strings.TrimPrefix(destination, "./")
	return destination
}

type lineIndex []int

func newLineIndex(source []byte) lineIndex {
	starts := []int{0}
	for offset, b := range source {
		if b == '\n' && offset+1 < len(source) {
			starts = append(starts, offset+1)
		}
	}
	return starts
}

func (idx lineIndex) lineForNode(node gast.Node) int {
	lines := node.Lines()
	if lines == nil || lines.Len() == 0 {
		return 0
	}
	return idx.lineForOffset(lines.At(0).Start)
}

func (idx lineIndex) lineForNodeDeep(node gast.Node) int {
	if line := idx.lineForNode(node); line != 0 {
		return line
	}
	for child := node.FirstChild(); child != nil; child = child.NextSibling() {
		if line := idx.lineForNodeDeep(child); line != 0 {
			return line
		}
	}
	return 0
}

func (idx lineIndex) lineForOffset(offset int) int {
	line := 1
	for _, start := range idx {
		if start > offset {
			break
		}
		line++
	}
	return line - 1
}
