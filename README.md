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
- **Text wrapping** - Wrap long text in columns for better readability (Ctrl+W)
- **Search functionality** - Search for text within cells and navigate through results
- **Column filter** - Filter rows based on column values (Ctrl+F)
- **Smart column type detection** - Automatically detects String, Number, and Date columns
- **Fast sorting** - Optimized sorting with pre-parsed values for numbers and dates

## To do

- [x] search string

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

For tv, there are three data types for every column: **string**, **number**, and **date**, which affect the sorting function and the stats. The data type of the current column is shown on the right of the footer bar.

**Automatic Type Detection**: When you load a file, tv automatically analyzes each column and detects whether it contains strings, numbers, or dates. The detection algorithm samples data intelligently and uses a 90% threshold to classify columns.

**Manual Type Changes**: You can change the data type of the current column by pressing **Ctrl+M**, which cycles through: String → Number → Date → String.

**Column Type Behavior**:
- **String columns**: Sorted alphabetically
- **Number columns**: Sorted numerically (handles integers, floats, scientific notation, thousand separators)
- **Date columns**: Sorted chronologically (supports multiple date formats including ISO-8601, US, EU formats, and more)

**Performance Optimizations**: Sorting is optimized by pre-parsing all values once and caching them, which makes sorting large datasets significantly faster than repeatedly parsing during comparison.

**Statistics**: For number columns, tv shows min/max values, mean, etc. For string columns, tv counts the frequency of each unique value.

## Key binding

All key bindings follow vim-like conventions for intuitive navigation and operation.

| Key               | description                                            |
| ----------------- | ------------------------------------------------------ |
| ?                 | Show help dialog (modal overlay)                       |
| q                 | Quit                                                   |
| Esc               | Close help/stats dialog, or clear search              |
| h                 | Move left                                              |
| l                 | Move right                                             |
| j                 | Move down                                              |
| k                 | Move up                                                |
| w                 | Move to next column (word forward)                     |
| b                 | Move to previous column (word backward)                |
| gg                | Go to first row (press g twice)                        |
| G                 | Go to last row                                         |
| 0                 | Go to first column                                     |
| $                 | Go to last column                                      |
| Ctrl-d            | Page down (half page)                                  |
| Ctrl-u            | Page up (half page)                                    |
| /                 | Search for text (case-insensitive)                     |
| n                 | Next search result                                     |
| N                 | Previous search result                                 |
| f                 | Filter rows by current column value                    |
| r                 | Reset/clear column filter                              |
| t                 | Toggle column data type (Str -> Num -> Date)          |
| s                 | Sort data by column (ascending)                        |
| S                 | Sort data by column (descending)                       |
| W                 | Toggle text wrapping for current column                |
| i                 | Show stats info for current column                     |

### Text Wrapping

For columns with long text content, press **Ctrl+W** to toggle text wrapping:
- **First press**: Wraps text at 25 characters with smart word breaking
- **Second press**: Unwraps to original single-line view
- **Per-column**: Each column can be wrapped independently
- **Smart wrapping**: Breaks at spaces/hyphens when possible

Example use case: Reading long descriptions or comments without horizontal scrolling.

### Search

Search for text within the table and navigate through results with highlighting:
- **Press `/`**: Opens search input dialog
- **Type your query**: Search is case-insensitive by default
- **Press Enter**: Executes the search and shows the number of matches
- **Press `n`**: Jump to next search result
- **Press `N`**: Jump to previous search result (capital N)
- **Press `Ctrl+/`**: Clear search highlighting

**Highlighting:**
- **Current match**: Highlighted with bright cyan background
- **Other matches**: Highlighted with gray background
- The footer shows your current position in the search results (e.g., "Match 1/5")

Example workflow:
1. Press `/` to open search
2. Type "error" and press Enter
3. All cells containing "error" are highlighted, cursor jumps to the first match
4. Use `n` and `N` to navigate through all cells containing "error"
5. Press `Ctrl+/` to clear the highlighting when done

### Column Filter

Filter table rows based on column values to focus on specific data:
- **Press `f`**: Opens column filter dialog for the current column
- **Type your filter text**: Filter is case-insensitive and matches partial text
- **Press Enter**: Apply the filter - only rows where the column contains the filter text are displayed
- **Press `Ctrl+R`**: Reset/clear the filter and show all rows again

**Features:**
- **Case-insensitive**: Filter matches are not case-sensitive
- **Partial matching**: Shows rows where the column value contains your filter text
- **Header preserved**: The header row is always visible, even when filtered
- **Status indication**: Footer shows how many rows match the filter
- **Easy reset**: Press `Ctrl+R` to return to the full dataset

Example workflow:
1. Navigate to a column you want to filter (e.g., a "Status" column)
2. Press `f` to open the filter dialog
3. Type "active" and press Enter
4. Only rows where that column contains "active" are displayed
5. Press `Ctrl+R` to show all rows again

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
