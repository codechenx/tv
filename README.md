# tv(table viewer) for delimited file in terminal
[![GitHub release](https://img.shields.io/github/release/codechenx/tv.svg)](http://GitHub.com/codechenx/tv/releases)
[![Go Report Card](https://goreportcard.com/badge/github.com/codechenx/tv)](https://goreportcard.com/report/github.com/codechenx/tv)
[![GoDoc](https://godoc.org/github.com/codechenx/tv?status.svg)](https://godoc.org/github.com/codechenx/tv)
[![GitHub license](https://img.shields.io/github/license/codechenx/tv.svg)](https://github.com/codechenx/tv/blob/master/LICENSE)

#### Description

tv is a tool to view the delimited file in terminal.


 ![Screenshot](screenshots/example.gif)


# Table of Contents

- [Description](#description)
- [Feature](#feature)
- [To do](#to-do)
- [Installation](#installation)
  - [Prebuilt binaries](#prebuilt-binariesonly-64bit)
  - [Build from source](#build-from-source)
- [Key binding](#key-binding)
- [Usage](#usage)
- [(Extra)Examples for common biological data](#extraexamples-for-common-biological-data)

# Feature

- Spreadsheet-like view for delimited data
- Vim-like key binding 
- Support for gzip compressed file
- Automatically identify seperator

# To do

- [ ] search string


# Installation

### Prebuilt binaries(only x86_64)

#### Homebrew(macOS)
```bash
brew install codechenx/tv/tv
```
#### Bash(Linux and macOS)
```bash
$ curl https://raw.githubusercontent.com/codechenx/tv/master/install.sh | bash
```
#### Debian package(Ubuntu, Debian, etc)
download from [releases](https://github.com/codechenx/tv/releases) 

#### RPM package(Centos, Fedora, etc)
download from [releases](https://github.com/codechenx/tv/releases) 

### Build from source

 Use go get to install and update:
```bash
$ go get -u github.com/codechenx/tv
```
# Key binding

| Key               | description              |
| ----------------- | ------------------------ |
| ? | Help page |
| h, left arrow     | Move left |
| l, right arrow    | Move right |
| j, down arrow     | Move down|
| k, up             | Move up     |
| g, home           | move to first cell of table        |
| G, end            | move to last cell of table      |
| Ctrl-f, page down | Move down by one page    |
| Ctrl-b, page up  | Move up by one page      |
| Ctrl-e | Move to end of current column |
| Ctrl-h | Move to head of current column |
| Ctrl-m | Change column data type to string or number |
| Ctrl-k | Sort data by column(ascend) |
| Ctrl-l | Sort data by column(descend) |
| Ctrl-y | Show basic stats of current column, back to data table |

# Usage
#### Usage:
  tv {File_Name} [flags]

#### Flags:
```bash
  --s string       Split symbol
  --nl int         Only display first N line
  --is strings     Ignore lines with specific prefix(support for multiple arguments, separated by comma
  --in int         Ignore first N row [default: 0]
  --dc ints        Only display certain columns(support for multiple arguments, separated by comma)
  --hc ints        Do not display certain columns(support for multiple arguments, separated by comma)
  --fi int         -1, Unfreeze first row and first column; 0, Freeze first row and first column; 1, Freeze first row; 2, Freeze first column [default: 0]
  --tr             Transpose and view data [default: false]
  -h, --help         help for tv
  -v, --version     version for tv
```


tv also can recive data form pipe as an input

  ```bash
  cat file.csv | tv
  ```



#### Sorting and Stats

For tv, there are two data types for every column, **string**, and **number**, which can affect the sorting function and the stats. The data type of the current column is shown on the right of the footer bar.You can change the data type of current column by Ctrl-m. the difference of column data type will determine how the column data would be sorted, as string or as number. In addition, for the column with the number data type, tv will show its minimal value, maximal value, and so on. But for the column with string data type. tv will count the number of every string.





# (Extra)Examples for common biological data

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
