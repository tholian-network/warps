package certificates

import "tholian-warps/console"
import "crypto/tls"
import _ "embed"

//go:embed Proxy.crt
var embedded_proxy_cert []byte

//go:embed Proxy.key
var embedded_proxy_key []byte

func init() {

	cert, err := tls.X509KeyPair(embedded_proxy_cert, embedded_proxy_key)

	if err == nil {
		Proxy = cert
	} else {
		console.Error(err.Error())
	}

}
