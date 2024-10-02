package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh/spinner"
)

const API string = "https://hub.docker.com/v2/namespaces/"

type Tag struct {
	Name        string `json:"name"`
	Tag_status  string `json:"tag_status"`
	LastUpdated string `json:"last_updated"`
	FullSize    int    `json:"full_size"`
}
type Search struct {
	Next    string `json:"next"`
	Results []Tag  `json:"results"`
}
type Model struct {
	tags []Tag
	list list.Model
	err  error
}

func (t Tag) FilterValue() string {
	return t.Name
}
func (t Tag) Title() string {
	return t.Name
}

func (m *Model) initList(width, height int) {
	m.list = list.New([]list.Item{}, list.NewDefaultDelegate(), width, height)
	m.list.Title = "Tags"
	var items []list.Item
	for _, tag := range m.tags {
		items = append(items, tag)
	}
	m.list.SetItems(items)
}
func (m Model) Init() tea.Cmd {
	return nil
}
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.initList(msg.Width, msg.Height)
	}
	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m Model) View() string {
	return m.list.View()
}
func (s *Search) Tags(url string) {
	var result Search
	resp, err := http.Get(url)
	if err != nil {
		log.Panicln(err.Error())
	}
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Panicln(err.Error())
	}
	errJson := json.Unmarshal(data, &result)
	if errJson != nil {
		log.Panicln(err.Error())
	}
	s.Results = append(s.Results, result.Results...)
	if result.Next != "" {
		s.Tags(result.Next)
	}
}

func (s *Search) doTagsExist(url string) bool {
	resp, err := http.Get(url)
	if err != nil {
		log.Panicln(err.Error())
	}
	if resp.StatusCode == 200 {
		return true
	}
	return false
}
func tagsToList(result []Tag) []string {
	var tags []string = []string{}
	for i := len(result) - 1; i >= 0; i-- {
		date, _ := time.Parse(time.RFC3339, result[i].LastUpdated)
		tags = append(tags, fmt.Sprintf("%-30s %-10s %-15s %-15s", result[i].Name, result[i].Tag_status, fmt.Sprintf("%.2f MB", float64(result[i].FullSize)/(1<<20)), date.Format("2006-01-02")))
	}
	return tags
}

func main() {
	var wg sync.WaitGroup
	var namespace, repository string
	flag.StringVar(&repository, "repository", "", "docker repository name, example: nginx, bash, ubuntu")
	flag.StringVar(&namespace, "namespace", "library", "your docker namespace")
	flag.Parse()

	repository = strings.TrimSpace(repository)
	if repository == "" && len(repository) == 0 {
		fmt.Fprintf(os.Stderr, "missing required -repository flag\n")
		os.Exit(2)
	}
	var search Search
	tagsURL := API + namespace + "/repositories/" + repository + "/tags?page=1&page_size=100"
	tagExistURL := API + namespace + "/repositories/" + repository + "/tags"
	if !search.doTagsExist(tagExistURL) {
		fmt.Println("Repository does not contain any tags.")
		os.Exit(0)
	}
	wg.Add(1)
	go func() {
		defer wg.Done()
		search.Tags(tagsURL)
	}()
	makeSpinnerWait := func() {
		wg.Wait()
	}
	err := spinner.New().
		Title("Fetching all image tags...").
		Action(makeSpinnerWait).
		Run()
	if err != nil {
		log.Panicln(err.Error())
	}

	sort.Slice(search.Results, func(i, j int) bool {
		timeI, err1 := time.Parse(time.RFC3339, search.Results[i].LastUpdated)
		timeJ, err2 := time.Parse(time.RFC3339, search.Results[j].LastUpdated)

		if err1 != nil || err2 != nil {
			return false
		}
		return timeI.Before(timeJ)
	})
	tags := tagsToList(search.Results)
	for _, tag := range tags {
		fmt.Println(tag)
	}
}
