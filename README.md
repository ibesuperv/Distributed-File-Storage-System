# Distributed-File-Storage-System

[![Go Version](https://img.shields.io/badge/Go-1.26+-00ADD8?style=flat&logo=go)](https://golang.org/)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)
[![GitHub Star](https://img.shields.io/github/stars/ibesuperv/Distributed-File-Storage-System?style=flat&logo=github)](https://github.com/ibesuperv/Distributed-File-Storage-System)
[![Architecture](https://img.shields.io/badge/Architecture-P2P%20%7C%20CAS%20%7C%20Streaming-blueviolet)]()

**Distributed-File-Storage-System** is a high-performance, decentralized file storage infrastructure engineered in Go. It solves the "Large File" problem in peer-to-peer networks by combining **Content-Addressable Storage (CAS)**, **Custom TCP Transport**, and **Streaming Cryptography** into a single, cohesive engine.

Designed for scalability and resilience, this system eliminates the need for central authority while providing industrial-grade data integrity and security—mirroring the architectural principles found in large-scale distributed systems like Cloud Object Storage or BitTorrent.

---

## 🏛️ System Architecture

Our design prioritizes **Zero-Buffer I/O** and **Stateless Discovery**. Unlike centralized storage, every node in this network acts as both an indexer and a provider.

```mermaid
graph TD
    subgraph "External API / CLI"
        A[User Request] -->|Upload/Download| B[FileServer Orchestrator]
    end

    subgraph "Internal Engine (Node)"
        B --> C{Orchestrator}
        C -->|Stream| D[Crypto Layer: AES-CTR]
        C -->|Hash & Path| E[CAS Store Engine]
        C -->|Broadcast| F[P2P Transport: TCP]
    end

    subgraph "Distributed Network"
        F <---|"Custom Binary Protocol" ---> G[Peer Cluster]
    end

    D -->|Encrypted Chunk| E
    E -->|O(1) Write| J[(Local Storage: Sharded Hierarchy)]
```

---

## 🚀 Engineering Deep Dives

### 1. Content-Addressable Storage (CAS)
Traditional storage uses file names. This system uses **Mathematical Identity**.
- **The Process**: Files are hashed using SHA1. The hash *is* the pointer.
- **The Value**: 
    - **Self-Verifying Integrity**: The data address is a cryptographic proof of its content.
    - **Global Deduplication**: Identical data across the network occupies exactly one address.
- **Scalability**: Uses a **Sharded Directory Hierarchy** (e.g., `/af32d/12c3b/...`) to prevent directory "hotspotting" and filesystem degradation.

### 2. High-Performance TCP Wire Protocol
Instead of standard HTTP, we implemented a custom, lightweight TCP protocol designed for **Large Object Transfers**.
- **Stateful Decoding**: Our `DefaultDecoder` peeks at the wire using a `Switch-Case` strategy, allowing it to transition seamlessly between **GOB-encoded metadata** and **raw binary streams** without resetting connections.
- **Multiplexed Logic**: Pause/Resume mechanics allow nodes to process control messages while data streams in the background.

### 3. Handling "Massive" Files ($O(1) Memory$)
A core requirement for production systems is that memory usage must not scale with file size.
- **Streaming Pipeline**: Uses Go's `io.Reader` and `io.Writer` interfaces throughout. Data moves through the crypto engine and out to the disk in small chunks (32KB buffers).
- **Constant Memory**: Whether you store a 10MB photo or a 100GB 8K video, the node's memory footprint remains nearly constant.

### 4. Cryptographic Security
- **AES-CTR (Counter Mode)**: Chosen for its parallelizable nature and compatibility with streaming data (no padding required).
- **IV Salt Propagation**: Every file transaction generates a unique Initialization Vector (IV), prepended to the data to prevent frequency analysis attacks.

---

## 📈 Scalability & Performance Metrics

- **Write Complexity**: $O(1)$ directory lookup via sharded CAS pathing.
- **Read Latency**: Minimized through parallel discovery broadcast.
- **Replication**: Fully decentralized replication across all interconnected nodes.

---

## 🚦 Getting Started

### 1. Installation
```bash
git clone https://github.com/ibesuperv/Distributed-File-Storage-System
cd Distributed-File-Storage-System
go build -o bin/dfs
```

### 2. Launch Local Cluster (Simulation)
```bash
# Start nodes and upload a sample file
make run ARGS="-u test_files/audio.mpeg"

# Download and verify integrity
make run ARGS="-d audio.mpeg"
```

---

## 🧪 Engineering Philosophy

- **Scalability**: Designed to handle arbitrarily large files without memory spikes.
- **Resilience**: A decentralized "Shared Nothing" architecture where any node can fail without data loss.
- **Observability**: Consistent logging and clear separation of transport vs. storage concerns.
- **Reliability**: Self-correcting stream logic with proper synchronization.

---

## 🛠️ Project Roadmap
Next-phase implementation targets:
- **Distributed Hash Table (DHT)**: Scoping for $O(\log n)$ node discovery.
- **Erasure Coding**: Resilience planning for multi-node failure scenarios.
