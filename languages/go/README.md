# 🐹 Go Reference

> Comprehensive Go reference — from fundamentals to production concurrency patterns. Based on Go 1.21+.

---

## Table of Contents

1. [What is Go](#what-is-go)
2. [Data Types](#data-types)
3. [Variables & Constants](#variables--constants)
4. [Control Flow](#control-flow)
5. [Arrays, Slices & Maps](#arrays-slices--maps)
6. [Functions](#functions)
7. [Structs & Methods](#structs--methods)
8. [Interfaces](#interfaces)
9. [Pointers](#pointers)
10. [Packages & Modules](#packages--modules)
11. [Error Handling](#error-handling)
12. [Goroutines & Channels](#goroutines--channels)
13. [Concurrency Patterns](#concurrency-patterns)
14. [Testing](#testing)
15. [Common Patterns](#common-patterns)
16. [CLI Tools](#cli-tools)
17. [Common Mistakes](#common-mistakes)
18. [Production Tips](#production-tips)
19. [References](#references)

---

## What is Go

Go (Golang) is a statically typed, compiled language designed at Google by Robert Griesemer, Rob Pike, and Ken Thompson. Released in 2009, Go emphasizes simplicity, reliability, and efficiency — especially for concurrent, networked systems.

### Why Go

- **Fast compilation** — compiles to a single static binary in seconds
- **Built-in concurrency** — goroutines and channels are first-class
- **Simple language** — ~25 keywords, one way to do things
- **Strong standard library** — HTTP server, JSON, crypto, testing built in
- **Cross-compilation** — `GOOS=linux GOARCH=amd64 go build`
- **Garbage collected** — no manual memory management

### Go Philosophy

- Simplicity over cleverness
- Composition over inheritance
- Explicit error handling over exceptions
- Interfaces are implicit (structural typing)
- "A little copying is better than a little dependency"

### Program Structure

Every Go program starts with a `package` declaration and a `main` function in `package main`:

```go
package main

import "fmt"

func main() {
    fmt.Println("Hello, World!")
}
```

Two types of Go programs:

1. **Executable** — `package main` with a `main()` function. Compiles to a binary.
2. **Library** — Any other package name. Imported and used by other packages.

```bash
go run main.go          # compile and run
go build -o myapp       # compile to binary
```

---

## Data Types

### Basic Types

| Type | Zero Value | Description |
|------|-----------|-------------|
| `bool` | `false` | Boolean |
| `string` | `""` | UTF-8 string (immutable) |
| `int`, `int8`, `int16`, `int32`, `int64` | `0` | Signed integers |
| `uint`, `uint8`, `uint16`, `uint32`, `uint64` | `0` | Unsigned integers |
| `float32`, `float64` | `0.0` | Floating point |
| `complex64`, `complex128` | `(0+0i)` | Complex numbers |
| `byte` | `0` | Alias for `uint8` |
| `rune` | `0` | Alias for `int32` (Unicode code point) |

> **Zero values:** When a variable is declared without initialization, Go assigns the type's zero value: `0` for numbers, `""` for strings, `false` for bools, `nil` for pointers/slices/maps/channels/interfaces.

### Type Conversion

Go requires explicit type conversions — there is no implicit casting:

```go
var i int = 42
var f float64 = float64(i)
var u uint = uint(f)

// String conversions
s := strconv.Itoa(42)            // int to string: "42"
n, err := strconv.Atoi("42")     // string to int: 42
f, err := strconv.ParseFloat("3.14", 64)  // string to float64

// Byte slice to/from string
b := []byte("hello")
s := string([]byte{'h', 'e', 'l', 'l', 'o'})
```

### Finding Types at Runtime

```go
import (
    "fmt"
    "reflect"
)

x := 10
fmt.Printf("Type: %T\n", x)                    // Type: int
fmt.Println("Type:", reflect.TypeOf(x))          // Type: int
```

---

## Variables & Constants

### Variable Declaration

```go
// Full declaration with type
var name string = "Go"

// Type inference
var version = 1.21

// Short declaration (inside functions only)
count := 42

// Multiple variables
var (
    host    = "localhost"
    port    = 8080
    debug   = false
)

// Multiple assignment
x, y := 10, 20
```

> **Important:** The Go compiler won't allow you to create variables that you never use. Unused variables cause a compilation error.

### Constants

```go
const pi = 3.14159
const greeting string = "Hello"

// Constant block
const (
    StatusOK    = 200
    StatusNotFound = 404
    StatusError = 500
)
```

### iota — Enumerated Constants

```go
type Weekday int

const (
    Sunday Weekday = iota  // 0
    Monday                  // 1
    Tuesday                 // 2
    Wednesday               // 3
    Thursday                // 4
    Friday                  // 5
    Saturday                // 6
)

// Bitmask pattern
type Permission uint8

const (
    Read    Permission = 1 << iota  // 1
    Write                            // 2
    Execute                          // 4
)

// Usage
perms := Read | Write  // 3
hasRead := perms&Read != 0  // true
```

---

## Control Flow

### if / else

Go's `if` does not require parentheses, but braces are mandatory. An `if` can include a short initialization statement:

```go
if x > 10 {
    fmt.Println("greater")
} else if x == 10 {
    fmt.Println("equal")
} else {
    fmt.Println("less")
}

// With initialization statement
if err := doSomething(); err != nil {
    log.Fatal(err)
}

// Common pattern: comma-ok idiom
if value, ok := myMap["key"]; ok {
    fmt.Println("Found:", value)
}
```

### for Loop

Go has only `for` — no `while` or `do-while`. It covers all looping patterns:

```go
// Classic for
for i := 0; i < 10; i++ {
    fmt.Println(i)
}

// While-style
i := 0
for i < 10 {
    fmt.Println(i)
    i++
}

// Infinite loop
for {
    // break when done
    if shouldStop() {
        break
    }
}

// Range over slice
fruits := []string{"apple", "banana", "cherry"}
for index, value := range fruits {
    fmt.Printf("%d: %s\n", index, value)
}

// Range over map
for key, value := range myMap {
    fmt.Printf("%s: %v\n", key, value)
}

// Range over string (iterates runes, not bytes)
for i, ch := range "Hello, 世界" {
    fmt.Printf("%d: %c\n", i, ch)
}

// Skip index or value with _
for _, value := range items {
    process(value)
}
```

### switch

```go
// No need for break — cases don't fall through by default
switch day {
case "Monday":
    fmt.Println("Start of week")
case "Friday":
    fmt.Println("Almost weekend")
case "Saturday", "Sunday":
    fmt.Println("Weekend!")
default:
    fmt.Println("Midweek")
}

// Switch with no condition (cleaner than if/else chains)
switch {
case score >= 90:
    grade = "A"
case score >= 80:
    grade = "B"
case score >= 70:
    grade = "C"
default:
    grade = "F"
}

// Type switch
switch v := value.(type) {
case string:
    fmt.Println("string:", v)
case int:
    fmt.Println("int:", v)
case bool:
    fmt.Println("bool:", v)
default:
    fmt.Printf("unknown type: %T\n", v)
}
```

### defer

`defer` schedules a function call to run after the surrounding function returns. Deferred calls execute in LIFO order.

```go
func readFile(path string) ([]byte, error) {
    f, err := os.Open(path)
    if err != nil {
        return nil, err
    }
    defer f.Close()  // guaranteed to run, even on panic

    return io.ReadAll(f)
}

// Multiple defers run in LIFO order
defer fmt.Println("first")   // prints third
defer fmt.Println("second")  // prints second
defer fmt.Println("third")   // prints first
```

**Advantages of defer:**

- Keeps cleanup close to resource acquisition
- Runs even if the function has multiple return paths
- Runs even if a runtime panic occurs

---

## Arrays, Slices & Maps

### Arrays (Fixed Length)

Arrays in Go have a fixed size that is part of the type. Arrays are rarely used directly — slices are preferred.

```go
var arr [5]int                    // [0, 0, 0, 0, 0]
arr[0] = 10

primes := [5]int{2, 3, 5, 7, 11}
auto := [...]int{1, 2, 3}        // compiler counts: [3]int
```

### Slices (Dynamic Length)

A slice is a dynamically-sized view into an underlying array. Slices are the workhorse of Go collections.

```go
// Create slices
s := []int{1, 2, 3, 4, 5}
s2 := make([]int, 5)             // length 5, capacity 5
s3 := make([]int, 0, 10)         // length 0, capacity 10

// Slice from array
arr := [5]int{1, 2, 3, 4, 5}
slice := arr[1:4]                // [2, 3, 4]

// Append (may allocate new underlying array)
s = append(s, 6)
s = append(s, 7, 8, 9)
s = append(s, otherSlice...)

// Copy
src := []int{1, 2, 3}
dst := make([]int, len(src))
copy(dst, src)

// Length and capacity
len(s)                           // number of elements
cap(s)                           // capacity of underlying array

// Nil vs empty slice
var nilSlice []int               // nil, len=0, cap=0
emptySlice := []int{}            // not nil, len=0, cap=0
// Both work with append, len, range

// Delete element (order preserved)
s = append(s[:i], s[i+1:]...)

// slices package (Go 1.21+)
import "slices"
slices.Sort(s)
slices.Contains(s, 42)
idx := slices.Index(s, 42)
```

### Maps

A map is an unordered collection of key-value pairs. Keys must be comparable types.

```go
// Create maps
m := map[string]int{
    "alice": 95,
    "bob":   87,
}
m2 := make(map[string]int)       // empty map
// var m3 map[string]int          // nil map — read OK, write panics!

// Access
score := m["alice"]              // 95
score = m["unknown"]             // 0 (zero value, no error)

// Check existence (comma-ok idiom)
score, ok := m["alice"]
if ok {
    fmt.Println("Found:", score)
}

// Set and delete
m["charlie"] = 92
delete(m, "bob")

// Iteration (order is random)
for key, value := range m {
    fmt.Printf("%s: %d\n", key, value)
}

// Length
fmt.Println(len(m))

// Nested maps
nested := map[string]map[string]string{
    "H": {
        "name":  "Hydrogen",
        "state": "gas",
    },
    "He": {
        "name":  "Helium",
        "state": "gas",
    },
}

// maps package (Go 1.21+)
import "maps"
keys := maps.Keys(m)
maps.Equal(m1, m2)
```

---

## Functions

### Basic Functions

```go
// Single return
func add(a, b int) int {
    return a + b
}

// Multiple returns
func divide(a, b float64) (float64, error) {
    if b == 0 {
        return 0, fmt.Errorf("division by zero")
    }
    return a / b, nil
}

result, err := divide(10, 3)

// Named returns
func swap(a, b int) (x, y int) {
    x = b
    y = a
    return  // naked return — returns x, y
}
```

### Variadic Functions

Functions can accept a variable number of arguments using `...` before the type:

```go
func sum(nums ...int) int {
    total := 0
    for _, n := range nums {
        total += n
    }
    return total
}

sum(1, 2, 3)                  // 6
sum(1, 2, 3, 4, 5)            // 15

// Pass a slice with ...
nums := []int{1, 2, 3}
sum(nums...)                   // 6
```

### Closures

Functions can capture variables from their enclosing scope:

```go
func counter() func() int {
    count := 0
    return func() int {
        count++
        return count
    }
}

c := counter()
fmt.Println(c())  // 1
fmt.Println(c())  // 2
fmt.Println(c())  // 3
```

### Recursion

```go
func factorial(n int) int {
    if n == 0 {
        return 1
    }
    return n * factorial(n-1)
}
```

### defer, panic & recover

```go
// panic terminates the program with a message
func mustParse(s string) int {
    n, err := strconv.Atoi(s)
    if err != nil {
        panic(fmt.Sprintf("failed to parse %q: %v", s, err))
    }
    return n
}

// recover catches panics in deferred functions
func safeDiv(a, b int) (result int, err error) {
    defer func() {
        if r := recover(); r != nil {
            err = fmt.Errorf("recovered: %v", r)
        }
    }()

    return a / b, nil
}
```

---

## Structs & Methods

### Struct Declaration

A struct is a composite type grouping named fields:

```go
type User struct {
    ID        int
    Name      string
    Email     string
    Active    bool
    CreatedAt time.Time
}

// Initialization
u1 := User{Name: "Alice", Email: "alice@example.com", Active: true}
u2 := User{ID: 1, Name: "Bob"}  // unset fields get zero values

// Access
fmt.Println(u1.Name)
u1.Active = false
```

### Methods (Receiver Functions)

Methods are functions with a receiver argument:

```go
// Value receiver — operates on a copy
func (u User) FullName() string {
    return u.Name
}

// Pointer receiver — can modify the original
func (u *User) Deactivate() {
    u.Active = false
}

// Pointer receiver — avoids copying large structs
func (u *User) UpdateEmail(email string) error {
    if !strings.Contains(email, "@") {
        return fmt.Errorf("invalid email: %s", email)
    }
    u.Email = email
    return nil
}

u := &User{Name: "Alice", Active: true}
u.Deactivate()
```

**When to use pointer receivers:**

- Method needs to modify the receiver
- Struct is large (avoids copying)
- Consistency — if any method uses a pointer receiver, all should

### Embedding (Composition)

Go doesn't have inheritance. Use embedding for composition:

```go
type Address struct {
    Street string
    City   string
    State  string
}

type Employee struct {
    User              // embedded — fields/methods promoted
    Address           // embedded
    Department string
}

emp := Employee{
    User:       User{Name: "Alice", Email: "alice@co.com"},
    Address:    Address{City: "Helsinki"},
    Department: "Engineering",
}

// Access promoted fields directly
fmt.Println(emp.Name)   // from User
fmt.Println(emp.City)   // from Address
```

### Constructor Pattern

Go doesn't have constructors. Use `New` functions:

```go
func NewUser(name, email string) (*User, error) {
    if name == "" {
        return nil, fmt.Errorf("name is required")
    }
    if !strings.Contains(email, "@") {
        return nil, fmt.Errorf("invalid email: %s", email)
    }
    return &User{
        Name:      name,
        Email:     email,
        Active:    true,
        CreatedAt: time.Now(),
    }, nil
}
```

---

## Interfaces

### Implicit Implementation

Interfaces in Go are satisfied implicitly — no `implements` keyword:

```go
type Shape interface {
    Area() float64
    Perimeter() float64
}

type Circle struct {
    Radius float64
}

func (c Circle) Area() float64 {
    return math.Pi * c.Radius * c.Radius
}

func (c Circle) Perimeter() float64 {
    return 2 * math.Pi * c.Radius
}

type Rectangle struct {
    Width, Height float64
}

func (r Rectangle) Area() float64 {
    return r.Width * r.Height
}

func (r Rectangle) Perimeter() float64 {
    return 2 * (r.Width + r.Height)
}

// Both Circle and Rectangle implement Shape
func printInfo(s Shape) {
    fmt.Printf("Area: %.2f, Perimeter: %.2f\n", s.Area(), s.Perimeter())
}

printInfo(Circle{Radius: 5})
printInfo(Rectangle{Width: 3, Height: 4})
```

### Empty Interface and any

```go
// any is an alias for interface{} (Go 1.18+)
func printValue(v any) {
    fmt.Printf("Value: %v (Type: %T)\n", v, v)
}

printValue(42)
printValue("hello")
printValue(true)
```

### Type Assertions and Type Switches

```go
// Type assertion
var i any = "hello"
s, ok := i.(string)
if ok {
    fmt.Println("String:", s)
}

// Type switch
func describe(i any) string {
    switch v := i.(type) {
    case string:
        return fmt.Sprintf("string of length %d", len(v))
    case int:
        return fmt.Sprintf("integer %d", v)
    case bool:
        return fmt.Sprintf("boolean %t", v)
    case nil:
        return "nil"
    default:
        return fmt.Sprintf("unknown type %T", v)
    }
}
```

### Common Standard Library Interfaces

```go
// io.Reader — anything that can be read from
type Reader interface {
    Read(p []byte) (n int, err error)
}

// io.Writer — anything that can be written to
type Writer interface {
    Write(p []byte) (n int, err error)
}

// fmt.Stringer — custom string representation
type Stringer interface {
    String() string
}

// error — the built-in error interface
type error interface {
    Error() string
}

// sort.Interface
type Interface interface {
    Len() int
    Less(i, j int) bool
    Swap(i, j int)
}
```

### Interface Best Practices

- Keep interfaces small (1-3 methods)
- Define interfaces where they are used, not where they are implemented
- Accept interfaces, return structs
- "The bigger the interface, the weaker the abstraction" — Rob Pike

---

## Pointers

Go is pass-by-value. Pointers allow functions to modify the original value and avoid copying large structures.

```go
// & takes the address, * dereferences
func increment(ptr *int) {
    *ptr++
}

x := 10
increment(&x)
fmt.Println(x)  // 11

// new() allocates and returns a pointer
ptr := new(int)   // *int, pointing to 0
*ptr = 42
```

### Value Types vs Reference Types

| Value Types | Reference Types |
|-------------|-----------------|
| `int`, `float`, `string`, `bool`, `struct` | `slice`, `map`, `channel`, `pointer`, `function` |
| Use pointers to modify in functions | Already passed by reference internally |

```go
// Slice — reference type, no pointer needed
func appendItem(s []string, item string) []string {
    return append(s, item)
}

// Struct — value type, use pointer to modify
func updateName(u *User, name string) {
    u.Name = name
}
```

---

## Packages & Modules

### Package Basics

Every Go file belongs to a package. Exported identifiers start with an uppercase letter:

```go
// math/calc.go
package math

// Exported — accessible from other packages
func Add(a, b int) int { return a + b }

// unexported — only accessible within this package
func helper() {}
```

### Go Modules

```bash
# Initialize a module
go mod init github.com/user/myproject

# Add a dependency
go get github.com/gin-gonic/gin@v1.9.1

# Tidy dependencies (remove unused, add missing)
go mod tidy

# Vendor dependencies
go mod vendor
```

### go.mod File

```go
module github.com/user/myproject

go 1.21

require (
    github.com/gin-gonic/gin v1.9.1
    go.uber.org/zap v1.26.0
)
```

### Project Structure

```text
myproject/
├── go.mod
├── go.sum
├── main.go                  # package main
├── cmd/
│   └── server/
│       └── main.go          # package main (alternative entrypoint)
├── internal/                # private packages (not importable externally)
│   ├── handler/
│   │   └── handler.go
│   └── service/
│       └── service.go
├── pkg/                     # public packages (importable)
│   └── client/
│       └── client.go
└── tests/
    └── integration_test.go
```

### Visibility

```go
// Uppercase first letter → exported (public)
func ProcessData() {}    // accessible from other packages
type Config struct {}    // accessible from other packages

// Lowercase first letter → unexported (private)
func helper() {}         // only accessible within this package
type config struct {}    // only accessible within this package
```

---

## Error Handling

### The error Interface

```go
// Go's error handling is explicit — no exceptions
type error interface {
    Error() string
}

result, err := doSomething()
if err != nil {
    return fmt.Errorf("doSomething failed: %w", err)  // wrap error
}
```

### Custom Errors

```go
// Sentinel errors
var (
    ErrNotFound     = errors.New("not found")
    ErrUnauthorized = errors.New("unauthorized")
    ErrTimeout      = errors.New("operation timed out")
)

// Structured errors
type ValidationError struct {
    Field   string
    Message string
}

func (e *ValidationError) Error() string {
    return fmt.Sprintf("validation error on %s: %s", e.Field, e.Message)
}

func validateAge(age int) error {
    if age < 0 || age > 150 {
        return &ValidationError{
            Field:   "age",
            Message: fmt.Sprintf("invalid age: %d", age),
        }
    }
    return nil
}
```

### Error Wrapping and Unwrapping (Go 1.13+)

```go
import "errors"

// Wrap errors with context
func getUser(id string) (*User, error) {
    user, err := db.FindUser(id)
    if err != nil {
        return nil, fmt.Errorf("getUser(%s): %w", id, err)
    }
    return user, nil
}

// Check error chain
if errors.Is(err, ErrNotFound) {
    // handle not found
}

// Extract specific error types
var validErr *ValidationError
if errors.As(err, &validErr) {
    fmt.Printf("Field: %s, Message: %s\n", validErr.Field, validErr.Message)
}

// Multiple error wrapping (Go 1.20+)
err := fmt.Errorf("operation failed: %w, also: %w", err1, err2)
```

### Error Handling Patterns

```go
// Early return pattern
func processOrder(orderID string) error {
    order, err := getOrder(orderID)
    if err != nil {
        return fmt.Errorf("get order: %w", err)
    }

    if err := validateOrder(order); err != nil {
        return fmt.Errorf("validate order: %w", err)
    }

    if err := chargePayment(order); err != nil {
        return fmt.Errorf("charge payment: %w", err)
    }

    return nil
}

// Must pattern (for initialization — panics on error)
func mustParseURL(raw string) *url.URL {
    u, err := url.Parse(raw)
    if err != nil {
        panic(fmt.Sprintf("invalid URL %q: %v", raw, err))
    }
    return u
}
```

---

## Goroutines & Channels

### Goroutines

A goroutine is a lightweight thread of execution managed by the Go runtime:

```go
// Start a goroutine with the go keyword
go func() {
    fmt.Println("running concurrently")
}()

go processItem(item)

// Goroutines are cheap — thousands are normal
for i := 0; i < 1000; i++ {
    go worker(i)
}
```

### Channels

Channels are typed conduits for communication between goroutines:

```go
// Unbuffered channel — sender blocks until receiver is ready
ch := make(chan string)

go func() {
    ch <- "hello"     // send
}()

msg := <-ch           // receive (blocks until value available)
fmt.Println(msg)      // "hello"

// Buffered channel — sender blocks only when buffer is full
ch := make(chan int, 10)

// Directional channels (for function signatures)
func producer(out chan<- int) {  // send-only
    out <- 42
}

func consumer(in <-chan int) {   // receive-only
    val := <-in
    fmt.Println(val)
}

// Close a channel to signal no more values
close(ch)

// Range over channel until closed
for msg := range ch {
    process(msg)
}
```

### select Statement

`select` lets a goroutine wait on multiple channel operations:

```go
select {
case msg := <-msgCh:
    fmt.Println("Received:", msg)
case err := <-errCh:
    fmt.Println("Error:", err)
case <-time.After(5 * time.Second):
    fmt.Println("Timeout")
case <-ctx.Done():
    fmt.Println("Cancelled")
}
```

### Done Pattern (Signaling Completion)

```go
func worker(done chan struct{}, jobs <-chan int) {
    defer close(done)
    for job := range jobs {
        process(job)
    }
}

done := make(chan struct{})
jobs := make(chan int, 100)

go worker(done, jobs)

// Send work
for i := 0; i < 50; i++ {
    jobs <- i
}
close(jobs)

// Wait for completion
<-done
```

### Fan-Out / Fan-In

```go
// Fan-out: distribute work across multiple goroutines
func fanOut(input <-chan int, workers int) []<-chan int {
    channels := make([]<-chan int, workers)
    for i := 0; i < workers; i++ {
        channels[i] = process(input)
    }
    return channels
}

// Fan-in: merge multiple channels into one
func fanIn(channels ...<-chan int) <-chan int {
    var wg sync.WaitGroup
    merged := make(chan int)

    for _, ch := range channels {
        wg.Add(1)
        go func(c <-chan int) {
            defer wg.Done()
            for val := range c {
                merged <- val
            }
        }(ch)
    }

    go func() {
        wg.Wait()
        close(merged)
    }()

    return merged
}
```

---

## Concurrency Patterns

### sync.WaitGroup

Wait for a collection of goroutines to finish:

```go
var wg sync.WaitGroup

for _, url := range urls {
    wg.Add(1)
    go func(u string) {
        defer wg.Done()
        fetch(u)
    }(url)
}

wg.Wait()  // blocks until all goroutines call Done()
```

### sync.Mutex

Protect shared state from concurrent access:

```go
type SafeCounter struct {
    mu sync.Mutex
    v  map[string]int
}

func (c *SafeCounter) Inc(key string) {
    c.mu.Lock()
    defer c.mu.Unlock()
    c.v[key]++
}

func (c *SafeCounter) Get(key string) int {
    c.mu.Lock()
    defer c.mu.Unlock()
    return c.v[key]
}

// RWMutex — allows concurrent reads
type Cache struct {
    mu   sync.RWMutex
    data map[string]string
}

func (c *Cache) Get(key string) (string, bool) {
    c.mu.RLock()
    defer c.mu.RUnlock()
    val, ok := c.data[key]
    return val, ok
}

func (c *Cache) Set(key, value string) {
    c.mu.Lock()
    defer c.mu.Unlock()
    c.data[key] = value
}
```

### sync.Once

Execute initialization exactly once, safely:

```go
var (
    instance *Database
    once     sync.Once
)

func GetDB() *Database {
    once.Do(func() {
        instance = connectDB()
    })
    return instance
}
```

### context.Context

Carry deadlines, cancellations, and values across API boundaries:

```go
import "context"

// With timeout
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()

// With cancellation
ctx, cancel := context.WithCancel(context.Background())
defer cancel()

// Use in functions
func fetchData(ctx context.Context, url string) ([]byte, error) {
    req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
    if err != nil {
        return nil, err
    }

    resp, err := http.DefaultClient.Do(req)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()

    return io.ReadAll(resp.Body)
}

// Check for cancellation in long operations
func process(ctx context.Context, items []Item) error {
    for _, item := range items {
        select {
        case <-ctx.Done():
            return ctx.Err()
        default:
            if err := handle(item); err != nil {
                return err
            }
        }
    }
    return nil
}
```

### errgroup

Run goroutines and collect errors:

```go
import "golang.org/x/sync/errgroup"

func fetchAll(ctx context.Context, urls []string) ([]string, error) {
    g, ctx := errgroup.WithContext(ctx)
    results := make([]string, len(urls))

    for i, url := range urls {
        i, url := i, url  // capture loop vars
        g.Go(func() error {
            body, err := fetchData(ctx, url)
            if err != nil {
                return fmt.Errorf("fetch %s: %w", url, err)
            }
            results[i] = string(body)
            return nil
        })
    }

    if err := g.Wait(); err != nil {
        return nil, err
    }
    return results, nil
}
```

---

## Testing

### Basic Tests

```go
// math_test.go
package math

import "testing"

func TestAdd(t *testing.T) {
    got := Add(2, 3)
    want := 5
    if got != want {
        t.Errorf("Add(2, 3) = %d, want %d", got, want)
    }
}
```

### Table-Driven Tests

The standard Go testing pattern:

```go
func TestAdd(t *testing.T) {
    tests := []struct {
        name string
        a, b int
        want int
    }{
        {"positive", 2, 3, 5},
        {"negative", -1, -1, -2},
        {"zero", 0, 0, 0},
        {"mixed", -1, 1, 0},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got := Add(tt.a, tt.b)
            if got != tt.want {
                t.Errorf("Add(%d, %d) = %d, want %d", tt.a, tt.b, got, tt.want)
            }
        })
    }
}
```

### Benchmarks

```go
func BenchmarkAdd(b *testing.B) {
    for i := 0; i < b.N; i++ {
        Add(100, 200)
    }
}

// Run benchmarks
// go test -bench=. -benchmem
```

### Test Helpers and Mocks

```go
// Test helper
func setupTestDB(t *testing.T) *sql.DB {
    t.Helper()  // marks as helper — errors report caller's line
    db, err := sql.Open("sqlite3", ":memory:")
    if err != nil {
        t.Fatal(err)
    }
    t.Cleanup(func() {
        db.Close()
    })
    return db
}

// Interface-based mocking
type UserRepository interface {
    GetUser(id string) (*User, error)
}

type mockRepo struct {
    users map[string]*User
}

func (m *mockRepo) GetUser(id string) (*User, error) {
    u, ok := m.users[id]
    if !ok {
        return nil, ErrNotFound
    }
    return u, nil
}

func TestService(t *testing.T) {
    repo := &mockRepo{
        users: map[string]*User{
            "1": {Name: "Alice"},
        },
    }
    svc := NewService(repo)
    // test svc...
}
```

### Running Tests

```bash
go test ./...                    # all tests recursively
go test -v ./pkg/...             # verbose
go test -run TestAdd             # specific test
go test -count=1 ./...           # bypass test cache
go test -race ./...              # detect race conditions
go test -cover ./...             # show coverage
go test -coverprofile=cover.out  # coverage profile
go tool cover -html=cover.out    # view in browser
go test -bench=. -benchmem       # run benchmarks
```

---

## Common Patterns

### Functional Options

A clean pattern for configurable constructors:

```go
type Server struct {
    host    string
    port    int
    timeout time.Duration
    logger  *log.Logger
}

type Option func(*Server)

func WithPort(port int) Option {
    return func(s *Server) {
        s.port = port
    }
}

func WithTimeout(d time.Duration) Option {
    return func(s *Server) {
        s.timeout = d
    }
}

func WithLogger(l *log.Logger) Option {
    return func(s *Server) {
        s.logger = l
    }
}

func NewServer(host string, opts ...Option) *Server {
    s := &Server{
        host:    host,
        port:    8080,                   // defaults
        timeout: 30 * time.Second,
        logger:  log.Default(),
    }
    for _, opt := range opts {
        opt(s)
    }
    return s
}

// Usage
srv := NewServer("localhost",
    WithPort(9090),
    WithTimeout(10*time.Second),
)
```

### Builder Pattern

```go
type QueryBuilder struct {
    table      string
    conditions []string
    orderBy    string
    limit      int
}

func NewQuery(table string) *QueryBuilder {
    return &QueryBuilder{table: table}
}

func (qb *QueryBuilder) Where(condition string) *QueryBuilder {
    qb.conditions = append(qb.conditions, condition)
    return qb
}

func (qb *QueryBuilder) OrderBy(field string) *QueryBuilder {
    qb.orderBy = field
    return qb
}

func (qb *QueryBuilder) Limit(n int) *QueryBuilder {
    qb.limit = n
    return qb
}

func (qb *QueryBuilder) Build() string {
    query := fmt.Sprintf("SELECT * FROM %s", qb.table)
    if len(qb.conditions) > 0 {
        query += " WHERE " + strings.Join(qb.conditions, " AND ")
    }
    if qb.orderBy != "" {
        query += " ORDER BY " + qb.orderBy
    }
    if qb.limit > 0 {
        query += fmt.Sprintf(" LIMIT %d", qb.limit)
    }
    return query
}

// Usage
query := NewQuery("users").
    Where("age > 18").
    Where("active = true").
    OrderBy("name").
    Limit(10).
    Build()
```

---

## CLI Tools

| Command | Description |
|---------|-------------|
| `go run main.go` | Compile and execute |
| `go build` | Compile packages and dependencies |
| `go build -o myapp` | Compile to a named binary |
| `go test ./...` | Run all tests |
| `go test -v -race ./...` | Verbose tests with race detection |
| `go test -bench=.` | Run benchmarks |
| `go fmt ./...` | Format all code |
| `go vet ./...` | Report suspicious constructs |
| `go mod init <path>` | Initialize a new module |
| `go mod tidy` | Add missing / remove unused deps |
| `go mod vendor` | Copy deps into vendor directory |
| `go get <pkg>@<ver>` | Add or update a dependency |
| `go doc <pkg>` | View package documentation |
| `go doc <pkg> <sym>` | View symbol documentation |
| `go install <pkg>@latest` | Install a binary |
| `go generate ./...` | Run go:generate directives |
| `go version` | Print Go version |

### Static Analysis

```bash
# Built-in vet
go vet ./...

# golangci-lint (aggregated linter)
golangci-lint run

# staticcheck
staticcheck ./...

# govulncheck (vulnerability scanner)
govulncheck ./...
```

### Cross-Compilation

```bash
# Build for Linux from macOS
GOOS=linux GOARCH=amd64 go build -o myapp-linux

# Build for Windows
GOOS=windows GOARCH=amd64 go build -o myapp.exe

# Build for ARM (e.g., Raspberry Pi)
GOOS=linux GOARCH=arm64 go build -o myapp-arm
```

---

## Common Mistakes

### Nil Pointer Dereference

```go
// ❌ WRONG — accessing fields on nil pointer
var u *User
fmt.Println(u.Name)  // panic: nil pointer dereference

// ✅ CORRECT — check before using
if u != nil {
    fmt.Println(u.Name)
}

// ✅ CORRECT — return early on error
u, err := getUser(id)
if err != nil {
    return err
}
fmt.Println(u.Name)  // u is guaranteed non-nil here
```

### Goroutine Leaks

```go
// ❌ WRONG — goroutine blocks forever if no one reads
func leak() {
    ch := make(chan int)
    go func() {
        ch <- 42  // blocks forever — no reader
    }()
    // function returns, goroutine leaks
}

// ✅ CORRECT — use context for cancellation
func noLeak(ctx context.Context) {
    ch := make(chan int, 1)  // buffered so goroutine won't block
    go func() {
        select {
        case ch <- 42:
        case <-ctx.Done():
            return
        }
    }()
}
```

### Race Conditions

```go
// ❌ WRONG — concurrent map writes
m := make(map[string]int)
for i := 0; i < 100; i++ {
    go func() {
        m["key"]++  // DATA RACE — panic: concurrent map writes
    }()
}

// ✅ CORRECT — use sync.Mutex
var mu sync.Mutex
m := make(map[string]int)
for i := 0; i < 100; i++ {
    go func() {
        mu.Lock()
        m["key"]++
        mu.Unlock()
    }()
}

// ✅ CORRECT — or use sync.Map for concurrent map access
var m sync.Map
m.Store("key", 42)
val, ok := m.Load("key")
```

### Loop Variable Capture

```go
// ❌ WRONG (before Go 1.22) — all goroutines share same variable
for _, url := range urls {
    go func() {
        fetch(url)  // all goroutines see the last url!
    }()
}

// ✅ CORRECT — capture the variable
for _, url := range urls {
    url := url  // shadow the loop variable
    go func() {
        fetch(url)
    }()
}

// ✅ Go 1.22+ — loop variables are per-iteration by default
```

### Not Handling Errors

```go
// ❌ WRONG — ignoring errors
result, _ := json.Marshal(data)

// ✅ CORRECT — always handle errors
result, err := json.Marshal(data)
if err != nil {
    return fmt.Errorf("marshal failed: %w", err)
}
```

### Nil Map Write

```go
// ❌ WRONG — writing to nil map panics
var m map[string]int
m["key"] = 1  // panic: assignment to entry in nil map

// ✅ CORRECT — initialize with make
m := make(map[string]int)
m["key"] = 1
```

---

## Production Tips

### Structured Logging

```go
import "go.uber.org/zap"

logger, _ := zap.NewProduction()
defer logger.Sync()

logger.Info("server started",
    zap.String("host", host),
    zap.Int("port", port),
)

logger.Error("request failed",
    zap.String("method", r.Method),
    zap.String("path", r.URL.Path),
    zap.Error(err),
)

// slog (Go 1.21+ standard library)
import "log/slog"

slog.Info("server started", "host", host, "port", port)
slog.Error("request failed", "err", err, "path", r.URL.Path)
```

### Graceful Shutdown

```go
func main() {
    srv := &http.Server{Addr: ":8080", Handler: mux}

    go func() {
        if err := srv.ListenAndServe(); err != http.ErrServerClosed {
            log.Fatal(err)
        }
    }()

    // Wait for interrupt signal
    quit := make(chan os.Signal, 1)
    signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
    <-quit

    log.Println("Shutting down...")
    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()

    if err := srv.Shutdown(ctx); err != nil {
        log.Fatal("Forced shutdown:", err)
    }
    log.Println("Server stopped")
}
```

### Health Checks

```go
mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
    w.WriteHeader(http.StatusOK)
    w.Write([]byte("ok"))
})

mux.HandleFunc("/readyz", func(w http.ResponseWriter, r *http.Request) {
    if err := db.Ping(); err != nil {
        w.WriteHeader(http.StatusServiceUnavailable)
        w.Write([]byte("not ready"))
        return
    }
    w.WriteHeader(http.StatusOK)
    w.Write([]byte("ready"))
})
```

### Build Flags

```bash
# Embed version info at build time
go build -ldflags "-X main.version=1.0.0 -X main.commit=$(git rev-parse HEAD)" -o myapp

# Strip debug info for smaller binary
go build -ldflags "-s -w" -o myapp

# Static binary (for scratch Docker images)
CGO_ENABLED=0 go build -o myapp
```

### Docker

```dockerfile
# Multi-stage build
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -ldflags="-s -w" -o /app/server ./cmd/server

FROM scratch
COPY --from=builder /app/server /server
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
EXPOSE 8080
ENTRYPOINT ["/server"]
```

---

## References

### Official Resources

- [Go Documentation](https://go.dev/doc/) — Official docs
- [Go Tour](https://go.dev/tour/) — Interactive tutorial
- [Go Playground](https://go.dev/play/) — Run Go in the browser
- [Go Standard Library](https://pkg.go.dev/std) — Package reference
- [Effective Go](https://go.dev/doc/effective_go) — Idiomatic Go guide
- [Go Blog](https://go.dev/blog/) — Official blog

### Books

- *The Go Programming Language* by Alan Donovan & Brian Kernighan
- *Go in Action* by William Kennedy, Brian Ketelsen & Erik St. Martin
- *Learning Go* (2nd ed.) by Jon Bodner
- *Concurrency in Go* by Katherine Cox-Buday
- *Go Web Programming* by Sau Sheong Chang

### Tools

- [golangci-lint](https://golangci-lint.run/) — Aggregated linter
- [staticcheck](https://staticcheck.io/) — Advanced static analysis
- [govulncheck](https://pkg.go.dev/golang.org/x/vuln/cmd/govulncheck) — Vulnerability scanner
- [delve](https://github.com/go-delve/delve) — Go debugger

### Related Topics

- [🧰 CLI / Docker](../../cli/docker/) — Containerizing Go apps
- [🧰 CLI / kubectl](../../cli/kubectl/) — Kubernetes deployments
- [🚀 DevOps](../../devops/) — CI/CD for Go projects
- [⚙️ Engineering / Testing](../../engineering/) — Testing best practices
- [🏗 System Design](../../system-design/) — Architecture patterns

---

*Part of [Tech-Byte Languages](../). Last updated: 2026-08.*
