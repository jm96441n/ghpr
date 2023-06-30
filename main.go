package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"

	tea "github.com/charmbracelet/bubbletea"

	retryablehttp "github.com/hashicorp/go-retryablehttp"
)

func main() {
	var (
		handle      string
		tokenEnvVar string
	)

	flag.StringVar(&handle, "handle", "", "your github handle")
	flag.StringVar(&tokenEnvVar, "tokenEnv", "", "override for env var for your github token, defaults to GITHUB_ACCESS_TOKEN")
	flag.Parse()

	didError := false
	if handle == "" {
		fmt.Println("You must supply your github handle")
		didError = true
	}

	if tokenEnvVar == "" {
		tokenEnvVar = "GITHUB_ACCESS_TOKEN"
	}
	token := os.Getenv(tokenEnvVar)
	if token == "" {
		fmt.Printf("Did not find token in env at %s\n", tokenEnvVar)
		didError = true
	}
	if didError {
		os.Exit(1)
	}

	if _, err := tea.NewProgram(model{handle: handle}).Run(); err != nil {
		fmt.Printf("Uh oh, there was an error: %v\n", err)
		os.Exit(1)
	}
}

type model struct {
	handle string
	token  string
	repos  []string
	status int
	err    error
}

type ghResp struct {
	status int
	Items  []item `json:"items"`
}

type item struct {
	URL string `json:"html_url"`
}

func (m model) Init() tea.Cmd {
	return getGHInfo(m.token, m.handle)
}

func getGHInfo(token, handle string) tea.Cmd {
	return func() tea.Msg {
		client := retryablehttp.NewClient()
		client.Logger = nil

		req, err := retryablehttp.NewRequest("GET", fmt.Sprintf("https://api.github.com/search/issues?q=assignee:%s+state:open", handle), nil)
		if err != nil {
			return err
		}

		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

		resp, err := client.Do(req)
		if err != nil {
			return err
		}

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return err
		}

		var parsed ghResp
		err = json.Unmarshal(body, &parsed)
		if err != nil {
			return err
		}

		parsed.status = resp.StatusCode

		return parsed
	}
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case ghResp:
		m.status = msg.status
		m.repos = make([]string, 0, len(msg.Items))
		for _, item := range msg.Items {
			m.repos = append(m.repos, item.URL)
		}

		return m, tea.Quit
	case error:
		m.err = msg
		return m, tea.Quit
	}
	return m, nil
}

func (m model) View() string {
	if m.err != nil {
		return fmt.Sprintf("something went wrong: %v\n\n", m.err)
	}

	s := "Checking github..."

	if m.status > 0 {
		if len(m.repos) == 0 {
			return "No open pull requests assigned to you!\n\n"
		}
		s = "Pull requests you're assigned to:\n"

		for _, url := range m.repos {
			s = fmt.Sprintf("%s%s\n", s, url)
		}
	}
	return s + "\n\n"
}
