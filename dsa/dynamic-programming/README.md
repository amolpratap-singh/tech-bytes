# Dynamic Programming

> **Solve complex problems by breaking them into overlapping subproblems.** Dynamic programming (DP) is an optimization technique that stores results of expensive computations to avoid redundant work — turning exponential brute force into polynomial efficiency.

---

## 📑 Table of Contents

- [What is Dynamic Programming?](#-what-is-dynamic-programming)
- [Top-Down vs Bottom-Up](#-top-down-vs-bottom-up)
- [Framework for Solving DP Problems](#-framework-for-solving-dp-problems)
- [1D Dynamic Programming](#-1d-dynamic-programming)
- [2D Dynamic Programming](#-2d-dynamic-programming)
- [String DP](#-string-dp)
- [State Machine DP](#-state-machine-dp)
- [Common Patterns and Templates](#-common-patterns-and-templates)
- [Tips for Interviews](#-tips-for-interviews)
- [Related Topics](#-related-topics)

---

## 🧠 What is Dynamic Programming?

DP applies when a problem has two properties:

1. **Overlapping subproblems:** The same subproblems are solved multiple times
2. **Optimal substructure:** The optimal solution contains optimal solutions to subproblems

```text
Fibonacci WITHOUT DP (O(2^n)):
                fib(5)
               /      \
          fib(4)       fib(3)      ← fib(3) computed twice!
         /    \        /    \
      fib(3)  fib(2) fib(2) fib(1)  ← fib(2) computed three times!
     /    \
  fib(2)  fib(1)

Fibonacci WITH DP (O(n)):
  fib(1)=1 → fib(2)=1 → fib(3)=2 → fib(4)=3 → fib(5)=5
  Each subproblem computed exactly once!
```

---

## ⬆️ Top-Down vs Bottom-Up

### Top-Down (Memoization)

Start with the original problem, recurse into subproblems, cache results.

```python
def fib_memo(n: int, memo: dict[int, int] = None) -> int:
    """Top-down Fibonacci with memoization. O(n) time, O(n) space."""
    if memo is None:
        memo = {}
    if n in memo:
        return memo[n]
    if n <= 1:
        return n
    memo[n] = fib_memo(n - 1, memo) + fib_memo(n - 2, memo)
    return memo[n]
```

### Bottom-Up (Tabulation)

Start with smallest subproblems, build up to the answer.

```python
def fib_tab(n: int) -> int:
    """Bottom-up Fibonacci with tabulation. O(n) time, O(n) space."""
    if n <= 1:
        return n
    dp = [0] * (n + 1)
    dp[1] = 1
    for i in range(2, n + 1):
        dp[i] = dp[i - 1] + dp[i - 2]
    return dp[n]

def fib_optimized(n: int) -> int:
    """Space-optimized Fibonacci. O(n) time, O(1) space."""
    if n <= 1:
        return n
    prev2, prev1 = 0, 1
    for _ in range(2, n + 1):
        prev2, prev1 = prev1, prev2 + prev1
    return prev1
```

### Comparison

| Aspect | Top-Down (Memoization) | Bottom-Up (Tabulation) |
|--------|----------------------|----------------------|
| Approach | Recursive + cache | Iterative |
| Computes | Only needed subproblems | All subproblems |
| Stack overflow risk | Yes (deep recursion) | No |
| Space optimization | Harder | Easier (often reduce to O(1)) |
| Easier to write | Usually (natural recursion) | Sometimes (need to figure out order) |

---

## 🔧 Framework for Solving DP Problems

### Step-by-Step

```text
1. DEFINE THE STATE
   What variables uniquely describe a subproblem?
   → dp[i] = "the answer for the first i elements"

2. DEFINE THE TRANSITION
   How does the current state relate to previous states?
   → dp[i] = dp[i-1] + dp[i-2]  (Fibonacci example)

3. DEFINE THE BASE CASE
   What are the smallest subproblems you can solve directly?
   → dp[0] = 0, dp[1] = 1

4. DEFINE THE ANSWER
   Which state gives you the final answer?
   → dp[n]

5. OPTIMIZE SPACE (optional)
   Can you reduce from O(n) to O(1) by only keeping recent states?
```

---

## 📏 1D Dynamic Programming

### Climbing Stairs

```text
You can climb 1 or 2 steps at a time. How many ways to reach step n?

State: dp[i] = number of ways to reach step i
Transition: dp[i] = dp[i-1] + dp[i-2]
Base: dp[0] = 1, dp[1] = 1
```

```python
def climb_stairs(n: int) -> int:
    """Number of ways to climb n stairs. O(n) time, O(1) space."""
    if n <= 2:
        return n
    prev2, prev1 = 1, 2
    for _ in range(3, n + 1):
        prev2, prev1 = prev1, prev2 + prev1
    return prev1
```

```go
func climbStairs(n int) int {
    if n <= 2 {
        return n
    }
    prev2, prev1 := 1, 2
    for i := 3; i <= n; i++ {
        prev2, prev1 = prev1, prev2+prev1
    }
    return prev1
}
```

### House Robber

```text
Can't rob two adjacent houses. Maximize total robbery.

State: dp[i] = max money robbing houses 0..i
Transition: dp[i] = max(dp[i-1], dp[i-2] + nums[i])
  → Either skip house i (dp[i-1]) or rob it (dp[i-2] + nums[i])
```

```python
def rob(nums: list[int]) -> int:
    """House robber. O(n) time, O(1) space."""
    if not nums:
        return 0
    if len(nums) == 1:
        return nums[0]

    prev2, prev1 = 0, 0
    for num in nums:
        prev2, prev1 = prev1, max(prev1, prev2 + num)
    return prev1
```

### Coin Change

```text
Find minimum coins to make amount. Coins can be used unlimited times.

State: dp[amount] = min coins needed for this amount
Transition: dp[i] = min(dp[i - coin] + 1) for each coin
Base: dp[0] = 0
```

```python
def coin_change(coins: list[int], amount: int) -> int:
    """Minimum coins to make amount. O(amount × len(coins))."""
    dp = [float('inf')] * (amount + 1)
    dp[0] = 0

    for i in range(1, amount + 1):
        for coin in coins:
            if coin <= i and dp[i - coin] != float('inf'):
                dp[i] = min(dp[i], dp[i - coin] + 1)

    return dp[amount] if dp[amount] != float('inf') else -1
```

```go
func coinChange(coins []int, amount int) int {
    dp := make([]int, amount+1)
    for i := 1; i <= amount; i++ {
        dp[i] = amount + 1 // Impossible sentinel
    }
    for i := 1; i <= amount; i++ {
        for _, coin := range coins {
            if coin <= i && dp[i-coin]+1 < dp[i] {
                dp[i] = dp[i-coin] + 1
            }
        }
    }
    if dp[amount] > amount {
        return -1
    }
    return dp[amount]
}
```

### Word Break

```text
Can the string be segmented into dictionary words?

State: dp[i] = True if s[0:i] can be segmented
Transition: dp[i] = True if dp[j] and s[j:i] in wordDict for some j < i
Base: dp[0] = True (empty string)
```

```python
def word_break(s: str, word_dict: list[str]) -> bool:
    """Can string be segmented into dictionary words? O(n² × k)."""
    words = set(word_dict)
    dp = [False] * (len(s) + 1)
    dp[0] = True

    for i in range(1, len(s) + 1):
        for j in range(i):
            if dp[j] and s[j:i] in words:
                dp[i] = True
                break

    return dp[len(s)]
```

---

## 📐 2D Dynamic Programming

### Unique Paths

```text
Grid from top-left to bottom-right, can only move right or down.

State: dp[i][j] = number of paths to (i, j)
Transition: dp[i][j] = dp[i-1][j] + dp[i][j-1]
Base: dp[0][j] = 1, dp[i][0] = 1
```

```python
def unique_paths(m: int, n: int) -> int:
    """Count unique paths in m×n grid. O(m × n) time, O(n) space."""
    dp = [1] * n
    for _ in range(1, m):
        for j in range(1, n):
            dp[j] += dp[j - 1]
    return dp[n - 1]
```

### Longest Common Subsequence (LCS)

```text
Find length of longest subsequence common to both strings.

State: dp[i][j] = LCS of text1[0:i] and text2[0:j]
Transition:
  If text1[i-1] == text2[j-1]: dp[i][j] = dp[i-1][j-1] + 1
  Else: dp[i][j] = max(dp[i-1][j], dp[i][j-1])
Base: dp[0][j] = 0, dp[i][0] = 0
```

```python
def longest_common_subsequence(text1: str, text2: str) -> int:
    """LCS length. O(m × n) time, O(m × n) space."""
    m, n = len(text1), len(text2)
    dp = [[0] * (n + 1) for _ in range(m + 1)]

    for i in range(1, m + 1):
        for j in range(1, n + 1):
            if text1[i - 1] == text2[j - 1]:
                dp[i][j] = dp[i - 1][j - 1] + 1
            else:
                dp[i][j] = max(dp[i - 1][j], dp[i][j - 1])

    return dp[m][n]
```

```go
func longestCommonSubsequence(text1, text2 string) int {
    m, n := len(text1), len(text2)
    dp := make([][]int, m+1)
    for i := range dp {
        dp[i] = make([]int, n+1)
    }
    for i := 1; i <= m; i++ {
        for j := 1; j <= n; j++ {
            if text1[i-1] == text2[j-1] {
                dp[i][j] = dp[i-1][j-1] + 1
            } else {
                dp[i][j] = max(dp[i-1][j], dp[i][j-1])
            }
        }
    }
    return dp[m][n]
}

func max(a, b int) int {
    if a > b {
        return a
    }
    return b
}
```

### Edit Distance (Levenshtein)

```text
Minimum operations (insert, delete, replace) to transform word1 to word2.

State: dp[i][j] = edit distance of word1[0:i] and word2[0:j]
Transition:
  If word1[i-1] == word2[j-1]: dp[i][j] = dp[i-1][j-1]
  Else: dp[i][j] = 1 + min(dp[i-1][j], dp[i][j-1], dp[i-1][j-1])
    Insert: dp[i][j-1], Delete: dp[i-1][j], Replace: dp[i-1][j-1]
```

```python
def edit_distance(word1: str, word2: str) -> int:
    """Minimum edit distance. O(m × n) time, O(m × n) space."""
    m, n = len(word1), len(word2)
    dp = [[0] * (n + 1) for _ in range(m + 1)]

    for i in range(m + 1):
        dp[i][0] = i
    for j in range(n + 1):
        dp[0][j] = j

    for i in range(1, m + 1):
        for j in range(1, n + 1):
            if word1[i - 1] == word2[j - 1]:
                dp[i][j] = dp[i - 1][j - 1]
            else:
                dp[i][j] = 1 + min(dp[i - 1][j], dp[i][j - 1], dp[i - 1][j - 1])

    return dp[m][n]
```

### 0/1 Knapsack

```text
Given items with weights and values, maximize value within weight capacity.
Each item can be taken at most once.

State: dp[i][w] = max value using items 0..i-1 with capacity w
Transition:
  Skip item: dp[i][w] = dp[i-1][w]
  Take item: dp[i][w] = dp[i-1][w - weight[i]] + value[i]
  dp[i][w] = max(skip, take if weight[i] <= w)
```

```python
def knapsack(weights: list[int], values: list[int], capacity: int) -> int:
    """0/1 Knapsack. O(n × capacity) time, O(capacity) space."""
    n = len(weights)
    dp = [0] * (capacity + 1)

    for i in range(n):
        # Traverse backwards to avoid using item i twice
        for w in range(capacity, weights[i] - 1, -1):
            dp[w] = max(dp[w], dp[w - weights[i]] + values[i])

    return dp[capacity]
```

---

## 🔤 String DP

### Longest Palindromic Substring

```python
def longest_palindrome(s: str) -> str:
    """Find longest palindromic substring. O(n²) time, O(n²) space."""
    n = len(s)
    if n < 2:
        return s

    # dp[i][j] = True if s[i:j+1] is a palindrome
    dp = [[False] * n for _ in range(n)]
    start, max_len = 0, 1

    # Base cases: single characters
    for i in range(n):
        dp[i][i] = True

    # Base cases: two characters
    for i in range(n - 1):
        if s[i] == s[i + 1]:
            dp[i][i + 1] = True
            start, max_len = i, 2

    # Fill for lengths 3+
    for length in range(3, n + 1):
        for i in range(n - length + 1):
            j = i + length - 1
            if s[i] == s[j] and dp[i + 1][j - 1]:
                dp[i][j] = True
                if length > max_len:
                    start, max_len = i, length

    return s[start:start + max_len]
```

### Expand from Center (Alternative, simpler)

```python
def longest_palindrome_expand(s: str) -> str:
    """Expand around center approach. O(n²) time, O(1) space."""
    def expand(left: int, right: int) -> str:
        while left >= 0 and right < len(s) and s[left] == s[right]:
            left -= 1
            right += 1
        return s[left + 1:right]

    result = ""
    for i in range(len(s)):
        # Odd length palindrome
        odd = expand(i, i)
        if len(odd) > len(result):
            result = odd
        # Even length palindrome
        even = expand(i, i + 1)
        if len(even) > len(result):
            result = even
    return result
```

---

## 🤖 State Machine DP

### Best Time to Buy and Sell Stock with Cooldown

```text
States: Holding, Not Holding (Cooldown), Not Holding (Ready)

Transitions:
  hold[i] = max(hold[i-1], ready[i-1] - prices[i])  # Keep holding or buy
  sold[i] = hold[i-1] + prices[i]                     # Sell today
  ready[i] = max(ready[i-1], sold[i-1])               # Cooldown or stay ready
```

```python
def max_profit_with_cooldown(prices: list[int]) -> int:
    """Stock buy/sell with cooldown. O(n) time, O(1) space."""
    if len(prices) < 2:
        return 0

    hold = -prices[0]   # Holding stock
    sold = 0            # Just sold
    ready = 0           # Ready to buy

    for i in range(1, len(prices)):
        prev_hold = hold
        prev_sold = sold
        prev_ready = ready

        hold = max(prev_hold, prev_ready - prices[i])
        sold = prev_hold + prices[i]
        ready = max(prev_ready, prev_sold)

    return max(sold, ready)
```

### Best Time to Buy and Sell Stock with K Transactions

```python
def max_profit_k_transactions(k: int, prices: list[int]) -> int:
    """At most k transactions. O(n × k) time, O(k) space."""
    if not prices or k == 0:
        return 0

    # If k >= n/2, unlimited transactions
    if k >= len(prices) // 2:
        return sum(max(0, prices[i + 1] - prices[i]) for i in range(len(prices) - 1))

    buy = [float('-inf')] * (k + 1)
    sell = [0] * (k + 1)

    for price in prices:
        for j in range(1, k + 1):
            buy[j] = max(buy[j], sell[j - 1] - price)
            sell[j] = max(sell[j], buy[j] + price)

    return sell[k]
```

---

## 📋 Common Patterns and Templates

| Pattern | Example Problems | Key Insight |
|---------|-----------------|-------------|
| **Linear DP** | Climbing stairs, house robber, coin change | dp[i] depends on dp[i-1], dp[i-2], etc. |
| **Grid DP** | Unique paths, minimum path sum | dp[i][j] depends on dp[i-1][j], dp[i][j-1] |
| **String DP** | LCS, edit distance, palindrome | dp[i][j] for two string indices |
| **Knapsack** | 0/1 knapsack, subset sum, coin change | Take/skip decision at each item |
| **Interval DP** | Matrix chain, burst balloons | dp[i][j] = best for range [i..j] |
| **State machine** | Stock problems, string matching | States with transitions |
| **Tree DP** | House robber III, diameter | dp on tree nodes |
| **Bitmask DP** | Traveling salesman, subset problems | dp[mask] = visited set as bitmask |

### Template: 1D DP

```python
def solve(arr):
    n = len(arr)
    dp = [base_value] * n  # or (n + 1)
    dp[0] = base_case

    for i in range(1, n):
        dp[i] = transition(dp, arr, i)

    return dp[n - 1]  # or max(dp), or dp[target]
```

### Template: 2D DP

```python
def solve(arr1, arr2):
    m, n = len(arr1), len(arr2)
    dp = [[base_value] * (n + 1) for _ in range(m + 1)]
    # Initialize base cases

    for i in range(1, m + 1):
        for j in range(1, n + 1):
            dp[i][j] = transition(dp, arr1, arr2, i, j)

    return dp[m][n]
```

---

## 💡 Tips for Interviews

1. **Start with brute force recursion** — Understand the recursive structure first
2. **Identify overlapping subproblems** — Draw the recursion tree, look for repeated calls
3. **Add memoization** — Cache results with a hash map or array
4. **Convert to bottom-up if asked** — Fill the DP table iteratively
5. **Optimize space** — If dp[i] only depends on dp[i-1] (and dp[i-2]), use O(1) space
6. **Trace through a small example** — Verify your transition with concrete values
7. **Watch out for off-by-one** — Clarify 0-indexed vs 1-indexed, check base cases
8. **Name your state clearly** — "dp[i] represents the minimum cost to reach step i"

**Common signals that DP applies:**
- "Find the minimum/maximum…"
- "Count the number of ways…"
- "Is it possible to…"
- "Find the longest/shortest…"
- "Given choices at each step…"
- Problem has optimal substructure (greedy doesn't work)

---

## 🔗 Related Topics

- [Complexity Analysis](../complexity.md) — DP reduces exponential to polynomial
- [Arrays](../arrays/) — Many DP problems operate on arrays
- [Hash Tables](../hash-table/) — Memoization uses hash tables
- [Graphs](../graph/) — Shortest path algorithms use DP concepts
- [Trees](../tree/) — Tree DP problems
- [System Design](../../system-design/) — Optimization in system design

---

> **DP is a skill, not a trick.** It takes practice to recognize DP problems and define states and transitions. Start with the classic problems (climbing stairs, coin change, LCS), understand the patterns, and gradually tackle harder ones. The framework (state → transition → base case → answer) works for every DP problem.
