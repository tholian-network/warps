package tunnel

import "tholian-endpoint/protocols/dns"
import net_url "net/url"

func EncodeFirstRequest(packet *dns.Packet, url *net_url.URL) {

	if packet.Type == "query" {

		domain := ToRecordName(url)
		bytes_domain := "0-.bytes." + domain
		headers_domain := "0-.headers." + domain

		packet.AddQuestion(dns.NewQuestion(bytes_domain, dns.TypeURI))
		packet.AddQuestion(dns.NewQuestion(bytes_domain, dns.TypeTXT))
		packet.AddQuestion(dns.NewQuestion(headers_domain, dns.TypeTXT))

		url_record := dns.NewRecord(bytes_domain, dns.TypeURI)
		url_record.SetURL(url.String())
		packet.AddAdditional(url_record)

	}

}
