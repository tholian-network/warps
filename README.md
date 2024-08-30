
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
2. All Warps instances use a local [DomainCache](./source/structs/DomainCache) and [WebCache](./source/structs/WebCache).
3. All Warps instances can modify network protocols and rotate encryption keys on-demand.
4. A Warps `gateway` or `peer` uses `optimizers` which optimize the web asset file sizes.
5. A Warps `gateway` or `peer` uses `DNS over TLS` to resolve relayed DNS queries.
6. A Warps `peer` instance can tunnel through another Warps `gateway` or `peer` instance.

![network-architecture.png](./assets/network-architecture.png)


# Usage

:construction: Highly Experimental at this point - Use at own risk! :construction:

```bash
bash build.sh;
cp ./build/tholian-warps /usr/bin/tholian-warps;
chmod +x /usr/bin/tholian-warps;

# Show CLI usage help
tholian-warps;
```

# License

AGPL3
