package actions

import "tholian-endpoint/types"
import "tholian-warps/console"
import "tholian-warps/protocols/dns"
import "tholian-warps/protocols/http"
// import "tholian-warps/protocols/https"
import "tholian-warps/structs"
import "strconv"

func Gateway(folder string, host string, port uint16, protocol types.Protocol) {

	console.Group("actions/Gateway")

	web_cache := structs.NewProxyCache(folder + "/proxy")

	if protocol == types.ProtocolDNS {

		proxy := dns.NewProxy(host, port, &web_cache)

		console.Log("Listening on dns://" + host + ":" + strconv.FormatUint(uint64(port), 10))

		err := proxy.Listen()

		if err != nil {
			console.Error(err.Error())
		}

	} else if protocol == types.ProtocolHTTP {

		proxy := http.NewProxy(host, port, &web_cache)

		console.Log("Listening on http://" + host + ":" + strconv.FormatUint(uint64(port), 10))

		err := proxy.Listen()

		if err != nil {
			console.Error(err.Error())
		}

	} else if protocol == types.ProtocolHTTPS {

		// TODO

		// proxy := https.NewProxy(host, port, &web_cache)

		// console.Log("Listening on https://" + host + ":" + strconv.FormatUint(uint64(port), 10))

		// err := proxy.Listen()

		// if err != nil {
		// 	console.Error(err.Error())
		// }

	}

	console.GroupEnd("actions/Gateway")

}
