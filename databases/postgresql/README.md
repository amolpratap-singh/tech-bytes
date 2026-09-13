# 🐘 PostgreSQL

> **PostgreSQL is the world's most advanced open-source relational database.** It combines SQL compliance, extensibility, and reliability with advanced features like JSONB, window functions, CTEs, and a rich ecosystem of extensions — making it suitable for everything from simple CRUD applications to complex analytical workloads.

---

## 📑 Table of Contents

- [What is PostgreSQL?](#-what-is-postgresql)
- [Data Types](#-data-types)
- [Essential SQL](#-essential-sql)
- [Indexes](#-indexes)
- [EXPLAIN & Query Optimization](#-explain--query-optimization)
- [JSON/JSONB Operations](#-jsonjsonb-operations)
- [Transactions & MVCC](#-transactions--mvcc)
- [Extensions](#-extensions)
- [Administration](#-administration)
- [Replication](#-replication)
- [Performance](#-performance)
- [Troubleshooting](#-troubleshooting)
- [Production Tips](#-production-tips)
- [Related Topics](#-related-topics)

---

## 🧠 What is PostgreSQL?

PostgreSQL (often called "Postgres") has been in development since 1986 (UC Berkeley). It's known for:

- **Standards compliance** — Most SQL-compliant open-source database
- **Extensibility** — Custom types, functions, operators, index types, languages
- **MVCC** — Multi-Version Concurrency Control for concurrent access without locking
- **Reliability** — WAL (Write-Ahead Logging), crash recovery, data integrity
- **JSON support** — First-class JSONB type with indexing and operators
- **Advanced features** — Window functions, CTEs, full-text search, partitioning, foreign data wrappers

---

## 📊 Data Types

### Common Types

| Type | Description | Example |
|------|-------------|---------|
| `integer` / `int` | 32-bit integer | `42` |
| `bigint` | 64-bit integer | `9223372036854775807` |
| `serial` / `bigserial` | Auto-incrementing integer | Auto-generated IDs |
| `numeric(p,s)` | Exact decimal | `99.99` |
| `real` / `double precision` | Floating point | `3.14159` |
| `text` | Variable-length string | `'hello world'` |
| `varchar(n)` | Variable-length with limit | `'hello'` (max n chars) |
| `boolean` | True/false | `true`, `false` |
| `date` | Calendar date | `'2025-01-01'` |
| `timestamp` | Date + time | `'2025-01-01 12:00:00'` |
| `timestamptz` | Timestamp with timezone | `'2025-01-01 12:00:00+00'` |
| `interval` | Time duration | `'1 day 2 hours'` |
| `uuid` | Universally unique identifier | `gen_random_uuid()` |
| `jsonb` | Binary JSON (indexed) | `'{"key": "value"}'` |
| `json` | Text JSON (not indexed) | `'{"key": "value"}'` |
| `array` | Array of any type | `'{1,2,3}'` or `ARRAY[1,2,3]` |
| `bytea` | Binary data | `'\xDEADBEEF'` |
| `inet` / `cidr` | IP address / network | `'192.168.1.1'` |
| `tsvector` | Full-text search document | `to_tsvector('english', 'hello world')` |

> **Tip:** Use `timestamptz` (not `timestamp`) for time data. Use `text` instead of `varchar` unless you need a length constraint. Use `jsonb` (not `json`) for JSON data.

---

## 📝 Essential SQL

### SELECT Queries

```sql
-- Basic select
SELECT id, name, email FROM users WHERE active = true;

-- Pagination
SELECT * FROM users ORDER BY created_at DESC LIMIT 20 OFFSET 40;

-- Aggregation
SELECT department, COUNT(*), AVG(salary)
FROM employees
GROUP BY department
HAVING COUNT(*) > 5
ORDER BY AVG(salary) DESC;

-- Distinct
SELECT DISTINCT department FROM employees;
```

### JOINs

```sql
-- INNER JOIN (only matching rows)
SELECT u.name, o.total
FROM users u
INNER JOIN orders o ON u.id = o.user_id;

-- LEFT JOIN (all left rows, matching right rows or NULL)
SELECT u.name, COUNT(o.id) AS order_count
FROM users u
LEFT JOIN orders o ON u.id = o.user_id
GROUP BY u.id, u.name;

-- Multiple JOINs
SELECT u.name, o.id AS order_id, p.name AS product_name
FROM users u
JOIN orders o ON u.id = o.user_id
JOIN order_items oi ON o.id = oi.order_id
JOIN products p ON oi.product_id = p.id
WHERE o.created_at > NOW() - INTERVAL '30 days';

-- Self JOIN
SELECT e.name AS employee, m.name AS manager
FROM employees e
LEFT JOIN employees m ON e.manager_id = m.id;
```

### Subqueries

```sql
-- Subquery in WHERE
SELECT * FROM users
WHERE id IN (SELECT user_id FROM orders WHERE total > 1000);

-- Correlated subquery
SELECT u.*, (SELECT COUNT(*) FROM orders o WHERE o.user_id = u.id) AS order_count
FROM users u;

-- EXISTS
SELECT * FROM users u
WHERE EXISTS (SELECT 1 FROM orders o WHERE o.user_id = u.id AND o.total > 1000);
```

### CTEs (Common Table Expressions)

```sql
-- Simple CTE
WITH active_users AS (
    SELECT id, name, email
    FROM users
    WHERE last_login > NOW() - INTERVAL '30 days'
)
SELECT au.name, COUNT(o.id) AS recent_orders
FROM active_users au
JOIN orders o ON au.id = o.user_id
GROUP BY au.name;

-- Recursive CTE (hierarchical data)
WITH RECURSIVE org_tree AS (
    -- Base case: top-level managers
    SELECT id, name, manager_id, 0 AS depth
    FROM employees
    WHERE manager_id IS NULL
    UNION ALL
    -- Recursive case
    SELECT e.id, e.name, e.manager_id, ot.depth + 1
    FROM employees e
    JOIN org_tree ot ON e.manager_id = ot.id
)
SELECT * FROM org_tree ORDER BY depth, name;
```

### Window Functions

```sql
-- Row number (ranking)
SELECT name, department, salary,
    ROW_NUMBER() OVER (PARTITION BY department ORDER BY salary DESC) AS rank
FROM employees;

-- Running total
SELECT date, amount,
    SUM(amount) OVER (ORDER BY date) AS running_total
FROM transactions;

-- Moving average
SELECT date, value,
    AVG(value) OVER (ORDER BY date ROWS BETWEEN 6 PRECEDING AND CURRENT ROW) AS moving_avg_7d
FROM metrics;

-- Lag/Lead (previous/next row)
SELECT date, revenue,
    LAG(revenue) OVER (ORDER BY date) AS prev_day_revenue,
    revenue - LAG(revenue) OVER (ORDER BY date) AS daily_change
FROM daily_revenue;

-- Percentile
SELECT name, salary,
    PERCENT_RANK() OVER (ORDER BY salary) AS percentile
FROM employees;

-- NTILE (divide into buckets)
SELECT name, salary,
    NTILE(4) OVER (ORDER BY salary) AS salary_quartile
FROM employees;
```

### INSERT, UPDATE, DELETE

```sql
-- Insert
INSERT INTO users (name, email) VALUES ('Jane', 'jane@example.com');

-- Insert with returning
INSERT INTO users (name, email)
VALUES ('Jane', 'jane@example.com')
RETURNING id, created_at;

-- Upsert (INSERT ... ON CONFLICT)
INSERT INTO users (email, name, login_count)
VALUES ('jane@example.com', 'Jane', 1)
ON CONFLICT (email)
DO UPDATE SET
    login_count = users.login_count + 1,
    last_login = NOW();

-- Update with JOIN
UPDATE orders o
SET status = 'cancelled'
FROM users u
WHERE o.user_id = u.id AND u.active = false;

-- Delete with RETURNING
DELETE FROM sessions WHERE expires_at < NOW() RETURNING id, user_id;
```

---

## 📇 Indexes

### Index Types

| Type | Use Case | Example |
|------|----------|---------|
| **B-tree** (default) | Equality, range, sorting | `CREATE INDEX idx ON t(col)` |
| **Hash** | Equality only | `CREATE INDEX idx ON t USING hash(col)` |
| **GIN** | Arrays, JSONB, full-text search | `CREATE INDEX idx ON t USING gin(jsonb_col)` |
| **GiST** | Geometric, full-text, range types | `CREATE INDEX idx ON t USING gist(geom_col)` |
| **BRIN** | Large, naturally ordered tables | `CREATE INDEX idx ON t USING brin(created_at)` |
| **SP-GiST** | Non-balanced data structures | `CREATE INDEX idx ON t USING spgist(point_col)` |

### Index Strategies

```sql
-- Partial index (only index subset of rows)
CREATE INDEX idx_active_users ON users(email) WHERE active = true;

-- Expression index (index on computed value)
CREATE INDEX idx_lower_email ON users(LOWER(email));

-- Multi-column index (order matters!)
CREATE INDEX idx_user_date ON orders(user_id, created_at DESC);

-- Covering index (include non-key columns)
CREATE INDEX idx_orders_cover ON orders(user_id) INCLUDE (total, status);

-- Unique index
CREATE UNIQUE INDEX idx_unique_email ON users(email);

-- Concurrent index creation (no locks!)
CREATE INDEX CONCURRENTLY idx_name ON users(name);
```

### When to Index

```text
✅ Index when:
  - Column used in WHERE clauses frequently
  - Column used in JOIN conditions
  - Column used in ORDER BY
  - Foreign key columns
  - Unique constraints

❌ Don't index when:
  - Table is small (< 10,000 rows)
  - Column has low cardinality (true/false)
  - Table is write-heavy with few reads
  - Column is rarely used in queries
```

---

## 🔍 EXPLAIN & Query Optimization

### Reading EXPLAIN Output

```sql
EXPLAIN ANALYZE SELECT * FROM users WHERE email = 'jane@example.com';

-- Output:
-- Index Scan using idx_users_email on users  (cost=0.42..8.44 rows=1 width=128) (actual time=0.023..0.024 rows=1 loops=1)
--   Index Cond: (email = 'jane@example.com'::text)
-- Planning Time: 0.082 ms
-- Execution Time: 0.042 ms
```

| Field | Description |
|-------|-------------|
| **cost** | Estimated startup..total cost (arbitrary units) |
| **rows** | Estimated number of rows |
| **actual time** | Real execution time (ms) |
| **loops** | Number of times this node executed |

### Scan Types (Best to Worst)

| Scan Type | Performance | When Used |
|-----------|-------------|-----------|
| **Index Only Scan** | ⭐⭐⭐⭐⭐ | All needed columns are in the index |
| **Index Scan** | ⭐⭐⭐⭐ | Index lookup + table access |
| **Bitmap Index Scan** | ⭐⭐⭐ | Multiple index conditions combined |
| **Seq Scan** | ⭐⭐ (for small tables) | Full table scan |

### Common Optimization Patterns

```sql
-- Bad: Function prevents index use
SELECT * FROM users WHERE LOWER(email) = 'jane@example.com';
-- Fix: Create expression index
CREATE INDEX idx_lower_email ON users(LOWER(email));

-- Bad: Leading wildcard prevents index use
SELECT * FROM users WHERE name LIKE '%jane%';
-- Fix: Use pg_trgm extension for trigram index
CREATE EXTENSION pg_trgm;
CREATE INDEX idx_name_trgm ON users USING gin(name gin_trgm_ops);

-- Bad: Implicit type cast
SELECT * FROM users WHERE id = '123';  -- id is integer, '123' is text
-- Fix: Use correct type
SELECT * FROM users WHERE id = 123;

-- Bad: SELECT * when only need few columns
SELECT * FROM large_table WHERE ...;
-- Fix: Select only needed columns
SELECT id, name FROM large_table WHERE ...;

-- Bad: N+1 queries
-- Fix: Use JOINs or batch queries
SELECT u.*, array_agg(o.id) AS order_ids
FROM users u
LEFT JOIN orders o ON u.id = o.user_id
GROUP BY u.id;
```

---

## 📋 JSON/JSONB Operations

```sql
-- Create table with JSONB column
CREATE TABLE events (
    id bigserial PRIMARY KEY,
    data jsonb NOT NULL,
    created_at timestamptz DEFAULT NOW()
);

-- Insert JSON data
INSERT INTO events (data) VALUES ('{"type": "click", "page": "/home", "user_id": 123}');

-- Access JSON fields
SELECT data->>'type' AS event_type,          -- text value
       data->'user_id' AS user_id_json,      -- JSON value
       (data->>'user_id')::int AS user_id    -- cast to int
FROM events;

-- Filter by JSON field
SELECT * FROM events WHERE data->>'type' = 'click';
SELECT * FROM events WHERE (data->>'user_id')::int = 123;

-- Containment operator (@>)
SELECT * FROM events WHERE data @> '{"type": "click"}';

-- Existence operator (?)
SELECT * FROM events WHERE data ? 'user_id';

-- Path access
SELECT data #>> '{metadata,browser}' FROM events;

-- Update JSON field
UPDATE events SET data = jsonb_set(data, '{processed}', 'true') WHERE id = 1;

-- Remove JSON key
UPDATE events SET data = data - 'temporary_field' WHERE id = 1;

-- JSONB indexing (GIN)
CREATE INDEX idx_events_data ON events USING gin(data);
-- Supports: @>, ?, ?|, ?&

-- JSONB path index (for specific fields)
CREATE INDEX idx_events_type ON events ((data->>'type'));
```

---

## 🔄 Transactions & MVCC

### Isolation Levels

| Level | Dirty Read | Non-Repeatable Read | Phantom Read | Performance |
|-------|-----------|-------------------|--------------|-------------|
| **Read Uncommitted** | Possible* | Possible | Possible | Fastest |
| **Read Committed** (default) | ❌ No | Possible | Possible | Fast |
| **Repeatable Read** | ❌ No | ❌ No | ❌ No** | Moderate |
| **Serializable** | ❌ No | ❌ No | ❌ No | Slowest |

\* PostgreSQL treats Read Uncommitted as Read Committed
\** PostgreSQL's Repeatable Read also prevents phantom reads (uses snapshot isolation)

```sql
-- Set isolation level for transaction
BEGIN ISOLATION LEVEL SERIALIZABLE;
  SELECT balance FROM accounts WHERE id = 1;
  UPDATE accounts SET balance = balance - 100 WHERE id = 1;
  UPDATE accounts SET balance = balance + 100 WHERE id = 2;
COMMIT;

-- Advisory locks (application-level)
SELECT pg_advisory_lock(12345);     -- Block until acquired
-- ... critical section ...
SELECT pg_advisory_unlock(12345);

-- Try lock (non-blocking)
SELECT pg_try_advisory_lock(12345); -- Returns true/false
```

### Deadlock Detection

```sql
-- Check for locks
SELECT pid, pg_blocking_pids(pid) AS blocked_by, query
FROM pg_stat_activity
WHERE cardinality(pg_blocking_pids(pid)) > 0;

-- Kill a blocking query
SELECT pg_cancel_backend(<pid>);     -- Graceful cancel
SELECT pg_terminate_backend(<pid>);  -- Force terminate
```

---

## 🧩 Extensions

```sql
-- List installed extensions
SELECT * FROM pg_available_extensions WHERE installed_version IS NOT NULL;

-- pg_stat_statements (query performance)
CREATE EXTENSION pg_stat_statements;
SELECT query, calls, mean_exec_time, rows
FROM pg_stat_statements
ORDER BY mean_exec_time DESC
LIMIT 20;

-- pgcrypto (encryption functions)
CREATE EXTENSION pgcrypto;
SELECT crypt('password', gen_salt('bf'));  -- bcrypt hash
SELECT gen_random_uuid();                  -- UUID generation

-- pg_trgm (fuzzy text matching)
CREATE EXTENSION pg_trgm;
SELECT * FROM users WHERE name % 'Jon';   -- similarity search
SELECT similarity('Jon', 'John');         -- returns 0.0-1.0

-- PostGIS (geospatial)
CREATE EXTENSION postgis;
SELECT ST_Distance(
    ST_MakePoint(-73.9857, 40.7484)::geography,
    ST_MakePoint(-0.1278, 51.5074)::geography
) / 1000 AS distance_km;
```

---

## 🛠 Administration

### psql Essential Commands

```bash
# Connect
psql -h localhost -U myuser -d mydb

# Inside psql:
\l              # List databases
\c mydb         # Connect to database
\dt             # List tables
\dt+            # List tables with sizes
\d tablename    # Describe table
\di             # List indexes
\du             # List users/roles
\df             # List functions
\x              # Toggle expanded display
\timing         # Toggle query timing
\q              # Quit
```

### Backup & Restore

```bash
# Dump a database
pg_dump -h localhost -U myuser mydb > backup.sql
pg_dump -h localhost -U myuser -Fc mydb > backup.dump  # Custom format (compressed)

# Dump schema only
pg_dump -h localhost -U myuser --schema-only mydb > schema.sql

# Dump specific tables
pg_dump -h localhost -U myuser -t users -t orders mydb > tables.sql

# Restore from SQL dump
psql -h localhost -U myuser mydb < backup.sql

# Restore from custom format
pg_restore -h localhost -U myuser -d mydb backup.dump

# Dump all databases
pg_dumpall -h localhost -U myuser > all_databases.sql
```

### Vacuuming

```sql
-- Manual vacuum (reclaim space, update statistics)
VACUUM ANALYZE users;

-- Full vacuum (compacts table, requires exclusive lock)
VACUUM FULL users;

-- Check autovacuum activity
SELECT relname, last_vacuum, last_autovacuum, n_dead_tup, n_live_tup
FROM pg_stat_user_tables
ORDER BY n_dead_tup DESC;

-- Autovacuum configuration (postgresql.conf)
-- autovacuum = on
-- autovacuum_vacuum_threshold = 50
-- autovacuum_vacuum_scale_factor = 0.2
-- autovacuum_analyze_threshold = 50
-- autovacuum_analyze_scale_factor = 0.1
```

---

## 🔁 Replication

### Streaming Replication

```text
Primary ──WAL stream──► Standby (read replica)
                         ├── Synchronous (waits for standby confirmation)
                         └── Asynchronous (default, doesn't wait)
```

```bash
# On primary (postgresql.conf)
# wal_level = replica
# max_wal_senders = 5
# synchronous_standby_names = ''

# On standby
pg_basebackup -h primary.example.com -U replication_user -D /var/lib/postgresql/data -P -R
# -R creates standby.signal and configures primary_conninfo
```

### Logical Replication

```sql
-- On publisher (source)
CREATE PUBLICATION my_pub FOR TABLE users, orders;

-- On subscriber (target)
CREATE SUBSCRIPTION my_sub
CONNECTION 'host=primary.example.com dbname=mydb user=repl_user'
PUBLICATION my_pub;

-- Check replication status
SELECT * FROM pg_stat_replication;
SELECT * FROM pg_stat_subscription;
```

---

## ⚡ Performance

### Connection Pooling (PgBouncer)

```ini
# pgbouncer.ini
[databases]
mydb = host=localhost port=5432 dbname=mydb

[pgbouncer]
listen_port = 6432
listen_addr = *
auth_type = md5
pool_mode = transaction    # Best for most apps
max_client_conn = 1000
default_pool_size = 25
```

### Table Partitioning

```sql
-- Range partitioning (by date)
CREATE TABLE events (
    id bigserial,
    data jsonb,
    created_at timestamptz NOT NULL
) PARTITION BY RANGE (created_at);

CREATE TABLE events_2025_01 PARTITION OF events
    FOR VALUES FROM ('2025-01-01') TO ('2025-02-01');
CREATE TABLE events_2025_02 PARTITION OF events
    FOR VALUES FROM ('2025-02-01') TO ('2025-03-01');

-- Hash partitioning
CREATE TABLE users (
    id bigserial,
    name text,
    email text
) PARTITION BY HASH (id);

CREATE TABLE users_0 PARTITION OF users FOR VALUES WITH (MODULUS 4, REMAINDER 0);
CREATE TABLE users_1 PARTITION OF users FOR VALUES WITH (MODULUS 4, REMAINDER 1);
CREATE TABLE users_2 PARTITION OF users FOR VALUES WITH (MODULUS 4, REMAINDER 2);
CREATE TABLE users_3 PARTITION OF users FOR VALUES WITH (MODULUS 4, REMAINDER 3);
```

---

## 🐛 Troubleshooting

### Slow Queries

```sql
-- Find slow queries (requires pg_stat_statements)
SELECT query, calls, mean_exec_time::numeric(10,2) AS avg_ms,
       total_exec_time::numeric(10,2) AS total_ms, rows
FROM pg_stat_statements
ORDER BY mean_exec_time DESC LIMIT 10;

-- Currently running queries
SELECT pid, now() - pg_stat_activity.query_start AS duration, query, state
FROM pg_stat_activity
WHERE state != 'idle'
ORDER BY duration DESC;

-- Kill long-running query
SELECT pg_cancel_backend(<pid>);
```

### Connection Issues

```sql
-- Check current connections
SELECT count(*) FROM pg_stat_activity;
SELECT datname, usename, state, count(*)
FROM pg_stat_activity GROUP BY datname, usename, state;

-- Check max connections setting
SHOW max_connections;

-- Check connection limits
SELECT rolname, rolconnlimit FROM pg_roles WHERE rolconnlimit > 0;
```

### Lock Contention

```sql
-- View active locks
SELECT l.pid, l.mode, l.granted, a.query
FROM pg_locks l
JOIN pg_stat_activity a ON l.pid = a.pid
WHERE NOT l.granted;

-- Find blocking queries
SELECT blocked.pid AS blocked_pid,
       blocked.query AS blocked_query,
       blocking.pid AS blocking_pid,
       blocking.query AS blocking_query
FROM pg_stat_activity blocked
JOIN pg_locks bl ON bl.pid = blocked.pid
JOIN pg_locks kl ON kl.locktype = bl.locktype
  AND kl.database IS NOT DISTINCT FROM bl.database
  AND kl.relation IS NOT DISTINCT FROM bl.relation
  AND kl.pid != bl.pid
JOIN pg_stat_activity blocking ON kl.pid = blocking.pid
WHERE NOT bl.granted;
```

### Table Bloat

```sql
-- Check table bloat
SELECT schemaname, tablename,
    pg_size_pretty(pg_total_relation_size(schemaname || '.' || tablename)) AS total_size,
    n_dead_tup, n_live_tup,
    ROUND(n_dead_tup * 100.0 / NULLIF(n_live_tup + n_dead_tup, 0), 2) AS dead_pct
FROM pg_stat_user_tables
ORDER BY n_dead_tup DESC LIMIT 10;
```

---

## 🏭 Production Tips

### Configuration Recommendations

| Parameter | Default | Recommended | Description |
|-----------|---------|-------------|-------------|
| `shared_buffers` | 128MB | 25% of RAM | Shared memory for caching |
| `effective_cache_size` | 4GB | 75% of RAM | Planner's estimate of OS cache |
| `work_mem` | 4MB | 16-64MB | Per-operation sort/hash memory |
| `maintenance_work_mem` | 64MB | 512MB-1GB | Vacuum, index creation memory |
| `wal_buffers` | -1 (auto) | 64MB | WAL write buffering |
| `max_connections` | 100 | 100-300 | Use PgBouncer for more |
| `random_page_cost` | 4.0 | 1.1 (SSD) | Cost estimate for random I/O |

### Monitoring Checklist

```text
□ Active connections vs max_connections
□ Query execution time (p95, p99)
□ Cache hit ratio (should be > 99%)
□ Dead tuples and autovacuum activity
□ Replication lag (if using replicas)
□ Disk usage and growth rate
□ Lock contention
□ Connection pool utilization (PgBouncer)
□ WAL generation rate
□ Table and index bloat
```

```sql
-- Cache hit ratio (should be > 99%)
SELECT sum(heap_blks_hit) / (sum(heap_blks_hit) + sum(heap_blks_read)) AS cache_hit_ratio
FROM pg_statio_user_tables;
```

---

## 🔗 Related Topics

- [Databases — Redis](../redis/) — Caching layer for PostgreSQL
- [Security — Vault](../../security/vault/) — Dynamic PostgreSQL credentials
- [Search — Elasticsearch](../../search/elasticsearch/) — Full-text search complementing PostgreSQL
- [Distributed Systems — Caching](../../distributed-systems/caching/) — Caching strategies for databases
- [DevOps — Docker](../../devops/docker/) — Running PostgreSQL in containers

---

> **PostgreSQL is the Swiss Army knife of databases.** Its combination of reliability, features, and extensibility makes it the default choice for most applications. Start simple and leverage advanced features as your needs grow.
