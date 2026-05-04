package autolink_pages

import (
	"strings"

	"github.com/emad-elsaid/xlog"
	"github.com/emad-elsaid/xlog/markdown/ast"
	"github.com/emad-elsaid/xlog/markdown/parser"
	"github.com/emad-elsaid/xlog/markdown/text"
	"github.com/emad-elsaid/xlog/markdown/util"
)

type pageLinkParser struct{}

func (*pageLinkParser) Trigger() []byte {
	// ' ' indicates any white spaces and a line head
	return []byte{' ', '*', '_', '~', '('}
}

func (s *pageLinkParser) Parse(parent ast.Node, block text.Reader, pc parser.Context) ast.Node {
	if pc.IsInLinkLabel() {
		return nil
	}

	if err := ensurePagesListInitialized(); err != nil {
		return nil
	}

	line, segment := block.PeekLine()
	if line == nil {
		return nil
	}

	consumes, start, line := advanceIfNeeded(line, segment.Start)
	found, matchLen := findMatchingPage(line)

	if !isValidMatch(found, matchLen, line) {
		block.Advance(consumes)
		return nil
	}

	if consumes != 0 {
		s := segment.WithStop(segment.Start + 1)
		ast.MergeOrAppendTextSegment(parent, s)
	}

	consumes += matchLen
	block.Advance(consumes)

	n := ast.NewTextSegment(text.NewSegment(start, start+matchLen))
	link := &PageLink{
		page: found,
	}
	link.AppendChild(link, n)
	return link
}

func ensurePagesListInitialized() error {
	if autolinkPages != nil {
		return nil
	}
	return UpdatePagesList(nil)
}

func advanceIfNeeded(line []byte, start int) (consumes, newStart int, newLine []byte) {
	c := line[0]
	if c == ' ' || c == '*' || c == '_' || c == '~' || c == '(' {
		return 1, start + 1, line[1:]
	}
	return 0, start, line
}

func findMatchingPage(line []byte) (page xlog.Page, matchLen int) {
	normalizedLine := strings.ToLower(string(line))

	for _, p := range autolinkPages {
		if len(line) < len(p.normalizedName) {
			continue
		}

		if strings.HasPrefix(normalizedLine, p.normalizedName) {
			return p.page, len(p.normalizedName)
		}
	}

	return nil, 0
}

func isValidMatch(found xlog.Page, matchLen int, line []byte) bool {
	if found == nil {
		return false
	}
	// Check if next character is a word character (which would invalidate the match)
	if len(line) > matchLen && util.IsAlphaNumeric(line[matchLen]) {
		return false
	}
	return true
}
