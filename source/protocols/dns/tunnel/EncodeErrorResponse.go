package tunnel

import "tholian-endpoint/protocols/dns"
import net_url "net/url"
import "strconv"

func EncodeErrorResponse(packet *dns.Packet, url *net_url.URL, from int, to int) {

	if packet.Type == "response" {

		domain := ToRecordName(url)
		bytes_domain := "0-.bytes." + domain
		headers_domain := "0-.headers." + domain

		if from == 0 && to == -1 {

			bytes_domain = "0-.bytes." + domain
			headers_domain = "0-.headers." + domain

		} else if from > 0 && to > from {

			bytes_domain = strconv.Itoa(from) + "-" + strconv.Itoa(to) + ".bytes." + domain
			headers_domain = strconv.Itoa(from) + "-" + strconv.Itoa(to) + ".headers." + domain

		}

		packet.AddQuestion(dns.NewQuestion(bytes_domain, dns.TypeURI))
		packet.AddQuestion(dns.NewQuestion(bytes_domain, dns.TypeTXT))
		packet.AddQuestion(dns.NewQuestion(headers_domain, dns.TypeTXT))
		packet.SetResponseCode(dns.ResponseCodeNonExistDomain)

	}

}
