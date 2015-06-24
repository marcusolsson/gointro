package main

import (
	"bufio"
	"bytes"
	"io"
	"unicode"
)

// Token represents a lexical token.
type Token int

// Tokens
const (
	ILLEGAL Token = iota
	WS
	EOF

	IDENT
	INTEGER

	LEFTPAREN
	RIGHTPAREN

	GAMETYPE
	ROMTYPE
)

func (t Token) String() string {
	switch t {
	case ILLEGAL:
		return "ILLEGAL"
	case WS:
		return "WS"
	case EOF:
		return "EOF"
	case IDENT:
		return "IDENT"
	case INTEGER:
		return "INTEGER"
	case LEFTPAREN:
		return "LEFTPAREN"
	case RIGHTPAREN:
		return "RIGHTPAREN"
	case GAMETYPE:
		return "GAMETYPE"
	case ROMTYPE:
		return "ROMTYPE"
	default:
		return ""
	}
}

var eof = rune(0)

// Scanner represents a lexical scanner.
type Scanner struct {
	r *bufio.Reader
}

// NewScanner returns a new instance of Scanner.
func NewScanner(r io.Reader) *Scanner {
	return &Scanner{
		r: bufio.NewReader(r),
	}
}

func (s *Scanner) read() rune {
	ch, _, err := s.r.ReadRune()
	if err != nil {
		return eof
	}
	return ch
}

func (s *Scanner) unread() {
	_ = s.r.UnreadRune()
}

// Scan returns the next token and its value.
func (s *Scanner) Scan() (Token, string) {
	ch := s.read()

	if unicode.IsSpace(ch) {
		s.unread()
		return s.scanSpace()
	} else if unicode.IsLetter(ch) || unicode.IsDigit(ch) {
		s.unread()
		return s.scanIdent()
	} else if isQuote(ch) {
		s.unread()
		return s.scanQuote()
	} else if isParen(ch) {
		if ch == '(' {
			return LEFTPAREN, "("
		} else if ch == ')' {
			return RIGHTPAREN, ")"
		}
	} else if ch == eof {
		return EOF, ""
	}

	return ILLEGAL, ""
}

func isQuote(ch rune) bool {
	return ch == '"'
}

func isParen(ch rune) bool {
	return ch == '(' || ch == ')'
}

func (s *Scanner) scanSpace() (Token, string) {
	var buf bytes.Buffer
	buf.WriteRune(s.read())

	for {
		if ch := s.read(); ch == eof {
			break
		} else if !unicode.IsSpace(ch) {
			s.unread()
			break
		} else {
			buf.WriteRune(ch)
		}
	}

	return WS, " "
}

func (s *Scanner) scanIdent() (Token, string) {
	var buf bytes.Buffer
	buf.WriteRune(s.read())

	for {
		if ch := s.read(); ch == eof {
			break
		} else if !unicode.IsLetter(ch) && !unicode.IsDigit(ch) && ch != '-' {
			s.unread()
			break
		} else {
			_, _ = buf.WriteRune(ch)
		}
	}

	return IDENT, buf.String()
}

func (s *Scanner) scanInt() (Token, string) {
	var buf bytes.Buffer
	buf.WriteRune(s.read())

	for {
		if ch := s.read(); ch == eof {
			break
		} else if !unicode.IsDigit(ch) {
			s.unread()
			break
		} else {
			_, _ = buf.WriteRune(ch)
		}
	}

	return INTEGER, buf.String()
}

func (s *Scanner) scanQuote() (Token, string) {
	var buf bytes.Buffer
	buf.WriteRune(s.read())

	for {
		if ch := s.read(); ch == eof {
			break
		} else if isQuote(ch) {
			_, _ = buf.WriteRune(ch)
			break
		} else {
			_, _ = buf.WriteRune(ch)
		}
	}

	return IDENT, buf.String()
}
