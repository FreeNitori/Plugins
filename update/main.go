package main

import (
	"bytes"
	"fmt"
	"git.randomchars.net/FreeNitori/FreeNitori/nitori/multiplexer"
	"git.randomchars.net/FreeNitori/FreeNitori/nitori/state"
	"os/exec"
)

// Setup sets up the plugin and returns route.
//goland:noinspection GoUnusedExportedFunction
func Setup() interface{} {
	return &multiplexer.Route{
		Pattern:       "update",
		AliasPatterns: []string{""},
		Description:   "",
		Category:      multiplexer.SystemCategory,
		Handler: func(context *multiplexer.Context) {
			if !context.IsAdministrator() {
				context.SendMessage(state.AdminOnly)
				return
			}

			context.SendMessage("Compiling...")
			command := exec.Command("/bin/sh", "-c", "make && make plugins")
			var output bytes.Buffer
			command.Stdout = &output
			command.Stderr = &output
			err := command.Run()
			if err != nil {
				context.SendMessage("Error occurred while compiling.")
				context.SendMessage(fmt.Sprintf("```\n%s\n```", output.String()))
				return
			}
			context.SendMessage("Finished compiling, attempting restart...")
			state.ExitCode <- -1
		},
	}
}
