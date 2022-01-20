package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"os"
	"sort"
	"strings"
	"sync"
)

type note struct {
	Body string `json:"body"`
}

type mergeRequest struct {
	Url       string `json:"web_url"`
	Title     string `json:"title"`
	Iid       int    `json:"iid"`
	ProjectID int    `json:"project_id"`
	State     string `json:"state"`
	CreatedAt string `json:"created_at"`
	approved  string
}

const (
	Reset = "\033[0m"
	Red   = "\033[31m"
	Green = "\033[32m"
	Blue  = "\033[34m"
)

func r(l string) string {
	return strings.Join([]string{Red, l, Reset}, "")
}
func g(l string) string {
	return strings.Join([]string{Green, l, Reset}, "")
}
func b(l string) string {
	return strings.Join([]string{Blue, l, Reset}, "")
}

func printHelp() {
	fmt.Println(`
	My MRs lister

	Lists all open merge requests authored by you (along with review status and links)

	Usage:
		my-mr [-a|-t]

	Options:
	-t	Gitlab token, requires scopes "read_api"; overrides default environment variable GITLAB_RO_TOKEN
	-a	Print also merged (green) and closed (red) merge requests
	`)
}

func main() {
	var (
		token *string = flag.String("t", os.Getenv("GITLAB_RO_TOKEN"), "gitlab token (required scopes: read_api)")
		all   *bool   = flag.Bool("a", false, "true: check all MRs, false: check only opened MRs")
		help  *bool   = flag.Bool("h", false, "print help")
	)
	flag.Parse()

	if *help {
		printHelp()
		return
	}

	getMergeRequests := func(c *http.Client) ([]*mergeRequest, error) {
		url := "https://gitlab.com/api/v4/merge_requests"
		if !*all {
			url += "?state=opened"
		}
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			return nil, err
		}
		req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", *token))

		resp, err := c.Do(req)
		if err != nil {
			return nil, err
		}
		if resp.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("tempo-api: invalid status code for fetching schedule: %d", resp.StatusCode)
		}
		defer resp.Body.Close()

		mergeRequests := new([]*mergeRequest)
		err = json.NewDecoder(resp.Body).Decode(mergeRequests)
		return *mergeRequests, nil
	}

	getMergeRequestNotes := func(c *http.Client, p, i int) ([]note, error) {
		url := fmt.Sprintf("https://gitlab.com/api/v4/projects/%d/merge_requests/%d/notes", p, i)
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			return nil, err
		}
		req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", *token))

		resp, err := c.Do(req)
		if err != nil {
			return nil, err
		}
		if resp.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("tempo-api: invalid status code for fetching schedule: %d", resp.StatusCode)
		}
		defer resp.Body.Close()

		notes := new([]note)
		err = json.NewDecoder(resp.Body).Decode(notes)
		return *notes, nil
	}

	c := &http.Client{}
	res := "Your merge requests:\n=====\n"

	mrs, e := getMergeRequests(c)
	if e != nil {
		fmt.Printf("Failed to get merge requests: %v", e)
		return
	}
	wg := sync.WaitGroup{}
	wg.Add(len(mrs))
	for _, mr := range mrs {
		go func(mr *mergeRequest) {
			notes, e := getMergeRequestNotes(c, mr.ProjectID, mr.Iid)
			if e != nil {
				fmt.Printf("Failed to get merge request notes: %v", e)
				return
			}
			mr.approved = r("✘")
			for _, n := range notes {
				if strings.Contains(n.Body, "approved this merge request") {
					mr.approved = g("✔")
					break
				}
			}

			if *all {
				switch mr.State {
				case "closed":
					mr.Title = r(mr.Title)
				case "merged":
					mr.Title = g(mr.Title)
				}
			}

			wg.Done()
		}(mr)
	}
	wg.Wait()

	sort.Slice(mrs[:], func(i, j int) bool {
		return mrs[i].CreatedAt > mrs[j].CreatedAt
	})

	for i, mr := range mrs {
		res += fmt.Sprintf("%d: %s [ review: %s ] (%s)\n", i+1, mr.Title, mr.approved, b(mr.Url))
	}

	res += "=====\n"
	fmt.Print(res)
}
