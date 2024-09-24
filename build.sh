#!/bin/bash

GO="$(which go)";
OPENSSL="$(which openssl)";
ROOT=$(pwd);
FLAG="${1}";

renew_certificate() {

	local rootca_key="${ROOT}/source/certificates/RootCA.key";
	local rootca_crt="${ROOT}/source/certificates/RootCA.crt";
	local rootca_pem="${ROOT}/source/certificates/RootCA.pem";
	local rootca_subject="/C=DE/ST=Berlin/L=Berlin/O=Tholian Network/OU=Warps Proxy";

	local server_key="${ROOT}/source/certificates/Proxy.key";
	local server_crt="${ROOT}/source/certificates/Proxy.crt";
	local server_csr="${ROOT}/source/certificates/Proxy.csr";
	local server_pem="${ROOT}/source/certificates/Proxy.pem";
	local server_subject="/C=DE/ST=Berlin/L=Berlin/O=Decentralized/OU=Warps Proxy";

	# Generate root certificate authority's private key and certificate
	openssl req -x509 -sha256 -nodes -newkey rsa:4096 -keyout "${rootca_key}" -out "${rootca_crt}" -days 3560 -subj "${rootca_subject}";
	openssl x509 -in "${rootca_crt}" -out "${rootca_pem}";

	if [[ "$?" == "0" ]]; then
		echo -e "- Renew Root CA certificate [\e[32mok\e[0m]";
	else
		echo -e "- Renew Root CA certificate [\e[31mfail\e[0m]";
	fi;

	# Generate server's private key and certificate
	openssl genrsa -out "${server_key}" 4096;
	openssl req -nodes -key "${server_key}" -new -out "${server_csr}" -subj "${server_subject}";

	if [[ "$?" == "0" ]]; then
		echo -e "- Renew Proxy certificate [\e[32mok\e[0m]";
	else
		echo -e "- Renew Proxy certificate [\e[31mfail\e[0m]";
	fi;

	# Sign the server's private key and certificate
	openssl x509 -req -CA "${rootca_crt}" -CAkey "${rootca_key}" -in "${server_csr}" -out "${server_crt}" -days 365 -CAcreateserial -extfile "build.ext";
	openssl x509 -in "${server_crt}" -out "${server_pem}";

	if [[ "$?" == "0" ]]; then
		echo -e "- Sign Proxy certificate [\e[32mok\e[0m]";
	else
		echo -e "- Sign Proxy certificate [\e[31mfail\e[0m]";
	fi;

}

build_warps() {

	local os="$1";
	local arch="$2";
	local binary="${ROOT}/build/tholian-warps-${os}_${arch}";

	if [[ "${os}" == "windows" ]]; then
		binary="${ROOT}/build/tholian-warps-${os}_${arch}.exe";
	fi;

	cd "${ROOT}/source";

	if [[ -f "${binary}" ]]; then
		rm "${binary}" > /dev/null;
	fi;

	env CGO_ENABLED=0 GOOS="${os}" GOARCH="${arch}" go build -ldflags="-s -w" -o "${binary}" "${ROOT}/source/cmds/tholian-warps/main.go";

	if [[ "$?" == "0" ]]; then
		echo -e "- Build tholian-warps (${os} / ${arch}) [\e[32mok\e[0m]";
	else
		echo -e "- Build tholian-warps (${os} / ${arch}) [\e[31mfail\e[0m]";
	fi;

	chmod +x "${binary}";
	strip "${binary}";

	if [[ "$?" == "0" ]]; then
		echo -e "- Strip tholian-warps (${os} / ${arch}) [\e[32mok\e[0m]";
	else
		echo -e "- Strip tholian-warps (${os} / ${arch}) [\e[31mfail\e[0m]";
	fi;

}

if [[ "${FLAG}" == "--renew-certificate" ]]; then

	if [[ "${OPENSSL}" != "" ]]; then

		renew_certificate;

	else
		echo "Please install OpenSSL to renew the Tholian Warps RootCA certificate.";
		exit 1;
	fi;

fi;

if [[ "${GO}" != "" ]]; then

	build_warps "linux" "amd64";
	build_warps "linux" "arm64";

else
	echo "Please install go(lang) compiler.";
	exit 1;
fi;

