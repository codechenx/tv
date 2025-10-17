# tv(table viewer) for delimited text file(csv,tsv,etc) in terminal

[![Go Report Card](https://goreportcard.com/badge/github.com/codechenx/tv)](https://goreportcard.com/report/github.com/codechenx/tv)
![test](https://github.com/codechenx/tv/workflows/test/badge.svg)
[![Coverage Status](https://coveralls.io/repos/github/codechenx/tv/badge.svg?branch=master)](https://coveralls.io/github/codechenx/tv?branch=master)
[![GitHub license](https://img.shields.io/github/license/codechenx/tv.svg)](https://github.com/codechenx/tv/blob/master/LICENSE)
[![GitHub release](https://img.shields.io/github/release/codechenx/tv.svg)](http://GitHub.com/codechenx/tv/releases)
[![codechenx-tv](https://snapcraft.io//codechenx-tv/badge.svg)](https://snapcraft.io/codechenx-tv)

<p align="center">
   <img src="data/icon-192x192.png" alt="tv icon"/>
</p>

## Introduction

[![asciicast](https://asciinema.org/a/347295.svg)](https://asciinema.org/a/347295)

## Table of Contents

- [Introduction](#introduction)
- [Feature](#feature)
- [To do](#to-do)
- [Installation](#installation)
  - [Prebuilt binaries](#prebuilt-binaries)
  - [Build from source](#build-from-source)
- [Usage](#usage)
- [Key binding](#key-binding)
- [(Extra)Examples for common biological data](#extraexamples-for-common-biological-data)

## Feature

- Spreadsheet-like view for delimited text data
- Support for gzip compressed file
- Automatically identify separator
- **Progressive loading for large files** - View data immediately as it loads, with responsive UI even for massive datasets

## To do

- [ ] search string

## Installation

### Prebuilt binaries

#### Bash(Linux and macOS, best choice for non-root user)

```bash
curl -fsSL https://raw.githubusercontent.com/codechenx/tv/master/install.sh | bash
```

\* This command will download tv binary file to your current directory, you need to run `sudo cp tv /usr/local/bin/` or copy tv binary file to any directory which is in the environment variable **PATH**

\* You also can download tv binaries manually, from [releases](https://github.com/codechenx/tv/releases)

#### Homebrew(Linux and macOS, only 64-bit)

```bash
brew install codechenx/tv/tv
```

#### Snap(Linux)

[![Get it from the Snap Store](https://snapcraft.io/static/images/badges/en/snap-store-white.svg)](https://snapcraft.io/codechenx-tv)

\*After installation, you need to run `sudo snap alias codechenx-tv tv`.This makes it possible to launch the application by `tv`

#### Debian package(Ubuntu, Debian, etc)

download from [releases](https://github.com/codechenx/tv/releases)

#### RPM package(Centos, Fedora, etc)

download from [releases](https://github.com/codechenx/tv/releases)

### Build from source

Use go get to install and update:

```bash
go get -u github.com/codechenx/tv
```

## Usage

### Usage

tv {File_Name} [flags]

### Flags

```
  --s string       (optional) Split symbol
  --is strings     (optional) Ignore lines with specific prefix(multiple arguments support, separated by comma)
  --in int         (optional) Ignore first N lines
  --nl int         (optional) Only display first N lines
  --dc ints        (optional) Only display specific columns(multiple arguments support, separated by comma)
  --hc ints        (optional) Do not display specific columns(multiple arguments support, separated by comma)
  --fi int         (optional) [default: 0]
                   -1, Unfreeze first row and first column
                    0, Freeze first row and first column
                    1, Freeze first row
                    2, Freeze first column
  --tr             (optional) Transpose data
  --strict         (optional) Check for missing data
  --async          (optional) Load data asynchronously for progressive rendering (default: true)
  -h, --help       help for tv
  -v, --version    version for tv
```

tv also can recive data from pipe as an input

```=
cat file.csv | tv
```

#### Progressive Loading

By default, tv uses asynchronous loading to display data as it's being read. This provides immediate feedback when viewing large files:

```bash
# Async loading (default) - UI appears immediately
tv large_file.csv

# Disable async loading if needed (original behavior)
tv --async=false large_file.csv
```

The footer will show loading progress:
- **For files**: Shows percentage: `Loading... 45.2%` → `Loading... 87.8%` → `Loaded N rows`
- **For pipes**: Shows row count: `Loading... 5234 rows` → `Loaded N rows`
- **Update frequency**: Table refreshes every 20ms (50 FPS) for ultra-smooth progressive rendering

### Sorting and Stats

For tv, there are two data types for every column, **string**, and **number**, which can affect the sorting function and the stats. The data type of the current column is shown on the right of the footer bar.You can change the data type of current column by Ctrl-m. the difference of column data type will determine how the column data would be sorted, as string or as number. In addition, for the column with the number data type, tv will show its minimal value, maximal value, and so on. But for the column with string data type. tv will count the number of every string.

## Key binding

| Key               | description                                            |
| ----------------- | ------------------------------------------------------ |
| ?                 | Help page                                              |
| h, left arrow     | Move left                                              |
| l, right arrow    | Move right                                             |
| j, down arrow     | Move down                                              |
| k, up             | Move up                                                |
| g, home           | Move to first cell of table                            |
| G, end            | Move to last cell of table                             |
| Ctrl-f, page down | Move down by one page                                  |
| Ctrl-b, page up   | Move up by one page                                    |
| Ctrl-e            | Move to end of current column                          |
| Ctrl-h            | Move to head of current column                         |
| Ctrl-m            | Change column data type to string or number            |
| Ctrl-k            | Sort data by column(ascend)                            |
| Ctrl-l            | Sort data by column(descend)                           |
| Ctrl-y            | Show basic stats of current column, back to data table |
| Ctrl-c            | Quit                                                   |

## (Extra)Examples for common biological data

```bash
#vcf or compressed vcf format
tv file.vcf --is "##"
tv file.vcf.gz --is "##"
#qiime otu table
tv file.txt --is "# "
#maf format
tv file.maf --is "#"
#interval list
tv file.interval_list --is "@HD","@SQ"
tv file.interval_list --is "@"
```
