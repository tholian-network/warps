package tunnel

import "tholian-endpoint/protocols/dns"
import net_url "net/url"
import "strconv"

func EncodeContentRange(packet *dns.Packet, url *net_url.URL, from int, to int, size int) bool {

	var result bool = false

	domain := ToRecordName(url)
	range_domain := ""

	if from != -1 && to == -1 && size == -1 {
		range_domain = "bytes." + strconv.Itoa(from) + "-.x." + domain
	} else if from != -1 && to != -1 && size == -1 {
		range_domain = "bytes." + strconv.Itoa(from) + "-" + strconv.Itoa(to) + ".x." + domain
	} else if from != -1 && to != -1 && size != -1 {
		range_domain = "bytes." + strconv.Itoa(from) + "-" + strconv.Itoa(to) + "." + strconv.Itoa(size) + "." + domain
	}

	if range_domain != "" {

		range_record := dns.NewRecord(range_domain, dns.TypeURI)
		range_record.SetURL(url.String())

		packet.AddQuestion(dns.NewQuestion(range_domain, dns.TypeURI))
		packet.AddAnswer(range_record)

		result = true

	}

	return result

}
