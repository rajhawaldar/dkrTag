package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
)

const API string = "https://hub.docker.com/v2/namespaces/"

type Image struct {
	Architecture string `json:"architecture"`
	Features     string `json:"os_features"`
	Digest       string `json:"digest"`
	Os           string `json:"os"`
	Size         int    `json:"size"`
	Status       string `json:"status"`
}
type Tag struct {
	Name       string  `json:"name"`
	Tag_status string  `json:"tag_status"`
	Images     []Image `json:"images"`
}
type Search struct {
	Next    string `json:"next"`
	Results []Tag  `json:"results"`
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

	return false
}
func main() {
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
	// /v2/namespaces/{namespace}/repositories/{repository}/tags
	search.Tags(tagsURL)
	fmt.Println(len(search.Results))
}
