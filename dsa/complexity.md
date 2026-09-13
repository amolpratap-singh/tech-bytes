# Complexity Analysis

> **Measure how algorithms scale.** Complexity analysis tells you how an algorithm's resource consumption (time, space) grows as input size increases. It's the universal language for comparing algorithmic efficiency.

---

## 📑 Table of Contents

- [Big-O Notation](#-big-o-notation)
- [Space vs Time Complexity](#-space-vs-time-complexity)
- [Amortized Analysis](#-amortized-analysis)
- [Complexity Comparison](#-complexity-comparison)
- [How to Analyze Complexity](#-how-to-analyze-complexity)
- [Common Complexities by Data Structure](#-common-complexities-by-data-structure)
- [Common Complexities by Algorithm](#-common-complexities-by-algorithm)
- [Related Topics](#-related-topics)

---

## 📐 Big-O Notation

Big-O describes the **upper bound** of an algorithm's growth rate. It answers: "In the worst case, how does the runtime scale with input size?"

### O(1) — Constant Time

Runtime doesn't change with input size.

```python
def get_first(arr: list) -> int:
    return arr[0]  # Always one operation, regardless of array size

# Hash table lookup
my_dict = {"key": "value"}
result = my_dict["key"]  # O(1) average
```

**Examples:** Array access by index, hash table lookup, push/pop on a stack.

### O(log n) — Logarithmic Time

Runtime grows logarithmically. Halving the problem each step.

```python
def binary_search(arr: list[int], target: int) -> int:
    """O(log n) — halves search space each iteration."""
    left, right = 0, len(arr) - 1
    while left <= right:
        mid = (left + right) // 2
        if arr[mid] == target:
            return mid
        elif arr[mid] < target:
            left = mid + 1
        else:
            right = mid - 1
    return -1
```

**Examples:** Binary search, balanced BST operations, binary exponentiation.

### O(n) — Linear Time

Runtime grows proportionally with input size.

```python
def find_max(arr: list[int]) -> int:
    """O(n) — must check every element."""
    max_val = arr[0]
    for num in arr:
        max_val = max(max_val, num)
    return max_val
```

**Examples:** Linear search, traversing an array/list, counting elements.

### O(n log n) — Linearithmic Time

Common for efficient sorting algorithms.

```python
def merge_sort(arr: list[int]) -> list[int]:
    """O(n log n) — divide and conquer."""
    if len(arr) <= 1:
        return arr
    mid = len(arr) // 2
    left = merge_sort(arr[:mid])
    right = merge_sort(arr[mid:])
    return merge(left, right)
```

**Examples:** Merge sort, heap sort, quicksort (average), many divide-and-conquer algorithms.

### O(n²) — Quadratic Time

Often involves nested loops over the input.

```python
def bubble_sort(arr: list[int]) -> list[int]:
    """O(n²) — nested loops."""
    n = len(arr)
    for i in range(n):
        for j in range(0, n - i - 1):
            if arr[j] > arr[j + 1]:
                arr[j], arr[j + 1] = arr[j + 1], arr[j]
    return arr
```

**Examples:** Bubble sort, insertion sort, brute-force pairwise comparison.

### O(2ⁿ) — Exponential Time

Doubles with each additional input element. Only feasible for small n.

```python
def fibonacci_naive(n: int) -> int:
    """O(2^n) — exponential without memoization."""
    if n <= 1:
        return n
    return fibonacci_naive(n - 1) + fibonacci_naive(n - 2)
```

**Examples:** Naive Fibonacci, generating all subsets, brute-force combinatorial problems.

### O(n!) — Factorial Time

Grows extremely fast. Only for very small inputs.

```python
def permutations(arr: list) -> list[list]:
    """O(n!) — generate all permutations."""
    if len(arr) <= 1:
        return [arr]
    result = []
    for i, item in enumerate(arr):
        rest = arr[:i] + arr[i + 1:]
        for perm in permutations(rest):
            result.append([item] + perm)
    return result
```

**Examples:** Generating all permutations, traveling salesman (brute force).

---

## 🧠 Space vs Time Complexity

**Time complexity:** How many operations the algorithm performs.
**Space complexity:** How much extra memory the algorithm uses.

| Algorithm | Time | Space | Trade-off |
|-----------|------|-------|-----------|
| Merge sort | O(n log n) | O(n) | Fast but uses extra memory |
| Quicksort (in-place) | O(n log n)* | O(log n) | Fast, memory-efficient |
| Hash table lookup | O(1) | O(n) | Fast lookup, uses memory for table |
| Binary search | O(log n) | O(1) | Fast, minimal memory |
| BFS traversal | O(V + E) | O(V) | Queue stores visited nodes |
| DFS traversal | O(V + E) | O(V) | Recursion stack or explicit stack |
| DP (memoization) | O(subproblems × work per) | O(subproblems) | Time for space |

**Common trade-off:** Use more space to gain speed (caching, hash tables, memoization).

### Auxiliary vs Total Space

```python
def sort_array(arr: list[int]) -> list[int]:
    # Total space: O(n) (the input itself)
    # Auxiliary space: O(1) (in-place sort)
    arr.sort()  # Timsort uses O(n) auxiliary space actually
    return arr

def merge_sort(arr: list[int]) -> list[int]:
    # Total space: O(n)
    # Auxiliary space: O(n) (creates new arrays at each level)
    ...
```

**In practice:** When discussing space complexity, we usually mean **auxiliary space** — extra memory beyond the input.

---

## 📈 Amortized Analysis

Some operations are occasionally expensive but cheap on average. Amortized analysis considers the **average cost over a sequence of operations**.

### Example: Dynamic Array (Python list)

```text
Appending to a dynamic array:
  - Most appends: O(1) — just write to pre-allocated space
  - Occasional append: O(n) — when array is full, allocate 2x and copy everything

Amortized cost:
  n appends → total cost = n + (1 + 2 + 4 + 8 + ... + n) = n + 2n = 3n
  Per operation: 3n / n = O(1) amortized
```

### Example: Hash Table

```text
Insert with rehashing:
  - Most inserts: O(1) — direct bucket placement
  - Occasional insert: O(n) — when load factor exceeds threshold, rehash all entries

Amortized: O(1) per insert
```

### Key Insight

Amortized ≠ average case. Amortized guarantees the average over **any** sequence of operations, while average case is over random inputs.

---

## 📊 Complexity Comparison

### Visual Growth Rates

```text
Operations vs Input Size (n):

n=10      n=100     n=1,000   n=10,000  n=100,000
─────     ─────     ─────     ─────     ─────
O(1)         1         1         1         1           1
O(log n)     3         7        10        13          17
O(n)        10       100     1,000    10,000     100,000
O(n log n)  33       664    10,000   132,877   1,660,964
O(n²)      100    10,000 1,000,000   10⁸         10¹⁰ ☠️
O(2ⁿ)    1,024    10³⁰☠️     ☠️        ☠️           ☠️
```

### Practical Limits

| Time Limit | Max n for O(n²) | Max n for O(n log n) | Max n for O(n) |
|-----------|-----------------|---------------------|----------------|
| 1 second | ~10,000 | ~10,000,000 | ~100,000,000 |
| 10 seconds | ~30,000 | ~100,000,000 | ~1,000,000,000 |

---

## 🔍 How to Analyze Complexity

### Simple Loops

```python
# O(n)
for i in range(n):
    do_something()  # O(1)

# O(n²)
for i in range(n):
    for j in range(n):
        do_something()  # O(1)

# O(n × m) — two different inputs
for i in range(n):
    for j in range(m):
        do_something()  # O(1)
```

### Logarithmic Patterns

```python
# O(log n) — halving
i = n
while i > 0:
    do_something()
    i //= 2

# O(n log n) — linear loop with log inner
for i in range(n):      # O(n)
    j = n
    while j > 0:        # O(log n)
        do_something()
        j //= 2
```

### Recursion

```python
# O(2^n) — two recursive calls, depth n
def fib(n):
    if n <= 1:
        return n
    return fib(n - 1) + fib(n - 2)

# O(n) — one recursive call, depth n
def factorial(n):
    if n <= 1:
        return 1
    return n * factorial(n - 1)

# O(n log n) — divide in half, linear work each level
def merge_sort(arr):
    if len(arr) <= 1:
        return arr
    mid = len(arr) // 2
    left = merge_sort(arr[:mid])      # T(n/2)
    right = merge_sort(arr[mid:])     # T(n/2)
    return merge(left, right)         # O(n) work
    # T(n) = 2T(n/2) + O(n) → O(n log n) by Master theorem
```

### Multiple Inputs

```python
# Don't simplify to O(n²) — it's O(n × m)
def compare_lists(list_a: list, list_b: list) -> list:
    result = []
    for a in list_a:       # O(n) where n = len(list_a)
        for b in list_b:   # O(m) where m = len(list_b)
            if a == b:
                result.append(a)
    return result
```

### Drop Constants and Lower Terms

```text
O(2n + 100)    → O(n)
O(n² + n)      → O(n²)
O(n + log n)   → O(n)
O(5n log n)    → O(n log n)
```

---

## 📋 Common Complexities by Data Structure

| Data Structure | Access | Search | Insert | Delete | Notes |
|---------------|--------|--------|--------|--------|-------|
| **Array** | O(1) | O(n) | O(n) | O(n) | Insert/delete shifts elements |
| **Dynamic Array** | O(1) | O(n) | O(1)* | O(n) | *Amortized append |
| **Singly Linked List** | O(n) | O(n) | O(1)† | O(1)† | †At head/given node |
| **Doubly Linked List** | O(n) | O(n) | O(1)† | O(1)† | †At head/tail/given node |
| **Hash Table** | — | O(1)* | O(1)* | O(1)* | *Average case; O(n) worst |
| **BST (balanced)** | — | O(log n) | O(log n) | O(log n) | AVL, Red-Black |
| **BST (unbalanced)** | — | O(n) | O(n) | O(n) | Degenerates to linked list |
| **Min/Max Heap** | O(1)‡ | O(n) | O(log n) | O(log n) | ‡Peek min/max only |
| **Trie** | — | O(k) | O(k) | O(k) | k = key length |
| **Skip List** | — | O(log n) | O(log n) | O(log n) | Probabilistic |

---

## 📋 Common Complexities by Algorithm

| Algorithm | Best | Average | Worst | Space | Stable |
|-----------|------|---------|-------|-------|--------|
| **Bubble Sort** | O(n) | O(n²) | O(n²) | O(1) | Yes |
| **Selection Sort** | O(n²) | O(n²) | O(n²) | O(1) | No |
| **Insertion Sort** | O(n) | O(n²) | O(n²) | O(1) | Yes |
| **Merge Sort** | O(n log n) | O(n log n) | O(n log n) | O(n) | Yes |
| **Quick Sort** | O(n log n) | O(n log n) | O(n²) | O(log n) | No |
| **Heap Sort** | O(n log n) | O(n log n) | O(n log n) | O(1) | No |
| **Counting Sort** | O(n + k) | O(n + k) | O(n + k) | O(k) | Yes |
| **Radix Sort** | O(nk) | O(nk) | O(nk) | O(n + k) | Yes |
| **Binary Search** | O(1) | O(log n) | O(log n) | O(1) | — |
| **BFS** | O(V + E) | O(V + E) | O(V + E) | O(V) | — |
| **DFS** | O(V + E) | O(V + E) | O(V + E) | O(V) | — |
| **Dijkstra** | O((V + E) log V) | O((V + E) log V) | O((V + E) log V) | O(V) | — |

---

## 🔗 Related Topics

- [Arrays](arrays/) — Array algorithms and patterns
- [Hash Tables](hash-table/) — Hash-based data structures
- [Trees](tree/) — Tree traversals and operations
- [Graphs](graph/) — Graph algorithms
- [Dynamic Programming](dynamic-programming/) — Optimization with memoization
- [System Design](../system-design/) — Complexity matters at scale
- [Engineering](../engineering/) — Performance optimization

---

> **Complexity analysis is a tool, not a religion.** O(n²) with a small constant factor can be faster than O(n log n) for small inputs. Always profile real-world performance. But understanding complexity gives you the vocabulary to reason about scalability before you hit the wall.
