package main

import (
	"strings"
	"testing"
)

type tokenValuePair struct {
	Token Token
	Value string
}

var tests = []struct {
	In  string
	Out []tokenValuePair
}{
	{In: ` `, Out: []tokenValuePair{{Token: WS, Value: " "}}},
	{In: `"`, Out: []tokenValuePair{{Token: IDENT, Value: `"`}}},
	{In: `%`, Out: []tokenValuePair{{Token: ILLEGAL, Value: `%`}}},
	{
		In: "test",
		Out: []tokenValuePair{
			{Token: IDENT, Value: "test"},
		},
	},
	{
		In: `key "value"`,
		Out: []tokenValuePair{
			{Token: IDENT, Value: "key"},
			{Token: WS, Value: " "},
			{Token: IDENT, Value: `"value"`},
		},
	},
	{
		In: `game (
			name "Test Name"
		)`,
		Out: []tokenValuePair{
			{Token: IDENT, Value: "game"},
			{Token: WS, Value: " "},
			{Token: LEFTPAREN, Value: `(`},
			{Token: WS, Value: " "},
			{Token: IDENT, Value: "name"},
			{Token: WS, Value: " "},
			{Token: IDENT, Value: `"Test Name"`},
			{Token: WS, Value: " "},
			{Token: RIGHTPAREN, Value: `)`},
		},
	},
	{
		In: `rom ( name "Test Name" size 2621440 )`,
		Out: []tokenValuePair{
			{Token: IDENT, Value: "rom"},
			{Token: WS, Value: " "},
			{Token: LEFTPAREN, Value: `(`},
			{Token: WS, Value: " "},
			{Token: IDENT, Value: "name"},
			{Token: WS, Value: " "},
			{Token: IDENT, Value: `"Test Name"`},
			{Token: WS, Value: " "},
			{Token: IDENT, Value: "size"},
			{Token: WS, Value: " "},
			{Token: IDENT, Value: "2621440"},
			{Token: WS, Value: " "},
			{Token: RIGHTPAREN, Value: `)`},
		},
	},
	{
		In: `game (
			name "Test Name"
			description "Test Description"
			rom ( name "Test Name" size 2621440 )
		)`,
		Out: []tokenValuePair{
			{Token: IDENT, Value: "game"},
			{Token: WS, Value: " "},
			{Token: LEFTPAREN, Value: `(`},
			{Token: WS, Value: " "},
			{Token: IDENT, Value: "name"},
			{Token: WS, Value: " "},
			{Token: IDENT, Value: `"Test Name"`},
			{Token: WS, Value: " "},
			{Token: IDENT, Value: "description"},
			{Token: WS, Value: " "},
			{Token: IDENT, Value: `"Test Description"`},
			{Token: WS, Value: " "},
			{Token: IDENT, Value: "rom"},
			{Token: WS, Value: " "},
			{Token: LEFTPAREN, Value: `(`},
			{Token: WS, Value: " "},
			{Token: IDENT, Value: "name"},
			{Token: WS, Value: " "},
			{Token: IDENT, Value: `"Test Name"`},
			{Token: WS, Value: " "},
			{Token: IDENT, Value: "size"},
			{Token: WS, Value: " "},
			{Token: IDENT, Value: "2621440"},
			{Token: WS, Value: " "},
			{Token: RIGHTPAREN, Value: `)`},
		},
	},
}

func TestScanner(t *testing.T) {
	for _, tt := range tests {
		r := strings.NewReader(tt.In)
		s := NewScanner(r)

		for _, p := range tt.Out {
			tok, str := s.Scan()
			if tok != p.Token {
				t.Fatalf("unexpected token: %s, expected %s", tok, p.Token)
			}

			if str != p.Value {
				t.Fatalf("unexpected value: %s, expected %s", str, p.Value)
			}
		}
	}
}
