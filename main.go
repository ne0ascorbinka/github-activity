package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
)

type EventType string

const (
	EVENT_TYPE_CREATE            EventType = "CreateEvent"
	EVENT_TYPE_PUSH              EventType = "PushEvent"
	EVENT_TYPE_ISSUES            EventType = "IssuesEvent"
	EVENT_TYPE_WATCH             EventType = "WatchEvent"
	EVENT_TYPE_PR                EventType = "PullRequestEvent"
	EVENT_TYPE_PR_REVIEW         EventType = "PullRequestReviewEvent"
	EVENT_TYPE_PR_REVIEW_COMMENT EventType = "PullRequestReviewCommentEvent"
	EVENT_TYPE_RELEASE           EventType = "ReleaseEvent"
)

type Event struct {
	// ID int `json:"id"`
	Type EventType `json:"type"`
	Repo struct {
		Name string `json:"name"`
	} `json:"repo"`
	Payload struct {
		Action  string `json:"action"`
		RefType string `json:"ref_type"`
	} `json:"payload"`
}

type Events []Event

func (e Event) ProcessPushEvent() {
	fmt.Printf("Pushed to %s\n", e.Repo.Name)
}

func (e Event) ProcessIssuesEvent() {
	switch e.Payload.Action {
	case "opened":
		fmt.Printf("Opened a new issue in %s\n", e.Repo.Name)
	case "closed":
		fmt.Printf("Closed an issue in %s\n", e.Repo.Name)
	case "reopened":
		fmt.Printf("Reopened an issue in %s\n", e.Repo.Name)
	}
}

func (e Event) ProcessCreateEvent() {
	switch e.Payload.RefType {
	case "branch":
		fmt.Printf("Created a branch in %s\n", e.Repo.Name)
	case "tag":
		fmt.Printf("Created a tag in %s\n", e.Repo.Name)
	case "repository":
		fmt.Printf("Created a repository %s\n", e.Repo.Name)
	}
}

func (e Event) ProcessWatchEvent() {
	action := e.Payload.Action
	switch action {
	case "started":
		fmt.Printf("Starred %s\n", e.Repo.Name)
	default:
		fmt.Printf("Unknown WatchEvent action: %s\n", action)
	}
}

func (e Event) ProcessPullRequestEvent() {
	action := e.Payload.Action
	switch action {
	case "opened":
		fmt.Printf("Opened a pull request in %s\n", e.Repo.Name)
	case "closed":
		fmt.Printf("Closed a pull request in %s\n", e.Repo.Name)
	case "merged":
		fmt.Printf("Merged a pull request in %s\n", e.Repo.Name)
	case "reopened":
		fmt.Printf("Reopened a pull request in %s\n", e.Repo.Name)
	case "assigned":
		fmt.Printf("Assigned a user to a pull request in %s\n", e.Repo.Name)
	case "unassigned":
		fmt.Printf("Unassigned a user from a pull request in %s\n", e.Repo.Name)
	case "labeled":
		fmt.Printf("Added a label to a pull request in %s\n", e.Repo.Name)
	case "unlabeled":
		fmt.Printf("Removed a label from a pull request in %s\n", e.Repo.Name)
	}
}

func (e Event) ProcessPullRequestReviewEvent() {
	action := e.Payload.Action
	switch action {
	case "created":
		fmt.Printf("Created a review to a pull request in %s\n", e.Repo.Name)
	case "updated":
		fmt.Printf("Updated a review to a pull request in %s\n", e.Repo.Name)
	case "dismissed":
		fmt.Printf("Dismissed a review from a pull request in %s\n", e.Repo.Name)
	}
}

func (e Event) ProcessPullRequestReviewCommentEvent() {
	action := e.Payload.Action
	switch action {
	case "created":
		fmt.Printf("Commented on a pull request in %s\n", e.Repo.Name)
	default:
		fmt.Printf("Unknown PullRequestReviewCommentEvent action: %q\n", action)
	}
}

func (e Event) ProcessReleaseEvent() {
	fmt.Printf("Published a release in %s\n", e.Repo.Name)
}

func (e Event) ProcessEvent() {
	switch e.Type {
	case EVENT_TYPE_PUSH:
		e.ProcessPushEvent()
	case EVENT_TYPE_ISSUES:
		e.ProcessIssuesEvent()
	case EVENT_TYPE_CREATE:
		e.ProcessCreateEvent()
	case EVENT_TYPE_WATCH:
		e.ProcessWatchEvent()
	case EVENT_TYPE_PR:
		e.ProcessPullRequestEvent()
	case EVENT_TYPE_PR_REVIEW:
		e.ProcessPullRequestReviewEvent()
	case EVENT_TYPE_PR_REVIEW_COMMENT:
		e.ProcessPullRequestReviewCommentEvent()
	case EVENT_TYPE_RELEASE:
		e.ProcessReleaseEvent()
	default:
		fmt.Printf("Skipping unknown event of type %s\n", e.Type)
	}
}

const url = "https://api.github.com/users/%s/events"

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "Script takes username as argument")
		os.Exit(1)
	}
	username := os.Args[1]

	r, err := http.Get(fmt.Sprintf(url, username))
	if err != nil {
		log.Fatalf("Error fetching GitHub data: %v", err)
	}
	defer r.Body.Close()

	var events Events
	err = json.NewDecoder(r.Body).Decode(&events)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Bad JSON: %s\n", err)
		os.Exit(1)
	}

	for _, event := range events {
		event.ProcessEvent()
	}
}
