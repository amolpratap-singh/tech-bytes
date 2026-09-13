# 💻 Programming Languages

> Comprehensive language references with concepts, practical code, examples, common mistakes, and production tips.

---

## 📚 Language Index

| Language | Use Cases | Key Features | Status |
|----------|-----------|--------------|--------|
| [Python](python/) | Web backends, data science, scripting, automation, ML/AI | Dynamic typing, rich ecosystem, readability, rapid prototyping | ✅ Complete |
| [Go](go/) | Microservices, CLI tools, cloud infrastructure, networking | Static typing, goroutines, fast compilation, single binary | ✅ Complete |
| [Bash](bash/) | Shell scripting, automation, CI/CD pipelines, sysadmin | Ubiquitous on Unix, pipes, process control, text processing | ✅ Complete |
| Java | Enterprise backends, Android, distributed systems | JVM, strong typing, mature ecosystem, Spring framework | 🔜 Planned |
| JavaScript | Frontend, Node.js backends, full-stack web | Event loop, async/await, V8 engine, npm ecosystem | 🔜 Planned |
| TypeScript | Type-safe JavaScript, large-scale web apps | Static types, interfaces, generics, tooling support | 🔜 Planned |
| Rust | Systems programming, WebAssembly, performance-critical | Memory safety, zero-cost abstractions, ownership model | 🔜 Planned |

---

## 🗺️ How to Use These References

Each language reference follows the same structure:

1. **Overview** — What is it, why use it, when to choose it
2. **Core Language** — Types, variables, control flow, functions
3. **Advanced Features** — OOP, concurrency, generics, patterns
4. **Error Handling** — Idiomatic error handling for the language
5. **Testing** — Testing frameworks, patterns, and best practices
6. **Common Mistakes** — Gotchas and pitfalls to avoid
7. **Production Tips** — Logging, profiling, packaging, deployment
8. **References** — Books, docs, and learning resources

---

## ⚡ Quick Comparison

### Hello World

**Python:**

```python
print("Hello, World!")
```

**Go:**

```go
package main

import "fmt"

func main() {
    fmt.Println("Hello, World!")
}
```

**Bash:**

```bash
#!/bin/bash
echo "Hello, World!"
```

---

### Error Handling

**Python** — Exceptions:

```python
try:
    result = risky_operation()
except ValueError as e:
    logger.error("Invalid value: %s", e)
```

**Go** — Explicit error returns:

```go
result, err := riskyOperation()
if err != nil {
    log.Printf("operation failed: %v", err)
    return err
}
```

**Bash** — Exit codes:

```bash
if ! risky_command; then
    echo "Command failed with exit code $?" >&2
    exit 1
fi
```

---

### Concurrency

| Feature | Python | Go | Bash |
|---------|--------|----|------|
| Model | asyncio / threading / multiprocessing | Goroutines + Channels | Background processes (`&`) |
| Parallelism | multiprocessing (GIL limits threads) | Native (GOMAXPROCS) | Multiple processes |
| Communication | Queue, asyncio.Queue | Channels | Pipes, signals, files |
| Typical Use | I/O-bound async, CPU-bound multiprocessing | High-concurrency servers | Parallel shell tasks |

---

## 🔗 Related Topics

- [🧰 CLI Tools](../cli/) — Command-line tool references
- [🚀 DevOps](../devops/) — Docker, Kubernetes, CI/CD
- [🧩 Frameworks](../frameworks/) — Flask, FastAPI, Django, Spring, Node.js
- [⚙️ Engineering](../engineering/) — Testing, logging, code review
- [🎯 Learning Paths](../learning/) — Guided learning paths

---

*Each language reference is a living document. Contributions welcome — see [CONTRIBUTING.md](../CONTRIBUTING.md).*
