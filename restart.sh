#!/bin/bash

ROOT=$(pwd);

bash build.sh;

"${ROOT}/build/tholian-warps" "peer";

