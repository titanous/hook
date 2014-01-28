package main

import (
	"encoding/json"
	"os"
	"sort"
	"strings"

	"fmt"
	"github.com/google/go-github/github"
)

var cmdDelete = &Command{
	Usage: "delete <flags> <repos>",
	Short: "delete a hook",
	Long:  "delete a hook",
}

func init() {
	// break init loop
	cmdDelete.Run = runDelete
}

var deleteHookFile = cmdDelete.Flag.String("file", "", "hook json config")

func runDelete(cmd *Command, args []string) {
	hook := parseHook(*deleteHookFile)
	client := mustClient()

	for _, repo := range getRepos(client, args) {
		ownerRepo := strings.SplitN(repo, "/", 2)
		hooks, _, err := client.Repositories.ListHooks(ownerRepo[0], ownerRepo[1], &github.ListOptions{PerPage: 100})
		if err != nil {
			fmt.Println(repo, "lookup failed:", err)
			continue
		}
		for _, h := range hooks {
			if *hook.Name == "web" && *h.Name == "web" && h.Config["url"] == hook.Config["url"] || *h.Name == *hook.Name {
				_, err = client.Repositories.DeleteHook(ownerRepo[0], ownerRepo[1], *h.ID)
				if err != nil {
					fmt.Println(repo, "remove failed:", err)
					continue
				}
				fmt.Println(repo, "deleted")
			}
		}
	}
}

func getRepos(client *github.Client, list []string) []string {
	var repos []string
	for _, repoName := range list {
		if strings.HasSuffix(repoName, "/") {
			res, err := listRepos(client, repoName[:len(repoName)-1])
			if err != nil {
				fatal("Error listing repos:", err)
			}
			repos = append(repos, res...)
		} else {
			repos = append(repos, repoName)
		}
	}
	return repos
}

func parseHook(filename string) *github.Hook {
	f, err := os.Open(filename)
	if err != nil {
		fatal("Error opening hook json:", err)
	}
	var hook github.Hook
	if err := json.NewDecoder(f).Decode(&hook); err != nil {
		fatal("Error reading hook json:", err)
	}
	sort.Strings(hook.Events)
	return &hook
}
