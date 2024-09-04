package structs

import "tholian-warps/console"
import "crypto/tls"
import _ "embed"

var Certificate *tls.Certificate

//go:embed Certificate.crt
var embedded_cert_pem []byte

//go:embed Certificate.key
var embedded_key_pem []byte

func init() {

	cert, err := tls.X509KeyPair(embedded_cert_pem, embedded_key_pem)

	if err == nil {

		Certificate = &cert

	} else {
		console.Error(err.Error())
	}

}
