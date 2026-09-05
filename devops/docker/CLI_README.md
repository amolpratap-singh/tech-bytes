# 🐳 Docker — Comprehensive Reference

> **Build, ship, and run containers.** From images and containers to networking, storage, Compose, and production best practices — the complete Docker CLI reference.

---

## 🖼️ Images

### Pulling Images

```bash
docker pull nginx                   # Pull latest from Docker Hub
docker pull redis:4.0               # Pull specific version/tag
docker pull busybox                 # Lightweight Linux with shell commands
```

> If `docker run` is used and the image is not found locally, Docker automatically pulls it from Docker Hub.

### Listing Images

```bash
docker images                       # List all local images
docker images -q                    # Only image IDs
```

**Example output:**
```
$ docker images
REPOSITORY          TAG       IMAGE ID       CREATED        SIZE
nginx               latest    a6bd71f48f68   2 days ago     187MB
redis               alpine    3900abf41552   5 days ago     37.8MB
postgres            14        ceccf204b270   1 week ago     379MB
```

### Image History

Shows the layered build architecture — each instruction in a Dockerfile creates a new layer with its own ID and size.

```bash
docker history nginx
docker history new_myansible
```

**Example output:**
```
$ docker history new_myansible
IMAGE               CREATED       CREATED BY                                      SIZE      COMMENT
a0967c827609        2 days ago    /bin/bash                                       1.52kB    test image
bf74bd9701cb        2 days ago    /bin/bash                                       983B
c0814b11c786        2 weeks ago   /bin/sh -c pip3 install paramiko && pip3...     78.9MB
86ce23341144        2 weeks ago   /bin/sh -c apt-get update && apt-get -y inst…   120MB
d44ab0648e76        2 weeks ago   /bin/sh -c pip3 install ansible==2.9.9 &…      152MB
fa38a11831cb        2 weeks ago   /bin/sh -c apt-get update && apt-get ins…       57.6MB
```
