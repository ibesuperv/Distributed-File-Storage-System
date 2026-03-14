# Crypto & Security Specification

Security in a distributed system is non-negotiable. This document explains the implementation of the streaming cryptographic layer.

## 1. The Threat Model

Our security layer protects against:

- **Passive Sniffing**: Files in transit over the network are encrypted.
- **Disk Theft**: Files stored on a node's physical disk are "at-rest" encrypted.
- **Data Tampering**: The Content-Addressable nature (SHA1) ensures that any bit-level modification of the encrypted block is detectable (integrity check).

---

## 2. AES-256 CTR Implementation

We chose **AES in Counter Mode (CTR)** for its efficiency in streaming environments.

### Why AES-CTR?

- **No Padding**: Unlike CBC mode, CTR does not require data to be a multiple of the block size. This is crucial for streaming arbitrary file sizes.
- **Parallelizable**: Each block can be decrypted independently, which is a significant performance win for large files.
- **Random Access**: We can seek to any part of the file and begin decryption without processing the whole stream.

---

## 3. Initialization Vector (IV) Management

An IV is essential to ensure that identical files result in different ciphertexts.

### Generation & Persistence:

1. For every `Store` operation, a **16-byte cryptographically secure random IV** is generated.
2. The IV is written as the **first 16 bytes** of the file on disk.
3. During `Read`, the first 16 bytes are pulled, the AES cipher is seeded, and the remainder of the file is streamed through the XOR engine.

---

## 4. Key Management

The current implementation uses a **32-byte (256-bit)** master key.

> [!NOTE]
> For production environments, this key should be derived from a Key Management Service (KMS) or a secure Vault. In the current simulation, it is provided via `FileServerOpts`.

---

## 5. Streaming Pipeline (`copyStream`)

The magic happens in `crypto.go`:

```go
func copyStream(stream cipher.Stream, blockSize int, src io.Reader, dst io.Writer)
```

This function reads from an `io.Reader` into a 32KB buffer, applies the `XORKeyStream` transformation, and writes to an `io.Writer`. This loop continues until `io.EOF`, ensuring that memory usage never exceeds the buffer size regardless of file size.
