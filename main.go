package main

import (
	"os"
)

// The CLI magic was mostly stolen from https://github.com/heroku/force and https://github.com/heroku/hk
var commands = []*Command{
	cmdHelp,
	cmdTemplates,
	cmdTemplate,
}

func main() {
	args := os.Args[1:]
	if len(args) < 1 {
		usage()
	}

	for _, cmd := range commands {
		if cmd.Name() == args[0] && cmd.Run != nil {
			cmd.Flag.Usage = func() {
				cmd.printUsage()
			}
			if err := cmd.Flag.Parse(args[1:]); err != nil {
				os.Exit(2)
			}
			cmd.Run(cmd, cmd.Flag.Args())
			return
		}
	}
	usage()
}

// operate on:
// list of repos
// prefix (user or organization)
//
// add hook
//  - url or json hook config
//  - force replace
// remove hook
//  - name or url
// deactivate hook
// activate hook
// show hooks
//  - dump json hook config
