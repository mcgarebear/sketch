# sketch

[![Go Report Card](https://goreportcard.com/badge/github.com/mcgarebear/sketch?style=flat-square)](https://goreportcard.com/report/github.com/mcgarebear/sketch)
[![Go Doc](https://img.shields.io/badge/godoc-reference-blue.svg?style=flat-square)](http://godoc.org/github.com/mcgarebear/sketch)
[![PkgGoDev](https://pkg.go.dev/badge/github.com/mcgarebear/sketch)](https://pkg.go.dev/github.com/mcgarebear/sketch)
[![Release](https://img.shields.io/github/release/mcgarebear/sketch.svg?style=flat-square)](https://github.com/mcgarebear/sketch/releases/latest)

A cool, little program that can play gifs in the terminal. Yes, it supports
fancy colors too. Shading is based on the color intensity of the color pixel
and can be configured by setting `SKETCH_SHADER` during invocation. This must
be a string of characters representing the gradient white to black.

## Usage

```bash
SKETCH_PATH='/path/to/a/gif' sketch
```

## Install

```bash
go install github.com/mcgarebear/sketch@latest
```

or

```bash
git clone git@github.com:mcgarebear/sketch.git
cd sketch
go install
```

## Other

Clone pokemon sprites from
[PokeAPI/sprites](https://github.com/PokeAPI/sprites)
for some dope pokemon sprites.

```
git clone git@github.com/PokeAPI/sprites
ls -l sprites/sprites/pokemon/versions/generation-v/black-white/
```
