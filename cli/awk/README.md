# 📊 awk — Text Processing Reference

> **The programmable text processor.** Extract columns, transform data, generate reports, and process structured text — from simple one-liners to complete data pipelines.

---

## 📑 Table of Contents

- [Concepts](#-concepts)
- [Basic Syntax](#-basic-syntax)
- [Fields and Records](#-fields-and-records)
- [Patterns and Actions](#-patterns-and-actions)
- [Built-in Variables](#-built-in-variables)
- [Operators & Expressions](#-operators--expressions)
- [Control Flow](#-control-flow)
- [String Functions](#-string-functions)
- [Math Functions](#-math-functions)
- [Arrays](#-arrays)
- [Output Formatting](#-output-formatting)
- [Practical Examples](#-practical-examples)
- [Troubleshooting](#-troubleshooting)
- [Related Topics](#-related-topics)

---

## 💡 Concepts

`awk` is a pattern-scanning and processing language. It reads input line by line, splits each line into fields, and applies pattern-action rules.

```
awk [OPTIONS] 'pattern { action }' [FILE...]
```

### Processing Model

```
For each line (record):
  1. Split line into fields ($1, $2, ..., $NF)
  2. For each pattern-action rule:
     - If pattern matches → execute action
  3. Read next line
```

### Common Flags

| Flag | Description |
|------|-------------|
| `-F` | Set field separator |
| `-v var=val` | Set variable before execution |
| `-f file` | Read program from file |

---

## 🔍 Basic Syntax

```bash
# Print entire line
awk '{print}' file.txt
awk '{print $0}' file.txt          # $0 = entire line

# Print specific fields
awk '{print $1}' file.txt          # First field
awk '{print $1, $3}' file.txt      # Fields 1 and 3 (space-separated)
awk '{print $NF}' file.txt         # Last field

# Custom output separator
awk '{print $1 ":" $2}' file.txt   # Colon-separated
awk -v OFS="\t" '{print $1, $2}' file.txt  # Tab-separated

# Print with text
awk '{print "Name:", $1, "Age:", $2}' file.txt

# From stdin
echo "hello world" | awk '{print $2}'
# world

ps aux | awk '{print $1, $11}'     # User and command
```

---

## 📐 Fields and Records

### Field Separators

```bash
# Default separator (whitespace)
awk '{print $1}' file.txt

# Custom separator
awk -F: '{print $1}' /etc/passwd           # Colon
awk -F, '{print $1, $3}' data.csv          # Comma
awk -F'\t' '{print $2}' data.tsv           # Tab
awk -F'[,;:]' '{print $1}' file.txt        # Multiple separators (regex)

# Set in BEGIN block
awk 'BEGIN{FS=":"} {print $1}' /etc/passwd

# Change output separator
awk -F: -v OFS="," '{print $1, $3, $6}' /etc/passwd
```

**Example:**
```
$ cat /etc/passwd | head -3 | awk -F: '{print $1, $3, $6}'
root 0 /root
daemon 1 /usr/sbin
bin 2 /bin

$ cat /etc/passwd | head -3 | awk -F: -v OFS="\t" '{print $1, $3, $6}'
root	0	/root
daemon	1	/usr/sbin
bin	2	/bin
```

### Record Separators

```bash
# Default record separator (newline)
# Change to paragraph mode (blank line separates records)
awk 'BEGIN{RS=""} {print NR, $0}' file.txt

# Custom record separator
awk 'BEGIN{RS=";"} {print NR, $0}' data.txt
```

### Modifying Fields

```bash
# Change a field value
awk -F: -v OFS=":" '{$7="/bin/zsh"; print}' /etc/passwd

# Add a new field
awk '{$(NF+1)="NEW"; print}' file.txt
```

---

## 🎯 Patterns and Actions

### Pattern Types

```bash
# No pattern (apply to all lines)
awk '{print $1}' file.txt

# Regular expression pattern
awk '/error/' file.txt                     # Lines containing "error"
awk '/^#/' file.txt                        # Lines starting with #
awk '!/^#/' file.txt                       # Lines NOT starting with #

# Field-specific regex
awk '$3 ~ /error/' file.txt               # Field 3 contains "error"
awk '$1 !~ /^test/' file.txt              # Field 1 doesn't start with "test"

# Comparison pattern
awk '$3 > 100' file.txt                   # Field 3 greater than 100
awk '$1 == "admin"' file.txt              # Field 1 equals "admin"
awk 'NR > 1' file.txt                     # Skip header (line > 1)

# Range pattern (from first match to second match)
awk '/START/,/END/' file.txt              # Print between START and END

# BEGIN and END blocks
awk 'BEGIN{print "Header"} {print} END{print "Footer"}' file.txt

# Compound patterns
awk '$3 > 50 && $4 == "active"' file.txt  # AND
awk '$3 > 90 || $4 == "critical"' file.txt # OR
```

### Multiple Patterns

```bash
awk '
    /ERROR/   { errors++ }
    /WARNING/ { warnings++ }
    END       { print "Errors:", errors, "Warnings:", warnings }
' app.log
```

---

## 📌 Built-in Variables

| Variable | Description |
|----------|-------------|
| `$0` | Entire current record (line) |
| `$1, $2, ...` | Individual fields |
| `$NF` | Last field |
| `NR` | Current record (line) number (global) |
| `NF` | Number of fields in current record |
| `FNR` | Record number in current file |
| `FS` | Input field separator (default: whitespace) |
| `OFS` | Output field separator (default: space) |
| `RS` | Input record separator (default: newline) |
| `ORS` | Output record separator (default: newline) |
| `FILENAME` | Current input filename |
| `ARGC` | Number of arguments |
| `ARGV` | Argument array |
| `ENVIRON` | Environment variable array |

```bash
# Print line numbers
awk '{print NR, $0}' file.txt

# Print number of fields per line
awk '{print NR, NF, $0}' file.txt

# Last field of each line
awk '{print $NF}' file.txt

# Second to last field
awk '{print $(NF-1)}' file.txt

# Total lines
awk 'END{print NR}' file.txt

# Multiple files — distinguish by FNR vs NR
awk 'FNR==1{print "--- File:", FILENAME, "---"} {print}' file1.txt file2.txt
```

**Example:**
```
$ echo -e "a b c\n1 2 3 4\nx y" | awk '{print "Line " NR ": " NF " fields, last=" $NF}'
Line 1: 3 fields, last=c
Line 2: 4 fields, last=4
Line 3: 2 fields, last=y
```

---

## ➕ Operators & Expressions

### Arithmetic

```bash
awk '{print $1 + $2}' file.txt            # Addition
awk '{print $1 * $2}' file.txt            # Multiplication
awk '{print $1 / $2}' file.txt            # Division
awk '{print $1 % $2}' file.txt            # Modulo
awk '{print $1 ^ 2}' file.txt             # Exponentiation
awk '{total += $1} END{print total}' file.txt  # Running total
```

### Assignment

```bash
awk '{count++} END{print count}' file.txt     # Increment
awk '{sum += $3} END{print "Avg:", sum/NR}' file.txt  # Average
```

### String Concatenation

```bash
# Strings concatenate by juxtaposition (no operator)
awk '{print $1 "-" $2}' file.txt
awk '{name = $1 " " $2; print name}' file.txt
```

### Ternary Operator

```bash
awk '{print ($3 > 80 ? "PASS" : "FAIL"), $0}' grades.txt
```

---

## 🔀 Control Flow

```bash
# If-else
awk '{
    if ($3 > 90) print $1, "A"
    else if ($3 > 80) print $1, "B"
    else if ($3 > 70) print $1, "C"
    else print $1, "F"
}' grades.txt

# While loop
awk '{
    i = 1
    while (i <= NF) {
        print $i
        i++
    }
}' file.txt

# For loop
awk '{
    for (i = 1; i <= NF; i++)
        print $i
}' file.txt

# For-in (iterate array)
awk '{
    words[$1]++
} END {
    for (w in words)
        print w, words[w]
}' file.txt

# Next (skip to next record)
awk '/^#/ {next} {print}' file.txt          # Skip comments

# Exit
awk 'NR > 100 {exit} {print}' file.txt     # Print first 100 lines
```

---

## 📝 String Functions

| Function | Description |
|----------|-------------|
| `length(s)` | String length |
| `substr(s, start, len)` | Substring |
| `index(s, target)` | Position of target in s (1-based, 0 if not found) |
| `split(s, arr, sep)` | Split string into array |
| `sub(regex, replacement, target)` | Replace first match |
| `gsub(regex, replacement, target)` | Replace all matches |
| `match(s, regex)` | Find regex match (sets RSTART, RLENGTH) |
| `sprintf(fmt, ...)` | Format string (like C sprintf) |
| `tolower(s)` | Lowercase |
| `toupper(s)` | Uppercase |

```bash
# String length
awk '{print length($0)}' file.txt

# Substring
awk '{print substr($0, 1, 10)}' file.txt          # First 10 chars

# Split
awk '{n = split($0, parts, ":"); print parts[1], n}' file.txt

# Replace
awk '{gsub(/old/, "new"); print}' file.txt
awk '{sub(/^[ \t]+/, ""); print}' file.txt         # Left trim

# Case conversion
awk '{print toupper($1), tolower($2)}' file.txt

# Match and extract
awk 'match($0, /[0-9]+/) {print substr($0, RSTART, RLENGTH)}' file.txt
```

---

## 🔢 Math Functions

| Function | Description |
|----------|-------------|
| `int(x)` | Truncate to integer |
| `sqrt(x)` | Square root |
| `sin(x)`, `cos(x)` | Trigonometry |
| `atan2(y, x)` | Arc tangent |
| `exp(x)` | Exponential (e^x) |
| `log(x)` | Natural logarithm |
| `rand()` | Random 0-1 |
| `srand(seed)` | Seed random generator |

```bash
# Round to 2 decimal places
awk '{printf "%.2f\n", $1}' numbers.txt

# Random sample (10% of lines)
awk 'BEGIN{srand()} rand() < 0.1' file.txt

# Absolute value
awk '{print ($1 < 0 ? -$1 : $1)}' numbers.txt
```

---

## 📦 Arrays

awk arrays are associative (like dictionaries/maps).

```bash
# Count occurrences
awk '{count[$1]++} END {for (k in count) print k, count[k]}' file.txt

# Sum by group
awk '{sum[$1] += $2} END {for (k in sum) print k, sum[k]}' sales.txt

# Check existence
awk '{if ($1 in seen) print "duplicate:", $0; seen[$1]=1}' file.txt

# Delete element
awk '{a[NR]=$0} END {delete a[5]; for (i=1; i<=NR; i++) if (i in a) print a[i]}' file.txt

# Multi-dimensional (simulated with SUBSEP)
awk '{data[$1,$2] += $3} END {
    for (key in data) {
        split(key, parts, SUBSEP)
        print parts[1], parts[2], data[key]
    }
}' file.txt

# Array length (gawk)
awk '{arr[$1]++} END {print length(arr), "unique values"}' file.txt

# Sort output (pipe to sort)
awk '{count[$1]++} END {for (k in count) print count[k], k}' file.txt | sort -rn
```

---

## 🖨️ Output Formatting

### printf

```bash
# Formatted output (like C printf)
awk '{printf "%-20s %10d %8.2f\n", $1, $2, $3}' data.txt

# Column alignment
awk 'BEGIN{printf "%-15s %10s %10s\n", "Name", "Age", "Score"; print "---"}
     {printf "%-15s %10d %10.1f\n", $1, $2, $3}' data.txt
```

**Format specifiers:**

| Specifier | Type |
|-----------|------|
| `%s` | String |
| `%d` | Integer |
| `%f` | Float |
| `%e` | Scientific notation |
| `%x` | Hexadecimal |
| `%o` | Octal |
| `%%` | Literal `%` |

**Modifiers:**
- `-` left-align
- `N` minimum width
- `.N` precision

```bash
# Table formatting example
awk -F: 'BEGIN{
    printf "%-20s %-6s %-30s\n", "User", "UID", "Home"
    printf "%-20s %-6s %-30s\n", "----", "---", "----"
}
$3 >= 1000 {
    printf "%-20s %-6d %-30s\n", $1, $3, $6
}' /etc/passwd
```

### Output Redirection

```bash
# Write to file from within awk
awk '{print $1 > "output.txt"}' file.txt

# Append to file
awk '{print $1 >> "output.txt"}' file.txt

# Pipe to command
awk '{print $1 | "sort -u"}' file.txt
```

---

## 💼 Practical Examples

### CSV Processing

```bash
# Print specific columns from CSV
awk -F, '{print $1, $3}' data.csv

# Skip header
awk -F, 'NR > 1 {print $1, $3}' data.csv

# Add header and format
awk -F, 'BEGIN{OFS=","} 
    NR==1 {print $0, "total"}
    NR>1  {print $0, $2+$3+$4}
' data.csv

# Convert CSV to TSV
awk -F, -v OFS="\t" '{$1=$1; print}' data.csv

# Filter rows by column value
awk -F, '$3 > 1000 && $5 == "active"' sales.csv

# CSV with quoted fields (simple approach)
awk -v FPAT='([^,]*)|("[^"]*")' '{print $1, $3}' data.csv
```

### Log Analysis

```bash
# Count requests per IP (Apache/Nginx access log)
awk '{print $1}' access.log | sort | uniq -c | sort -rn | head -10

# Count HTTP status codes
awk '{print $9}' access.log | sort | uniq -c | sort -rn

# Average response time (if logged as last field)
awk '{sum += $NF; count++} END {print "Avg:", sum/count, "ms"}' access.log

# Requests per hour
awk -F'[: ]' '{print $5}' access.log | sort | uniq -c

# Slow requests (>1 second)
awk '$NF > 1000 {print}' access.log

# Error rate per minute
awk '/ERROR/ {
    split($2, t, ":")
    minute = t[1] ":" t[2]
    errors[minute]++
} END {
    for (m in errors) print m, errors[m]
}' app.log | sort

# Top 10 most requested URLs
awk '{print $7}' access.log | sort | uniq -c | sort -rn | head -10

# Bandwidth per IP
awk '{bytes[$1] += $10} END {for (ip in bytes) printf "%15s %10.2f MB\n", ip, bytes[ip]/1048576}' access.log | sort -k2 -rn
```

### Report Generation

```bash
# Summary statistics
awk '
BEGIN {
    min = 999999; max = 0
}
NR > 1 {
    sum += $2
    count++
    if ($2 > max) max = $2
    if ($2 < min) min = $2
}
END {
    printf "Count:   %d\n", count
    printf "Sum:     %.2f\n", sum
    printf "Average: %.2f\n", sum/count
    printf "Min:     %.2f\n", min
    printf "Max:     %.2f\n", max
}' data.txt

# Pivot table
awk -F, 'NR > 1 {
    revenue[$1] += $3
    count[$1]++
} END {
    printf "%-20s %10s %10s %10s\n", "Category", "Revenue", "Count", "Average"
    printf "%-20s %10s %10s %10s\n", "--------", "-------", "-----", "-------"
    for (cat in revenue)
        printf "%-20s %10.2f %10d %10.2f\n", cat, revenue[cat], count[cat], revenue[cat]/count[cat]
}' sales.csv

# Generate HTML table
awk 'BEGIN{
    print "<table border=\"1\">"
    print "<tr><th>Name</th><th>Score</th><th>Grade</th></tr>"
}
{
    grade = ($2 >= 90 ? "A" : ($2 >= 80 ? "B" : ($2 >= 70 ? "C" : "F")))
    printf "<tr><td>%s</td><td>%d</td><td>%s</td></tr>\n", $1, $2, grade
}
END{print "</table>"}' grades.txt
```

### System Administration

```bash
# Disk usage by filesystem (parse df output)
df -h | awk 'NR>1 {print $5, $6}' | sort -rn

# Find large processes
ps aux | awk 'NR>1 && $4 > 5 {printf "%-10s %6s%% %6s%% %s\n", $1, $3, $4, $11}'

# Parse /etc/passwd
awk -F: '$3 >= 1000 {print $1, $6}' /etc/passwd    # Regular users

# Connection count by state
ss -t | awk 'NR>1 {print $1}' | sort | uniq -c | sort -rn

# Monitor file growth
ls -l file.log | awk '{print strftime("%H:%M:%S"), $5, "bytes"}'
```

### Multi-File Processing

```bash
# Compare two files (like join)
awk 'NR==FNR {a[$1]=$2; next} $1 in a {print $1, a[$1], $2}' file1.txt file2.txt

# Subtract file2 lines from file1
awk 'NR==FNR {a[$0]; next} !($0 in a)' file2.txt file1.txt

# Merge files by key
awk -F, 'NR==FNR {name[$1]=$2; next} {print $0, name[$1]}' names.csv data.csv
```

---

## 🐛 Troubleshooting

| Problem | Solution |
|---------|----------|
| Fields not splitting correctly | Check `-F` separator, use `BEGIN{FS=...}` |
| Floating point issues | Use `printf "%.2f"` for precision |
| "not a number" errors | Validate input: `$1 + 0` forces numeric |
| Large file memory issues | Process line by line, avoid storing everything in arrays |
| Different awk versions | Use `gawk` for full features, `mawk` for speed |
| Regex escaping | Use `\\` for literal backslash |
| Single quotes in awk | Use `'"'"'` or pass via `-v` variable |

---

## 🔗 Related Topics

- [🧰 CLI Command Center](../) — All CLI tool references
- [grep](../grep/) — Pattern matching (often piped to awk)
- [sed](../sed/) — Stream editing
- [jq](../jq/) — JSON processing (awk for JSON)
- [Linux](../linux/) — Core system commands
- [⚙️ Engineering](../../engineering/) — Automation and scripting

---

> **Tip:** For complex data processing, start with `awk '{print NF, $0}'` to understand your data structure, then build up the logic incrementally.
