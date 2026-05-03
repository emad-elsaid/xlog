package extension

import (
	"bytes"
	"regexp"

	"github.com/emad-elsaid/xlog/markdown"
	"github.com/emad-elsaid/xlog/markdown/ast"
	"github.com/emad-elsaid/xlog/markdown/parser"
	"github.com/emad-elsaid/xlog/markdown/text"
	"github.com/emad-elsaid/xlog/markdown/util"
)

var wwwURLRegxp = regexp.MustCompile(`^www\.[-a-zA-Z0-9@:%._\+~#=]{1,256}\.[a-z]+(?:[/#?][-a-zA-Z0-9@:%_\+.~#!?&/=\(\);,'">\^{}\[\]` + "`" + `]*)?`)

var urlRegexp = regexp.MustCompile(`^(?:http|https|ftp)://[-a-zA-Z0-9@:%._\+~#=]{1,256}\.[a-z]+(?::\d+)?(?:[/#?][-a-zA-Z0-9@:%_+.~#$!?&/=\(\);,'">\^{}\[\]` + "`" + `]*)?`)

// An LinkifyConfig struct is a data structure that holds configuration of the
// Linkify extension.
type LinkifyConfig struct {
	AllowedProtocols [][]byte
	URLRegexp        *regexp.Regexp
	WWWRegexp        *regexp.Regexp
	EmailRegexp      *regexp.Regexp
}

const (
	optLinkifyAllowedProtocols parser.OptionName = "LinkifyAllowedProtocols"
	optLinkifyURLRegexp        parser.OptionName = "LinkifyURLRegexp"
	optLinkifyWWWRegexp        parser.OptionName = "LinkifyWWWRegexp"
	optLinkifyEmailRegexp      parser.OptionName = "LinkifyEmailRegexp"
)

// SetOption implements SetOptioner.
func (c *LinkifyConfig) SetOption(name parser.OptionName, value any) {
	switch name {
	case optLinkifyAllowedProtocols:
		c.AllowedProtocols = value.([][]byte)
	case optLinkifyURLRegexp:
		c.URLRegexp = value.(*regexp.Regexp)
	case optLinkifyWWWRegexp:
		c.WWWRegexp = value.(*regexp.Regexp)
	case optLinkifyEmailRegexp:
		c.EmailRegexp = value.(*regexp.Regexp)
	}
}

// A LinkifyOption interface sets options for the LinkifyOption.
type LinkifyOption interface {
	parser.Option
	SetLinkifyOption(*LinkifyConfig)
}

type withLinkifyAllowedProtocols struct {
	value [][]byte
}

func (o *withLinkifyAllowedProtocols) SetParserOption(c *parser.Config) {
	c.Options[optLinkifyAllowedProtocols] = o.value
}

func (o *withLinkifyAllowedProtocols) SetLinkifyOption(p *LinkifyConfig) {
	p.AllowedProtocols = o.value
}

// WithLinkifyAllowedProtocols is a functional option that specify allowed
// protocols in autolinks. Each protocol must end with ':' like
// 'http:' .
func WithLinkifyAllowedProtocols[T []byte | string](value []T) LinkifyOption {
	opt := &withLinkifyAllowedProtocols{}
	for _, v := range value {
		opt.value = append(opt.value, []byte(v))
	}
	return opt
}

type withLinkifyURLRegexp struct {
	value *regexp.Regexp
}

func (o *withLinkifyURLRegexp) SetParserOption(c *parser.Config) {
	c.Options[optLinkifyURLRegexp] = o.value
}

func (o *withLinkifyURLRegexp) SetLinkifyOption(p *LinkifyConfig) {
	p.URLRegexp = o.value
}

// WithLinkifyURLRegexp is a functional option that specify
// a pattern of the URL including a protocol.
func WithLinkifyURLRegexp(value *regexp.Regexp) LinkifyOption {
	return &withLinkifyURLRegexp{
		value: value,
	}
}

type withLinkifyWWWRegexp struct {
	value *regexp.Regexp
}

func (o *withLinkifyWWWRegexp) SetParserOption(c *parser.Config) {
	c.Options[optLinkifyWWWRegexp] = o.value
}

func (o *withLinkifyWWWRegexp) SetLinkifyOption(p *LinkifyConfig) {
	p.WWWRegexp = o.value
}

// WithLinkifyWWWRegexp is a functional option that specify
// a pattern of the URL without a protocol.
// This pattern must start with 'www.' .
func WithLinkifyWWWRegexp(value *regexp.Regexp) LinkifyOption {
	return &withLinkifyWWWRegexp{
		value: value,
	}
}

type withLinkifyEmailRegexp struct {
	value *regexp.Regexp
}

func (o *withLinkifyEmailRegexp) SetParserOption(c *parser.Config) {
	c.Options[optLinkifyEmailRegexp] = o.value
}

func (o *withLinkifyEmailRegexp) SetLinkifyOption(p *LinkifyConfig) {
	p.EmailRegexp = o.value
}

// WithLinkifyEmailRegexp is a functional otpion that specify
// a pattern of the email address.
func WithLinkifyEmailRegexp(value *regexp.Regexp) LinkifyOption {
	return &withLinkifyEmailRegexp{
		value: value,
	}
}

type linkifyParser struct {
	LinkifyConfig
}

// NewLinkifyParser return a new InlineParser can parse
// text that seems like a URL.
func NewLinkifyParser(opts ...LinkifyOption) parser.InlineParser {
	p := &linkifyParser{
		LinkifyConfig: LinkifyConfig{
			AllowedProtocols: nil,
			URLRegexp:        urlRegexp,
			WWWRegexp:        wwwURLRegxp,
		},
	}
	for _, o := range opts {
		o.SetLinkifyOption(&p.LinkifyConfig)
	}
	return p
}

func (s *linkifyParser) Trigger() []byte {
	// ' ' indicates any white spaces and a line head
	return []byte{' ', '*', '_', '~', '('}
}

var (
	protoHTTP  = []byte("http:")
	protoHTTPS = []byte("https:")
	protoFTP   = []byte("ftp:")
	domainWWW  = []byte("www.")
)

// tryMatchURL attempts to match a URL with protocol in the line.
func (s *linkifyParser) tryMatchURL(line []byte) []int {
	if s.AllowedProtocols == nil {
		if bytes.HasPrefix(line, protoHTTP) || bytes.HasPrefix(line, protoHTTPS) || bytes.HasPrefix(line, protoFTP) {
			return s.URLRegexp.FindSubmatchIndex(line)
		}
	} else {
		for _, prefix := range s.AllowedProtocols {
			if bytes.HasPrefix(line, prefix) {
				return s.URLRegexp.FindSubmatchIndex(line)
			}
		}
	}
	return nil
}

// trimURLTrailingPunctuation removes trailing punctuation from URL matches.
func (s *linkifyParser) trimURLTrailingPunctuation(line []byte, m []int) {
	lastChar := line[m[1]-1]
	switch lastChar {
	case '.':
		m[1]--
	case ')':
		closing := 0
		for i := m[1] - 1; i >= m[0]; i-- {
			switch line[i] {
			case ')':
				closing++
			case '(':
				closing--
			}
		}
		if closing > 0 {
			m[1] -= closing
		}
	case ';':
		i := m[1] - 2
		for ; i >= m[0]; i-- {
			if util.IsAlphaNumeric(line[i]) {
				continue
			}
			break
		}
		if i != m[1]-2 {
			if line[i] == '&' {
				m[1] -= m[1] - i
			}
		}
	}
}

// tryMatchEmail attempts to match an email address in the line.
func (s *linkifyParser) tryMatchEmail(line []byte) []int {
	if len(line) > 0 && util.IsPunct(line[0]) {
		return nil
	}
	stop := -1
	if s.EmailRegexp == nil {
		stop = util.FindEmailIndex(line)
	} else {
		m := s.EmailRegexp.FindSubmatchIndex(line)
		if m != nil && m[0] == 0 {
			stop = m[1]
		}
	}
	if stop < 0 {
		return nil
	}
	at := bytes.IndexByte(line, '@')
	m := []int{0, stop, at, stop - 1}
	if m == nil || bytes.IndexByte(line[m[2]:m[3]], '.') < 0 {
		return nil
	}
	lastChar := line[m[1]-1]
	if lastChar == '.' {
		m[1]--
	}
	if m[1] < len(line) {
		nextChar := line[m[1]]
		if nextChar == '-' || nextChar == '_' {
			return nil
		}
	}
	return m
}

// trimTrailingPunctuation removes trailing punctuation from the matched text.
func trimTrailingPunctuation(line []byte, endIdx int) int {
	i := endIdx - 1
	for ; i > 0; i-- {
		c := line[i]
		switch c {
		case '?', '!', '.', ',', ':', '*', '_', '~':
		default:
			return i + 1
		}
	}
	return i + 1
}

func (s *linkifyParser) Parse(parent ast.Node, block text.Reader, pc parser.Context) ast.Node {
	if pc.IsInLinkLabel() {
		return nil
	}
	line, segment := block.PeekLine()
	consumes := 0
	start := segment.Start
	c := line[0]
	// Advance if current position is not a line head.
	if c == ' ' || c == '*' || c == '_' || c == '~' || c == '(' {
		consumes++
		start++
		line = line[1:]
	}

	var m []int
	var protocol []byte
	typ := ast.AutoLinkURL

	// Try to match URL with protocol.
	m = s.tryMatchURL(line)
	if m == nil && bytes.HasPrefix(line, domainWWW) {
		m = s.WWWRegexp.FindSubmatchIndex(line)
		protocol = []byte("http")
	}
	if m != nil && m[0] != 0 {
		m = nil
	}
	if m != nil && m[0] == 0 {
		s.trimURLTrailingPunctuation(line, m)
	}

	// Try to match email if no URL was found.
	if m == nil {
		typ = ast.AutoLinkEmail
		m = s.tryMatchEmail(line)
	}

	if m == nil {
		return nil
	}

	if consumes != 0 {
		s := segment.WithStop(segment.Start + 1)
		ast.MergeOrAppendTextSegment(parent, s)
	}

	i := trimTrailingPunctuation(line, m[1])
	consumes += i
	block.Advance(consumes)
	n := ast.NewTextSegment(text.NewSegment(start, start+i))
	link := ast.NewAutoLink(typ, n)
	link.Protocol = protocol
	return link
}

func (s *linkifyParser) CloseBlock(parent ast.Node, pc parser.Context) {
	// nothing to do
}

type linkify struct {
	options []LinkifyOption
}

// Linkify is an extension that allow you to parse text that seems like a URL.
var Linkify = &linkify{}

// NewLinkify creates a new [markdown.Extender] that
// allow you to parse text that seems like a URL.
func NewLinkify(opts ...LinkifyOption) markdown.Extender {
	return &linkify{
		options: opts,
	}
}

func (e *linkify) Extend(m markdown.Markdown) {
	m.Parser().AddOptions(
		parser.WithInlineParsers(
			util.Prioritized(NewLinkifyParser(e.options...), 999),
		),
	)
}
