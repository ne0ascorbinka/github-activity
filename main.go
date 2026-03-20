package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	event "github.com/ne0ascorbinka/github-activity/internal"
)

const url = "https://api.github.com/users/%s/events"

func main() {
	flag.Usage = func() {
		fmt.Fprint(os.Stderr, "Usage: github-activity <username>\n\n")
		fmt.Fprint(os.Stderr, "A CLI tool to list recent GitHub activity for a specific user.\n")
		flag.PrintDefaults()
	}

	flag.Parse()
	args := flag.Args()

	if len(args) != 1 {
		flag.Usage()
		os.Exit(1)
	}
	
	username := args[0]

	r, err := http.Get(fmt.Sprintf(url, username))
	if err != nil {
		log.Fatalf("Error fetching GitHub data: %v", err)
	}
	defer r.Body.Close()

	if r.StatusCode == http.StatusNotFound {
		fmt.Printf("Error: User '%s' not found.\n", username)
		os.Exit(1)
	}

	if r.StatusCode != http.StatusOK {
		fmt.Printf("GitHub API failed with status: %s\n", r.Status)
		os.Exit(1)
	}

	var events event.Events
	err = json.NewDecoder(r.Body).Decode(&events)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Bad JSON: %s\n", err)
		os.Exit(1)
	}

	for _, event := range events {
		event.ProcessEvent()
	}
}
