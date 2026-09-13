# ⚡ Caching in Distributed Systems

> **Caching is the most impactful performance optimization in distributed systems.** By storing copies of frequently accessed data closer to where it's needed, caching reduces latency, decreases load on backend systems, and improves throughput — often by orders of magnitude.

---

## 📑 Table of Contents

- [What is Caching?](#-what-is-caching)
- [Cache Hierarchy](#-cache-hierarchy)
- [Caching Strategies](#-caching-strategies)
- [Cache Invalidation](#-cache-invalidation)
- [Distributed Caching](#-distributed-caching)
- [Cache Problems & Solutions](#-cache-problems--solutions)
- [Production Tips](#-production-tips)
- [Related Topics](#-related-topics)

---

## 🧠 What is Caching?

Caching stores copies of data in a faster storage layer so that future requests can be served more quickly. The principle is simple: **accessing data from memory is orders of magnitude faster than accessing it from disk or over the network.**

**Latency comparison:**

| Operation | Latency | Relative |
|-----------|---------|----------|
| L1 cache reference | ~1 ns | 1x |
| L2 cache reference | ~4 ns | 4x |
| Main memory (RAM) | ~100 ns | 100x |
| Redis/Memcached (network) | ~0.5 ms | 500,000x |
| SSD read | ~0.1 ms | 100,000x |
| HDD seek | ~10 ms | 10,000,000x |
| Network round trip (same DC) | ~0.5 ms | 500,000x |
| Network round trip (cross-region) | ~50 ms | 50,000,000x |

**When to use caching:**

- Read-heavy workloads (high read-to-write ratio)
- Expensive computations that are repeated frequently
- Data that doesn't change frequently (or staleness is acceptable)
- High-latency data sources (external APIs, distant databases)
- Hot data with skewed access patterns (20% of keys get 80% of traffic)

---

## 📊 Cache Hierarchy

```text
Cache Hierarchy (fastest → slowest, smallest → largest):

┌─────────────────────────────┐
│      CPU L1/L2/L3 Cache     │  ~1-10 ns    Managed by hardware
├─────────────────────────────┤
│      Application Cache       │  ~1 μs       In-process (local memory)
│      (HashMap, Guava, etc.)  │
├─────────────────────────────┤
│     Distributed Cache        │  ~0.5-2 ms   Network call to Redis/Memcached
│     (Redis, Memcached)       │
├─────────────────────────────┤
│          CDN Cache           │  ~5-50 ms    Edge servers near users
│     (CloudFront, Akamai)     │
├─────────────────────────────┤
│       Database Cache         │  ~1-10 ms    Buffer pool, query cache
│     (PostgreSQL, MySQL)      │
├─────────────────────────────┤
│       Origin/Database        │  ~10-100 ms  Full query execution
│     (PostgreSQL, MongoDB)    │
└─────────────────────────────┘
```

### Multi-Level Cache Example

```text
Request flow:
User → CDN → Application Cache → Redis → Database

1. Check CDN cache         → HIT? Return immediately
2. Check application cache → HIT? Return immediately
3. Check Redis             → HIT? Return, update app cache
4. Query database          → Return, update Redis + app cache + CDN
```

---

## 🔄 Caching Strategies

### Cache-Aside (Lazy Loading)

The application manages the cache directly. Most common strategy.

```text
Read:
┌──────┐     ┌───────┐     ┌──────────┐
│ App  │──1──│ Cache │     │ Database │
│      │◄─2──│       │     │          │
│      │     └───────┘     │          │
│      │──3───────────────►│          │  (cache miss)
│      │◄─4────────────────│          │
│      │──5──►┌───────┐   │          │
│      │      │ Cache │   │          │  (populate cache)
└──────┘      └───────┘   └──────────┘

Write:
┌──────┐     ┌───────┐     ┌──────────┐
│ App  │──1──────────────►│ Database │  (write to DB)
│      │──2──►┌───────┐   │          │
│      │      │ Cache │   │          │  (invalidate/update cache)
└──────┘      └───────┘   └──────────┘
```

```python
def get_user(user_id):
    # 1. Check cache
    cached = cache.get(f"user:{user_id}")
    if cached:
        return cached  # Cache hit
    
    # 2. Cache miss - query database
    user = db.query("SELECT * FROM users WHERE id = %s", user_id)
    
    # 3. Populate cache with TTL
    cache.set(f"user:{user_id}", user, ttl=300)
    
    return user

def update_user(user_id, data):
    # 1. Update database
    db.execute("UPDATE users SET ... WHERE id = %s", user_id, data)
    
    # 2. Invalidate cache
    cache.delete(f"user:{user_id}")
```

**Pros:** Simple, application controls caching logic, only caches what's needed
**Cons:** Cache miss penalty (extra round trip), potential for stale data

### Read-Through

The cache itself loads data from the database on a miss. Application only talks to cache.

```text
┌──────┐     ┌───────┐     ┌──────────┐
│ App  │──1──│ Cache │     │ Database │
│      │     │       │──2──│          │  (cache loads from DB on miss)
│      │◄─3──│       │◄─── │          │
└──────┘     └───────┘     └──────────┘
```

**Pros:** Simpler application code, cache handles loading logic
**Cons:** Cache must know about data source, cold start slower

### Write-Through

Every write goes through the cache to the database. Cache always has fresh data.

```text
┌──────┐     ┌───────┐     ┌──────────┐
│ App  │──1──│ Cache │──2──│ Database │
│      │◄─3──│       │◄────│          │
└──────┘     └───────┘     └──────────┘
             (cache + DB written synchronously)
```

**Pros:** Cache always consistent with database, simple reads
**Cons:** Write latency (two writes per operation), cache may hold rarely-read data

### Write-Behind (Write-Back)

Writes go to cache immediately, then asynchronously flushed to database.

```text
┌──────┐     ┌───────┐  async  ┌──────────┐
│ App  │──1──│ Cache │ ──2──► │ Database │
│      │◄─── │       │        │          │
└──────┘     └───────┘        └──────────┘
             (immediate ack,   (batched writes)
              DB updated later)
```

**Pros:** Very low write latency, can batch database writes
**Cons:** Data loss risk if cache crashes before flush, complexity

### Refresh-Ahead

Cache proactively refreshes entries before they expire.

```text
Key TTL: 60s
Refresh at: 50s (refresh-ahead factor = 0.83)

Time 0s:  [Cache SET, TTL=60s]
Time 50s: [Background refresh triggered] → DB query → cache updated
Time 60s: [Key would have expired, but was already refreshed]

Result: No cache miss for frequently accessed keys
```

**Pros:** No cache miss for hot data, predictable latency
**Cons:** Wastes resources refreshing data that may not be read, complexity

### Strategy Comparison

| Strategy | Read Latency | Write Latency | Consistency | Complexity | Best For |
|----------|-------------|--------------|-------------|------------|----------|
| **Cache-Aside** | Miss: high | Low | Eventual | Low | General purpose |
| **Read-Through** | Miss: high | Low | Eventual | Medium | Read-heavy |
| **Write-Through** | Low (always hit) | High (2x) | Strong | Medium | Read-heavy, consistency needed |
| **Write-Behind** | Low | Very low | Eventual | High | Write-heavy |
| **Refresh-Ahead** | Very low | Low | Near real-time | High | Hot data with predictable patterns |

---

## 🗑 Cache Invalidation

> "There are only two hard things in Computer Science: cache invalidation and naming things." — Phil Karlton

### TTL-Based (Time-to-Live)

```python
# Set with TTL
cache.set("user:123", user_data, ttl=300)  # Expires in 5 minutes

# Different TTLs for different data
CACHE_TTLS = {
    "user_profile": 3600,      # 1 hour (rarely changes)
    "user_session": 1800,      # 30 minutes
    "product_price": 60,       # 1 minute (changes often)
    "static_content": 86400,   # 24 hours
}
```

**Pros:** Simple, automatic cleanup, bounded staleness
**Cons:** Data can be stale for up to TTL duration, thundering herd at expiry

### Event-Based Invalidation

```text
Database change → Event (Kafka/CDC) → Cache invalidation

┌──────────┐     ┌───────┐     ┌───────┐     ┌───────┐
│ Database │──►  │ Kafka │──►  │ App   │──►  │ Cache │
│ (CDC)    │     │ Topic │     │ (sub) │     │ (del) │
└──────────┘     └───────┘     └───────┘     └───────┘
```

**Pros:** Near real-time invalidation, no polling
**Cons:** Infrastructure complexity (CDC, message broker), eventual consistency

### Version-Based

```python
# Include version in cache key
version = db.get_version("products")  # e.g., "v42"
cache_key = f"products:{version}"

# When data changes, increment version
db.increment_version("products")  # Now "v43"
# Old cache key "products:v42" becomes unused
# New requests use "products:v43" → cache miss → fresh data
```

**Pros:** No explicit invalidation needed, atomic switch
**Cons:** Old versions occupy memory until evicted, version management overhead

---

## 🌐 Distributed Caching

### Redis vs Memcached

| Feature | Redis | Memcached |
|---------|-------|-----------|
| **Data structures** | Strings, Lists, Sets, Hashes, Sorted Sets, Streams | Strings only |
| **Persistence** | RDB + AOF | None |
| **Replication** | Master-replica | None (client-side) |
| **Cluster mode** | Hash slots (16384) | Client-side consistent hashing |
| **Pub/Sub** | Yes | No |
| **Lua scripting** | Yes | No |
| **Max value size** | 512 MB | 1 MB (default) |
| **Threading** | Single-threaded (I/O threads in 6.0+) | Multi-threaded |
| **Use case** | Cache + data store + message broker | Pure caching |

### Consistent Hashing

Distributes keys across cache nodes while minimizing redistribution when nodes are added/removed:

```text
Hash ring (0 to 2^32):

        Node A (hash: 100)
           ╱
          ○
     ╱         ╲
  ○               ○  ← Node B (hash: 300)
  Node D            
  (hash: 900)  ╲   ╱
                ○
           Node C (hash: 600)

Key "user:123" → hash = 250 → clockwise → Node B
Key "user:456" → hash = 700 → clockwise → Node D

Adding Node E (hash: 500):
- Only keys between C(600) and E(500) need to move
- Other keys stay in place
```

**Virtual nodes** improve balance: each physical node gets multiple points on the ring.

### Cache Sharding

```text
Strategy 1: Hash-based
  shard = hash(key) % num_shards
  Simple but resharding moves many keys

Strategy 2: Consistent hashing
  Minimal key movement during scaling
  Used by Redis Cluster (hash slots)

Strategy 3: Range-based
  Keys a-m → shard 1, n-z → shard 2
  Can cause hot spots
```

---

## 🐛 Cache Problems & Solutions

### Cache Stampede (Thundering Herd)

**Problem:** Many requests arrive simultaneously for an expired key, all causing cache misses and hitting the database.

```text
Time 0: Key "popular-item" TTL expires
Time 0.001s: 1000 concurrent requests → all cache miss → 1000 DB queries!

              ┌───────────┐
              │ App × 1000│──── all miss ────► Database (OVERLOADED)
              └───────────┘
```

**Solutions:**

```python
# Solution 1: Locking (only one request loads from DB)
def get_with_lock(key):
    value = cache.get(key)
    if value:
        return value
    
    lock_key = f"lock:{key}"
    if cache.set(lock_key, "1", nx=True, ttl=10):  # Acquire lock
        value = db.query(key)
        cache.set(key, value, ttl=300)
        cache.delete(lock_key)
        return value
    else:
        time.sleep(0.1)  # Wait for lock holder to populate
        return cache.get(key)  # Retry

# Solution 2: Probabilistic early refresh
def get_with_early_refresh(key, ttl=300, beta=1.0):
    value, expiry = cache.get_with_ttl(key)
    remaining_ttl = expiry - time.time()
    
    # Probabilistically refresh before expiry
    if remaining_ttl - beta * math.log(random.random()) <= 0:
        value = db.query(key)
        cache.set(key, value, ttl=ttl)
    
    return value
```

### Cache Penetration

**Problem:** Requests for data that doesn't exist bypass the cache and always hit the database.

```text
Request for user_id=999999 (doesn't exist)
→ Cache miss → DB query → not found → no cache write → repeat forever!
```

**Solutions:**

```python
# Solution 1: Cache null values
def get_user(user_id):
    cached = cache.get(f"user:{user_id}")
    if cached == "NULL_MARKER":
        return None  # Known to not exist
    if cached:
        return cached
    
    user = db.query_user(user_id)
    if user:
        cache.set(f"user:{user_id}", user, ttl=300)
    else:
        cache.set(f"user:{user_id}", "NULL_MARKER", ttl=60)  # Short TTL
    return user

# Solution 2: Bloom filter (pre-check)
# Before querying, check if key could exist
bloom_filter = BloomFilter(expected_items=1_000_000, fp_rate=0.01)

def get_user(user_id):
    if not bloom_filter.might_contain(user_id):
        return None  # Definitely doesn't exist
    # ... normal cache/DB lookup
```

### Cache Avalanche

**Problem:** Many cache keys expire at the same time, causing a massive spike in database load.

```text
Time 0: Bulk cache population (all with TTL=3600)
Time 3600: ALL keys expire simultaneously → DB overloaded
```

**Solutions:**

```python
# Solution 1: Jittered TTL
import random

def set_with_jitter(key, value, base_ttl=3600):
    jitter = random.randint(0, 600)  # 0-10 minutes random jitter
    cache.set(key, value, ttl=base_ttl + jitter)

# Solution 2: Staggered warm-up
# Pre-warm cache in batches, not all at once

# Solution 3: Circuit breaker
# If DB load exceeds threshold, serve stale cache data
```

### Cold Start

**Problem:** After a restart or new deployment, the cache is empty and all requests hit the database.

**Solutions:**

- **Pre-warming:** Load popular keys into cache before serving traffic
- **Gradual rollout:** Shift traffic gradually to new instances
- **Cache serialization:** Dump cache to disk before restart, reload after

```python
# Pre-warm cache at startup
def warm_cache():
    popular_keys = db.query("SELECT id FROM users ORDER BY access_count DESC LIMIT 1000")
    for key in popular_keys:
        user = db.query_user(key)
        cache.set(f"user:{key}", user, ttl=3600)
```

---

## 🏭 Production Tips

### Sizing

| Factor | Consideration |
|--------|---------------|
| **Working set** | Cache should hold your most frequently accessed data |
| **Memory** | Monitor memory usage; set `maxmemory` in Redis |
| **Hit rate** | Target 90%+ hit rate for most applications |
| **Key count** | Each key has overhead (~70 bytes in Redis); factor this in |

### Eviction Policies

| Policy | Description | Best For |
|--------|-------------|----------|
| **LRU** (Least Recently Used) | Evict least recently accessed | General purpose |
| **LFU** (Least Frequently Used) | Evict least frequently accessed | Stable access patterns |
| **TTL** | Evict by expiration time | Time-sensitive data |
| **Random** | Evict random keys | When all keys equally important |
| **noeviction** | Return errors when full | When data loss is unacceptable |

### Monitoring Checklist

```text
□ Cache hit rate (target: > 90%)
□ Cache miss rate and miss reasons
□ Memory usage and eviction rate
□ Latency (p50, p95, p99)
□ Connection pool utilization
□ Key expiration rate
□ Biggest keys (memory hogs)
□ Hottest keys (access pattern)
□ Error rate (connection failures, timeouts)
```

### Common Patterns

```python
# Pattern: Cache key naming convention
# {entity}:{id}:{field}
"user:123:profile"
"user:123:settings"
"product:456:details"
"session:abc-def-ghi"

# Pattern: Cache warming on deployment
# 1. Deploy new version
# 2. Warm cache with popular keys
# 3. Shift traffic gradually

# Pattern: Multi-get for batch reads
users = cache.mget(["user:1", "user:2", "user:3"])
missing = [k for k, v in zip(keys, users) if v is None]
if missing:
    db_results = db.query_users(missing)
    cache.mset({f"user:{u.id}": u for u in db_results})
```

### Anti-Patterns

| Anti-Pattern | Problem | Fix |
|-------------|---------|-----|
| Caching everything | Wastes memory, low hit rate | Cache only hot/expensive data |
| No TTL | Stale data forever | Always set a TTL |
| Cache as primary store | Data loss on failure | Database is source of truth |
| Huge values | Slow serialization, memory | Compress or split large values |
| No monitoring | Blind to performance issues | Monitor hit rate, latency, memory |
| Same TTL everywhere | Avalanche risk | Use jittered TTLs |

---

## 🔗 Related Topics

- [Databases — Redis](../../databases/redis/) — Redis as a distributed cache and data structure store
- [Distributed Systems — Consensus](../consensus/) — Consistency models affecting cache design
- [Distributed Systems — Kafka](../kafka/) — Event-driven cache invalidation
- [System Design — Caching](../../system-design/caching.md) — Caching in system design interviews
- [DevOps — Kubernetes](../../devops/kubernetes/) — Caching in containerized environments

---

> **Caching is not free.** Every cache introduces complexity: invalidation logic, consistency trade-offs, memory management, and operational overhead. Start with the simplest strategy that meets your requirements, measure its effectiveness, and optimize from there.
