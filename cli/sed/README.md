# ✂️ sed — Stream Editor Reference

> **Transform text streams with surgical precision.** Find and replace, delete, insert, and manipulate text — in files or pipelines — using powerful pattern matching.

---

## 📑 Table of Contents

- [Concepts](#-concepts)
- [Substitution](#-substitution)
- [Deletion](#-deletion)
- [Insertion & Appending](#-insertion--appending)
- [In-Place Editing](#-in-place-editing)
- [Address Ranges](#-address-ranges)
- [Regular Expressions in sed](#-regular-expressions-in-sed)
- [Multiple Commands](#-multiple-commands)
- [Hold Space & Advanced](#-hold-space--advanced)
- [Practical Examples](#-practical-examples)
- [Troubleshooting](#-troubleshooting)
- [Related Topics](#-related-topics)

---

## 💡 Concepts

`sed` (Stream Editor) reads input line by line, applies editing commands, and writes to stdout. It doesn't modify files by default — use `-i` for in-place editing.

```
sed [OPTIONS] 'COMMAND' [FILE...]
```

### Command Structure

```
[address]command[options]
```

- **Address**: Which lines to operate on (line number, range, or pattern)
- **Command**: What to do (`s` substitute, `d` delete, `i` insert, `a` append, `p` print, etc.)

### Common Flags

| Flag | Description |
|------|-------------|
| `-i` | In-place editing (modifies file directly) |
| `-i.bak` | In-place with backup (creates `.bak` file) |
| `-n` | Suppress default output (only print with `p`) |
| `-e` | Multiple commands |
| `-E` / `-r` | Extended regular expressions |
| `-f` | Read commands from file |

---

## 🔄 Substitution

The most common sed operation: `s/pattern/replacement/flags`

```bash
# Basic substitution (first occurrence per line)
sed 's/old/new/' file.txt

# Replace ALL occurrences per line (global)
sed 's/old/new/g' file.txt

# Replace Nth occurrence
sed 's/old/new/2' file.txt          # 2nd occurrence only
sed 's/old/new/3g' file.txt         # 3rd and all subsequent

# Case-insensitive substitution
sed 's/error/WARNING/gI' file.txt

# Different delimiters (useful with paths)
sed 's|/usr/local/bin|/opt/bin|g' file.txt
sed 's#http://#https://#g' file.txt
sed 's@old@new@g' file.txt

# Print only changed lines
sed -n 's/error/ERROR/gp' file.txt

# Substitution with captured groups
sed 's/\(Hello\) \(World\)/\2 \1/' file.txt
# "Hello World" → "World Hello"

# Extended regex groups (no backslash needed)
sed -E 's/(Hello) (World)/\2 \1/' file.txt

# Add prefix/suffix
sed 's/^/PREFIX: /' file.txt         # Add prefix to every line
sed 's/$/ SUFFIX/' file.txt          # Add suffix to every line

# Remove leading/trailing whitespace
sed 's/^[[:space:]]*//' file.txt     # Leading
sed 's/[[:space:]]*$//' file.txt     # Trailing
sed 's/^[[:space:]]*//;s/[[:space:]]*$//' file.txt  # Both
```

**Example:**
```
$ echo "Hello World 2024" | sed 's/[0-9]\+/YEAR/'
Hello World YEAR

$ echo "foo bar baz" | sed 's/\b\w/\U&/g'    # (GNU sed) Capitalize first letter
Foo Bar Baz
```

### Case Conversion (GNU sed)

```bash
# Uppercase
sed 's/[a-z]/\U&/g' file.txt

# Lowercase
sed 's/[A-Z]/\L&/g' file.txt

# Capitalize first letter of each word
sed 's/\b\w/\U&/g' file.txt

# Uppercase captured group
sed -E 's/(error|warning)/\U\1/g' file.txt
```

---

## 🗑️ Deletion

```bash
# Delete specific line
sed '5d' file.txt                   # Delete line 5

# Delete line range
sed '10,20d' file.txt               # Delete lines 10-20

# Delete first line
sed '1d' file.txt

# Delete last line
sed '$d' file.txt

# Delete empty lines
sed '/^$/d' file.txt

# Delete lines matching pattern
sed '/^#/d' file.txt                # Delete comment lines
sed '/DEBUG/d' app.log              # Delete debug lines

# Delete lines NOT matching pattern
sed '/ERROR/!d' app.log             # Keep only ERROR lines (same as grep)

# Delete from pattern to end of file
sed '/END_MARKER/,$d' file.txt

# Delete from start to pattern
sed '1,/START_MARKER/d' file.txt
```

---

## ➕ Insertion & Appending

```bash
# Insert BEFORE line (i command)
sed '1i\# This is a header' file.txt          # Before first line
sed '5i\New line here' file.txt                # Before line 5
sed '/pattern/i\Inserted before match' file.txt  # Before matching lines

# Append AFTER line (a command)
sed '1a\Line after first' file.txt             # After first line
sed '$a\# End of file' file.txt                # After last line
sed '/pattern/a\Appended after match' file.txt # After matching lines

# Change (replace) entire line (c command)
sed '3c\Completely new line 3' file.txt
sed '/old_text/c\Replacement line' file.txt

# Read file and insert after pattern (r command)
sed '/MARKER/r header.txt' file.txt
```

---

## 📝 In-Place Editing

```bash
# Edit file directly (GNU sed)
sed -i 's/old/new/g' file.txt

# Edit with backup
sed -i.bak 's/old/new/g' file.txt    # Creates file.txt.bak

# macOS/BSD sed requires backup extension
sed -i '' 's/old/new/g' file.txt      # No backup on macOS
sed -i '.bak' 's/old/new/g' file.txt  # With backup on macOS

# Edit multiple files
sed -i 's/old/new/g' *.conf
sed -i 's/old/new/g' file1.txt file2.txt file3.txt

# Combine with find for bulk operations
find . -name "*.py" -exec sed -i 's/old_func/new_func/g' {} +
find . -name "*.yml" -exec sed -i 's/apiVersion: v1beta1/apiVersion: v1/g' {} +
```

> **Warning:** Always test sed commands without `-i` first, or use `-i.bak` to create backups. In-place editing is irreversible.

---

## 📍 Address Ranges

Addresses specify which lines a command applies to.

```bash
# Line number
sed '5s/old/new/' file.txt          # Only line 5

# Line range
sed '10,20s/old/new/' file.txt      # Lines 10-20

# Pattern address
sed '/ERROR/s/old/new/' file.txt    # Lines matching "ERROR"

# Pattern range (between two patterns, inclusive)
sed '/START/,/END/s/old/new/' file.txt

# From pattern to end
sed '/MARKER/,$s/old/new/' file.txt

# First line to pattern
sed '1,/MARKER/s/old/new/' file.txt

# Every Nth line (step)
sed '0~2s/old/new/' file.txt        # Every 2nd line (GNU sed)
sed '1~3s/old/new/' file.txt        # Every 3rd line starting from 1

# Last line
sed '$s/old/new/' file.txt

# Negation (lines NOT matching)
sed '/^#/!s/old/new/g' file.txt     # Non-comment lines only
```

---

## 📐 Regular Expressions in sed

### BRE (Default)

```bash
sed 's/[0-9]\{3\}-[0-9]\{4\}/REDACTED/' file.txt   # Phone numbers
sed 's/\(hello\)/[\1]/' file.txt                     # Capture group
```

### ERE (Extended — with `-E` or `-r`)

```bash
sed -E 's/[0-9]{3}-[0-9]{4}/REDACTED/' file.txt     # No backslash escaping
sed -E 's/(hello) (world)/\2 \1/' file.txt           # Capture groups
sed -E 's/https?:\/\///' file.txt                    # Optional character
sed -E 's/(error|warning|critical)/[\1]/gi' file.txt # Alternation
```

### Special Characters

| In Pattern | Matches |
|-----------|---------|
| `\n` | Newline (in some contexts) |
| `\t` | Tab (GNU sed) |
| `\w` | Word character `[a-zA-Z0-9_]` (GNU sed) |
| `\b` | Word boundary (GNU sed) |
| `&` | Entire matched text (in replacement) |
| `\1`-`\9` | Captured group (in replacement) |

```bash
# Use & to reference entire match
sed 's/[0-9]\+/(&)/' file.txt       # Wrap numbers in parentheses
# "Score: 95" → "Score: (95)"

# Capture and rearrange
sed -E 's/([A-Z]+),([A-Z]+)/\2 \1/' names.txt
# "SMITH,JOHN" → "JOHN SMITH"
```

---

## 🔗 Multiple Commands

```bash
# Semicolon separation
sed 's/old/new/g; s/foo/bar/g' file.txt

# -e flag (multiple expressions)
sed -e 's/old/new/g' -e 's/foo/bar/g' file.txt

# Multi-line (backslash continuation)
sed -e 's/old/new/g' \
    -e '/^#/d' \
    -e 's/foo/bar/g' \
    file.txt

# From a script file
cat > commands.sed << 'EOF'
s/old/new/g
/^#/d
s/foo/bar/g
EOF
sed -f commands.sed file.txt

# Curly braces (apply multiple commands to same address)
sed '/ERROR/ { s/old/new/; s/foo/bar/; }' file.txt
```

---

## 🧠 Hold Space & Advanced

sed has two buffers:
- **Pattern space**: Current line being processed
- **Hold space**: Auxiliary buffer for multi-line operations

| Command | Description |
|---------|-------------|
| `h` | Copy pattern space to hold space |
| `H` | Append pattern space to hold space |
| `g` | Copy hold space to pattern space |
| `G` | Append hold space to pattern space |
| `x` | Exchange pattern and hold space |
| `N` | Append next line to pattern space |
| `P` | Print first line of pattern space |
| `D` | Delete first line of pattern space |

```bash
# Reverse line order (like tac)
sed -n '1!G;h;$p' file.txt

# Double-space a file
sed 'G' file.txt

# Remove double-spacing
sed 'n;d' file.txt

# Print lines between two patterns (exclusive)
sed -n '/START/,/END/{/START/!{/END/!p}}' file.txt

# Join continuation lines (lines ending with \)
sed ':a; /\\$/ { N; s/\\\n//; ba; }' file.txt
```

---

## 💼 Practical Examples

### Config File Editing

```bash
# Update a setting
sed -i 's/^port=.*/port=8080/' config.ini

# Uncomment a line
sed -i 's/^#\(listen_address\)/\1/' config.yaml

# Comment out a line
sed -i 's/^\(dangerous_setting\)/#\1/' config.yaml

# Change value between quotes
sed -i 's/"database_host": ".*"/"database_host": "newhost.example.com"/' config.json

# Update YAML value
sed -i 's/replicas: [0-9]*/replicas: 5/' deployment.yaml
```

### Log Processing

```bash
# Remove ANSI color codes
sed 's/\x1b\[[0-9;]*m//g' colored-output.log

# Extract specific field from structured log
sed -n 's/.*timestamp="\([^"]*\)".*/\1/p' app.log

# Redact sensitive data
sed -E 's/(password|secret|token)=[^ ]*/\1=REDACTED/g' app.log

# Remove duplicate empty lines
sed '/^$/N;/^\n$/d' file.txt

# Prefix each line with timestamp
sed "s/^/$(date '+%Y-%m-%d %H:%M:%S') /" input.txt
```

### Bulk File Changes

```bash
# Rename imports across project
find . -name "*.py" -exec sed -i 's/from old_module/from new_module/g' {} +

# Update API version in Kubernetes manifests
find . -name "*.yaml" -exec sed -i 's/apiVersion: apps\/v1beta1/apiVersion: apps\/v1/g' {} +

# Fix line endings (Windows → Unix)
sed -i 's/\r$//' file.txt

# Add shebang to all .sh files missing it
find . -name "*.sh" -exec sed -i '1{/^#!/!i\#!/bin/bash\n}' {} +

# Replace string in specific file sections only
sed '/\[production\]/,/\[/s/host=.*/host=prod-db.example.com/' config.ini
```

### Text Transformation

```bash
# Number lines (similar to nl)
sed = file.txt | sed 'N;s/\n/\t/'

# Remove HTML tags
sed 's/<[^>]*>//g' page.html

# Extract URLs
sed -n 's/.*\(https\?:\/\/[^ "]*\).*/\1/p' file.txt

# Convert tabs to spaces
sed 's/\t/    /g' file.txt

# Wrap long lines at 80 characters
sed 's/.\{80\}/&\n/g' file.txt
```

---

## 🐛 Troubleshooting

| Problem | Solution |
|---------|----------|
| "unterminated `s' command" | Missing closing delimiter — check `s/old/new/` |
| Escaping `/` in paths | Use different delimiter: `s|/old/path|/new/path|g` |
| `-i` doesn't work on macOS | Use `sed -i '' 's/...'` (empty backup extension) |
| Regex not matching | Check BRE vs ERE — try `-E` flag |
| Newlines not matching | Use `N` to load multiple lines, or `tr` for simple cases |
| Special chars in replacement | Escape `&`, `\`, `/` with backslash |
| "Permission denied" on `-i` | Check file permissions, ensure file is writable |
| Pattern includes `!` in bash | Use single quotes: `sed 's/pattern/replace/'` |

### Testing Before Applying

```bash
# Always test without -i first
sed 's/old/new/g' file.txt          # Preview changes

# Test with diff
sed 's/old/new/g' file.txt | diff file.txt -

# Create backup
sed -i.bak 's/old/new/g' file.txt
# If something went wrong:
mv file.txt.bak file.txt
```

---

## 🔗 Related Topics

- [🧰 CLI Command Center](../) — All CLI tool references
- [grep](../grep/) — Find patterns (often piped to sed)
- [awk](../awk/) — Column-based processing
- [Linux](../linux/) — Core Linux commands
- [⚙️ Engineering](../../engineering/) — Automation and scripting

---

> **Tip:** When building complex sed commands, start simple and add complexity one step at a time. Test each step before combining.
