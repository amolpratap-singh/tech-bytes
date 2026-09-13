# 📊 Grafana

> **The open-source analytics and visualization platform.** Grafana turns your metrics, logs, and traces into beautiful, interactive dashboards. It connects to virtually any data source and provides a unified view across your entire observability stack.

---

## 📑 Table of Contents

- [What is Grafana?](#-what-is-grafana)
- [Core Concepts](#-core-concepts)
- [Data Sources](#-data-sources)
- [Dashboard Design](#-dashboard-design)
- [Panel Types](#-panel-types)
- [Variables](#-variables)
- [Alerting](#-alerting)
- [Provisioning](#-provisioning)
- [Dashboard Best Practices](#-dashboard-best-practices)
- [Useful Dashboards](#-useful-dashboards)
- [Troubleshooting](#-troubleshooting)
- [Related Topics](#-related-topics)

---

## 🧠 What is Grafana?

Grafana is an open-source platform for monitoring and observability visualization. It doesn't store data — it queries data sources in real time and renders dashboards, graphs, and alerts.

**Key features:**

- **Multi-datasource** — Query Prometheus, VictoriaMetrics, Elasticsearch, Loki, InfluxDB, CloudWatch, and 100+ more
- **Rich visualizations** — Time series, heatmaps, tables, stats, gauges, logs, node graphs
- **Template variables** — Dynamic dashboards that adapt to different environments, clusters, namespaces
- **Alerting** — Built-in alert evaluation and notification routing
- **Provisioning** — Dashboards as code via YAML configuration and JSON models
- **Plugins** — Extensible with community and enterprise plugins

### Versions

| Edition | Features | License |
|---------|----------|---------|
| Grafana OSS | Core dashboarding, alerting, data sources | AGPL v3 |
| Grafana Enterprise | RBAC, audit, reporting, advanced plugins | Commercial |
| Grafana Cloud | Managed Grafana + Mimir + Loki + Tempo | SaaS |

---

## 🧱 Core Concepts

### Hierarchy

```text
Organization
├── Data Sources (Prometheus, Loki, etc.)
├── Folders
│   └── Dashboards
│       ├── Variables
│       ├── Annotations
│       ├── Rows
│       │   └── Panels
│       │       ├── Query (PromQL, LogQL, etc.)
│       │       ├── Transformations
│       │       └── Overrides
│       └── Links
└── Alert Rules
    ├── Alert Folders
    ├── Notification Policies
    ├── Contact Points
    └── Silences
```

### Key Terms

| Term | Description |
|------|-------------|
| **Dashboard** | A collection of panels organized in rows, focused on a topic |
| **Panel** | A single visualization (graph, stat, table, etc.) |
| **Data Source** | External system Grafana queries (Prometheus, Loki, etc.) |
| **Query** | The expression sent to a data source (PromQL, LogQL, SQL, etc.) |
| **Variable** | A dynamic parameter used in queries and titles (dropdown selector) |
| **Annotation** | Vertical marker on time series showing events (deploys, incidents) |
| **Row** | A horizontal grouping of panels (collapsible) |
| **Transformation** | Post-query data manipulation (join, filter, calculate) |
| **Override** | Per-series customization of display options |
| **Link** | Navigation to other dashboards or external URLs |

---

## 🔌 Data Sources

### Configuration via UI

Settings → Data Sources → Add data source

### Prometheus / VictoriaMetrics

```yaml
# Provisioning: datasources/prometheus.yml
apiVersion: 1
datasources:
  - name: Prometheus
    type: prometheus
    access: proxy
    url: http://prometheus:9090
    isDefault: true
    jsonData:
      timeInterval: "15s"       # Default scrape interval
      queryTimeout: "60s"
      httpMethod: POST          # Better for large queries

  - name: VictoriaMetrics
    type: prometheus
    access: proxy
    url: http://victoriametrics:8428
    jsonData:
      timeInterval: "15s"
```

### Loki (Logs)

```yaml
  - name: Loki
    type: loki
    access: proxy
    url: http://loki:3100
    jsonData:
      derivedFields:
        - datasourceUid: tempo-uid
          matcherRegex: "trace_id=(\\w+)"
          name: TraceID
          url: '$${__value.raw}'
```

### Elasticsearch

```yaml
  - name: Elasticsearch
    type: elasticsearch
    access: proxy
    url: http://elasticsearch:9200
    database: "logs-*"
    jsonData:
      esVersion: "8.0.0"
      timeField: "@timestamp"
      logMessageField: message
      logLevelField: level
```

### InfluxDB

```yaml
  - name: InfluxDB
    type: influxdb
    access: proxy
    url: http://influxdb:8086
    jsonData:
      version: Flux
      organization: myorg
      defaultBucket: metrics
    secureJsonData:
      token: "<influxdb-token>"
```

---

## 🎨 Dashboard Design

### Dashboard JSON Structure

```json
{
  "dashboard": {
    "title": "Service Overview",
    "uid": "service-overview",
    "tags": ["service", "production"],
    "timezone": "browser",
    "refresh": "30s",
    "time": {
      "from": "now-1h",
      "to": "now"
    },
    "templating": {
      "list": []
    },
    "panels": [],
    "rows": []
  }
}
```

### Annotations

Mark events on dashboards (deployments, incidents, config changes):

```json
{
  "annotations": {
    "list": [
      {
        "datasource": "Prometheus",
        "enable": true,
        "expr": "changes(deployment_info[1m]) > 0",
        "name": "Deployments",
        "iconColor": "green",
        "titleFormat": "Deploy: {{version}}"
      }
    ]
  }
}
```

### Dashboard Links

```json
{
  "links": [
    {
      "title": "Logs",
      "url": "/d/logs-dashboard/logs?var-service=${service}",
      "type": "link",
      "icon": "doc"
    },
    {
      "title": "Related Dashboards",
      "type": "dashboards",
      "tags": ["service"]
    }
  ]
}
```

---

## 📊 Panel Types

### Time Series (Default)

The most common panel. Shows metrics over time.

```text
Query: rate(http_requests_total{job="$service"}[5m])
Legend: {{handler}} - {{status}}
```

**Key options:**
- Line width, fill opacity, gradient mode
- Stack series (normal, percent)
- Show points, bars, or lines
- Axis scale (linear, log2, log10)
- Thresholds (color zones)
- Override per-series display

### Stat

Single large number with optional sparkline. Ideal for KPI panels.

```text
Query: sum(rate(http_requests_total{job="$service"}[5m]))
Display: Value + sparkline
Thresholds: Green < 100, Yellow < 500, Red >= 500
```

### Gauge

Circular gauge showing current value against min/max.

```text
Query: (1 - node_memory_MemAvailable_bytes / node_memory_MemTotal_bytes) * 100
Min: 0, Max: 100
Thresholds: Green < 70, Yellow < 85, Red >= 85
```

### Bar Gauge

Horizontal or vertical bars.

```text
Query: topk(10, sum by (handler) (rate(http_requests_total[5m])))
Display: Horizontal bars
Sort: Descending by value
```

### Table

Tabular data display with sorting, filtering, and column formatting.

```text
Query A: sum by (pod) (rate(container_cpu_usage_seconds_total[5m]))
Query B: sum by (pod) (container_memory_working_set_bytes)
Transform: Merge (join on pod label)
```

### Heatmap

Visualize histogram distributions over time.

```text
Query: sum by (le) (rate(http_request_duration_seconds_bucket[5m]))
Format: Heatmap
Y-axis: Bucket boundaries (le values)
Color: Scheme → Temperature
```

### Logs Panel

Display log lines from Loki or Elasticsearch.

```text
Query (LogQL): {job="$service"} |= "$search" | logfmt
Display: Log lines with severity colors
```

### Node Graph

Visualize relationships and dependencies (service mesh, network topology).

### Geomap

Geographic visualization for location-based metrics.

---

## 🔧 Variables

Variables make dashboards dynamic and reusable.

### Variable Types

| Type | Description | Example |
|------|-------------|---------|
| **Query** | Values from a data source query | Namespace list from Prometheus |
| **Custom** | Static list of values | `production, staging, development` |
| **Text box** | Free-text input | Search filter |
| **Interval** | Time interval selection | `1m, 5m, 15m, 1h` |
| **Data source** | Select between data sources | Choose Prometheus instance |
| **Ad hoc filters** | Dynamic label filters | Add any label filter on the fly |
| **Constant** | Hidden constant value | Base URL, cluster name |

### Query Variable Examples

```promql
# List all namespaces
label_values(kube_pod_info, namespace)

# List pods in selected namespace
label_values(kube_pod_info{namespace="$namespace"}, pod)

# List all jobs
label_values(up, job)

# List instances for a job
label_values(up{job="$job"}, instance)

# List handler values
label_values(http_requests_total{job="$service"}, handler)

# List node names
label_values(kube_node_info, node)
```

### Variable Configuration

```json
{
  "templating": {
    "list": [
      {
        "name": "datasource",
        "type": "datasource",
        "query": "prometheus",
        "current": { "text": "Prometheus", "value": "Prometheus" }
      },
      {
        "name": "namespace",
        "type": "query",
        "datasource": "$datasource",
        "query": "label_values(kube_pod_info, namespace)",
        "refresh": 2,
        "sort": 1,
        "multi": true,
        "includeAll": true,
        "allValue": ".*"
      },
      {
        "name": "service",
        "type": "query",
        "datasource": "$datasource",
        "query": "label_values(up{namespace=~\"$namespace\"}, job)",
        "refresh": 2,
        "sort": 1
      },
      {
        "name": "interval",
        "type": "interval",
        "query": "1m,5m,15m,30m,1h",
        "current": { "text": "5m", "value": "5m" }
      }
    ]
  }
}
```

### Using Variables in Queries

```promql
# In PromQL queries
rate(http_requests_total{job="$service", namespace=~"$namespace"}[$__rate_interval])

# Multi-value variable (use regex matcher)
rate(http_requests_total{namespace=~"$namespace"}[5m])

# Built-in variables
$__rate_interval    # Recommended rate() range (auto-calculated)
$__interval         # Dashboard auto-interval
$__range            # Selected time range
$__from / $__to     # Start/end timestamps
```

> **Best practice:** Always use `$__rate_interval` instead of hardcoded ranges like `[5m]` in `rate()` calls. It automatically adjusts based on the dashboard time range and scrape interval.

---

## 🚨 Alerting

### Grafana Alerting (Unified Alerting)

Since Grafana 9, alerting is built into Grafana with a unified model.

### Alert Rule

```yaml
# Provisioning: alerting/rules.yml
apiVersion: 1
groups:
  - orgId: 1
    name: service_alerts
    folder: Production
    interval: 1m
    rules:
      - uid: high-error-rate
        title: High Error Rate
        condition: C
        data:
          - refId: A
            datasourceUid: prometheus
            model:
              expr: sum(rate(http_requests_total{status=~"5.."}[5m]))
              instant: false
          - refId: B
            datasourceUid: prometheus
            model:
              expr: sum(rate(http_requests_total[5m]))
              instant: false
          - refId: C
            datasourceUid: __expr__
            model:
              type: math
              expression: "$A / $B > 0.05"
        for: 5m
        labels:
          severity: critical
        annotations:
          summary: "Error rate is above 5%"
```

### Contact Points

```yaml
# Provisioning: alerting/contactpoints.yml
apiVersion: 1
contactPoints:
  - orgId: 1
    name: slack-critical
    receivers:
      - uid: slack-1
        type: slack
        settings:
          url: "https://hooks.slack.com/services/xxx/yyy/zzz"
          recipient: "#alerts-critical"
          title: '{{ template "slack.default.title" . }}'
          text: '{{ template "slack.default.text" . }}'

  - orgId: 1
    name: pagerduty
    receivers:
      - uid: pd-1
        type: pagerduty
        settings:
          integrationKey: "<service-key>"
          severity: critical
```

### Notification Policies

```yaml
# Provisioning: alerting/policies.yml
apiVersion: 1
policies:
  - orgId: 1
    receiver: slack-default
    group_by: ["alertname", "namespace"]
    group_wait: 30s
    group_interval: 5m
    repeat_interval: 4h
    routes:
      - receiver: pagerduty
        matchers:
          - severity = critical
        continue: true
      - receiver: slack-critical
        matchers:
          - severity = critical
      - receiver: slack-warnings
        matchers:
          - severity = warning
```

### Grafana Alerting vs External Alerting

| Feature | Grafana Alerting | Prometheus + Alertmanager |
|---------|-----------------|--------------------------|
| Query Language | Multi-datasource | PromQL only |
| Configuration | UI + provisioning | YAML files |
| Visualization | Alerts shown on panels | Separate UI |
| Multi-datasource | Yes | No |
| Alert State | Stored in Grafana DB | Stored in Alertmanager |
| Silence/Mute | Built-in UI | Alertmanager UI |
| Best for | Teams using Grafana UI | GitOps, Prometheus-native teams |

> **Recommendation:** For Prometheus-only setups, use Prometheus alerting rules + Alertmanager. For multi-datasource environments, Grafana Unified Alerting is more flexible.

---

## 📦 Provisioning

### Directory Structure

```text
/etc/grafana/provisioning/
├── dashboards/
│   └── dashboards.yml          # Dashboard provider config
├── datasources/
│   └── datasources.yml         # Data source definitions
├── alerting/
│   ├── rules.yml               # Alert rules
│   ├── contactpoints.yml       # Notification channels
│   └── policies.yml            # Routing policies
├── notifiers/                   # Legacy notification channels
└── plugins/
    └── plugins.yml             # Plugin installation
```

### Dashboard Provider

```yaml
# provisioning/dashboards/dashboards.yml
apiVersion: 1
providers:
  - name: "default"
    orgId: 1
    folder: "Provisioned"
    type: file
    disableDeletion: false
    editable: true
    updateIntervalSeconds: 30
    options:
      path: /var/lib/grafana/dashboards
      foldersFromFilesStructure: true   # Create folders from subdirectories
```

### Dashboard as JSON

```bash
# Export dashboard
curl -H "Authorization: Bearer <api-key>" \
  http://grafana:3000/api/dashboards/uid/my-dashboard | jq '.dashboard' > dashboard.json

# Import dashboard
curl -X POST -H "Content-Type: application/json" \
  -H "Authorization: Bearer <api-key>" \
  -d '{
    "dashboard": '"$(cat dashboard.json)"',
    "overwrite": true,
    "folderId": 0
  }' http://grafana:3000/api/dashboards/db
```

### Grafana API — Common Operations

```bash
# List dashboards
curl -H "Authorization: Bearer <api-key>" http://grafana:3000/api/search

# Get dashboard by UID
curl -H "Authorization: Bearer <api-key>" http://grafana:3000/api/dashboards/uid/<uid>

# Create/update dashboard
curl -X POST -H "Content-Type: application/json" \
  -H "Authorization: Bearer <api-key>" \
  -d @dashboard-payload.json http://grafana:3000/api/dashboards/db

# List data sources
curl -H "Authorization: Bearer <api-key>" http://grafana:3000/api/datasources

# Create folder
curl -X POST -H "Content-Type: application/json" \
  -H "Authorization: Bearer <api-key>" \
  -d '{"title": "Production"}' http://grafana:3000/api/folders

# Create API key
curl -X POST -H "Content-Type: application/json" \
  -d '{"name": "ci-cd", "role": "Editor"}' \
  http://admin:admin@grafana:3000/api/auth/keys
```

### Grafana with Docker Compose

```yaml
version: '3.8'
services:
  grafana:
    image: grafana/grafana:10.3.0
    ports:
      - "3000:3000"
    environment:
      - GF_SECURITY_ADMIN_PASSWORD=admin
      - GF_USERS_ALLOW_SIGN_UP=false
      - GF_SERVER_ROOT_URL=https://grafana.example.com
      - GF_INSTALL_PLUGINS=grafana-piechart-panel
    volumes:
      - grafana-data:/var/lib/grafana
      - ./provisioning:/etc/grafana/provisioning
      - ./dashboards:/var/lib/grafana/dashboards
    restart: unless-stopped

volumes:
  grafana-data:
```

---

## ✅ Dashboard Best Practices

### Naming Convention

```text
[Environment] Service Name - Focus Area
Examples:
  [Prod] API Gateway - Overview
  [Prod] API Gateway - Latency Deep Dive
  [Prod] Node Resources - CPU & Memory
  [All] Kubernetes - Cluster Overview
```

### Dashboard Hierarchy

```text
Level 1: Overview (all services at a glance)
  ├── Total request rate, error rate, p99 latency
  ├── Service health status (up/down)
  └── Links to Level 2

Level 2: Service Dashboard (per-service details)
  ├── RED metrics (Rate, Errors, Duration)
  ├── Resource usage (CPU, memory)
  ├── Dependencies health
  └── Links to Level 3

Level 3: Deep Dive (detailed debugging)
  ├── Per-endpoint latency breakdown
  ├── Trace examples
  ├── Log panel with filters
  └── Specific component metrics
```

### Design Guidelines

1. **Use variables** — Make dashboards work across environments (namespace, cluster, service)
2. **Use `$__rate_interval`** — Instead of hardcoded rate ranges
3. **Include a data source variable** — For multi-cluster setups
4. **Add dashboard links** — Connect overview → detail → logs → traces
5. **Use consistent colors** — Red = error/critical, Yellow = warning, Green = healthy
6. **Add descriptions** — Panel descriptions explain what the metric means and what to look for
7. **Use annotations** — Show deployments, incidents, config changes on graphs
8. **Use rows** — Group related panels, make rows collapsible for large dashboards
9. **Set sensible defaults** — Default time range (1h), refresh interval (30s)
10. **Include thresholds** — Visual indicators for acceptable ranges
11. **Tag dashboards** — Use consistent tags: `service`, `infrastructure`, `kubernetes`
12. **Use repeat panels** — Repeat a panel for each value of a variable (one per pod, per namespace)

### Color Scheme

```text
Green:  #73BF69  — Healthy, within bounds
Yellow: #FADE2A  — Warning, approaching limit
Orange: #FF9830  — High warning
Red:    #F2495C  — Critical, SLO breach
Blue:   #5794F2  — Informational
Purple: #B877D9  — Auxiliary metric
```

---

## 📋 Useful Dashboards

### Community Dashboards (grafana.com/dashboards)

| Dashboard | ID | Description |
|-----------|----|-------------|
| Node Exporter Full | 1860 | Comprehensive host metrics |
| Kubernetes Cluster | 6417 | Cluster-wide resource overview |
| Kubernetes Pods | 6336 | Per-pod CPU, memory, network |
| Docker Dashboard | 893 | Container resource metrics |
| NGINX | 12708 | NGINX request metrics |
| PostgreSQL | 9628 | PostgreSQL performance |
| Redis | 11835 | Redis performance and memory |
| JVM Micrometer | 4701 | JVM heap, GC, threads |
| Go Processes | 6671 | Go runtime metrics |
| Prometheus Stats | 2 | Prometheus self-monitoring |
| VictoriaMetrics | 10229 | VM self-monitoring |

### Import a Dashboard

```bash
# Via UI: Dashboards → Import → Enter ID → Load → Select data source → Import

# Via API
curl -X POST -H "Content-Type: application/json" \
  -H "Authorization: Bearer <api-key>" \
  -d '{
    "dashboard": { "id": null },
    "folderId": 0,
    "overwrite": true,
    "inputs": [{"name": "DS_PROMETHEUS", "type": "datasource", "pluginId": "prometheus", "value": "Prometheus"}],
    "pluginId": "grafana"
  }' http://grafana:3000/api/dashboards/import
```

### Essential Dashboard Panels

**Service Overview (RED):**

```promql
# Request Rate
sum(rate(http_requests_total{job="$service"}[$__rate_interval]))

# Error Rate (%)
sum(rate(http_requests_total{job="$service", status=~"5.."}[$__rate_interval]))
/
sum(rate(http_requests_total{job="$service"}[$__rate_interval]))
* 100

# P50 / P95 / P99 Latency
histogram_quantile(0.50, sum by (le) (rate(http_request_duration_seconds_bucket{job="$service"}[$__rate_interval])))
histogram_quantile(0.95, sum by (le) (rate(http_request_duration_seconds_bucket{job="$service"}[$__rate_interval])))
histogram_quantile(0.99, sum by (le) (rate(http_request_duration_seconds_bucket{job="$service"}[$__rate_interval])))
```

---

## 🐛 Troubleshooting

### No Data in Panel

```text
1. Check time range — is it set correctly? Try "Last 1 hour"
2. Check data source — is the correct source selected?
3. Check variable values — are filters too restrictive?
4. Test query in Explore — go to Explore, paste query, run
5. Check data source connectivity — Settings → Data Sources → Test
6. Check Prometheus targets — ensure target is UP
7. Check metric exists — search in Explore with {__name__=~"partial_name.*"}
```

### Template Variable Issues

```text
Problem: Variable dropdown is empty
1. Check variable query — test in Explore
2. Check dependent variables — parent variable must have a value
3. Check regex filter — too restrictive?
4. Check data source — does the variable use the right source?
5. Check refresh setting — set to "On time range change" or "On dashboard load"

Problem: Variable shows "All" but no data
1. Check allValue setting — should be ".*" for regex matchers
2. Check query uses =~ (regex) not = (exact) for multi-value variables
```

### Dashboard Performance

```text
Symptoms: Slow loading, timeout errors, high browser memory

1. Reduce time range — 24h → 6h → 1h
2. Increase min interval — Panel → Query options → Min interval = 1m
3. Use recording rules — pre-compute expensive queries
4. Reduce panel count — split into sub-dashboards
5. Use $__rate_interval — avoids unnecessarily high resolution
6. Limit topk/bottomk — topk(10, ...) instead of showing all series
7. Use instant queries for stat panels — toggle "Instant" in query options
8. Check query inspector — Panel menu → Inspect → Query to see duration
9. Disable auto-refresh on exploration dashboards
```

### Grafana Configuration Issues

```bash
# Check Grafana logs
docker logs grafana 2>&1 | tail -50
journalctl -u grafana-server -f

# Reset admin password
grafana-cli admin reset-admin-password newpassword

# Check provisioning errors
# Grafana logs on startup: "provisioning" errors indicate YAML issues

# Test data source connectivity
curl -X POST http://admin:admin@grafana:3000/api/datasources/proxy/1/api/v1/query \
  -d 'query=up'
```

---

## 🔗 Related Topics

- [🔥 Prometheus](../prometheus/) — Metrics collection and storage
- [📐 PromQL](../promql/) — Query language for Prometheus panels
- [📈 VictoriaMetrics](../victoriametrics/) — Alternative metrics backend
- [🔭 OpenTelemetry](../opentelemetry/) — Traces and metrics instrumentation
- [📊 Observability Overview](../) — Three pillars and tool comparison

---

> **Grafana is the window into your infrastructure.** Invest in good dashboard design, use variables for reusability, leverage provisioning for dashboards-as-code, and build a clear hierarchy from overview to deep-dive. A well-built dashboard saves hours during incident response.
