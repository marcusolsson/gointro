# gotro

[![Build Status](https://travis-ci.org/marcusolsson/gotro.svg?branch=master)](https://travis-ci.org/marcusolsson/gotro)
[![GoDoc](https://img.shields.io/badge/godoc-reference-blue.svg?style=flat)](https://godoc.org/github.com/marcusolsson/gotro)
[![License MIT](https://img.shields.io/badge/license-MIT-lightgrey.svg?style=flat)](LICENSE)

Check if a file belongs to a datset from [DAT-o-MATIC](http://datomatic.no-intro.org).

## Usage

```bash
gotro -d <datset> -f <file>
```

Output is the real name according to the datset. If the file was not in the datset, the output will be empty.

## Installation

```bash
go get github.com/marcusolsson/gotro
```
