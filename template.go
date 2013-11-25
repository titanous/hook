package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

type HookTemplate struct {
	Name            string                 `json:"name,omitempty"`
	Events          []string               `json:"events,omitempty"`
	SupportedEvents []string               `json:"supported_events,omitempty"`
	Schema          [][]string             `json:"schema,omitempty"`
	Config          map[string]interface{} `json:"config,omitempty"`
}

var cmdTemplates = &Command{
	Usage: "templates",
	Short: "List available hook templates",
	Long:  "List available hook templates (to be used with `hook template`)",
	Run:   runTemplates,
}

func runTemplates(cmd *Command, args []string) {
	var hooks []HookTemplate
	if err := getJSON("https://api.github.com/hooks", &hooks); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	for _, h := range hooks {
		fmt.Println(h.Name)
	}
}

func getJSON(url string, data interface{}) error {
	res, err := http.Get(url)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		return fmt.Errorf("got status code %d", res.StatusCode)
	}
	return json.NewDecoder(res.Body).Decode(data)
}

var cmdTemplate = &Command{
	Usage: "template <name>",
	Short: "Dump hook template JSON",
	Long:  "Output a named hook template (list is available from `hook templates`)",
	Run:   runTemplate,
}

func runTemplate(cmd *Command, args []string) {
	if len(args) == 0 {
		fmt.Println("no template named")
		os.Exit(1)
	}

	var hook HookTemplate
	if err := getJSON("https://api.github.com/hooks/"+args[0], &hook); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	if len(hook.SupportedEvents) == 0 || len(hook.SupportedEvents) == len(hook.Events) {
		hook.SupportedEvents = nil
	}
	hook.Config = make(map[string]interface{}, len(hook.Schema))
	for _, schema := range hook.Schema {
		switch schema[0] {
		case "string", "password":
			hook.Config[schema[1]] = ""
		case "boolean":
			hook.Config[schema[1]] = "0"
		}
	}
	hook.Schema = nil
	data, _ := json.MarshalIndent(&hook, "", "  ")
	os.Stdout.Write(data)
	os.Stdout.Write([]byte{'\n'})
}
