# 🐍 Python Reference

> Comprehensive Python reference — from fundamentals to production patterns. Focused on Python 3.10+.

---

## Table of Contents

1. [What is Python](#what-is-python)
2. [Data Types](#data-types)
3. [Control Flow](#control-flow)
4. [Functions](#functions)
5. [Object-Oriented Programming](#object-oriented-programming)
6. [Decorators](#decorators)
7. [Generators & Iterators](#generators--iterators)
8. [Context Managers](#context-managers)
9. [Error Handling](#error-handling)
10. [Type Hints](#type-hints)
11. [Async / Await](#async--await)
12. [Virtual Environments & Packaging](#virtual-environments--packaging)
13. [Testing](#testing)
14. [Common Patterns](#common-patterns)
15. [Common Mistakes & Gotchas](#common-mistakes--gotchas)
16. [Production Tips](#production-tips)
17. [Related Topics](#related-topics)
18. [References](#references)

---

## What is Python

Python is a high-level, interpreted, dynamically typed programming language emphasizing readability and developer productivity. Created by Guido van Rossum (1991), it's used across web backends, data science, machine learning, automation, scripting, and DevOps.

### Why Python

- **Readable syntax** — code reads close to English
- **Massive ecosystem** — PyPI has 500k+ packages
- **Rapid prototyping** — less boilerplate than compiled languages
- **Cross-domain** — web, data, ML, automation, infrastructure
- **Strong community** — extensive documentation, tutorials, Stack Overflow answers

### Python Versions (Focus: 3.10+)

| Version | Key Features |
|---------|-------------|
| 3.10 | Structural pattern matching (`match/case`), better error messages |
| 3.11 | 10-25% faster CPython, exception groups, `tomllib` |
| 3.12 | Per-interpreter GIL (experimental), f-string improvements, `type` statement |
| 3.13 | Free-threaded mode (experimental), JIT compiler (experimental) |

```bash
# Check version
python3 --version

# Use specific version
python3.12 -m venv .venv
```

---

## Data Types

### Strings

```python
# String creation
name = "Python"
multiline = """This is
a multiline string"""
raw = r"No \n escape here"

# f-strings (3.6+)
version = 3.12
greeting = f"Hello {name} {version}"

# Common operations
"hello".upper()                  # "HELLO"
"Hello World".split()            # ["Hello", "World"]
", ".join(["a", "b", "c"])       # "a, b, c"
"  spaces  ".strip()             # "spaces"
"hello world".replace("world", "python")  # "hello python"
"python" in "I love python"     # True

# String formatting
f"{3.14159:.2f}"                 # "3.14"
f"{'hello':>20}"                 # "               hello"
f"{1000000:,}"                   # "1,000,000"

# Slicing
s = "Hello, World"
s[0:5]     # "Hello"
s[-5:]     # "World"
s[::-1]    # "dlroW ,olleH" (reversed)
```

### Lists

```python
# Creation
numbers = [1, 2, 3, 4, 5]
mixed = [1, "two", 3.0, True]
nested = [[1, 2], [3, 4]]

# Operations
numbers.append(6)               # [1, 2, 3, 4, 5, 6]
numbers.extend([7, 8])          # [1, 2, 3, 4, 5, 6, 7, 8]
numbers.insert(0, 0)            # [0, 1, 2, 3, ...]
numbers.pop()                   # removes and returns last element
numbers.remove(3)               # removes first occurrence of 3

# Slicing
numbers[1:4]                    # elements 1, 2, 3
numbers[::2]                    # every other element
numbers[::-1]                   # reversed

# Sorting
sorted_nums = sorted(numbers)               # returns new list
numbers.sort()                               # sorts in place
sorted(items, key=lambda x: x["name"])       # sort by key

# Unpacking
first, *rest = [1, 2, 3, 4, 5]  # first=1, rest=[2, 3, 4, 5]
first, *_, last = [1, 2, 3, 4]  # first=1, last=4
```

### Tuples

```python
# Immutable sequences
point = (3, 4)
single = (42,)                  # trailing comma for single-element tuple
x, y = point                    # tuple unpacking

# Named tuples
from collections import namedtuple
Point = namedtuple("Point", ["x", "y"])
p = Point(3, 4)
print(p.x, p.y)                # 3 4

# As dict keys (hashable, unlike lists)
cache = {(0, 0): "origin", (1, 1): "diagonal"}
```

### Dictionaries

```python
# Creation
user = {"name": "Alice", "age": 30, "active": True}
from_keys = dict.fromkeys(["a", "b", "c"], 0)  # {"a": 0, "b": 0, "c": 0}

# Access
user["name"]                    # "Alice" (KeyError if missing)
user.get("email", "N/A")       # "N/A" (default if missing)

# Modification
user["email"] = "alice@example.com"
user.update({"age": 31, "role": "admin"})
user.pop("active")              # removes and returns value
user.setdefault("score", 100)   # sets only if key missing

# Iteration
for key, value in user.items():
    print(f"{key}: {value}")

# Merge (3.9+)
defaults = {"theme": "dark", "lang": "en"}
overrides = {"lang": "fi", "debug": True}
config = defaults | overrides   # {"theme": "dark", "lang": "fi", "debug": True}

# Dictionary comprehension
squares = {n: n**2 for n in range(10)}
```

### Sets

```python
# Unique, unordered collections
fruits = {"apple", "banana", "cherry"}
empty_set = set()               # not {} — that's an empty dict

# Operations
fruits.add("date")
fruits.discard("banana")        # no error if missing
fruits.remove("apple")          # KeyError if missing

# Set math
a = {1, 2, 3, 4}
b = {3, 4, 5, 6}
a | b                           # union: {1, 2, 3, 4, 5, 6}
a & b                           # intersection: {3, 4}
a - b                           # difference: {1, 2}
a ^ b                           # symmetric difference: {1, 2, 5, 6}

# Membership testing (O(1) average)
if "apple" in fruits:
    print("Found it")
```

---

## Control Flow

### Conditionals

```python
# if / elif / else
status_code = 404

if status_code == 200:
    print("OK")
elif status_code == 404:
    print("Not Found")
elif 500 <= status_code < 600:
    print("Server Error")
else:
    print(f"Status: {status_code}")

# Ternary operator
result = "even" if x % 2 == 0 else "odd"

# Walrus operator (3.8+)
if (n := len(data)) > 10:
    print(f"Processing {n} items")
```

### Pattern Matching (3.10+)

```python
match command.split():
    case ["quit"]:
        quit_game()
    case ["go", direction]:
        move(direction)
    case ["get", item] if item in valid_items:
        pick_up(item)
    case _:
        print("Unknown command")

# Matching types and structures
match event:
    case {"type": "click", "x": x, "y": y}:
        handle_click(x, y)
    case {"type": "keypress", "key": str(key)}:
        handle_key(key)
    case _:
        pass
```

### Loops

```python
# for loop
for item in collection:
    process(item)

# for with index
for i, item in enumerate(collection, start=1):
    print(f"{i}. {item}")

# for with zip
for name, score in zip(names, scores):
    print(f"{name}: {score}")

# while loop
while queue:
    item = queue.pop(0)
    process(item)

# Loop control
for n in range(100):
    if n % 2 == 0:
        continue            # skip even numbers
    if n > 50:
        break               # exit loop
    print(n)

# for/else — else runs if loop completed without break
for item in items:
    if item.matches(query):
        result = item
        break
else:
    result = None           # no match found
```

### Comprehensions

```python
# List comprehension
squares = [x**2 for x in range(10)]
evens = [x for x in numbers if x % 2 == 0]

# Dict comprehension
word_lengths = {word: len(word) for word in words}

# Set comprehension
unique_lengths = {len(word) for word in words}

# Nested comprehension
matrix = [[1, 2, 3], [4, 5, 6], [7, 8, 9]]
flat = [val for row in matrix for val in row]  # [1, 2, 3, 4, 5, 6, 7, 8, 9]

# Generator expression (lazy — uses parentheses)
total = sum(x**2 for x in range(1000000))  # memory efficient
```

---

## Functions

### Basic Functions

```python
def greet(name: str, greeting: str = "Hello") -> str:
    """Return a greeting string.

    Args:
        name: The person's name.
        greeting: The greeting word. Defaults to "Hello".

    Returns:
        A formatted greeting string.
    """
    return f"{greeting}, {name}!"

# Calling
greet("Alice")                       # "Hello, Alice!"
greet("Bob", greeting="Hi")          # "Hi, Bob!"
```

### *args and **kwargs

```python
def log_event(event: str, *tags: str, **metadata: str) -> dict:
    """Log an event with optional tags and metadata."""
    return {
        "event": event,
        "tags": list(tags),
        "metadata": metadata,
    }

log_event("deploy", "prod", "v2.1", author="alice", ticket="JIRA-123")
# {"event": "deploy", "tags": ["prod", "v2.1"], "metadata": {"author": "alice", "ticket": "JIRA-123"}}
```

### Keyword-Only and Positional-Only Arguments

```python
# Keyword-only (after *)
def connect(host: str, port: int, *, timeout: int = 30, ssl: bool = True):
    ...

connect("db.example.com", 5432, timeout=10)  # OK
connect("db.example.com", 5432, 10)           # TypeError

# Positional-only (before /) — Python 3.8+
def pow(base: float, exp: float, /) -> float:
    return base ** exp

pow(2, 10)        # OK
pow(base=2, exp=10)  # TypeError
```

### Lambda Functions

```python
# Short anonymous functions
square = lambda x: x**2

# Common use: sorting keys
users.sort(key=lambda u: u["last_login"])

# Common use: filtering
active = list(filter(lambda u: u["active"], users))

# Prefer named functions for anything non-trivial
```

### First-Class Functions

```python
def apply(func, value):
    """Apply a function to a value."""
    return func(value)

from functools import partial, reduce

# Partial application
double = partial(lambda x, y: x * y, 2)
double(5)  # 10

# reduce
total = reduce(lambda acc, x: acc + x, [1, 2, 3, 4], 0)  # 10
```

---

## Object-Oriented Programming

### Classes

```python
class User:
    """Represents a system user."""

    # Class variable (shared across instances)
    _instance_count: int = 0

    def __init__(self, name: str, email: str) -> None:
        self.name = name          # instance variable
        self.email = email
        self._active = True       # convention: private
        User._instance_count += 1

    def deactivate(self) -> None:
        """Deactivate the user account."""
        self._active = False

    @property
    def is_active(self) -> bool:
        """Check if user is active."""
        return self._active

    @classmethod
    def get_count(cls) -> int:
        """Return total number of User instances created."""
        return cls._instance_count

    @staticmethod
    def validate_email(email: str) -> bool:
        """Basic email validation."""
        return "@" in email and "." in email

    def __repr__(self) -> str:
        return f"User(name={self.name!r}, email={self.email!r})"

    def __str__(self) -> str:
        return f"{self.name} <{self.email}>"

    def __eq__(self, other: object) -> bool:
        if not isinstance(other, User):
            return NotImplemented
        return self.email == other.email

    def __hash__(self) -> int:
        return hash(self.email)
```

### Inheritance

```python
class AdminUser(User):
    """User with administrative privileges."""

    def __init__(self, name: str, email: str, permissions: list[str]) -> None:
        super().__init__(name, email)
        self.permissions = permissions

    def has_permission(self, perm: str) -> bool:
        return perm in self.permissions

    def __repr__(self) -> str:
        return f"AdminUser(name={self.name!r}, permissions={self.permissions!r})"
```

### Abstract Base Classes

```python
from abc import ABC, abstractmethod

class Repository(ABC):
    """Abstract repository interface."""

    @abstractmethod
    def get(self, id: str) -> dict:
        ...

    @abstractmethod
    def save(self, entity: dict) -> None:
        ...

    @abstractmethod
    def delete(self, id: str) -> None:
        ...

class PostgresRepository(Repository):
    def get(self, id: str) -> dict:
        # actual implementation
        ...

    def save(self, entity: dict) -> None:
        ...

    def delete(self, id: str) -> None:
        ...
```

### Dataclasses

```python
from dataclasses import dataclass, field
from datetime import datetime

@dataclass
class Config:
    """Application configuration."""
    host: str
    port: int = 8080
    debug: bool = False
    tags: list[str] = field(default_factory=list)
    created_at: datetime = field(default_factory=datetime.now)

    def __post_init__(self) -> None:
        if self.port < 1 or self.port > 65535:
            raise ValueError(f"Invalid port: {self.port}")

# Immutable dataclass
@dataclass(frozen=True)
class Point:
    x: float
    y: float

config = Config(host="localhost", tags=["dev", "test"])
p = Point(3.0, 4.0)
# p.x = 5.0  # FrozenInstanceError
```

---

## Decorators

### Function Decorators

```python
import functools
import time

def timer(func):
    """Log execution time of the decorated function."""
    @functools.wraps(func)
    def wrapper(*args, **kwargs):
        start = time.perf_counter()
        result = func(*args, **kwargs)
        elapsed = time.perf_counter() - start
        print(f"{func.__name__} took {elapsed:.4f}s")
        return result
    return wrapper

@timer
def fetch_data(url: str) -> dict:
    ...
```

### Decorators with Arguments

```python
def retry(max_attempts: int = 3, delay: float = 1.0):
    """Retry a function on failure."""
    def decorator(func):
        @functools.wraps(func)
        def wrapper(*args, **kwargs):
            last_exception = None
            for attempt in range(1, max_attempts + 1):
                try:
                    return func(*args, **kwargs)
                except Exception as e:
                    last_exception = e
                    print(f"Attempt {attempt}/{max_attempts} failed: {e}")
                    if attempt < max_attempts:
                        time.sleep(delay)
            raise last_exception
        return wrapper
    return decorator

@retry(max_attempts=5, delay=2.0)
def call_external_api(endpoint: str) -> dict:
    ...
```

### Practical Decorator Patterns

```python
# Caching (use functools.lru_cache for simple cases)
from functools import lru_cache

@lru_cache(maxsize=128)
def expensive_computation(n: int) -> int:
    return sum(i * i for i in range(n))

# Require authentication
def require_auth(func):
    @functools.wraps(func)
    def wrapper(request, *args, **kwargs):
        if not request.user.is_authenticated:
            raise PermissionError("Authentication required")
        return func(request, *args, **kwargs)
    return wrapper

# Singleton via decorator
def singleton(cls):
    instances = {}
    @functools.wraps(cls)
    def get_instance(*args, **kwargs):
        if cls not in instances:
            instances[cls] = cls(*args, **kwargs)
        return instances[cls]
    return get_instance
```

---

## Generators & Iterators

### Generators

```python
def fibonacci(limit: int):
    """Generate Fibonacci numbers up to limit."""
    a, b = 0, 1
    while a < limit:
        yield a
        a, b = b, a + b

# Lazy evaluation — memory efficient
for num in fibonacci(1000):
    print(num)

# Generator as a pipeline
def read_lines(filepath: str):
    with open(filepath) as f:
        for line in f:
            yield line.strip()

def filter_errors(lines):
    for line in lines:
        if "ERROR" in line:
            yield line

def extract_timestamps(lines):
    for line in lines:
        yield line.split()[0]

# Compose pipelines
log_lines = read_lines("/var/log/app.log")
errors = filter_errors(log_lines)
timestamps = extract_timestamps(errors)
```

### Generator Expressions

```python
# Like list comprehension but lazy
squares = (x**2 for x in range(1_000_000))

# Memory comparison
import sys
list_comp = [x**2 for x in range(1_000_000)]
gen_exp = (x**2 for x in range(1_000_000))
print(sys.getsizeof(list_comp))  # ~8 MB
print(sys.getsizeof(gen_exp))    # ~200 bytes
```

### itertools

```python
import itertools

# Chain multiple iterables
combined = itertools.chain([1, 2], [3, 4], [5, 6])  # 1, 2, 3, 4, 5, 6

# Group consecutive equal elements
data = [("a", 1), ("a", 2), ("b", 3), ("b", 4)]
for key, group in itertools.groupby(data, key=lambda x: x[0]):
    print(key, list(group))

# Cartesian product
for combo in itertools.product("AB", range(3)):
    print(combo)  # ("A", 0), ("A", 1), ("A", 2), ("B", 0), ...

# Sliding window (3.12+)
from itertools import pairwise
for a, b in pairwise([1, 2, 3, 4]):
    print(a, b)  # (1, 2), (2, 3), (3, 4)

# Batched (3.12+)
from itertools import batched
for batch in batched(range(10), 3):
    print(batch)  # (0, 1, 2), (3, 4, 5), (6, 7, 8), (9,)
```

---

## Context Managers

### Using with Statement

```python
# File handling
with open("data.txt", "r") as f:
    content = f.read()

# Multiple context managers
with open("input.txt") as src, open("output.txt", "w") as dst:
    dst.write(src.read())

# Database connection
with psycopg2.connect(dsn) as conn:
    with conn.cursor() as cur:
        cur.execute("SELECT * FROM users")
        rows = cur.fetchall()
```

### Custom Context Manager (Class-Based)

```python
class Timer:
    """Context manager that measures execution time."""

    def __enter__(self):
        self.start = time.perf_counter()
        return self

    def __exit__(self, exc_type, exc_val, exc_tb):
        self.elapsed = time.perf_counter() - self.start
        print(f"Elapsed: {self.elapsed:.4f}s")
        return False  # don't suppress exceptions

with Timer() as t:
    heavy_computation()
print(f"Took {t.elapsed:.2f}s")
```

### contextlib

```python
from contextlib import contextmanager, suppress

@contextmanager
def temp_directory():
    """Create a temporary directory and clean up after use."""
    import tempfile
    import shutil
    dirpath = tempfile.mkdtemp()
    try:
        yield dirpath
    finally:
        shutil.rmtree(dirpath)

with temp_directory() as tmpdir:
    # use tmpdir...
    pass  # cleaned up automatically

# Suppress specific exceptions
with suppress(FileNotFoundError):
    os.remove("maybe_exists.tmp")
```

---

## Error Handling

### try / except / else / finally

```python
def parse_config(filepath: str) -> dict:
    """Parse a JSON configuration file."""
    try:
        with open(filepath) as f:
            config = json.load(f)
    except FileNotFoundError:
        print(f"Config file not found: {filepath}")
        return {}
    except json.JSONDecodeError as e:
        print(f"Invalid JSON in {filepath}: {e}")
        return {}
    else:
        # Runs only if no exception occurred
        print(f"Loaded config from {filepath}")
        return config
    finally:
        # Always runs — cleanup goes here
        print("Config parsing complete")
```

### Custom Exceptions

```python
class AppError(Exception):
    """Base exception for the application."""
    pass

class NotFoundError(AppError):
    """Raised when a resource is not found."""
    def __init__(self, resource: str, id: str) -> None:
        self.resource = resource
        self.id = id
        super().__init__(f"{resource} not found: {id}")

class ValidationError(AppError):
    """Raised when input validation fails."""
    def __init__(self, field: str, message: str) -> None:
        self.field = field
        super().__init__(f"Validation error on '{field}': {message}")

# Usage
def get_user(user_id: str) -> dict:
    user = db.find_user(user_id)
    if user is None:
        raise NotFoundError("User", user_id)
    return user

try:
    user = get_user("abc123")
except NotFoundError as e:
    print(f"Could not find {e.resource}: {e.id}")
```

### Exception Groups (3.11+)

```python
# Raise multiple exceptions at once
exceptions = []
for task in tasks:
    try:
        process(task)
    except Exception as e:
        exceptions.append(e)

if exceptions:
    raise ExceptionGroup("batch processing failed", exceptions)

# Handle with except*
try:
    run_batch()
except* ValueError as eg:
    for e in eg.exceptions:
        print(f"Value error: {e}")
except* TypeError as eg:
    for e in eg.exceptions:
        print(f"Type error: {e}")
```

---

## Type Hints

### Basic Type Hints

```python
# Variables
name: str = "Alice"
age: int = 30
scores: list[int] = [95, 87, 92]
config: dict[str, str] = {"host": "localhost"}

# Functions
def process(data: list[str], limit: int = 100) -> dict[str, int]:
    return {item: len(item) for item in data[:limit]}
```

### Advanced Types

```python
from typing import Optional, Union, TypeAlias, TypeVar, Protocol

# Optional (equivalent to X | None in 3.10+)
def find_user(id: str) -> User | None:
    ...

# Union types (3.10+ syntax)
def parse_input(value: str | int | float) -> str:
    return str(value)

# Type aliases
JSON: TypeAlias = dict[str, "JSON"] | list["JSON"] | str | int | float | bool | None
UserId: TypeAlias = str

# TypeVar for generics
T = TypeVar("T")

def first(items: list[T]) -> T | None:
    return items[0] if items else None

# Protocol (structural subtyping / duck typing)
class Closable(Protocol):
    def close(self) -> None:
        ...

def cleanup(resource: Closable) -> None:
    resource.close()

# Literal types
from typing import Literal

def set_mode(mode: Literal["read", "write", "append"]) -> None:
    ...

# TypedDict
from typing import TypedDict

class UserDict(TypedDict):
    name: str
    email: str
    age: int
```

### Running Type Checks

```bash
# Install mypy
pip install mypy

# Run type checks
mypy src/
mypy --strict src/main.py

# Inline type ignore
x = some_untyped_func()  # type: ignore[no-untyped-call]
```

---

## Async / Await

### asyncio Basics

```python
import asyncio

async def fetch_url(url: str) -> str:
    """Simulate fetching a URL."""
    print(f"Fetching {url}...")
    await asyncio.sleep(1)  # simulates I/O
    return f"Response from {url}"

async def main():
    # Run concurrently
    results = await asyncio.gather(
        fetch_url("https://api.example.com/users"),
        fetch_url("https://api.example.com/posts"),
        fetch_url("https://api.example.com/comments"),
    )
    for result in results:
        print(result)

asyncio.run(main())
```

### aiohttp

```python
import aiohttp
import asyncio

async def fetch_json(session: aiohttp.ClientSession, url: str) -> dict:
    async with session.get(url) as response:
        response.raise_for_status()
        return await response.json()

async def main():
    async with aiohttp.ClientSession() as session:
        users, posts = await asyncio.gather(
            fetch_json(session, "https://api.example.com/users"),
            fetch_json(session, "https://api.example.com/posts"),
        )
        print(f"Got {len(users)} users and {len(posts)} posts")

asyncio.run(main())
```

### Async Context Managers and Iterators

```python
class AsyncDBPool:
    """Async database connection pool."""

    async def __aenter__(self):
        self.pool = await create_pool(dsn="postgresql://localhost/mydb")
        return self

    async def __aexit__(self, exc_type, exc_val, exc_tb):
        await self.pool.close()

async def stream_results(query: str):
    """Async generator for streaming large query results."""
    async with AsyncDBPool() as db:
        async with db.pool.acquire() as conn:
            async for row in conn.cursor(query):
                yield row

# Usage
async for row in stream_results("SELECT * FROM events"):
    process(row)
```

### asyncio Task Management

```python
async def worker(name: str, queue: asyncio.Queue):
    while True:
        item = await queue.get()
        try:
            await process(item)
        finally:
            queue.task_done()

async def main():
    queue = asyncio.Queue()

    # Create worker tasks
    workers = [asyncio.create_task(worker(f"w-{i}", queue)) for i in range(5)]

    # Add work items
    for item in work_items:
        await queue.put(item)

    # Wait for all items to be processed
    await queue.join()

    # Cancel workers
    for w in workers:
        w.cancel()
```

---

## Virtual Environments & Packaging

### venv and pip

```bash
# Create virtual environment
python3 -m venv .venv

# Activate
source .venv/bin/activate       # Linux/macOS
# .venv\Scripts\activate        # Windows

# Install packages
pip install requests==2.31.0
pip install -r requirements.txt

# Freeze dependencies
pip freeze > requirements.txt

# Deactivate
deactivate
```

### requirements.txt

```text
# Pin exact versions for reproducibility
requests==2.31.0
flask==3.0.0
sqlalchemy==2.0.23
pytest==7.4.3

# Separate dev dependencies
# requirements-dev.txt
pytest==7.4.3
mypy==1.7.1
ruff==0.1.8
```

### pyproject.toml (Modern Standard)

```toml
[build-system]
requires = ["setuptools>=68.0", "wheel"]
build-backend = "setuptools.backends._legacy:_Backend"

[project]
name = "my-project"
version = "1.0.0"
description = "A useful project"
requires-python = ">=3.10"
dependencies = [
    "requests>=2.31.0,<3.0",
    "pydantic>=2.0,<3.0",
]

[project.optional-dependencies]
dev = [
    "pytest>=7.4",
    "mypy>=1.7",
    "ruff>=0.1",
]

[tool.pytest.ini_options]
testpaths = ["tests"]
addopts = "-v --tb=short"

[tool.mypy]
strict = true
python_version = "3.12"

[tool.ruff]
line-length = 120
target-version = "py312"
```

---

## Testing

### pytest Basics

```python
# tests/test_calculator.py
import pytest
from src.calculator import add, divide

def test_add():
    assert add(2, 3) == 5

def test_add_negative():
    assert add(-1, 1) == 0

def test_divide():
    assert divide(10, 2) == 5.0

def test_divide_by_zero():
    with pytest.raises(ZeroDivisionError):
        divide(10, 0)

# Parametrized tests
@pytest.mark.parametrize("a, b, expected", [
    (1, 1, 2),
    (0, 0, 0),
    (-1, 1, 0),
    (100, 200, 300),
])
def test_add_parametrized(a, b, expected):
    assert add(a, b) == expected
```

### Fixtures

```python
import pytest
from src.database import Database

@pytest.fixture
def db():
    """Create a test database and clean up after."""
    database = Database(":memory:")
    database.create_tables()
    yield database
    database.close()

@pytest.fixture
def sample_user(db):
    """Insert and return a sample user."""
    user = {"name": "Alice", "email": "alice@example.com"}
    db.insert("users", user)
    return user

def test_find_user(db, sample_user):
    found = db.find("users", {"email": "alice@example.com"})
    assert found["name"] == "Alice"

# Shared fixture (conftest.py)
# tests/conftest.py
@pytest.fixture(scope="session")
def app_config():
    return {"host": "localhost", "port": 5432, "db": "test_db"}
```

### Mocking

```python
from unittest.mock import patch, MagicMock, AsyncMock

# Patch external dependency
@patch("src.service.requests.get")
def test_fetch_users(mock_get):
    mock_get.return_value.status_code = 200
    mock_get.return_value.json.return_value = [{"name": "Alice"}]

    result = fetch_users()
    assert len(result) == 1
    mock_get.assert_called_once_with("https://api.example.com/users")

# Mock async functions
@patch("src.service.fetch_data", new_callable=AsyncMock)
async def test_async_service(mock_fetch):
    mock_fetch.return_value = {"status": "ok"}
    result = await process_data()
    assert result["status"] == "ok"
```

### Running Tests

```bash
# Run all tests
pytest

# Run specific file
pytest tests/test_calculator.py

# Run specific test
pytest tests/test_calculator.py::test_add

# Run with verbose output
pytest -v

# Run with coverage
pip install pytest-cov
pytest --cov=src --cov-report=html

# Run only marked tests
pytest -m "not slow"

# Stop on first failure
pytest -x
```

---

## Common Patterns

### Singleton

```python
class DatabaseConnection:
    """Thread-safe singleton database connection."""
    _instance = None
    _lock = threading.Lock()

    def __new__(cls, *args, **kwargs):
        if cls._instance is None:
            with cls._lock:
                if cls._instance is None:
                    cls._instance = super().__new__(cls)
        return cls._instance

    def __init__(self, connection_string: str = ""):
        if not hasattr(self, "_initialized"):
            self._conn = connect(connection_string)
            self._initialized = True
```

### Factory

```python
from dataclasses import dataclass
from enum import Enum

class NotificationType(Enum):
    EMAIL = "email"
    SMS = "sms"
    PUSH = "push"

@dataclass
class Notification:
    recipient: str
    message: str

class NotificationFactory:
    """Create notification senders by type."""
    _registry: dict[NotificationType, type] = {}

    @classmethod
    def register(cls, ntype: NotificationType):
        def decorator(klass):
            cls._registry[ntype] = klass
            return klass
        return decorator

    @classmethod
    def create(cls, ntype: NotificationType, **kwargs):
        if ntype not in cls._registry:
            raise ValueError(f"Unknown notification type: {ntype}")
        return cls._registry[ntype](**kwargs)

@NotificationFactory.register(NotificationType.EMAIL)
class EmailSender:
    def send(self, notification: Notification) -> None:
        ...

@NotificationFactory.register(NotificationType.SMS)
class SMSSender:
    def send(self, notification: Notification) -> None:
        ...
```

### Observer

```python
from collections import defaultdict
from typing import Callable

class EventBus:
    """Simple publish-subscribe event system."""

    def __init__(self):
        self._subscribers: dict[str, list[Callable]] = defaultdict(list)

    def subscribe(self, event: str, callback: Callable) -> None:
        self._subscribers[event].append(callback)

    def unsubscribe(self, event: str, callback: Callable) -> None:
        self._subscribers[event].remove(callback)

    def publish(self, event: str, **data) -> None:
        for callback in self._subscribers[event]:
            callback(**data)

# Usage
bus = EventBus()
bus.subscribe("user.created", lambda name, email: send_welcome_email(name, email))
bus.subscribe("user.created", lambda name, **_: log_new_user(name))
bus.publish("user.created", name="Alice", email="alice@example.com")
```

---

## Common Mistakes & Gotchas

### Mutable Default Arguments

```python
# ❌ WRONG — default list is shared across calls
def add_item(item, items=[]):
    items.append(item)
    return items

add_item("a")  # ["a"]
add_item("b")  # ["a", "b"] — surprise!

# ✅ CORRECT
def add_item(item, items=None):
    if items is None:
        items = []
    items.append(item)
    return items
```

### Late Binding Closures

```python
# ❌ WRONG — all functions reference the same variable
funcs = [lambda: i for i in range(5)]
[f() for f in funcs]  # [4, 4, 4, 4, 4]

# ✅ CORRECT — capture the value
funcs = [lambda i=i: i for i in range(5)]
[f() for f in funcs]  # [0, 1, 2, 3, 4]
```

### Modifying a List While Iterating

```python
# ❌ WRONG — skips elements
items = [1, 2, 3, 4, 5]
for item in items:
    if item % 2 == 0:
        items.remove(item)

# ✅ CORRECT — use list comprehension or iterate over a copy
items = [item for item in items if item % 2 != 0]
```

### The GIL (Global Interpreter Lock)

```python
# The GIL prevents true parallel execution of Python threads.
# Threads are fine for I/O-bound work, not CPU-bound.

# ❌ Threads don't speed up CPU-bound work
import threading
# This won't be faster than single-threaded

# ✅ Use multiprocessing for CPU-bound work
from concurrent.futures import ProcessPoolExecutor

with ProcessPoolExecutor() as pool:
    results = list(pool.map(cpu_intensive_func, data))

# ✅ Use threads/asyncio for I/O-bound work
from concurrent.futures import ThreadPoolExecutor

with ThreadPoolExecutor(max_workers=10) as pool:
    results = list(pool.map(fetch_url, urls))
```

### is vs ==

```python
# 'is' checks identity (same object), '==' checks equality (same value)
a = [1, 2, 3]
b = [1, 2, 3]
a == b   # True (same value)
a is b   # False (different objects)

# Use 'is' only for None, True, False
if result is None:
    ...
```

### Shallow vs Deep Copy

```python
import copy

original = [[1, 2], [3, 4]]

shallow = original.copy()       # or original[:]
shallow[0][0] = 99
print(original[0][0])           # 99 — inner lists are shared!

deep = copy.deepcopy(original)
deep[0][0] = 42
print(original[0][0])           # 99 — not affected
```

---

## Production Tips

### Logging

```python
import logging

# Configure logging — do this once at startup
logging.basicConfig(
    level=logging.INFO,
    format="%(asctime)s %(levelname)s %(name)s %(message)s",
    datefmt="%Y-%m-%d %H:%M:%S",
)

logger = logging.getLogger(__name__)

# Use structured logging in production
logger.info("User logged in", extra={"user_id": user.id, "ip": request.remote_addr})
logger.error("Failed to process payment", exc_info=True)

# ❌ Don't use print() in production
# ❌ Don't use f-strings in logger calls (defeats lazy evaluation)
logger.info("User %s logged in", user.id)  # ✅ lazy formatting
```

### Profiling

```python
# Quick profiling with cProfile
import cProfile
cProfile.run("main()", sort="cumulative")

# Line profiling
# pip install line_profiler
# kernprof -l -v script.py

# Memory profiling
# pip install memory_profiler
# python -m memory_profiler script.py

# Timing specific blocks
import time

start = time.perf_counter()
result = expensive_operation()
elapsed = time.perf_counter() - start
print(f"Took {elapsed:.3f}s")
```

### Packaging for Production

```bash
# Build a wheel
pip install build
python -m build

# Install from wheel
pip install dist/my_project-1.0.0-py3-none-any.whl

# Docker best practices
# Use multi-stage builds, pin base image versions
```

```dockerfile
FROM python:3.12-slim AS builder
WORKDIR /app
COPY requirements.txt .
RUN pip install --no-cache-dir --target=/deps -r requirements.txt

FROM python:3.12-slim
WORKDIR /app
COPY --from=builder /deps /usr/local/lib/python3.12/site-packages
COPY src/ ./src/
USER nobody
CMD ["python", "-m", "src.main"]
```

### Environment Configuration

```python
import os
from dataclasses import dataclass

@dataclass
class Settings:
    """Load settings from environment variables."""
    db_host: str = os.getenv("DB_HOST", "localhost")
    db_port: int = int(os.getenv("DB_PORT", "5432"))
    db_name: str = os.getenv("DB_NAME", "myapp")
    debug: bool = os.getenv("DEBUG", "false").lower() == "true"
    log_level: str = os.getenv("LOG_LEVEL", "INFO")

settings = Settings()
```

---

## Related Topics

- [🧰 CLI / grep](../../cli/grep/) — Search code and logs
- [🧰 CLI / sed](../../cli/sed/) — Stream editing
- [🧰 CLI / jq](../../cli/jq/) — JSON processing
- [🧩 Frameworks / Flask](../../frameworks/) — Python web framework
- [🧩 Frameworks / FastAPI](../../frameworks/) — Async Python API framework
- [🗄 Databases / PostgreSQL](../../databases/) — Database integration
- [⚙️ Engineering / Testing](../../engineering/) — Testing best practices
- [🚀 DevOps / Docker](../../devops/) — Containerizing Python apps
- [📊 Observability](../../observability/) — Application monitoring

---

## References

### Official Documentation

- [Python Docs](https://docs.python.org/3/) — Official documentation
- [Python Tutorial](https://docs.python.org/3/tutorial/) — Official tutorial
- [PEP Index](https://peps.python.org/) — Python Enhancement Proposals
- [PyPI](https://pypi.org/) — Python Package Index

### Books

- *Fluent Python* (2nd ed.) by Luciano Ramalho — Idiomatic Python deep-dive
- *Effective Python* (2nd ed.) by Brett Slatkin — 90 specific ways to write better Python
- *Python Cookbook* (3rd ed.) by David Beazley & Brian K. Jones — Practical recipes
- *Architecture Patterns with Python* by Harry Percival & Bob Gregory — DDD and patterns

### Tools

- [mypy](https://mypy-lang.org/) — Static type checker
- [ruff](https://github.com/astral-sh/ruff) — Fast Python linter and formatter
- [pytest](https://docs.pytest.org/) — Testing framework
- [uv](https://github.com/astral-sh/uv) — Fast Python package manager
- [poetry](https://python-poetry.org/) — Dependency management

---

*Part of [Tech-Byte Languages](../). Last updated: 2026-08.*
