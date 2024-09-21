
# Tholian® Warps

Tholian® Warps is a proactive and adaptive Mesh Network Router.

It tries to automatically detect and bypass censorship and throttling
measurements at all cost. Programmable routing, traffic compression,
traffic scattering, dynamic encryption rotation and other features
are part of this experimental research project.

The goal of this project is to find out how feasible common NAT breaking
and firewall bypassing techniques are and whether they can be used to
build a reliable mesh network that's based on a peer-to-peer architecture.


# Network Architecture

1. A Warps `tunnel` tunnels network traffic through a `peer` or `gateway` to access the internet.
2. A Warps `tunnel` uses the initial configured network protocol.
3. A Warps `peer` tunnels network traffic through other `peer` or `gateway` instances.
4. A Warps `peer` rotates the used network protocols.

5. All Warps `peer` or `gateway` instances use `optimizers` to reduce web asset file sizes.
6. All Warps `peer` or `gateway` instances use `DNS over TLS` to resolve relayed DNS queries.
7. All Warps instances use a local [DomainCache](./source/structs/DomainCache.go) and [WebCache](./source/structs/WebCache.go).
8. All Warps instances can rotate encryption keys and can scatter network traffic on-demand.


## How to use Tunnels and Gateways

The easiest way to use Warps is with running a Warps `gateway` on your own VPS that is connected to the internet,
and a locally running Warps `tunnel`.

As a defaulted network protocol, it is best to use `dns`, as that usually works to bypass typical firewall setups.
Alternative supported network protocols are documented further down in this document.

![network-architecture.png](https://github.com/tholian-network/warps/blob/master/assets/network-chart.png?raw=true)

```bash
# On your VPS server
tholian-warps gateway --protocol=dns --port=1053;

# On your local machine
tholian-warps tunnel --protocol=dns --host=1.3.3.7 --port=1053;
curl -x localhost:8080 http://google.com;
```

## How to use Peers

Warps can be used peer-to-peer and discover locally and globally running Warps `gateway` instances via multicast DNS-SD.
This mode currently needs a Warps `gateway` or `peer` instance that is known among all peers, so that keys can be exchanged.

Note that a Warps `peer` instance can change network protocols on-demand, and that behavior is different from the `tunnel` mode
which prevails the given initial network protocol.

```bash
# On your VPS server
tholian-warps gateway --protocol=dns --port=1053;

# On your local machine
tholian-warps peer --protocol=dns --host=1.3.3.7 --port=1053; # start a local peer, and exchange public keys
curl -x localhost:8080 http://google.com;
```


## How to use Proxy Chains

Warps can be chained via multiple proxies, without a limit on how many network hops you want to the public internet.
In this example, we are routing local web traffic through 3 instances before the traffic hits the clearnet.

```bash
# On your VPS server
tholian-warps gateway --protocol=dns --port=1053;

# On your second VPS server
tholian-warps peer --protocol=dns --host=1.3.3.7 --port=1053;

# On your third VPS server
tholian-warps peer --protocol=dns --host=1.3.3.8 --port=1053;

# On your local machine
# local -> 1.3.3.9 -> 1.3.3.8 -> 1.3.3.7 -> internet
tholian-warps tunnel --protocol=dns --host=1.3.3.9 --port=1053;
curl -x localhost:8080 http://google.com;
```


# Usage

:construction: Highly Experimental at this point - Use at own risk! :construction:

```bash
bash build.sh;
sudo cp ./build/tholian-warps /usr/bin/tholian-warps;

# Show CLI usage help
tholian-warps;
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

- [x] `dns` or [protocols/dns](/source/protocols/dns) implements DNS Exfiltration
- [x] `http` or [protocols/http](/source/protocols/http) implements HTTP Smuggling and DNS over HTTP
- [ ] `https` or [protocols/https](/source/protocols/https) implements HTTPS Smuggling and DNS over HTTPS
- [ ] `icmp` or [protocols/icmp](/source/protocols/icmp) implements ICMP Knocking and DNS over ICMP
- [ ] `ssh` or [protocols/ssh](/source/protocols/ssh) implements SSH Tunneling and DNS over SSH
- [ ] `socks` or [protocols/socks](/source/protocols/socks) implements SOCKS Routing (e.g. for TOR/I2P usage)


# Test Coverage

These are the `go test` files that have been implemented:

- [x] [structs/DomainCache](/source/structs/DomainCache_test.go)
- [x] [structs/WebCache](/source/structs/WebCache_test.go)
- [ ] [utils/net/url/IsTrackingParameter](/source/utils/net/url/IsTrackingParameter_test.go)
- [ ] [utils/net/url/IsXSSParameter](/source/utils/net/url/IsXSSParameter_test.go)
- [x] [utils/net/url/ResolveWebCache](/source/utils/net/url/ResolveWebCache_test.go)
- [x] [utils/net/url/ToHostAndPort](/source/utils/net/url/ToHostAndPort_test.go)
- [x] [utils/net/url/ToHost](/source/utils/net/url/ToHost_test.go)
- [ ] [utils/protocols/http/IsFilteredHeader](/source/protocols/http/IsFilteredHeader_test.go)
- [ ] [protocols/dns/Resolver](/source/protocols/dns/Resolver_test.go)
- [ ] [protocols/dns/Proxy](/source/protocols/dns/Proxy_test.go)
- [ ] [protocols/dns/Tunnel](/source/protocols/dns/Tunnel_test.go)
- [x] [protocols/dns/tunnel/ToRecordName](/source/protocols/dns/tunnel/ToRecordName_test.go)
- [ ] [protocols/http/Proxy](/source/protocols/http/Proxy_test.go)
- [ ] [protocols/http/Tunnel](/source/protocols/http/Tunnel_test.go)


# License

AGPL3
