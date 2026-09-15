# Architecture Overview

## Network topology

```
[tunnel] --(initial protocol)--> [forward] --(next protocol)--> [gateway] --> internet
```

- A `tunnel` instance tunnels network traffic through a `forward` or `gateway`
  instance to access the internet.
- A `forward` instance tunnels network traffic through other `forward` or
  `gateway` instances.
- A `gateway` instance terminates the tunnel and executes real HTTP/HTTPS/DNS
  traffic against the internet.
- Every instance uses a local `ProxyCache` and `ResolverCache`.

## Supported protocols

The implemented transport protocols and their wire-format codecs:

| Scheme  | Transport package            | Wire codec        | Notes                                   |
|---------|------------------------------|-------------------|-----------------------------------------|
| `dns`   | `protocols/dnstunnel`        | `protocols/dns`   | DNS exfiltration                        |
| `http`  | `protocols/httptunnel`       | `protocols/http`  | HTTP smuggling + DNS over HTTP          |
| `https` | `protocols/https`            | `protocols/http`  | HTTPS smuggling + DNS over HTTPS        |
| `icmp`  | `protocols/icmptunnel`       | `protocols/icmp`  | ICMP knocking + DNS over ICMP           |
| `socks` | `protocols/socks`            | -                 | SOCKS5 routing (e.g. TOR/I2P)           |

## Protocol layering

Warps separates most network protocols into two layers:

1. **Wire codec** (`protocols/<name>`) - pure, hand-written, dependency-free
   encode/parse logic for the on-wire format: `protocols/dns`, `protocols/http`,
   and `protocols/icmp`.
2. **Transport** (`protocols/<name>tunnel`, or `protocols/https` / `protocols/socks`)
   - implements the `Proxy` / `Tunnel` / `Resolver` interfaces on top of the
   codec. The package naming varies: `dnstunnel`, `httptunnel`, `icmptunnel`,
   `https`, and `socks`.

`socks` is the exception: it is a raw stream protocol, so it has no dedicated
`tunnel/` codec subpackage. It performs a SOCKS5 handshake and then relays raw
bytes, still exposing itself through the shared `ResolvePacket` /
`RequestPacket` interfaces.

## Interfaces

Defined in `source/interfaces`:

- `Tunnel` - `ResolvePacket(dns.Packet) dns.Packet` and
  `RequestPacket(http.Packet) http.Packet`. The outbound half of a transport.
- `Proxy` - a `Tunnel` plus `Destroy() error` and `Listen() error`. The inbound
  half that accepts connections and serves requests.
- `Resolver` - `Resolve(string) dns.Packet` and
  `ResolvePacket(dns.Packet) dns.Packet`.
- `ProxyCache` / `ResolverCache` - storage backends for web responses and DNS
  records.

## Request / response flow

Every transport follows the same flow, regardless of the underlying wire
protocol:

1. A `Tunnel.RequestPacket` encodes an `http.Packet` into the transport's wire
   format and sends it to the remote instance.
2. A `Proxy.Listen` loop receives the wire packet, decodes it back into an
   `http.Packet`, and dispatches it through `RequestPacket`.
3. `Proxy.RequestPacket` checks the `ProxyCache`, then delegates to its
   configured `Tunnel` or `Resolver`, or executes the request directly via
   `http.RequestPacket`.
4. The response is encoded back into the wire format and returned.

DNS queries follow the identical pattern through `ResolvePacket`.

## Data structures

- `types.Server` describes a remote endpoint (domain, addresses, port,
  protocol, schema).
- `types.Protocol` enumerates the supported protocols (`dns`, `http`, `https`,
  `icmp`, `socks`, `ssh`, `whois`, ...).
- `arguments.Config` parses `protocol://host:port` URLs into a typed config and
  applies per-protocol default ports.

## ICMP specifics

ICMP is a connectionless datagram protocol, so the transport uses a small
`socket` abstraction instead of `net.ListenUDP` / `net.ListenTCP`:

```go
type Socket interface {
    Listen(host string) error
    ReadPacket() (Packet, net.Addr, error)
    WritePacket(Packet, net.Addr) error
    Close() error
}
```

The production implementation (`PingSocket`) uses an unprivileged ICMP
datagram socket (Linux `SOCK_DGRAM` + `IPPROTO_ICMP`), which requires root or
the kernel's `net.ipv4.ping_group_range` sysctl to allow it. Tests inject an
in-memory `LoopbackSocket` pair so end-to-end ICMP tests run without root.

The ICMP echo data carries a pingtunnel `MyMsg` (protobuf framing, hand-written
to stay dependency-free). HTTP requests use `target = request URL` and DNS
queries use `target = "dns-query"` so the receiver can distinguish the two
payload kinds.

## SOCKS specifics

SOCKS is a TCP stream protocol. The `socks.Tunnel` acts as a SOCKS5 client: it
performs the greeting/method selection and a `CONNECT` to the target, then
`RequestPacket` writes the HTTP request and reads the response over the
established stream. The `socks.Proxy` acts as a SOCKS5 server: it accepts a
connection, performs the handshake, dials the requested target, and relays
bytes bidirectionally.

## Testing

- Codec unit tests cover encode/decode round-trips (`protocols/dns`,
  `protocols/icmp`, `protocols/icmptunnel/tunnel`).
- End-to-end tests run a `Tunnel` against a `Proxy` (and optionally a
  `Resolver`) over real loopback sockets, or over `LoopbackSocket` for ICMP.
- `protocols/test` provides `Spy` implementations of the caches, resolver, and
  tunnel for deterministic assertions.
