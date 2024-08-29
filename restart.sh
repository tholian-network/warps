#!/bin/bash

ROOT=$(pwd);

rm ./build/tholian-warps 2> /dev/null;

bash build.sh;

"${ROOT}/build/tholian-warps" "peer";

