# 🌐 curl — Comprehensive HTTP Client Reference

> **The universal HTTP/API testing tool.** Make requests, test APIs, debug connections, handle authentication, upload files, and diagnose network issues — all from the command line.

---

## 📑 Table of Contents

- [Concepts](#-concepts)
- [GET Requests](#-get-requests)
- [POST Requests](#-post-requests)
- [PUT & PATCH Requests](#-put--patch-requests)
- [DELETE Requests](#-delete-requests)
- [Headers](#-headers)
- [Authentication](#-authentication)
- [File Upload & Download](#-file-upload--download)
- [Redirects & Cookies](#-redirects--cookies)
- [Timing & Performance](#-timing--performance)
- [Verbose & Debug Output](#-verbose--debug-output)
- [Common API Testing Patterns](#-common-api-testing-patterns)
- [Troubleshooting](#-troubleshooting)
- [Production Tips](#-production-tips)
- [Related Topics](#-related-topics)

---

## 💡 Concepts

`curl` (Client URL) transfers data to or from a server using supported protocols (HTTP, HTTPS, FTP, SCP, SFTP, etc.). It's the go-to tool for testing APIs, debugging HTTP issues, and scripting HTTP interactions.

### Common Flags Quick Reference

| Flag | Description |
|------|-------------|
| `-X` | HTTP method (GET, POST, PUT, DELETE, PATCH) |
| `-H` | Add request header |
| `-d` | Request body data |
| `-o` | Write output to file |
| `-O` | Save with remote filename |
| `-s` | Silent mode (no progress) |
| `-S` | Show errors in silent mode |
| `-v` | Verbose output |
| `-k` | Allow insecure SSL |
| `-L` | Follow redirects |
| `-w` | Write-out format string |
| `-u` | User:password for auth |
| `-b` | Send cookies |
| `-c` | Save cookies to file |
| `-i` | Include response headers |
| `-I` | HEAD request (headers only) |
| `--connect-timeout` | Connection timeout (seconds) |
| `--max-time` | Total operation timeout |

---

## 📥 GET Requests

```bash
# Basic GET
curl https://api.example.com/users

# Silent (no progress bar)
curl -s https://api.example.com/users

# With response headers
curl -i https://api.example.com/users

# Headers only (HEAD request)
curl -I https://api.example.com/users

# Pretty-print JSON with jq
curl -s https://api.example.com/users | jq '.'

# With query parameters
curl -s "https://api.example.com/users?page=2&limit=10"

# URL-encode query params
curl -s -G https://api.example.com/search \
  --data-urlencode "q=hello world" \
  --data-urlencode "lang=en"

# Get HTTP status code only
curl -s -o /dev/null -w "%{http_code}" https://api.example.com/health
```

**Example output:**
```
$ curl -s https://api.example.com/users | jq '.'
{
  "users": [
    {
      "id": 1,
      "name": "Alice",
      "email": "alice@example.com"
    },
    {
      "id": 2,
      "name": "Bob",
      "email": "bob@example.com"
    }
  ],
  "total": 2
}
```

---

## 📤 POST Requests

```bash
# POST JSON
curl -X POST https://api.example.com/users \
  -H "Content-Type: application/json" \
  -d '{"name": "Alice", "email": "alice@example.com"}'

# POST from file
curl -X POST https://api.example.com/users \
  -H "Content-Type: application/json" \
  -d @payload.json

# POST form data
curl -X POST https://api.example.com/login \
  -d "username=admin&password=secret"

# POST form data (explicit)
curl -X POST https://api.example.com/login \
  --data-urlencode "username=admin" \
  --data-urlencode "password=p@ss w0rd"

# POST multipart form (file upload)
curl -X POST https://api.example.com/upload \
  -F "file=@/path/to/document.pdf" \
  -F "description=My document"
```

---

## ✏️ PUT & PATCH Requests

```bash
# PUT (full resource replacement)
curl -X PUT https://api.example.com/users/1 \
  -H "Content-Type: application/json" \
  -d '{"name": "Alice Updated", "email": "alice.new@example.com"}'

# PATCH (partial update)
curl -X PATCH https://api.example.com/users/1 \
  -H "Content-Type: application/json" \
  -d '{"email": "newemail@example.com"}'
```

---

## 🗑️ DELETE Requests

```bash
# Simple DELETE
curl -X DELETE https://api.example.com/users/1

# DELETE with auth
curl -X DELETE https://api.example.com/users/1 \
  -H "Authorization: Bearer <token>"

# DELETE with body (some APIs require it)
curl -X DELETE https://api.example.com/bulk \
  -H "Content-Type: application/json" \
  -d '{"ids": [1, 2, 3]}'
```

---

## 📋 Headers

```bash
# Set custom headers
curl -H "Content-Type: application/json" \
     -H "Accept: application/json" \
     -H "X-Request-ID: abc-123" \
     https://api.example.com/data

# View request and response headers
curl -v https://api.example.com/data 2>&1 | grep -E '^[<>]'

# Common headers
curl -H "Content-Type: application/json" URL        # JSON body
curl -H "Accept: text/xml" URL                      # Request XML response
curl -H "Cache-Control: no-cache" URL                # Skip cache
curl -H "User-Agent: MyApp/1.0" URL                 # Custom user agent
curl -H "X-Forwarded-For: 10.0.0.1" URL             # Forwarded IP
```

---

## 🔐 Authentication

### Basic Auth

```bash
curl -u username:password https://api.example.com/secure
curl -u "admin:p@ssw0rd" https://api.example.com/secure

# Prompt for password (don't expose in shell history)
curl -u username https://api.example.com/secure
```

### Bearer Token (OAuth/JWT)

```bash
curl -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIs..." \
  https://api.example.com/protected

# Token from variable
TOKEN=$(curl -s -X POST https://auth.example.com/token \
  -d "client_id=myapp&client_secret=secret&grant_type=client_credentials" \
  | jq -r '.access_token')

curl -H "Authorization: Bearer $TOKEN" https://api.example.com/data
```

### Client Certificates (mTLS)

```bash
# Client certificate authentication
curl --cert client.pem --key client-key.pem https://secure.api.com

# With CA certificate
curl --cacert ca.pem --cert client.pem --key client-key.pem \
  https://secure.api.com

# PFX/P12 format
curl --cert-type P12 --cert client.p12:password https://secure.api.com
```

### API Key

```bash
# In header
curl -H "X-API-Key: your-api-key-here" https://api.example.com/data

# In query parameter
curl "https://api.example.com/data?api_key=your-key"
```

---

## 📁 File Upload & Download

### Download

```bash
# Save with custom name
curl -o output.tar.gz https://example.com/file.tar.gz

# Save with remote filename
curl -O https://example.com/file.tar.gz

# Resume interrupted download
curl -C - -O https://example.com/large-file.iso

# Download silently
curl -sLO https://example.com/script.sh

# Download multiple files
curl -O https://example.com/file1.txt -O https://example.com/file2.txt

# Limit download speed
curl --limit-rate 1M -O https://example.com/large-file.iso
```

### Upload

```bash
# Upload file
curl -X POST -F "file=@photo.jpg" https://api.example.com/upload

# Multiple files
curl -X POST \
  -F "file1=@photo1.jpg" \
  -F "file2=@photo2.jpg" \
  https://api.example.com/upload

# Upload with PUT
curl -T file.txt https://api.example.com/files/file.txt

# Upload binary data
curl -X POST --data-binary @image.png \
  -H "Content-Type: image/png" \
  https://api.example.com/images
```

---

## 🔄 Redirects & Cookies

### Follow Redirects

```bash
curl -L https://example.com           # Follow 3xx redirects
curl -L --max-redirs 5 https://short.url/abc  # Limit redirect count
```

### Cookies

```bash
# Send cookie
curl -b "session=abc123" https://api.example.com/dashboard

# Save cookies from response
curl -c cookies.txt https://api.example.com/login -d "user=admin&pass=secret"

# Use saved cookies
curl -b cookies.txt https://api.example.com/dashboard

# Save and send in one session
curl -c cookies.txt -b cookies.txt https://api.example.com/page
```

---

## ⏱️ Timing & Performance

### Write-Out Variables

```bash
# HTTP status code
curl -s -o /dev/null -w "%{http_code}\n" https://api.example.com

# Full timing breakdown
curl -s -o /dev/null -w "\
    DNS Lookup:  %{time_namelookup}s\n\
    TCP Connect: %{time_connect}s\n\
    TLS Handshake: %{time_appconnect}s\n\
    Start Transfer: %{time_starttransfer}s\n\
    Total Time:    %{time_total}s\n\
    Download Size: %{size_download} bytes\n\
    HTTP Code:     %{http_code}\n" \
  https://api.example.com
```

**Example output:**
```
    DNS Lookup:  0.012s
    TCP Connect: 0.045s
    TLS Handshake: 0.123s
    Start Transfer: 0.234s
    Total Time:    0.456s
    Download Size: 15234 bytes
    HTTP Code:     200
```

### Timeouts

```bash
curl --connect-timeout 5 https://api.example.com     # Connection timeout
curl --max-time 30 https://api.example.com            # Total request timeout
curl --connect-timeout 5 --max-time 30 https://api.example.com  # Both
curl --retry 3 --retry-delay 2 https://api.example.com  # Retry on failure
```

---

## 🔎 Verbose & Debug Output

```bash
# Verbose (shows request/response headers)
curl -v https://api.example.com

# Even more verbose (includes TLS handshake)
curl -vvv https://api.example.com

# Trace full network data
curl --trace trace.log https://api.example.com
curl --trace-ascii trace.txt https://api.example.com
```

**Example verbose output:**
```
$ curl -v https://api.example.com/health
*   Trying 93.184.216.34:443...
* Connected to api.example.com (93.184.216.34) port 443 (#0)
* SSL connection using TLSv1.3 / TLS_AES_256_GCM_SHA384
> GET /health HTTP/2
> Host: api.example.com
> User-Agent: curl/8.1.2
> Accept: */*
>
< HTTP/2 200
< content-type: application/json
< date: Thu, 15 Aug 2024 10:30:00 GMT
<
{"status":"healthy","version":"1.2.3"}
```

---

## 🧪 Common API Testing Patterns

### Health Check Script

```bash
#!/bin/bash
STATUS=$(curl -s -o /dev/null -w "%{http_code}" https://api.example.com/health)
if [ "$STATUS" -eq 200 ]; then
    echo "✅ Service is healthy"
else
    echo "❌ Service returned HTTP $STATUS"
    exit 1
fi
```

### CRUD Operations

```bash
# Create
curl -s -X POST https://api.example.com/items \
  -H "Content-Type: application/json" \
  -d '{"name":"Widget","price":9.99}' | jq '.'

# Read
curl -s https://api.example.com/items/1 | jq '.'

# Update
curl -s -X PUT https://api.example.com/items/1 \
  -H "Content-Type: application/json" \
  -d '{"name":"Widget Pro","price":19.99}' | jq '.'

# Delete
curl -s -X DELETE https://api.example.com/items/1 -w "\n%{http_code}\n"
```

### Test Multiple Endpoints

```bash
#!/bin/bash
ENDPOINTS=(
    "https://api.example.com/health"
    "https://api.example.com/users"
    "https://api.example.com/products"
)

for url in "${ENDPOINTS[@]}"; do
    STATUS=$(curl -s -o /dev/null -w "%{http_code}" --max-time 10 "$url")
    printf "%-50s %s\n" "$url" "$STATUS"
done
```

### Webhook Testing

```bash
# Send a test webhook payload
curl -X POST https://webhook.example.com/callback \
  -H "Content-Type: application/json" \
  -H "X-Webhook-Signature: sha256=..." \
  -d '{
    "event": "order.completed",
    "data": {
      "order_id": "12345",
      "amount": 99.99
    }
  }'
```

---

## 🐛 Troubleshooting

### SSL/TLS Errors

```bash
# Skip certificate verification (INSECURE — dev/testing only)
curl -k https://self-signed.example.com

# Specify CA certificate
curl --cacert /path/to/ca-cert.pem https://api.example.com

# Check SSL certificate details
curl -vI https://api.example.com 2>&1 | grep -A5 "Server certificate"

# Force specific TLS version
curl --tlsv1.2 https://api.example.com
curl --tlsv1.3 https://api.example.com
```

### Common Issues

| Problem | Diagnosis | Solution |
|---------|-----------|----------|
| `curl: (6) Could not resolve host` | DNS failure | Check DNS, try IP directly, verify hostname |
| `curl: (7) Failed to connect` | Server unreachable | Check port, firewall, security groups |
| `curl: (28) Operation timed out` | Timeout | Increase `--max-time`, check network |
| `curl: (35) SSL connect error` | TLS handshake failure | Check cert, try `-k` for testing, verify TLS version |
| `curl: (51) SSL peer certificate` | Certificate mismatch | Verify hostname matches cert, use `--cacert` |
| `curl: (56) Recv failure` | Connection reset | Server closed connection, check server logs |
| `curl: (60) SSL certificate problem` | Untrusted CA | Use `--cacert` or update CA bundle |
| HTTP 401 | Authentication failed | Check credentials, token expiry |
| HTTP 403 | Forbidden | Check permissions, API key, IP allowlist |
| HTTP 429 | Rate limited | Add delay between requests, check rate limit headers |

### Debug Network Issues

```bash
# Resolve but don't connect
curl -v --resolve api.example.com:443:10.0.0.1 https://api.example.com

# Use specific DNS server (via --resolve)
curl --resolve "api.example.com:443:$(dig +short api.example.com @8.8.8.8)" \
  https://api.example.com

# Connect to different IP (bypass DNS)
curl --connect-to api.example.com:443:10.0.0.1:443 https://api.example.com
```

---

## 🏭 Production Tips

- **Never put credentials in command line** in shared environments — use `-u user` (prompts) or env vars
- **Always use `--max-time`** in scripts to prevent hanging
- **Use `-sS`** (silent + show errors) in scripts for clean output
- **Store curl format strings** in `~/.curlrc` or a file for reuse
- **Use `--retry` with `--retry-delay`** for resilient scripts
- **Pipe to `jq`** for readable JSON output
- **Set `--connect-timeout`** separately from `--max-time` for better error diagnosis

### ~/.curlrc (Default Options)

```
# ~/.curlrc
--max-time 30
--connect-timeout 10
--retry 2
--silent
--show-error
--location
```

---

## 🔗 Related Topics

- [🧰 CLI Command Center](../) — All CLI tool references
- [jq](../jq/) — JSON processing for API responses
- [SSH](../ssh/) — Secure tunneling for API access
- [Networking](../networking/) — DNS, connectivity debugging
- [🔐 Security](../../security/) — TLS, certificates, OAuth, JWT

---

> **Tip:** Use `curl -w "%{http_code}" -s -o /dev/null URL` as a quick health check in scripts. Combine with `jq` for powerful API testing one-liners.
