package main

import (
	"fmt"
	"os"

	"code.google.com/p/goauth2/oauth"
	"github.com/google/go-github/github"
)

// The CLI magic was mostly stolen from https://github.com/heroku/force and https://github.com/heroku/hk
var commands = []*Command{
	cmdHelp,
	cmdTemplates,
	cmdTemplate,
	cmdAdd,
	cmdDelete,
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

func mustClient() *github.Client {
	if os.Getenv("GITHUB_TOKEN") == "" {
		fatal("The GITHUB_TOKEN environment variable must be set")
	}
	t := &oauth.Transport{Token: &oauth.Token{AccessToken: os.Getenv("GITHUB_TOKEN")}}
	return github.NewClient(t.Client())
}

func fatal(args ...interface{}) {
	fmt.Println(args...)
	os.Exit(1)
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
