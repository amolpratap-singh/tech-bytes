# 💻 CPU Troubleshooting

> **Diagnosing and resolving high CPU usage, throttling, and load issues on Linux systems and containers.**

---

## 📋 Table of Contents

- [Understanding CPU Metrics](#understanding-cpu-metrics)
- [High CPU Usage](#high-cpu-usage)
- [Process Consuming 100% CPU](#process-consuming-100-cpu)
- [System vs User CPU](#system-vs-user-cpu)
- [CPU Throttling in Containers](#cpu-throttling-in-containers)
- [Load Average](#load-average)
- [Context Switching](#context-switching)
- [General CPU Diagnostic Commands](#general-cpu-diagnostic-commands)

---

## Understanding CPU Metrics

Before debugging, understand what the metrics mean.

### Key CPU Metrics

| Metric | Meaning | Healthy Range |
|--------|---------|---------------|
| `%us` (user) | Time in user-space code (application) | Depends on workload |
| `%sy` (system) | Time in kernel-space code (system calls, I/O) | < 20% typically |
| `%wa` (iowait) | Waiting for I/O (disk, network) | < 10% |
| `%id` (idle) | CPU doing nothing | > 20% (some headroom) |
| `%si` (softirq) | Handling software interrupts | < 5% |
| `%st` (steal) | Time stolen by hypervisor (VMs) | < 5% |
| Load Average | Running + waiting processes | ≤ CPU count |

### CPU Time Breakdown

```
Total CPU = User + System + IOWait + Idle + Steal + IRQ + SoftIRQ + Nice

High User%   → Application is CPU-intensive
High System% → Too many system calls, context switches, I/O
High IOWait% → Process waiting for disk/network I/O
High Steal%  → VM is starved by the hypervisor
```

---

## High CPU Usage

Overall CPU usage is high, affecting system responsiveness.

### Symptoms

```bash
$ uptime
 14:30:00 up 5 days,  load average: 8.50, 7.20, 5.10
# (on a 4-core system, this is overloaded)

$ top
%Cpu(s): 95.2 us,  3.1 sy,  0.0 ni,  1.2 id,  0.3 wa,  0.0 hi,  0.2 si,  0.0 st
```

### Diagnostic Commands

```bash
# Step 1: Quick overview
top -bn1 | head -20
# Shows overall CPU usage and top processes

# Step 2: Identify the top CPU consumers
ps aux --sort=-%cpu | head -20
# USER  PID  %CPU %MEM  VSZ   RSS TTY  STAT START  TIME COMMAND
# app   1234 98.5  2.3  ...

# Step 3: Real-time monitoring with htop (better than top)
htop
# Features: tree view, per-core usage, search, filter, kill

# Step 4: Per-process CPU usage over time
pidstat 1 5
# Shows CPU usage per process every 1 second, 5 times
# Linux 5.15.0    08/31/2026   _x86_64_
# 14:30:01  UID  PID  %usr %system %guest %wait %CPU  CPU  Command
# 14:30:02  1000 1234 95.0    3.0    0.0   0.0  98.0   2  myapp

# Step 5: Per-CPU core usage
mpstat -P ALL 1 5
# Shows usage for each individual CPU core
# Useful for detecting CPU imbalance

# Step 6: Thread-level CPU usage
ps -eLo pid,tid,pcpu,comm --sort=-%cpu | head -20
# Shows per-thread CPU usage

# Step 7: Check what the process is doing
strace -p <pid> -c              # System call summary
strace -p <pid> -e trace=open   # Trace specific calls

# Step 8: Profile with perf (advanced)
perf top                        # Real-time function profiling
perf record -p <pid> -g -- sleep 30   # Record profile
perf report                     # Analyze recording
```

### Resolution

```bash
# Option 1: Identify and fix the hot code path
perf record -p <pid> -g -- sleep 30
perf report
# Look for functions consuming the most CPU

# Option 2: Limit CPU usage with cgroups
# Using systemd:
systemctl set-property <service> CPUQuota=50%

# Using cgroupv2 directly:
echo '50000 100000' > /sys/fs/cgroup/<group>/cpu.max  # 50% of one CPU

# Option 3: Change process priority (nice)
renice 19 -p <pid>             # Lowest priority
nice -n 19 <command>           # Start with low priority

# Option 4: Scale horizontally
# Add more instances behind a load balancer

# Option 5: Optimize the application
# - Reduce unnecessary work (caching, algorithmic improvements)
# - Fix busy loops or polling
# - Use async I/O instead of threads for I/O-bound work
```

### Prevention

- Set up CPU usage alerts (e.g., >80% sustained for 5 minutes)
- Monitor per-process CPU in production
- Profile applications under load before deploying
- Use resource limits in containers and cgroups

---

## Process Consuming 100% CPU

A single process is pegging one or more CPU cores.

### Symptoms

```bash
$ top
  PID USER      PR  NI    VIRT    RES    SHR S  %CPU %MEM     TIME+ COMMAND
 1234 app       20   0  1.2g   256m    12m R  99.7  3.2  45:22.10 myapp
```

### Common Causes

| Cause | Diagnostic |
|-------|------------|
| Infinite loop in code | CPU stays at 100% indefinitely |
| Busy-wait / spin lock | Thread spinning without sleep |
| Regex backtracking | Pathological regex on specific input |
| Decompression bomb | Processing malicious data |
| GC storm | JVM/Go spending all time in garbage collection |
| Hash collision attack | Hash table performance degrades to O(n) |

### Diagnostic Commands

```bash
# Step 1: Identify the process
top -H -p <pid>          # Show threads within the process
ps -eLo pid,tid,pcpu,comm -p <pid>   # Per-thread view

# Step 2: Get a thread dump (application-specific)
# Java:
kill -3 <pid>                     # Thread dump to stdout
jstack <pid>                      # Thread dump
jcmd <pid> Thread.print           # Alternative

# Go:
kill -SIGQUIT <pid>               # Goroutine dump to stderr
# Or send SIGABRT for full stack trace

# Python:
# Add signal handler in code, or use py-spy:
py-spy top --pid <pid>
py-spy dump --pid <pid>

# Step 3: Profile the process
perf record -p <pid> -g -- sleep 30
perf report

# Step 4: Generate flame graph
perf script | stackcollapse-perf.pl | flamegraph.pl > flame.svg

# Step 5: Check if it's GC
# Java:
jstat -gcutil <pid> 1000          # GC stats every 1s
# If GC% is >50%, it's a GC problem

# Go:
GODEBUG=gctrace=1 ./myapp        # GC trace output

# Step 6: Trace system calls
strace -p <pid> -c -S calls      # Summary sorted by call count
# High count of futex or poll = likely busy-waiting
```

### Resolution

**Infinite Loop:**

```bash
# Generate a core dump for analysis
gcore <pid>                       # Creates core.<pid>
gdb <binary> core.<pid>           # Analyze

# Or get stack trace without core dump
gdb -p <pid> -ex "thread apply all bt" -ex quit
```

**GC Storm (Java):**

```bash
# Check heap usage
jcmd <pid> GC.heap_info

# Trigger manual GC to see if it helps
jcmd <pid> GC.run

# Generate heap dump for analysis
jmap -dump:format=b,file=heap.hprof <pid>
# Analyze with Eclipse MAT, VisualVM, or jhat

# Fix: Increase heap, fix memory leaks, tune GC
# -Xmx4g -XX:+UseG1GC -XX:MaxGCPauseMillis=200
```

**Busy-Wait:**

```bash
# Identify: High CPU but low throughput
strace -p <pid> -c
# If mostly futex(), poll(), or epoll_wait() with high call count:
# = busy-waiting or spin-lock contention

# Fix: Add sleep/yield in loops, use proper synchronization
# Fix: Use blocking I/O or event-driven patterns
```

### Prevention

- Add timeouts to all loops and retries
- Use connection pools with proper limits
- Monitor GC metrics (Java: JMX, Go: runtime/metrics)
- Profile under realistic load before production
- Set CPU limits on containers to contain runaway processes

---

## System vs User CPU

Understanding whether CPU time is spent in application code or kernel operations.

### High User CPU (`%us`)

Application code is consuming CPU. This is normal for compute-intensive workloads.

```bash
# Identify the hot code path
perf top                          # Real-time function profiling
perf record -p <pid> -g -- sleep 30
perf report

# Common causes:
# - Algorithmic inefficiency (O(n²) instead of O(n log n))
# - Unnecessary computation (redundant calculations, no caching)
# - Serialization/deserialization overhead
# - Regex backtracking
# - Compression/encryption
```

### High System CPU (`%sy`)

Kernel is consuming CPU. Often indicates excessive I/O or system call overhead.

```bash
# Identify system calls
strace -p <pid> -c
# Look for high-count system calls

# Common causes and their system calls:
# - Too many small I/O operations → read(), write() (batch I/O instead)
# - Excessive file operations → open(), close(), stat() (cache file handles)
# - Network I/O overhead → sendto(), recvfrom() (use larger buffers)
# - Context switch overhead → futex(), sched_yield() (reduce thread count)
# - Page faults → mmap(), brk() (memory allocation issues)

# Check context switches
pidstat -w 1 5
# cswch/s = voluntary (I/O wait)
# nvcswch/s = involuntary (preempted by scheduler)
```

### High IOWait (`%wa`)

CPU is idle but processes are waiting for I/O.

```bash
# Identify I/O-heavy processes
iotop                             # Real-time I/O usage per process
pidstat -d 1 5                    # I/O stats per process

# Check disk I/O
iostat -x 1 5                    # Disk I/O statistics
# Look for: %util > 80%, await > 10ms, queue length > 2

# Common causes:
# - Slow disk (check iostat await)
# - Too many random I/O operations
# - Swapping (check free -h for swap usage)
# - Logging too much (synchronous writes)
```

### High Steal (`%st`)

Hypervisor is taking CPU time from your VM.

```bash
# Check steal time
top                               # Look at %st

# Causes:
# - Overcommitted host
# - Noisy neighbor on shared infrastructure
# - Burstable instance type exceeded baseline (AWS T-series)

# Resolution:
# - Migrate to dedicated instance
# - Use non-burstable instance type
# - Request different host placement
```

---

## CPU Throttling in Containers

Containers with CPU limits can be throttled even when the host has available CPU.

### Understanding Container CPU Limits

```yaml
# Kubernetes resource spec
resources:
  requests:
    cpu: "250m"     # 0.25 CPU — scheduler guarantee
  limits:
    cpu: "1000m"    # 1.0 CPU — maximum allowed
```

How CFS (Completely Fair Scheduler) enforces limits:

```
CPU limit = cpu.cfs_quota_us / cpu.cfs_period_us

1000m = 100000us / 100000us = 1.0 CPU
500m  = 50000us  / 100000us = 0.5 CPU
250m  = 25000us  / 100000us = 0.25 CPU

If a container uses its quota within the period, it's throttled
until the next period starts.
```

### Detecting Throttling

```bash
# In Kubernetes
kubectl top pod <pod> --containers

# Check cgroup throttling stats
# Inside the container or on the node:
cat /sys/fs/cgroup/cpu/cpu.stat
# nr_periods 12345
# nr_throttled 1000      # Number of times throttled
# throttled_time 50000   # Total throttle time in nanoseconds

# cgroupv2:
cat /sys/fs/cgroup/cpu.stat
# usage_usec 123456789
# nr_periods 12345
# nr_throttled 1000
# throttled_usec 50000000

# Calculate throttle percentage
# throttled% = nr_throttled / nr_periods * 100
# If > 25%, the container is significantly throttled

# Prometheus metrics (if cadvisor is running)
# container_cpu_cfs_throttled_periods_total
# container_cpu_cfs_periods_total
# Throttle ratio = throttled_periods / total_periods
```

### Resolution

```yaml
# Option 1: Increase CPU limits
resources:
  requests:
    cpu: "500m"
  limits:
    cpu: "2000m"     # Increase limit

# Option 2: Remove CPU limits (use only requests)
# Some teams (e.g., Google) recommend NOT setting CPU limits
# because throttling can cause worse latency than sharing
resources:
  requests:
    cpu: "500m"
  # No limit — container can burst

# Option 3: Optimize the application
# - Reduce CPU usage during request processing
# - Use async patterns, batch operations
# - Cache expensive computations
```

### Understanding Throttling Spikes

```
A container with 500m CPU limit (50ms per 100ms period):

Time 0ms:    ┌──────── 50ms of work ────────┐
Time 50ms:   ├──────── THROTTLED ──────────────┤  (waiting for next period)
Time 100ms:  ├──────── 50ms of work ────────┐
Time 150ms:  ├──────── THROTTLED ──────────────┤

Even if the container only needs 60ms total, if it uses
it all at once, it's throttled for 40ms. This causes
latency spikes that look like periodic slowdowns.
```

### Prevention

- Monitor throttling metrics (`container_cpu_cfs_throttled_periods_total`)
- Set CPU requests based on actual P95 usage
- Consider not setting CPU limits if latency matters more than isolation
- Alert on high throttle ratios (>25%)

---

## Load Average

Understanding load average and what it means for system health.

### What Load Average Means

```bash
$ uptime
 14:30:00 up 5 days,  load average: 4.50, 3.20, 2.10
#                                    1min  5min  15min
```

Load average = average number of processes in **running** or **uninterruptible sleep** (waiting for I/O) state.

### Interpreting Load Average

```
Rule of thumb: Compare load average to CPU count

$ nproc
4    # 4 CPU cores

Load avg = 4.0 on 4 cores → 100% utilized, no queuing
Load avg = 2.0 on 4 cores → 50% utilized, healthy
Load avg = 8.0 on 4 cores → 200% — processes waiting in queue
Load avg = 1.0 on 4 cores → 25% utilized, plenty of headroom
```

| Load / CPUs | Meaning |
|-------------|---------|
| < 0.7 | Healthy, good headroom |
| 0.7 - 1.0 | Moderate usage, acceptable |
| 1.0 - 1.5 | High usage, some queuing |
| > 2.0 | Overloaded, processes waiting |
| >> CPU count | Severely overloaded |

### Trend Analysis

```bash
# 1-min, 5-min, 15-min comparison:
# 1min > 5min > 15min → Load is increasing (something just happened)
# 1min < 5min < 15min → Load is decreasing (problem resolving)
# All similar           → Steady state

# Example: 8.50, 3.20, 2.10
# Load spiked recently — investigate what started in the last few minutes
```

### Load Average vs CPU Usage

```bash
# High load but low CPU% → I/O bound (processes waiting for disk/network)
# Check with:
iostat -x 1 5            # Disk I/O
vmstat 1 5               # System overview including I/O wait

# High load and high CPU% → CPU bound
# Check with:
top                      # CPU consumers
ps aux --sort=-%cpu | head -20

# High load but many processes in D state (uninterruptible sleep)
ps aux | awk '$8 ~ /D/ {print}'   # Processes in D state
# These are waiting for I/O and contribute to load average
```

### Prevention

- Set alerts on load average relative to CPU count
- Monitor 1-min vs 15-min trends to detect spikes
- Distinguish CPU-bound from I/O-bound load
- Scale horizontally when consistently at >70% capacity

---

## Context Switching

Excessive context switching can waste CPU cycles.

### What Are Context Switches?

```
Voluntary context switch:  Process voluntarily gives up CPU (I/O wait, sleep)
Involuntary context switch: Scheduler preempts a running process (time slice expired)

Some context switching is normal and healthy.
Excessive involuntary switching = too many threads competing for too few CPUs.
```

### Detecting Context Switch Issues

```bash
# System-wide context switches
vmstat 1 5
# procs -----------memory---------- ---swap-- -----io---- -system-- ------cpu-----
#  r  b   swpd   free   buff  cache   si   so    bi    bo   in    cs us sy id wa st
#  3  0      0 1234567 12345 234567    0    0     0    12  500 15000 45 20 34  1  0
# cs = 15000 context switches per second

# Per-process context switches
pidstat -w 1 5
#          PID   cswch/s nvcswch/s  Command
#         1234    500.00    200.00  myapp
# cswch/s   = voluntary (process chose to yield)
# nvcswch/s = involuntary (scheduler forced)

# Check /proc directly
cat /proc/<pid>/status | grep ctxt
# voluntary_ctxt_switches:    12345
# nonvoluntary_ctxt_switches: 6789
```

### Resolution

```bash
# High involuntary context switches:
# = Too many threads/processes for available CPUs
# Fix: Reduce thread count, increase CPU, use thread pools

# High voluntary context switches:
# = Lots of I/O operations (usually normal)
# Fix: Batch I/O, use async I/O, reduce lock contention

# Reduce thread count in application
# Java: Use thread pools with bounded size
# Go: Goroutines are cheap but contention isn't
# Python: GIL limits true parallelism anyway

# Pin processes to CPUs for latency-sensitive workloads
taskset -c 0,1 ./myapp           # Run on CPUs 0 and 1
```

### Prevention

- Size thread pools based on available CPUs (CPU-bound: N threads, I/O-bound: 2N threads)
- Monitor context switch rates as part of system health
- Use CPU affinity for latency-sensitive applications
- Profile lock contention in multi-threaded applications

---

## General CPU Diagnostic Commands

### Quick Reference

```bash
# Overview
uptime                             # Load average
top -bn1 | head -20               # Quick snapshot
htop                               # Interactive (better)
nproc                              # CPU count

# Per-process
ps aux --sort=-%cpu | head -20     # Top CPU consumers
pidstat 1 5                        # Per-process CPU over time
pidstat -t 1 5                     # Per-thread CPU

# Per-CPU core
mpstat -P ALL 1 5                  # Per-core usage

# Profiling
perf top                           # Real-time profiling
perf record -p <pid> -g -- sleep 30  # Record profile
perf report                        # Analyze recording
perf stat -p <pid> -- sleep 10     # Performance counters

# System calls
strace -p <pid> -c                 # System call summary
strace -p <pid> -e trace=open      # Specific syscalls

# Context switches
vmstat 1 5                         # System-wide
pidstat -w 1 5                     # Per-process

# Container CPU
cat /sys/fs/cgroup/cpu.stat        # CPU throttling stats
cat /sys/fs/cgroup/cpu.max         # CPU limit (cgroupv2)
```

---

## 🔗 Related Topics

- [🚀 DevOps / Linux](../../devops/linux/) — Linux system administration
- [🧰 CLI / Linux](../../cli/linux/) — Linux command reference
- [🧠 Memory Troubleshooting](../memory/) — Memory issues (often related to CPU)
- [🚀 Kubernetes Troubleshooting](../kubernetes/) — Container CPU throttling
- [🐳 Docker Troubleshooting](../docker/) — Docker performance issues
- [📊 Observability / Prometheus](../../observability/prometheus/) — CPU metrics and alerting

---

> **Pro tip:** High load average doesn't always mean high CPU. Check if the load is CPU-bound (`%us` + `%sy` high) or I/O-bound (`%wa` high). The fix is completely different for each case.
