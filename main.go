package main

import (
	"archive/zip"
	"crypto/md5"
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
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

	h1, h2 := md5.New(), sha1.New()

	if filepath.Ext(input) == ".zip" {
		r, err := zip.OpenReader(input)
		if err != nil {
			panic(err)
		}
		defer r.Close()

		if len(r.File) == 0 {
			fmt.Println("archive is empty")
			os.Exit(1)
		}

		if len(r.File) > 1 {
			fmt.Println("multiple file archives are currently not supported")
			os.Exit(1)
		}

		for _, f := range r.File {
			rc, err := f.Open()
			if err != nil {
				panic(err)
			}
			defer rc.Close()

			io.Copy(io.MultiWriter(h1, h2), rc)
		}
	} else {
		in, err := os.Open(input)
		if err != nil {
			panic(err)
		}
		defer in.Close()

		io.Copy(io.MultiWriter(h1, h2), in)
	}

	md5hash := strings.ToUpper(hex.EncodeToString(h1.Sum(nil)))
	sha1hash := strings.ToUpper(hex.EncodeToString(h2.Sum(nil)))

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

	for _, g := range col.Games {
		for _, r := range g.ROM {
			if r.MD5 != md5hash && r.SHA1 != sha1hash {
				continue
			}

			b, err := json.MarshalIndent(r, "", "  ")
			if err != nil {
				panic(err)
			}

			fmt.Println(string(b))
		}
	}
}
