package https

import "tholian-warps/protocols/http"
import "crypto/rand"
import "crypto/rsa"
import "crypto/tls"
import "crypto/x509"
import "crypto/x509/pkix"
import "math/big"
import "net"
import net_http "net/http"
import net_url "net/url"
import "strconv"
import "testing"
import "time"

func generateCertificate(t *testing.T) tls.Certificate {

	key, err := rsa.GenerateKey(rand.Reader, 2048)

	if err != nil {
		t.Fatalf("Cannot generate key: %s", err.Error())
	}

	template := x509.Certificate{
		SerialNumber:          big.NewInt(1),
		Subject:               pkix.Name{CommonName: "localhost"},
		NotBefore:             time.Now().Add(-time.Hour),
		NotAfter:              time.Now().Add(time.Hour),
		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
		IsCA:                  true,
		IPAddresses:           []net.IP{net.ParseIP("127.0.0.1")},
		DNSNames:              []string{"localhost"},
	}

	der, err := x509.CreateCertificate(rand.Reader, &template, &template, &key.PublicKey, key)

	if err != nil {
		t.Fatalf("Cannot create certificate: %s", err.Error())
	}

	leaf, err := x509.ParseCertificate(der)

	if err != nil {
		t.Fatalf("Cannot parse certificate: %s", err.Error())
	}

	return tls.Certificate{
		Certificate: [][]byte{der},
		PrivateKey:  key,
		Leaf:        leaf,
	}

}

func TestTunnel(t *testing.T) {

	t.Run("Tunnel with HTTPS Payload", func(t *testing.T) {

		certificate := generateCertificate(t)

		listener, err := net.ListenTCP("tcp", &net.TCPAddr{
			IP:   net.ParseIP("127.0.0.1"),
			Port: 0,
		})

		if err != nil {
			t.Fatalf("Cannot listen: %s", err.Error())
		}

		port := listener.Addr().(*net.TCPAddr).Port

		tls_listener := tls.NewListener(listener, &tls.Config{
			Certificates: []tls.Certificate{certificate},
		})

		server := &net_http.Server{
			Handler: net_http.HandlerFunc(func(response net_http.ResponseWriter, request *net_http.Request) {
				response.WriteHeader(200)
				response.Write([]byte("Hello, world!"))
			}),
		}

		go server.Serve(tls_listener)
		defer server.Close()

		tunnel := NewTunnel("127.0.0.1", uint16(port))
		tunnel.SetCertificate(certificate)

		url, _ := net_url.Parse("https://127.0.0.1:" + strconv.Itoa(port) + "/index.html")

		request := http.NewPacket()
		request.SetMethod(http.MethodGet)
		request.SetURL(*url)
		request.SetHeader("Host", url.Host)

		response := tunnel.RequestPacket(request)
		response.Decode()

		if response.Status != http.StatusOK {
			t.Errorf("Expected HTTP response status '%s' but got '%s'", http.Status(http.StatusOK).String(), http.Status(response.Status).String())
		}

		if string(response.Payload) != "Hello, world!" {
			t.Errorf("Expected HTTP response payload '%s' but got '%s'", "Hello, world!", string(response.Payload))
		}

	})

}
