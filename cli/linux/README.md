# 🐧 Linux — Comprehensive Command Reference

> **Essential Linux commands for developers and sysadmins.** File system operations, process management, text processing, services, users, packages, and disk management.

---

## 📑 Table of Contents

- [File System Navigation](#-file-system-navigation)
- [File Operations](#-file-operations)
- [File Permissions](#-file-permissions)
- [Text Processing](#-text-processing)
- [Process Management](#-process-management)
- [System Information](#-system-information)
- [Service Management](#-service-management)
- [User Management](#-user-management)
- [Package Management](#-package-management)
- [Disk & Storage](#-disk--storage)
- [Troubleshooting](#-troubleshooting)
- [Production Tips](#-production-tips)
- [Related Topics](#-related-topics)

---

## 📂 File System Navigation

### ls — List Directory Contents

```bash
ls                          # Basic listing
ls -la                      # Long format, including hidden files
ls -lh                      # Human-readable sizes
ls -lt                      # Sort by modification time (newest first)
ls -ltr                     # Sort by time (oldest first)
ls -lS                      # Sort by file size (largest first)
ls -R                       # Recursive listing
ls -1                       # One file per line
```

**Example output:**
```
$ ls -lah
total 48K
drwxr-xr-x  5 user group 4.0K Aug 15 10:30 .
drwxr-xr-x 12 user group 4.0K Aug 14 09:00 ..
-rw-r--r--  1 user group  256 Aug 15 10:30 README.md
drwxr-xr-x  3 user group 4.0K Aug 15 10:25 src
-rwxr-xr-x  1 user group 1.2K Aug 14 15:00 deploy.sh
```

| Column | Meaning |
|--------|---------|
| `drwxr-xr-x` | Type + permissions (d=directory, -=file, l=link) |
| `5` | Hard link count |
| `user` | Owner |
| `group` | Group |
| `4.0K` | File size |
| `Aug 15 10:30` | Last modification time |

### cd — Change Directory

```bash
cd /var/log                 # Absolute path
cd src/                     # Relative path
cd ~                        # Home directory
cd -                        # Previous directory
cd ..                       # Parent directory
cd ../..                    # Two levels up
```

### find — Search Files

```bash
# By name
find / -name "*.log"                        # Find all .log files
find . -name "*.py" -type f                 # Python files only
find . -iname "readme*"                     # Case-insensitive

# By type
find . -type f                              # Files only
find . -type d                              # Directories only
find . -type l                              # Symbolic links

# By time
find . -mtime -1                            # Modified in last 24 hours
find . -mtime +30                           # Modified more than 30 days ago
find . -mmin -60                            # Modified in last 60 minutes
find . -newer reference.txt                 # Newer than reference file

# By size
find . -size +100M                          # Larger than 100MB
find . -size -1k                            # Smaller than 1KB
find / -size +1G -type f                    # Files larger than 1GB

# By permissions
find . -perm 644                            # Exact permissions
find . -perm -u+x                           # User executable
find / -perm -4000 -type f                  # SUID files (security audit)

# With actions
find . -name "*.tmp" -delete                # Delete matching files
find . -name "*.sh" -exec chmod +x {} \;    # Execute command on results
find . -name "*.log" -exec grep -l "ERROR" {} \;  # Grep in found files
find . -name "*.py" | xargs wc -l           # Count lines in Python files

# Combining conditions
find . -name "*.log" -mtime +7 -size +10M   # Old, large log files
find . \( -name "*.py" -o -name "*.js" \)   # Python OR JavaScript files
find . -name "*.bak" ! -path "*/vendor/*"   # Exclude vendor directory
```

### du — Disk Usage

```bash
du -sh *                    # Size of each item in current directory
du -sh /var/log             # Total size of directory
du -ah . | sort -rh | head -20  # Top 20 largest files/dirs
du -sh --max-depth=1 /      # Size of top-level directories
du -sh --exclude="*.log" .  # Exclude log files
```

### df — Filesystem Disk Space

```bash
df -h                       # Human-readable filesystem usage
df -hT                      # Include filesystem type
df -i                       # Inode usage (important when running out of inodes)
```

**Example output:**
```
$ df -h
Filesystem      Size  Used Avail Use% Mounted on
/dev/sda1       100G   45G   55G  45% /
/dev/sdb1       500G  200G  300G  40% /data
tmpfs           16G     0   16G   0% /dev/shm
```

---

## 📁 File Operations

### Copy, Move, Remove

```bash
# Copy
cp file.txt backup.txt              # Copy file
cp -r src/ dest/                    # Copy directory recursively
cp -p file.txt backup.txt          # Preserve permissions and timestamps
cp -a src/ dest/                    # Archive (preserve everything)
cp -i file.txt dest/               # Interactive (prompt before overwrite)

# Move / Rename
mv file.txt new-name.txt           # Rename
mv file.txt /path/to/dest/         # Move
mv -i src dest                     # Interactive (prompt before overwrite)

# Remove
rm file.txt                        # Remove file
rm -r directory/                   # Remove directory recursively
rm -rf directory/                  # Force remove (no prompts)
rm -i file.txt                     # Interactive (confirm before delete)
rmdir empty-dir/                   # Remove empty directory only
```

> **Warning:** `rm -rf /` or `rm -rf *` in the wrong directory can be catastrophic. Always double-check your current directory and the target path.

### Create

```bash
mkdir mydir                         # Create directory
mkdir -p path/to/nested/dir         # Create nested directories
touch file.txt                      # Create empty file / update timestamp
ln -s /path/to/target link-name    # Create symbolic link
ln /path/to/target hard-link       # Create hard link
```

### File Content

```bash
cat file.txt                        # Display entire file
less file.txt                       # Paginated viewer (q to quit)
more file.txt                       # Basic paginated viewer
head -20 file.txt                   # First 20 lines
tail -20 file.txt                   # Last 20 lines
tail -f /var/log/syslog             # Follow file (live updates)
tail -f -n 100 app.log             # Follow, starting from last 100 lines
```

---

## 🔐 File Permissions

### chmod — Change Permissions

```bash
# Symbolic mode
chmod +x script.sh                  # Add execute for everyone
chmod u+rwx file.txt                # Owner: read+write+execute
chmod go-w file.txt                 # Remove write from group and others
chmod u=rwx,g=rx,o=r file.txt      # Set exact permissions

# Numeric mode
chmod 755 script.sh                 # rwxr-xr-x
chmod 644 config.yml                # rw-r--r--
chmod 600 ~/.ssh/id_rsa             # rw------- (private key)
chmod 700 ~/.ssh                    # rwx------ (ssh directory)

# Recursive
chmod -R 755 directory/
```

| Digit | Permission | Meaning |
|-------|-----------|---------|
| 7 | rwx | Read + Write + Execute |
| 6 | rw- | Read + Write |
| 5 | r-x | Read + Execute |
| 4 | r-- | Read only |
| 3 | -wx | Write + Execute |
| 2 | -w- | Write only |
| 1 | --x | Execute only |
| 0 | --- | No permissions |

### chown — Change Ownership

```bash
chown user file.txt                 # Change owner
chown user:group file.txt           # Change owner and group
chown -R user:group directory/      # Recursive
chown --reference=ref.txt target.txt  # Copy ownership from reference
```

---

## 📝 Text Processing

> **See also:** [grep](../grep/), [sed](../sed/), [awk](../awk/) for comprehensive text processing references.

### Basic Commands

```bash
# Word, line, character counts
wc file.txt                         # Lines, words, characters
wc -l file.txt                      # Line count only
wc -w file.txt                      # Word count only
find . -name "*.py" | xargs wc -l  # Count lines across files

# Sort
sort file.txt                       # Alphabetical sort
sort -n file.txt                    # Numeric sort
sort -r file.txt                    # Reverse sort
sort -u file.txt                    # Sort and remove duplicates
sort -t: -k3 -n /etc/passwd        # Sort by 3rd field (: delimiter)
sort -h file.txt                    # Human-readable numbers (1K, 2M)

# Unique
uniq file.txt                       # Remove consecutive duplicates
uniq -c file.txt                    # Count occurrences
uniq -d file.txt                    # Show only duplicates
sort file.txt | uniq -c | sort -rn  # Frequency count (most common first)

# Transform
tr 'a-z' 'A-Z' < file.txt          # Uppercase
tr -d '\r' < file.txt               # Remove carriage returns
tr -s ' ' < file.txt               # Squeeze repeated spaces
tr ',' '\t' < file.csv              # CSV to TSV

# Cut — extract columns
cut -d: -f1 /etc/passwd             # First field, : delimiter
cut -d, -f1,3 data.csv             # Fields 1 and 3 from CSV
cut -c1-10 file.txt                 # First 10 characters per line

# Paste — merge files
paste file1.txt file2.txt           # Side by side (tab separated)
paste -d, file1.txt file2.txt       # Side by side (comma separated)
```

### Comparing Files

```bash
diff file1.txt file2.txt            # Show differences
diff -u file1.txt file2.txt         # Unified format (like git diff)
diff -y file1.txt file2.txt         # Side-by-side
diff -r dir1/ dir2/                 # Recursive directory diff
colordiff file1.txt file2.txt       # Colored diff (if installed)
comm file1.txt file2.txt            # Compare sorted files (common/unique)
```

---

## ⚙️ Process Management

### ps — Process Status

```bash
ps                                  # Current user's processes
ps aux                              # All processes (BSD style)
ps -ef                              # All processes (System V style)
ps aux --sort=-%mem | head          # Top memory consumers
ps aux --sort=-%cpu | head          # Top CPU consumers
ps -ef | grep nginx                 # Find specific process
ps -p <pid> -o pid,ppid,cmd,%mem,%cpu  # Specific process details
ps --forest                         # Process tree
```

**Example output:**
```
$ ps aux --sort=-%mem | head -5
USER       PID %CPU %MEM    VSZ   RSS TTY      STAT START   TIME COMMAND
postgres   512  2.1 15.3  524288 251412 ?      Ss   Aug10  45:23 postgres
java      1024  5.2 12.1  2048576 198432 ?     Sl   Aug12  30:15 java -jar app.jar
nginx      256  0.5  3.2  128512  52488 ?      S    Aug10   8:45 nginx: worker
```

### top / htop — Interactive Process Monitor

```bash
top                                 # Interactive process monitor
top -bn1 | head -20                 # Non-interactive snapshot (for scripts)
htop                                # Enhanced interactive monitor (if installed)
```

**top key commands:**

| Key | Action |
|-----|--------|
| `M` | Sort by memory |
| `P` | Sort by CPU |
| `T` | Sort by time |
| `k` | Kill a process |
| `u` | Filter by user |
| `1` | Toggle per-CPU view |
| `q` | Quit |

### kill — Signal Processes

```bash
kill <pid>                          # Send SIGTERM (graceful)
kill -9 <pid>                       # Send SIGKILL (force)
kill -HUP <pid>                     # Send SIGHUP (reload config)
killall nginx                       # Kill by process name
pkill -f "python app.py"           # Kill by command pattern

# Common signals
kill -l                             # List all signals
```

| Signal | Number | Description |
|--------|--------|-------------|
| `SIGHUP` | 1 | Hangup / reload |
| `SIGINT` | 2 | Interrupt (Ctrl+C) |
| `SIGTERM` | 15 | Graceful termination (default) |
| `SIGKILL` | 9 | Force kill (cannot be caught) |
| `SIGUSR1/2` | 10/12 | User-defined |

### Job Control

```bash
command &                           # Run in background
Ctrl+Z                              # Suspend current process
bg                                  # Resume suspended process in background
fg                                  # Bring background process to foreground
jobs                                # List background jobs
fg %1                               # Bring job #1 to foreground
kill %1                             # Kill job #1
nohup command &                     # Run immune to hangup signal
disown %1                           # Detach job from terminal
```

### nice / renice — Process Priority

```bash
nice -n 10 command                  # Run with lower priority (10)
nice -n -5 command                  # Run with higher priority (-5, needs root)
renice -n 10 -p <pid>              # Change priority of running process
```

---

## 📊 System Information

```bash
# OS and kernel
uname -a                            # All system info
uname -r                            # Kernel version
cat /etc/os-release                 # Distribution info
hostname                            # System hostname
hostnamectl                         # Detailed host info

# Uptime and load
uptime                              # Uptime and load averages
w                                   # Who's logged in + what they're doing
```

**Example output:**
```
$ uptime
 14:23:15 up 45 days,  3:12,  2 users,  load average: 0.52, 0.48, 0.51
```

> **Load average:** 3 numbers = 1, 5, 15 minute averages. Values should be below the number of CPU cores.

### Memory

```bash
free -h                             # Memory usage (human-readable)
free -m                             # Memory in MB
cat /proc/meminfo                   # Detailed memory info
```

**Example output:**
```
$ free -h
              total        used        free      shared  buff/cache   available
Mem:           16Gi       8.2Gi       1.5Gi       256Mi       6.3Gi       7.2Gi
Swap:          4.0Gi       512Mi       3.5Gi
```

### CPU and Performance

```bash
nproc                               # Number of CPU cores
lscpu                               # Detailed CPU info

# vmstat — virtual memory statistics
vmstat 1 5                          # Print every 1 second, 5 times
```

**Example output:**
```
$ vmstat 1 3
procs -----------memory---------- ---swap-- -----io---- -system-- ------cpu-----
 r  b   swpd   free   buff  cache   si   so    bi    bo   in   cs us sy id wa st
 1  0 524288 1548672 132096 6619136  0    0    12    45   1024  2048 12  3 84  1  0
 0  0 524288 1547900 132096 6619200  0    0     0    16    980  1900 10  2 88  0  0
```

| Column | Meaning |
|--------|---------|
| `r` | Runnable processes |
| `b` | Blocked processes |
| `us` | User CPU % |
| `sy` | System CPU % |
| `id` | Idle CPU % |
| `wa` | I/O wait % |

```bash
# iostat — I/O statistics
iostat -x 1 3                       # Extended I/O stats every 1 second
```

---

## 🔧 Service Management

### systemctl — systemd Service Manager

```bash
# Service status
systemctl status nginx              # Detailed status
systemctl is-active nginx           # Quick check (active/inactive)
systemctl is-enabled nginx          # Check if starts on boot

# Start / Stop / Restart
systemctl start nginx
systemctl stop nginx
systemctl restart nginx             # Stop and start
systemctl reload nginx              # Reload config without stopping

# Enable / Disable (boot behavior)
systemctl enable nginx              # Start on boot
systemctl disable nginx             # Don't start on boot
systemctl enable --now nginx        # Enable and start immediately

# List services
systemctl list-units --type=service             # Active services
systemctl list-units --type=service --all       # All services
systemctl list-unit-files --type=service        # All service files

# Check failed services
systemctl --failed
```

### journalctl — systemd Logs

```bash
journalctl                          # All logs
journalctl -u nginx                 # Logs for specific service
journalctl -u nginx -f              # Follow service logs
journalctl -u nginx --since "1 hour ago"
journalctl -u nginx --since "2024-08-01" --until "2024-08-02"
journalctl -u nginx -n 100          # Last 100 lines
journalctl -p err                   # Error priority and above
journalctl -p warning               # Warning priority and above
journalctl -b                       # Current boot logs
journalctl -b -1                    # Previous boot logs
journalctl --disk-usage             # Log storage size
journalctl --vacuum-size=500M       # Limit log storage to 500MB
```

| Priority | Level |
|----------|-------|
| 0 | emerg |
| 1 | alert |
| 2 | crit |
| 3 | err |
| 4 | warning |
| 5 | notice |
| 6 | info |
| 7 | debug |

---

## 👤 User Management

```bash
# Current user info
whoami                              # Current username
id                                  # User ID, group ID, groups
groups                              # Group membership

# User management (requires root)
useradd -m -s /bin/bash newuser     # Create user with home dir and shell
usermod -aG docker newuser          # Add user to group
usermod -s /bin/zsh newuser         # Change shell
userdel -r olduser                  # Delete user and home directory
passwd newuser                      # Set/change password

# Group management
groupadd developers                 # Create group
groupdel developers                 # Delete group

# Switch user
su - otheruser                      # Switch user (login shell)
sudo -u otheruser command           # Run command as another user
sudo -i                             # Root shell

# Who's logged in
who                                 # Current users
w                                   # Users and activity
last                                # Login history
last -n 10                          # Last 10 logins
```

---

## 📦 Package Management

### APT (Debian/Ubuntu)

```bash
apt update                          # Update package list
apt upgrade                         # Upgrade all packages
apt install nginx                   # Install package
apt remove nginx                    # Remove package (keep config)
apt purge nginx                     # Remove package + config
apt autoremove                      # Remove unused dependencies
apt search nginx                    # Search for packages
apt show nginx                      # Package details
apt list --installed                # List installed packages
apt list --upgradable               # List upgradable packages
```

### YUM/DNF (RHEL/CentOS/Fedora)

```bash
dnf check-update                    # Check for updates
dnf upgrade                         # Upgrade all packages
dnf install nginx                   # Install package
dnf remove nginx                    # Remove package
dnf search nginx                    # Search packages
dnf info nginx                      # Package details
dnf list installed                  # List installed packages
dnf provides "*/nginx"              # Find which package provides a file
dnf history                         # Transaction history
dnf history undo <id>               # Undo a transaction
```

---

## 💽 Disk & Storage

### Block Devices

```bash
lsblk                               # List block devices (tree view)
lsblk -f                            # With filesystem info
blkid                                # Block device IDs and types
fdisk -l                             # List partitions (requires root)
```

**Example output:**
```
$ lsblk
NAME   MAJ:MIN RM   SIZE RO TYPE MOUNTPOINT
sda      8:0    0   100G  0 disk
├─sda1   8:1    0    99G  0 part /
└─sda2   8:2    0     1G  0 part [SWAP]
sdb      8:16   0   500G  0 disk
└─sdb1   8:17   0   500G  0 part /data
```

### Mount Operations

```bash
mount /dev/sdb1 /mnt/data           # Mount device
mount -t nfs server:/share /mnt/nfs # Mount NFS share
umount /mnt/data                    # Unmount
mount | grep sda                    # Show mounts for device
cat /etc/fstab                      # Permanent mount config
```

### Filesystem Operations

```bash
mkfs.ext4 /dev/sdb1                 # Create ext4 filesystem
mkfs.xfs /dev/sdb1                  # Create XFS filesystem
fsck /dev/sdb1                      # Filesystem check (unmount first!)
tune2fs -l /dev/sda1                # Filesystem info
resize2fs /dev/sda1                 # Resize ext filesystem
```

---

## 🐛 Troubleshooting

### Quick System Health Check

```bash
# One-line system overview
uptime && free -h && df -h

# Top resource consumers
ps aux --sort=-%cpu | head -10      # CPU
ps aux --sort=-%mem | head -10      # Memory

# Disk space issues
du -sh /* 2>/dev/null | sort -rh | head  # Largest top-level dirs
find / -xdev -size +100M -type f 2>/dev/null  # Large files

# Open files / file descriptors
lsof | wc -l                        # Total open files
lsof -u user                        # Open files by user
lsof +D /var/log                    # Files open in directory
ulimit -n                           # Max open files limit
```

### Common Issues

| Problem | Diagnosis | Solution |
|---------|-----------|----------|
| Disk full | `df -h`, `du -sh /*` | Find and remove large files, clear logs |
| Out of inodes | `df -i` | Remove many small files, find with `find / -xdev -printf '%h\n' \| sort \| uniq -c \| sort -rn \| head` |
| High CPU | `top`, `ps aux --sort=-%cpu` | Identify process, optimize or kill |
| High memory | `free -h`, `ps aux --sort=-%mem` | Check for leaks, increase swap, kill process |
| Service won't start | `systemctl status svc`, `journalctl -u svc` | Check config, permissions, port conflicts |
| Permission denied | `ls -la`, `namei -l /path` | Fix ownership/permissions, check SELinux |
| "Too many open files" | `ulimit -n`, `lsof \| wc -l` | Increase limits in `/etc/security/limits.conf` |

---

## 🏭 Production Tips

### Essential Aliases

```bash
alias ll='ls -lah'
alias ..='cd ..'
alias ...='cd ../..'
alias ports='ss -tlnp'
alias meminfo='free -h'
alias diskinfo='df -h'
alias psmem='ps aux --sort=-%mem | head -20'
alias pscpu='ps aux --sort=-%cpu | head -20'
```

### Security Hardening

- Disable root SSH login (`PermitRootLogin no` in `/etc/ssh/sshd_config`)
- Use SSH keys instead of passwords
- Keep system updated (`apt upgrade` / `dnf upgrade`)
- Audit SUID files: `find / -perm -4000 -type f 2>/dev/null`
- Monitor auth logs: `journalctl -u sshd` or `/var/log/auth.log`

### Log Rotation

Most systems use `logrotate`. Check config at `/etc/logrotate.conf` and `/etc/logrotate.d/`.

```bash
logrotate --debug /etc/logrotate.conf    # Test configuration
logrotate -f /etc/logrotate.d/myapp      # Force rotation
```

---

## 🔗 Related Topics

- [🧰 CLI Command Center](../) — All CLI tool references
- [grep](../grep/) — Text pattern searching
- [sed](../sed/) — Stream editing
- [awk](../awk/) — Text processing and reporting
- [SSH](../ssh/) — Remote access
- [Networking](../networking/) — Network tools
- [🐛 Troubleshooting — CPU](../../troubleshooting/cpu/) — CPU debugging
- [🐛 Troubleshooting — Memory](../../troubleshooting/memory/) — Memory debugging

---

> **Tip:** Combine commands with pipes (`|`) for powerful one-liners. Example: `find . -name "*.log" -mtime +7 | xargs rm -f` removes log files older than 7 days.
