package console

import "strings"

func Apply(message Message) {

	if message.Method == "Clear" {
		Clear()
	} else if message.Method == "Error" {
		Error(message.Value)
	} else if message.Method == "Group" {
		Group(message.Value)
	} else if message.Method == "GroupEnd" {
		GroupEnd(message.Value)
	} else if message.Method == "GroupEndResult" {

		if strings.HasSuffix(message.Value, " succeeded") {
			GroupEndResult(true, message.Value[0:len(message.Value)-10])
		} else if message.Value == "succeeded" {
			GroupEndResult(true, "")
		} else if strings.HasSuffix(message.Value, " failed") {
			GroupEndResult(false, message.Value[0:len(message.Value)-7])
		} else if message.Value == "failed" {
			GroupEndResult(false, "")
		}

	} else if message.Method == "Info" {
		Info(message.Value)
	} else if message.Method == "Inspect" {
		Inspect(message.Value)
	} else if message.Method == "Log" {
		Log(message.Value)
	} else if message.Method == "Progress" {
		Progress(message.Value)
	} else if message.Method == "Warn" {
		Warn(message.Value)
	}

}
