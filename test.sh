#!/bin/bash

ROOT=$(pwd);
CACERT="${ROOT}/source/config/Certificate.crt";

# DNS Tests
dig A cookie.engineer @127.0.0.1 -p 1053;
dig A tholian.network @127.0.0.1 -p 1053;
dig AAAA tholian.network @127.0.0.1 -p 1053;

# HTTP Tests
curl -x localhost:1080 -L http://tholian.network/;
curl -x localhost:1443 -L --insecure --cacert "${CACERT}" https://tholian.network/;

