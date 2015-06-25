package main

import (
	"crypto/md5"
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
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

	var err error

	in, err := os.Open(input)
	if err != nil {
		panic(err)
	}
	defer in.Close()

	h1, h2 := md5.New(), sha1.New()

	io.Copy(io.MultiWriter(h1, h2), in)

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
		if g.ROM.MD5 != md5hash && g.ROM.SHA1 != sha1hash {
			continue
		}

		b, err := json.MarshalIndent(g, "", "  ")
		if err != nil {
			panic(err)
		}

		fmt.Println(string(b))
	}
}
