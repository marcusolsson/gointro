package main

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
	ROM         []ROM  `json:"rom"`
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
