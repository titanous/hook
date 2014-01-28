package main

import (
	"fmt"
	"reflect"
	"sort"
	"strings"

	"github.com/google/go-github/github"
)

var cmdAdd = &Command{
	Usage: "add <flags> <repos>",
	Short: "Add a hook",
	Long:  "Add a hook",
}

func init() {
	// break init loop
	cmdAdd.Run = runAdd
}

var addHookFile = cmdAdd.Flag.String("file", "", "hook json config")

func runAdd(cmd *Command, args []string) {
	hook := parseHook(*addHookFile)
	client := mustClient()

outer:
	for _, repo := range getRepos(client, args) {
		ownerRepo := strings.SplitN(repo, "/", 2)
		hooks, _, err := client.Repositories.ListHooks(ownerRepo[0], ownerRepo[1], &github.ListOptions{PerPage: 100})
		if err != nil {
			fmt.Println(repo, "lookup failed:", err)
			continue
		}
		var replace int
		for _, h := range hooks {
			if *hook.Name == "web" && *h.Name == "web" && h.Config["url"] == hook.Config["url"] || *h.Name == *hook.Name {
				if len(h.Events) != len(hook.Events) || len(hook.Config) != len(h.Config) {
					replace = *h.ID
					break
				} else {
					sort.Strings(h.Events)
					if !reflect.DeepEqual(h.Events, hook.Events) || !reflect.DeepEqual(h.Config, hook.Config) {
						replace = *h.ID
						break
					}
				}
				fmt.Println(repo, "exists")
				continue outer
			}
		}
		if replace > 0 {
			_, err = client.Repositories.DeleteHook(ownerRepo[0], ownerRepo[1], replace)
			if err != nil {
				fmt.Println(repo, "remove failed:", err)
				continue
			}
		}
		_, _, err = client.Repositories.CreateHook(ownerRepo[0], ownerRepo[1], hook)
		if err != nil {
			fmt.Println(repo, "failed:", err)
			continue
		}
		if replace > 0 {
			fmt.Println(repo, "replaced")
		} else {
			fmt.Println(repo, "added")
		}
	}
}

func listRepos(client *github.Client, name string) ([]string, error) {
	var names []string
	nextPage := 1
	for {
		repos, res, err := client.Repositories.ListByOrg(name, &github.RepositoryListByOrgOptions{ListOptions: github.ListOptions{Page: nextPage, PerPage: 100}})
		if err != nil {
			break
		}
		for _, repo := range repos {
			names = append(names, *repo.Owner.Login+"/"+*repo.Name)
		}
		if res.NextPage == 0 {
			return names, nil
		}
		nextPage = res.NextPage
	}
	for {
		repos, res, err := client.Repositories.List(name, &github.RepositoryListOptions{ListOptions: github.ListOptions{Page: nextPage, PerPage: 100}})
		if err != nil {
			return nil, err
		}
		for _, repo := range repos {
			names = append(names, *repo.Owner.Login+"/"+*repo.Name)
		}
		if res.NextPage == 0 {
			return names, nil
		}
		nextPage = res.NextPage
	}
}
