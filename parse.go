package main

import (
	"fmt"
	"io"
)

// Parser represents a parser.
type Parser struct {
	s   *Scanner
	buf struct {
		tok  Token
		val  string
		size int
	}
}

// NewParser returns a new instance of Parser.
func NewParser(r io.Reader) *Parser {
	return &Parser{
		s: NewScanner(r),
	}
}

// Parse parses a Collection.
func (p *Parser) Parse() (*Collection, error) {
	col := &Collection{
		Games: make([]Game, 0),
	}

	for {
		tok, val := p.scanIgnoreWhitespace()
		if tok == EOF {
			break
		} else if tok != IDENT {
			return nil, errUnexpectedToken(val)
		}

		switch val {
		case "clrmamepro":
			p.unscan()

			info, err := p.parseFileInfo()
			if err != nil {
				return nil, err
			}

			col.FileInfo = info
		case "game":
			p.unscan()

			g, err := p.parseGame()
			if err != nil {
				return nil, err
			}

			col.Games = append(col.Games, g)
		default:
			return nil, errUnexpectedToken(val)
		}
	}

	return col, nil
}

func (p *Parser) parseFileInfo() (FileInfo, error) {
	info := FileInfo{}

	var (
		token Token
		value string
	)

	token, value = p.scanIgnoreWhitespace()
	if token != IDENT || value != "clrmamepro" {
		return FileInfo{}, errUnexpectedToken(value)
	}

	token, value = p.scanIgnoreWhitespace()
	if token != LEFTPAREN {
		return FileInfo{}, errUnexpectedToken(value)
	}

	for {
		token, value = p.scanIgnoreWhitespace()
		if token == EOF {
			return FileInfo{}, fmt.Errorf("missing paren")
		} else if token == RIGHTPAREN {
			break
		} else if token != IDENT {
			return FileInfo{}, errUnexpectedToken(value)
		}

		switch value {
		case "name":
			t, v := p.scanIgnoreWhitespace()
			if t != IDENT {
				return FileInfo{}, errUnexpectedToken(v)
			}
			info.Name = v
		case "description":
			t, v := p.scanIgnoreWhitespace()
			if t != IDENT {
				return FileInfo{}, errUnexpectedToken(v)
			}
			info.Description = v
		case "version":
			t, v := p.scanIgnoreWhitespace()
			if t != IDENT {
				return FileInfo{}, errUnexpectedToken(v)
			}
			info.Version = v
		case "comment":
			t, v := p.scanIgnoreWhitespace()
			if t != IDENT {
				return FileInfo{}, errUnexpectedToken(v)
			}
			info.Comment = v
		case "forcenodump":
			t, v := p.scanIgnoreWhitespace()
			if t != IDENT {
				return FileInfo{}, errUnexpectedToken(v)
			}
			info.Forcenodump = v
		default:
			return FileInfo{}, errUnexpectedToken(value)
		}
	}

	return info, nil
}

func (p *Parser) parseGame() (Game, error) {
	game := Game{}

	var (
		token Token
		value string
	)

	token, value = p.scanIgnoreWhitespace()
	if token != IDENT || value != "game" {
		return Game{}, errUnexpectedToken(value)
	}

	token, value = p.scanIgnoreWhitespace()
	if token != LEFTPAREN {
		return Game{}, errUnexpectedToken(value)
	}

	for {
		token, value = p.scanIgnoreWhitespace()
		if token == EOF {
			return Game{}, fmt.Errorf("missing paren")
		} else if token == RIGHTPAREN {
			break
		} else if token != IDENT {
			return Game{}, errUnexpectedToken(value)
		}

		switch value {
		case "name":
			t, v := p.scanIgnoreWhitespace()
			if t != IDENT {
				return Game{}, errUnexpectedToken(v)
			}
			game.Name = v
		case "description":
			t, v := p.scanIgnoreWhitespace()
			if t != IDENT {
				return Game{}, errUnexpectedToken(v)
			}
			game.Description = v
		case "serial":
			t, v := p.scanIgnoreWhitespace()
			if t != IDENT {
				return Game{}, errUnexpectedToken(v)
			}
			game.Serial = v
		case "rom":
			p.unscan()
			r, err := p.parseROM()
			if err != nil {
				return Game{}, err
			}
			game.ROM = append(game.ROM, r)
		default:
			return Game{}, errUnexpectedToken(value)
		}
	}

	return game, nil
}

func (p *Parser) parseROM() (ROM, error) {
	rom := ROM{}

	var (
		token Token
		value string
	)

	token, value = p.scanIgnoreWhitespace()
	if token != IDENT || value != "rom" {
		return ROM{}, errUnexpectedToken(value)
	}

	token, value = p.scanIgnoreWhitespace()
	if token != LEFTPAREN {
		return ROM{}, errUnexpectedToken(value)
	}

	for {
		token, value = p.scanIgnoreWhitespace()
		if token == EOF {
			return ROM{}, fmt.Errorf("missing paren")
		} else if token == RIGHTPAREN {
			break
		} else if token != IDENT {
			return ROM{}, errUnexpectedToken(value)
		}

		switch value {
		case "name":
			t, v := p.scanIgnoreWhitespace()
			if t != IDENT {
				return ROM{}, errUnexpectedToken(v)
			}
			rom.Name = v
		case "size":
			t, v := p.scanIgnoreWhitespace()
			if t != IDENT {
				return ROM{}, errUnexpectedToken(v)
			}
			rom.Size = v
		case "crc":
			t, v := p.scanIgnoreWhitespace()
			if t != IDENT {
				return ROM{}, errUnexpectedToken(v)
			}
			rom.CRC = v
		case "md5":
			t, v := p.scanIgnoreWhitespace()
			if t != IDENT {
				return ROM{}, errUnexpectedToken(v)
			}
			rom.MD5 = v
		case "sha1":
			t, v := p.scanIgnoreWhitespace()
			if t != IDENT {
				return ROM{}, errUnexpectedToken(v)
			}
			rom.SHA1 = v
		case "flags":
			t, v := p.scanIgnoreWhitespace()
			if t != IDENT {
				return ROM{}, errUnexpectedToken(v)
			}
			rom.Flags = v
		default:
			return ROM{}, errUnexpectedToken(value)
		}
	}

	return rom, nil
}

func (p *Parser) scanIgnoreWhitespace() (Token, string) {
	tok, val := p.scan()
	if tok == WS {
		return p.scan()
	}
	return tok, val
}

func (p *Parser) scan() (tok Token, val string) {
	if p.buf.size != 0 {
		p.buf.size = 0
		return p.buf.tok, p.buf.val
	}

	tok, val = p.s.Scan()

	p.buf.tok, p.buf.val = tok, val

	return
}

func (p *Parser) unscan() {
	p.buf.size = 1
}

func errUnexpectedToken(t string) error {
	return fmt.Errorf("unexpected token %s", t)
}
