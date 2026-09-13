# 🐳 Docker Troubleshooting

> **Diagnosing and fixing container issues — from build failures to runtime problems to disk exhaustion.**

---

## 📋 Table of Contents

- [Container Won't Start](#container-wont-start)
- [Container Exits Immediately](#container-exits-immediately)
- [Port Conflicts](#port-conflicts)
- [Image Build Failures](#image-build-failures)
- [Volume Permission Issues](#volume-permission-issues)
- [Network Connectivity Issues](#network-connectivity-issues)
- [Disk Space Issues](#disk-space-issues)
- [Performance Issues](#performance-issues)
- [General Diagnostic Commands](#general-diagnostic-commands)

---

## Container Won't Start

The container fails to start or exits with an error immediately.

### Symptoms

```bash
$ docker run myapp
docker: Error response from daemon: OCI runtime create failed: ...

$ docker ps -a
CONTAINER ID   IMAGE   COMMAND   CREATED   STATUS                     NAMES
abc123         myapp   "/app"    5s ago    Exited (1) 3 seconds ago   myapp
```

### Understanding Exit Codes

| Exit Code | Meaning | Common Cause |
|-----------|---------|--------------|
| 0 | Success | Process completed normally (not a daemon?) |
| 1 | General error | Application error, exception, misconfiguration |
| 2 | Misuse of shell command | Wrong shell syntax in CMD/ENTRYPOINT |
| 126 | Permission denied | Binary not executable |
| 127 | Command not found | Wrong ENTRYPOINT/CMD, binary missing |
| 128+N | Killed by signal N | 137 = SIGKILL (OOM), 143 = SIGTERM |
| 137 | OOM Killed or `kill -9` | Container exceeded memory limit |
| 139 | Segmentation fault (SIGSEGV) | Application crash |
| 143 | Graceful termination (SIGTERM) | Docker stop or orchestrator shutdown |

### Diagnostic Commands

```bash
# Step 1: Check container status and exit code
docker ps -a --filter name=<container>
docker inspect <container> --format='{{.State.ExitCode}}'
docker inspect <container> --format='{{.State.Error}}'

# Step 2: Check logs
docker logs <container>
docker logs <container> --tail 50

# Step 3: Check the entrypoint and command
docker inspect <image> --format='{{.Config.Entrypoint}}'
docker inspect <image> --format='{{.Config.Cmd}}'

# Step 4: Try running interactively
docker run --rm -it <image> /bin/sh
# If /bin/sh doesn't exist:
docker run --rm -it <image> /bin/bash

# Step 5: Check if the binary exists and is executable
docker run --rm <image> ls -la /app
docker run --rm <image> which <entrypoint-binary>

# Step 6: Check environment variables
docker run --rm <image> env
```

### Resolution

**Exit Code 127 — Command Not Found:**

```dockerfile
# Wrong: Binary doesn't exist in the image
CMD ["myapp"]

# Fix: Ensure binary is in PATH or use absolute path
CMD ["/usr/local/bin/myapp"]

# Common with multi-stage builds — forgot to copy binary
FROM builder AS build
RUN go build -o /app

FROM alpine
COPY --from=build /app /app      # Don't forget this!
CMD ["/app"]
```

**Exit Code 126 — Permission Denied:**

```dockerfile
# Ensure binary is executable
RUN chmod +x /app/entrypoint.sh
CMD ["/app/entrypoint.sh"]
```

**Exit Code 1 — Application Error:**

```bash
# Read the logs carefully
docker logs <container>

# Common causes:
# - Missing environment variables
# - Cannot connect to database/service
# - Configuration file missing or invalid
# - Port already in use inside container

# Run with environment variables
docker run --rm -e DB_HOST=localhost -e DB_PORT=5432 <image>
```

### Prevention

- Test images locally with `docker run` before deploying
- Use health checks in Dockerfiles
- Pin base image versions (don't use `latest`)
- Use multi-stage builds to minimize image size and dependencies

---

## Container Exits Immediately

Container starts, does its work (or nothing), and exits with code 0.

### Symptoms

```bash
$ docker run -d myapp
<container_id>

$ docker ps    # Container not listed!

$ docker ps -a
CONTAINER ID   IMAGE   STATUS                   NAMES
abc123         myapp   Exited (0) 2 seconds ago myapp
```

### Possible Causes

| Cause | Solution |
|-------|----------|
| Process runs in background (daemonizes) | Run in foreground |
| CMD runs a script that completes | Keep the process running |
| No CMD or ENTRYPOINT | Add proper CMD |
| Shell form CMD with no lasting process | Use exec form |

### Diagnostic Commands

```bash
# Check what the container was supposed to run
docker inspect <image> --format='Entrypoint: {{.Config.Entrypoint}} CMD: {{.Config.Cmd}}'

# Check logs
docker logs <container>

# Try running interactively
docker run --rm -it <image> /bin/sh
```

### Resolution

**Process Daemonizes (Detaches):**

```dockerfile
# Wrong: nginx daemonizes by default
CMD ["nginx"]

# Fix: Run nginx in foreground
CMD ["nginx", "-g", "daemon off;"]
```

**Script Completes Without Keeping Process Alive:**

```dockerfile
# Wrong: Script runs and exits
CMD ["./setup.sh"]

# Fix: Script should exec the main process
# In setup.sh:
#!/bin/sh
./setup-config.sh
exec ./myapp    # exec replaces the shell, keeps PID 1
```

**Shell Form vs Exec Form:**

```dockerfile
# Shell form: Runs under /bin/sh -c
# The shell may not forward signals properly
CMD echo "hello" && myapp

# Exec form: Runs directly, gets signals correctly
CMD ["myapp", "--config", "/etc/myapp/config.yaml"]
```

### Prevention

- Always use exec form for CMD and ENTRYPOINT
- Ensure the main process runs in the foreground
- Use `exec` in shell scripts to replace the shell with the main process
- Add a HEALTHCHECK to detect unhealthy containers

---

## Port Conflicts

Container can't bind to a port because it's already in use.

### Symptoms

```bash
$ docker run -p 8080:80 nginx
Error: Bind for 0.0.0.0:8080 failed: port is already allocated

# Or
Error starting userland proxy: listen tcp4 0.0.0.0:8080: bind: address already in use
```

### Diagnostic Commands

```bash
# Find what's using the port
ss -tlnp | grep 8080
lsof -i :8080
netstat -tlnp | grep 8080

# Check running Docker containers using the port
docker ps --format '{{.Ports}}' | grep 8080
```

### Resolution

```bash
# Option 1: Use a different host port
docker run -p 8081:80 nginx

# Option 2: Stop the conflicting process
kill <pid>                       # Graceful
sudo kill -9 <pid>               # Force

# Option 3: Stop the conflicting container
docker stop <container>

# Option 4: Bind to a specific interface
docker run -p 127.0.0.1:8080:80 nginx    # Only localhost
```

### Prevention

- Use docker-compose to manage port assignments across services
- Document port assignments for your stack
- Use environment variables for port configuration
- Bind to specific interfaces in production

---

## Image Build Failures

Docker build fails at various stages.

### Symptoms

```bash
$ docker build -t myapp .
ERROR: failed to solve: process "/bin/sh -c apt-get install -y curl" did not complete successfully: exit code: 100
```

### Common Build Failures

**Package Installation Fails:**

```dockerfile
# Wrong: Cache is stale or missing
RUN apt-get install -y curl

# Fix: Always update package lists first
RUN apt-get update && apt-get install -y curl \
    && rm -rf /var/lib/apt/lists/*
```

**COPY/ADD File Not Found:**

```bash
# Error: COPY failed: file not found in build context
$ docker build -t myapp .

# Check: Is the file in the build context?
ls -la ./path/to/file

# Check: Is it excluded by .dockerignore?
cat .dockerignore
```

```dockerfile
# Wrong: Absolute path (outside build context)
COPY /home/user/config.yaml /app/

# Fix: Use relative path within build context
COPY config.yaml /app/
```

**Layer Cache Issues:**

```bash
# Force rebuild without cache
docker build --no-cache -t myapp .

# Rebuild from a specific stage
docker build --target=build -t myapp .
```

```dockerfile
# Optimize cache: Put rarely-changing layers first
FROM python:3.12-slim

# These change rarely — cached
COPY requirements.txt .
RUN pip install -r requirements.txt

# This changes often — not cached but above layers are
COPY . .
CMD ["python", "app.py"]
```

**Multi-stage Build — Missing Artifact:**

```dockerfile
# Wrong: Copy from wrong stage or wrong path
FROM golang:1.22 AS builder
RUN go build -o /app/server ./cmd/server

FROM alpine:3.19
COPY --from=builder /app/myserver /app/server     # Wrong name!

# Fix: Match the exact path and filename
COPY --from=builder /app/server /app/server
```

### Diagnostic Commands

```bash
# Build with verbose output
docker build --progress=plain -t myapp .

# Build with no cache
docker build --no-cache -t myapp .

# Check build context size
du -sh .
cat .dockerignore

# Inspect intermediate layers (BuildKit)
docker buildx build --target=<stage> -t test .
docker run --rm -it test /bin/sh
```

### Prevention

- Order Dockerfile instructions from least to most frequently changing
- Use `.dockerignore` to exclude unnecessary files
- Pin package versions in Dockerfiles
- Use multi-stage builds to minimize final image size
- Test builds in CI before merging

---

## Volume Permission Issues

Files in volumes are inaccessible due to permission mismatches.

### Symptoms

```bash
$ docker run -v $(pwd)/data:/app/data myapp
Permission denied: /app/data/output.log

# Or
$ ls -la data/
-rw------- 1 root root 1024 Aug 31 10:00 output.log
# Host user can't read files created by container
```

### Possible Causes

| Cause | Details |
|-------|---------|
| Container runs as root, host user is non-root | Files created as root in mounted volume |
| Container runs as non-root, volume owned by root | Container can't write to volume |
| SELinux/AppArmor blocking access | Security module denying access |
| UID/GID mismatch | Container user UID ≠ host user UID |

### Diagnostic Commands

```bash
# Check what user the container runs as
docker inspect <image> --format='{{.Config.User}}'
docker run --rm <image> id

# Check volume file permissions
docker run --rm -v $(pwd)/data:/app/data <image> ls -la /app/data

# Check SELinux context (if applicable)
ls -laZ data/
```

### Resolution

**Match UID/GID:**

```dockerfile
# Create a user with the same UID as the host user
ARG UID=1000
ARG GID=1000
RUN groupadd -g $GID appgroup && \
    useradd -u $UID -g $GID -m appuser
USER appuser
```

```bash
# Build with matching UID
docker build --build-arg UID=$(id -u) --build-arg GID=$(id -g) -t myapp .
```

**Fix Permissions in Entrypoint:**

```bash
#!/bin/sh
# entrypoint.sh — fix ownership at startup
chown -R appuser:appgroup /app/data
exec gosu appuser "$@"
```

**SELinux (RHEL/CentOS/Fedora):**

```bash
# Add :z (shared) or :Z (private) label
docker run -v $(pwd)/data:/app/data:z myapp
```

### Prevention

- Define a non-root USER in Dockerfiles
- Use named volumes instead of bind mounts when possible
- Document expected UID/GID for volume access
- Use init containers to fix permissions before the app starts

---

## Network Connectivity Issues

Container can't reach other containers, host services, or external networks.

### Symptoms

```bash
# Container can't reach another container
$ docker exec myapp curl http://mydb:5432
curl: (6) Could not resolve host: mydb

# Container can't reach the internet
$ docker exec myapp curl https://api.example.com
curl: (28) Connection timed out
```

### Possible Causes

| Cause | Diagnostic |
|-------|------------|
| Containers on different networks | `docker network inspect` |
| DNS resolution failing | Container can't resolve names |
| Firewall blocking traffic | iptables rules, host firewall |
| Docker network misconfigured | Bridge, overlay issues |
| Host networking mode not set | Using wrong network driver |

### Diagnostic Commands

```bash
# Check container's network settings
docker inspect <container> --format='{{json .NetworkSettings.Networks}}' | jq

# List Docker networks
docker network ls

# Inspect a network
docker network inspect <network>

# Test DNS from inside container
docker exec <container> nslookup <hostname>
docker exec <container> cat /etc/resolv.conf

# Test connectivity
docker exec <container> ping <target>
docker exec <container> curl -v http://<target>:<port>

# Check iptables rules (may affect Docker networking)
sudo iptables -L -n | grep -i docker
```

### Resolution

**Containers on Different Networks:**

```bash
# Connect container to the right network
docker network connect <network> <container>

# Or use docker-compose (containers in same compose file share a network)
```

```yaml
# docker-compose.yml
services:
  app:
    image: myapp
    networks:
      - backend
  db:
    image: postgres:16
    networks:
      - backend

networks:
  backend:
    driver: bridge
```

**DNS Resolution:**

```bash
# Docker's embedded DNS only works with user-defined networks
# Default bridge network does NOT have DNS resolution

# Wrong: Using default bridge
docker run -d --name mydb postgres
docker run -d --name myapp myapp    # Can't resolve "mydb"

# Fix: Use a user-defined network
docker network create mynet
docker run -d --name mydb --network mynet postgres
docker run -d --name myapp --network mynet myapp
```

**Can't Reach External Network:**

```bash
# Check if container can reach gateway
docker exec <container> ip route
docker exec <container> ping 8.8.8.8

# Check Docker daemon DNS settings
cat /etc/docker/daemon.json
# Add DNS servers:
# {"dns": ["8.8.8.8", "8.8.4.4"]}

# Restart Docker daemon after changing config
sudo systemctl restart docker
```

### Prevention

- Always use user-defined networks (never the default bridge)
- Use docker-compose for multi-container applications
- Document network architecture and port mappings
- Test container networking as part of CI/CD

---

## Disk Space Issues

Docker is consuming too much disk space, causing failures.

### Symptoms

```bash
$ docker build -t myapp .
ERROR: no space left on device

$ docker run myapp
Error: no space left on device

$ df -h
Filesystem      Size  Used Avail Use% Mounted on
/dev/sda1       50G   48G   2G   96%  /
```

### Diagnostic Commands

```bash
# Docker disk usage overview
docker system df
# TYPE            TOTAL   ACTIVE  SIZE      RECLAIMABLE
# Images          45      10      12.5GB    8.2GB (65%)
# Containers      15      3       2.1GB     1.8GB (85%)
# Local Volumes   12      5       5.3GB     3.1GB (58%)
# Build Cache     -       -       4.2GB     4.2GB

# Detailed breakdown
docker system df -v

# Find large images
docker images --format '{{.Size}}\t{{.Repository}}:{{.Tag}}' | sort -hr | head -20

# Find large containers (including stopped)
docker ps -a --size --format '{{.Size}}\t{{.Names}}' | sort -hr | head -20

# Find large volumes
docker volume ls
docker system df -v | grep -A100 "Local Volumes"
```

### Resolution — Prune Strategies

```bash
# Conservative: Remove only dangling (untagged) images
docker image prune

# Moderate: Remove stopped containers, dangling images, unused networks
docker system prune

# Aggressive: Remove ALL unused images, containers, volumes, networks
docker system prune -a --volumes
# ⚠️ WARNING: This removes all unused volumes including named volumes!

# Targeted cleanup:
docker container prune              # Remove stopped containers
docker image prune -a               # Remove all unused images
docker volume prune                 # Remove unused volumes
docker builder prune                # Remove build cache

# Remove images older than 24 hours
docker image prune -a --filter "until=24h"

# Remove specific old images
docker rmi $(docker images -q --filter "dangling=true")
```

### Prevention

- Run `docker system prune` on a schedule (cron job)
- Use multi-stage builds to minimize image sizes
- Set `--storage-opt size=10G` for container storage limits
- Monitor disk usage with alerts
- Configure log rotation for container logs:

```json
// /etc/docker/daemon.json
{
  "log-driver": "json-file",
  "log-opts": {
    "max-size": "10m",
    "max-file": "3"
  }
}
```

---

## Performance Issues

Containers running slowly or consuming excessive resources.

### Symptoms

```bash
# High CPU or memory usage
$ docker stats
CONTAINER ID   NAME    CPU %   MEM USAGE / LIMIT   MEM %   NET I/O       BLOCK I/O
abc123         myapp   98.5%   480MiB / 512MiB     93.8%   10MB / 5MB    50MB / 20MB
```

### Diagnostic Commands

```bash
# Live resource usage
docker stats
docker stats <container>

# Check resource limits
docker inspect <container> --format='CPU: {{.HostConfig.NanoCpus}} Memory: {{.HostConfig.Memory}}'

# Check processes inside container
docker exec <container> ps aux
docker exec <container> top -bn1

# Check I/O
docker exec <container> iostat 1 5    # If available
```

### Resolution

**Set Resource Limits:**

```bash
# Run with CPU and memory limits
docker run -d \
  --cpus="1.0" \
  --memory="512m" \
  --memory-swap="512m" \     # Same as memory = no swap
  myapp

# Update running container limits
docker update --cpus="2.0" --memory="1g" <container>
```

**CPU Issues:**

```bash
# CPU pinning (for consistent performance)
docker run -d --cpuset-cpus="0,1" myapp

# CPU shares (relative weight)
docker run -d --cpu-shares=512 myapp    # Default is 1024
```

**Memory Issues:**

```bash
# Check if OOM killer acted
docker inspect <container> --format='{{.State.OOMKilled}}'

# Check memory events
dmesg | grep -i oom | tail -10
```

### Prevention

- Always set memory limits in production
- Set CPU limits to prevent noisy-neighbor issues
- Monitor resource usage trends, not just current values
- Use `--memory-swap` equal to `--memory` to prevent swap (better: disable swap)
- Optimize application for container environments (JVM, Python, Node.js memory settings)

---

## General Diagnostic Commands

### Quick Health Check

```bash
# Docker daemon status
docker info
systemctl status docker

# All containers (including stopped)
docker ps -a

# Resource usage
docker stats --no-stream

# Disk usage
docker system df

# System events (live)
docker events
docker events --filter type=container --since "1h"
```

### Container Inspection

```bash
# Full container configuration
docker inspect <container>

# Specific fields
docker inspect <container> --format='{{.State.Status}}'
docker inspect <container> --format='{{.Config.Env}}'
docker inspect <container> --format='{{.NetworkSettings.IPAddress}}'
docker inspect <container> --format='{{json .Mounts}}' | jq

# Logs
docker logs <container> --tail 100
docker logs <container> --since "1h"
docker logs <container> -f              # Follow

# File system changes
docker diff <container>

# Execute commands inside
docker exec -it <container> /bin/sh
docker exec <container> cat /etc/hosts
```

### Image Inspection

```bash
# Image details
docker inspect <image>
docker history <image>              # Show layers

# Check image size
docker images <image>

# Image filesystem
docker run --rm -it <image> /bin/sh
```

---

## 🔗 Related Topics

- [🚀 DevOps / Docker](../../devops/docker/) — Docker concepts and operations
- [🧰 CLI / Docker](../../cli/docker/) — Docker command reference
- [🚀 Kubernetes Troubleshooting](../kubernetes/) — Kubernetes-specific container issues
- [🌐 Network Troubleshooting](../networking/) — Network debugging
- [🧠 Memory Troubleshooting](../memory/) — Memory and OOM issues
- [💻 CPU Troubleshooting](../cpu/) — CPU performance issues

---

> **Pro tip:** When a container won't start, check the exit code first. It tells you whether the problem is OOM (137), command not found (127), permission denied (126), or an application error (1). Then check the logs for details.
