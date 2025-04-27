# id3stat

## Summary

Checks if a proper ID3 tag is added to the specified MP3 file.

## Motivation

My SatNav/Player cannot read MP3 ID3 tag version 2.x but only version 1.0,
meanwhile some MP3 files come without ID3v1 tag.  This makes me write
a tool to test and report if an MP3 file has an ID3v1 tag.

## Usage

    id3stat mp3file [...]
    id3stat --files=<list> --encoding=<encoding>
    id3stat --dir=<directory>
    id3stat -L
    id3stat -V
    id3stat -H

The first syntax checks if a specified MP3 file has an ID3v1 tag and prints
the file name if the file has no ID3v1 tag.  Two or more files can be
specified.

The second syntax gives a list of files to test, with specifying the encoding
of the content of the list file.  The _list_ parameter specifies the file
name of a text file consisting lines that have an MP3 file name on each.
The _encoding_ parameter is the encoding of the list file itself, neither
the encoding of MP3 files nor of MP3 file names on the file system.
Currently `UTF-8` (default) and `ShiftJIS` are supported.

The third syntax gives a _directory_ to test files in.  All MP3 files are
tested in this directory and descendants.

The `-L` flag indicates to display a licensing notice.  The `-V` flag
indicates to display the version number of `id3stat`.  The `-H` flag
indicates to display the usage help.

## Limitation

By the original requirement, `id3stat` supports ID3v1 tag only.

## Licensing

This tool is licensed under the
[Apache License, Version 2.0](http://www.apache.org/licenses/LICENSE-2.0).
Different licences may apply to depending packages.  See `NOTICE.txt`
file for details.

## How to build

    cd $(GOPATH)/src
    go get -u github.com/dhowden/tag \
              golang.org/x/text/encoding \
              golang.org/x/text/encoding/japanese \
              golang.org/x/text/transform
    go build github.com/upperstream/id3stat

When you modify `NOTICE.txt`, you have to execute `go generate` so that `notice.go` gets updated:

    go generate github.com/upperstream/id3stat
