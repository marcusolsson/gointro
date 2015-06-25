package main

import (
	"strings"
	"testing"
)

func TestParseInvalidFileInfo(t *testing.T) {
	r := strings.NewReader(`clrmamepro ( invalid )`)

	p := NewParser(r)

	_, err := p.Parse()
	if err == nil {
		t.Fatal("unexpected success")
	}
}

func TestParseInvalidGame(t *testing.T) {
	r := strings.NewReader(`game ( invalid )`)

	p := NewParser(r)

	_, err := p.Parse()
	if err == nil {
		t.Fatal("unexpected success")
	}
}

func TestParser(t *testing.T) {
	r := strings.NewReader(`clrmamepro (
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
)`)

	p := NewParser(r)

	col, err := p.Parse()
	if err != nil {
		t.Fatal(err)
	}

	if col.FileInfo.Name != `"Test Name"` {
		t.Fatalf("unexpected collection name: %s", col.FileInfo.Name)
	}

	if col.FileInfo.Description != `"Test Description"` {
		t.Fatalf("unexpected collection description: %s", col.FileInfo.Description)
	}

	if col.FileInfo.Version != `20080101-123456` {
		t.Fatalf("unexpected collection version: %s", col.FileInfo.Version)
	}

	if col.FileInfo.Comment != `"Test Comment"` {
		t.Fatalf("unexpected collection comment: %s", col.FileInfo.Comment)
	}

	if len(col.Games) != 2 {
		t.Fatalf("unexpected number of games: %v", len(col.Games))
	}

	if col.Games[0].Name != `"Test Name"` {
		t.Fatalf("unexpected game name: %s", col.Games[0].Name)
	}

	if col.Games[0].Description != `"Test Description"` {
		t.Fatalf("unexpected game description: %s", col.Games[0].Description)
	}

	if col.Games[0].ROM[0].Name != `"Test Name"` {
		t.Fatalf("unexpected ROM name: %s", col.Games[0].ROM[0].Name)
	}

	if col.Games[0].ROM[0].Size != `2621440` {
		t.Fatalf("unexpected ROM size: %s", col.Games[0].ROM[0].Size)
	}

}
