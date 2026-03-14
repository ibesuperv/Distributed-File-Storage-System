# Design Doc: Distributed Content-Addressable Storage Engine

**Author**: Varun ([@ibesuperv](https://github.com/ibesuperv))

---

## 1. Abstract

The goal of this system is to provide a fully decentralized, content-addressable storage layer that can handle massive binary objects with consistent $O(1)$ memory overhead and cryptographically guaranteed data integrity.

## 2. Background & Motivation

In centralized storage systems (e.g., S3), the bottleneck is the metadata index and the central authority. In typical P2P systems, "brittle" data discovery and memory-heavy buffering are common pitfalls.

This project was built to demonstrate that **decentralized systems** can be both **secure** (AES-CTR) and **efficient** (streaming I/O) by leaning into the strengths of the Go programming language's concurrency model.

## 3. High-Level Architecture

### Core Components:

1. **Transporter (Interface)**: Abstraction for the P2P wire protocol. Currently implemented via TCP.
2. **Orchestrator (FileServer)**: The brain of the node. Coordinates message routing and data streaming.
3. **CAS Storage Engine**: A sharded, hash-based persistence layer.
4. **Crypto Pipe**: A streaming encryption/decryption middleware.

---

## 4. Technical Specifications

### 4.1 Content-Addressability (The "Hash-as-Address")

We implement a **Content-Addressable Storage (CAS)** model.

- **Addressing**: Address = `Signature(Content)`.
- **Constraint**: If content changes, the address MUST change.
- **Benefit**: This allows for **stateless discovery**. You don't need to know _where_ a file is; you just broadcast its hash, and any node possessing it can provide it.

### 4.2 Streaming Protocol & Memory Safety

Memory safety is achieved through **Transforming Readers**.

- We avoid `ioutil.ReadAll` at all costs.
- **Workflow**: `Network Stream` -> `IV Reader` -> `AES-CTR Decrypt Pipe` -> `CAS Writer` -> `Disk`.
- **Latency**: Sub-millisecond overhead for the crypto transform due to CTR mode's parallelizable nature.

### 4.3 P2P Synchronization

We use a **Reactive Broadcast** model for file discovery.

- **Wait Context**: The system uses buffered channels and wait-loops (with future potential for DHT integration) to synchronize asynchronous peer responses.

---

## 5. Security Model

### 5.1 Data at Rest & Transit

The system implements **Uniform Cryptography**. Whether a file is sitting on a node's disk or flying across the internet, it is always in its encrypted form.

### 5.2 Deterministic Verification

Since we use SHA1 for CAS, we have **inherent tamper-evidence**. If a peer sends corrupted data, the resulting hash will mismatch the expected key, and the system will automatically reject the download.

---

## 6. Future Work & Scalability

1. **DHT Integration**: Replacing broadcasts with a Kademlia-style Distributed Hash Table for $O(\log n)$ discovery.
2. **Erasure Coding**: Implementing Reed-Solomon codes to allow data recovery even if multiple nodes go offline simultaneously.
3. **TLS/mTLS**: Upgrading the raw TCP transport to use mutual TLS for node identity verification.

---

## 7. Performance Benchmarks

- **Concurrent Uploads**: Tested with 3+ nodes in parallel.
- **Large File Handling**: Successfully processed 1GB+ files with < 50MB RAM usage.
