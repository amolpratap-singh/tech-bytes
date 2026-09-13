# 🧠 Memory Troubleshooting

> **Diagnosing and resolving memory issues — from OOM kills to leaks to container memory limits.**

---

## 📋 Table of Contents

- [Understanding Linux Memory](#understanding-linux-memory)
- [High Memory Usage](#high-memory-usage)
- [OOMKiller](#oomkiller)
- [Memory Leaks](#memory-leaks)
- [Cache vs Used Memory](#cache-vs-used-memory)
- [Container OOMKilled](#container-oomkilled)
- [Swap](#swap)
- [General Memory Diagnostic Commands](#general-memory-diagnostic-commands)

---

## Understanding Linux Memory

### Memory Types

```
Physical Memory (RAM)
├── Used
│   ├── Application memory (RSS — Resident Set Size)
│   ├── Shared memory (shared libraries, tmpfs)
│   ├── Kernel memory (slab, page tables)
│   └── Cache/Buffers (file system cache — reclaimable!)
├── Free (unused, immediately available)
└── Available (free + reclaimable cache)

Virtual Memory
├── RSS (Resident Set Size) — physical memory actually used
├── VSZ (Virtual Size) — total virtual address space (not all in RAM)
└── Swap — pages moved to disk when RAM is full
```

### Key Metric: Available vs Free

```bash
$ free -h
               total        used        free      shared  buff/cache   available
Mem:            16Gi        8.5Gi       512Mi       256Mi       7.0Gi       7.2Gi
Swap:           4.0Gi       0.0Gi       4.0Gi

# free (512Mi) ≠ available (7.2Gi)
# available = free + reclaimable cache
# The system is healthy! 7.2Gi available.
# Don't panic about low "free" — Linux uses free RAM as file cache.
```

---

## High Memory Usage

System is running low on available memory.

### Symptoms

```bash
$ free -h
               total        used        free      shared  buff/cache   available
Mem:            16Gi       14.5Gi       128Mi        64Mi       1.3Gi       1.1Gi
# available < 10% of total → getting tight

$ dmesg | tail -5
[12345.678] Out of memory: Killed process 1234 (myapp) total-vm:4096000kB ...
```

### Diagnostic Commands

```bash
# Step 1: Overview
free -h

# Step 2: Top memory consumers (processes)
ps aux --sort=-%mem | head -20
# USER  PID  %CPU %MEM  VSZ   RSS    TTY  STAT START  TIME COMMAND
# app   1234  2.0 45.0  4.0g  7.2g   ?    Sl   08:00 5:00 java -jar myapp.jar

# Step 3: Detailed memory info
cat /proc/meminfo
# MemTotal:       16384000 kB
# MemFree:          131072 kB
# MemAvailable:    1126400 kB
# Buffers:          102400 kB
# Cached:          1228800 kB
# SwapTotal:       4096000 kB
# SwapFree:        4096000 kB
# Slab:             256000 kB

# Step 4: Per-process detailed memory
# RSS = physical memory; VSZ = virtual (may be much larger)
ps -eo pid,rss,vsz,comm --sort=-rss | head -20

# Step 5: Even more detail with smaps
cat /proc/<pid>/smaps_rollup
# Rss:              7340032 kB
# Pss:              7200000 kB     (Proportional — shared pages divided)
# Shared_Clean:      140032 kB
# Shared_Dirty:            0 kB
# Private_Clean:     102400 kB
# Private_Dirty:    7097600 kB     (This is "your" memory)

# Step 6: Memory map of a process
pmap -x <pid> | tail -5
# Total mapped memory for the process

# Step 7: Track memory over time
# Using pidstat:
pidstat -r 1 10
# Columns: minflt/s, majflt/s, VSZ, RSS, %MEM

# Step 8: Check kernel memory (slab)
slabtop                            # Real-time slab usage
cat /proc/slabinfo                 # Raw slab data
```

### Resolution

```bash
# Identify the top consumer and investigate
# Is it expected? (Database, cache, JVM heap)
# Is it growing? (Memory leak)

# Quick relief:
# 1. Restart the process
systemctl restart <service>

# 2. Drop file system caches (safe, recoverable)
sync                               # Flush dirty pages first
echo 3 > /proc/sys/vm/drop_caches # Free cached pages

# 3. Kill a specific process
kill <pid>

# Long-term:
# - Fix memory leaks
# - Right-size application memory settings
# - Add more RAM
# - Tune JVM heap, GC settings
```

### Prevention

- Monitor available memory (not just free)
- Alert when available < 20% of total
- Set memory limits for all services
- Profile memory usage during load tests
- Use memory-efficient data structures

---

## OOMKiller

Linux OOM Killer activates when the system runs out of memory, killing processes to reclaim RAM.

### How OOMKiller Works

```
When available memory → 0:
1. Kernel triggers OOM Killer
2. Kernel scores each process (oom_score)
3. Process with highest score gets killed (SIGKILL)
4. Score based on: memory usage, process age, priority

oom_score_adj ranges from -1000 (never kill) to +1000 (kill first)
```

### Symptoms

```bash
# System logs
$ dmesg | grep -i oom
[12345.678] myapp invoked oom-killer: gfp_mask=0x100cca, order=0
[12345.679] Out of memory: Killed process 1234 (myapp) total-vm:4096000kB, anon-rss:3800000kB

$ journalctl -k | grep -i oom
Aug 31 14:30:00 myhost kernel: Out of memory: Killed process 1234 (myapp)

# Process was killed with signal 9 (SIGKILL)
# Exit code = 137 (128 + 9)
```

### Diagnostic Commands

```bash
# Step 1: Check if OOM killer was triggered
dmesg | grep -i "out of memory"
dmesg | grep -i "oom"
journalctl -k --since "1 hour ago" | grep -i oom

# Step 2: Find which process was killed
dmesg | grep "Killed process"
# Killed process 1234 (myapp) total-vm:4096000kB, anon-rss:3800000kB

# Step 3: Check OOM scores for current processes
# Higher score = more likely to be killed
for pid in $(ls /proc/*/oom_score 2>/dev/null | grep -oP '\d+'); do
  echo "PID=$pid Score=$(cat /proc/$pid/oom_score 2>/dev/null) $(cat /proc/$pid/cmdline 2>/dev/null | tr '\0' ' ')"
done | sort -t= -k3 -n -r | head -20

# Step 4: Check OOM score adjustment
cat /proc/<pid>/oom_score_adj
# -1000 = never kill (e.g., init, sshd)
# 0     = default
# +1000 = kill first

# Step 5: Check system memory at time of OOM
# dmesg output includes memory state at OOM time:
dmesg | grep -A20 "Out of memory"
```

### Resolution

```bash
# Immediate: Restart the killed process
systemctl restart <service>

# Protect critical processes from OOM killer
echo -1000 > /proc/<pid>/oom_score_adj     # Never kill
# Or in systemd:
# [Service]
# OOMScoreAdjust=-1000

# Increase memory
# - Add physical RAM
# - Increase VM memory
# - Enable swap (temporary relief)

# Fix the root cause
# - Reduce memory usage of the offending process
# - Fix memory leaks
# - Set proper memory limits
```

### Prevention

- Set OOM score adjustments to protect critical services
- Monitor memory usage trends
- Configure alerting for memory pressure
- Use cgroups/containers to limit per-service memory
- Add swap as a safety net (not a primary solution)

---

## Memory Leaks

Process memory grows continuously without being freed.

### Symptoms

```bash
# Memory grows over time
$ pidstat -r -p <pid> 60 10
# RSS increases with each sample, never goes down

# Or watch with top
$ top -d 60 -p <pid>
# RES column keeps growing
```

### Identifying Memory Leaks

```bash
# Step 1: Confirm memory is growing
# Record RSS over time:
while true; do
  echo "$(date): $(ps -o rss= -p <pid>)" >> /tmp/mem_usage.log
  sleep 60
done

# Step 2: Differentiate leak from cache
# A leak: RSS grows AND doesn't decrease after idle period
# Not a leak: RSS grows under load but stabilizes

# Step 3: Language-specific profiling

# Python:
pip install memory-profiler
python -m memory_profiler myscript.py
# Or use tracemalloc:
# import tracemalloc
# tracemalloc.start()
# ... code ...
# snapshot = tracemalloc.take_snapshot()
# for stat in snapshot.statistics('lineno')[:10]:
#     print(stat)

# Python — objgraph:
pip install objgraph
# import objgraph
# objgraph.show_most_common_types(limit=20)
# objgraph.show_growth()

# Go:
# Built-in profiling:
# import _ "net/http/pprof"
# go tool pprof http://localhost:6060/debug/pprof/heap
# go tool pprof -diff_base=base.prof current.prof

# Java:
# Heap dump and analysis:
jmap -dump:format=b,file=heap.hprof <pid>
# Open with Eclipse MAT or VisualVM

# C/C++:
valgrind --leak-check=full ./myapp
# AddressSanitizer:
# gcc -fsanitize=address -g myapp.c -o myapp

# Step 4: General approach — take snapshots over time
# Compare memory profiles at different points:
# - After startup (baseline)
# - After load test
# - After idle period (should return near baseline)
# Growing objects that aren't freed = likely leak
```

### Common Memory Leak Patterns

| Pattern | Language | Fix |
|---------|----------|-----|
| Unbounded cache | All | Set max size and eviction policy |
| Event listener not removed | JS, Python | Remove listeners on cleanup |
| Connection not closed | All | Use context managers / defer / try-finally |
| Growing list/map | All | Clear after processing |
| Circular references | Python, JS | Break cycles, use weak references |
| ThreadLocal not cleaned | Java | Remove in finally block |
| Global state accumulation | All | Bound all collections |

### Resolution

```bash
# Short-term: Schedule periodic restarts
# In systemd:
# RuntimeMaxSec=86400   # Restart after 24 hours

# In Kubernetes:
# Use a CronJob to periodically restart pods
# Or set memory limits so OOMKill triggers restart

# Long-term: Fix the leak
# 1. Profile to identify the leaking allocation
# 2. Trace the code path that allocates without freeing
# 3. Add proper cleanup (close, dispose, remove)
# 4. Add memory limit checks and circuit breakers
```

### Prevention

- Profile memory under load as part of CI/CD
- Set memory limits in containers (acts as a safety net)
- Use memory-bounded collections (LRU caches, ring buffers)
- Review code for resource cleanup in error paths
- Monitor process RSS trends in production

---

## Cache vs Used Memory

Linux aggressively caches disk reads in unused RAM. This is a feature, not a problem.

### Understanding Linux Page Cache

```bash
$ free -h
               total        used        free      shared  buff/cache   available
Mem:            16Gi        4.0Gi       512Mi       256Mi      11.5Gi      11.7Gi

# Interpretation:
# used (4.0Gi)      = application memory
# buff/cache (11.5Gi) = file system cache (reclaimable)
# free (512Mi)       = completely unused
# available (11.7Gi) = free + reclaimable cache
#
# This system is HEALTHY! Linux is using free RAM as cache.
# When applications need more memory, cache is automatically reclaimed.
```

### When Cache Is a Problem

```bash
# Cache is a problem ONLY when:
# 1. available memory is low (< 10% of total)
# 2. Swap is being used AND applications are slow

# Check if cache is being properly reclaimed
vmstat 1 5
# Look at si/so (swap in/out):
# si > 0 or so > 0 = swapping, cache isn't being reclaimed fast enough

# Check actual pressure
cat /proc/pressure/memory
# some avg10=0.00 avg60=0.00 avg300=0.00 total=0
# full avg10=0.00 avg60=0.00 avg300=0.00 total=0
# Non-zero values = memory pressure
```

### Manually Dropping Cache

```bash
# Only do this for debugging/testing. Linux manages cache well.
# This is NOT a fix for memory problems.

sync                              # Flush dirty pages to disk first
echo 1 > /proc/sys/vm/drop_caches  # Drop page cache
echo 2 > /proc/sys/vm/drop_caches  # Drop dentries and inodes
echo 3 > /proc/sys/vm/drop_caches  # Drop both

# After dropping:
free -h                            # available should not change much
# Because "available" already counted cache as reclaimable
```

### Prevention

- Educate your team: high buff/cache is normal and healthy
- Monitor **available** memory, not free memory
- Don't add swap "because free memory is low" — check available first
- Use `free -h` and look at the available column

---

## Container OOMKilled

Containers killed by the kernel when they exceed their memory limit.

### How Container Memory Limits Work

```yaml
# Kubernetes
resources:
  requests:
    memory: "256Mi"    # Scheduler guarantee (affects pod placement)
  limits:
    memory: "512Mi"    # Hard limit — container killed if exceeded
```

```bash
# Docker
docker run --memory=512m myapp

# The container's cgroup has a hard memory limit.
# When the container tries to allocate beyond the limit:
# 1. Kernel attempts to reclaim memory (page cache, etc.)
# 2. If that fails, OOM killer kills the container process
# 3. Exit code = 137 (SIGKILL)
```

### Detecting Container OOMKilled

```bash
# Kubernetes
kubectl describe pod <pod>
# Look for:
#   Last State:     Terminated
#     Reason:       OOMKilled
#     Exit Code:    137

kubectl get pod <pod> -o jsonpath='{.status.containerStatuses[0].lastState.terminated.reason}'
# OOMKilled

# Docker
docker inspect <container> --format='{{.State.OOMKilled}}'
# true

# On the host
dmesg | grep "memory cgroup"
# memory cgroup out of memory: Killed process 1234 (java)
```

### JVM in Containers (Common Pitfall)

```
Problem: JVM doesn't see container memory limit by default (older JVM).
JVM sets heap based on HOST memory, not container limit.

Host: 16GB RAM → JVM sets -Xmx ~4GB
Container limit: 512MB
Result: JVM tries to use 4GB → OOMKilled!
```

```bash
# Java 8u191+ and Java 11+: Container-aware by default
# But always set explicit limits:
java -Xmx384m -Xms256m -jar app.jar

# Memory budget:
# Container limit: 512Mi
# -Xmx:           384Mi   (heap)
# Off-heap:        ~128Mi  (metaspace, threads, native, GC overhead)
# Rule of thumb: Xmx = container_limit * 0.75
```

```yaml
# Kubernetes spec for JVM app
containers:
  - name: myapp
    image: myapp:latest
    resources:
      requests:
        memory: "512Mi"
      limits:
        memory: "512Mi"     # Request = Limit for predictable behavior
    env:
      - name: JAVA_OPTS
        value: "-Xmx384m -Xms256m -XX:MaxMetaspaceSize=64m"
```

### Other Language Considerations

```bash
# Python: No hard memory limit, but can use resource module
import resource
resource.setrlimit(resource.RLIMIT_AS, (512 * 1024 * 1024, -1))

# Go: Memory is managed by GC, but GOMEMLIMIT helps
# Set GOMEMLIMIT to ~90% of container limit
GOMEMLIMIT=460MiB ./myapp

# Node.js: Default heap is ~1.5GB, set explicitly
node --max-old-space-size=384 app.js
```

### Resolution

```yaml
# Step 1: Check actual memory usage
kubectl top pod <pod> --containers

# Step 2: Right-size limits based on actual usage
# Set limits at P99 usage + 20% headroom
resources:
  requests:
    memory: "256Mi"    # P50 usage
  limits:
    memory: "512Mi"    # P99 usage + headroom

# Step 3: For JVM apps, tune heap size
# Step 4: For all apps, profile memory usage under load
# Step 5: Fix memory leaks if usage grows unbounded
```

### Prevention

- Always set memory limits in Kubernetes
- Monitor memory usage trends (not just current values)
- Profile applications under load to determine correct limits
- Set memory requests = limits for predictable behavior (Guaranteed QoS)
- For JVM: Always set explicit `-Xmx` and account for off-heap memory
- Use Vertical Pod Autoscaler (VPA) for automatic right-sizing recommendations

---

## Swap

Understanding when swap helps and when it hurts.

### What Is Swap?

```
Swap = disk space used as overflow for RAM.
When RAM is full, the kernel moves inactive pages to swap (swap out).
When those pages are needed, they're moved back (swap in).

Swap is ~100-1000x slower than RAM.
```

### Checking Swap Status

```bash
# Overview
free -h
# Swap:           4.0Gi       1.2Gi       2.8Gi

# Detailed
swapon --show
# NAME      TYPE  SIZE  USED  PRIO
# /dev/sda2 part  4.0G  1.2G  -2

# Per-process swap usage
for pid in $(ls /proc/*/status 2>/dev/null | grep -oP '\d+'); do
  swap=$(grep VmSwap /proc/$pid/status 2>/dev/null | awk '{print $2}')
  if [ -n "$swap" ] && [ "$swap" -gt 0 ]; then
    name=$(cat /proc/$pid/cmdline 2>/dev/null | tr '\0' ' ')
    echo "${swap}kB PID=$pid $name"
  fi
done | sort -nr | head -20

# Swap activity (real-time)
vmstat 1 5
# si = swap in (pages read from swap)
# so = swap out (pages written to swap)
# Both > 0 = active swapping = performance degradation
```

### When Swap Is Good

```
1. Safety net: Prevents OOM kills during temporary memory spikes
2. Moving inactive pages to disk frees RAM for active use
3. Required for hibernation
4. Allows systems to handle more processes than RAM allows
```

### When Swap Is Bad

```
1. Databases: Random access to swapped pages = terrible performance
2. Latency-sensitive apps: Swap adds unpredictable latency
3. Containers: Kubernetes usually expects no swap (and may disable it)
4. Active swapping: If si/so are consistently > 0, you need more RAM
```

### Tuning Swappiness

```bash
# Check current swappiness (0-200, default 60)
cat /proc/sys/vm/swappiness
# 60

# swappiness values:
# 0   = Swap only to avoid OOM (prefer OOM to swap)
# 1   = Minimal swapping
# 60  = Default (balanced)
# 100 = Aggressively swap
# 200 = Maximum swapping (cgroupv2)

# For servers with databases or latency-sensitive apps:
sudo sysctl vm.swappiness=10

# Persist across reboots:
echo "vm.swappiness=10" | sudo tee -a /etc/sysctl.d/99-swappiness.conf
```

### Swap in Kubernetes

```bash
# Kubernetes traditionally requires swap to be disabled
sudo swapoff -a
# Remove swap entries from /etc/fstab

# Since Kubernetes 1.28: Swap support is available (beta)
# Enable with: --feature-gates=NodeSwap=true
# Configure in kubelet config:
# memorySwap:
#   swapBehavior: LimitedSwap   # or NoSwap
```

### Prevention

- Monitor swap usage and alert on active swapping (si/so > 0)
- Set appropriate swappiness for your workload
- If swapping is frequent, add more RAM rather than more swap
- In Kubernetes, disable swap unless you specifically configure swap support
- For databases: set swappiness to 1 or disable swap entirely

---

## General Memory Diagnostic Commands

### Quick Reference

```bash
# Overview
free -h                              # Memory overview
cat /proc/meminfo                    # Detailed memory info
vmstat 1 5                           # System activity

# Per-process
ps aux --sort=-%mem | head -20       # Top memory consumers
top -o %MEM                          # Sort by memory in top
pidstat -r 1 5                       # Memory stats over time
pmap -x <pid>                        # Process memory map
cat /proc/<pid>/smaps_rollup         # Detailed RSS breakdown

# OOM
dmesg | grep -i oom                  # OOM events
journalctl -k | grep -i oom         # OOM in journal
cat /proc/<pid>/oom_score            # OOM score

# Swap
swapon --show                        # Swap devices
vmstat 1 5                           # si/so columns
cat /proc/swaps                      # Swap usage

# Container memory
cat /sys/fs/cgroup/memory.current    # Current usage (cgroupv2)
cat /sys/fs/cgroup/memory.max        # Limit (cgroupv2)
cat /sys/fs/cgroup/memory.stat       # Detailed stats

# Cache
slabtop                              # Kernel slab cache
cat /proc/buddyinfo                  # Memory fragmentation
```

---

## 🔗 Related Topics

- [🚀 DevOps / Linux](../../devops/linux/) — Linux system administration
- [🧰 CLI / Linux](../../cli/linux/) — Linux command reference
- [💻 CPU Troubleshooting](../cpu/) — CPU and performance issues
- [🚀 Kubernetes Troubleshooting](../kubernetes/) — Container OOMKilled
- [🐳 Docker Troubleshooting](../docker/) — Docker memory issues
- [📊 Observability / Prometheus](../../observability/prometheus/) — Memory metrics and alerting
- [📝 Real-World Incidents](../../real-world/production-incidents/) — OOM incident case study

---

> **Pro tip:** When someone says "the server is out of memory," always check `free -h` and look at the **available** column, not the free column. Linux using free RAM as file cache is normal and healthy — it's not a problem unless available memory is low.
