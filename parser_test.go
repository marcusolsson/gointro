package main

import (
	"fmt"
	"strings"
	"testing"
)

var parseTests = []struct {
	In  string
	Out *Collection
	Err error
}{
	{In: `unknown ()`, Err: errUnexpectedToken("unknown")},
	{In: `clrmamepro ()`, Out: &Collection{}},
	{In: `clrmamepro (`, Err: fmt.Errorf("missing paren")},
	{In: `clrmamepro (()`, Err: errUnexpectedToken("(")},
	{In: `clrmamepro )`, Err: errUnexpectedToken(")")},
	{In: `clrmamepro ( invalid )`, Err: errUnexpectedToken("invalid")},
	{In: `clrmamepro ( name )`, Err: errUnexpectedToken(")")},
	{In: `clrmamepro ( description )`, Err: errUnexpectedToken(")")},
	{In: `clrmamepro ( version )`, Err: errUnexpectedToken(")")},
	{In: `clrmamepro ( comment )`, Err: errUnexpectedToken(")")},
	{In: `game ()`, Out: &Collection{Games: make([]Game, 1)}},
	{In: `game (`, Err: fmt.Errorf("missing paren")},
	{In: `game (()`, Err: errUnexpectedToken("(")},
	{In: `game ())`, Err: errUnexpectedToken(")")},
	{In: `game )`, Err: errUnexpectedToken(")")},
	{In: `game ( invalid )`, Err: errUnexpectedToken("invalid")},
	{In: `game ( name )`, Err: errUnexpectedToken(")")},
	{In: `game ( description )`, Err: errUnexpectedToken(")")},
	{In: `game ( serial )`, Err: errUnexpectedToken(")")},
	{In: `game ( rom ())`, Out: &Collection{Games: []Game{{ROM: make([]ROM, 1)}}}},
	{In: `game ( rom (`, Err: fmt.Errorf("missing paren")},
	{In: `game ( rom (()`, Err: errUnexpectedToken("(")},
	{In: `game ( rom )`, Err: errUnexpectedToken(")")},
	{In: `game ( rom ( invalid ) )`, Err: errUnexpectedToken("invalid")},
	{In: `game ( rom ( name ) )`, Err: errUnexpectedToken(")")},
	{In: `game ( rom ( size ) )`, Err: errUnexpectedToken(")")},
	{In: `game ( rom ( crc ) )`, Err: errUnexpectedToken(")")},
	{In: `game ( rom ( md5 ) )`, Err: errUnexpectedToken(")")},
	{In: `game ( rom ( sha1 ) )`, Err: errUnexpectedToken(")")},
	{In: `game ( rom ( flags ) )`, Err: errUnexpectedToken(")")},
	{
		In: `clrmamepro (
        name "Test Name"
        description "Test Description"
        version 20080101-123456
        comment "Test Comment"
)

game (
	name "First Game"
	description "First Game Description"
	serial 123
	rom ( name "Test Name" size 2621440 crc C167987D md5 A990AE4416DD75F7C68C5DB06425D648 sha1 21286747D360C03E3BF86CD4504508CE55DEFF8F flags verified)
)

game (
	name "Second Game"
	description "Second Game Description"
	rom ( name "Test Name" size 2621440 crc C167987D md5 A990AE4416DD75F7C68C5DB06425D648 sha1 21286747D360C03E3BF86CD4504508CE55DEFF8F )
)`,
		Out: &Collection{
			FileInfo: FileInfo{
				Name:        `"Test Name"`,
				Description: `"Test Description"`,
				Version:     "20080101-123456",
				Comment:     `"Test Comment"`,
			},
			Games: []Game{
				{
					Name:        `"First Game"`,
					Description: `"First Game Description"`,
					Serial:      "123",
					ROM: []ROM{
						{
							Name:  `"Test Name"`,
							Size:  "2621440",
							CRC:   "C167987D",
							MD5:   "A990AE4416DD75F7C68C5DB06425D648",
							SHA1:  "21286747D360C03E3BF86CD4504508CE55DEFF8F",
							Flags: "verified",
						},
					},
				},
				{
					Name:        `"Second Game"`,
					Description: `"Second Game Description"`,
					ROM: []ROM{
						{
							Name: `"Test Name"`,
							Size: "2621440",
							CRC:  "C167987D",
							MD5:  "A990AE4416DD75F7C68C5DB06425D648",
							SHA1: "21286747D360C03E3BF86CD4504508CE55DEFF8F",
						},
					},
				},
			},
		},
		Err: nil,
	},
}

func TestParser(t *testing.T) {
	for _, tt := range parseTests {
		r := strings.NewReader(tt.In)
		p := NewParser(r)
		col, err := p.Parse()
		if err != nil {
			if err.Error() != tt.Err.Error() {
				t.Fatalf("unexpected error: %q, expected: %q", err, tt.Err)
			}
			continue
		}

		if col.FileInfo != tt.Out.FileInfo {
			t.Fatalf("expected %q, got %q", tt.Out.FileInfo, col.FileInfo)
		}

		if len(col.Games) != len(tt.Out.Games) {
			t.Fatal("unexpected number of games")
		}

		for i, g := range col.Games {
			e := tt.Out.Games[i]
			if g.Name != e.Name {
				t.Errorf("expected %v, got %v", e.Name, g.Name)
			}
			if g.Description != e.Description {
				t.Errorf("expected %v, got %v", e.Description, g.Description)
			}
			if g.Serial != e.Serial {
				t.Errorf("expected %v, got %v", e.Serial, g.Serial)
			}

			if len(g.ROM) != len(e.ROM) {
				t.Fatal("unexpected number of ROMs")
			}

			for j, rom := range g.ROM {
				if rom != e.ROM[j] {
					t.Errorf("expected %v, got %v", e.ROM[j], rom)
				}
			}
		}
	}
}
