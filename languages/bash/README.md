# 🐚 Bash Reference

> Comprehensive Bash scripting reference — from fundamentals to production-grade shell scripts.

---

## Table of Contents

1. [What is Bash](#what-is-bash)
2. [Variables](#variables)
3. [Quoting](#quoting)
4. [Control Flow](#control-flow)
5. [Functions](#functions)
6. [Arrays](#arrays)
7. [String Operations](#string-operations)
8. [File Operations](#file-operations)
9. [Input / Output](#input--output)
10. [Process Control](#process-control)
11. [Text Processing](#text-processing)
12. [Script Patterns](#script-patterns)
13. [Best Practices](#best-practices)
14. [Common Mistakes](#common-mistakes)
15. [Debugging](#debugging)
16. [Production Tips](#production-tips)
17. [References](#references)

---

## What is Bash

Bash (Bourne Again SHell) is a Unix shell and command language. It's the default shell on most Linux distributions and macOS (pre-Catalina). Bash scripts automate system tasks, CI/CD pipelines, deployments, and infrastructure management.

### Shebang and Script Execution

```bash
#!/bin/bash
# or for portability:
#!/usr/bin/env bash

echo "Hello, World!"
```

```bash
# Make executable and run
chmod +x script.sh
./script.sh

# Or run with bash directly
bash script.sh

# Check syntax without executing
bash -n script.sh
```

### When to Use Bash

| Use Case | Recommendation |
|----------|---------------|
| Glue code connecting CLI tools | ✅ Bash |
| File manipulation and text processing | ✅ Bash |
| Quick automation (< 100 lines) | ✅ Bash |
| CI/CD pipeline steps | ✅ Bash |
| Complex data structures | ❌ Use Python/Go |
| Error handling with retries | ❌ Use Python/Go |
| Cross-platform portability | ❌ Use Python/Go |
| Anything > 300 lines | ❌ Consider Python/Go |

---

## Variables

### Declaration and Assignment

```bash
# No spaces around =
name="Alice"
count=42
path="/var/log"

# Read-only
readonly DB_HOST="localhost"

# Unset a variable
unset name
```

### Environment Variables

```bash
# Export to child processes
export APP_ENV="production"
export PORT=8080

# Use in child process
echo "Running in $APP_ENV on port $PORT"

# Set for a single command
DB_HOST=localhost DB_PORT=5432 ./myapp
```

### Special Variables

| Variable | Description |
|----------|-------------|
| `$0` | Script name |
| `$1`, `$2`, ... | Positional parameters |
| `$#` | Number of arguments |
| `$@` | All arguments (as separate words) |
| `$*` | All arguments (as single string) |
| `$?` | Exit code of last command |
| `$$` | PID of current shell |
| `$!` | PID of last background process |
| `$_` | Last argument of previous command |
| `$LINENO` | Current line number |
| `$FUNCNAME` | Current function name |

```bash
#!/bin/bash
echo "Script: $0"
echo "Args: $@"
echo "Count: $#"
echo "First: $1"
echo "PID: $$"

some_command
echo "Exit code: $?"
```

### Variable Defaults

```bash
# Use default if unset or empty
echo "${NAME:-World}"              # prints "World" if NAME is unset/empty

# Set default if unset or empty
echo "${NAME:=World}"              # also assigns "World" to NAME

# Error if unset or empty
echo "${NAME:?'NAME is required'}" # exits with error if NAME is unset/empty

# Use alternative if set
echo "${NAME:+found}"              # prints "found" if NAME is set, empty otherwise
```

---

## Quoting

Quoting controls how the shell interprets special characters.

### Single Quotes — Literal

```bash
# Everything inside is literal — no variable expansion, no escaping
echo 'Hello $USER'        # Hello $USER
echo 'No \n escaping'     # No \n escaping
echo 'Even "quotes" work' # Even "quotes" work
```

### Double Quotes — Variable Expansion

```bash
# Variables and command substitution are expanded
echo "Hello $USER"           # Hello alice
echo "Home: ${HOME}"         # Home: /home/alice
echo "Date: $(date)"         # Date: Mon Aug 31 ...
echo "Tab:\there"            # Tab:   here (escape sequences work)
```

### Command Substitution

```bash
# Modern syntax (preferred)
today=$(date +%Y-%m-%d)
file_count=$(ls -1 | wc -l)

# Legacy backtick syntax (avoid — harder to nest)
today=`date +%Y-%m-%d`

# Nesting
files_modified=$(find . -newer "$(date -d '1 hour ago' +%Y%m%d%H%M)")
```

### When to Quote

```bash
# ✅ ALWAYS quote variables — prevents word splitting and globbing
file="my file.txt"
cp "$file" /tmp/              # correct
# cp $file /tmp/              # WRONG — splits into "my" and "file.txt"

# ✅ Quote command substitution
result="$(some_command)"

# Array expansion with quotes preserves elements
files=("file one.txt" "file two.txt")
for f in "${files[@]}"; do    # correct — preserves spaces
    echo "$f"
done
```

---

## Control Flow

### if / then / else / fi

```bash
if [[ "$status" == "active" ]]; then
    echo "Service is running"
elif [[ "$status" == "stopped" ]]; then
    echo "Service is stopped"
else
    echo "Unknown status: $status"
fi
```

### test vs [[ (Double Bracket)

```bash
# [[ is preferred over [ (test) — safer, more features
# [[ doesn't word-split or glob-expand variables

# String comparison
[[ "$name" == "Alice" ]]       # equal
[[ "$name" != "Bob" ]]         # not equal
[[ "$name" == A* ]]            # glob match
[[ "$name" =~ ^[A-Z] ]]       # regex match
[[ -z "$var" ]]                # empty/unset
[[ -n "$var" ]]                # non-empty

# Numeric comparison
[[ "$count" -eq 5 ]]           # equal
[[ "$count" -ne 0 ]]           # not equal
[[ "$count" -gt 10 ]]          # greater than
[[ "$count" -lt 100 ]]         # less than
[[ "$count" -ge 1 ]]           # greater or equal
[[ "$count" -le 50 ]]          # less or equal

# Logical operators
[[ "$a" -gt 0 && "$b" -gt 0 ]]  # AND
[[ "$a" -gt 0 || "$b" -gt 0 ]]  # OR
[[ ! -f "$file" ]]               # NOT

# Arithmetic comparison (alternative)
(( count > 10 ))
(( count >= 1 && count <= 100 ))
```

### for Loop

```bash
# Iterate over list
for fruit in apple banana cherry; do
    echo "$fruit"
done

# C-style for loop
for (( i=0; i<10; i++ )); do
    echo "$i"
done

# Iterate over files
for file in /var/log/*.log; do
    echo "Processing: $file"
done

# Iterate over command output
for user in $(getent passwd | cut -d: -f1); do
    echo "User: $user"
done

# Iterate over array
servers=("web1" "web2" "db1")
for server in "${servers[@]}"; do
    echo "Pinging $server..."
    ping -c 1 "$server"
done

# Iterate with index
for i in "${!servers[@]}"; do
    echo "$i: ${servers[$i]}"
done
```

### while and until

```bash
# while loop
count=0
while [[ "$count" -lt 5 ]]; do
    echo "Count: $count"
    (( count++ ))
done

# Read file line by line
while IFS= read -r line; do
    echo "Line: $line"
done < input.txt

# Read from command output
while IFS= read -r pid; do
    echo "Killing $pid"
    kill "$pid"
done < <(pgrep -f "myprocess")

# until loop (runs until condition is true)
until ping -c 1 google.com &>/dev/null; do
    echo "Waiting for network..."
    sleep 1
done
echo "Network is up!"
```

### case Statement

```bash
case "$1" in
    start)
        echo "Starting service..."
        start_service
        ;;
    stop)
        echo "Stopping service..."
        stop_service
        ;;
    restart)
        echo "Restarting..."
        stop_service
        start_service
        ;;
    status)
        check_status
        ;;
    *)
        echo "Usage: $0 {start|stop|restart|status}"
        exit 1
        ;;
esac

# Pattern matching in case
case "$filename" in
    *.tar.gz|*.tgz)
        tar xzf "$filename"
        ;;
    *.zip)
        unzip "$filename"
        ;;
    *.deb)
        dpkg -i "$filename"
        ;;
    *)
        echo "Unknown file type: $filename"
        ;;
esac
```

---

## Functions

### Declaration and Arguments

```bash
# Function declaration
greet() {
    local name="$1"
    local greeting="${2:-Hello}"
    echo "${greeting}, ${name}!"
}

greet "Alice"            # Hello, Alice!
greet "Bob" "Hi"         # Hi, Bob!

# Alternative syntax
function greet {
    echo "Hello, $1!"
}
```

### Return Values

Functions return exit codes (0-255). Use `echo` for data and command substitution to capture:

```bash
# Return exit code
is_file() {
    [[ -f "$1" ]]     # returns 0 (true) or 1 (false)
}

if is_file "/etc/hosts"; then
    echo "File exists"
fi

# Return data via stdout
get_timestamp() {
    date +"%Y-%m-%d %H:%M:%S"
}

timestamp=$(get_timestamp)
echo "Timestamp: $timestamp"

# Return multiple values
get_dimensions() {
    echo "1920 1080"
}

read -r width height <<< "$(get_dimensions)"
```

### Local Variables

```bash
# Always use local to avoid polluting global scope
process() {
    local input="$1"
    local result
    result=$(echo "$input" | tr '[:lower:]' '[:upper:]')
    echo "$result"
}

# Global variable gotcha
count=0
increment() {
    (( count++ ))    # modifies the global — intentional here
}
increment
echo "$count"        # 1
```

---

## Arrays

### Indexed Arrays

```bash
# Declaration
fruits=("apple" "banana" "cherry")
declare -a numbers=(1 2 3 4 5)

# Access
echo "${fruits[0]}"              # apple
echo "${fruits[-1]}"             # cherry (last element)

# All elements
echo "${fruits[@]}"              # apple banana cherry
echo "${fruits[*]}"              # apple banana cherry

# Length
echo "${#fruits[@]}"             # 3

# Append
fruits+=("date")

# Delete
unset 'fruits[1]'                # removes banana (leaves gap)

# Slice
echo "${fruits[@]:1:2}"          # 2 elements starting at index 1

# Iterate
for fruit in "${fruits[@]}"; do
    echo "$fruit"
done
```

### Associative Arrays (Bash 4+)

```bash
# Must declare with -A
declare -A config
config[host]="localhost"
config[port]="5432"
config[db]="myapp"

# Or initialize
declare -A colors=(
    [red]="#FF0000"
    [green]="#00FF00"
    [blue]="#0000FF"
)

# Access
echo "${config[host]}"

# Keys
echo "${!config[@]}"             # host port db

# Iterate
for key in "${!config[@]}"; do
    echo "$key = ${config[$key]}"
done

# Check key exists
if [[ -v config[host] ]]; then
    echo "host is set"
fi
```

---

## String Operations

### Substring and Length

```bash
str="Hello, World!"

# Length
echo "${#str}"                   # 13

# Substring (offset, length)
echo "${str:7}"                  # World!
echo "${str:7:5}"                # World
echo "${str: -6}"                # orld!  (note the space before -)
```

### Search and Replace

```bash
path="/home/user/documents/file.tar.gz"

# Remove prefix (shortest match)
echo "${path#*/}"                # home/user/documents/file.tar.gz

# Remove prefix (longest match)
echo "${path##*/}"               # file.tar.gz (like basename)

# Remove suffix (shortest match)
echo "${path%.*}"                # /home/user/documents/file.tar

# Remove suffix (longest match)
echo "${path%%.*}"               # /home/user/documents/file

# Replace first occurrence
echo "${path/user/admin}"        # /home/admin/documents/file.tar.gz

# Replace all occurrences
echo "${path//o/0}"              # /h0me/user/d0cuments/file.tar.gz

# Case conversion (Bash 4+)
name="hello world"
echo "${name^}"                  # Hello world (capitalize first)
echo "${name^^}"                 # HELLO WORLD (uppercase all)

upper="HELLO"
echo "${upper,}"                 # hELLO (lowercase first)
echo "${upper,,}"                # hello (lowercase all)
```

---

## File Operations

### Test Operators

```bash
# File existence and type
[[ -e "$path" ]]        # exists (file or directory)
[[ -f "$path" ]]        # is a regular file
[[ -d "$path" ]]        # is a directory
[[ -L "$path" ]]        # is a symbolic link
[[ -s "$path" ]]        # exists and is not empty
[[ -p "$path" ]]        # is a named pipe

# Permissions
[[ -r "$path" ]]        # is readable
[[ -w "$path" ]]        # is writable
[[ -x "$path" ]]        # is executable

# Comparison
[[ "$file1" -nt "$file2" ]]   # file1 is newer than file2
[[ "$file1" -ot "$file2" ]]   # file1 is older than file2

# Practical patterns
if [[ ! -f "$config_file" ]]; then
    echo "Error: Config file not found: $config_file" >&2
    exit 1
fi

if [[ ! -d "$output_dir" ]]; then
    mkdir -p "$output_dir"
fi

if [[ ! -x "$binary" ]]; then
    echo "Error: $binary is not executable" >&2
    exit 1
fi
```

### File Operations

```bash
# Create temporary file safely
tmpfile=$(mktemp)
trap 'rm -f "$tmpfile"' EXIT

# Create temporary directory
tmpdir=$(mktemp -d)
trap 'rm -rf "$tmpdir"' EXIT

# Atomic file write (write to temp, then move)
echo "$content" > "${config_file}.tmp"
mv "${config_file}.tmp" "$config_file"

# Find and process files
find /var/log -name "*.log" -mtime +30 -delete
find . -name "*.py" -exec grep -l "TODO" {} +
```

---

## Input / Output

### read

```bash
# Read user input
read -rp "Enter your name: " name
echo "Hello, $name!"

# Read with timeout
if read -rt 10 -p "Continue? [y/N] " answer; then
    [[ "$answer" == [yY] ]] && echo "Continuing..."
else
    echo "Timed out, aborting."
fi

# Read password (no echo)
read -rsp "Password: " password
echo

# Read into array
read -ra words <<< "one two three"
echo "${words[1]}"               # two
```

### echo and printf

```bash
# echo — simple output
echo "Hello, World"
echo -n "No newline"            # no trailing newline
echo -e "Tab:\there"            # interpret escape sequences

# printf — formatted output (preferred for portability)
printf "Name: %s, Age: %d\n" "Alice" 30
printf "%.2f\n" 3.14159
printf "%05d\n" 42               # 00042

# printf to variable
printf -v formatted "%-20s %d" "$name" "$score"
```

### Redirection

```bash
# Standard streams
# 0 = stdin, 1 = stdout, 2 = stderr

# Redirect stdout to file
command > output.txt             # overwrite
command >> output.txt            # append

# Redirect stderr
command 2> errors.txt
command 2>> errors.txt

# Redirect both
command > output.txt 2>&1        # both to same file
command &> output.txt            # shorthand (same thing)

# Redirect stderr to stdout (for piping)
command 2>&1 | grep "ERROR"

# Discard output
command > /dev/null 2>&1         # discard everything
command &> /dev/null             # shorthand

# Redirect stdin
command < input.txt
sort < unsorted.txt > sorted.txt
```

### Here Documents and Here Strings

```bash
# Here document — multi-line input
cat << 'EOF'
This is a multi-line string.
Variables are NOT expanded with single-quoted delimiter.
EOF

cat << EOF
Hello, $USER
Variables ARE expanded with unquoted delimiter.
EOF

# Indented here document (<<- strips leading tabs)
	cat <<- EOF
	Tabs are stripped from the beginning.
	Useful in functions and loops.
	EOF

# Here string — single-line input to stdin
grep "error" <<< "$log_output"
read -r first last <<< "Alice Smith"
```

---

## Process Control

### Exit Codes

```bash
# Exit with specific code
exit 0          # success
exit 1          # general error

# Exit codes convention
# 0   = success
# 1   = general error
# 2   = misuse of shell command
# 126 = command not executable
# 127 = command not found
# 128+n = killed by signal n

# Check last exit code
if command; then
    echo "Success"
else
    echo "Failed with exit code: $?"
fi
```

### Trap and Signals

```bash
# Trap cleanup on exit
cleanup() {
    echo "Cleaning up..."
    rm -f "$tmpfile"
    [[ -n "$child_pid" ]] && kill "$child_pid" 2>/dev/null
}
trap cleanup EXIT

# Trap specific signals
trap 'echo "Interrupted!"; exit 130' INT
trap 'echo "Terminated!"; exit 143' TERM

# Ignore signal
trap '' HUP

# Common signals
# SIGHUP  (1)  — terminal closed
# SIGINT  (2)  — Ctrl+C
# SIGTERM (15) — graceful termination
# SIGKILL (9)  — force kill (cannot be trapped)
```

### Background Processes

```bash
# Run in background
long_running_command &
child_pid=$!

# Wait for specific process
wait "$child_pid"
echo "Process exited with: $?"

# Wait for all background jobs
wait

# Run multiple in parallel
for server in web1 web2 web3; do
    deploy "$server" &
done
wait
echo "All deployments complete"

# Process substitution
diff <(sort file1.txt) <(sort file2.txt)
```

---

## Text Processing

### Combining grep, sed, awk

```bash
# Find errors in logs, extract timestamps
grep "ERROR" /var/log/app.log | awk '{print $1, $2}' | sort | uniq -c | sort -rn

# Replace config values
sed -i 's/DEBUG=true/DEBUG=false/' config.env

# Extract fields from CSV
awk -F',' '{print $1, $3}' data.csv

# Count lines matching a pattern
grep -c "ERROR" /var/log/app.log

# Extract unique IPs from access log
awk '{print $1}' access.log | sort -u

# Sum a column
awk '{sum += $3} END {print sum}' data.txt

# Multi-line log processing
awk '/START/,/END/' logfile.txt

# Replace in multiple files
find . -name "*.conf" -exec sed -i 's/old_value/new_value/g' {} +
```

→ See [CLI / grep](../../cli/grep/), [CLI / sed](../../cli/sed/), [CLI / awk](../../cli/awk/) for comprehensive references.

---

## Script Patterns

### Option Parsing with getopts

```bash
#!/bin/bash
set -euo pipefail

usage() {
    cat << EOF
Usage: $(basename "$0") [OPTIONS] <input_file>

Options:
    -o FILE    Output file (default: stdout)
    -v         Verbose mode
    -n NUM     Number of lines to process (default: all)
    -h         Show this help message

Example:
    $(basename "$0") -v -n 100 -o output.txt input.txt
EOF
}

output="/dev/stdout"
verbose=false
num_lines=0

while getopts ":o:vn:h" opt; do
    case "$opt" in
        o) output="$OPTARG" ;;
        v) verbose=true ;;
        n) num_lines="$OPTARG" ;;
        h) usage; exit 0 ;;
        :) echo "Error: -$OPTARG requires an argument" >&2; exit 1 ;;
        *) echo "Error: Unknown option -$OPTARG" >&2; usage; exit 1 ;;
    esac
done
shift $((OPTIND - 1))

# Validate required arguments
if [[ $# -lt 1 ]]; then
    echo "Error: Input file required" >&2
    usage
    exit 1
fi

input_file="$1"

if [[ ! -f "$input_file" ]]; then
    echo "Error: File not found: $input_file" >&2
    exit 1
fi

$verbose && echo "Processing $input_file..."
```

### Logging Function

```bash
# Colored log output
readonly RED='\033[0;31m'
readonly GREEN='\033[0;32m'
readonly YELLOW='\033[0;33m'
readonly BLUE='\033[0;34m'
readonly NC='\033[0m'  # No Color

log_info()  { echo -e "${BLUE}[INFO]${NC}  $(date '+%H:%M:%S') $*"; }
log_ok()    { echo -e "${GREEN}[OK]${NC}    $(date '+%H:%M:%S') $*"; }
log_warn()  { echo -e "${YELLOW}[WARN]${NC}  $(date '+%H:%M:%S') $*" >&2; }
log_error() { echo -e "${RED}[ERROR]${NC} $(date '+%H:%M:%S') $*" >&2; }

die() {
    log_error "$@"
    exit 1
}

# Usage
log_info "Starting deployment..."
log_ok "Service started successfully"
log_warn "Config file uses deprecated format"
log_error "Connection refused"
die "Cannot continue without database connection"
```

### Error Handling Pattern

```bash
#!/bin/bash
set -euo pipefail

# Cleanup on exit
cleanup() {
    local exit_code=$?
    if [[ "$exit_code" -ne 0 ]]; then
        log_error "Script failed with exit code $exit_code"
    fi
    # Remove temp files
    rm -f "${TMPFILES[@]}" 2>/dev/null || true
    exit "$exit_code"
}
trap cleanup EXIT

TMPFILES=()

# Safe temp file creation
create_temp() {
    local tmpfile
    tmpfile=$(mktemp)
    TMPFILES+=("$tmpfile")
    echo "$tmpfile"
}

# Retry pattern
retry() {
    local max_attempts="$1"
    local delay="$2"
    shift 2
    local attempt=1

    until "$@"; do
        if (( attempt >= max_attempts )); then
            log_error "Command failed after $max_attempts attempts: $*"
            return 1
        fi
        log_warn "Attempt $attempt/$max_attempts failed, retrying in ${delay}s..."
        sleep "$delay"
        (( attempt++ ))
    done
}

# Usage
retry 5 2 curl -sf "https://api.example.com/health"
```

---

## Best Practices

### Always Start With

```bash
#!/bin/bash
set -euo pipefail

# set -e    — exit on any command failure
# set -u    — treat unset variables as errors
# set -o pipefail — pipe fails if any command in pipe fails
```

### ShellCheck

Always run [ShellCheck](https://www.shellcheck.net/) on your scripts:

```bash
# Install
# apt-get install shellcheck    # Debian/Ubuntu
# brew install shellcheck       # macOS

# Run
shellcheck script.sh
shellcheck -x script.sh         # follow sourced files

# Inline directives to suppress warnings
# shellcheck disable=SC2086
echo $unquoted_var
```

### Quoting Rules

```bash
# ✅ Always quote variables
echo "$variable"
cp "$source" "$destination"

# ✅ Always quote command substitutions
result="$(some_command)"

# ✅ Quote array expansions
for item in "${array[@]}"; do
    process "$item"
done

# ✅ Exception: arithmetic context doesn't need quotes
(( count++ ))
if (( count > 10 )); then ...
```

### Portable Scripts

```bash
# Use POSIX-compatible constructs when possible
# Prefer $(command) over `command`
# Prefer [[ ]] over [ ] in bash-specific scripts
# Use printf over echo for portability

# Check for required commands
require_command() {
    if ! command -v "$1" &>/dev/null; then
        echo "Error: Required command '$1' not found" >&2
        exit 1
    fi
}

require_command curl
require_command jq
require_command docker
```

---

## Common Mistakes

### Word Splitting

```bash
# ❌ WRONG — fails if filename has spaces
file=my file.txt
cat $file            # cat tries "my" and "file.txt" separately

# ✅ CORRECT
file="my file.txt"
cat "$file"
```

### Glob Expansion

```bash
# ❌ WRONG — *.log expands to matching files
message="Check *.log files"
echo $message        # "Check access.log error.log files"

# ✅ CORRECT
echo "$message"      # "Check *.log files"
```

### Missing Quotes in Conditionals

```bash
# ❌ WRONG — breaks if $var is empty
if [ $var = "value" ]; then ...  # syntax error if var is empty

# ✅ CORRECT
if [[ "$var" = "value" ]]; then ...
```

### Testing Empty String

```bash
# ❌ WRONG — test -n without quotes
if [ -n $maybe_empty ]; then ...  # always true!

# ✅ CORRECT
if [[ -n "$maybe_empty" ]]; then ...
```

### cd in Scripts

```bash
# ❌ WRONG — if cd fails, rm runs in wrong directory
cd /some/path
rm -rf ./build

# ✅ CORRECT
cd /some/path || exit 1
rm -rf ./build

# ✅ BETTER — use subshell to restore directory
(
    cd /some/path || exit 1
    rm -rf ./build
)
```

### Parsing ls Output

```bash
# ❌ WRONG — parsing ls is fragile
for file in $(ls *.txt); do ...

# ✅ CORRECT — use globbing
for file in *.txt; do
    [[ -e "$file" ]] || continue    # handle no matches
    process "$file"
done
```

---

## Debugging

### Debug Mode

```bash
# Enable debug output (prints every command before executing)
set -x

# Disable debug output
set +x

# Run script in debug mode
bash -x script.sh

# Debug specific sections
set -x
problematic_function
set +x
```

### Custom Debug Trace

```bash
# Customize the trace prompt
export PS4='+${BASH_SOURCE}:${LINENO}:${FUNCNAME[0]:-main}: '
set -x

# Now traces show file:line:function
# +script.sh:42:process_data: grep "ERROR" /var/log/app.log
```

### Syntax Check

```bash
# Check syntax without executing
bash -n script.sh

# Verbose mode (print each line as it's read)
bash -v script.sh
```

### Debugging Tips

```bash
# Print variable state
declare -p my_array              # show array details
echo "DEBUG: var=$var" >&2       # debug to stderr

# Trap ERR for error reporting
trap 'echo "Error on line $LINENO, exit code $?" >&2' ERR

# Trace function calls
trap 'echo "TRACE: ${FUNCNAME[0]:-main} line $LINENO" >&2' DEBUG
```

---

## Production Tips

### Signal Handling and Cleanup

```bash
#!/bin/bash
set -euo pipefail

PID_FILE="/var/run/myservice.pid"
LOG_FILE="/var/log/myservice.log"

# Write PID file
echo $$ > "$PID_FILE"

# Cleanup on any exit
cleanup() {
    local exit_code=$?
    log_info "Shutting down (exit code: $exit_code)..."
    rm -f "$PID_FILE"
    # Kill child processes
    jobs -p | xargs -r kill 2>/dev/null || true
    wait 2>/dev/null || true
    exit "$exit_code"
}

trap cleanup EXIT INT TERM
```

### Lock Files (Prevent Duplicate Execution)

```bash
LOCK_FILE="/var/lock/myscript.lock"

acquire_lock() {
    if ! mkdir "$LOCK_FILE" 2>/dev/null; then
        echo "Error: Another instance is running (lock: $LOCK_FILE)" >&2
        exit 1
    fi
    trap 'rm -rf "$LOCK_FILE"' EXIT
}

acquire_lock
# ... rest of script
```

### Log Rotation Pattern

```bash
rotate_log() {
    local log_file="$1"
    local max_size="${2:-10485760}"  # 10 MB default
    local max_files="${3:-5}"

    if [[ -f "$log_file" ]] && (( $(stat -f%z "$log_file" 2>/dev/null || stat -c%s "$log_file") > max_size )); then
        for (( i=max_files-1; i>=1; i-- )); do
            [[ -f "${log_file}.$i" ]] && mv "${log_file}.$i" "${log_file}.$((i+1))"
        done
        mv "$log_file" "${log_file}.1"
        touch "$log_file"
    fi
}
```

### Configuration File Parsing

```bash
# Parse key=value config file
parse_config() {
    local config_file="$1"
    while IFS='=' read -r key value; do
        # Skip empty lines and comments
        [[ -z "$key" || "$key" == \#* ]] && continue
        # Trim whitespace
        key=$(echo "$key" | xargs)
        value=$(echo "$value" | xargs)
        # Export as variable
        export "$key"="$value"
    done < "$config_file"
}

parse_config /etc/myapp/config.env
echo "DB Host: $DB_HOST"
```

### Health Check Script

```bash
#!/bin/bash
set -euo pipefail

check_http() {
    local url="$1"
    local timeout="${2:-5}"
    if curl -sf --max-time "$timeout" "$url" > /dev/null; then
        echo "✅ $url"
        return 0
    else
        echo "❌ $url"
        return 1
    fi
}

check_port() {
    local host="$1"
    local port="$2"
    if timeout 3 bash -c "echo > /dev/tcp/$host/$port" 2>/dev/null; then
        echo "✅ $host:$port is open"
        return 0
    else
        echo "❌ $host:$port is closed"
        return 1
    fi
}

failed=0
check_http "http://localhost:8080/health" || (( failed++ ))
check_port "localhost" 5432                || (( failed++ ))
check_port "localhost" 6379                || (( failed++ ))

if (( failed > 0 )); then
    echo "⚠️  $failed check(s) failed"
    exit 1
fi
echo "All checks passed!"
```

---

## References

### Official Resources

- [GNU Bash Manual](https://www.gnu.org/software/bash/manual/) — The definitive reference
- [Bash Reference Manual](https://www.gnu.org/software/bash/manual/bashref.html) — Complete syntax reference
- [POSIX Shell Specification](https://pubs.opengroup.org/onlinepubs/9699919799/utilities/V3_chap02.html) — POSIX standard

### Guides and Tutorials

- [Bash Guide](https://mywiki.wooledge.org/BashGuide) — Comprehensive community guide
- [Bash FAQ](https://mywiki.wooledge.org/BashFAQ) — Common questions and answers
- [Bash Pitfalls](https://mywiki.wooledge.org/BashPitfalls) — Common mistakes
- [ShellCheck](https://www.shellcheck.net/) — Online shell script linter
- [Google Shell Style Guide](https://google.github.io/styleguide/shellguide.html) — Google's style conventions

### Books

- *Learning the bash Shell* (3rd ed.) by Cameron Newham
- *bash Cookbook* (2nd ed.) by Carl Albing & JP Vossen
- *Classic Shell Scripting* by Arnold Robbins & Nelson Beebe

### Related Topics

- [🧰 CLI / Linux](../../cli/linux/) — Linux command reference
- [🧰 CLI / grep](../../cli/grep/) — Pattern searching
- [🧰 CLI / sed](../../cli/sed/) — Stream editing
- [🧰 CLI / awk](../../cli/awk/) — Text processing
- [🧰 CLI / curl](../../cli/curl/) — HTTP requests
- [🧰 CLI / jq](../../cli/jq/) — JSON processing
- [🧰 CLI / SSH](../../cli/ssh/) — Remote shell
- [🚀 DevOps](../../devops/) — CI/CD pipelines

---

*Part of [Tech-Byte Languages](../). Last updated: 2026-08.*
