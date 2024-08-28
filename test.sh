#!/bin/bash

dig A tholian.network @127.0.0.1 -p 8053;
dig AAAA tholian.network @127.0.0.1 -p 8053;

curl -x localhost:8080 http://tholian.network/index.html;
curl -x localhost:8080 https://tholian.network/index.html;

