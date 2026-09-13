# Graphs

> **Model relationships and connections between entities.** Graphs represent networks of nodes (vertices) connected by edges — used in social networks, maps, dependency resolution, network routing, and countless other applications.

---

## 📑 Table of Contents

- [What are Graphs?](#-what-are-graphs)
- [Representations](#-representations)
- [BFS (Breadth-First Search)](#-bfs-breadth-first-search)
- [DFS (Depth-First Search)](#-dfs-depth-first-search)
- [Dijkstra's Algorithm](#-dijkstras-algorithm)
- [Topological Sort](#-topological-sort)
- [Union-Find (Disjoint Set)](#-union-find-disjoint-set)
- [Minimum Spanning Tree](#-minimum-spanning-tree)
- [Common Patterns](#-common-patterns)
- [Related Topics](#-related-topics)

---

## 🧠 What are Graphs?

A graph G = (V, E) consists of a set of **vertices** (nodes) V and a set of **edges** E connecting them.

### Types of Graphs

```text
Undirected:          Directed:            Weighted:
  A --- B            A → B                A --3-- B
  |     |            ↑   ↓                |       |
  C --- D            C ← D                5       2
                                          |       |
                                          C --4-- D

DAG (Directed Acyclic Graph):    Cyclic:
  A → B → D                      A → B
  ↓       ↑                      ↑   ↓
  C ------+                      D ← C
```

| Type | Description | Example |
|------|-------------|---------|
| **Undirected** | Edges have no direction | Social network (friendship) |
| **Directed** | Edges have direction | Twitter follows, dependencies |
| **Weighted** | Edges have costs/weights | Road network (distances) |
| **DAG** | Directed, no cycles | Build dependencies, course prerequisites |
| **Cyclic** | Contains at least one cycle | Circular dependencies |
| **Connected** | Path exists between all vertex pairs | Single network component |
| **Bipartite** | Vertices split into two groups, edges only between groups | Matching problems |

---

## 📐 Representations

### Adjacency List

Each vertex stores a list of its neighbors. Most common representation.

```python
# Unweighted graph
graph = {
    'A': ['B', 'C'],
    'B': ['A', 'D'],
    'C': ['A', 'D'],
    'D': ['B', 'C'],
}

# Weighted graph
weighted_graph = {
    'A': [('B', 3), ('C', 5)],
    'B': [('A', 3), ('D', 2)],
    'C': [('A', 5), ('D', 4)],
    'D': [('B', 2), ('C', 4)],
}

# Using defaultdict for easier building
from collections import defaultdict
graph = defaultdict(list)
edges = [(0, 1), (0, 2), (1, 3), (2, 3)]
for u, v in edges:
    graph[u].append(v)
    graph[v].append(u)  # Undirected
```

```go
// Adjacency list in Go
graph := map[int][]int{
    0: {1, 2},
    1: {0, 3},
    2: {0, 3},
    3: {1, 2},
}

// Weighted adjacency list
type Edge struct {
    To     int
    Weight int
}
weightedGraph := map[int][]Edge{
    0: {{1, 3}, {2, 5}},
    1: {{0, 3}, {3, 2}},
}
```

### Adjacency Matrix

2D array where `matrix[i][j]` indicates an edge between vertex i and j.

```python
# Adjacency matrix
# Vertices: 0, 1, 2, 3
matrix = [
    [0, 1, 1, 0],  # 0 connects to 1, 2
    [1, 0, 0, 1],  # 1 connects to 0, 3
    [1, 0, 0, 1],  # 2 connects to 0, 3
    [0, 1, 1, 0],  # 3 connects to 1, 2
]
```

### Edge List

List of edges. Simplest representation.

```python
edges = [(0, 1), (0, 2), (1, 3), (2, 3)]
weighted_edges = [(0, 1, 3), (0, 2, 5), (1, 3, 2), (2, 3, 4)]
```

### Representation Comparison

| Feature | Adjacency List | Adjacency Matrix | Edge List |
|---------|---------------|-----------------|-----------|
| Space | O(V + E) | O(V²) | O(E) |
| Check edge exists | O(degree) | O(1) | O(E) |
| Find neighbors | O(degree) | O(V) | O(E) |
| Add edge | O(1) | O(1) | O(1) |
| Best for | Sparse graphs | Dense graphs | Edge-centric algorithms |
| Most common? | ✅ Yes | Rare (dense graphs only) | Kruskal's MST |

---

## 🔵 BFS (Breadth-First Search)

Explore all neighbors at current depth before moving to next depth. Uses a **queue**.

### Implementation

```python
from collections import deque

def bfs(graph: dict, start: int) -> list[int]:
    """BFS traversal. O(V + E)."""
    visited = set([start])
    queue = deque([start])
    order = []

    while queue:
        node = queue.popleft()
        order.append(node)
        for neighbor in graph[node]:
            if neighbor not in visited:
                visited.add(neighbor)
                queue.append(neighbor)

    return order
```

```go
func bfs(graph map[int][]int, start int) []int {
    visited := map[int]bool{start: true}
    queue := []int{start}
    var order []int

    for len(queue) > 0 {
        node := queue[0]
        queue = queue[1:]
        order = append(order, node)
        for _, neighbor := range graph[node] {
            if !visited[neighbor] {
                visited[neighbor] = true
                queue = append(queue, neighbor)
            }
        }
    }
    return order
}
```

### Shortest Path (Unweighted Graph)

```python
def shortest_path_bfs(graph: dict, start: int, end: int) -> list[int]:
    """Find shortest path in unweighted graph using BFS. O(V + E)."""
    if start == end:
        return [start]

    visited = set([start])
    queue = deque([(start, [start])])

    while queue:
        node, path = queue.popleft()
        for neighbor in graph[node]:
            if neighbor == end:
                return path + [neighbor]
            if neighbor not in visited:
                visited.add(neighbor)
                queue.append((neighbor, path + [neighbor]))

    return []  # No path found
```

**BFS guarantees shortest path in unweighted graphs** because it explores all nodes at distance d before exploring nodes at distance d+1.

---

## 🔴 DFS (Depth-First Search)

Explore as deep as possible before backtracking. Uses a **stack** (or recursion).

### Implementation (Recursive)

```python
def dfs(graph: dict, start: int) -> list[int]:
    """DFS traversal (recursive). O(V + E)."""
    visited = set()
    order = []

    def _dfs(node: int) -> None:
        visited.add(node)
        order.append(node)
        for neighbor in graph[node]:
            if neighbor not in visited:
                _dfs(neighbor)

    _dfs(start)
    return order
```

### Implementation (Iterative)

```python
def dfs_iterative(graph: dict, start: int) -> list[int]:
    """DFS traversal (iterative). O(V + E)."""
    visited = set()
    stack = [start]
    order = []

    while stack:
        node = stack.pop()
        if node in visited:
            continue
        visited.add(node)
        order.append(node)
        for neighbor in reversed(graph[node]):
            if neighbor not in visited:
                stack.append(neighbor)

    return order
```

```go
func dfs(graph map[int][]int, start int) []int {
    visited := map[int]bool{}
    var order []int

    var dfsHelper func(node int)
    dfsHelper = func(node int) {
        visited[node] = true
        order = append(order, node)
        for _, neighbor := range graph[node] {
            if !visited[neighbor] {
                dfsHelper(neighbor)
            }
        }
    }

    dfsHelper(start)
    return order
}
```

### Cycle Detection (Directed Graph)

```python
def has_cycle(graph: dict[int, list[int]], num_nodes: int) -> bool:
    """Detect cycle in directed graph using DFS. O(V + E)."""
    WHITE, GRAY, BLACK = 0, 1, 2
    color = [WHITE] * num_nodes

    def dfs(node: int) -> bool:
        color[node] = GRAY  # Being processed
        for neighbor in graph.get(node, []):
            if color[neighbor] == GRAY:
                return True  # Back edge → cycle!
            if color[neighbor] == WHITE:
                if dfs(neighbor):
                    return True
        color[node] = BLACK  # Fully processed
        return False

    for node in range(num_nodes):
        if color[node] == WHITE:
            if dfs(node):
                return True
    return False
```

### Connected Components

```python
def count_components(n: int, edges: list[tuple[int, int]]) -> int:
    """Count connected components in undirected graph. O(V + E)."""
    graph = defaultdict(list)
    for u, v in edges:
        graph[u].append(v)
        graph[v].append(u)

    visited = set()
    components = 0

    def dfs(node: int) -> None:
        visited.add(node)
        for neighbor in graph[node]:
            if neighbor not in visited:
                dfs(neighbor)

    for node in range(n):
        if node not in visited:
            dfs(node)
            components += 1

    return components
```

---

## 📏 Dijkstra's Algorithm

Find the shortest path in a **weighted graph with non-negative weights**.

```python
import heapq

def dijkstra(graph: dict[int, list[tuple[int, int]]], start: int) -> dict[int, int]:
    """Dijkstra's shortest path. O((V + E) log V) with min-heap."""
    distances = {start: 0}
    min_heap = [(0, start)]  # (distance, node)

    while min_heap:
        dist, node = heapq.heappop(min_heap)

        if dist > distances.get(node, float('inf')):
            continue  # Already found a shorter path

        for neighbor, weight in graph.get(node, []):
            new_dist = dist + weight
            if new_dist < distances.get(neighbor, float('inf')):
                distances[neighbor] = new_dist
                heapq.heappush(min_heap, (new_dist, neighbor))

    return distances
```

```go
func dijkstra(graph map[int][]Edge, start int) map[int]int {
    distances := map[int]int{start: 0}
    // MinHeap of (distance, node)
    h := &MinHeap{{0, start}}
    heap.Init(h)

    for h.Len() > 0 {
        item := heap.Pop(h).(Edge)
        dist, node := item.Weight, item.To

        if d, ok := distances[node]; ok && dist > d {
            continue
        }

        for _, edge := range graph[node] {
            newDist := dist + edge.Weight
            if d, ok := distances[edge.To]; !ok || newDist < d {
                distances[edge.To] = newDist
                heap.Push(h, Edge{edge.To, newDist})
            }
        }
    }
    return distances
}
```

**Key points:**
- Only works with **non-negative** weights (use Bellman-Ford for negative weights)
- Greedy: always processes the closest unvisited vertex
- Time: O((V + E) log V) with binary heap

### Dijkstra with Path Reconstruction

```python
def dijkstra_with_path(graph: dict, start: int, end: int) -> tuple[int, list[int]]:
    """Return shortest distance and path. O((V + E) log V)."""
    distances = {start: 0}
    previous = {start: None}
    min_heap = [(0, start)]

    while min_heap:
        dist, node = heapq.heappop(min_heap)
        if node == end:
            break
        if dist > distances.get(node, float('inf')):
            continue
        for neighbor, weight in graph.get(node, []):
            new_dist = dist + weight
            if new_dist < distances.get(neighbor, float('inf')):
                distances[neighbor] = new_dist
                previous[neighbor] = node
                heapq.heappush(min_heap, (new_dist, neighbor))

    # Reconstruct path
    if end not in distances:
        return float('inf'), []
    path = []
    node = end
    while node is not None:
        path.append(node)
        node = previous[node]
    return distances[end], list(reversed(path))
```

---

## 📊 Topological Sort

Linear ordering of vertices in a DAG such that for every edge u → v, u comes before v.

### Kahn's Algorithm (BFS-based)

```python
from collections import deque

def topological_sort_kahn(graph: dict[int, list[int]], num_nodes: int) -> list[int]:
    """Kahn's algorithm for topological sort. O(V + E)."""
    # Calculate in-degrees
    in_degree = [0] * num_nodes
    for node in graph:
        for neighbor in graph[node]:
            in_degree[neighbor] += 1

    # Start with nodes having zero in-degree
    queue = deque([i for i in range(num_nodes) if in_degree[i] == 0])
    result = []

    while queue:
        node = queue.popleft()
        result.append(node)
        for neighbor in graph.get(node, []):
            in_degree[neighbor] -= 1
            if in_degree[neighbor] == 0:
                queue.append(neighbor)

    if len(result) != num_nodes:
        return []  # Cycle detected — topological sort impossible
    return result
```

### DFS-based Topological Sort

```python
def topological_sort_dfs(graph: dict[int, list[int]], num_nodes: int) -> list[int]:
    """DFS-based topological sort. O(V + E)."""
    visited = set()
    stack = []

    def dfs(node: int) -> None:
        visited.add(node)
        for neighbor in graph.get(node, []):
            if neighbor not in visited:
                dfs(neighbor)
        stack.append(node)  # Add after all descendants processed

    for node in range(num_nodes):
        if node not in visited:
            dfs(node)

    return list(reversed(stack))
```

```go
func topologicalSort(graph map[int][]int, numNodes int) []int {
    visited := make(map[int]bool)
    var stack []int

    var dfs func(node int)
    dfs = func(node int) {
        visited[node] = true
        for _, neighbor := range graph[node] {
            if !visited[neighbor] {
                dfs(neighbor)
            }
        }
        stack = append(stack, node)
    }

    for i := 0; i < numNodes; i++ {
        if !visited[i] {
            dfs(i)
        }
    }

    // Reverse the stack
    for i, j := 0, len(stack)-1; i < j; i, j = i+1, j-1 {
        stack[i], stack[j] = stack[j], stack[i]
    }
    return stack
}
```

**Use cases:** Build systems (make), course scheduling, dependency resolution, task ordering.

---

## 🤝 Union-Find (Disjoint Set)

Efficiently track which elements belong to the same group. Two operations: **union** (merge groups) and **find** (which group is an element in).

```python
class UnionFind:
    """Disjoint Set with path compression and union by rank. Near O(1) amortized."""

    def __init__(self, n: int):
        self.parent = list(range(n))
        self.rank = [0] * n
        self.count = n  # Number of components

    def find(self, x: int) -> int:
        """Find root with path compression. Near O(1) amortized."""
        if self.parent[x] != x:
            self.parent[x] = self.find(self.parent[x])  # Path compression
        return self.parent[x]

    def union(self, x: int, y: int) -> bool:
        """Union by rank. Returns False if already connected."""
        root_x = self.find(x)
        root_y = self.find(y)
        if root_x == root_y:
            return False
        # Union by rank
        if self.rank[root_x] < self.rank[root_y]:
            root_x, root_y = root_y, root_x
        self.parent[root_y] = root_x
        if self.rank[root_x] == self.rank[root_y]:
            self.rank[root_x] += 1
        self.count -= 1
        return True

    def connected(self, x: int, y: int) -> bool:
        return self.find(x) == self.find(y)
```

```go
type UnionFind struct {
    Parent []int
    Rank   []int
    Count  int
}

func NewUnionFind(n int) *UnionFind {
    parent := make([]int, n)
    rank := make([]int, n)
    for i := range parent {
        parent[i] = i
    }
    return &UnionFind{Parent: parent, Rank: rank, Count: n}
}

func (uf *UnionFind) Find(x int) int {
    if uf.Parent[x] != x {
        uf.Parent[x] = uf.Find(uf.Parent[x])
    }
    return uf.Parent[x]
}

func (uf *UnionFind) Union(x, y int) bool {
    rootX, rootY := uf.Find(x), uf.Find(y)
    if rootX == rootY {
        return false
    }
    if uf.Rank[rootX] < uf.Rank[rootY] {
        rootX, rootY = rootY, rootX
    }
    uf.Parent[rootY] = rootX
    if uf.Rank[rootX] == uf.Rank[rootY] {
        uf.Rank[rootX]++
    }
    uf.Count--
    return true
}
```

**Complexity with both optimizations:** O(α(n)) per operation where α is the inverse Ackermann function — effectively O(1).

**Use cases:** Connected components, cycle detection in undirected graphs, Kruskal's MST, network connectivity.

---

## 🌲 Minimum Spanning Tree

A tree that connects all vertices with minimum total edge weight.

### Kruskal's Algorithm

Sort edges by weight, add edges greedily using Union-Find to avoid cycles.

```python
def kruskal(num_nodes: int, edges: list[tuple[int, int, int]]) -> list[tuple[int, int, int]]:
    """Kruskal's MST. O(E log E)."""
    edges.sort(key=lambda e: e[2])  # Sort by weight
    uf = UnionFind(num_nodes)
    mst = []

    for u, v, weight in edges:
        if uf.union(u, v):
            mst.append((u, v, weight))
            if len(mst) == num_nodes - 1:
                break

    return mst
```

### Prim's Algorithm

Start from any vertex, greedily add the cheapest edge connecting the tree to an outside vertex.

```python
def prim(graph: dict[int, list[tuple[int, int]]], start: int = 0) -> list[tuple[int, int, int]]:
    """Prim's MST. O((V + E) log V)."""
    visited = set([start])
    edges = [(weight, start, neighbor) for neighbor, weight in graph[start]]
    heapq.heapify(edges)
    mst = []

    while edges and len(visited) < len(graph):
        weight, u, v = heapq.heappop(edges)
        if v in visited:
            continue
        visited.add(v)
        mst.append((u, v, weight))
        for neighbor, w in graph[v]:
            if neighbor not in visited:
                heapq.heappush(edges, (w, v, neighbor))

    return mst
```

| Algorithm | Time | Space | Best For |
|-----------|------|-------|----------|
| **Kruskal** | O(E log E) | O(V) | Sparse graphs, edge list |
| **Prim** | O((V + E) log V) | O(V) | Dense graphs, adjacency list |

---

## 🔄 Common Patterns

### Number of Islands (Grid DFS/BFS)

```python
def num_islands(grid: list[list[str]]) -> int:
    """Count islands in a grid. O(rows × cols)."""
    if not grid:
        return 0

    rows, cols = len(grid), len(grid[0])
    count = 0

    def dfs(r: int, c: int) -> None:
        if r < 0 or r >= rows or c < 0 or c >= cols or grid[r][c] != '1':
            return
        grid[r][c] = '0'  # Mark visited
        dfs(r + 1, c)
        dfs(r - 1, c)
        dfs(r, c + 1)
        dfs(r, c - 1)

    for r in range(rows):
        for c in range(cols):
            if grid[r][c] == '1':
                dfs(r, c)
                count += 1

    return count
```

### Course Schedule (Cycle Detection + Topological Sort)

```python
def can_finish(num_courses: int, prerequisites: list[list[int]]) -> bool:
    """Check if all courses can be finished (no cycle). O(V + E)."""
    graph = defaultdict(list)
    in_degree = [0] * num_courses

    for course, prereq in prerequisites:
        graph[prereq].append(course)
        in_degree[course] += 1

    queue = deque([i for i in range(num_courses) if in_degree[i] == 0])
    count = 0

    while queue:
        node = queue.popleft()
        count += 1
        for neighbor in graph[node]:
            in_degree[neighbor] -= 1
            if in_degree[neighbor] == 0:
                queue.append(neighbor)

    return count == num_courses
```

### Network Delay Time (Dijkstra)

```python
def network_delay_time(times: list[list[int]], n: int, k: int) -> int:
    """Find time for signal to reach all nodes from k. O((V + E) log V)."""
    graph = defaultdict(list)
    for u, v, w in times:
        graph[u].append((v, w))

    distances = dijkstra(graph, k)

    if len(distances) < n:
        return -1  # Not all nodes reachable
    return max(distances.values())
```

---

## 🔗 Related Topics

- [Complexity Analysis](../complexity.md) — Graph algorithm complexities
- [Trees](../tree/) — Trees are special cases of graphs (connected, acyclic)
- [Dynamic Programming](../dynamic-programming/) — DP on graphs (shortest paths, etc.)
- [Hash Tables](../hash-table/) — Adjacency lists use hash maps
- [System Design](../../system-design/) — Graph algorithms in distributed systems
- [Distributed Systems](../../distributed-systems/) — Network topology, consensus
- [Databases](../../databases/) — Query planners use graph algorithms

---

> **Graphs model the real world.** Social networks, road maps, internet routing, dependency management, recommendation engines — once you see graphs, you see them everywhere. Master BFS, DFS, Dijkstra, and Union-Find, and you can solve the vast majority of graph problems.
