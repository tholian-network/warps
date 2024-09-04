#!/bin/bash

ROOT=$(pwd);
CACERT="${ROOT}/source/structs/Certificate.crt";

# DNS Tests
dig A cookie.engineer @127.0.0.1 -p 8053;
dig A tholian.network @127.0.0.1 -p 8053;
dig AAAA tholian.network @127.0.0.1 -p 8053;

# HTTP Tests
curl -x localhost:8080 -L --insecure --cacert "${CACERT}" http://tholian.network/;
curl -x localhost:8080 -L --insecure --cacert "${CACERT}" https://tholian.network/;

