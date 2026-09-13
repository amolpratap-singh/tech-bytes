# Arrays

> **The most fundamental data structure.** Arrays store elements in contiguous memory locations, providing O(1) random access — the building block for nearly every other data structure and algorithm.

---

## 📑 Table of Contents

- [What are Arrays?](#-what-are-arrays)
- [Operations and Complexity](#-operations-and-complexity)
- [Static vs Dynamic Arrays](#-static-vs-dynamic-arrays)
- [Two Pointer Technique](#-two-pointer-technique)
- [Sliding Window](#-sliding-window)
- [Prefix Sum](#-prefix-sum)
- [Kadane's Algorithm](#-kadanes-algorithm)
- [Common Patterns](#-common-patterns)
- [Common Interview Problems](#-common-interview-problems)
- [Common Mistakes](#-common-mistakes)
- [Related Topics](#-related-topics)

---

## 🧠 What are Arrays?

An array is a collection of elements stored at **contiguous memory locations**. Each element is accessible by its index in O(1) time.

```text
Memory layout:
Index:    [0]   [1]   [2]   [3]   [4]
Value:  | 10  | 20  | 30  | 40  | 50  |
Address: 0x00  0x04  0x08  0x0C  0x10

Address of element i = base_address + (i × element_size)
```

**Key properties:**
- Fixed size (static) or dynamically resizable
- Contiguous memory → excellent cache locality
- O(1) random access by index
- O(n) insert/delete (must shift elements)

---

## ⚙️ Operations and Complexity

| Operation | Time Complexity | Notes |
|-----------|----------------|-------|
| Access by index | O(1) | Direct address calculation |
| Search (unsorted) | O(n) | Must scan all elements |
| Search (sorted) | O(log n) | Binary search |
| Insert at end | O(1)* | *Amortized for dynamic arrays |
| Insert at index | O(n) | Shift elements right |
| Delete at end | O(1) | |
| Delete at index | O(n) | Shift elements left |
| Append | O(1)* | *Amortized for dynamic arrays |

---

## 📦 Static vs Dynamic Arrays

| Feature | Static Array | Dynamic Array |
|---------|-------------|---------------|
| Size | Fixed at creation | Grows automatically |
| Memory | Exact allocation | Over-allocates (typically 2x) |
| Append | Not supported (fixed) | O(1) amortized |
| Languages | C arrays, Go arrays | Python list, Java ArrayList, Go slices |

### Python Lists (Dynamic Array)

```python
# Python list is a dynamic array
arr = [1, 2, 3, 4, 5]

# Access: O(1)
print(arr[2])  # 3

# Append: O(1) amortized
arr.append(6)

# Insert at index: O(n) — shifts elements
arr.insert(0, 0)  # [0, 1, 2, 3, 4, 5, 6]

# Delete by index: O(n) — shifts elements
del arr[0]  # [1, 2, 3, 4, 5, 6]

# Slicing: O(k) where k = slice size
sub = arr[1:4]  # [2, 3, 4]
```

### Go Slices (Dynamic Array)

```go
// Go slice is a dynamic array backed by an array
arr := []int{1, 2, 3, 4, 5}

// Access: O(1)
fmt.Println(arr[2]) // 3

// Append: O(1) amortized
arr = append(arr, 6)

// Length and capacity
fmt.Println(len(arr), cap(arr))

// Slicing: O(1) — creates a view, not a copy
sub := arr[1:4] // [2, 3, 4]
```

---

## 👆 Two Pointer Technique

Use two pointers to traverse the array, reducing O(n²) brute force to O(n).

### Pattern 1: Opposite Ends (Converging Pointers)

Used on **sorted arrays** or when checking properties from both ends.

```python
def two_sum_sorted(arr: list[int], target: int) -> tuple[int, int]:
    """Find two numbers that sum to target in a sorted array. O(n)."""
    left, right = 0, len(arr) - 1
    while left < right:
        current_sum = arr[left] + arr[right]
        if current_sum == target:
            return (left, right)
        elif current_sum < target:
            left += 1
        else:
            right -= 1
    return (-1, -1)
```

```go
func twoSumSorted(arr []int, target int) (int, int) {
    left, right := 0, len(arr)-1
    for left < right {
        sum := arr[left] + arr[right]
        if sum == target {
            return left, right
        } else if sum < target {
            left++
        } else {
            right--
        }
    }
    return -1, -1
}
```

### Pattern 2: Same Direction (Fast/Slow)

Used for removing duplicates, partitioning, or cycle detection.

```python
def remove_duplicates(arr: list[int]) -> int:
    """Remove duplicates from sorted array in-place. O(n) time, O(1) space."""
    if not arr:
        return 0
    slow = 0
    for fast in range(1, len(arr)):
        if arr[fast] != arr[slow]:
            slow += 1
            arr[slow] = arr[fast]
    return slow + 1  # Length of unique portion
```

```go
func removeDuplicates(arr []int) int {
    if len(arr) == 0 {
        return 0
    }
    slow := 0
    for fast := 1; fast < len(arr); fast++ {
        if arr[fast] != arr[slow] {
            slow++
            arr[slow] = arr[fast]
        }
    }
    return slow + 1
}
```

---

## 🪟 Sliding Window

Maintain a window over a contiguous subarray, sliding it across to find optimal results.

### Fixed-Size Window

```python
def max_sum_subarray(arr: list[int], k: int) -> int:
    """Find maximum sum of any contiguous subarray of size k. O(n)."""
    if len(arr) < k:
        return 0

    # Calculate sum of first window
    window_sum = sum(arr[:k])
    max_sum = window_sum

    # Slide the window
    for i in range(k, len(arr)):
        window_sum += arr[i] - arr[i - k]  # Add right, remove left
        max_sum = max(max_sum, window_sum)

    return max_sum
```

```go
func maxSumSubarray(arr []int, k int) int {
    if len(arr) < k {
        return 0
    }
    windowSum := 0
    for i := 0; i < k; i++ {
        windowSum += arr[i]
    }
    maxSum := windowSum
    for i := k; i < len(arr); i++ {
        windowSum += arr[i] - arr[i-k]
        if windowSum > maxSum {
            maxSum = windowSum
        }
    }
    return maxSum
}
```

### Variable-Size Window

```python
def min_subarray_len(target: int, arr: list[int]) -> int:
    """Find minimum length subarray with sum >= target. O(n)."""
    left = 0
    current_sum = 0
    min_len = float('inf')

    for right in range(len(arr)):
        current_sum += arr[right]

        while current_sum >= target:
            min_len = min(min_len, right - left + 1)
            current_sum -= arr[left]
            left += 1

    return min_len if min_len != float('inf') else 0
```

```go
func minSubarrayLen(target int, arr []int) int {
    left, currentSum, minLen := 0, 0, len(arr)+1
    for right := 0; right < len(arr); right++ {
        currentSum += arr[right]
        for currentSum >= target {
            if right-left+1 < minLen {
                minLen = right - left + 1
            }
            currentSum -= arr[left]
            left++
        }
    }
    if minLen == len(arr)+1 {
        return 0
    }
    return minLen
}
```

---

## ➕ Prefix Sum

Pre-compute cumulative sums to answer range sum queries in O(1).

```python
def build_prefix_sum(arr: list[int]) -> list[int]:
    """Build prefix sum array. O(n)."""
    prefix = [0] * (len(arr) + 1)
    for i in range(len(arr)):
        prefix[i + 1] = prefix[i] + arr[i]
    return prefix

def range_sum(prefix: list[int], left: int, right: int) -> int:
    """Sum of arr[left:right+1] in O(1)."""
    return prefix[right + 1] - prefix[left]

# Example
arr = [1, 2, 3, 4, 5]
prefix = build_prefix_sum(arr)  # [0, 1, 3, 6, 10, 15]
print(range_sum(prefix, 1, 3))   # sum(2,3,4) = 9
```

```go
func buildPrefixSum(arr []int) []int {
    prefix := make([]int, len(arr)+1)
    for i := 0; i < len(arr); i++ {
        prefix[i+1] = prefix[i] + arr[i]
    }
    return prefix
}

func rangeSum(prefix []int, left, right int) int {
    return prefix[right+1] - prefix[left]
}
```

**Use cases:** Range sum queries, subarray sum equals K, product except self.

---

## 📈 Kadane's Algorithm

Find the maximum sum contiguous subarray in O(n).

```python
def max_subarray_sum(arr: list[int]) -> int:
    """Kadane's algorithm — maximum subarray sum. O(n) time, O(1) space."""
    max_ending_here = arr[0]
    max_so_far = arr[0]

    for i in range(1, len(arr)):
        # Either extend the current subarray or start a new one
        max_ending_here = max(arr[i], max_ending_here + arr[i])
        max_so_far = max(max_so_far, max_ending_here)

    return max_so_far

# Example
arr = [-2, 1, -3, 4, -1, 2, 1, -5, 4]
print(max_subarray_sum(arr))  # 6 (subarray: [4, -1, 2, 1])
```

```go
func maxSubarraySum(arr []int) int {
    maxEndingHere := arr[0]
    maxSoFar := arr[0]
    for i := 1; i < len(arr); i++ {
        if arr[i] > maxEndingHere+arr[i] {
            maxEndingHere = arr[i]
        } else {
            maxEndingHere += arr[i]
        }
        if maxEndingHere > maxSoFar {
            maxSoFar = maxEndingHere
        }
    }
    return maxSoFar
}
```

**Key insight:** At each position, decide whether to extend the current subarray or start fresh. If the current sum is negative, starting fresh is always better.

---

## 🔄 Common Patterns

### Reverse an Array

```python
def reverse_array(arr: list[int]) -> list[int]:
    """Reverse in-place. O(n) time, O(1) space."""
    left, right = 0, len(arr) - 1
    while left < right:
        arr[left], arr[right] = arr[right], arr[left]
        left += 1
        right -= 1
    return arr
```

### Rotate an Array

```python
def rotate_array(arr: list[int], k: int) -> list[int]:
    """Rotate array right by k positions. O(n) time, O(1) space."""
    n = len(arr)
    k = k % n  # Handle k > n

    def reverse(start: int, end: int) -> None:
        while start < end:
            arr[start], arr[end] = arr[end], arr[start]
            start += 1
            end -= 1

    # Reverse entire array, then reverse each part
    reverse(0, n - 1)
    reverse(0, k - 1)
    reverse(k, n - 1)
    return arr

# [1,2,3,4,5] rotate by 2 → [4,5,1,2,3]
```

### Merge Two Sorted Arrays

```python
def merge_sorted(arr1: list[int], arr2: list[int]) -> list[int]:
    """Merge two sorted arrays. O(n + m) time, O(n + m) space."""
    result = []
    i, j = 0, 0
    while i < len(arr1) and j < len(arr2):
        if arr1[i] <= arr2[j]:
            result.append(arr1[i])
            i += 1
        else:
            result.append(arr2[j])
            j += 1
    result.extend(arr1[i:])
    result.extend(arr2[j:])
    return result
```

```go
func mergeSorted(arr1, arr2 []int) []int {
    result := make([]int, 0, len(arr1)+len(arr2))
    i, j := 0, 0
    for i < len(arr1) && j < len(arr2) {
        if arr1[i] <= arr2[j] {
            result = append(result, arr1[i])
            i++
        } else {
            result = append(result, arr2[j])
            j++
        }
    }
    result = append(result, arr1[i:]...)
    result = append(result, arr2[j:]...)
    return result
}
```

---

## 🎯 Common Interview Problems

### Two Sum

```python
def two_sum(nums: list[int], target: int) -> list[int]:
    """Find indices of two numbers that sum to target. O(n)."""
    seen = {}
    for i, num in enumerate(nums):
        complement = target - num
        if complement in seen:
            return [seen[complement], i]
        seen[num] = i
    return []
```

### Three Sum

```python
def three_sum(nums: list[int]) -> list[list[int]]:
    """Find all unique triplets that sum to zero. O(n²)."""
    nums.sort()
    result = []
    for i in range(len(nums) - 2):
        if i > 0 and nums[i] == nums[i - 1]:
            continue  # Skip duplicates
        left, right = i + 1, len(nums) - 1
        while left < right:
            total = nums[i] + nums[left] + nums[right]
            if total == 0:
                result.append([nums[i], nums[left], nums[right]])
                while left < right and nums[left] == nums[left + 1]:
                    left += 1
                while left < right and nums[right] == nums[right - 1]:
                    right -= 1
                left += 1
                right -= 1
            elif total < 0:
                left += 1
            else:
                right -= 1
    return result
```

### Product of Array Except Self

```python
def product_except_self(nums: list[int]) -> list[int]:
    """Product of all elements except self. O(n) time, O(1) extra space."""
    n = len(nums)
    result = [1] * n

    # Left products
    left_product = 1
    for i in range(n):
        result[i] = left_product
        left_product *= nums[i]

    # Right products
    right_product = 1
    for i in range(n - 1, -1, -1):
        result[i] *= right_product
        right_product *= nums[i]

    return result
```

### Container With Most Water

```python
def max_area(height: list[int]) -> int:
    """Find two lines that together with x-axis form the max water container. O(n)."""
    left, right = 0, len(height) - 1
    max_water = 0
    while left < right:
        width = right - left
        h = min(height[left], height[right])
        max_water = max(max_water, width * h)
        if height[left] < height[right]:
            left += 1
        else:
            right -= 1
    return max_water
```

---

## ⚠️ Common Mistakes

| Mistake | Example | Fix |
|---------|---------|-----|
| Off-by-one errors | `for i in range(len(arr))` accessing `arr[i+1]` | Check bounds: `range(len(arr)-1)` |
| Modifying while iterating | Deleting elements while looping | Use a copy, or iterate backwards |
| Not handling empty arrays | `arr[0]` on empty array | Check `if not arr:` first |
| Forgetting to sort for two pointers | Two pointer on unsorted array | Sort first: `arr.sort()` |
| Integer overflow in sum | Summing large arrays | Use `long` in Java, Python handles natively |
| Slice copies in Python | `arr2 = arr` (reference, not copy) | `arr2 = arr[:]` or `arr2 = arr.copy()` |
| Go slice pitfalls | Appending to slice shared by multiple variables | Use `copy()` when sharing backing array |

---

## 🔗 Related Topics

- [Complexity Analysis](../complexity.md) — Understanding array operation costs
- [Hash Tables](../hash-table/) — Often used with arrays (two sum, frequency counting)
- [Dynamic Programming](../dynamic-programming/) — Many DP problems use arrays as state
- [Trees](../tree/) — Arrays can represent heaps and segment trees
- [System Design](../../system-design/) — Array-backed ring buffers, circular queues
- [Languages: Python](../../languages/python/) — Python list internals
- [Languages: Go](../../languages/go/) — Go slices and arrays

---

> **Arrays are simple, but array problems aren't.** Master the core techniques (two pointers, sliding window, prefix sum) and you'll have tools that apply to dozens of different problems. The key is recognizing which pattern fits.
