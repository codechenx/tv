# tv(table viewer) for delimited file in terminal
[![GitHub release](https://img.shields.io/github/release/codechenx/tv.svg)](http://GitHub.com/codechenx/tv/releases)
[![Build Status](https://travis-ci.org/codechenx/tv.svg?branch=master)](https://travis-ci.org/codechenx/tv)
[![codecov](https://codecov.io/gh/codechenx/tv/branch/master/graph/badge.svg)](https://codecov.io/gh/codechenx/tv)
[![Go Report Card](https://goreportcard.com/badge/github.com/codechenx/tv)](https://goreportcard.com/report/github.com/codechenx/tv)
[![GitHub license](https://img.shields.io/github/license/codechenx/tv.svg)](https://github.com/codechenx/tv/blob/master/LICENSE)

#### Description

tv is a tool to view the delimited file in terminal

 ![Screenshot](screenshots/example.png)


# Table of Contents

- [Description](#description)
- [Feature](#feature)
- [To do](#to-do)
- [Installation](#installation)
  - [Prebuilt binaries](#prebuilt-binariesonly-64bit)
  - [Build from source](#build-from-source)
- [Key binding](#key-binding)
- [Usage](#usage)
- [(Extra)Examples for common biological data](#examples-for-common-biological-data)

# Feature

- Spreadsheet-like view for delimited data
- Vim-like key binding 
- Support for gzip compressed file
- Automatically identify tsv and csv format(Experimental)

# To do

- [ ] search string
- [ ] sort values of column


# Installation

### Prebuilt binaries(only 64bit)

```bash
$ curl https://raw.githubusercontent.com/codechenx/tv/master/install.sh | bash
```

### Build from source

 Use go get to install and update:
```bash
$ go get -u github.com/codechenx/tv
```
# Key binding

  | Key               | description              |
  | ----------------- | ------------------------ |
  | h, left arrow     | Move left by one column  |
  | l, right arrow    | Move right by one column |
  | j, down arrow     | Move down by one row     |
  | k, up             | Move up by one row       |
  | g, home           | Move to the top          |
  | G, end            | Move to the bottom       |
  | Ctrl-F, page down | Move down by one page    |
  | Ctrl-B, page up   | Move up by one page      |

# Usage

Usage: tv [--sep SEP] [--ss SS] [--h H] [--t] FILENAME

Positional arguments:
  FILENAME

Options:
  - --sep SEP, -s SEP      split symbol
  - --ss SS                ignore lines with specific prefix
  - --sn SN                ignore first n lines
  - --h H(default:0)       -1, no column name and row name; 0, use first row as row name; 1, use first column as column name; 2, use first column as column name and first row as row name
  - --t                    transpose and view data
  - --help, -h             display this help and exit
  - --version              display version and exit
  
# Examples for common biological data
```bash
#vcf or compressed vcf format
tv file.vcf --ss "##"
tv file.vcf.gz --ss "##"
#qiime otu table
tv file.txt --ss "# "
#maf format
tv file.maf --ss "#"
```