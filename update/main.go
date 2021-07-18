package main

import (
	"bytes"
	"fmt"
	"git.randomchars.net/freenitori/freenitori/v2/nitori/state"
	"git.randomchars.net/freenitori/multiplexer"
	"os/exec"
)

// Setup sets up the plugin and returns route.
//goland:noinspection GoUnusedExportedFunction
func Setup() interface{} {
	return &multiplexer.Route{
		Pattern:       "update",
		AliasPatterns: []string{},
		Description:   "",
		Category:      multiplexer.SystemCategory,
		Handler: func(context *multiplexer.Context) {
			if !context.IsAdministrator() {
				context.SendMessage(multiplexer.AdminOnly)
				return
			}

			context.SendMessage("Compiling...")
			command := exec.Command("/bin/sh", "-c", "make && cd plugins && make")
			var output bytes.Buffer
			command.Stdout = &output
			command.Stderr = &output
			err := command.Run()
			if err != nil {
				context.SendMessage("Error occurred while compiling.")
				context.SendMessage(fmt.Sprintf("```\n%s\n```", output.String()))
				return
			}
			message := context.SendMessage("Finished compiling, attempting restart...")
			if message != nil {
				state.Reincarnation = message.ChannelID + "\t" + message.ID + "\t" + "Update complete."
			}
			state.Exit <- -1
		},
	}
}
