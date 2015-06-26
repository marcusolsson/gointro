package main

import (
	"archive/zip"
	"crypto/md5"
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	var (
		input  string
		datset string
	)

	flag.StringVar(&input, "f", "", "input file")
	flag.StringVar(&datset, "d", "", "datset")
	flag.Parse()

	if input == "" || datset == "" {
		flag.Usage()
		os.Exit(1)
	}

	// Parse datset.
	f, err := os.Open(datset)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	p := NewParser(f)

	col, err := p.Parse()
	if err != nil {
		panic(err)
	}

	roms := []ROM{}

	if filepath.Ext(input) == ".zip" {
		r, err := zip.OpenReader(input)
		if err != nil {
			panic(err)
		}
		defer r.Close()

		for _, f := range r.File {
			rc, err := f.Open()
			if err != nil {
				panic(err)
			}
			defer rc.Close()

			h1, h2 := hashReader(rc)
			r, err := findROM(col, h1, h2)
			if err == nil {
				roms = append(roms, r)
			}
		}
	} else {
		in, err := os.Open(input)
		if err != nil {
			panic(err)
		}
		defer in.Close()

		h1, h2 := hashReader(in)
		r, err := findROM(col, h1, h2)
		if err == nil {
			roms = append(roms, r)
		}
	}

	b, err := json.MarshalIndent(roms, "", "  ")
	if err != nil {
		panic(err)
	}

	fmt.Println(string(b))
}

func findROM(col *Collection, md5hash, sha1hash string) (ROM, error) {
	for _, g := range col.Games {
		for _, r := range g.ROM {
			if r.MD5 != md5hash && r.SHA1 != sha1hash {
				continue
			}
			return r, nil
		}
	}

	return ROM{}, errors.New("rom not found")
}

func hashReader(r io.Reader) (md5hash string, sha1hash string) {
	h1, h2 := md5.New(), sha1.New()

	io.Copy(io.MultiWriter(h1, h2), r)

	md5hash = strings.ToUpper(hex.EncodeToString(h1.Sum(nil)))
	sha1hash = strings.ToUpper(hex.EncodeToString(h2.Sum(nil)))

	return
}
