package actions

import "tholian-warps/structs"
import "tholian-warps/types"
import "os"

func Gateway(host string, port uint16, protocol types.Protocol) {

	pwd, err := os.Getwd()

	if err == nil {

		web_cache := structs.NewWebCache(pwd + "/warps/gateway")

		proxy := structs.NewProxy(host, port, &web_cache, nil, protocol)
		proxy.Listen()

	}

}
