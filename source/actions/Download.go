package actions

import endpoint_http "tholian-endpoint/protocols/http"
import "tholian-endpoint/types"
import "tholian-warps/certificates"
import "tholian-warps/console"
import "tholian-warps/protocols/dns"
import "tholian-warps/protocols/http"
import "tholian-warps/protocols/https"
import "tholian-warps/protocols/socks"
import "tholian-warps/utils/arguments"
import net_url "net/url"
import "os"
import "path"
import "strings"

func Download(_ string, tunnel *arguments.Config, download *net_url.URL) {

	console.Group("actions/Download")

	request := endpoint_http.NewPacket()
	request.SetURL(*download)
	request.SetMethod(endpoint_http.MethodGet)
	request.SetHeader("Range", "bytes=0-")

	var response endpoint_http.Packet

	if tunnel.Protocol == types.ProtocolDNS {

		tmp := dns.NewTunnel(tunnel.Host, tunnel.Port)
		tmp.SetDebug(true)
		response = tmp.RequestPacket(request)

	} else if tunnel.Protocol == types.ProtocolHTTP {

		tmp := http.NewTunnel(tunnel.Host, tunnel.Port)
		response = tmp.RequestPacket(request)

	} else if tunnel.Protocol == types.ProtocolHTTPS {

		tmp := https.NewTunnel(tunnel.Host, tunnel.Port)
		tmp.SetCertificate(certificates.Proxy)
		response = tmp.RequestPacket(request)

	} else if tunnel.Protocol == types.ProtocolSOCKS {

		tmp := socks.NewTunnel(tunnel.Host, tunnel.Port)
		response = tmp.RequestPacket(request)

	}

	if response.Type == "response" {

		response.Decode()

		cwd, err1 := os.Getwd()

		if err1 == nil {

			basename := path.Base(download.Path)

			if basename == "/" || basename == "." {
				basename = "index.html"
			}

			if strings.Contains(basename, ".") {

				err2 := os.WriteFile(cwd + "/" + basename, response.Payload, 0666)

				if err2 == nil {
					console.Log("File \"" + cwd + "/" + basename + "\" was written")
				} else {
					console.Error("File \"" + cwd + "/" + basename + "\" cannot be written")
				}

			}

		}

	}

	console.GroupEnd("actions/Download")

}
