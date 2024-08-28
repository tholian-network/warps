#!/bin/bash

ROOT=$(pwd);


cd "${ROOT}/source";

go build -o "${ROOT}/build/tholian-warps" "${ROOT}/source/cmds/tholian-warps/main.go";

