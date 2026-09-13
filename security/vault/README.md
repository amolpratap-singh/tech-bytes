# 🔐 HashiCorp Vault

> **Vault is the industry standard for secrets management.** It provides a unified interface to manage secrets, encrypt data in transit, and control access to sensitive information. Vault handles the lifecycle of secrets — from generation to rotation to revocation — so your applications don't have to.

---

## 📑 Table of Contents

- [What is Vault?](#-what-is-vault)
- [Architecture](#-architecture)
- [Secret Engines](#-secret-engines)
- [Authentication Methods](#-authentication-methods)
- [Policies](#-policies)
- [CLI Reference](#-cli-reference)
- [Dynamic Secrets](#-dynamic-secrets)
- [Transit Engine](#-transit-engine)
- [Kubernetes Integration](#-kubernetes-integration)
- [High Availability](#-high-availability)
- [Production Tips](#-production-tips)
- [Troubleshooting](#-troubleshooting)
- [Related Topics](#-related-topics)

---

## 🧠 What is Vault?

HashiCorp Vault is a tool for securely accessing secrets. A "secret" is anything you want to tightly control access to: API keys, passwords, certificates, encryption keys, database credentials, cloud tokens, etc.

**Key characteristics:**

- **Centralized secret management** — Single source of truth for all secrets
- **Dynamic secrets** — Generate credentials on-demand with automatic expiry
- **Encryption as a service** — Encrypt/decrypt data without managing keys
- **Identity-based access** — Policies tied to authenticated identities
- **Audit logging** — Every secret access is logged
- **Lease management** — Secrets have TTLs and can be revoked

**What Vault replaces:**

| Before Vault | With Vault |
|-------------|-----------|
| Secrets in environment variables | Dynamic secrets from Vault API |
| Hardcoded credentials in code | Application authenticates to Vault |
| Shared password spreadsheets | Role-based access with audit trail |
| Manual key rotation | Automatic rotation and dynamic credentials |
| Scattered secret stores | Unified secret management |

---

## 🏗 Architecture

```mermaid
graph TB
    subgraph "Vault Server"
        API[HTTP API<br/>Port 8200]
        BARRIER[Encryption Barrier]
        
        subgraph "Core"
            AUTH[Auth Methods<br/>Token, AppRole,<br/>K8s, LDAP, OIDC]
            POLICY[Policy Engine<br/>ACL Policies]
            SECRET[Secret Engines<br/>KV, PKI, Transit,<br/>Database, AWS]
            AUDIT[Audit Devices<br/>File, Syslog, Socket]
        end
        
        STORAGE[(Storage Backend<br/>Raft, Consul,<br/>S3, DynamoDB)]
    end
    
    CLIENT[Client / Application] -->|HTTPS| API
    API --> BARRIER
    BARRIER --> AUTH
    AUTH --> POLICY
    POLICY --> SECRET
    SECRET --> STORAGE
    API --> AUDIT
    
    UNSEAL[Unseal Keys<br/>Shamir's Secret Sharing] -.->|unseal| BARRIER
```

### Seal/Unseal Process

```text
Vault starts SEALED → cannot read/write any secrets

Unsealing:
┌──────────────┐
│ Master Key   │ ── split via Shamir's Secret Sharing ──►
└──────────────┘
    ┌─────────┐  ┌─────────┐  ┌─────────┐  ┌─────────┐  ┌─────────┐
    │ Share 1 │  │ Share 2 │  │ Share 3 │  │ Share 4 │  │ Share 5 │
    └─────────┘  └─────────┘  └─────────┘  └─────────┘  └─────────┘
    
    Need 3 of 5 shares to unseal (threshold = 3)

Auto-unseal (recommended for production):
  Uses cloud KMS (AWS KMS, GCP KMS, Azure Key Vault)
  Vault automatically unseals on restart
```

---

## 🗄 Secret Engines

### KV (Key-Value) v2

```bash
# Enable KV v2 engine
vault secrets enable -path=secret kv-v2

# Write a secret
vault kv put secret/myapp/config \
  db_host="db.example.com" \
  db_port="5432" \
  db_password="s3cur3-p@ss"

# Read a secret
vault kv get secret/myapp/config

# Read specific field
vault kv get -field=db_password secret/myapp/config

# Read as JSON
vault kv get -format=json secret/myapp/config

# List secrets
vault kv list secret/myapp/

# Delete (soft delete — can be undeleted)
vault kv delete secret/myapp/config

# Undelete
vault kv undelete -versions=2 secret/myapp/config

# Permanently destroy
vault kv destroy -versions=1,2 secret/myapp/config

# View version history
vault kv metadata get secret/myapp/config

# Rollback to previous version
vault kv rollback -version=1 secret/myapp/config
```

### PKI (Certificate Authority)

```bash
# Enable PKI engine
vault secrets enable pki

# Configure max TTL
vault secrets tune -max-lease-ttl=87600h pki

# Generate root CA
vault write pki/root/generate/internal \
  common_name="Example Root CA" \
  ttl=87600h

# Enable intermediate CA
vault secrets enable -path=pki_int pki
vault secrets tune -max-lease-ttl=43800h pki_int

# Create role for issuing certificates
vault write pki_int/roles/example-dot-com \
  allowed_domains="example.com" \
  allow_subdomains=true \
  max_ttl=720h

# Issue a certificate
vault write pki_int/issue/example-dot-com \
  common_name="app.example.com" \
  ttl=24h
```

### Database

```bash
# Enable database engine
vault secrets enable database

# Configure PostgreSQL connection
vault write database/config/mydb \
  plugin_name=postgresql-database-plugin \
  allowed_roles="readonly,readwrite" \
  connection_url="postgresql://{{username}}:{{password}}@db.example.com:5432/mydb?sslmode=require" \
  username="vault_admin" \
  password="vault_admin_password"

# Create a role (generates dynamic credentials)
vault write database/roles/readonly \
  db_name=mydb \
  creation_statements="CREATE ROLE \"{{name}}\" WITH LOGIN PASSWORD '{{password}}' VALID UNTIL '{{expiration}}'; GRANT SELECT ON ALL TABLES IN SCHEMA public TO \"{{name}}\";" \
  default_ttl=1h \
  max_ttl=24h

# Get dynamic credentials
vault read database/creds/readonly
# Key                Value
# ---                -----
# lease_id           database/creds/readonly/abc123
# lease_duration     1h
# username           v-token-readonly-abc123-1234567890
# password           A1b2C3d4E5-randomly-generated
```

### Transit (Encryption as a Service)

See dedicated section below.

### Other Secret Engines

| Engine | Description | Use Case |
|--------|-------------|----------|
| **AWS** | Dynamic AWS IAM credentials | Temporary AWS access |
| **GCP** | Dynamic GCP service account keys | Temporary GCP access |
| **SSH** | SSH certificates and OTP | Secure SSH access |
| **TOTP** | Time-based one-time passwords | MFA generation |
| **Consul** | Dynamic Consul tokens | Consul ACL management |
| **Nomad** | Dynamic Nomad tokens | Nomad access management |
| **RabbitMQ** | Dynamic RabbitMQ credentials | Message queue access |

---

## 🔑 Authentication Methods

### Token Auth (Default)

```bash
# Create a token
vault token create -ttl=1h -policy=my-policy

# Create a periodic token (auto-renews)
vault token create -period=24h -policy=my-policy

# Create an orphan token (no parent)
vault token create -orphan -ttl=8h -policy=my-policy

# Lookup token information
vault token lookup

# Renew token
vault token renew

# Revoke token
vault token revoke <token>

# Revoke all tokens for an accessor
vault token revoke -accessor <accessor>
```

### AppRole (Machine-to-Machine)

```bash
# Enable AppRole auth
vault auth enable approle

# Create a role
vault write auth/approle/role/my-app \
  token_policies="my-app-policy" \
  token_ttl=1h \
  token_max_ttl=4h \
  secret_id_ttl=720h

# Get Role ID (deploy-time, stable)
vault read auth/approle/role/my-app/role-id

# Generate Secret ID (runtime, rotated)
vault write -f auth/approle/role/my-app/secret-id

# Login with AppRole
vault write auth/approle/login \
  role_id="abc-123-def" \
  secret_id="xyz-789-uvw"
```

### Kubernetes Auth

```bash
# Enable Kubernetes auth
vault auth enable kubernetes

# Configure Kubernetes auth
vault write auth/kubernetes/config \
  kubernetes_host="https://kubernetes.default.svc:443"

# Create a role
vault write auth/kubernetes/role/my-app \
  bound_service_account_names=my-app-sa \
  bound_service_account_namespaces=production \
  policies=my-app-policy \
  ttl=1h
```

### Other Auth Methods

| Method | Use Case | Identity Source |
|--------|----------|----------------|
| **LDAP** | Corporate directory | LDAP/AD groups |
| **OIDC** | SSO with identity providers | OIDC provider (Okta, Auth0) |
| **GitHub** | GitHub-based teams | GitHub org/teams |
| **AWS** | AWS IAM roles/instances | IAM role or EC2 instance profile |
| **GCP** | GCP service accounts | GCP IAM |
| **TLS Certificate** | Client certificate auth | X.509 certificates |
| **Userpass** | Username/password | Local accounts (dev/testing) |

---

## 📝 Policies

Policies define what a token can access. Written in HCL (HashiCorp Configuration Language).

### Policy Syntax

```hcl
# Read-only access to application secrets
path "secret/data/myapp/*" {
  capabilities = ["read", "list"]
}

# Full access to a specific path
path "secret/data/myapp/config" {
  capabilities = ["create", "read", "update", "delete", "list"]
}

# Generate database credentials
path "database/creds/readonly" {
  capabilities = ["read"]
}

# Issue certificates
path "pki_int/issue/example-dot-com" {
  capabilities = ["create", "update"]
}

# Deny access to specific path
path "secret/data/production/admin/*" {
  capabilities = ["deny"]
}

# Token self-management
path "auth/token/renew-self" {
  capabilities = ["update"]
}
path "auth/token/lookup-self" {
  capabilities = ["read"]
}
```

### Capabilities

| Capability | HTTP Verb | Description |
|-----------|-----------|-------------|
| `create` | POST | Create new data |
| `read` | GET | Read data |
| `update` | PUT/POST | Update existing data |
| `delete` | DELETE | Delete data |
| `list` | LIST | List keys/paths |
| `deny` | - | Explicitly deny access |
| `sudo` | - | Root-level operations |

### Policy Management

```bash
# Write a policy
vault policy write my-app-policy - <<EOF
path "secret/data/myapp/*" {
  capabilities = ["read", "list"]
}
path "database/creds/readonly" {
  capabilities = ["read"]
}
EOF

# Read a policy
vault policy read my-app-policy

# List policies
vault policy list

# Delete a policy
vault policy delete my-app-policy

# Test a policy (with token)
vault token capabilities <token> secret/data/myapp/config
```

---

## 💻 CLI Reference

### Core Commands

```bash
# Login
vault login <token>
vault login -method=ldap username=myuser
vault login -method=oidc

# Status
vault status

# Seal/unseal
vault operator seal
vault operator unseal <share>

# Initialize new Vault
vault operator init -key-shares=5 -key-threshold=3

# Auto-unseal (check seal status)
vault status | grep "Seal Type"
```

### Secret Operations

```bash
# KV operations
vault kv put secret/path key=value
vault kv get secret/path
vault kv get -field=key secret/path
vault kv get -format=json secret/path
vault kv list secret/
vault kv delete secret/path
vault kv metadata get secret/path
vault kv metadata delete secret/path

# Enable/disable secret engines
vault secrets enable -path=myengine kv-v2
vault secrets disable myengine
vault secrets list
vault secrets tune -max-lease-ttl=72h myengine
```

### Auth Operations

```bash
# Enable/disable auth methods
vault auth enable -path=myauth approle
vault auth disable myauth
vault auth list

# Token operations
vault token create -policy=my-policy -ttl=1h
vault token lookup
vault token renew
vault token revoke <token>
```

### Lease Operations

```bash
# List leases
vault list sys/leases/lookup/database/creds/readonly

# Renew a lease
vault lease renew <lease-id>
vault lease renew -increment=3600 <lease-id>

# Revoke a lease
vault lease revoke <lease-id>

# Revoke all leases under a prefix
vault lease revoke -prefix database/creds/readonly
```

---

## 🔄 Dynamic Secrets

### Database Credentials

```bash
# Request dynamic PostgreSQL credentials
vault read database/creds/readonly

# Output:
# Key                Value
# ---                -----
# lease_id           database/creds/readonly/abcd-1234
# lease_duration     1h
# lease_renewable    true
# password           A1b2C3d4-generated
# username           v-approle-readonly-abc123

# Renew before expiry
vault lease renew database/creds/readonly/abcd-1234

# Revoke when done
vault lease revoke database/creds/readonly/abcd-1234
```

### AWS Dynamic Credentials

```bash
# Enable AWS engine
vault secrets enable aws

# Configure root credentials
vault write aws/config/root \
  access_key="AKIA..." \
  secret_key="..." \
  region="us-east-1"

# Create a role
vault write aws/roles/s3-readonly \
  credential_type=iam_user \
  policy_document=-<<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Action": "s3:GetObject",
      "Resource": "*"
    }
  ]
}
EOF

# Get temporary AWS credentials
vault read aws/creds/s3-readonly
```

### PKI Certificates

```bash
# Issue a short-lived certificate
vault write pki_int/issue/example-dot-com \
  common_name="api.example.com" \
  alt_names="api2.example.com" \
  ttl=24h

# Output includes:
# certificate      -----BEGIN CERTIFICATE-----...
# issuing_ca       -----BEGIN CERTIFICATE-----...
# private_key      -----BEGIN RSA PRIVATE KEY-----...
# serial_number    abc123
```

---

## 🔐 Transit Engine

Vault Transit provides **encryption as a service** — encrypt/decrypt data without exposing encryption keys.

```bash
# Enable Transit engine
vault secrets enable transit

# Create an encryption key
vault write -f transit/keys/my-app-key

# Encrypt data (input must be base64-encoded)
vault write transit/encrypt/my-app-key \
  plaintext=$(echo -n "sensitive-data" | base64)
# Output: ciphertext = vault:v1:abc123...

# Decrypt data
vault write transit/decrypt/my-app-key \
  ciphertext="vault:v1:abc123..."
# Output: plaintext = c2Vuc2l0aXZlLWRhdGE= (base64)
echo "c2Vuc2l0aXZlLWRhdGE=" | base64 -d
# Output: sensitive-data

# Rotate encryption key
vault write -f transit/keys/my-app-key/rotate

# Rewrap ciphertext with latest key version
vault write transit/rewrap/my-app-key \
  ciphertext="vault:v1:abc123..."
# Output: ciphertext = vault:v2:def456... (encrypted with new key)

# Key configuration
vault write transit/keys/my-app-key/config \
  min_decryption_version=1 \
  min_encryption_version=2
```

---

## ☸️ Kubernetes Integration

### Vault Agent Injector

```yaml
# Pod annotation-based secret injection
apiVersion: apps/v1
kind: Deployment
metadata:
  name: my-app
spec:
  template:
    metadata:
      annotations:
        vault.hashicorp.com/agent-inject: "true"
        vault.hashicorp.com/role: "my-app"
        vault.hashicorp.com/agent-inject-secret-config: "secret/data/myapp/config"
        vault.hashicorp.com/agent-inject-template-config: |
          {{- with secret "secret/data/myapp/config" -}}
          export DB_HOST={{ .Data.data.db_host }}
          export DB_PASSWORD={{ .Data.data.db_password }}
          {{- end }}
    spec:
      serviceAccountName: my-app-sa
      containers:
        - name: my-app
          image: my-app:latest
          command: ["/bin/sh", "-c", "source /vault/secrets/config && ./start.sh"]
```

### CSI Provider

```yaml
# SecretProviderClass for Vault CSI
apiVersion: secrets-store.csi.x-k8s.io/v1
kind: SecretProviderClass
metadata:
  name: vault-db-creds
spec:
  provider: vault
  parameters:
    roleName: "my-app"
    vaultAddress: "https://vault.example.com:8200"
    objects: |
      - objectName: "db-password"
        secretPath: "secret/data/myapp/config"
        secretKey: "db_password"

---
# Pod using CSI volume
apiVersion: v1
kind: Pod
metadata:
  name: my-app
spec:
  serviceAccountName: my-app-sa
  containers:
    - name: my-app
      image: my-app:latest
      volumeMounts:
        - name: secrets
          mountPath: "/mnt/secrets"
          readOnly: true
  volumes:
    - name: secrets
      csi:
        driver: secrets-store.csi.k8s.io
        readOnly: true
        volumeAttributes:
          secretProviderClass: "vault-db-creds"
```

---

## 🔄 High Availability

### Raft Storage (Recommended)

```hcl
# vault.hcl configuration
storage "raft" {
  path = "/vault/data"
  node_id = "vault-1"
}

listener "tcp" {
  address     = "0.0.0.0:8200"
  tls_cert_file = "/vault/tls/server.crt"
  tls_key_file  = "/vault/tls/server.key"
}

api_addr     = "https://vault-1.example.com:8200"
cluster_addr = "https://vault-1.example.com:8201"

# Auto-unseal with AWS KMS
seal "awskms" {
  region     = "us-east-1"
  kms_key_id = "alias/vault-unseal"
}
```

```bash
# Join additional nodes to Raft cluster
vault operator raft join https://vault-1.example.com:8200

# List Raft peers
vault operator raft list-peers

# Take Raft snapshot
vault operator raft snapshot save /backup/vault-$(date +%Y%m%d).snap

# Restore Raft snapshot
vault operator raft snapshot restore /backup/vault-20250101.snap
```

---

## 🏭 Production Tips

### Deployment Checklist

```text
□ Auto-unseal configured (AWS KMS, GCP KMS, Azure Key Vault)
□ TLS enabled for all communication
□ Audit logging enabled (file and/or syslog)
□ HA storage backend (Raft or Consul)
□ Regular backups (Raft snapshots)
□ Monitoring configured (sealed status, token count, lease count)
□ Policies follow least privilege
□ Root token revoked after initial setup
□ Namespaces configured (Enterprise) for multi-tenancy
□ Tested disaster recovery procedure
□ Secret engine TTLs configured appropriately
□ Token TTLs enforced (no indefinite tokens)
```

### Audit Logging

```bash
# Enable file audit device
vault audit enable file file_path=/var/log/vault/audit.log

# Enable syslog audit device
vault audit enable syslog tag="vault" facility="AUTH"

# List audit devices
vault audit list

# Audit log entry (JSON):
# {
#   "time": "2025-01-01T12:00:00Z",
#   "type": "request",
#   "auth": { "token": "hmac-sha256:...", "policies": ["my-policy"] },
#   "request": { "path": "secret/data/myapp/config", "operation": "read" }
# }
```

### Key Metrics

| Metric | Description | Alert Threshold |
|--------|-------------|-----------------|
| `vault.core.unsealed` | Vault seal status | = 0 (sealed!) |
| `vault.token.count` | Active tokens | Unusually high growth |
| `vault.expire.num_leases` | Active leases | Approaching limits |
| `vault.audit.log_response` | Audit log errors | Any failures |
| `vault.runtime.alloc_bytes` | Memory usage | > 80% of limit |
| `vault.raft.leader.lastContact` | Raft leader contact | > 200ms |

---

## 🐛 Troubleshooting

### Sealed Vault

```bash
# Check seal status
vault status

# If sealed, unseal with shares
vault operator unseal <share-1>
vault operator unseal <share-2>
vault operator unseal <share-3>

# If auto-unseal fails, check:
# 1. Cloud KMS connectivity
# 2. IAM permissions for KMS access
# 3. KMS key state (not disabled/deleted)
```

### Authentication Failures

```bash
# Check if auth method is enabled
vault auth list

# For Kubernetes auth:
# 1. Check service account exists
kubectl get sa my-app-sa -n production

# 2. Check Vault role configuration
vault read auth/kubernetes/role/my-app

# 3. Check Kubernetes API connectivity
vault write auth/kubernetes/config \
  kubernetes_host="https://kubernetes.default.svc:443"
```

### Policy Denied

```bash
# Check token policies
vault token lookup

# Check what capabilities the token has on a path
vault token capabilities <token> secret/data/myapp/config

# Test policy without applying
vault policy read my-app-policy

# Common issues:
# - Missing "data" in KV v2 paths: secret/data/myapp not secret/myapp
# - Wrong capabilities: "read" needed for GET, "create/update" for PUT
# - Policy not attached to auth role
```

### Token Expired

```bash
# Check token TTL
vault token lookup

# Renew before expiry
vault token renew

# If expired, re-authenticate:
vault login -method=approle role_id=... secret_id=...
```

---

## 🔗 Related Topics

- [Security — TLS](../tls/) — TLS certificates for Vault and from Vault PKI
- [Security — JWT](../jwt/) — JWT signing key storage in Vault
- [Security — OAuth](../oauth/) — Client secret storage in Vault
- [DevOps — Kubernetes](../../devops/kubernetes/) — Vault integration with Kubernetes
- [Databases — PostgreSQL](../../databases/postgresql/) — Dynamic database credentials
- [Databases — Redis](../../databases/redis/) — Redis password management with Vault

---

> **Never store secrets in code, environment variables, or configuration files.** Use Vault to centralize secrets management, enable dynamic credentials, and maintain a complete audit trail. Start simple with the KV engine and progressively adopt dynamic secrets and transit encryption.
