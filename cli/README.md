# 🧰 CLI Command Center

> **Your daily command-line reference.** Quick access to every CLI tool, organized by category with essential commands and links to comprehensive guides.

---

## 📑 Table of Contents

- [Version Control](#-version-control)
- [Containers & Orchestration](#-containers--orchestration)
- [Text Processing](#-text-processing)
- [Networking](#-networking)
- [System Administration](#-system-administration)
- [Quick Lookup Table](#-quick-lookup-table)

---

## 🔀 Version Control

### [Git](git/) — Distributed Version Control

The essential tool for source code management, collaboration, and history tracking.

```bash
git status                          # Check working tree status
git log --oneline -10               # Last 10 commits, compact
git diff --staged                   # View staged changes before commit
git stash && git stash pop          # Stash and restore changes
git rebase -i HEAD~3                # Interactive rebase last 3 commits
git branch -a                       # List all branches (local + remote)
git cherry-pick <commit>            # Apply a specific commit to current branch
```

→ **[Full Git Reference](git/)**

---

## 🐳 Containers & Orchestration

### [Docker](docker/) — Container Runtime & Image Management

Build, ship, and run containers. The foundation of modern containerized applications.

```bash
docker ps                           # List running containers
docker logs -f <container>          # Follow container logs
docker exec -it <container> sh      # Shell into container
docker build -t myapp:v1 .          # Build image from Dockerfile
docker system prune -a              # Clean up unused resources
docker stats                        # Live resource usage
```

→ **[Full Docker Reference](docker/)**

### [kubectl](kubectl/) — Kubernetes Cluster Management

Interact with Kubernetes clusters: deploy, inspect, debug, and manage workloads.

```bash
kubectl get pods -A                 # All pods across namespaces
kubectl describe pod <pod>          # Pod details and events
kubectl logs <pod> -f --tail=100    # Follow pod logs
kubectl exec -it <pod> -- sh        # Shell into pod
kubectl top pods                    # Pod resource usage
kubectl get events --sort-by='.lastTimestamp'  # Recent cluster events
```

→ **[Full kubectl Reference](kubectl/)**

### [Helm](helm/) — Kubernetes Package Manager

Manage Kubernetes applications through Helm charts — install, upgrade, rollback, and version releases.

```bash
helm list -A                        # All releases across namespaces
helm install myapp ./chart          # Install a chart
helm upgrade myapp ./chart          # Upgrade a release
helm rollback myapp 1               # Rollback to revision 1
helm template myapp ./chart         # Render templates locally
helm get values myapp               # Get current release values
```

→ **[Full Helm Reference](helm/)**

---

## 📝 Text Processing

### [grep](grep/) — Pattern Search

Search files and streams for text patterns using regular expressions.

```bash
grep -rn "TODO" src/                # Recursive search with line numbers
grep -i "error" /var/log/syslog     # Case-insensitive search
grep -v "^#" config.yaml            # Exclude comment lines
grep -c "pattern" file.txt          # Count matches
grep -A3 -B1 "Exception" app.log   # Context around matches
```

→ **[Full grep Reference](grep/)**

### [sed](sed/) — Stream Editor

Transform text streams — find and replace, delete, insert, and modify files in-place.

```bash
sed 's/old/new/g' file.txt          # Replace all occurrences
sed -i 's/foo/bar/g' *.conf         # In-place edit across files
sed -n '10,20p' file.txt            # Print lines 10-20
sed '/^$/d' file.txt                # Delete empty lines
sed '1i\# Header' file.txt         # Insert line at top
```

→ **[Full sed Reference](sed/)**

### [awk](awk/) — Text Processing Language

Column-based processing, report generation, and data transformation.

```bash
awk '{print $1, $3}' file.txt      # Print columns 1 and 3
awk -F: '{print $1}' /etc/passwd   # Custom field separator
awk '/ERROR/ {count++} END {print count}' app.log  # Count errors
awk '{sum+=$2} END {print sum}' data.txt            # Sum a column
awk 'NR>=10 && NR<=20' file.txt    # Print line range
```

→ **[Full awk Reference](awk/)**

### [jq](jq/) — JSON Processing

Parse, filter, and transform JSON data from APIs, config files, and command output.

```bash
cat data.json | jq '.'             # Pretty-print JSON
jq '.items[].name' response.json   # Extract nested field
jq '.[] | select(.status=="active")' data.json  # Filter by value
kubectl get pods -o json | jq '.items[].metadata.name'  # K8s + jq
jq -r '.results[] | [.name, .age] | @csv' data.json     # JSON to CSV
```

→ **[Full jq Reference](jq/)**

---

## 🌐 Networking

### [curl](curl/) — HTTP Client & API Testing

Make HTTP requests, test APIs, download files, and debug connections.

```bash
curl -s https://api.example.com/health        # Simple GET
curl -X POST -H "Content-Type: application/json" \
  -d '{"key":"value"}' https://api.example.com  # POST JSON
curl -o file.tar.gz https://example.com/file   # Download file
curl -w "%{http_code}" -s -o /dev/null URL     # Get HTTP status code
curl -k --cert client.pem https://secure.api   # mTLS request
```

→ **[Full curl Reference](curl/)**

### [SSH](ssh/) — Secure Shell & Tunneling

Remote access, key management, tunneling, and file transfer.

```bash
ssh user@host                       # Basic connection
ssh -L 8080:localhost:80 user@host  # Local port forward
ssh -J jumphost user@target         # Jump through bastion
scp file.txt user@host:/path/       # Copy file to remote
ssh-keygen -t ed25519               # Generate key pair
```

→ **[Full SSH Reference](ssh/)**

### [Networking Tools](networking/) — DNS, Connectivity & Traffic

Diagnose DNS, test connectivity, inspect ports, and analyze traffic.

```bash
dig example.com                     # DNS lookup
ping -c 4 host                      # Test connectivity
ss -tlnp                            # Listening TCP ports
traceroute host                     # Trace network path
tcpdump -i eth0 port 443            # Capture HTTPS traffic
ip addr show                        # Show network interfaces
```

→ **[Full Networking Reference](networking/)**

---

## 🖥️ System Administration

### [Linux](linux/) — Core System Commands

File system navigation, process management, services, users, packages, and disks.

```bash
find / -name "*.log" -mtime -1      # Files modified in last day
ps aux --sort=-%mem | head           # Top memory consumers
systemctl status nginx               # Check service status
journalctl -u myservice -f --since "1h ago"  # Recent service logs
df -h && free -h                     # Disk and memory overview
chmod 755 script.sh && chown user:group file  # Permissions
```

→ **[Full Linux Reference](linux/)**

---

## 📋 Quick Lookup Table

| Tool | Primary Use | Category |
|------|------------|----------|
| [git](git/) | Version control, collaboration | Version Control |
| [docker](docker/) | Container lifecycle management | Containers |
| [kubectl](kubectl/) | Kubernetes cluster operations | Orchestration |
| [helm](helm/) | Kubernetes package management | Orchestration |
| [grep](grep/) | Text pattern searching | Text Processing |
| [sed](sed/) | Stream editing, find & replace | Text Processing |
| [awk](awk/) | Column processing, reporting | Text Processing |
| [jq](jq/) | JSON parsing & transformation | Text Processing |
| [curl](curl/) | HTTP requests, API testing | Networking |
| [ssh](ssh/) | Remote access, tunneling | Networking |
| [networking](networking/) | DNS, connectivity, traffic | Networking |
| [linux](linux/) | System administration | System Admin |

---

## 🔗 Related Topics

- [🚀 DevOps](../devops/) — Docker, Kubernetes, Helm, CI/CD in depth
- [🐛 Troubleshooting](../troubleshooting/) — Symptom-based debugging guides
- [🔐 Security](../security/) — TLS, certificates, authentication
- [📊 Observability](../observability/) — Monitoring, alerting, dashboards
- [💻 Languages](../languages/) — Programming language references

---

> **Tip:** Most CLI tools support `--help` and have detailed `man` pages. When in doubt: `man <command>` or `<command> --help`.
