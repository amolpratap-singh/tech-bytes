# ⚡ Redis

> **Redis is an in-memory data structure store used as a database, cache, message broker, and streaming engine.** Its sub-millisecond latency, rich data structures, and versatile functionality make it one of the most widely used technologies in modern infrastructure.

---

## 📑 Table of Contents

- [What is Redis?](#-what-is-redis)
- [Data Types & Commands](#-data-types--commands)
- [Key Management](#-key-management)
- [Pub/Sub](#-pubsub)
- [Streams](#-streams)
- [Transactions](#-transactions)
- [Lua Scripting](#-lua-scripting)
- [Persistence](#-persistence)
- [Replication](#-replication)
- [Cluster Mode](#-cluster-mode)
- [Performance](#-performance)
- [Common Patterns](#-common-patterns)
- [Troubleshooting](#-troubleshooting)
- [Production Tips](#-production-tips)
- [Related Topics](#-related-topics)

---

## 🧠 What is Redis?

Redis (Remote Dictionary Server) was created by Salvatore Sanfilippo in 2009. It stores data entirely in memory, providing extremely fast access times.

**Why Redis is fast:**

- **In-memory** — All data resides in RAM (100-1000x faster than disk)
- **Single-threaded** — No locking overhead for data operations
- **Efficient data structures** — Optimized C implementations
- **I/O multiplexing** — Handles thousands of connections concurrently (epoll/kqueue)
- **I/O threads** — Redis 6.0+ uses threads for network I/O (data ops still single-threaded)

**Common use cases:**

| Use Case | Why Redis | Data Structure |
|----------|----------|----------------|
| **Caching** | Sub-ms latency, TTL support | String, Hash |
| **Session store** | Fast reads, automatic expiry | String, Hash |
| **Rate limiting** | Atomic counters, TTL | String (INCR) |
| **Leaderboard** | Sorted sets with scores | Sorted Set |
| **Pub/Sub messaging** | Built-in publish/subscribe | Pub/Sub |
| **Event streaming** | Persistent, consumer groups | Stream |
| **Distributed lock** | Atomic operations, TTL | String (SET NX) |
| **Real-time analytics** | HyperLogLog for cardinality | HyperLogLog |
| **Geospatial** | Location-based queries | Geospatial |
| **Queue** | FIFO with blocking pops | List |

---

## 📊 Data Types & Commands

### String

The simplest type — can hold text, numbers, or binary data (max 512 MB).

```bash
# Set and get
SET user:1:name "Jane Doe"
GET user:1:name                    # "Jane Doe"

# Set with expiry
SET session:abc "user-data" EX 3600          # Expires in 1 hour
SET session:abc "user-data" PX 60000         # Expires in 60 seconds (ms)
SETEX session:abc 3600 "user-data"           # Same as SET ... EX

# Set only if not exists (distributed lock)
SET lock:resource "owner-1" NX EX 30         # Returns OK or nil

# Atomic increment/decrement
SET counter 10
INCR counter                       # 11
INCRBY counter 5                   # 16
DECR counter                       # 15
INCRBYFLOAT counter 1.5            # 16.5

# Multiple keys
MSET key1 "val1" key2 "val2" key3 "val3"
MGET key1 key2 key3                # ["val1", "val2", "val3"]

# String operations
APPEND key1 " appended"
STRLEN key1                        # Length of value
GETRANGE key1 0 3                  # Substring
```

### Hash

A map of field-value pairs — perfect for objects.

```bash
# Set fields
HSET user:1 name "Jane" email "jane@example.com" age 30
HSETNX user:1 name "Won't overwrite"  # Set only if field doesn't exist

# Get fields
HGET user:1 name                   # "Jane"
HMGET user:1 name email            # ["Jane", "jane@example.com"]
HGETALL user:1                     # {name: "Jane", email: "jane@example.com", age: "30"}

# Field operations
HINCRBY user:1 age 1               # 31
HDEL user:1 email
HEXISTS user:1 name                # 1 (true)
HKEYS user:1                       # ["name", "age"]
HVALS user:1                       # ["Jane", "31"]
HLEN user:1                        # 2
```

### List

Ordered sequence of strings — implemented as a doubly linked list.

```bash
# Push (add elements)
LPUSH queue "first"                # Push to head (left)
RPUSH queue "last"                 # Push to tail (right)
LPUSH queue "newer"                # ["newer", "first", "last"]

# Pop (remove and return)
LPOP queue                         # "newer" (from head)
RPOP queue                         # "last" (from tail)
BLPOP queue 30                     # Blocking pop (wait up to 30s)

# Access
LRANGE queue 0 -1                  # All elements
LRANGE queue 0 9                   # First 10 elements
LINDEX queue 0                     # Element at index 0
LLEN queue                         # Length

# Trim (keep only range)
LTRIM queue 0 99                   # Keep first 100 elements
```

### Set

Unordered collection of unique strings.

```bash
# Add members
SADD tags:post:1 "python" "redis" "backend"

# Check membership
SISMEMBER tags:post:1 "python"     # 1 (true)

# Get members
SMEMBERS tags:post:1               # {"python", "redis", "backend"}
SCARD tags:post:1                  # 3 (count)

# Set operations
SADD tags:post:2 "python" "django" "api"
SINTER tags:post:1 tags:post:2     # {"python"} (intersection)
SUNION tags:post:1 tags:post:2     # {"python", "redis", "backend", "django", "api"}
SDIFF tags:post:1 tags:post:2      # {"redis", "backend"} (in 1 but not 2)

# Random member
SRANDMEMBER tags:post:1 2          # 2 random members
SPOP tags:post:1                   # Remove and return random member
```

### Sorted Set

Like a Set, but each member has a score for ordering.

```bash
# Add members with scores
ZADD leaderboard 100 "alice" 85 "bob" 92 "charlie"

# Get by rank (ascending)
ZRANGE leaderboard 0 -1 WITHSCORES          # All, lowest first
ZREVRANGE leaderboard 0 2 WITHSCORES        # Top 3, highest first

# Get by score range
ZRANGEBYSCORE leaderboard 80 95              # Score between 80-95
ZCOUNT leaderboard 80 95                     # Count in range

# Rank and score
ZRANK leaderboard "alice"                    # 2 (0-indexed, ascending)
ZREVRANK leaderboard "alice"                 # 0 (highest score)
ZSCORE leaderboard "alice"                   # 100

# Update score
ZINCRBY leaderboard 10 "bob"                # bob's score: 95

# Remove
ZREM leaderboard "charlie"
ZREMRANGEBYSCORE leaderboard 0 50            # Remove scores 0-50

# Cardinality
ZCARD leaderboard                            # Total members
```

### Stream

Append-only log with consumer groups — like a lightweight Kafka.

```bash
# Add entries
XADD events * type "click" page "/home" user_id "123"
XADD events * type "purchase" amount "99.99" user_id "456"
# Returns: "1704067200000-0" (timestamp-sequence)

# Read entries
XRANGE events - +                            # All entries
XRANGE events - + COUNT 10                   # First 10
XREVRANGE events + - COUNT 5                 # Last 5

# Read new entries (blocking)
XREAD BLOCK 5000 STREAMS events $            # Wait for new entries

# Stream info
XLEN events                                  # Entry count
XINFO STREAM events                          # Stream details
```

### Other Types

```bash
# Bitmap (bit-level operations)
SETBIT online:2025-01-01 1001 1    # User 1001 was online
GETBIT online:2025-01-01 1001      # 1
BITCOUNT online:2025-01-01         # Count of online users

# HyperLogLog (approximate cardinality)
PFADD unique_visitors:today "user-1" "user-2" "user-3"
PFADD unique_visitors:today "user-1"  # Duplicate, not counted
PFCOUNT unique_visitors:today          # ~3 (approximate)

# Geospatial
GEOADD locations -73.9857 40.7484 "NYC" -0.1278 51.5074 "London"
GEODIST locations "NYC" "London" km    # ~5570 km
GEOSEARCH locations FROMLONLAT -73.9857 40.7484 BYRADIUS 100 km ASC
```

---

## 🔑 Key Management

```bash
# TTL and expiry
EXPIRE key 3600                    # Set TTL (seconds)
PEXPIRE key 60000                  # Set TTL (milliseconds)
EXPIREAT key 1735689600            # Expire at Unix timestamp
TTL key                            # Remaining TTL (-1 = no expiry, -2 = doesn't exist)
PTTL key                           # Remaining TTL in ms
PERSIST key                        # Remove TTL (make permanent)

# Key operations
EXISTS key                         # 1 if exists
TYPE key                           # Data type
RENAME key newkey
DEL key1 key2 key3                 # Delete (blocking)
UNLINK key1 key2                   # Delete (non-blocking, async)

# Key scanning (production-safe iteration)
SCAN 0 MATCH "user:*" COUNT 100   # Iterate keys matching pattern
# Returns: [cursor, [keys...]]
# Continue with returned cursor until cursor = 0

# ⚠️ NEVER use KEYS in production (blocks server!)
# KEYS user:*                      # DON'T DO THIS IN PRODUCTION
```

---

## 📢 Pub/Sub

```bash
# Subscribe to channels
SUBSCRIBE news:tech news:sports

# Subscribe with pattern
PSUBSCRIBE news:*

# Publish a message
PUBLISH news:tech "Redis 8.0 released!"

# Message flow:
# Publisher → PUBLISH news:tech "msg" → Redis → all subscribers to news:tech
```

> **Limitation:** Pub/Sub messages are fire-and-forget. If no subscriber is listening, the message is lost. Use **Streams** for persistent messaging.

---

## 📺 Streams

### Consumer Groups

```bash
# Create consumer group
XGROUP CREATE events mygroup 0    # Start from beginning
XGROUP CREATE events mygroup $    # Start from new messages

# Read as consumer in group
XREADGROUP GROUP mygroup consumer-1 COUNT 10 STREAMS events >
# '>' means only new (unacknowledged) messages

# Acknowledge processed messages
XACK events mygroup "1704067200000-0"

# Check pending (unacknowledged) messages
XPENDING events mygroup - + 10

# Claim abandoned messages (if consumer crashed)
XCLAIM events mygroup consumer-2 3600000 "1704067200000-0"
# Claims messages idle for > 1 hour

# Delete consumer group
XGROUP DESTROY events mygroup
```

---

## 💱 Transactions

```bash
# MULTI/EXEC (atomic execution)
MULTI
SET account:1:balance 900
SET account:2:balance 1100
EXEC
# Both commands execute atomically

# WATCH (optimistic locking)
WATCH account:1:balance
balance = GET account:1:balance
MULTI
SET account:1:balance (balance - 100)
EXEC
# Returns nil if account:1:balance changed between WATCH and EXEC

# DISCARD (abort transaction)
MULTI
SET key "value"
DISCARD                            # Cancel transaction
```

---

## 📜 Lua Scripting

```bash
# Inline script
EVAL "return redis.call('SET', KEYS[1], ARGV[1])" 1 mykey myvalue

# Rate limiter script
EVAL "
local current = redis.call('INCR', KEYS[1])
if current == 1 then
    redis.call('EXPIRE', KEYS[1], ARGV[1])
end
return current
" 1 "ratelimit:user:123" 60
# Atomically: increment counter, set 60s TTL on first increment

# Load script (returns SHA)
SCRIPT LOAD "return redis.call('GET', KEYS[1])"
# Returns: "abc123..."

# Execute by SHA (faster, no script transfer)
EVALSHA "abc123..." 1 mykey
```

---

## 💾 Persistence

### RDB (Snapshots)

```text
Point-in-time snapshots at configured intervals.

redis.conf:
  save 900 1      # Snapshot if 1+ keys changed in 900 seconds
  save 300 10     # Snapshot if 10+ keys changed in 300 seconds
  save 60 10000   # Snapshot if 10000+ keys changed in 60 seconds

Pros: Compact, fast restart, good for backups
Cons: Data loss between snapshots (minutes)
```

### AOF (Append Only File)

```text
Logs every write operation.

redis.conf:
  appendonly yes
  appendfsync everysec     # Fsync every second (recommended)
  # appendfsync always     # Fsync every write (safest, slowest)
  # appendfsync no         # OS decides when to fsync (fastest)

Pros: Minimal data loss (1 second with everysec)
Cons: Larger files, slower restart
```

### Hybrid (RDB + AOF) — Recommended

```text
redis.conf:
  aof-use-rdb-preamble yes    # Default in Redis 7+

On restart: loads RDB for speed, then replays AOF for completeness
Best of both worlds
```

---

## 🔁 Replication

### Master-Replica

```text
Master ──async replication──► Replica 1 (read-only)
                           └──► Replica 2 (read-only)
```

```bash
# On replica
REPLICAOF master.example.com 6379

# Check replication info
INFO replication

# Promote replica to master
REPLICAOF NO ONE
```

### Redis Sentinel (HA)

```text
Sentinel 1 ──monitor──► Master
Sentinel 2 ──monitor──►   ├──► Replica 1
Sentinel 3 ──monitor──►   └──► Replica 2

If Master fails:
  Sentinels vote → promote Replica 1 to Master
  Redirect clients to new Master
```

```bash
# sentinel.conf
sentinel monitor mymaster master.example.com 6379 2   # 2 = quorum
sentinel down-after-milliseconds mymaster 5000
sentinel failover-timeout mymaster 60000
sentinel parallel-syncs mymaster 1
```

---

## 🌐 Cluster Mode

Redis Cluster provides automatic sharding across multiple nodes.

```text
16384 hash slots distributed across nodes:

Node 1: slots 0-5460       (keys hashing to these slots)
Node 2: slots 5461-10922
Node 3: slots 10923-16383

Key → CRC16(key) mod 16384 → slot → node
```

```bash
# Create cluster
redis-cli --cluster create \
  node1:6379 node2:6379 node3:6379 \
  node4:6379 node5:6379 node6:6379 \
  --cluster-replicas 1

# Cluster info
redis-cli -c CLUSTER INFO
redis-cli -c CLUSTER NODES

# Reshard
redis-cli --cluster reshard node1:6379

# Key tags (force related keys to same slot)
SET {user:123}:profile "..."
SET {user:123}:settings "..."
# Both go to same slot because {user:123} is the hash tag
```

---

## ⚡ Performance

### Pipelining

```bash
# Without pipelining: 100 commands = 100 round trips
# With pipelining: 100 commands = 1 round trip

# Example (redis-cli)
echo "SET key1 val1\nSET key2 val2\nGET key1" | redis-cli --pipe
```

### Memory Optimization

```bash
# Check memory usage
INFO memory
MEMORY USAGE key                   # Bytes used by specific key
MEMORY DOCTOR                      # Memory health report

# Redis configuration
maxmemory 4gb
maxmemory-policy allkeys-lru       # Evict least recently used keys

# Eviction policies:
# noeviction      — Return errors when memory is full
# allkeys-lru     — Evict least recently used (any key)
# allkeys-lfu     — Evict least frequently used (any key)
# volatile-lru    — Evict LRU among keys with TTL
# volatile-lfu    — Evict LFU among keys with TTL
# volatile-ttl    — Evict keys with shortest TTL
# allkeys-random  — Random eviction
```

---

## 🧩 Common Patterns

### Caching

```bash
# Cache-aside pattern
GET user:123               # Check cache
# If miss → query DB → SET user:123 "..." EX 3600
```

### Rate Limiting

```bash
# Fixed window rate limiter
EVAL "
local count = redis.call('INCR', KEYS[1])
if count == 1 then redis.call('EXPIRE', KEYS[1], ARGV[1]) end
if count > tonumber(ARGV[2]) then return 0 end
return 1
" 1 "ratelimit:user:123:minute" 60 100
# Allow 100 requests per 60 seconds
```

### Distributed Lock (Redlock Pattern)

```bash
# Acquire lock
SET lock:resource "owner-uuid" NX EX 30

# Release lock (only if you own it — use Lua)
EVAL "
if redis.call('GET', KEYS[1]) == ARGV[1] then
    return redis.call('DEL', KEYS[1])
end
return 0
" 1 "lock:resource" "owner-uuid"
```

### Session Store

```bash
HSET session:abc user_id "123" role "admin" created_at "2025-01-01"
EXPIRE session:abc 1800   # 30 minute TTL
```

### Leaderboard

```bash
ZADD game:scores 1500 "player-1" 1200 "player-2" 1800 "player-3"
ZREVRANGE game:scores 0 9 WITHSCORES    # Top 10
ZREVRANK game:scores "player-1"          # Player's rank
```

---

## 🐛 Troubleshooting

### High Memory Usage

```bash
INFO memory                        # Check used_memory
MEMORY DOCTOR                      # Recommendations

# Find big keys
redis-cli --bigkeys                # Scan for largest keys
redis-cli --memkeys --samples 100  # Sample key memory

# Solutions:
# - Set maxmemory + eviction policy
# - Compress values before storing
# - Use shorter key names
# - Set TTLs on all keys
```

### Slow Commands

```bash
# Enable slow log
CONFIG SET slowlog-log-slower-than 10000   # Log commands > 10ms
CONFIG SET slowlog-max-len 128

SLOWLOG GET 10                     # View slow commands
SLOWLOG RESET                      # Clear slow log

# Common slow commands to avoid:
# KEYS *          → Use SCAN
# SMEMBERS (large set) → Use SSCAN
# HGETALL (large hash) → Use HSCAN
# SORT (large data) → Pre-sort in application
```

### Connection Issues

```bash
INFO clients                       # Connected clients info
CLIENT LIST                        # All connected clients
CONFIG GET maxclients              # Max connections allowed

# Solutions:
# - Use connection pooling (most Redis libraries support this)
# - Increase maxclients if needed
# - Check for connection leaks in application code
# - Monitor client connections
```

---

## 🏭 Production Tips

### Configuration Checklist

```text
□ maxmemory set with appropriate eviction policy
□ Persistence configured (RDB + AOF hybrid recommended)
□ Replication configured (at least 1 replica)
□ Sentinel or Cluster for high availability
□ AUTH password or ACLs configured
□ TLS enabled for client connections
□ Rename dangerous commands (FLUSHDB, FLUSHALL, DEBUG)
□ maxclients set appropriately
□ Slow log configured for monitoring
□ Key expiry strategy defined
□ Connection pooling in application
□ Monitoring configured (memory, connections, hit rate)
```

### Monitoring Metrics

| Metric | Description | Alert Threshold |
|--------|-------------|-----------------|
| `used_memory` | Total memory used | > 80% of maxmemory |
| `mem_fragmentation_ratio` | Memory fragmentation | > 1.5 or < 1.0 |
| `connected_clients` | Current client count | Unusual spikes |
| `keyspace_hits/misses` | Cache hit ratio | Hit rate < 90% |
| `evicted_keys` | Keys evicted by policy | Any (indicates memory pressure) |
| `rejected_connections` | Refused connections | > 0 |
| `rdb_last_bgsave_status` | Last RDB save status | != "ok" |
| `master_link_status` | Replication link | != "up" |
| `instantaneous_ops_per_sec` | Commands per second | Unusual spikes |

---

## 🔗 Related Topics

- [Distributed Systems — Caching](../../distributed-systems/caching/) — Caching strategies and patterns
- [Distributed Systems — etcd](../../distributed-systems/etcd/) — Comparison for distributed KV stores
- [Databases — PostgreSQL](../postgresql/) — Redis as a caching layer for PostgreSQL
- [Security — Vault](../../security/vault/) — Redis credential management
- [Observability — Prometheus](../../observability/prometheus/) — Monitoring Redis with exporters

---

> **Redis is incredibly powerful but requires discipline.** Set memory limits, use TTLs, choose the right data structure, and monitor actively. Treat Redis as a cache by default — if you need persistence, configure it deliberately and understand the trade-offs.
