# Wire Protocol Specification

This document defines the custom P2P wire protocol used by the **Distributed-File-Storage-System**. 

## 1. Frame Structure

The protocol uses a **Type-Length-Value (TLV)** inspired framing system optimized for Go's `io.Reader/Writer` patterns.

### Frame Anatomy:
| Byte Offset | Field | Type | Description |
|---|---|---|---|
| 0 | **Type Identifier** | 1 byte | Indicates the nature of the following payload (Message vs Stream). |
| 1+ | **Payload** | Body | Either a GOB-encoded message or a raw binary stream. |

---

## 2. Type Identifiers

To prevent expensive full-packet buffering, we use a single-byte prefix to switch the state of the `Decoder`.

| Constant | Value | Description |
|---|---|---|
| `IncomingMessage` | `0x1` | The following payload is a control message. |
| `IncomingStream` | `0x2` | The following bytes are a raw data stream. |

---

## 3. The Multi-Phase Decoding Logic

### Phase 1: The Peek
The `DefaultDecoder` reads only **one byte** from the `net.Conn`.
- If the byte is `IncomingStream`, the `RPC` object is flagged as `Stream = true` and returned immediately. This pauses the control loop and allows the binary stream to be piped directly to disk.

### Phase 2: Implementation-Specific Decoding
If the byte is `IncomingMessage`, the decoder reads the remaining bytes into a buffer and populates the `RPC.Payload`. This payload typically contains GOB-encoded structs:

#### `MessageStoreFile`:
- Used to signify a new file upload to the network.
- Contains: `ID`, `Key`, and `Size`.

#### `MessageGetFile`:
- Used to request a file.
- Contains: `ID` and `Key`.

---

## 4. Why custom TCP vs HTTP?

1. **Overhead**: HTTP headers (User-Agent, Accept, etc.) add significant bloat to small control messages in a P2P network.
2. **Streaming Control**: HTTP/1.1 is strictly Request-Response. Our custom TCP protocol allows for **Bidirectional Streaming** where a peer can send a control message and immediately follow it with 100GB of raw data without closing the connection or renegotiating headers.
3. **Simplicity**: By working directly with `net.Conn`, we leverage Go's internal buffer management more efficiently than high-level web frameworks.
