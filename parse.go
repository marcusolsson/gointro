package gotro

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log"
	"strings"
)

// Collection describes a collection of games.
type Collection struct {
	FileInfo FileInfo `json:"clrmamepro"`
	Games    []Game   `json:"games"`
}

// FileInfo contains info about the dat file.
type FileInfo struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Version     string `json:"version"`
	Comment     string `json:"comment,omitempty"`
}

// Game describes a game ...
type Game struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Serial      string `json:"serial,omitempty"`
	ROM         ROM    `json:"rom"`
}

// ROM contains means to validate a ROM file.
type ROM struct {
	Name  string `json:"name"`
	Size  string `json:"size"`
	CRC   string `json:"crc,omitempty"`
	MD5   string `json:"md5,omitempty"`
	SHA1  string `json:"sha1,omitempty"`
	Flags string `json:"flags,omitempty"`
}

// Parse reads a dat file from Dat-O-Matic and returns the parsed collection..
func Parse(r io.Reader) (Collection, error) {

	s := &tokenStream{
		tokens:  tokenize(r),
		current: -1,
	}

	c := Collection{
		Games: make([]Game, 0),
	}

	for {
		val, err := s.peek()
		if err != nil {
			break
		}

		switch val {
		case "game":
			s.consume("game")
			s.consume("(")

			g, err := parseGame(s)
			if err != nil {
				return Collection{}, err
			}

			s.consume(")")

			c.Games = append(c.Games, g)
		case "clrmamepro":
			s.consume("clrmamepro")
			s.consume("(")

			info, err := parseFileInfo(s)
			if err != nil {
				return Collection{}, err
			}

			s.consume(")")

			c.FileInfo = info
		}
	}

	return c, nil
}

var errEndOfStream = errors.New("end of stream")

type tokenStream struct {
	tokens  []string
	current int
}

func (s *tokenStream) next() (string, error) {
	if s.current < len(s.tokens)-1 {
		s.current = s.current + 1
		return s.tokens[s.current], nil
	}

	return "", errEndOfStream
}

func (s *tokenStream) consume(str string) {
	if s.tokens[s.current+1] == str {
		s.current = s.current + 1
	}
}

func (s *tokenStream) peek() (string, error) {
	if s.current < len(s.tokens)-1 {
		return s.tokens[s.current+1], nil
	}

	return "", errEndOfStream
}

type parseState int

const (
	initialState parseState = iota
	quoteState
)

func tokenize(reader io.Reader) []string {
	var tokens []string
	var buf string

	state := initialState

	r := bufio.NewReader(reader)

	for {
		b, err := r.ReadByte()
		if err == io.EOF {
			out := strings.TrimSpace(buf)
			if len(out) > 0 {
				tokens = append(tokens, out)
			}
			buf = ""
			break
		}

		s := string(b)

		switch state {
		case initialState:
			switch s {
			case "\"":
				state = quoteState
			default:
				if len(strings.TrimSpace(s)) == 0 {
					out := strings.TrimSpace(buf)
					if len(out) > 0 {
						tokens = append(tokens, out)
					}
					buf = ""
				}
			}
		case quoteState:
			switch s {
			case "\"":
				state = initialState
			}
		}

		buf = buf + s
	}

	return tokens
}

func parseFileInfo(s *tokenStream) (FileInfo, error) {
	var info FileInfo

	parsing := true
	for parsing {
		var err error

		val, err := s.peek()
		if err != nil {
			log.Fatal(err)
		}

		switch val {
		case "name":
			if info.Name, err = parseKey("name", s); err != nil {
				return FileInfo{}, parseErr("name")
			}
		case "description":
			if info.Description, err = parseKey("description", s); err != nil {
				return FileInfo{}, parseErr("description")
			}
		case "version":
			if info.Version, err = parseKey("version", s); err != nil {
				return FileInfo{}, parseErr("version")
			}
		case "comment":
			if info.Comment, err = parseKey("comment", s); err != nil {
				return FileInfo{}, parseErr("comment")
			}
		case ")":
			parsing = false
		default:
			return FileInfo{}, fmt.Errorf("unexpected token, %v", val)
		}
	}

	return info, nil
}

func parseGame(s *tokenStream) (Game, error) {
	var g Game

	parsing := true
	for parsing {
		var err error

		val, err := s.peek()
		if err != nil {
			log.Fatal(err)
		}

		switch val {
		case "name":
			if g.Name, err = parseKey("name", s); err != nil {
				return Game{}, parseErr("name")
			}
		case "description":
			if g.Description, err = parseKey("description", s); err != nil {
				return Game{}, parseErr("description")
			}
		case "serial":
			if g.Serial, err = parseKey("serial", s); err != nil {
				return Game{}, parseErr("serial")
			}
		case "rom":
			if g.ROM, err = parseROM(s); err != nil {
				return Game{}, parseErr("rom")
			}
		case ")":
			parsing = false
		default:
			return Game{}, fmt.Errorf("unexpected token, %v", val)
		}
	}

	return g, nil
}

func parseROM(s *tokenStream) (ROM, error) {
	s.consume("rom")
	s.consume("(")

	var r ROM

	parsing := true
	for parsing {
		var err error

		val, err := s.peek()
		if err != nil {
			log.Fatal(err)
		}

		switch val {
		case "name":
			if r.Name, err = parseKey("name", s); err != nil {
				return ROM{}, parseErr("name")
			}
		case "size":
			if r.Size, err = parseKey("size", s); err != nil {
				return ROM{}, parseErr("size")
			}
		case "crc":
			if r.CRC, err = parseKey("crc", s); err != nil {
				return ROM{}, parseErr("crc")
			}
		case "md5":
			if r.MD5, err = parseKey("md5", s); err != nil {
				return ROM{}, parseErr("md5")
			}
		case "sha1":
			if r.SHA1, err = parseKey("sha1", s); err != nil {
				return ROM{}, parseErr("sha1")
			}
		case "flags":
			if r.Flags, err = parseKey("flags", s); err != nil {
				return ROM{}, parseErr("flags")
			}
		case ")":
			parsing = false
		default:
			return ROM{}, fmt.Errorf("unexpected token, %v", val)
		}
	}

	s.consume(")")

	return r, nil
}

func parseKey(key string, s *tokenStream) (string, error) {
	s.consume(key)
	return s.next()
}

func parseErr(s string) error {
	return fmt.Errorf("could not parse token, %s", s)
}
