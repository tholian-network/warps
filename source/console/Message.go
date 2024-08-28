package console

import "runtime"
import "strings"

type Message struct {
	Method string `json:"type"`
	Value  string `json:"value"`
	Caller struct {
		File string `json:"file"`
		Line int    `json:"line"`
	} `json:"caller"`
}

func NewMessage(method string, value string) Message {

	var message Message

	message.Method = method
	message.Value = value

	// skip NewMessage()
	// skip console.<Method>()
	_, file, line, ok := runtime.Caller(2)

	if ok == true {

		// XXX: Currently there's no way to get the source code / symbols file path
		if strings.Contains(file, "/Software/tholian-network/endpoint/") {
			file = file[strings.Index(file, "/Software/tholian-network/endpoint/")+35:]
		}

		message.Caller.File = file
		message.Caller.Line = line

	} else {

		message.Caller.File = "???"
		message.Caller.Line = 0

	}

	return message

}
