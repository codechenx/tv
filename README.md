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
  - [Statistics and Visualization](#statistics-and-visualization)
  - [Search](#search)
  - [Column Filter](#column-filter)
  - [Text Wrapping](#text-wrapping)
- [Filter Operators Guide](FILTER_OPERATORS.md)
- [Advanced Examples](#advanced-examples)
  - [Biological Data Formats](#biological-data-formats)

## Features

tv brings spreadsheet-like functionality to your terminal with vim-inspired controls.

- **Spreadsheet interface** - Navigate and view tabular data with frozen headers
- **Smart parsing** - Automatically detects delimiters (CSV, TSV, custom separators)
- **Progressive loading** - Start viewing large files immediately while they load
- **Gzip support** - Read compressed files directly
- **Powerful search** - Find text across all cells with highlighting
- **Advanced filtering** - Filter rows with OR, AND, and ROR operators for complex queries
- **Flexible sorting** - Sort by any column with intelligent type detection
- **Text wrapping** - Wrap long cell content for better readability
- **Statistics & plots** - View column statistics with visual distribution charts
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
wget https://github.com/codechenx/tv/releases/download/v0.6.2/tv_0.6.3_linux_amd64.deb
sudo dpkg -i tv_*.deb
```

**CentOS/Fedora (.rpm)**

```bash
# Download from releases page
wget https://github.com/codechenx/tv/releases/download/v0.6.2/tv_0.6.3_linux_amd64.rpm
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

### Statistics and Visualization

Analyze your data with comprehensive statistics and modern ASCII plots.

**View Statistics:** Press `i` on any column to open an interactive statistics dialog that displays:

**For numeric columns:**
- Summary stats: count, min, max, range, sum
- Central tendency: mean, median, mode
- Dispersion: standard deviation, variance
- Quartiles: Q1, Q2, Q3, and IQR
- **Visual distribution:** Histogram plot showing data distribution

**For categorical/string columns:**
- Total values, unique values, missing/empty count
- Frequency distribution with percentages
- **Visual distribution:** Bar chart of top 15 most frequent values

The statistics dialog features a split-pane layout with numerical stats on the left and an ASCII graph visualization on the right, powered by `asciigraph` for modern, clean plots.

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

Show only rows where specific columns match your criteria. **Supports filtering on multiple columns simultaneously** with powerful OR, AND, and ROR operators for complex filtering.

**How to filter:**

1. Navigate to the column you want to filter
2. Press `f` to open the filter dialog
3. Type filter text (case-insensitive, partial match)
4. Use operators for complex queries (must be UPPERCASE with spaces)
5. Press Enter to apply
6. **Repeat on other columns to add more filters**
7. Press `r` on a filtered column to remove that column's filter

**Multi-Column Filtering:**

- Apply filters to **multiple columns** by pressing `f` on each column and entering criteria
- All filters are combined with AND logic (rows must match all active filters)
- Each column can have different filter criteria including operators
- Press `f` on a filtered column to edit or remove its filter (empty query removes the filter)
- The footer shows how many filters are active
- Filtered column headers display 🔎 icons with orange background

**Operators:**

| Operator | Syntax | Behavior | Example |
|----------|--------|----------|---------|
| **Simple** | `term` | Matches cells containing the term | `active` → matches "Active", "Inactive" |
| **OR** | `term1 OR term2` | Same cell contains either term | `error OR warning` → cell has "error" or "warning" |
| **AND** | `term1 AND term2` | Same cell contains both terms | `user AND admin` → cell has both words |
| **ROR** | `term1 ROR term2` | Keeps all rows matching any term | `high ROR critical` → rows with "high" + rows with "critical" |
| **>** | `>value` | Numeric: greater than (number columns only) | `>30` → values greater than 30 |
| **<** | `<value` | Numeric: less than (number columns only) | `<50` → values less than 50 |
| **>=** | `>=value` | Numeric: greater than or equal (number columns only) | `>=100` → values 100 or more |
| **<=** | `<=value` | Numeric: less than or equal (number columns only) | `<=75` → values 75 or less |

**Key Differences:**
- **OR vs ROR**: OR checks if a single cell contains either term. ROR combines rows where any cell matches any term (row-level union).
- **Numeric operators** (`>`, `<`, `>=`, `<=`): Only work on numeric and date columns (automatically detected). Perform numeric comparisons instead of text matching.
- All operators must be **UPPERCASE** and surrounded by spaces (except numeric operators)
- Search terms are case-insensitive (except numeric comparisons)
- All matching is partial (substring) for text, exact comparison for numeric operators
- **Visual indicator**: Filtered column headers show 🔎 icons and orange background

**Examples:**

```bash
# Simple filter - partial match
Navigate to "Status" column → f → type "pending" → Enter
# Result: Matches "Pending", "Pending Review", etc.

# OR filter - same cell matches either term
Navigate to "Status" column → f → type "active OR pending" → Enter
# Result: Rows where Status contains "active" OR "pending"

# AND filter - same cell must contain both terms
Navigate to "Description" column → f → type "user AND admin" → Enter
# Result: Rows where Description has both "user" AND "admin"

# ROR filter - combines separate result sets
Navigate to "Priority" column → f → type "high ROR critical" → Enter
# Result: All rows with "high" + all rows with "critical" (union)

# Numeric comparison - greater than
Navigate to "Age" column → f → type ">30" → Enter
# Result: Rows where Age is greater than 30

# Numeric comparison - less than or equal
Navigate to "Score" column → f → type "<=85" → Enter
# Result: Rows where Score is 85 or less

# Numeric comparison - greater than or equal
Navigate to "Salary" column → f → type ">=50000" → Enter
# Result: Rows where Salary is 50000 or more

# Multi-column filtering
Navigate to "City" column → f → type "New York" → Enter
Navigate to "Department" column → f → type "Engineering" → Enter
# Result: Rows where City="New York" AND Department contains "Engineering"

# Multi-column with numeric filter
Navigate to "Age" column → f → type ">25" → Enter
Navigate to "Score" column → f → type ">80" → Enter
# Result: Rows where Age > 25 AND Score > 80

# Edit existing filter
Navigate to filtered column → f → modify text → Enter
# Or enter empty text to remove that column's filter

# Remove a specific filter
Navigate to filtered column → Press r
# Result: That column's filter removed, other filters remain

# Remove all filters one by one
Navigate to each filtered column → Press r on each
# Or press f and enter empty text
```

**Use Cases:**
- Use **OR** when a single field can have alternative values
- Use **AND** when a single field must meet multiple criteria
- Use **ROR** when you want to combine different categories of results
- Use **numeric operators** (`>`, `<`, `>=`, `<=`) to filter by numeric ranges on number or date columns
- Use **multi-column filters** to narrow down data by multiple dimensions (e.g., location AND department AND status)

**Visual Feedback:**
- When a filter is active, the filtered column header displays 🔎 icons and an orange background
- A dedicated **filter strip** appears above the main footer **only when cursor is on the filtered column**
- Filter strip format: `🔎 Filter Active: [Column Name] = "query"  |  Press 'r' to clear`
- The strip automatically hides when you move to a different column
- Press `r` to clear the filter and return to normal view

For more details, see [FILTER_OPERATORS.md](FILTER_OPERATORS.md).

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
- **Complex filtering?** Use `OR` for alternatives, `AND` for requirements, `ROR` to combine results
- **Need insights?** Press `i` for comprehensive statistics with visual plots - histograms for numeric data, frequency charts for categorical data

## License

Apache License - see [LICENSE](LICENSE) file for details.
