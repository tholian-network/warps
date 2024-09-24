package tunnel

import "tholian-endpoint/protocols/dns"
import net_url "net/url"
import "strconv"

func EncodeFrameRequest(packet *dns.Packet, url *net_url.URL, from int, to int) {

	if packet.Type == "query" {

		domain := ToRecordName(url)
		bytes_domain := strconv.Itoa(from) + "-" + strconv.Itoa(to) + ".bytes." + domain
		headers_domain := strconv.Itoa(from) + "-" + strconv.Itoa(to) + ".headers." + domain

		packet.AddQuestion(dns.NewQuestion(bytes_domain, dns.TypeURI))
		packet.AddQuestion(dns.NewQuestion(bytes_domain, dns.TypeTXT))
		packet.AddQuestion(dns.NewQuestion(headers_domain, dns.TypeTXT))

		url_record := dns.NewRecord(bytes_domain, dns.TypeURI)
		url_record.SetURL(url.String())
		packet.AddAnswer(url_record)

	}

}
