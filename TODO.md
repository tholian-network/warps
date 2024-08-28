
# TODO

**structs**

- [ ] WebCache Exists/Read/Write

**structs/DomainCache**

- [ ] Store A/AAAA/MX entries on filesystem, per domain
- [ ] `Exists(domain, entry_type)`
- [ ] `Read(domain, entry_type)`
- [ ] `Write(domain, entry_type, data)` ?

**structs/Resolver**

- [ ] Resolver can listen to `https`, `dns`, `tcp` protocols
- [ ] Resolve DNS requests locally, via DNS over TLS Ronin
- [ ] Resolve DNS requests through specified tunneling protocol to other peer
- [ ] Always use TTL 0 for DNS responses
- [ ] Lookup requests in Domain Cache
- [ ] `Resolve(domain, entry_type)`?

**structs/Proxy**

- [ ] Proxy can use `https`, `http`, `dns` protocols
- [ ] Proxy can use `icmp`, `tcp` protocols (later)
- [ ] Listen to HTTPS requests if specified as protocol
- [ ] Listen to HTTP requests if specified as protocol
- [ ] Listen to DNS requests if specified as protocol
- [ ] Lookup requests in Web Cache

**compressors**

- [ ] image/jpeg compressor
- [ ] image/png compressor
- [ ] text/html compressor
- [ ] text/css compressor


- [ ] Proxy to accept local requests via HTTP or HTTPS

- [ ] Gateway to accept local requests via specified tunneling protocol
- [ ] Tunnel to send requests via specified tunneling protocol

- [ ] WebProxy for tunneling HTTPS requests through HTTP/HTTPS/DNS tunnels
- [ ] DomainResolver for DNS cache and to tunnel requests through DNS tunnels

