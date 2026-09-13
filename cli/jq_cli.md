# 🔧 jq — JSON Processing Reference

> **The Swiss Army knife for JSON.** Parse, filter, transform, and format JSON data from APIs, config files, logs, and command-line tool output.

---

## 📑 Table of Contents

- [Concepts](#-concepts)
- [Basic Filters](#-basic-filters)
- [Types & Values](#-types--values)
- [Pipe & Composition](#-pipe--composition)
- [String Operations](#-string-operations)
- [Array Operations](#-array-operations)
- [Object Operations](#-object-operations)
- [Conditionals & Comparisons](#-conditionals--comparisons)
- [Map, Select & Reduce](#-map-select--reduce)
- [Built-in Functions](#-built-in-functions)
- [Advanced Features](#-advanced-features)
- [Practical Examples](#-practical-examples)
- [Troubleshooting](#-troubleshooting)
- [Related Topics](#-related-topics)

---

## 💡 Concepts

`jq` is a command-line JSON processor. It takes JSON input, applies filters, and produces output. Filters can be piped together like Unix commands.

```
Input JSON → jq 'filter expression' → Output
```

### Common Flags

| Flag | Description |
|------|-------------|
| `-r` | Raw output (no quotes around strings) |
| `-e` | Set exit status based on output (false/null = 1) |
| `-s` | Slurp: read entire input as single array |
| `-n` | Null input (don't read input, use for generating JSON) |
| `-c` | Compact output (single line) |
| `-S` | Sort keys |
| `--arg name val` | Pass string variable |
| `--argjson name val` | Pass JSON variable |
| `--slurpfile name file` | Load file as variable |
| `--rawfile name file` | Load file as raw string |

---

## 🔍 Basic Filters

### Identity and Field Access

```bash
# Pretty-print JSON
echo '{"name":"Alice","age":30}' | jq '.'

# Access field
echo '{"name":"Alice","age":30}' | jq '.name'
# "Alice"

# Raw string output (no quotes)
echo '{"name":"Alice"}' | jq -r '.name'
# Alice

# Nested field access
echo '{"user":{"name":"Alice","address":{"city":"NYC"}}}' | jq '.user.address.city'
# "NYC"

# Optional field (no error if missing)
echo '{"name":"Alice"}' | jq '.age // "N/A"'
# "N/A"
```

### Array Access

```bash
# Access array element
echo '[1,2,3,4,5]' | jq '.[0]'
# 1

echo '[1,2,3,4,5]' | jq '.[-1]'
# 5

# Array slice
echo '[1,2,3,4,5]' | jq '.[2:4]'
# [3,4]

# Iterate array (unwrap elements)
echo '[1,2,3]' | jq '.[]'
# 1
# 2
# 3

# Access field from array of objects
echo '[{"name":"Alice"},{"name":"Bob"}]' | jq '.[].name'
# "Alice"
# "Bob"

# Wrap results back into array
echo '[{"name":"Alice"},{"name":"Bob"}]' | jq '[.[].name]'
# ["Alice","Bob"]
```

---

## 📦 Types & Values

```bash
# Type checking
echo '{"a":1,"b":"hello","c":null,"d":true,"e":[1,2]}' | jq '.a | type'
# "number"

echo '"hello"' | jq 'type'
# "string"

# Length
echo '[1,2,3]' | jq 'length'
# 3

echo '{"a":1,"b":2}' | jq 'length'
# 2

echo '"hello"' | jq 'length'
# 5

# Null handling
echo '{"a":null}' | jq '.a // "default"'
# "default"

# Type conversion
echo '"42"' | jq 'tonumber'
# 42

echo '42' | jq 'tostring'
# "42"
```

---

## 🔗 Pipe & Composition

```bash
# Pipe filters together (like Unix pipes)
echo '{"users":[{"name":"Alice","age":30},{"name":"Bob","age":25}]}' | \
  jq '.users[] | .name'
# "Alice"
# "Bob"

# Construct new objects
echo '{"first":"Alice","last":"Smith","age":30}' | \
  jq '{fullName: (.first + " " + .last), age: .age}'
# {"fullName":"Alice Smith","age":30}

# Multiple outputs with comma
echo '{"a":1,"b":2,"c":3}' | jq '.a, .c'
# 1
# 3
```

---

## 📝 String Operations

```bash
# String interpolation
echo '{"name":"Alice","age":30}' | jq '"Name: \(.name), Age: \(.age)"'
# "Name: Alice, Age: 30"

# String functions
echo '"Hello World"' | jq 'ascii_downcase'
# "hello world"

echo '"hello world"' | jq 'ascii_upcase'
# "HELLO WORLD"

echo '"  hello  "' | jq 'ltrimstr(" ") | rtrimstr(" ")'

# Split and join
echo '"a,b,c,d"' | jq 'split(",")'
# ["a","b","c","d"]

echo '["a","b","c"]' | jq 'join("-")'
# "a-b-c"

# Test (regex match)
echo '"hello world"' | jq 'test("hello")'
# true

echo '"error: disk full"' | jq 'test("error|warning")'
# true

# Match (regex capture)
echo '"2024-08-15"' | jq 'match("([0-9]{4})-([0-9]{2})-([0-9]{2})") | .captures[].string'
# "2024"
# "08"
# "15"

# Sub (regex replace)
echo '"Hello World"' | jq 'sub("World";"jq")'
# "Hello jq"

echo '"foo-bar-baz"' | jq 'gsub("-";"_")'
# "foo_bar_baz"

# Starts/ends with
echo '"hello world"' | jq 'startswith("hello")'
# true

echo '"file.json"' | jq 'endswith(".json")'
# true

# Contains
echo '"hello world"' | jq 'contains("world")'
# true
```

---

## 📋 Array Operations

```bash
# Length
echo '[1,2,3,4,5]' | jq 'length'
# 5

# Add element
echo '[1,2,3]' | jq '. + [4,5]'
# [1,2,3,4,5]

# Flatten nested arrays
echo '[[1,2],[3,[4,5]]]' | jq 'flatten'
# [1,2,3,4,5]

# Unique values
echo '[1,2,2,3,3,3]' | jq 'unique'
# [1,2,3]

# Sort
echo '[3,1,4,1,5,9]' | jq 'sort'
# [1,1,3,4,5,9]

# Sort objects by field
echo '[{"name":"Bob","age":25},{"name":"Alice","age":30}]' | jq 'sort_by(.name)'
# [{"name":"Alice","age":30},{"name":"Bob","age":25}]

# Reverse
echo '[1,2,3]' | jq 'reverse'
# [3,2,1]

# First / Last
echo '[1,2,3,4,5]' | jq 'first'
# 1
echo '[1,2,3,4,5]' | jq 'last'
# 5

# Group by
echo '[{"dept":"eng","name":"Alice"},{"dept":"eng","name":"Bob"},{"dept":"sales","name":"Carol"}]' | \
  jq 'group_by(.dept)'
# [[{"dept":"eng","name":"Alice"},{"dept":"eng","name":"Bob"}],[{"dept":"sales","name":"Carol"}]]

# Min / Max
echo '[3,1,4,1,5,9]' | jq 'min'
# 1
echo '[{"name":"Alice","age":30},{"name":"Bob","age":25}]' | jq 'min_by(.age)'
# {"name":"Bob","age":25}

# Contains (array)
echo '[1,2,3]' | jq 'contains([2,3])'
# true

# Index / indices
echo '["a","b","c","b"]' | jq 'index("b")'
# 1
```

---

## 🗂️ Object Operations

```bash
# Keys and values
echo '{"name":"Alice","age":30,"city":"NYC"}' | jq 'keys'
# ["age","city","name"]

echo '{"name":"Alice","age":30}' | jq 'values'
# ["Alice",30]

# Key-value pairs
echo '{"name":"Alice","age":30}' | jq 'to_entries'
# [{"key":"name","value":"Alice"},{"key":"age","value":30}]

echo '[{"key":"name","value":"Alice"}]' | jq 'from_entries'
# {"name":"Alice"}

# Add/merge objects
echo '{"a":1}' | jq '. + {"b":2}'
# {"a":1,"b":2}

echo '{"a":1,"b":2}' | jq '. * {"b":3,"c":4}'
# {"a":1,"b":3,"c":4}

# Delete key
echo '{"a":1,"b":2,"c":3}' | jq 'del(.b)'
# {"a":1,"c":3}

# Has key
echo '{"name":"Alice"}' | jq 'has("name")'
# true

# Select specific keys
echo '{"name":"Alice","age":30,"email":"a@b.com","phone":"123"}' | \
  jq '{name, email}'
# {"name":"Alice","email":"a@b.com"}

# Rename keys
echo '{"first_name":"Alice"}' | jq '{name: .first_name}'
# {"name":"Alice"}

# With entries (transform keys/values)
echo '{"name":"Alice","age":"30"}' | \
  jq 'with_entries(.key |= ascii_upcase)'
# {"NAME":"Alice","AGE":"30"}
```

---

## ❓ Conditionals & Comparisons

```bash
# If-then-else
echo '{"age":25}' | jq 'if .age >= 18 then "adult" else "minor" end'
# "adult"

# Comparison operators: ==, !=, <, >, <=, >=
echo '{"score":85}' | jq '.score > 80'
# true

# And / Or / Not
echo '{"age":25,"role":"admin"}' | jq '.age > 18 and .role == "admin"'
# true

echo '{"active":false}' | jq '.active | not'
# true

# Alternative operator (default values)
echo '{"name":"Alice"}' | jq '.age // 0'
# 0

echo 'null' | jq '. // "default"'
# "default"

# Try-catch (handle errors)
echo '"not-a-number"' | jq 'try tonumber catch "invalid"'
# "invalid"
```

---

## 🔄 Map, Select & Reduce

### Map

```bash
# Transform each element
echo '[1,2,3,4,5]' | jq 'map(. * 2)'
# [2,4,6,8,10]

echo '[{"name":"Alice","age":30},{"name":"Bob","age":25}]' | \
  jq 'map(.name)'
# ["Alice","Bob"]

# Map with condition
echo '[1,2,3,4,5]' | jq 'map(select(. > 3))'
# [4,5]

# Map objects (transform values)
echo '{"a":1,"b":2,"c":3}' | jq 'map_values(. * 10)'
# {"a":10,"b":20,"c":30}
```

### Select

```bash
# Filter array elements
echo '[{"name":"Alice","age":30},{"name":"Bob","age":25},{"name":"Carol","age":35}]' | \
  jq '.[] | select(.age > 28)'
# {"name":"Alice","age":30}
# {"name":"Carol","age":35}

# Keep as array
echo '[{"name":"Alice","age":30},{"name":"Bob","age":25}]' | \
  jq '[.[] | select(.age >= 30)]'
# [{"name":"Alice","age":30}]

# Select by string match
echo '[{"status":"active"},{"status":"inactive"},{"status":"active"}]' | \
  jq '[.[] | select(.status == "active")]'
# [{"status":"active"},{"status":"active"}]

# Select with regex
echo '[{"name":"test-pod-1"},{"name":"prod-pod-1"},{"name":"test-pod-2"}]' | \
  jq '[.[] | select(.name | test("^test"))]'
```

### Reduce

```bash
# Sum values
echo '[1,2,3,4,5]' | jq 'reduce .[] as $x (0; . + $x)'
# 15

# Build object from array
echo '[{"key":"a","val":1},{"key":"b","val":2}]' | \
  jq 'reduce .[] as $x ({}; . + {($x.key): $x.val})'
# {"a":1,"b":2}

# Shorthand for common reductions
echo '[1,2,3,4,5]' | jq 'add'
# 15

echo '[[1,2],[3,4]]' | jq 'add'
# [1,2,3,4]
```

---

## 🧮 Built-in Functions

| Function | Description | Example |
|----------|-------------|---------|
| `length` | Length of array/string/object | `[1,2,3] \| length` → `3` |
| `keys` | Object keys | `{"a":1} \| keys` → `["a"]` |
| `values` | Object values | `{"a":1} \| values` → `[1]` |
| `type` | Value type | `42 \| type` → `"number"` |
| `empty` | Produce no output | Used in conditionals |
| `error` | Raise error | `error("bad input")` |
| `debug` | Print to stderr | `.items \| debug \| length` |
| `env` | Environment variables | `env.HOME` |
| `input` | Read next input | For multi-document processing |
| `range(n)` | Generate 0..n-1 | `[range(5)]` → `[0,1,2,3,4]` |
| `limit(n;f)` | First n results | `limit(3; .items[])` |
| `ascii_downcase` | Lowercase | `"ABC" \| ascii_downcase` |
| `ascii_upcase` | Uppercase | `"abc" \| ascii_upcase` |
| `tostring` | Convert to string | `42 \| tostring` → `"42"` |
| `tonumber` | Convert to number | `"42" \| tonumber` → `42` |
| `nan` / `isinfinite` / `isnan` | Numeric checks | Special value tests |
| `now` | Current Unix timestamp | `now` |
| `strftime` | Format timestamp | `now \| strftime("%Y-%m-%d")` |

---

## 🔬 Advanced Features

### Variables

```bash
# Bind variables with as
echo '{"items":[1,2,3]}' | jq '.items | length as $count | "Total: \($count)"'
# "Total: 3"

# Pass variables from shell
echo '{"users":["Alice","Bob","Carol"]}' | \
  jq --arg name "Bob" '.users[] | select(. == $name)'
# "Bob"

# Pass JSON value from shell
echo '{"scores":[80,90,70]}' | \
  jq --argjson threshold 85 '[.scores[] | select(. >= $threshold)]'
# [90]
```

### Output Formatting

```bash
# Compact output (for piping/storage)
echo '{"name":"Alice","age":30}' | jq -c '.'
# {"name":"Alice","age":30}

# Tab-separated values
echo '[{"name":"Alice","age":30},{"name":"Bob","age":25}]' | \
  jq -r '.[] | [.name, .age] | @tsv'
# Alice	30
# Bob	25

# CSV format
echo '[{"name":"Alice","age":30},{"name":"Bob","age":25}]' | \
  jq -r '.[] | [.name, .age] | @csv'
# "Alice",30
# "Bob",25

# HTML encoding
echo '"<script>alert(1)</script>"' | jq '@html'
# "&lt;script&gt;alert(1)&lt;/script&gt;"

# URI encoding
echo '"hello world & more"' | jq '@uri'
# "hello%20world%20%26%20more"

# Base64
echo '"hello"' | jq '@base64'
# "aGVsbG8="

echo '"aGVsbG8="' | jq '@base64d'
# "hello"
```

### Slurp Mode

```bash
# Process multiple JSON documents as array
echo '{"a":1}
{"a":2}
{"a":3}' | jq -s 'map(.a) | add'
# 6

# Process JSONL (newline-delimited JSON)
cat events.jsonl | jq -s 'group_by(.type) | map({type: .[0].type, count: length})'
```

---

## 💼 Practical Examples

### Parsing API Responses

```bash
# Extract users from paginated API
curl -s https://api.example.com/users | \
  jq '.data[] | {id, name: .attributes.name, email: .attributes.email}'

# Get specific field from response
curl -s https://api.example.com/status | jq -r '.version'
```

### Kubernetes JSON Output

```bash
# Get all pod names
kubectl get pods -o json | jq -r '.items[].metadata.name'

# Get pods with their status
kubectl get pods -o json | \
  jq -r '.items[] | "\(.metadata.name)\t\(.status.phase)"'

# Find pods not running
kubectl get pods -o json | \
  jq '[.items[] | select(.status.phase != "Running") | .metadata.name]'

# Get container images across all pods
kubectl get pods -A -o json | \
  jq -r '[.items[].spec.containers[].image] | unique | .[]'

# Get resource requests/limits
kubectl get pods -o json | \
  jq '.items[] | {
    name: .metadata.name,
    containers: [.spec.containers[] | {
      name: .name,
      cpu_request: .resources.requests.cpu,
      mem_request: .resources.requests.memory,
      cpu_limit: .resources.limits.cpu,
      mem_limit: .resources.limits.memory
    }]
  }'

# Get pods sorted by restart count
kubectl get pods -o json | \
  jq '[.items[] | {
    name: .metadata.name,
    restarts: ([.status.containerStatuses[]?.restartCount] | add)
  }] | sort_by(.restarts) | reverse'
```

### Log Processing

```bash
# Parse JSON logs
cat app.log | jq 'select(.level == "ERROR") | {timestamp: .ts, message: .msg}'

# Count log levels
cat app.log | jq -s 'group_by(.level) | map({level: .[0].level, count: length})'

# Extract errors in time range
cat app.log | jq 'select(.level == "ERROR" and .ts >= "2024-08-15T00:00:00")'

# Top error messages
cat app.log | jq -s '[.[] | select(.level == "ERROR")] | group_by(.msg) | 
  map({message: .[0].msg, count: length}) | sort_by(.count) | reverse | .[0:10]'
```

### Configuration Manipulation

```bash
# Merge config files
jq -s '.[0] * .[1]' defaults.json overrides.json

# Update nested value
cat config.json | jq '.database.host = "new-host.example.com"'

# Add to array
cat config.json | jq '.allowed_ips += ["10.0.0.5"]'

# Remove sensitive fields before sharing
cat config.json | jq 'del(.password, .api_key, .secret)'

# Convert environment variables to JSON
env | grep "^APP_" | jq -R 'split("=") | {(.[0]): .[1:] | join("=")}' | jq -s 'add'
```

### Data Transformation

```bash
# JSON to CSV
echo '[{"name":"Alice","age":30},{"name":"Bob","age":25}]' | \
  jq -r '(.[0] | keys_unsorted) as $keys | $keys, (.[] | [.[$keys[]]]) | @csv'
# "name","age"
# "Alice",30
# "Bob",25

# Reshape data
echo '{"results":[{"id":1,"attrs":{"color":"red"}},{"id":2,"attrs":{"color":"blue"}}]}' | \
  jq '[.results[] | {id, color: .attrs.color}]'
# [{"id":1,"color":"red"},{"id":2,"color":"blue"}]
```

---

## 🐛 Troubleshooting

| Problem | Solution |
|---------|----------|
| `null` output | Field doesn't exist — check path with `.` |
| String with quotes in output | Use `-r` flag for raw output |
| `parse error` | Input is not valid JSON — validate with `jq '.' file.json` |
| Multiple JSON objects not parsed | Use `-s` (slurp) or process line by line |
| Shell variable expansion issues | Use `--arg` instead of string interpolation |
| `Cannot iterate over null` | Add `?` (optional): `.items[]?` or use `// []` |

---

## 🔗 Related Topics

- [🧰 CLI Command Center](../) — All CLI tool references
- [curl](../curl/) — HTTP client (pipe responses to jq)
- [kubectl](../kubectl/) — Kubernetes (use with `-o json`)
- [grep](../grep/) — Text pattern searching
- [awk](../awk/) — Text processing

---

> **Tip:** Build jq filters incrementally. Start with `.` to see the structure, add one filter at a time, and pipe to build complex transformations step by step.
