# 🔎 grep — Text Search Reference

> **Find patterns in files and streams.** The essential tool for searching logs, code, configuration files, and command output using regular expressions.

---

## 📑 Table of Contents

- [Concepts](#-concepts)
- [Basic Usage](#-basic-usage)
- [Common Flags](#-common-flags)
- [Regular Expressions](#-regular-expressions)
- [Extended Grep (ERE)](#-extended-grep-ere)
- [Perl-Compatible Regex (PCRE)](#-perl-compatible-regex-pcre)
- [Context Lines](#-context-lines)
- [Recursive Search](#-recursive-search)
- [Practical Patterns](#-practical-patterns)
- [Performance Tips](#-performance-tips)
- [Troubleshooting](#-troubleshooting)
- [Related Topics](#-related-topics)

---

## 💡 Concepts

`grep` (Global Regular Expression Print) searches input files or stdin for lines matching a pattern and prints them.

```
grep [OPTIONS] PATTERN [FILE...]
```

### grep Variants

| Command | Description |
|---------|-------------|
| `grep` | Basic Regular Expressions (BRE) |
| `grep -E` / `egrep` | Extended Regular Expressions (ERE) |
| `grep -P` | Perl-Compatible Regular Expressions (PCRE) |
| `grep -F` / `fgrep` | Fixed strings (no regex, fastest) |
| `grep -r` | Recursive search |

---

## 🔍 Basic Usage

```bash
# Search in file
grep "error" app.log

# Search in multiple files
grep "TODO" *.py

# Search from stdin (pipe)
cat app.log | grep "error"
ps aux | grep nginx
kubectl get pods | grep Running

# Case-insensitive
grep -i "error" app.log

# Inverted match (lines NOT matching)
grep -v "debug" app.log

# Count matches
grep -c "error" app.log

# Show line numbers
grep -n "error" app.log

# Show only filenames with matches
grep -l "TODO" src/*.py

# Show only filenames WITHOUT matches
grep -L "TODO" src/*.py

# Show only the matching part (not full line)
grep -o "[0-9]\{1,3\}\.[0-9]\{1,3\}\.[0-9]\{1,3\}\.[0-9]\{1,3\}" access.log
```

**Example output:**
```
$ grep -n "ERROR" app.log
15:2024-08-15 10:23:45 ERROR Failed to connect to database
89:2024-08-15 10:45:12 ERROR Timeout waiting for response
234:2024-08-15 11:02:33 ERROR Out of memory
```

---

## 🚩 Common Flags

| Flag | Long Form | Description |
|------|-----------|-------------|
| `-i` | `--ignore-case` | Case-insensitive matching |
| `-v` | `--invert-match` | Print non-matching lines |
| `-n` | `--line-number` | Print line numbers |
| `-c` | `--count` | Count matching lines per file |
| `-l` | `--files-with-matches` | Print only filenames with matches |
| `-L` | `--files-without-match` | Print only filenames without matches |
| `-o` | `--only-matching` | Print only matched portion |
| `-w` | `--word-regexp` | Match whole words only |
| `-x` | `--line-regexp` | Match whole lines only |
| `-r` | `--recursive` | Search directories recursively |
| `-R` | `--dereference-recursive` | Recursive, following symlinks |
| `-h` | `--no-filename` | Suppress filename prefix |
| `-H` | `--with-filename` | Always show filename |
| `-m N` | `--max-count=N` | Stop after N matches |
| `-q` | `--quiet` | Quiet (exit status only, no output) |
| `-s` | `--no-messages` | Suppress error messages |
| `-A N` | `--after-context=N` | Print N lines after match |
| `-B N` | `--before-context=N` | Print N lines before match |
| `-C N` | `--context=N` | Print N lines before and after |
| `-E` | `--extended-regexp` | Extended regular expressions |
| `-P` | `--perl-regexp` | Perl-compatible regex |
| `-F` | `--fixed-strings` | Literal string (no regex) |
| `--include` | | Search only matching filenames |
| `--exclude` | | Skip matching filenames |
| `--exclude-dir` | | Skip matching directories |
| `--color` | | Highlight matches |

---

## 📐 Regular Expressions

### Basic Regular Expressions (BRE)

| Pattern | Matches |
|---------|---------|
| `.` | Any single character |
| `*` | Zero or more of preceding |
| `^` | Start of line |
| `$` | End of line |
| `[abc]` | Character class (a, b, or c) |
| `[^abc]` | Negated class (not a, b, c) |
| `[a-z]` | Character range |
| `\b` | Word boundary |
| `\<` | Start of word |
| `\>` | End of word |
| `\{n\}` | Exactly n repetitions |
| `\{n,m\}` | Between n and m repetitions |
| `\( \)` | Group (capturing) |
| `\1` | Back-reference |

```bash
# Lines starting with "Error"
grep "^Error" app.log

# Lines ending with a number
grep "[0-9]$" data.txt

# Lines with IP-like patterns
grep "[0-9]\{1,3\}\.[0-9]\{1,3\}\.[0-9]\{1,3\}\.[0-9]\{1,3\}" access.log

# Whole word match
grep -w "error" app.log       # Won't match "errors" or "errored"

# Lines with exactly 3 digits
grep "^[0-9]\{3\}$" data.txt
```

---

## 🔤 Extended Grep (ERE)

Extended regex removes the need for backslash escaping and adds more operators. Use `grep -E` or `egrep`.

| Pattern | Matches |
|---------|---------|
| `+` | One or more of preceding |
| `?` | Zero or one of preceding |
| `{n}` | Exactly n repetitions |
| `{n,m}` | Between n and m repetitions |
| `\|` becomes `|` | Alternation (OR) |
| `( )` | Grouping (no backslash needed) |

```bash
# Match "error" or "warning"
grep -E "error|warning" app.log

# Match repeated patterns
grep -E "([0-9]{1,3}\.){3}[0-9]{1,3}" access.log    # IP addresses

# Optional prefix
grep -E "https?://" urls.txt        # http:// or https://

# One or more digits
grep -E "[0-9]+" data.txt

# Email-like pattern
grep -E "[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}" contacts.txt

# HTTP status codes (4xx and 5xx)
grep -E "HTTP/[0-9.]+ [45][0-9]{2}" access.log
```

---

## 🐪 Perl-Compatible Regex (PCRE)

Use `grep -P` for advanced regex features.

```bash
# Lookahead: match "error" followed by a number
grep -P "error(?=.*[0-9])" app.log

# Lookbehind: match numbers preceded by "code:"
grep -P "(?<=code:)\d+" app.log

# Non-greedy match
grep -oP '"name":".*?"' data.json

# Named groups
grep -oP '(?P<ip>\d+\.\d+\.\d+\.\d+)' access.log

# \d, \w, \s shortcuts
grep -P "\d{4}-\d{2}-\d{2}" dates.txt       # Date pattern
grep -P "\bfoo\w+" code.py                   # Words starting with "foo"
```

---

## 📖 Context Lines

Show surrounding lines for context when searching logs or code.

```bash
# N lines AFTER each match
grep -A 3 "Exception" app.log

# N lines BEFORE each match
grep -B 2 "FATAL" app.log

# N lines BEFORE and AFTER (context)
grep -C 5 "segfault" kern.log
```

**Example output:**
```
$ grep -B1 -A2 "ERROR" app.log
2024-08-15 10:23:44 INFO  Processing request #1234
2024-08-15 10:23:45 ERROR Failed to connect to database
2024-08-15 10:23:45 ERROR   Connection refused: 10.0.0.5:5432
2024-08-15 10:23:46 INFO  Retry attempt 1
--
2024-08-15 10:45:11 INFO  Sending notification
2024-08-15 10:45:12 ERROR Timeout waiting for response
2024-08-15 10:45:12 WARN  Notification delivery failed
2024-08-15 10:45:13 INFO  Added to retry queue
```

---

## 📁 Recursive Search

```bash
# Search in directory tree
grep -r "TODO" src/

# With line numbers
grep -rn "TODO" src/

# Filter by file extension
grep -rn "import" --include="*.py" src/
grep -rn "require" --include="*.{js,ts}" src/

# Exclude directories
grep -rn "TODO" --exclude-dir=node_modules --exclude-dir=.git src/
grep -rn "password" --exclude-dir={node_modules,.git,vendor} .

# Exclude file types
grep -rn "error" --exclude="*.log" --exclude="*.bak" /var/
```

> **Tip:** For large codebases, consider using [`ripgrep`](https://github.com/BurntSushi/ripgrep) (`rg`) — it's significantly faster and respects `.gitignore` by default.

---

## 💼 Practical Patterns

### Searching Logs

```bash
# Errors in last hour (combine with other tools)
grep "$(date '+%Y-%m-%d %H')" app.log | grep -i "error"

# Count error types
grep -i "error" app.log | sort | uniq -c | sort -rn | head -10

# Extract timestamps of errors
grep -oP "\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2}" app.log | head

# Find slow queries (>1000ms)
grep -E "took [0-9]{4,}ms" db.log

# HTTP 5xx errors in access logs
grep -E '" [5][0-9]{2} ' access.log
```

### Searching Code

```bash
# Find function definitions
grep -rn "def " --include="*.py" src/
grep -rn "func " --include="*.go" src/
grep -rn "function " --include="*.js" src/

# Find TODO/FIXME/HACK comments
grep -rn "TODO\|FIXME\|HACK\|XXX" src/

# Find imports of a module
grep -rn "from requests import\|import requests" --include="*.py" src/

# Find hardcoded IPs
grep -rnE "[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}" --include="*.py" src/

# Find console.log (for cleanup)
grep -rn "console\.log" --include="*.{js,ts,tsx}" src/
```

### Searching Config Files

```bash
# Find non-comment, non-empty lines in config
grep -v "^#" /etc/nginx/nginx.conf | grep -v "^$"

# Find a setting across configs
grep -rn "max_connections" /etc/

# Check for specific port usage
grep -rn "8080" /etc/nginx/ /etc/apache2/
```

### Combining grep with Other Tools

```bash
# grep + awk: extract fields from matching lines
grep "ERROR" app.log | awk '{print $1, $2, $NF}'

# grep + sort + uniq: frequency analysis
grep -oP "error: \K[^$]+" app.log | sort | uniq -c | sort -rn

# grep + wc: count matches
grep -c "ERROR" app.log

# grep + xargs: act on matching files
grep -rl "deprecated_func" src/ | xargs sed -i 's/deprecated_func/new_func/g'

# Process check
ps aux | grep -v grep | grep nginx
# Better alternative:
pgrep -a nginx
```

---

## ⚡ Performance Tips

```bash
# Use fixed strings when no regex needed (much faster)
grep -F "exact string" largefile.txt

# Use --include to limit file scanning
grep -r "pattern" --include="*.log" /var/

# Limit matches per file
grep -m 1 "pattern" *.log          # Stop at first match per file

# Use -q for existence check (stops at first match)
if grep -q "error" app.log; then
    echo "Errors found!"
fi

# Set locale for speed with binary-safe search
LC_ALL=C grep "pattern" hugefile.txt

# Parallel grep with xargs
find . -name "*.log" -print0 | xargs -0 -P 4 grep -l "ERROR"
```

---

## 🐛 Troubleshooting

| Problem | Solution |
|---------|----------|
| "Binary file matches" | Use `-a` to treat as text, or `--binary-files=text` |
| No output but file has matches | Check encoding, try `LC_ALL=C grep` |
| Backslash issues | Use single quotes around pattern: `grep 'pattern'` |
| Special characters in pattern | Escape with `\` or use `-F` for fixed strings |
| Too many results | Add `-m N`, use more specific pattern, combine with `head` |
| Regex not matching | Check if using BRE vs ERE — try `-E` flag |
| "Argument list too long" | Use `find ... -exec grep` or pipe with `xargs` |

### Exit Codes

| Code | Meaning |
|------|---------|
| `0` | Match found |
| `1` | No match found |
| `2` | Error (bad pattern, unreadable file, etc.) |

---

## 🔗 Related Topics

- [🧰 CLI Command Center](../) — All CLI tool references
- [sed](../sed/) — Stream editing (often used after grep to modify matches)
- [awk](../awk/) — Column-based processing of grep results
- [Linux](../linux/) — Core Linux commands
- [jq](../jq/) — JSON processing (alternative to grep for JSON)
- [🐛 Troubleshooting](../../troubleshooting/) — Log analysis guides

---

> **Tip:** Chain `grep` with `sort | uniq -c | sort -rn` for quick frequency analysis of any repeated pattern in logs or data.
