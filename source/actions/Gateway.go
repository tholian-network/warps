package actions

import "tholian-warps/console"
import "tholian-warps/structs"
import "tholian-warps/types"

func Gateway(folder string, host string, port uint16, protocol types.Protocol) {

	console.Group("actions/Gateway")

	web_cache := structs.NewWebCache(folder + "/proxy")

	proxy := structs.NewProxy(host, port, &web_cache, nil, protocol)

	err := proxy.Listen()

	if err != nil {
		console.Error(err.Error())
	}

	console.GroupEnd("actions/Gateway")

}
