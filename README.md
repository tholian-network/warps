
# Tholian® Warps

Tholian® Warps is a proactive and adaptive Mesh Network Router.

It tries to automatically detect and bypass censorship and throttling
measurements at all cost. Programmable routing, traffic compression,
traffic scattering, dynamic encryption rotation and other features
are part of this experimental research project.

The goal of this project is to find out how feasible common NAT breaking
and firewall bypassing techniques are and whether they can be used to
build a reliable mesh network that's based on a peer-to-peer architecture.

### Network Architecture

![network-architecture.png](https://github.com/tholian-network/warps/blob/master/assets/network-chart.png?raw=true)

1. A `tunnel` instance tunnels network traffic through a `forward` or `gateway` instance to access the internet.
2. A `tunnel` uses the initial configured network protocol.
3. A `forward` instance tunnels network traffic through other `forward` or `gateway` instances.
4. All `gateway` instances use `optimizers` to reduce web asset file sizes.
5. All `gateway` instances use `DNS over TLS` to resolve relayed DNS queries.
6. All instances use a local [ProxyCache](./source/structs/ProxyCache.go) and [ResolverCache](./source/structs/ResolverCache.go).
7. All instances can rotate encryption keys and can scatter network traffic on-demand.

### Building

:construction: Highly Experimental at this point. Use software at your own risk! :construction:

```bash
bash build.sh;
sudo cp ./build/tholian-warps /usr/bin/tholian-warps;

# Show CLI usage help
tholian-warps;
```

### Usage: Tunnels and Gateways

The easiest way to use Warps is with running a Warps `gateway` on your own VPS that is connected to the internet,
and a locally running Warps `tunnel`.

As a defaulted network protocol, it is best to use `dns`, as that usually works to bypass typical firewall setups.
Alternative supported network protocols are documented further down in this document.

```bash
# On your VPS server (1.3.3.7)
tholian-warps gateway "dns://0.0.0.0:1053";

# On your local machine
tholian-warps tunnel "any" "dns://1.3.3.7:1053";
curl -x localhost:1080 http://google.com;
```

## Usage: Proxy Chains

Warps can be chained via multiple proxies, without a limit on how many network hops you want to the public internet.
In this example, we are routing local web traffic through 3 instances before the traffic hits the clearnet.

```bash
# On your first VPS server (1.3.3.7)
tholian-warps gateway "dns://0.0.0.0:1337";

# On your second VPS server (1.3.3.8)
tholian-warps forward "http://1.3.3.8:1338" "dns://1.3.3.7:1337";

# On your third VPS server (1.3.3.9)
tholian-warps forward "dns://1.3.3.9:1339" "http://1.3.3.8:1338";

# On your local machine
# local -> dns -> 1.3.3.9 -> http -> 1.3.3.8 -> dns -> 1.3.3.7 -> * -> internet
tholian-warps tunnel "any" "dns://1.3.3.9:1339";
curl -x localhost:1080 http://google.com;
```

# Data Compressors

These are the data compressors that have been implemented:

- [ ] [compressors/text/css](/source/compressors/text/css)
- [ ] [compressors/text/html](/source/compressors/text/html)
- [ ] [compressors/text/js](/source/compressors/text/js)
- [ ] [compressors/image/jpeg](/source/compressors/image/jpeg)
- [ ] [compressors/image/png](/source/compressors/image/png)


# Network Protocols

These are the transport protocols that have been implemented:

- [x] `dns` or [protocols/dnstunnel](/source/protocols/dnstunnel) implements DNS Exfiltration
- [x] `http` or [protocols/httptunnel](/source/protocols/httptunnel) implements HTTP Smuggling and DNS over HTTP
- [x] `https` or [protocols/https](/source/protocols/https) implements HTTPS Smuggling and DNS over HTTPS
- [x] `icmp` or [protocols/icmptunnel](/source/protocols/icmptunnel) implements ICMP Knocking and DNS over ICMP
- [ ] `ssh` or [protocols/ssh](/source/protocols/ssh) implements SSH Tunneling and DNS over SSH
- [x] `socks` or [protocols/socks](/source/protocols/socks) implements SOCKS Routing (e.g. for TOR/I2P usage)
- [x] [protocols/test](/source/protocols/test) implements the `Spy` testing data structures

The `dns`, `http` and `icmp` wire-format codec libraries are vendored in
[protocols/dns](/source/protocols/dns), [protocols/http](/source/protocols/http) and
[protocols/icmp](/source/protocols/icmp), with their shared helpers in
[types](/source/types) and [utils](/source/utils).


# Test Coverage

These are the `go test` files that have been implemented:

- [x] [types/ASN](/source/types/ASN_test.go)
- [x] [types/Datetime](/source/types/Datetime_test.go)
- [x] [types/Domain](/source/types/Domain_test.go)
- [x] [types/IPv4](/source/types/IPv4_test.go)
- [x] [types/IPv6](/source/types/IPv6_test.go)
- [x] [types/Time](/source/types/Time_test.go)
- [x] [structs/ProxyCache](/source/structs/ProxyCache_test.go)
- [x] [structs/ResolverCache](/source/structs/ResolverCache_test.go)
- [ ] [utils/net/url/IsTrackingParameter](/source/utils/net/url/IsTrackingParameter_test.go)
- [ ] [utils/net/url/IsXSSParameter](/source/utils/net/url/IsXSSParameter_test.go)
- [x] [utils/net/url/ResolveCache](/source/utils/net/url/ResolveCache_test.go)
- [x] [utils/net/url/ToHostAndPort](/source/utils/net/url/ToHostAndPort_test.go)
- [x] [utils/net/url/ToHost](/source/utils/net/url/ToHost_test.go)
- [ ] [utils/protocols/http/IsFilteredHeader](/source/protocols/http/IsFilteredHeader_test.go)
- [x] [protocols/dns/Type](/source/protocols/dns/Type_test.go)
- [x] [protocols/dns/Class](/source/protocols/dns/Class_test.go)
- [x] [protocols/dns/OperationCode](/source/protocols/dns/OperationCode_test.go)
- [x] [protocols/dns/ResponseCode](/source/protocols/dns/ResponseCode_test.go)
- [x] [protocols/dns/Question](/source/protocols/dns/Question_test.go)
- [x] [protocols/dns/Record](/source/protocols/dns/Record_test.go)
- [x] [protocols/dns/Packet](/source/protocols/dns/Packet_test.go)
- [x] [protocols/dns/toPTRName](/source/protocols/dns/toPTRName_test.go)
- [x] [protocols/http/Parse](/source/protocols/http/Parse_test.go)
- [x] [protocols/http/Request](/source/protocols/http/Request_test.go)
- [x] [protocols/icmp/Message](/source/protocols/icmp/Message_test.go)
- [x] [protocols/icmp/Packet](/source/protocols/icmp/Packet_test.go)
- [x] [protocols/icmp/Socket](/source/protocols/icmp/Socket_test.go)
- [x] [protocols/icmptunnel/Resolver](/source/protocols/icmptunnel/Resolver_test.go)
- [x] [protocols/icmptunnel/Proxy](/source/protocols/icmptunnel/Proxy_test.go)
- [x] [protocols/icmptunnel/tunnel/EncodeRequest](/source/protocols/icmptunnel/tunnel/EncodeRequest_test.go)
- [x] [protocols/icmptunnel/tunnel/EncodeResolveRequest](/source/protocols/icmptunnel/tunnel/EncodeResolveRequest_test.go)
- [x] [protocols/dnstunnel/Resolver](/source/protocols/dnstunnel/Resolver_test.go)
- [x] [protocols/dnstunnel/Proxy](/source/protocols/dnstunnel/Proxy_test.go)
- [x] [protocols/dnstunnel/Tunnel](/source/protocols/dnstunnel/Tunnel_test.go)
- [x] [protocols/dnstunnel/tunnel/ToRecordName](/source/protocols/dnstunnel/tunnel/ToRecordName_test.go)
- [x] [protocols/httptunnel/Proxy](/source/protocols/httptunnel/Proxy_test.go)
- [x] [protocols/httptunnel/Tunnel](/source/protocols/httptunnel/Tunnel_test.go)
- [x] [protocols/https/Proxy](/source/protocols/https/Proxy_test.go)
- [x] [protocols/https/Tunnel](/source/protocols/https/Tunnel_test.go)
- [x] [protocols/socks/Proxy](/source/protocols/socks/Proxy_test.go)
- [x] [protocols/socks/Tunnel](/source/protocols/socks/Tunnel_test.go)


# License

AGPL3
