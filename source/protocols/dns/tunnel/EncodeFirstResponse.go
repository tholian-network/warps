package tunnel

import "tholian-endpoint/protocols/dns"
import "encoding/json"
import net_url "net/url"
import "strconv"

func EncodeFirstResponse(packet *dns.Packet, url *net_url.URL, headers map[string]string, payload []byte) {

	if packet.Type == "response" {

		domain := ToRecordName(url)
		bytes_domain := "0-.bytes." + domain
		headers_domain := "0-.headers." + domain

		packet.AddQuestion(dns.NewQuestion(bytes_domain, dns.TypeURI))
		packet.AddQuestion(dns.NewQuestion(bytes_domain, dns.TypeTXT))
		packet.AddQuestion(dns.NewQuestion(headers_domain, dns.TypeTXT))

		url_record := dns.NewRecord(bytes_domain, dns.TypeURI)
		url_record.SetURL(url.String())
		packet.AddAnswer(url_record)

		if len(payload) > 512 {

			headers["Content-Range"] = "bytes 0-511/" + strconv.Itoa(len(payload))
			// headers["Content-Length"] = strconv.Itoa(len(payload))

			bytes_record := dns.NewRecord(bytes_domain, dns.TypeTXT)
			bytes_record.SetData(payload[0:512])
			packet.AddAnswer(bytes_record)

			headers_record := dns.NewRecord(headers_domain, dns.TypeTXT)
			buffer, _ := json.Marshal(headers)
			headers_record.SetData(buffer)
			packet.AddAnswer(headers_record)

		} else {

			headers["Content-Range"] = "bytes 0-" + strconv.Itoa(len(payload) - 1) + "/" + strconv.Itoa(len(payload))
			// headers["Content-Length"] = strconv.Itoa(len(payload))

			bytes_record := dns.NewRecord(bytes_domain, dns.TypeTXT)
			bytes_record.SetData(payload)

			headers_record := dns.NewRecord(headers_domain, dns.TypeTXT)
			buffer, _ := json.Marshal(headers)
			headers_record.SetData(buffer)
			packet.AddAnswer(headers_record)

		}

	}

}
