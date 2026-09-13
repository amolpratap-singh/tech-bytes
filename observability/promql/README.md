# 📐 PromQL — Prometheus Query Language

> **The query language for Prometheus and compatible systems.** PromQL is a functional, expressive language designed for selecting, aggregating, and transforming time-series data. It's the core skill for building dashboards, alerting rules, and debugging metrics-based issues.

---

## 📑 Table of Contents

- [What is PromQL?](#-what-is-promql)
- [Data Types](#-data-types)
- [Selectors](#-selectors)
- [Operators](#-operators)
- [Vector Matching](#-vector-matching)
- [Aggregation Operators](#-aggregation-operators)
- [Functions: Rate](#-functions-rate)
- [Functions: Aggregation over Time](#-functions-aggregation-over-time)
- [Functions: Math](#-functions-math)
- [Functions: Time and Date](#-functions-time-and-date)
- [Functions: Other](#-functions-other)
- [Histogram and Summary](#-histogram-and-summary)
- [Subqueries](#-subqueries)
- [Recording Rules](#-recording-rules)
- [Real-World Query Examples](#-real-world-query-examples)
- [Performance and Optimization](#-performance-and-optimization)
- [Common Mistakes](#-common-mistakes)
- [Troubleshooting](#-troubleshooting)
- [Related Topics](#-related-topics)

---

## 🧠 What is PromQL?

PromQL (Prometheus Query Language) is a functional language that lets you select and aggregate time-series data in real time. It's used in:

- **Prometheus UI** — Ad-hoc queries and exploration
- **Grafana dashboards** — Panel queries and variables
- **Alert rules** — Conditions that trigger alerts
- **Recording rules** — Pre-computed queries stored as new time series
- **HTTP API** — Programmatic queries via `/api/v1/query` and `/api/v1/query_range`

### Query Endpoints

```bash
# Instant query — evaluate at a single point in time
curl 'http://localhost:9090/api/v1/query?query=up&time=2024-01-15T10:00:00Z'

# Range query — evaluate over a time range
curl 'http://localhost:9090/api/v1/query_range?query=rate(http_requests_total[5m])&start=2024-01-15T00:00:00Z&end=2024-01-15T12:00:00Z&step=60s'
```

---

## 📊 Data Types

PromQL has four data types:

### 1. Instant Vector

A set of time series, each with a **single sample** at the current evaluation time.

```promql
# Returns one sample per matching time series
http_requests_total{method="GET"}
```

```text
http_requests_total{method="GET", handler="/api/users", status="200"} 12345
http_requests_total{method="GET", handler="/api/users", status="404"} 67
http_requests_total{method="GET", handler="/api/health", status="200"} 89012
```

### 2. Range Vector

A set of time series, each with a **range of samples** over a time window. Used as input to functions like `rate()`.

```promql
# Returns samples from the last 5 minutes for each matching series
http_requests_total{method="GET"}[5m]
```

```text
http_requests_total{method="GET", handler="/api/users", status="200"}
  12300 @1705312200
  12320 @1705312215
  12345 @1705312230
  ...
```

> **Important:** Range vectors cannot be graphed directly. They must be passed through a function (like `rate()`) that reduces them to an instant vector.

### 3. Scalar

A simple numeric floating-point value:

```promql
42
3.14
-1.5e3
```

### 4. String

A simple string value (rarely used directly):

```promql
"hello world"
```

### Time Durations

| Unit | Meaning |
|------|---------|
| `ms` | Milliseconds |
| `s` | Seconds |
| `m` | Minutes |
| `h` | Hours |
| `d` | Days (24h) |
| `w` | Weeks (7d) |
| `y` | Years (365d) |

```promql
rate(http_requests_total[5m])     # 5 minute range
rate(http_requests_total[1h30m])  # 1 hour 30 minutes
avg_over_time(up[7d])             # 7 days
```

---

## 🎯 Selectors

### Metric Name Selector

```promql
# Select by metric name
http_requests_total

# Equivalent using __name__ label
{__name__="http_requests_total"}

# Regex on metric name
{__name__=~"http_requests_.*"}
```

### Label Matchers

| Operator | Description | Example |
|----------|-------------|---------|
| `=` | Exact match | `{method="GET"}` |
| `!=` | Not equal | `{status!="200"}` |
| `=~` | Regex match | `{status=~"5.."}` |
| `!~` | Regex not match | `{handler!~"/health\|/ready"}` |

```promql
# Exact match
http_requests_total{method="GET", status="200"}

# Not equal
http_requests_total{status!="200"}

# Regex match (5xx errors)
http_requests_total{status=~"5.."}

# Regex NOT match (exclude health endpoints)
http_requests_total{handler!~"/health|/ready|/metrics"}

# Combine matchers
http_requests_total{method="GET", status=~"[45].."}

# All series with a specific label (any value)
http_requests_total{handler!=""}

# Match against __name__ for multi-metric queries
{__name__=~"http_(requests|errors)_total"}
```

### Offset Modifier

Query data from the past:

```promql
# Value 1 hour ago
http_requests_total offset 1h

# Rate 1 day ago
rate(http_requests_total[5m] offset 1d)

# Compare current to previous day
rate(http_requests_total[5m]) / rate(http_requests_total[5m] offset 1d)
```

### @ Modifier

Evaluate at a specific timestamp:

```promql
# Value at specific Unix timestamp
http_requests_total @ 1705312200

# Value at start of query range
http_requests_total @ start()

# Value at end of query range
http_requests_total @ end()
```

---

## ⚙️ Operators

### Arithmetic Operators

| Operator | Description | Example |
|----------|-------------|---------|
| `+` | Addition | `node_memory_MemTotal_bytes - node_memory_MemAvailable_bytes` |
| `-` | Subtraction | `node_filesystem_size_bytes - node_filesystem_avail_bytes` |
| `*` | Multiplication | `rate(http_requests_total[5m]) * 100` |
| `/` | Division | `node_memory_MemAvailable_bytes / node_memory_MemTotal_bytes` |
| `%` | Modulo | `time() % 3600` |
| `^` | Power | `2 ^ 10` |

```promql
# Memory usage percentage
(1 - node_memory_MemAvailable_bytes / node_memory_MemTotal_bytes) * 100

# Disk usage percentage
(1 - node_filesystem_avail_bytes / node_filesystem_size_bytes) * 100

# Bytes to megabytes
process_resident_memory_bytes / 1024 / 1024
```

### Comparison Operators

| Operator | Description |
|----------|-------------|
| `==` | Equal |
| `!=` | Not equal |
| `>` | Greater than |
| `<` | Less than |
| `>=` | Greater than or equal |
| `<=` | Less than or equal |

```promql
# Filter: only series where memory > 90%
(1 - node_memory_MemAvailable_bytes / node_memory_MemTotal_bytes) > 0.9

# Bool modifier: return 0/1 instead of filtering
(1 - node_memory_MemAvailable_bytes / node_memory_MemTotal_bytes) > bool 0.9
```

> **Without `bool`:** Comparison operators filter (drop non-matching series).
> **With `bool`:** Return 1 (true) or 0 (false) for each series.

### Logical/Set Operators

| Operator | Description |
|----------|-------------|
| `and` | Intersection — returns LHS elements that have matching label sets in RHS |
| `or` | Union — returns all elements from both sides |
| `unless` | Complement — returns LHS elements that do NOT have matching label sets in RHS |

```promql
# Pods that are running AND have high memory
kube_pod_status_phase{phase="Running"} and on(pod, namespace) (
  container_memory_working_set_bytes > 1e9
)

# All targets, but mark failed ones
up == 1 or up == 0

# Pods that are running but NOT ready
kube_pod_status_phase{phase="Running"} unless on(pod, namespace) (
  kube_pod_status_ready{condition="true"}
)
```

### Operator Precedence (highest to lowest)

1. `^`
2. `*`, `/`, `%`, `atan2`
3. `+`, `-`
4. `==`, `!=`, `<=`, `<`, `>=`, `>`
5. `and`, `unless`
6. `or`

---

## 🔗 Vector Matching

When binary operators are applied between two instant vectors, PromQL needs to determine which elements match.

### One-to-One Matching

Default behavior — each element on one side matches exactly one element on the other.

```promql
# Match by all labels
method_code:http_errors:rate5m / method_code:http_requests:rate5m

# Match using specific labels only
method_code:http_errors:rate5m / ignoring(code) method:http_requests:rate5m

# Match on specific labels
method_code:http_errors:rate5m / on(method) method:http_requests:rate5m
```

### Many-to-One and One-to-Many

When one side has multiple matches for a single element on the other.

```promql
# Many errors (different codes) → one total request count
# group_left: "many" is on the left side
method_code:http_errors:rate5m / ignoring(code) group_left method:http_requests:rate5m

# group_right: "many" is on the right side
method:http_requests:rate5m / ignoring(code) group_right method_code:http_errors:rate5m
```

**Visual explanation:**

```text
Left side (many):                       Right side (one):
{method="GET", code="200"} 100          {method="GET"} 1000
{method="GET", code="404"} 50
{method="GET", code="500"} 10

With: left / ignoring(code) group_left right

Result:
{method="GET", code="200"} 0.1     (100/1000)
{method="GET", code="404"} 0.05    (50/1000)
{method="GET", code="500"} 0.01    (10/1000)
```

### Keywords Summary

| Keyword | Purpose |
|---------|---------|
| `on(labels)` | Match only on specified labels |
| `ignoring(labels)` | Match on all labels except specified |
| `group_left(labels)` | Allow many-to-one, "many" on left, optionally copy labels from right |
| `group_right(labels)` | Allow one-to-many, "many" on right, optionally copy labels from left |

---

## 📊 Aggregation Operators

Aggregation operators take an instant vector and aggregate across dimensions (labels).

### Syntax

```promql
<aggr_op>([parameter,] <vector>) [by|without (<label list>)]
```

- `by(labels)` — Keep only specified labels (aggregate everything else)
- `without(labels)` — Remove specified labels (keep everything else)

### Operators

| Operator | Description | Example |
|----------|-------------|---------|
| `sum` | Sum of values | `sum(rate(http_requests_total[5m]))` |
| `min` | Minimum value | `min(node_filesystem_avail_bytes)` |
| `max` | Maximum value | `max(container_memory_working_set_bytes) by (pod)` |
| `avg` | Average value | `avg(rate(http_request_duration_seconds_sum[5m]))` |
| `count` | Count of elements | `count(up == 1)` |
| `count_values` | Count elements with same value | `count_values("version", node_exporter_build_info)` |
| `topk` | Top K elements | `topk(10, rate(http_requests_total[5m]))` |
| `bottomk` | Bottom K elements | `bottomk(5, up)` |
| `quantile` | Quantile over dimensions | `quantile(0.95, rate(http_request_duration_seconds_sum[5m]))` |
| `stddev` | Standard deviation | `stddev(rate(http_requests_total[5m])) by (job)` |
| `stdvar` | Standard variance | `stdvar(rate(http_requests_total[5m]))` |
| `group` | Returns 1 for each group | `group(kube_pod_info) by (namespace)` |

### Examples

```promql
# Total request rate across all instances
sum(rate(http_requests_total[5m]))

# Request rate per handler
sum by (handler) (rate(http_requests_total[5m]))

# Request rate per job, excluding status label
sum without (status, instance) (rate(http_requests_total[5m]))

# Top 10 pods by CPU usage
topk(10, sum by (pod) (rate(container_cpu_usage_seconds_total[5m])))

# Count of running pods per namespace
count by (namespace) (kube_pod_status_phase{phase="Running"})

# 95th percentile of request rate across instances
quantile(0.95, rate(http_requests_total[5m]))

# Count how many instances run each version
count_values("go_version", go_info)
```

---

## 📈 Functions: Rate

Rate functions are essential for working with counters. They calculate the per-second rate of increase.

### rate()

Average per-second rate of increase over a range. **The most commonly used function.**

```promql
# Requests per second over the last 5 minutes
rate(http_requests_total[5m])
```

**How it works:**
1. Takes the first and last values in the range
2. Calculates `(last - first) / time_elapsed`
3. Handles counter resets automatically

**When to use:** Dashboard visualizations, recording rules, alerts. Use with counters.

> **Rule of thumb:** Range should be at least 4× the scrape interval. For 15s scrape interval, use `[1m]` minimum. For dashboards, `[5m]` is a good default.

### irate()

Instantaneous per-second rate — uses only the **last two data points** in the range.

```promql
# Instantaneous request rate
irate(http_requests_total[5m])
```

**When to use:** When you need to see spikes and rapid changes. More volatile than `rate()`.

> **Caution:** `irate()` is sensitive to scrape interval and can produce misleading results if samples are missed. Not suitable for alerting rules.

### rate() vs irate()

| Feature | rate() | irate() |
|---------|--------|---------|
| Calculation | Average over entire range | Last two points only |
| Smoothing | Smooth, averaged | Volatile, shows spikes |
| Missing scrapes | Resilient | Sensitive |
| Alerting | **Recommended** | Not recommended |
| Dashboards | General use | Spike detection |
| Range purpose | Determines averaging window | Only needs 2 samples in range |

### increase()

Total increase over a range. Equivalent to `rate() * range_seconds`.

```promql
# Total requests in the last hour
increase(http_requests_total[1h])

# Equivalent to:
rate(http_requests_total[1h]) * 3600
```

**When to use:** When you want absolute counts rather than per-second rates. Good for "how many errors in the last hour?" type questions.

> **Note:** `increase()` may return non-integer values because it extrapolates.

### delta()

Change in value over a range. For **gauges** (not counters).

```promql
# Temperature change over the last hour
delta(temperature_celsius[1h])

# Memory change in the last 30 minutes
delta(process_resident_memory_bytes[30m])
```

### idelta()

Instantaneous delta — uses the last two data points.

```promql
# Instantaneous change in memory
idelta(process_resident_memory_bytes[5m])
```

### deriv()

Per-second derivative using simple linear regression over the range.

```promql
# Rate of memory growth (bytes per second)
deriv(process_resident_memory_bytes[1h])

# Predict disk full (extrapolate current trend)
node_filesystem_avail_bytes + deriv(node_filesystem_avail_bytes[6h]) * 3600 * 24
```

**When to use:** Trend analysis, prediction. Works on gauges. More stable than `delta()` for trend detection.

---

## ⏱ Functions: Aggregation over Time

`_over_time` functions aggregate samples within a range vector for each time series independently.

| Function | Description |
|----------|-------------|
| `avg_over_time(v[range])` | Average value |
| `min_over_time(v[range])` | Minimum value |
| `max_over_time(v[range])` | Maximum value |
| `sum_over_time(v[range])` | Sum of values |
| `count_over_time(v[range])` | Count of samples |
| `quantile_over_time(φ, v[range])` | φ-quantile (0 ≤ φ ≤ 1) |
| `stddev_over_time(v[range])` | Standard deviation |
| `stdvar_over_time(v[range])` | Standard variance |
| `last_over_time(v[range])` | Most recent value |
| `present_over_time(v[range])` | Value 1 if any sample exists |
| `changes(v[range])` | Number of times value changed |
| `resets(v[range])` | Number of counter resets |

```promql
# Average CPU usage over the last hour
avg_over_time(instance:node_cpu_utilisation:ratio_rate5m[1h])

# Maximum memory usage in the last 24 hours
max_over_time(container_memory_working_set_bytes[24h])

# 95th percentile of latency over the last hour
quantile_over_time(0.95, http_request_duration_seconds[1h])

# How many times did the pod restart in the last 24h?
resets(kube_pod_container_status_restarts_total[24h])

# Is the target present at all?
present_over_time(up{job="my-app"}[5m])

# How many times did the config change in the last hour?
changes(config_hash[1h])
```

---

## 🔢 Functions: Math

| Function | Description | Example |
|----------|-------------|---------|
| `abs(v)` | Absolute value | `abs(delta(temperature[1h]))` |
| `ceil(v)` | Round up | `ceil(rate(http_requests_total[5m]))` |
| `floor(v)` | Round down | `floor(rate(http_requests_total[5m]))` |
| `round(v [, to_nearest])` | Round to nearest | `round(rate(http_requests_total[5m]), 0.1)` |
| `clamp(v, min, max)` | Clamp between min and max | `clamp(cpu_usage, 0, 1)` |
| `clamp_min(v, min)` | Set lower bound | `clamp_min(available_slots, 0)` |
| `clamp_max(v, max)` | Set upper bound | `clamp_max(progress, 100)` |
| `ln(v)` | Natural logarithm | `ln(http_requests_total)` |
| `log2(v)` | Base-2 logarithm | `log2(heap_size_bytes)` |
| `log10(v)` | Base-10 logarithm | `log10(http_requests_total)` |
| `exp(v)` | Exponential (e^v) | `exp(predicted_growth)` |
| `sqrt(v)` | Square root | `sqrt(stdvar_over_time(latency[1h]))` |
| `sgn(v)` | Sign (-1, 0, 1) | `sgn(delta(temperature[1h]))` |

```promql
# Round memory to GB
round(node_memory_MemTotal_bytes / 1024 / 1024 / 1024, 0.1)

# Ensure metric never goes below 0
clamp_min(predicted_available_disk_bytes, 0)

# Standard deviation (sqrt of variance)
sqrt(stdvar_over_time(http_request_duration_seconds[1h]))
```

---

## 🕐 Functions: Time and Date

| Function | Description | Return Value |
|----------|-------------|-------------|
| `time()` | Current Unix timestamp | `1705312200` |
| `timestamp(v)` | Timestamp of each sample | Unix timestamp per series |
| `day_of_week(v)` | Day of the week | 0 (Sunday) - 6 (Saturday) |
| `day_of_month(v)` | Day of the month | 1-31 |
| `day_of_year(v)` | Day of the year | 1-366 |
| `days_in_month(v)` | Days in the month | 28-31 |
| `month(v)` | Month of the year | 1-12 |
| `year(v)` | Year | e.g. 2024 |
| `hour(v)` | Hour of the day | 0-23 |
| `minute(v)` | Minute of the hour | 0-59 |

```promql
# How long since a metric was last updated (staleness detection)
time() - timestamp(my_metric)

# Age of last successful backup in hours
(time() - last_backup_timestamp_seconds) / 3600

# Only alert during business hours (Mon-Fri, 9-17)
ALERTS{alertname="NonCritical"}
  and on() (hour() >= 9 < 17)
  and on() (day_of_week() >= 1 <= 5)

# Certificate expiry in days
(tls_certificate_not_after - time()) / 86400

# Uptime in days
(time() - process_start_time_seconds) / 86400
```

---

## 🛠 Functions: Other

### absent() and absent_over_time()

Returns 1 if the selector matches no time series. Essential for detecting missing metrics.

```promql
# Alert if no data from a target
absent(up{job="my-app"})

# Alert if no data in the last 15 minutes
absent_over_time(up{job="my-app"}[15m])

# Alert if a specific metric has disappeared
absent(http_requests_total{job="api-server"})
```

> **Use case:** Detect when a target stops reporting, when a metric disappears, or when a new deployment breaks instrumentation.

### sort() and sort_desc()

```promql
# Top consumers sorted descending
sort_desc(sum by (pod) (container_memory_working_set_bytes))

# Lowest first
sort(node_filesystem_avail_bytes)
```

### label_join() and label_replace()

Manipulate labels:

```promql
# Join labels into a new label
# label_join(v, dst_label, separator, src_label1, src_label2, ...)
label_join(up{job="api"}, "target_info", ":", "job", "instance")
# Result: target_info="api:10.0.0.1:8080"

# Replace label with regex
# label_replace(v, dst_label, replacement, src_label, regex)
label_replace(up, "short_instance", "$1", "instance", "(.+):\\d+")
# Extracts hostname from "hostname:port" → "hostname"

# Extract namespace from pod name
label_replace(
  kube_pod_info,
  "app",
  "$1",
  "pod",
  "(.*)-[a-z0-9]+-[a-z0-9]+"
)
```

### vector() and scalar()

```promql
# Create a constant vector
vector(1)

# Convert single-element vector to scalar
scalar(count(up))

# Use in expressions
up == bool 1 or vector(0)
```

### histogram_quantile()

See [Histogram and Summary section](#-histogram-and-summary).

---

## 📊 Histogram and Summary

### histogram_quantile()

Calculate quantiles from histogram buckets:

```promql
# Syntax: histogram_quantile(φ, rate(histogram_bucket[range]))

# P50 (median) request latency
histogram_quantile(0.5,
  sum by (le) (rate(http_request_duration_seconds_bucket[5m]))
)

# P95 request latency
histogram_quantile(0.95,
  sum by (le) (rate(http_request_duration_seconds_bucket[5m]))
)

# P99 request latency per handler
histogram_quantile(0.99,
  sum by (handler, le) (rate(http_request_duration_seconds_bucket[5m]))
)
```

> **Critical:** Always use `rate()` on histogram buckets before `histogram_quantile()`. The `le` label must be preserved in the `by` clause.

### Average from Histogram

```promql
# Average request duration
sum(rate(http_request_duration_seconds_sum[5m]))
/
sum(rate(http_request_duration_seconds_count[5m]))

# Average per handler
sum by (handler) (rate(http_request_duration_seconds_sum[5m]))
/
sum by (handler) (rate(http_request_duration_seconds_count[5m]))
```

### Apdex Score from Histogram

```promql
# Apdex with T=0.5s (satisfied ≤ 0.5s, tolerating ≤ 2s)
(
  sum(rate(http_request_duration_seconds_bucket{le="0.5"}[5m]))
  +
  sum(rate(http_request_duration_seconds_bucket{le="2.0"}[5m]))
)
/
(2 * sum(rate(http_request_duration_seconds_count[5m])))
```

### Pre-calculated Quantiles (Summary)

Summaries have pre-calculated quantiles. Just select them:

```promql
# P99 from summary
rpc_duration_seconds{quantile="0.99"}

# But you CANNOT aggregate summaries across instances!
# This is WRONG:
# avg(rpc_duration_seconds{quantile="0.99"})
# Average of percentiles ≠ percentile of averages
```

---

## 🔄 Subqueries

Subqueries allow you to run a range query within an instant query.

### Syntax

```promql
<instant_query>[<range>:<resolution>]
```

- `range` — How far back to look
- `resolution` — Step interval for the inner query (optional; defaults to global evaluation interval)

### Examples

```promql
# Maximum rate of requests over the last hour, evaluated every 1 minute
max_over_time(rate(http_requests_total[5m])[1h:1m])

# Minimum CPU usage over the last day, sampled every 5 minutes
min_over_time(instance:node_cpu_utilisation:ratio_rate5m[1d:5m])

# Standard deviation of error rate over 24 hours
stddev_over_time(
  sum by (job) (rate(http_requests_total{status=~"5.."}[5m]))[24h:5m]
)

# Detect if error rate spiked at any point in the last hour
max_over_time(
  sum(rate(http_requests_total{status=~"5.."}[5m]))
  /
  sum(rate(http_requests_total[5m]))
  [1h:1m]
) > 0.1
```

> **Performance tip:** Subqueries can be expensive. Use recording rules instead when the same subquery appears in multiple places.

---

## 📏 Recording Rules

Recording rules pre-compute expensive PromQL expressions and store results as new time series.

### Naming Convention

```text
level:metric:operations
```

| Part | Description | Examples |
|------|-------------|---------|
| `level` | Aggregation level | `job`, `instance`, `namespace`, `cluster`, `node` |
| `metric` | Base metric name | `http_requests_total`, `node_cpu_seconds_total` |
| `operations` | Applied operations | `rate5m`, `sum`, `ratio`, `p99` |

### Examples

```yaml
groups:
  - name: http_rules
    rules:
      # Level: job, Metric: http_requests, Operations: rate5m
      - record: job:http_requests_total:rate5m
        expr: sum by (job) (rate(http_requests_total[5m]))

      # Error ratio
      - record: job:http_errors:ratio_rate5m
        expr: |
          sum by (job) (rate(http_requests_total{status=~"5.."}[5m]))
          /
          sum by (job) (rate(http_requests_total[5m]))

      # P99 latency
      - record: job:http_request_duration_seconds:p99_rate5m
        expr: |
          histogram_quantile(0.99,
            sum by (job, le) (rate(http_request_duration_seconds_bucket[5m]))
          )

  - name: node_rules
    rules:
      # CPU utilization
      - record: instance:node_cpu_utilisation:ratio_rate5m
        expr: 1 - avg by (instance) (rate(node_cpu_seconds_total{mode="idle"}[5m]))

      # Memory utilization
      - record: instance:node_memory_utilisation:ratio
        expr: 1 - (node_memory_MemAvailable_bytes / node_memory_MemTotal_bytes)

  - name: k8s_rules
    rules:
      # Pod CPU usage
      - record: namespace_pod:container_cpu_usage_seconds_total:sum_rate5m
        expr: |
          sum by (namespace, pod) (
            rate(container_cpu_usage_seconds_total{container!=""}[5m])
          )
```

---

## 🌍 Real-World Query Examples

### CPU Usage

```promql
# Node CPU usage (%)
100 - (avg by (instance) (rate(node_cpu_seconds_total{mode="idle"}[5m])) * 100)

# Per-mode CPU breakdown
sum by (mode) (rate(node_cpu_seconds_total[5m])) * 100

# Container CPU usage (cores)
sum by (pod) (rate(container_cpu_usage_seconds_total{container!=""}[5m]))

# CPU throttling
sum by (pod) (rate(container_cpu_cfs_throttled_seconds_total[5m]))
/
sum by (pod) (rate(container_cpu_usage_seconds_total[5m]))

# CPU usage vs request
sum by (pod) (rate(container_cpu_usage_seconds_total{container!=""}[5m]))
/
sum by (pod) (kube_pod_container_resource_requests{resource="cpu"})
```

### Memory Usage

```promql
# Node memory usage (%)
(1 - node_memory_MemAvailable_bytes / node_memory_MemTotal_bytes) * 100

# Node memory available in GB
node_memory_MemAvailable_bytes / 1024 / 1024 / 1024

# Container memory usage
sum by (pod) (container_memory_working_set_bytes{container!=""})

# Memory usage vs limit
sum by (pod) (container_memory_working_set_bytes{container!=""})
/
sum by (pod) (kube_pod_container_resource_limits{resource="memory"})

# OOM kill detection
increase(kube_pod_container_status_last_terminated_reason{reason="OOMKilled"}[1h])
```

### Disk Usage and I/O

```promql
# Disk usage (%)
(1 - node_filesystem_avail_bytes{fstype!~"tmpfs|overlay"} / node_filesystem_size_bytes) * 100

# Disk space remaining in GB
node_filesystem_avail_bytes{mountpoint="/"} / 1024 / 1024 / 1024

# Predict disk full (hours)
- node_filesystem_avail_bytes{mountpoint="/"}
/
deriv(node_filesystem_avail_bytes{mountpoint="/"}[6h])
/ 3600

# Disk I/O utilization
rate(node_disk_io_time_seconds_total[5m])

# Disk read/write throughput (MB/s)
rate(node_disk_read_bytes_total[5m]) / 1024 / 1024
rate(node_disk_written_bytes_total[5m]) / 1024 / 1024

# IOPS
rate(node_disk_reads_completed_total[5m])
rate(node_disk_writes_completed_total[5m])
```

### Network Traffic

```promql
# Network receive rate (MB/s)
rate(node_network_receive_bytes_total{device!~"lo|veth.*|docker.*|br.*"}[5m]) / 1024 / 1024

# Network transmit rate (MB/s)
rate(node_network_transmit_bytes_total{device!~"lo|veth.*|docker.*|br.*"}[5m]) / 1024 / 1024

# Network errors
rate(node_network_receive_errs_total[5m])
rate(node_network_transmit_errs_total[5m])

# TCP connections by state
node_netstat_Tcp_CurrEstab
node_sockstat_TCP_tw         # TIME_WAIT connections
```

### HTTP: RED Method

```promql
# Request Rate (R)
sum by (handler) (rate(http_requests_total[5m]))

# Error Rate (E) — as ratio
sum by (handler) (rate(http_requests_total{status=~"5.."}[5m]))
/
sum by (handler) (rate(http_requests_total[5m]))

# Error Rate (E) — as percentage
sum by (handler) (rate(http_requests_total{status=~"5.."}[5m]))
/
sum by (handler) (rate(http_requests_total[5m]))
* 100

# Duration / Latency (D) — P50, P95, P99
histogram_quantile(0.50, sum by (handler, le) (rate(http_request_duration_seconds_bucket[5m])))
histogram_quantile(0.95, sum by (handler, le) (rate(http_request_duration_seconds_bucket[5m])))
histogram_quantile(0.99, sum by (handler, le) (rate(http_request_duration_seconds_bucket[5m])))

# Average latency
sum by (handler) (rate(http_request_duration_seconds_sum[5m]))
/
sum by (handler) (rate(http_request_duration_seconds_count[5m]))
```

### Kubernetes-Specific

```promql
# Pod restarts in last hour
increase(kube_pod_container_status_restarts_total[1h]) > 0

# Pods not ready
kube_pod_status_ready{condition="false"} == 1

# Pending pods
kube_pod_status_phase{phase="Pending"} == 1

# Deployments with unavailable replicas
kube_deployment_status_replicas_unavailable > 0

# Resource requests vs limits (over-provisioned)
sum by (namespace) (kube_pod_container_resource_requests{resource="cpu"})
/
sum by (namespace) (kube_pod_container_resource_limits{resource="cpu"})

# Namespace CPU usage vs requests
sum by (namespace) (rate(container_cpu_usage_seconds_total{container!=""}[5m]))
/
sum by (namespace) (kube_pod_container_resource_requests{resource="cpu"})

# Node allocatable vs capacity
kube_node_status_allocatable{resource="cpu"}
/
kube_node_status_capacity{resource="cpu"}

# HPA status
kube_horizontalpodautoscaler_status_current_replicas
/
kube_horizontalpodautoscaler_spec_max_replicas
```

### SLI / SLO Examples

```promql
# Availability SLI (% of successful requests)
sum(rate(http_requests_total{status!~"5.."}[30d]))
/
sum(rate(http_requests_total[30d]))

# Latency SLI (% of requests faster than 500ms)
sum(rate(http_request_duration_seconds_bucket{le="0.5"}[30d]))
/
sum(rate(http_request_duration_seconds_count[30d]))

# Error budget remaining (SLO = 99.9%)
1 - (
  (1 - (
    sum(rate(http_requests_total{status!~"5.."}[30d]))
    /
    sum(rate(http_requests_total[30d]))
  ))
  /
  (1 - 0.999)
)

# Error budget burn rate
(
  1 - sum(rate(http_requests_total{status!~"5.."}[1h])) / sum(rate(http_requests_total[1h]))
)
/
(1 - 0.999)

# Multi-window burn rate alert
# Fast burn: high error rate in short window
(
  1 - sum(rate(http_requests_total{status!~"5.."}[5m])) / sum(rate(http_requests_total[5m]))
) > 14.4 * (1 - 0.999)
and
(
  1 - sum(rate(http_requests_total{status!~"5.."}[1h])) / sum(rate(http_requests_total[1h]))
) > 14.4 * (1 - 0.999)
```

---

## ⚡ Performance and Optimization

### General Rules

1. **Limit time range** — Query `[5m]` instead of `[1h]` when possible
2. **Use recording rules** — Pre-compute dashboard queries
3. **Avoid high cardinality** — Don't use unbounded label values
4. **Use `without` for small exclusions** — More efficient than `by` with many labels
5. **Avoid `{__name__=~".+"}`** — Matches ALL metrics (very expensive)
6. **Avoid regex when exact match works** — `{status="200"}` faster than `{status=~"200"}`

### Query Cost Hierarchy (cheapest to most expensive)

```text
1. Direct metric name + label matchers   http_requests_total{job="api"}
2. Recording rule lookup                 job:http_requests:rate5m
3. rate() on single metric               rate(http_requests_total[5m])
4. Aggregation over instant vector       sum by (job) (rate(...))
5. Binary operations between vectors     rate(...) / rate(...)
6. Regex label matchers                  {status=~"5.."}
7. Subqueries                           max_over_time(rate(...)[1h:1m])
8. __name__ regex match                  {__name__=~"http_.*"}
```

### Recording Rules for Dashboard Performance

```yaml
# Instead of this expensive query in every dashboard panel:
# histogram_quantile(0.99, sum by (handler, le) (rate(http_request_duration_seconds_bucket[5m])))

# Create a recording rule:
- record: handler:http_request_duration_seconds:p99_rate5m
  expr: |
    histogram_quantile(0.99,
      sum by (handler, le) (rate(http_request_duration_seconds_bucket[5m]))
    )
```

---

## ❌ Common Mistakes

### 1. Using rate() on a Gauge

```promql
# WRONG — rate() is for counters only
rate(node_memory_MemAvailable_bytes[5m])

# CORRECT — use delta() or deriv() for gauges
delta(node_memory_MemAvailable_bytes[5m])
deriv(node_memory_MemAvailable_bytes[1h])
```

### 2. Forgetting rate() on a Counter

```promql
# WRONG — raw counter value is meaningless (always increasing)
http_requests_total

# CORRECT — use rate() to get requests per second
rate(http_requests_total[5m])
```

### 3. Wrong Range for rate()

```promql
# WRONG — range shorter than scrape interval means no samples
# (with 15s scrape interval)
rate(http_requests_total[10s])

# CORRECT — range should be ≥ 4× scrape interval
rate(http_requests_total[1m])  # 4 × 15s = 60s
```

### 4. Missing le Label in histogram_quantile()

```promql
# WRONG — le label is not preserved
histogram_quantile(0.99,
  sum by (handler) (rate(http_request_duration_seconds_bucket[5m]))
)

# CORRECT — always include le in by clause
histogram_quantile(0.99,
  sum by (handler, le) (rate(http_request_duration_seconds_bucket[5m]))
)
```

### 5. Averaging Percentiles

```promql
# WRONG — average of percentiles ≠ percentile
avg(histogram_quantile(0.99,
  sum by (instance, le) (rate(http_request_duration_seconds_bucket[5m]))
))

# CORRECT — aggregate buckets first, then calculate percentile
histogram_quantile(0.99,
  sum by (le) (rate(http_request_duration_seconds_bucket[5m]))
)
```

### 6. Label Mismatch in Binary Operations

```promql
# WRONG — left side has {method, handler, status}, right side has {method}
# No matching because label sets differ
http_requests_total{status="500"} / http_requests_total

# CORRECT — use ignoring() or on()
sum by (method) (rate(http_requests_total{status=~"5.."}[5m]))
/
sum by (method) (rate(http_requests_total[5m]))
```

### 7. Comparing Gauge to Rate

```promql
# WRONG — comparing raw gauge (bytes) to a rate (bytes/second)
container_memory_working_set_bytes > rate(container_network_receive_bytes_total[5m])

# These are different units! Always compare like with like.
```

### 8. Using sum() without by()

```promql
# DANGEROUS — sums everything into a single number
sum(rate(http_requests_total[5m]))

# Usually you want grouping
sum by (job) (rate(http_requests_total[5m]))
```

---

## 🐛 Troubleshooting

### No Data Returned

```promql
# 1. Check if metric exists
{__name__=~"http_request.*"}

# 2. Check if labels match
http_requests_total                              # All label combinations
count by (__name__) ({job="my-app"})             # All metrics from job

# 3. Check if target is up
up{job="my-app"}

# 4. Check scrape info
scrape_samples_scraped{job="my-app"}

# 5. Check if metric was dropped by relabeling
# Review metric_relabel_configs in prometheus.yml

# 6. Check staleness (metric not updated recently)
time() - timestamp(http_requests_total{job="my-app"}) > 300
```

### Unexpected Results

```promql
# Rate returns nothing?
# → Range might be too short. Increase range window.
rate(http_requests_total[1m])   # Try [5m] instead

# Value is NaN?
# → Division by zero. Guard with > 0
sum(rate(http_requests_total{status="500"}[5m]))
/
(sum(rate(http_requests_total[5m])) > 0)

# histogram_quantile returns NaN?
# → No samples in the range, or all in +Inf bucket
# Check: sum(rate(http_request_duration_seconds_bucket[5m])) by (le)

# Values seem wrong?
# → Check units (seconds vs milliseconds, bytes vs kilobytes)
# → Check if you need rate() (counter) vs direct value (gauge)
```

### Slow Queries

```promql
# 1. Check query duration
prometheus_engine_query_duration_seconds

# 2. Reduce time range
# Change [1h] to [5m] where possible

# 3. Add label matchers to narrow selection
# Instead of: rate(http_requests_total[5m])
# Use: rate(http_requests_total{job="api"}[5m])

# 4. Use recording rules for repeated expensive queries

# 5. Check cardinality
count by (__name__) ({__name__=~".+"})
# Look for metrics with >10K series
```

---

## 🔗 Related Topics

- [🔥 Prometheus](../prometheus/) — Prometheus architecture, configuration, deployment
- [📊 Grafana](../grafana/) — Visualization and dashboards using PromQL
- [📈 VictoriaMetrics](../victoriametrics/) — MetricsQL extensions to PromQL
- [📊 Observability Overview](../) — Three pillars and tool comparison

---

> **PromQL mastery is the single most impactful skill for metrics-based observability.** Start with basic selectors and `rate()`, build up to aggregations and histogram queries, and use recording rules to make everything fast. Practice in the Prometheus UI before putting queries into dashboards or alerts.
