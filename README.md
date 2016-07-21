# id3stat

## Summary

Checks is a proper ID3 tag is added to the specified MP3 file.

## Motivation

My Satnav/Player cannot read MP3 ID3 tag version 2.x but only version 1.0,
meanwhile some MP3 files come with ID3v2.x tag only.  This makes me write
a tool to test and report if an MP3 file has an ID3v1 tag.

## Usage

    id3stat mp3file [...]
    id3stat -L
    id3stat -V
    id3stat -H

The first syntax checks if a specified MP3 file has an ID3v1 tag.  Two or
more files can be specified.

The `-L` flag indicates to display a licensing notice.  The `-V` flag
indicates to display the version number of `id3stat`.  The `-H` flag
indicates to display the usage help.

## Limination

By the original requirement, `id3stat` only supports ID3v1 tag.

## Licensing

This tool is licensed under the
[Apache License, Version 2.0](http://www.apache.org/licenses/LICENSE-2.0).
Different licences may apply to depending packages.  See `NOTICE.txt`
file for details.
