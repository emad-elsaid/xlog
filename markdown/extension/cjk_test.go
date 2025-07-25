package extension

import (
	"testing"

	"github.com/emad-elsaid/xlog/markdown"
	"github.com/emad-elsaid/xlog/markdown/renderer/html"
	"github.com/emad-elsaid/xlog/markdown/testutil"
)

func TestEscapedSpace(t *testing.T) {
	md := markdown.New(markdown.WithRendererOptions(
		html.WithXHTML(),
		html.WithUnsafe(),
	))
	no := 1
	testutil.DoTestCase(
		md,
		testutil.MarkdownTestCase{
			No:          no,
			Description: "Without spaces around an emphasis started with east asian punctuations, it is not interpreted as an emphasis(as defined in CommonMark spec)",
			Markdown:    "太郎は**「こんにちわ」**と言った\nんです",
			Expected:    "<p>太郎は**「こんにちわ」**と言った\nんです</p>",
		},
		t,
	)

	no = 2
	testutil.DoTestCase(
		md,
		testutil.MarkdownTestCase{
			No:          no,
			Description: "With spaces around an emphasis started with east asian punctuations, it is interpreted as an emphasis(but remains unnecessary spaces)",
			Markdown:    "太郎は **「こんにちわ」** と言った\nんです",
			Expected:    "<p>太郎は <strong>「こんにちわ」</strong> と言った\nんです</p>",
		},
		t,
	)

	// Enables EscapedSpace
	md = markdown.New(markdown.WithRendererOptions(
		html.WithXHTML(),
		html.WithUnsafe(),
	),
		markdown.WithExtensions(NewCJK(WithEscapedSpace())),
	)

	no = 3
	testutil.DoTestCase(
		md,
		testutil.MarkdownTestCase{
			No:          no,
			Description: "With spaces around an emphasis started with east asian punctuations,it is interpreted as an emphasis",
			Markdown:    "太郎は\\ **「こんにちわ」**\\ と言った\nんです",
			Expected:    "<p>太郎は<strong>「こんにちわ」</strong>と言った\nんです</p>",
		},
		t,
	)

	// ' ' triggers Linkify extension inline parser.
	// Escaped spaces should not trigger the inline parser.

	md = markdown.New(markdown.WithRendererOptions(
		html.WithXHTML(),
		html.WithUnsafe(),
	),
		markdown.WithExtensions(
			NewCJK(WithEscapedSpace()),
			Linkify,
		),
	)

	no = 4
	testutil.DoTestCase(
		md,
		testutil.MarkdownTestCase{
			No:          no,
			Description: "Escaped space and linkfy extension",
			Markdown:    "太郎は\\ **「こんにちわ」**\\ と言った\nんです",
			Expected:    "<p>太郎は<strong>「こんにちわ」</strong>と言った\nんです</p>",
		},
		t,
	)
}

func TestEastAsianLineBreaks(t *testing.T) {
	md := markdown.New(markdown.WithRendererOptions(
		html.WithXHTML(),
		html.WithUnsafe(),
	))
	no := 1
	testutil.DoTestCase(
		md,
		testutil.MarkdownTestCase{
			No:          no,
			Description: "Soft line breaks are rendered as a newline, so some asian users will see it as an unnecessary space",
			Markdown:    "太郎は\\ **「こんにちわ」**\\ と言った\nんです",
			Expected:    "<p>太郎は\\ <strong>「こんにちわ」</strong>\\ と言った\nんです</p>",
		},
		t,
	)

	// Enables EastAsianLineBreaks

	md = markdown.New(markdown.WithRendererOptions(
		html.WithXHTML(),
		html.WithUnsafe(),
	),
		markdown.WithExtensions(NewCJK(WithEastAsianLineBreaks())),
	)

	no = 2
	testutil.DoTestCase(
		md,
		testutil.MarkdownTestCase{
			No:          no,
			Description: "Soft line breaks between east asian wide characters are ignored",
			Markdown:    "太郎は\\ **「こんにちわ」**\\ と言った\nんです",
			Expected:    "<p>太郎は\\ <strong>「こんにちわ」</strong>\\ と言ったんです</p>",
		},
		t,
	)

	no = 3
	testutil.DoTestCase(
		md,
		testutil.MarkdownTestCase{
			No:          no,
			Description: "Soft line breaks between western characters are rendered as a newline",
			Markdown:    "太郎は\\ **「こんにちわ」**\\ と言ったa\nbんです",
			Expected:    "<p>太郎は\\ <strong>「こんにちわ」</strong>\\ と言ったa\nbんです</p>",
		},
		t,
	)

	no = 4
	testutil.DoTestCase(
		md,
		testutil.MarkdownTestCase{
			No:          no,
			Description: "Soft line breaks between a western character and an east asian wide character are rendered as a newline",
			Markdown:    "太郎は\\ **「こんにちわ」**\\ と言ったa\nんです",
			Expected:    "<p>太郎は\\ <strong>「こんにちわ」</strong>\\ と言ったa\nんです</p>",
		},
		t,
	)

	no = 5
	testutil.DoTestCase(
		md,
		testutil.MarkdownTestCase{
			No:          no,
			Description: "Soft line breaks between an east asian wide character and a western character are rendered as a newline",
			Markdown:    "太郎は\\ **「こんにちわ」**\\ と言った\nbんです",
			Expected:    "<p>太郎は\\ <strong>「こんにちわ」</strong>\\ と言った\nbんです</p>",
		},
		t,
	)

	// WithHardWraps take precedence over WithEastAsianLineBreaks
	md = markdown.New(markdown.WithRendererOptions(
		html.WithHardWraps(),
		html.WithXHTML(),
		html.WithUnsafe(),
	),
		markdown.WithExtensions(NewCJK(WithEastAsianLineBreaks())),
	)
	no = 6
	testutil.DoTestCase(
		md,
		testutil.MarkdownTestCase{
			No:          no,
			Description: "WithHardWraps take precedence over WithEastAsianLineBreaks",
			Markdown:    "太郎は\\ **「こんにちわ」**\\ と言った\nんです",
			Expected:    "<p>太郎は\\ <strong>「こんにちわ」</strong>\\ と言った<br />\nんです</p>",
		},
		t,
	)

	// Tests with EastAsianLineBreaksStyleSimple
	md = markdown.New(markdown.WithRendererOptions(
		html.WithXHTML(),
		html.WithUnsafe(),
	),
		markdown.WithExtensions(
			NewCJK(WithEastAsianLineBreaks()),
			Linkify,
		),
	)
	no = 7
	testutil.DoTestCase(
		md,
		testutil.MarkdownTestCase{
			No:          no,
			Description: "WithEastAsianLineBreaks and linkfy extension",
			Markdown:    "太郎は\\ **「こんにちわ」**\\ と言った\r\nんです",
			Expected:    "<p>太郎は\\ <strong>「こんにちわ」</strong>\\ と言ったんです</p>",
		},
		t,
	)
	no = 8
	testutil.DoTestCase(
		md,
		testutil.MarkdownTestCase{
			No:          no,
			Description: "Soft line breaks between east asian wide characters or punctuations are ignored",
			Markdown:    "太郎は\\ **「こんにちわ」**\\ と、\r\n言った\r\nんです",
			Expected:    "<p>太郎は\\ <strong>「こんにちわ」</strong>\\ と、言ったんです</p>",
		},
		t,
	)
	no = 9
	testutil.DoTestCase(
		md,
		testutil.MarkdownTestCase{
			No:          no,
			Description: "Soft line breaks between an east asian wide character and a western character are ignored",
			Markdown:    "私はプログラマーです。\n東京の会社に勤めています。\nGoでWebアプリケーションを開発しています。",
			Expected:    "<p>私はプログラマーです。東京の会社に勤めています。\nGoでWebアプリケーションを開発しています。</p>",
		},
		t,
	)

	// Tests with EastAsianLineBreaksCSS3Draft
	md = markdown.New(markdown.WithRendererOptions(
		html.WithXHTML(),
		html.WithUnsafe(),
	),
		markdown.WithExtensions(
			NewCJK(WithEastAsianLineBreaks(EastAsianLineBreaksCSS3Draft)),
		),
	)
	no = 10
	testutil.DoTestCase(
		md,
		testutil.MarkdownTestCase{
			No:          no,
			Description: "Soft line breaks between a western character and an east asian wide character are ignored",
			Markdown:    "太郎は\\ **「こんにちわ」**\\ と言ったa\nんです",
			Expected:    "<p>太郎は\\ <strong>「こんにちわ」</strong>\\ と言ったaんです</p>",
		},
		t,
	)

	no = 11
	testutil.DoTestCase(
		md,
		testutil.MarkdownTestCase{
			No:          no,
			Description: "Soft line breaks between an east asian wide character and a western character are ignored",
			Markdown:    "太郎は\\ **「こんにちわ」**\\ と言った\nbんです",
			Expected:    "<p>太郎は\\ <strong>「こんにちわ」</strong>\\ と言ったbんです</p>",
		},
		t,
	)

	no = 12
	testutil.DoTestCase(
		md,
		testutil.MarkdownTestCase{
			No:          no,
			Description: "Soft line breaks between an east asian wide character and a western character are ignored",
			Markdown:    "私はプログラマーです。\n東京の会社に勤めています。\nGoでWebアプリケーションを開発しています。",
			Expected:    "<p>私はプログラマーです。東京の会社に勤めています。GoでWebアプリケーションを開発しています。</p>",
		},
		t,
	)

}
