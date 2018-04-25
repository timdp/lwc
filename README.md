# lwc

[![Go Report Card](https://goreportcard.com/badge/github.com/timdp/lwc)](https://goreportcard.com/report/github.com/timdp/lwc)
[![Build Status](https://img.shields.io/circleci/project/github/timdp/lwc/master.svg)](https://circleci.com/gh/timdp/lwc)
[![Release](https://img.shields.io/github/release/timdp/lwc.svg)](https://github.com/timdp/lwc/releases/latest)

A live-updating version of the UNIX [`wc` command](https://en.wikipedia.org/wiki/Wc_(Unix)).

![](demo.gif)

## Installation

You can get a prebuilt binary for every major platform from the
[Releases page](https://github.com/timdp/lwc/releases). Just extract it
somewhere under your `PATH` and you're good to go.

Alternatively, use `go get` to build from source:

```bash
go get -u github.com/timdp/lwc/cmd/lwc
```

## Usage

```
lwc [OPTION]...
```

Without any options, `lwc` will count the number of lines, words, and bytes
in standard input, and write them to standard output. Contrary to `wc`, it will
also update standard output while it is still counting.

The following [`wc` options](https://en.wikipedia.org/wiki/Wc_(Unix)) are
currently supported:

- `--lines` or `-l`
- `--words` or `-w`
- `--chars` or `-m`
- `--bytes` or `-c`
- `--max-line-length` or `-L`
- `--help`
- `--version`

In addition, the output update interval can be configured by passing either
`--interval=TIME` or `-i TIME`, where `TIME` is a duration in milliseconds.
The default update interval is 100 ms.

## Examples

Count the number of lines in a big file:

```bash
lwc --lines < big-file
```

Run a slow command and count the number of bytes logged:

```bash
slow-command | lwc --bytes
```

## TODO

- Support `--files0-from`
- Add tests

## JavaScript Version

This utility briefly existed as a
[Node.js package](https://github.com/timdp/lwc-nodejs). I'm keeping the code
around for educational purposes, but I will no longer be maintaining it.

## Author

[Tim De Pauw](https://tmdpw.eu/)

## License

MIT
