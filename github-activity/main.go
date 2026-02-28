package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
)

type Event struct {
	Type string `json:"type"`
	Repo struct {
		Name string `json:"name"`
	} `json:"repo"`
	Payload json.RawMessage `json:"payload"`
}

func fetch_event(username string, perPage int) ([]Event, error) {
	u := url.URL{
		Scheme: "https",
		Host:   "api.github.com",
		Path:   "/users/" + username + "/events",
	}
	q := u.Query()
	if perPage > 0 {
		q.Set("per_page", strconv.Itoa(perPage))
	}
	u.RawQuery = q.Encode()
	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", "my-cli-test")
	req.Header.Set("Accept", "application/vnd.github+json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("github api error: %s", resp.Status)
	}

	var events []Event
	if err := json.NewDecoder(resp.Body).Decode(&events); err != nil {
		return nil, err
	}

	return events, nil
}

func format_event(e Event) string {
	switch e.Type {
	case "PushEvent":
		return fmt.Sprintf("Pushed to %s", e.Repo.Name)
	case "IssuesEvent":
		return fmt.Sprintf("Interacted with an issue in %s", e.Repo.Name)
	case "WatchEvent":
		return fmt.Sprintf("Starred %s", e.Repo.Name)
	case "ForkEvent":
		return fmt.Sprintf("Forked %s", e.Repo.Name)
	case "CreateEvent":
		return fmt.Sprintf("Created something in %s", e.Repo.Name)
	case "PullRequestEvent":
		return fmt.Sprintf("Worked on a pull request in %s", e.Repo.Name)
	case "PullRequestReviewEvent":
		return fmt.Sprintf("Reviewed a pull request in %s", e.Repo.Name)
	case "ReleaseEvent":
		return fmt.Sprintf("Released something in %s", e.Repo.Name)
	case "DeleteEvent":
		return fmt.Sprintf("Deleted something in %s", e.Repo.Name)
	default:
		return fmt.Sprintf("%s in %s", strings.TrimSuffix(e.Type, "Event"), e.Repo.Name)
	}
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("github-activity need: username")
		return
	}
	username := os.Args[1]
	events, err := fetch_event(username, 10)
	if err != nil {
		fmt.Println("err:", err)
		return
	}
	for i, e := range events {
		s := format_event(e)
		fmt.Printf("%d): %s\n", i+1, s)
	}
}
