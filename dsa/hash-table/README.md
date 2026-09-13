# Hash Tables

> **O(1) lookup, insert, and delete — the workhorse of efficient programming.** Hash tables map keys to values using a hash function, enabling constant-time average performance for core operations.

---

## 📑 Table of Contents

- [What is a Hash Table?](#-what-is-a-hash-table)
- [How Hash Maps Work Internally](#-how-hash-maps-work-internally)
- [Operations and Complexity](#-operations-and-complexity)
- [Python dict and Go map](#-python-dict-and-go-map)
- [Common Patterns](#-common-patterns)
- [Implementation from Scratch](#-implementation-from-scratch)
- [Common Interview Problems](#-common-interview-problems)
- [When to Use](#-when-to-use)
- [Related Topics](#-related-topics)

---

## 🧠 What is a Hash Table?

A hash table stores key-value pairs. It uses a **hash function** to compute an index (bucket) where the value is stored.

```text
Key "alice" → hash("alice") = 748293 → 748293 % 8 = 5 → Bucket 5

Bucket:  [0]  [1]  [2]  [3]  [4]  [5]       [6]  [7]
                                     "alice"→42
```

**Components:**

1. **Hash function:** Converts key to an integer
2. **Buckets:** Array of storage locations
3. **Collision resolution:** Handles multiple keys mapping to the same bucket

### Hash Function Properties

A good hash function:
- **Deterministic:** Same key always produces the same hash
- **Uniform distribution:** Keys spread evenly across buckets
- **Fast to compute:** O(1) for fixed-size keys
- **Avalanche effect:** Small change in key → large change in hash

---

## ⚙️ How Hash Maps Work Internally

### Chaining (Separate Chaining)

Each bucket holds a linked list (or other collection) of entries.

```text
Bucket 0: → [("apple", 1)] → [("grape", 7)]
Bucket 1: → [("banana", 2)]
Bucket 2: → (empty)
Bucket 3: → [("cherry", 3)] → [("date", 4)] → [("fig", 9)]
Bucket 4: → (empty)
```

**Pros:** Simple, handles high load gracefully.
**Cons:** Extra memory for linked list pointers, poor cache locality.

### Open Addressing

All entries are stored in the bucket array itself. On collision, probe for the next empty slot.

**Linear probing:**

```text
hash("alice") = 3, but bucket 3 is occupied
Try bucket 4 → occupied
Try bucket 5 → empty → store here

Probe sequence: 3, 4, 5, 6, 7, ... (wraps around)
```

**Quadratic probing:** Try positions i, i+1², i+2², i+3², ...

**Double hashing:** Use a second hash function: `(h1(key) + i × h2(key)) % size`

| Method | Clustering | Cache Performance | Deletion |
|--------|-----------|------------------|----------|
| Linear probing | Primary clustering | Best (sequential) | Needs tombstones |
| Quadratic probing | Secondary clustering | Good | Needs tombstones |
| Double hashing | Minimal clustering | Fair | Needs tombstones |

### Load Factor and Resizing

```text
Load factor = number_of_entries / number_of_buckets

Too high (> 0.75): Many collisions, performance degrades
Too low (< 0.25):  Wasted memory

When load factor exceeds threshold:
1. Create new array (typically 2x size)
2. Rehash all entries into new array
3. Cost: O(n) — but amortized O(1) per insert
```

---

## 📊 Operations and Complexity

| Operation | Average Case | Worst Case | Notes |
|-----------|-------------|-----------|-------|
| Insert | O(1) | O(n) | Worst case: all keys hash to same bucket |
| Lookup | O(1) | O(n) | |
| Delete | O(1) | O(n) | |
| Resize | O(n) | O(n) | Amortized O(1) per insert |

**Worst case** happens when hash function is poor or adversarial input causes all keys to collide. In practice with good hash functions, operations are O(1).

---

## 🐍 Python dict and Go map

### Python dict

```python
# Python dict uses open addressing with random probing
# Guaranteed insertion order since Python 3.7+

# Create
scores = {"alice": 95, "bob": 87, "charlie": 92}

# Access: O(1) average
print(scores["alice"])     # 95
print(scores.get("dave"))  # None (no KeyError)
print(scores.get("dave", 0))  # 0 (default value)

# Insert/Update: O(1) average
scores["dave"] = 88

# Delete: O(1) average
del scores["bob"]
popped = scores.pop("charlie", None)  # Remove and return

# Iteration
for key, value in scores.items():
    print(f"{key}: {value}")

# Membership test: O(1) average
if "alice" in scores:
    print("Found!")

# Dictionary comprehension
squared = {x: x**2 for x in range(10)}

# defaultdict — auto-initialize missing keys
from collections import defaultdict
freq = defaultdict(int)
for char in "hello":
    freq[char] += 1  # No KeyError for missing keys

# Counter — specialized frequency counter
from collections import Counter
freq = Counter("hello")  # Counter({'l': 2, 'h': 1, 'e': 1, 'o': 1})
```

### Go map

```go
// Go map uses hash table with chaining (buckets)
// Iteration order is randomized by design

// Create
scores := map[string]int{
    "alice":   95,
    "bob":     87,
    "charlie": 92,
}

// Access: O(1) average
score := scores["alice"]     // 95
score, ok := scores["dave"]  // 0, false (zero value + exists check)

// Insert/Update: O(1) average
scores["dave"] = 88

// Delete: O(1) average
delete(scores, "bob")

// Iteration (random order)
for key, value := range scores {
    fmt.Printf("%s: %d\n", key, value)
}

// Membership test: O(1) average
if _, ok := scores["alice"]; ok {
    fmt.Println("Found!")
}

// map is not safe for concurrent use — use sync.Map or mutex
var mu sync.Mutex
mu.Lock()
scores["alice"] = 100
mu.Unlock()
```

---

## 🔄 Common Patterns

### Frequency Counting

```python
def char_frequency(s: str) -> dict[str, int]:
    """Count character frequencies. O(n)."""
    freq = {}
    for char in s:
        freq[char] = freq.get(char, 0) + 1
    return freq

# Or with Counter:
from collections import Counter
freq = Counter("abracadabra")
# Counter({'a': 5, 'b': 2, 'r': 2, 'c': 1, 'd': 1})
```

```go
func charFrequency(s string) map[rune]int {
    freq := make(map[rune]int)
    for _, ch := range s {
        freq[ch]++
    }
    return freq
}
```

### Grouping / Bucketing

```python
def group_by_length(words: list[str]) -> dict[int, list[str]]:
    """Group words by their length. O(n)."""
    groups = defaultdict(list)
    for word in words:
        groups[len(word)].append(word)
    return dict(groups)
```

### Two-Sum Pattern

```python
def two_sum(nums: list[int], target: int) -> list[int]:
    """Find indices of two numbers that sum to target. O(n)."""
    seen = {}  # value → index
    for i, num in enumerate(nums):
        complement = target - num
        if complement in seen:
            return [seen[complement], i]
        seen[num] = i
    return []
```

### Caching / Memoization

```python
def fibonacci(n: int, memo: dict[int, int] = None) -> int:
    """Fibonacci with hash table memoization. O(n)."""
    if memo is None:
        memo = {}
    if n in memo:
        return memo[n]
    if n <= 1:
        return n
    memo[n] = fibonacci(n - 1, memo) + fibonacci(n - 2, memo)
    return memo[n]
```

---

## 🔧 Implementation from Scratch

```python
class HashTable:
    """Simple hash table with chaining."""

    def __init__(self, capacity: int = 16):
        self.capacity = capacity
        self.size = 0
        self.buckets: list[list[tuple]] = [[] for _ in range(capacity)]
        self.load_factor_threshold = 0.75

    def _hash(self, key: str) -> int:
        """Simple hash function."""
        hash_value = 0
        for char in key:
            hash_value = (hash_value * 31 + ord(char)) % self.capacity
        return hash_value

    def put(self, key: str, value) -> None:
        """Insert or update a key-value pair. O(1) amortized."""
        if self.size / self.capacity >= self.load_factor_threshold:
            self._resize()

        index = self._hash(key)
        bucket = self.buckets[index]

        # Update existing key
        for i, (k, v) in enumerate(bucket):
            if k == key:
                bucket[i] = (key, value)
                return

        # Insert new key
        bucket.append((key, value))
        self.size += 1

    def get(self, key: str, default=None):
        """Retrieve value by key. O(1) average."""
        index = self._hash(key)
        for k, v in self.buckets[index]:
            if k == key:
                return v
        return default

    def delete(self, key: str) -> bool:
        """Remove a key-value pair. O(1) average."""
        index = self._hash(key)
        bucket = self.buckets[index]
        for i, (k, v) in enumerate(bucket):
            if k == key:
                bucket.pop(i)
                self.size -= 1
                return True
        return False

    def _resize(self) -> None:
        """Double capacity and rehash all entries. O(n)."""
        old_buckets = self.buckets
        self.capacity *= 2
        self.buckets = [[] for _ in range(self.capacity)]
        self.size = 0
        for bucket in old_buckets:
            for key, value in bucket:
                self.put(key, value)

    def __contains__(self, key: str) -> bool:
        return self.get(key) is not None

    def __len__(self) -> int:
        return self.size
```

---

## 🎯 Common Interview Problems

### Group Anagrams

```python
def group_anagrams(strs: list[str]) -> list[list[str]]:
    """Group strings that are anagrams. O(n × k log k) where k = max string length."""
    groups = defaultdict(list)
    for s in strs:
        key = tuple(sorted(s))  # Anagrams have same sorted form
        groups[key].append(s)
    return list(groups.values())

# O(n × k) version using character counts as key:
def group_anagrams_optimal(strs: list[str]) -> list[list[str]]:
    groups = defaultdict(list)
    for s in strs:
        count = [0] * 26
        for c in s:
            count[ord(c) - ord('a')] += 1
        groups[tuple(count)].append(s)
    return list(groups.values())
```

### Longest Substring Without Repeating Characters

```python
def length_of_longest_substring(s: str) -> int:
    """Sliding window + hash set. O(n)."""
    char_index = {}
    max_len = 0
    left = 0
    for right, char in enumerate(s):
        if char in char_index and char_index[char] >= left:
            left = char_index[char] + 1
        char_index[char] = right
        max_len = max(max_len, right - left + 1)
    return max_len
```

### LRU Cache

```python
from collections import OrderedDict

class LRUCache:
    """Least Recently Used cache using OrderedDict. O(1) for get/put."""

    def __init__(self, capacity: int):
        self.capacity = capacity
        self.cache = OrderedDict()

    def get(self, key: int) -> int:
        if key not in self.cache:
            return -1
        self.cache.move_to_end(key)  # Mark as recently used
        return self.cache[key]

    def put(self, key: int, value: int) -> None:
        if key in self.cache:
            self.cache.move_to_end(key)
        self.cache[key] = value
        if len(self.cache) > self.capacity:
            self.cache.popitem(last=False)  # Remove least recently used
```

---

## 🤔 When to Use

### Hash Table vs Other Data Structures

| Use Hash Table When | Use Alternative When |
|--------------------|---------------------|
| Need O(1) lookup by key | Need sorted order → balanced BST |
| Counting frequencies | Need prefix matching → Trie |
| Grouping/bucketing | Need range queries → sorted array + binary search |
| Caching results | Need min/max efficiently → Heap |
| Detecting duplicates | Memory is extremely constrained → bit array/bloom filter |
| Two-sum style problems | Need to iterate in order → sorted structure |

### Hash Table vs Sorted Array

| Operation | Hash Table | Sorted Array |
|-----------|-----------|-------------|
| Search | O(1) avg | O(log n) |
| Insert | O(1) avg | O(n) |
| Delete | O(1) avg | O(n) |
| Min/Max | O(n) | O(1) |
| Range query | O(n) | O(log n + k) |
| Ordered iteration | O(n log n) | O(n) |
| Space | O(n) + overhead | O(n) |

---

## 🔗 Related Topics

- [Complexity Analysis](../complexity.md) — Hash table complexity analysis
- [Arrays](../arrays/) — Arrays as underlying storage for hash tables
- [Trees](../tree/) — BSTs as ordered alternative to hash tables
- [Dynamic Programming](../dynamic-programming/) — Memoization uses hash tables
- [Graphs](../graph/) — Adjacency lists often use hash maps
- [System Design: Caching](../../system-design/caching.md) — Hash tables power caching systems
- [Distributed Systems](../../distributed-systems/) — Consistent hashing for distributed systems
- [Languages: Python](../../languages/python/) — Python dict internals
- [Languages: Go](../../languages/go/) — Go map internals

---

> **Hash tables are everywhere.** Database indexes, DNS caches, compiler symbol tables, HTTP headers, JSON objects, API rate limiters — they all use hash-based data structures. Understanding how they work internally makes you a better engineer, not just a better coder.
