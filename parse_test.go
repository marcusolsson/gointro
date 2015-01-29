package gotro

import (
	"io"
	"reflect"
	"strings"
	"testing"
)

func TestTokenStream(t *testing.T) {
	s := tokenStream{
		tokens:  []string{"a", "ab", "abc"},
		current: -1,
	}

	if tok, err := s.next(); err != nil || tok != "a" {
		t.Error(err)
	}

	if tok, err := s.next(); err != nil || tok != "ab" {
		t.Error(err)
	}

	if tok, err := s.next(); err != nil || tok != "abc" {
		t.Error(err)
	}

	if _, err := s.next(); err == nil {
		t.Error(err)
	}
}

func TestConsumeFromTokenStream(t *testing.T) {
	s := tokenStream{
		tokens:  []string{"a", "ab", "abc"},
		current: -1,
	}

	s.consume("a")

	if s.current != 0 {
		t.Fail()
	}

	if tok, _ := s.next(); tok != "ab" {
		t.Fail()
	}
}

func TestPeekFromTokenStream(t *testing.T) {
	s := tokenStream{
		tokens:  []string{"a", "ab", "abc"},
		current: -1,
	}

	if tok, err := s.peek(); tok != "a" {
		t.Error(err)
	}

	if s.current != -1 {
		t.Fail()
	}

	if tok, err := s.next(); tok != "a" {
		t.Error(err)
	}
}

func testReader() io.Reader {
	s := `clrmamepro (
        name "Test Name"
        description "Test Description"
        version 20080101-123456
        comment "Test Comment"
)

game (
	name "Test Name"
	description "Test Description"
	rom ( name "Test Name" size 2621440 crc C167987D md5 A990AE4416DD75F7C68C5DB06425D648 sha1 21286747D360C03E3BF86CD4504508CE55DEFF8F )
)

game (
	name "Test Name"
	description "Test Description"
	rom ( name "Test Name" size 2621440 crc C167987D md5 A990AE4416DD75F7C68C5DB06425D648 sha1 21286747D360C03E3BF86CD4504508CE55DEFF8F )
)`

	return strings.NewReader(s)
}

func TestTokenize(t *testing.T) {
	tokens := tokenize(testReader())
	expected := []string{"clrmamepro", "(", "name", "\"Test Name\"", "description", "\"Test Description\"", "version", "20080101-123456", "comment", "\"Test Comment\"", ")", "game", "(", "name", "\"Test Name\"", "description", "\"Test Description\"", "rom", "(", "name", "\"Test Name\"", "size", "2621440", "crc", "C167987D", "md5", "A990AE4416DD75F7C68C5DB06425D648", "sha1", "21286747D360C03E3BF86CD4504508CE55DEFF8F", ")", ")", "game", "(", "name", "\"Test Name\"", "description", "\"Test Description\"", "rom", "(", "name", "\"Test Name\"", "size", "2621440", "crc", "C167987D", "md5", "A990AE4416DD75F7C68C5DB06425D648", "sha1", "21286747D360C03E3BF86CD4504508CE55DEFF8F", ")", ")"}

	if !reflect.DeepEqual(tokens, expected) {
		t.Errorf("%#v", tokens)
	}
}

func TestParse(t *testing.T) {
	expected := Collection{
		FileInfo: FileInfo{
			Name:        "\"Test Name\"",
			Description: "\"Test Description\"",
			Version:     "20080101-123456",
			Comment:     "\"Test Comment\"",
		},
		Games: []Game{
			{
				Name:        "\"Test Name\"",
				Description: "\"Test Description\"",
				ROM: ROM{
					Name: "\"Test Name\"",
					Size: "2621440",
					CRC:  "C167987D",
					MD5:  "A990AE4416DD75F7C68C5DB06425D648",
					SHA1: "21286747D360C03E3BF86CD4504508CE55DEFF8F",
				},
			},
			{
				Name:        "\"Test Name\"",
				Description: "\"Test Description\"",
				ROM: ROM{
					Name: "\"Test Name\"",
					Size: "2621440",
					CRC:  "C167987D",
					MD5:  "A990AE4416DD75F7C68C5DB06425D648",
					SHA1: "21286747D360C03E3BF86CD4504508CE55DEFF8F",
				},
			},
		},
	}

	c, err := Parse(testReader())
	if err != nil {
		t.Error(err)
	}

	if !reflect.DeepEqual(c, expected) {
		t.Errorf("%#v", expected)
	}
}
