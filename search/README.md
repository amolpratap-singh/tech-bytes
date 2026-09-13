# 🔎 Search

> **Search engines power discovery across modern applications.** From full-text search to log analytics to real-time monitoring, search technologies like Elasticsearch and OpenSearch provide the ability to index, query, and analyze massive volumes of data with sub-second response times.

---

## 📑 Table of Contents

- [When to Use a Search Engine](#-when-to-use-a-search-engine)
- [How Search Engines Work](#-how-search-engines-work)
- [OpenSearch vs Elasticsearch](#-opensearch-vs-elasticsearch)
- [Topics](#-topics)
- [Quick Reference](#-quick-reference)
- [Related Topics](#-related-topics)

---

## 🧠 When to Use a Search Engine

| Use Case | Why a Search Engine | Example |
|----------|-------------------|---------|
| **Full-text search** | Inverted index, relevance scoring, tokenization | Product search, document search |
| **Log analytics** | High-volume ingestion, time-series queries | ELK stack, application logs |
| **Metrics & monitoring** | Aggregations, dashboards, alerting | Infrastructure metrics, APM |
| **Autocomplete** | Prefix matching, fuzzy search, suggestions | Search-as-you-type |
| **Faceted navigation** | Multi-dimension filtering with counts | E-commerce category filters |
| **Geospatial search** | Geo queries, distance, bounding box | Location-based services |

### When NOT to Use a Search Engine

- **Primary data store** — Use a database; search engines are secondary indexes
- **ACID transactions** — No transaction support; use PostgreSQL
- **Frequent updates** — High update rate causes segment merges; use a database
- **Strong consistency** — Near-real-time by design (refresh interval); not immediately consistent

---

## ⚙️ How Search Engines Work

```text
Document Ingestion:
  Document → Analyzer → Inverted Index
  
  "The quick brown fox" → tokenize → ["the", "quick", "brown", "fox"]
                        → lowercase → ["the", "quick", "brown", "fox"]
                        → stop words → ["quick", "brown", "fox"]

Inverted Index:
  Term      → Documents
  ─────────────────────
  "quick"   → [doc1, doc5, doc12]
  "brown"   → [doc1, doc7]
  "fox"     → [doc1, doc3, doc7]

Query: "quick fox"
  "quick" → [doc1, doc5, doc12]
  "fox"   → [doc1, doc3, doc7]
  Intersection → [doc1] (highest relevance)
  Union ranked → [doc1, doc3, doc5, doc7, doc12]
```

---

## ⚖️ OpenSearch vs Elasticsearch

| Feature | Elasticsearch | OpenSearch |
|---------|--------------|-----------|
| **License** | SSPL (Server Side Public License) | Apache 2.0 (open source) |
| **Maintained by** | Elastic | AWS + community |
| **Fork point** | Original | Fork of ES 7.10.2 (2021) |
| **Managed service** | Elastic Cloud | Amazon OpenSearch Service |
| **Kibana equivalent** | Kibana | OpenSearch Dashboards |
| **API compatibility** | ES native | ES-compatible (7.10 baseline) |
| **Security** | Paid (X-Pack) or basic | Built-in (free) |
| **Alerting** | Paid (X-Pack) | Built-in (free) |
| **SQL support** | Built-in | Built-in |
| **ML features** | Paid (X-Pack) | Basic anomaly detection (free) |

### When to Choose Which

| Scenario | Recommendation | Why |
|----------|---------------|-----|
| AWS infrastructure | OpenSearch | Native integration, managed service |
| Multi-cloud / on-prem | Either | Both work well |
| Need latest ES features | Elasticsearch | Active feature development |
| Open-source requirement | OpenSearch | Apache 2.0 license |
| Existing Elastic investment | Elasticsearch | Continuity |
| Budget-conscious | OpenSearch | Free security and alerting |

---

## 📚 Topics

| Topic | Description | Key Concepts |
|-------|-------------|-------------|
| [Elasticsearch](elasticsearch/) | Distributed search and analytics engine | Indexing, Query DSL, aggregations, cluster management |

---

## ⚡ Quick Reference

### Common Search Operations

```bash
# Check cluster health
curl -X GET "localhost:9200/_cluster/health?pretty"

# List indices
curl -X GET "localhost:9200/_cat/indices?v"

# Index a document
curl -X POST "localhost:9200/my-index/_doc" \
  -H "Content-Type: application/json" \
  -d '{"title": "Hello", "content": "World"}'

# Search
curl -X GET "localhost:9200/my-index/_search?pretty" \
  -H "Content-Type: application/json" \
  -d '{"query": {"match": {"content": "world"}}}'

# Delete index
curl -X DELETE "localhost:9200/my-index"
```

### Port Defaults

| Service | Port | Description |
|---------|------|-------------|
| Elasticsearch HTTP | 9200 | REST API |
| Elasticsearch Transport | 9300 | Node-to-node communication |
| OpenSearch HTTP | 9200 | REST API |
| Kibana | 5601 | Visualization dashboard |
| OpenSearch Dashboards | 5601 | Visualization dashboard |

---

## 🔗 Related Topics

- [Databases — PostgreSQL](../databases/postgresql/) — Full-text search with `tsvector`/`tsquery` |
- [Distributed Systems](../distributed-systems/) — Distributed architecture concepts
- [Observability](../observability/) — Log aggregation with search engines
- [DevOps — Kubernetes](../devops/kubernetes/) — Running search clusters on Kubernetes
- [CLI — curl](../cli/curl/) — REST API interaction with search engines
- [CLI — jq](../cli/jq/) — JSON response processing

---

> **Search engines are complementary to databases, not replacements.** Use your database as the source of truth and the search engine as a specialized index for discovery, analytics, and full-text search.
