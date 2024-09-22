package main

import "tholian-endpoint/types"
import "tholian-warps/actions"
import "tholian-warps/console"
import "tholian-warps/utils/arguments"
import "os"
import "strings"

func showUsage() {

	console.Info("")
	console.Info("Tholian Warps")
	console.Info("Adaptive Exfil Network Protocol Router")
	console.Info("")

	console.Group("Usage: tholian-warps [Action] [Listen URL]")
	console.Log("       tholian-warps [Action] [Listen URL] [Tunnel URL]")
	console.GroupEnd("------")

	console.Log("")
	console.Group("Listen and Tunnel URLs                                                    |")
	console.Log("The \"protocol://host:port\" scheme supports the protocols \"dns\", \"http\", \"https\". |")
	console.Log("The \"0.0.0.0\" host listens to incoming traffic on a given protocol and port.     |")
	console.Log("The \"any\" value listens to incoming traffic on the defaulted protocols, which    |")
	console.Log("are \"dns://0.0.0.0:1053\", \"http://0.0.0.0:1080\", and \"https://0.0.0.0:1443\".     |")
	console.GroupEnd("---------------------------------------------------------------------------------|")

	console.Log("")
	console.Group("Action  | Description                                                                    |")
	console.Log("--------|--------------------------------------------------------------------------------|")
	console.Log("tunnel  | Starts local Proxy and tunnels network traffic to a specified Warps instance.  |")
	console.Log("forward | Forwards network traffic from a Warps instance to another Warps instance.      |")
	console.Log("gateway | Starts local Tunnel Proxy and executes HTTPS/HTTP/DNS traffic to the internet. |")
	console.GroupEnd("--------|--------------------------------------------------------------------------------|")

	console.Log("")
	console.Group("Tunnel Example")
	console.Log("# VPS (IP 1.3.3.7): Start gateway to the internet")
	console.Log("tholian-warps gateway \"dns://0.0.0.0:1337\";")
	console.Log("")
	console.Log("# Desktop System: Start a local tunnel to 1.3.3.7")
	console.Log("tholian-warps tunnel \"any\" \"dns://1.3.3.7:1337\";")
	console.GroupEnd("--------------")

	console.Log("")
	console.Group("Forward Example")
	console.Log("# First VPS (IP 1.3.3.7): Start gateway to the internet")
	console.Log("tholian-warps gateway \"http://0.0.0.0:1337\";")
	console.Log("")
	console.Log("# Second VPS (IP 1.3.3.8): Forward incoming traffic to 1.3.3.7")
	console.Log("tholian-warps forward \"dns://1.3.3.8:1338\" \"http://1.3.3.7:1337\"")
	console.Log("")
	console.Log("# Desktop System: Start a local tunnel to 1.3.3.8")
	console.Log("tholian-warps tunnel \"any\" \"dns://1.3.3.8:1338\";")
	console.GroupEnd("---------------")

}

func main() {

	var action string
	var folder string = "/tmp/tholian-warps"
	var listen string = "any"
	var tunnel string = ""

	xdg_cache_home := os.Getenv("XDG_CACHE_HOME")
	home := os.Getenv("HOME")

	if xdg_cache_home != "" {
		folder = xdg_cache_home + "/tholian-warps"
	} else if home != "" {
		folder = home + "/tholian-warps"
	}

	if len(os.Args) > 2 {

		action = strings.ToLower(os.Args[1])

		if len(os.Args) >= 3 && os.Args[2] != "" {

			tmp := strings.ToLower(os.Args[2])

			if strings.HasPrefix(tmp, "\"") && strings.HasSuffix(tmp, "\"") {
				tmp = tmp[1:len(tmp)-1]
			}

			if tmp != "" {
				listen = tmp
			}

		}

		if len(os.Args) == 4 && os.Args[3] != "" {

			tmp := strings.ToLower(os.Args[3])

			if strings.HasPrefix(tmp, "\"") && strings.HasSuffix(tmp, "\"") {
				tmp = tmp[1:len(tmp)-1]
			}

			if tmp != "" {
				tunnel = tmp
			}

		}

	}

	console.Clear()

	console.Group("tholian-warps: Command-Line Arguments")
	console.Inspect(struct{
		Action string
		Folder string
		Listen string
		Tunnel string
	}{
		Action: action,
		Folder: folder,
		Listen: listen,
		Tunnel: tunnel,
	})
	console.GroupEnd("")

	if action == "gateway" {

		listen := arguments.ParseConfig(listen)

		if listen != nil {

			actions.Gateway(folder, listen)

		} else {

			console.Error("Cannot parse Configs from provided Listen URL")
			
			showUsage()
			os.Exit(1)

		}

	} else if action == "forward" {

		listen := arguments.ParseConfig(listen)
		tunnel := arguments.ParseConfig(tunnel)

		if listen != nil && tunnel != nil && tunnel.Protocol != types.ProtocolANY {

			actions.Forward(folder, listen, tunnel)

		} else if listen != nil {

			console.Error("Cannot parse Configs from provided Tunnel URL")
			console.Error("Did you intend to use the \"gateway\" action instead?")

			showUsage()
			os.Exit(1)

		} else {

			console.Error("Cannot parse Configs from provided Listen URL")
			
			showUsage()
			os.Exit(1)

		}

	} else if action == "tunnel" {

		listen := arguments.ParseConfig(listen)
		tunnel := arguments.ParseConfig(tunnel)

		if listen != nil && tunnel != nil && tunnel.Protocol != types.ProtocolANY {

			actions.Tunnel(folder, listen, tunnel)

		} else if listen != nil {

			console.Error("Cannot parse Configs from provided Tunnel URL")
			console.Error("Did you intend to use the \"gateway\" action instead?")

			showUsage()
			os.Exit(1)

		} else {

			console.Error("Cannot parse Configs from provided Listen URL")
			
			showUsage()
			os.Exit(1)

		}

	} else {

		console.Clear()
		showUsage()

		os.Exit(1)

	}

}

