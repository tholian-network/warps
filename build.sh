#!/bin/bash

ROOT=$(pwd);

build_certificate() {

	cd "${ROOT}/source";

	openssl genrsa -out "${ROOT}/source/structs/Certificate.key" 2048;
	openssl req -new -key "${ROOT}/source/structs/Certificate.key" -x509 -sha256 -subj "/C=DE/ST=Berlin/L=Berlin/O=Tholian Network/OU=Warps Proxy" -out "${ROOT}/source/structs/Certificate.crt";

}

build_warps() {

	cd "${ROOT}/source";

	go build -o "${ROOT}/build/tholian-warps" "${ROOT}/source/cmds/tholian-warps/main.go";
	chmod +x "${ROOT}/build/tholian-warps";

}


build_certificate;
build_warps;

