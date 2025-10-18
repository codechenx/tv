# tv - Table Viewer for Terminal

**A fast, feature-rich CSV/TSV/delimited file viewer for the command line**

[![Go Report Card](https://goreportcard.com/badge/github.com/codechenx/tv)](https://goreportcard.com/report/github.com/codechenx/tv)
![test](https://github.com/codechenx/tv/workflows/test/badge.svg)
[![Coverage Status](https://coveralls.io/repos/github/codechenx/tv/badge.svg?branch=master)](https://coveralls.io/github/codechenx/tv?branch=master)
[![GitHub license](https://img.shields.io/github/license/codechenx/tv.svg)](https://github.com/codechenx/tv/blob/master/LICENSE)
[![GitHub release](https://img.shields.io/github/release/codechenx/tv.svg)](http://GitHub.com/codechenx/tv/releases)


<p align="center">
   <img src="data/icon-192x192.png" alt="tv icon"/>
</p>

## Demo

[![asciicast](https://asciinema.org/a/347295.svg)](https://asciinema.org/a/347295)

## Table of Contents

- [Features](#features)
- [Installation](#installation)
  - [Quick Install (Linux/macOS)](#quick-install-linuxmacos)
  - [Package Managers](#package-managers)
  - [Build from Source](#build-from-source)
- [Quick Start](#quick-start)
- [Command Line Flags](#command-line-flags)
- [Key Bindings](#key-bindings)
- [Features in Detail](#features-in-detail)
  - [Progressive Loading](#progressive-loading)
  - [Data Types and Sorting](#data-types-and-sorting)
  - [Search](#search)
  - [Column Filter](#column-filter)
  - [Text Wrapping](#text-wrapping)
- [Advanced Examples](#advanced-examples)
  - [Biological Data Formats](#biological-data-formats)

## Features

tv brings spreadsheet-like functionality to your terminal with vim-inspired controls.

- **Spreadsheet interface** - Navigate and view tabular data with frozen headers
- **Smart parsing** - Automatically detects delimiters (CSV, TSV, custom separators)
- **Progressive loading** - Start viewing large files immediately while they load
- **Gzip support** - Read compressed files directly
- **Powerful search** - Find text across all cells with highlighting
- **Column filtering** - Show only rows matching specific criteria
- **Flexible sorting** - Sort by any column with intelligent type detection
- **Text wrapping** - Wrap long cell content for better readability
- **Statistics** - View column stats (min/max, mean for numbers; frequency for strings)
- **Vim keybindings** - Navigate naturally with h/j/k/l and more
- **Pipe support** - Read from stdin for seamless integration with shell pipelines


## Installation

### Recommended: Install Script (Linux/macOS)

```bash
curl -sSL https://raw.githubusercontent.com/codechenx/tv/master/install.sh | bash
sudo mv tv /usr/local/bin/
```

### Package Managers

**Snap (Linux)**

[![Get it from the Snap Store](https://snapcraft.io/static/images/badges/en/snap-store-white.svg)](https://snapcraft.io/codechenx-tv)

```bash
sudo snap install codechenx-tv
sudo snap alias codechenx-tv tv
```

**Go Install**

```bash
go install github.com/codechenx/tv@latest
```

### Linux Package Managers

**Debian/Ubuntu (.deb)**

```bash
# Download from releases page
wget https://github.com/codechenx/tv/releases/latest/download/tv_0.6.2_Linux_x86_64.deb
sudo dpkg -i tv_*.deb
```

**CentOS/Fedora (.rpm)**

```bash
# Download from releases page
wget https://github.com/codechenx/tv/releases/latest/download/tv_0.6.2_Linux_x86_64.rpm
sudo rpm -i tv_*.rpm
```

### Manual Download

Download pre-built binaries from [releases](https://github.com/codechenx/tv/releases) for:
- Linux (x86_64, ARM, ARM64)
- macOS (Intel, Apple Silicon)
- Windows (x86_64, i386)

### Build from Source

Requires Go 1.21 or later:

```bash
git clone https://github.com/codechenx/tv.git
cd tv
go build -ldflags="-s -w" -o tv
```

## Quick Start

View a CSV file:

```bash
tv data.csv
```

View a TSV file (tab-separated):

```bash
tv data.tsv
```

Read from stdin:

```bash
cat data.csv | tv
ps aux | tv
```

Specify a custom delimiter:

```bash
tv data.txt -s "|"
```

View only specific columns:

```bash
tv data.csv --columns 1,3,5
```

Skip header lines (e.g., for VCF files):

```bash
tv file.vcf --skip-prefix "##"
```

## Command Line Flags

## Command Line Flags

**Syntax:** `tv [FILE] [flags]`

| Flag | Short | Description |
|------|-------|-------------|
| `--separator` | `-s` | Delimiter character (use `\t` for tab) |
| `--lines` | `-n` | Display only first N lines |
| `--skip-prefix` | | Skip lines starting with prefix (comma-separated) |
| `--skip-lines` | | Skip first N lines |
| `--columns` | | Show only specified columns (comma-separated) |
| `--hide-columns` | | Hide specified columns (comma-separated) |
| `--freeze` | `-f` | Freeze mode: `-1`=none, `0`=row+col, `1`=row only, `2`=col only |
| `--transpose` | `-t` | Transpose rows and columns |
| `--strict` | | Strict mode: fail on missing/inconsistent data |
| `--async` | | Progressive rendering while loading (default: `true`) |
| `--help` | `-h` | Show help |
| `--version` | `-v` | Show version |

**Examples:**

```bash
# Use custom delimiter
tv data.txt -s ","

# View only columns 1, 3, and 5
tv data.csv --columns 1,3,5

# Skip lines starting with "#"
tv data.txt --skip-prefix "#"

# Disable header freezing
tv data.csv -f -1

# Transpose data (swap rows and columns)
tv data.csv -t

# Disable async loading for slow systems
tv large.csv --async=false
```

## Key Bindings

tv uses vim-inspired keybindings for intuitive navigation.

### Navigation

| Key | Action |
|-----|--------|
| `h` / `←` | Move left |
| `l` / `→` | Move right |
| `j` / `↓` | Move down |
| `k` / `↑` | Move up |
| `w` | Jump to next column |
| `b` | Jump to previous column |
| `gg` | Go to first row |
| `G` | Go to last row |
| `0` | Go to first column |
| `$` | Go to last column |
| `Ctrl-d` | Page down (half page) |
| `Ctrl-u` | Page up (half page) |

### Operations

| Key | Action |
|-----|--------|
| `/` | Search |
| `n` | Next search result |
| `N` | Previous search result |
| `Ctrl-/` | Clear search |
| `f` | Filter by column value |
| `r` | Reset/clear filter |
| `s` | Sort ascending |
| `S` | Sort descending |
| `t` | Toggle column type (String → Number → Date) |
| `W` | Toggle text wrapping |
| `i` | Show column statistics |
| `?` | Show help |
| `Esc` | Close dialogs / clear search |
| `q` | Quit |

## Features in Detail

### Progressive Loading

### Progressive Loading

Start viewing large files instantly without waiting for them to fully load. The UI appears immediately and updates as data streams in.

```bash
# Default behavior - UI appears instantly
tv huge_dataset.csv

# Disable if you prefer traditional loading
tv huge_dataset.csv --async=false
```

**Progress indicators:**
- Files show percentage: `Loading... 45.2%` → `Loaded 1,000,000 rows`
- Pipes show row count: `Loading... 5,234 rows` → `Loaded 10,000 rows`
- Updates at 50 FPS for smooth rendering

### Data Types and Sorting

tv automatically detects column types and provides intelligent sorting.

**Type Detection:** When loading data, tv analyzes each column to determine if it contains strings, numbers, or dates using a 90% confidence threshold.

**Manual Type Toggle:** Press `t` to cycle through types for the current column:
- String → Number → Date → String

**Sorting Behavior:**
- **Strings:** Alphabetical order
- **Numbers:** Numeric order (supports integers, floats, scientific notation, thousands separators)
- **Dates:** Chronological order (supports ISO-8601, US format, EU format, and more)

**View Statistics:** Press `i` to open a comprehensive statistics dialog showing:
- **For numeric columns:** count, min, max, range, sum, mean, median, mode, standard deviation, variance, quartiles (Q1, Q2, Q3), and IQR
- **For string columns:** total values, unique values, frequency distribution with percentages for each value

### Search

### Search

Find text anywhere in your table with full highlighting support.

**How to search:**

1. Press `/` to open the search dialog
2. Type your search term (case-insensitive)
3. Press Enter to execute
4. Navigate results with `n` (next) and `N` (previous)
5. Press `Ctrl-/` to clear highlighting

**Visual feedback:**
- Current match: bright cyan highlight
- Other matches: gray highlight
- Footer shows position: `Match 3/12`

**Example:**
```
/ → type "error" → Enter → n → n → N → Ctrl-/
```

### Column Filter

Show only rows where a specific column matches your criteria.

**How to filter:**

1. Navigate to the column you want to filter
2. Press `f` to open the filter dialog
3. Type filter text (case-insensitive, partial match)
4. Press Enter to apply
5. Press `r` to reset and show all rows

**Features:**
- Partial matching: "act" matches "active", "action", "react"
- Header always visible
- Footer shows filtered row count

**Example:**
```
Navigate to "Status" column → f → type "pending" → Enter
```

Now only rows where Status contains "pending" are displayed.

### Text Wrapping

Handle long cell content without horizontal scrolling.

**How to wrap:**
- Press `W` (capital W) on any column to toggle wrapping
- First press: wraps at 25 characters with smart word breaks
- Second press: unwraps back to single line

**Smart wrapping:**
- Breaks at spaces and hyphens when possible
- Each column can be wrapped independently
- Useful for comments, descriptions, URLs

## Advanced Examples

### Biological Data Formats

tv handles common bioinformatics file formats with comment/header prefixes.

**VCF files:**
```bash
# Skip VCF metadata lines
tv sample.vcf --skip-prefix "##"
tv sample.vcf.gz --skip-prefix "##"
```

**QIIME OTU tables:**
```bash
tv otu_table.txt --skip-prefix "# "
```

**MAF (Mutation Annotation Format):**
```bash
tv mutations.maf --skip-prefix "#"
```

**Interval lists:**
```bash
# Skip SAM header lines
tv intervals.interval_list --skip-prefix "@HD","@SQ"

# Or skip all @ lines
tv intervals.interval_list --skip-prefix "@"
```

**BED files with headers:**
```bash
tv peaks.bed --skip-prefix "track","browser"
```

### General Examples

**Large log files:**
```bash
# View first 1000 lines only
tv app.log -n 1000

# Skip timestamp lines
tv app.log --skip-prefix "2024"
```

**CSV with specific columns:**
```bash
# Show only columns 1, 3, and 5
tv data.csv --columns 1,3,5

# Hide sensitive columns 2 and 4
tv data.csv --hide-columns 2,4
```

**Pipeline integration:**
```bash
# View process list
ps aux | tv

# View git log as table
git log --pretty=format:"%h,%an,%ar,%s" | tv -s ","

# Parse JSON with jq, view as table
cat data.json | jq -r '.[] | [.id, .name, .value] | @csv' | tv
```

**Custom delimiters:**
```bash
# Pipe-separated
tv data.txt -s "|"

# Semicolon-separated
tv data.txt -s ";"

# Multiple spaces
tv data.txt -s "  "
```

---

## Tips and Tricks

- **Large files?** Let async loading work its magic - the UI appears instantly
- **Can't find data?** Use `/` to search across all cells
- **Too many columns?** Use `--columns` to show only what you need
- **Long text?** Press `W` to wrap the current column
- **Wrong sort order?** Press `t` to change the column type, then `s` to re-sort
- **Need stats?** Press `i` for comprehensive statistics including mean, median, quartiles, std dev, and frequency distributions

## License

Apache License - see [LICENSE](LICENSE) file for details.
