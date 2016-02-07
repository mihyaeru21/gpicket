package main

import (
	"fmt"
	"os"

	"github.com/codegangsta/cli"
	"github.com/mihyaeru21/gpicket/command"
)

var GlobalFlags = []cli.Flag{}

var Commands = []cli.Command{
	{
		Name:   "log",
		Usage:  "Logging messages",
		Action: command.CmdLog,
		Flags:  []cli.Flag{},
	},
	{
		Name:   "cat",
		Usage:  "Output messages",
		Action: command.CmdCat,
		Flags:  []cli.Flag{},
	},
}

func CommandNotFound(c *cli.Context, command string) {
	fmt.Fprintf(os.Stderr, "%s: '%s' is not a %s command. See '%s --help'.", c.App.Name, command, c.App.Name, c.App.Name)
	os.Exit(2)
}
