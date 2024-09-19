package main

import "tholian-endpoint/types"
import "tholian-warps/actions"
import "tholian-warps/console"
import "os"
import "strconv"
import "strings"

func showUsage() {

	console.Info("")
	console.Info("Tholian Warps")
	console.Info("Adaptive Mesh Network Router")
	console.Info("")

	console.Group("Usage: tholian-warps [Action]")
	console.GroupEnd("------")

	console.Group("Action  | Description                                                                               |")
	console.Log("--------|-------------------------------------------------------------------------------------------|")
	console.Log("tunnel  | Starts local HTTPS/HTTP/DNS proxy and tunnels traffic to a target Tholian Warps instance. |")
	console.Log("gateway | Starts local tunnel listener and executes HTTPS/HTTP/DNS requests.                        |")
	console.Log("peer    | Serves local HTTPS/HTTP/DNS Proxy and executes HTTPS/HTTP/DNS requests.                   |")
	console.GroupEnd("--------|-------------------------------------------------------------------------------------------|")

	console.Group("Examples")
	console.Log("# Start a gateway on your unblocked and internet-connected machine")
	console.Log("sudo tholian-warps gateway --protocol=dns --port=1053;")
	console.Log("")
	console.Log("# Start a local proxy and tunnel its traffic through another Warps gateway")
	console.Log("sudo tholian-warps tunnel --protocol=dns --host=1.3.3.7 --port=5353;")
	console.GroupEnd("--------")

}

func main() {

	var action string
	var folder string = "/tmp/tholian-warps"
	var host string
	var port uint16
	var protocol types.Protocol

	xdg_cache_home := os.Getenv("XDG_CACHE_HOME")
	home := os.Getenv("HOME")

	if xdg_cache_home != "" {
		folder = xdg_cache_home + "/tholian-warps"
	} else if home != "" {
		folder = home + "/tholian-warps"
	}

	if len(os.Args) > 2 {

		action = strings.ToLower(os.Args[1])

		for a := 2; a < len(os.Args); a++ {

			arg := os.Args[a]

			if strings.HasPrefix(arg, "--host=") {

				tmp := strings.TrimSpace(arg[7:])

				if strings.HasPrefix(tmp, "\"") && strings.HasSuffix(tmp, "\"") {
					host = tmp[1:len(tmp)-1]
				} else if strings.HasPrefix(tmp, "'") && strings.HasSuffix(tmp, "'") {
					host = tmp[1:len(tmp)-1]
				} else {
					host = tmp
				}

			} else if strings.HasPrefix(arg, "--port=") {

				tmp := strings.TrimSpace(arg[7:])

				if strings.HasPrefix(tmp, "\"") && strings.HasSuffix(tmp, "\"") {
					tmp = tmp[1:len(tmp)-1]
				} else if strings.HasPrefix(tmp, "'") && strings.HasSuffix(tmp, "'") {
					tmp = tmp[1:len(tmp)-1]
				}

				num, err := strconv.ParseUint(tmp, 10, 16)

				if err == nil {
					port = uint16(num)
				}

			} else if strings.HasPrefix(arg, "--protocol=") {

				tmp := strings.TrimSpace(arg[11:])

				if strings.HasPrefix(tmp, "\"") && strings.HasSuffix(tmp, "\"") {
					tmp = tmp[1:len(tmp)-1]
				} else if strings.HasPrefix(tmp, "'") && strings.HasSuffix(tmp, "'") {
					tmp = tmp[1:len(tmp)-1]
				}

				if tmp == "dns" {
					protocol = types.ProtocolDNS
				} else if tmp == "http" {
					protocol = types.ProtocolHTTP
				} else if tmp == "https" {
					protocol = types.ProtocolHTTPS
				}

			}

		}

	} else if len(os.Args) == 2 {
		action = strings.ToLower(os.Args[1])
	}

	console.Clear()

	if host == "" {
		host = "localhost"
	}

	console.Group("tholian-warps: Command-Line Arguments")
	console.Inspect(struct{
		Action   string
		Folder   string
		Host     string
		Port     uint16
		Protocol string
	}{
		Action:   action,
		Folder:   folder,
		Host:     host,
		Port:     port,
		Protocol: protocol.String(),
	})
	console.GroupEnd("")

	if action == "gateway" {

		actions.Gateway(folder, host, port, protocol)

	} else if action == "tunnel" {

		actions.Tunnel(folder, host, port, protocol)

	} else if action == "peer" {

		actions.Peer(folder, host)

	} else {

		console.Clear()
		showUsage()

		os.Exit(1)

	}

}

