package main

import (
	"errors"
	"fmt"

	"github.com/cirss/geist/cli"
)

func handleHelpSubcommand(cc *cli.CommandContext) (err error) {
	if len(cc.Args) < 2 {
		fmt.Fprintln(Main.OutWriter)
		cc.ShowProgramUsage()
		return
	}
	commandName := cc.Args[1]
	if commandName == "help" {
		return
	}
	if c, exists := commandCollection.Lookup(commandName); exists {
		cc.Descriptor = c
		cc.Args = []string{commandName, "help"}
		c.Handler(cc)
	} else {
		fmt.Fprintf(Main.ErrWriter, "\nnot a blazegraph command: %s\n\n", commandName)
		cc.ShowProgramUsage()
		err = errors.New("Not a blazegraph command")
	}
	return
}
