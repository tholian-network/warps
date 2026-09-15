# Protocol Implementation Guide

This guide describes how to add a new network protocol to Tholian Warps by
following the conventions of the existing `dns` / `dnstunnel`, `http` /
`httptunnel`, `https`, `icmp` / `icmptunnel`, and `socks` implementations.

## Step 1: Wire codec (`source/protocols/<name>`)

Create a package that knows how to build and parse the protocol's raw bytes.

- `Packet.go` - the core packet struct, a `NewPacket()` constructor, setters,
  and a `Bytes()` method.
- `Parse.go` - a `Parse([]byte) Packet` function that reads raw bytes into a
  `Packet`.
- Any sub-codecs the protocol needs (for example `Message.go` for the
  pingtunnel `MyMsg` framing used by ICMP).

The codec must be pure and dependency-free. Do not put networking in this
package. Connectionless protocols (like ICMP) also define a `Socket` interface
here so the transport can be tested without raw sockets.

## Step 2: Transport (`source/protocols/<name>tunnel`)

Create a package that implements `interfaces.Tunnel`, `interfaces.Proxy`, and
(optionally) `interfaces.Resolver`. Some packages are named after the scheme
instead (`protocols/https`, `protocols/socks`) and reuse an existing codec.

### `Tunnel.go`

- `type Tunnel struct { Host string; Port uint16; ... }`
- `NewTunnel(host string, port uint16) Tunnel`
- `(tunnel *Tunnel) ResolvePacket(dns.Packet) dns.Packet`
- `(tunnel *Tunnel) RequestPacket(http.Packet) http.Packet`

The Tunnel is the outbound half. It encodes requests into the wire format,
sends them, decodes the response, and returns it. For connectionless protocols
(like ICMP) inject a socket via `SetSocket(...)`.

### `Proxy.go`

- `type Proxy struct { Host, Port, Cache, Tunnel, Resolver, ... }`
- `NewProxy(host string, port uint16, cache interfaces.ProxyCache) Proxy`
- `(proxy *Proxy) ResolvePacket(dns.Packet) dns.Packet`
- `(proxy *Proxy) RequestPacket(http.Packet) http.Packet`
- `(proxy *Proxy) SetResolver(...)`, `SetTunnel(...)`
- `(proxy *Proxy) Destroy() error`
- `(proxy *Proxy) Listen() error`

The `RequestPacket` method follows this order:

1. If the packet is a transport-level DNS query, decode and dispatch to
   `ResolvePacket`.
2. Otherwise, check the `ProxyCache`.
3. Otherwise, delegate to the configured `Tunnel`.
4. Otherwise, execute the request directly (via `Resolver` and/or
   `http.RequestPacket`).

### `Resolver.go`

- `NewResolver(host string, port uint16, cache interfaces.ResolverCache) Resolver`
- `(resolver *Resolver) Resolve(domain string) dns.Packet`
- `(resolver *Resolver) ResolvePacket(dns.Packet) dns.Packet`
- `(resolver *Resolver) SetTunnel(...)`

`Listen()` is optional and only implemented where the resolver is also a
standalone server (e.g. `dnstunnel.Resolver`). The `interfaces.Resolver` type
only requires `Resolve` and `ResolvePacket`.

### `tunnel/` subpackage

For request/response protocols, put the encode/decode helpers in a subpackage,
one function per file:

- `EncodeRequest.go` / `DecodeRequest.go`
- `EncodeResponse.go` / `DecodeResponse.go`
- `IsRequest.go` / `IsResponse.go`
- `EncodeError.go`

Import it with an alias, for example
`icmp_tunnel "tholian-warps/protocols/icmptunnel/tunnel"`.

Stream protocols (like `socks`) do not need this subpackage; they expose the
same interfaces but relay raw bytes after a handshake.

## Step 3: Register the protocol

1. Add the protocol constant to `types.Protocol.go` (`ProtocolX Protocol = "x"`
   plus `IsProtocol` and `ParseProtocol` cases).
2. Add the scheme to `utils/arguments/Config.go` (`case "x":` plus a default
   port in the default-port switch).
3. Add tunnel and/or listen branches to `actions/Tunnel.go`, `actions/Forward.go`,
   and `actions/Gateway.go`.
4. Update the CLI usage text in `cmds/tholian-warps/main.go`.
5. Update `README.md` (network protocols and test coverage lists) and
   `docs/architecture.md`.

## Step 4: Tests

- Unit tests for the codec: encode/decode round-trips.
- End-to-end tests that run a Tunnel against a Proxy (and a Resolver) using
  real loopback sockets (or an in-memory socket for ICMP).
- Use the `test.NewSpy*` helpers for caches, resolvers, and tunnels.
- Follow the `t.Run("...", ...)` naming style and the
  `t.Errorf("Expected '%s' but got '%s'", ...)` assertion style.

For connectionless protocols, add a `LoopbackSocket` to `protocols/test` that
implements the transport socket interface over channels.

## Conventions

- PascalCase file names (`Packet.go`, `ResolvePacket.go`).
- Imports are grouped and aliased for disambiguation
  (`net_url "net/url"`, `dns_tunnel "tholian-warps/protocols/dnstunnel/tunnel"`).
- Struct fields use `json:"..."` tags matching their exported names.
- Zero external dependencies.
