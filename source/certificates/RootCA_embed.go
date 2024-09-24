package certificates

import "tholian-warps/console"
import "crypto/tls"
import _ "embed"

//go:embed RootCA.crt
var embedded_rootca_cert []byte

//go:embed RootCA.key
var embedded_rootca_key []byte

func init() {

	cert, err := tls.X509KeyPair(embedded_rootca_cert, embedded_rootca_key)

	if err == nil {
		Proxy = cert
	} else {
		console.Error(err.Error())
	}

}
