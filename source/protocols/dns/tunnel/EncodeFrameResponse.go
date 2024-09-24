package tunnel

import "tholian-endpoint/protocols/dns"
import "encoding/json"
import net_url "net/url"
import "strconv"

func EncodeFrameResponse(packet *dns.Packet, url *net_url.URL, headers map[string]string, payload []byte, from int, to int, size int) {

	if packet.Type == "response" {

		domain := ToRecordName(url)
		bytes_domain := strconv.Itoa(from) + "-" + strconv.Itoa(to) + ".bytes." + domain
		headers_domain := strconv.Itoa(from) + "-" + strconv.Itoa(to) + ".headers." + domain

		packet.AddQuestion(dns.NewQuestion(bytes_domain, dns.TypeURI))
		packet.AddQuestion(dns.NewQuestion(bytes_domain, dns.TypeTXT))
		packet.AddQuestion(dns.NewQuestion(headers_domain, dns.TypeTXT))

		url_record := dns.NewRecord(bytes_domain, dns.TypeURI)
		url_record.SetURL(url.String())
		packet.AddAdditional(url_record)

		headers["Content-Range"] = "bytes " + strconv.Itoa(from) + "-" + strconv.Itoa(to) + "/" + strconv.Itoa(size)

		_, ok := headers["Content-Length"]

		if ok == true {
			delete(headers, "Content-Length")
		}

		bytes_record := dns.NewRecord(bytes_domain, dns.TypeTXT)
		bytes_record.SetData(payload)
		packet.AddAnswer(bytes_record)

		headers_record := dns.NewRecord(headers_domain, dns.TypeTXT)
		buffer, _ := json.Marshal(headers)
		headers_record.SetData(buffer)
		packet.AddAnswer(headers_record)

	}

}
