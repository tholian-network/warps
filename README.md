
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

1. A Warps `tunnel` instance needs a `gateway` or `peer` to access the internet.
2. All Warps instances use a local [DomainCache](./source/structs/DomainCache.go) and [WebCache](./source/structs/WebCache.go).
3. All Warps instances can modify network protocols and rotate encryption keys on-demand.
4. A Warps `gateway` or `peer` uses `optimizers` which optimize the web asset file sizes.
5. A Warps `gateway` or `peer` uses `DNS over TLS` to resolve relayed DNS queries.
6. A Warps `peer` instance can tunnel through another Warps `gateway` or `peer` instance.


## How to use Tunnels and Gateways

The easiest way to use Warps is with running a Warps `gateway` on your own VPS that is connected to the internet,
and a locally running Warps `tunnel`.

As a defaulted network protocol, it is best to use `dns`, as that usually works to bypass typical firewall setups.
Alternative supported network protocols are documented in [Protocols.go](./source/types/Protocols.go).

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

# License

AGPL3
