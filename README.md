current status: phase 3 in progress

# 🏗️ Phase 1 — Single Node Runtime (MVP)

**Goal:** Build a platform that can deploy and manage containers on **one machine**.

This is what we're currently building.

### Components

```text
Models
Registry
Runtime Interface
Docker Runtime
Orchestrator
REST API
```

### Features

* ✅ Service model
* ✅ In-memory registry
* ✅ Runtime abstraction
* ✅ Docker SDK integration
* ✅ Deploy container
* ✅ Stop container
* ✅ Delete container
* ✅ View logs
* ✅ Basic REST API

Architecture

```text
Client

↓

REST API

↓

Orchestrator

↓

Docker Runtime

↓

Docker SDK

↓

Docker Daemon
```

---

# 🌐 Phase 2 — Multi-Node Cluster

Now we move from

```text
1 Machine
```

to

```text
Many Machines
```

### New Components

```text
Node Agent

Scheduler

Heartbeat Service
```

Architecture

```text
          Orchestrator

        /      |      \

    Agent    Agent    Agent
```

Features

* Node registration
* Heartbeats
* CPU/RAM monitoring
* Scheduling
* Deploy on best node

---

# 🔒 Phase 3 — Secure Private Mesh Network

Machines should communicate securely.

We'll build

```text
WireGuard Mesh
```

Features

* Private IPs
* Encrypted traffic
* Automatic peer discovery
* Node authentication

Architecture

```text
Laptop

═══════

WireGuard

═══════

Cloud VPS
```

---

# ⚡ Phase 4 — Serverless Functions

Instead of only Docker containers,

support

```text
WASM

Firecracker
```

Features

* Scale to zero
* Millisecond startup
* Runtime selection
* Cold starts

Architecture

```text
Runtime

├── Docker

├── WASM

└── Firecracker
```

---

# 🧠 Phase 5 — Intelligent Orchestrator

Now the orchestrator becomes smart.

Instead of

```text
Deploy Anywhere
```

it decides

```text
Best Machine
```

Features

* Resource-aware scheduling
* Bin packing
* Least-loaded scheduling
* Affinity rules
* Anti-affinity

---

# 🌍 Phase 6 — Networking & Gateway

Expose services to users.

Features

* Reverse proxy
* Load balancer
* Automatic routing
* Custom domains
* HTTPS

Architecture

```text
Internet

↓

Gateway

↓

Service A

↓

Node 2
```

---

# 📊 Phase 7 — Observability

Professional cloud platforms always include monitoring.

Features

* Live logs
* Metrics
* CPU graphs
* Memory graphs
* Distributed tracing
* Dashboard

Architecture

```text
Containers

↓

Metrics

↓

Dashboard
```

---

# 👥 Phase 8 — Authentication & Multi-Tenancy

Now Vortex becomes a platform for teams.

Features

* Login
* JWT Authentication
* RBAC
* Workspaces
* Team quotas
* Resource limits

---

# 🚀 Phase 9 — GitOps & CI/CD

Developer experience.

Instead of

```bash
docker run
```

users simply

```bash
git push
```

Flow

```text
Git Push

↓

Webhook

↓

Build

↓

Deploy

↓

Traffic Switch
```

Features

* Git deployments
* Zero downtime
* Rollbacks
* Version history

---

# ☁️ Phase 10 — Production Cloud Platform

This is the "wow" phase.

Features

* High Availability Orchestrator
* Raft Consensus
* Distributed Registry
* Auto Scaling
* Self Healing
* Rolling Updates
* Canary Deployments
* Image Signing
* eBPF Security
* Distributed Storage

Architecture

```text
          Leader

        /    |    \

Follower Follower Follower

        ↓

Distributed Cluster
```

---

# Final Architecture

```text
                Dashboard
                     │
                     ▼
                REST API
                     │
                     ▼
               Orchestrator
                     │
      ┌──────────────┼──────────────┐
      ▼              ▼              ▼
 Registry       Scheduler      Runtime
      │              │              │
      │              ▼              ▼
      │          Node Agents    Docker/WASM
      │              │
      ▼              ▼
 Distributed Nodes (Cluster)
```

---

# Learning Progression

Notice how each phase introduces new computer science concepts:

| Phase  | Main Concept                   | Skills You'll Learn                                            |
| ------ | ------------------------------ | -------------------------------------------------------------- |
| **1**  | Single-node container platform | Go, Docker SDK, REST APIs, interfaces, project architecture    |
| **2**  | Distributed systems            | gRPC, node agents, scheduling, heartbeats                      |
| **3**  | Networking                     | VPNs, WireGuard, service discovery, secure communication       |
| **4**  | Serverless runtimes            | WASM, Firecracker, isolation, runtime abstraction              |
| **5**  | Scheduling algorithms          | Resource management, placement strategies                      |
| **6**  | Reverse proxies                | Load balancing, routing, TLS, gateways                         |
| **7**  | Observability                  | Metrics, tracing, logging, monitoring                          |
| **8**  | Security                       | Authentication, authorization, multi-tenancy                   |
| **9**  | DevOps                         | GitOps, CI/CD, rolling deployments                             |
| **10** | Cloud architecture             | Consensus, high availability, self-healing, production systems |

---

