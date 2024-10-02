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
func main() {
	var image, repository string
	flag.StringVar(&image, "image", "", "docker image name")
	flag.StringVar(&repository, "repository", "library", "docker image name")
	flag.Parse()

	image = strings.TrimSpace(image)
	if image == "" && len(image) == 0 {
		fmt.Fprintf(os.Stderr, "missing required -image flag\n")
		os.Exit(2)
	}
	var search Search
	url := "https://hub.docker.com/v2/namespaces/" + repository + "/repositories/" + image + "/tags?page=1&page_size=100"

	search.Tags(url)
	fmt.Println(len(search.Results))
}
