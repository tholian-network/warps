package actions

import "tholian-warps/structs"
import "tholian-warps/types"

func Gateway(folder string, host string, port uint16, protocol types.Protocol) {

	web_cache := structs.NewWebCache(folder + "/proxy")

	proxy := structs.NewProxy(host, port, &web_cache, nil, protocol)
	proxy.Listen()

}
