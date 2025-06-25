# Telemetry Agent Documentation

This document outlines the design, rationale, and architecture of the telemetry software currently being developed. Its primary role is to serve as an **Intrusion Detection System (IDS)**, with room for future expansion into an **Intrusion Prevention System (IPS)**.

---

## 1. 🎯 Purpose and Scope

The software is designed to monitor and analyze network activity across multiple protocols (starting with HTTP and ICMP). With further optimization and deeper packet inspection capabilities, it can later be extended to monitor:

- SSH (Secure Shell)
- SMTP (Simple Mail Transfer Protocol)
- DNS, FTP, and other critical services

The system is split into modular components to allow clean separation of concerns and operational scalability.

---

## 2. 🚀 Why Go?

### ✅ Advantages of Go:

- **Optimized Concurrency**: Native support for goroutines and channels is ideal for network-heavy and concurrent tasks like reverse proxying and traffic analysis.
- **Stability**: Go's static type system and built-in safety features reduce the risk of critical runtime errors.
- **Performance**: Offers good balance between execution speed and developer productivity.
- **Ecosystem**: Many modern security tools (e.g., CrowdSec, Caddy) are written in Go, making it a natural choice for integration and hiring.

### ❌ Why not C++?

- **Overpowered for HTTP Use Case**: C++ is too low-level for the initial telemetry requirements.
- **Risk of Crashes**: Its memory-unsafe features mean a small bug (e.g., CrowdStrike’s infamous crash incident) could take down systems.
- **Onboarding Difficulty**: The average security engineer is far more likely to be productive in Go than in C++, which is more commonly found in high-performance or game development roles.
- **Maintenance Burden**: Code safety, testing, and onboarding are significantly more complex in C++.

> 🧠 **Note**: C++ will still be used in targeted components for memory-level analysis (e.g., malware detection or memory forensics), but is not the core agent language.

---

## 3. 🏗 System Architecture

### 🔹 Host 1: Application Backend

- Standard web application (HTTP API)
- The final recipient of proxied requests
- Can be any web app/service the agent sits in front of

### 🔹 Host 2: Telemetry Agent

- Reverse proxy + packet tracker
- Configurable via `config.yaml`
  - `listen_host`, `listen_port`
  - `forward_host`, `forward_port`
- Performs in-memory telemetry aggregation (per-path stats)
- Buffers data for 60 seconds (or configurable interval)
- Flushes to the data collector endpoint
- Implements retry and fallback strategy for resiliency:
  - If collector is unreachable, retry up to 5 times (e.g., 10s → 30s → 1min → 3min → 5min)
  - On final failure, dump unsent data to a local file (`fallback-queue.jsonl`) with timestamp
  - On next successful connection, resend buffered batches and wipe the cache
  - Data is not stored permanently — maximum retention is expected to be under 24 hours (as infrastructure will be restored quickly)

### 🔹 Host 3: Central Data Collector

- Secured HTTPS server
- Accepts batched telemetry JSON
- Analyzes for:
  - Abnormal path access
  - Suspicious frequency spikes
  - Malformed requests
- Stores history and visualizations for blue team use

---

## 4. 📈 Statistical Metrics

- Count per route (e.g., `/ping` accessed 1,000 times/minute)
- Timestamp bucketing (minute-based time windows)
- Client metadata (optional future addition: IP, User-Agent)
- Aggregated data sent every `flush_interval_seconds`

> This allows for effective trend detection, anomaly scoring, and early-warning indicators of compromise.

---

## 5. 🔐 Security Model

- Agents are issued **private keys** to authenticate to the data collector
- On startup, agent requests a **session token** from the central server
- All communication is over **HTTPS** only
- Central collector supports **whitelisting** per agent ID
- Config fetch and telemetry flush are authorized per-session

---

## 6. 🔧 Configuration

The `config.yaml` file will include:

```yaml
backend:
  listen_host: "127.0.0.1"
  listen_port: "8080"
  forward_host: "127.0.0.1"
  forward_port: "8081"
telemetry:
  flush_interval_seconds: 60
  collector_host: "https://collector.example.com"
  api_key: "<secure-api-key>"
```

This configuration is **internal-only**, and should not be exposed or modifiable by untrusted clients.

---

## 7. 📦 Future Enhancements

- Remote configuration pull via signed JWT tokens
- Add path whitelisting/blacklisting logic
- Include real-time alert pushback (e.g., webhook callback)
- Extend tracking beyond HTTP: DNS, TLS handshake analysis, SMTP
- Automatic backoff / retry queue for offline batch flushing

---

## ✅ Summary

This telemetry agent is designed for secure, scalable IDS-style monitoring with a developer-friendly language (Go), clear system boundaries, and room for defensive or even preventive features (IPS). Modularity allows for each part (agent, backend, collector) to evolve independently.

It minimizes risk via in-memory buffering, controlled config, and tokenized communication — suitable for production-grade environments with internal observability and forensic capabilities.