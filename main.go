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
type Result struct {
	Name       string  `json:"name"`
	Tag_status string  `json:"tag_status"`
	Images     []Image `json:"images"`
}
type Tag struct {
	Next    string   `json:"next"`
	Results []Result `json:"results"`
}

func GetTags(url string, tag *Tag) {
	var tags Tag
	enc := json.NewEncoder(os.Stdout)
	enc.SetEscapeHTML(false)
	resp, err := http.Get(url)
	if err != nil {
		log.Panicln(err.Error())
	}
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Panicln(err.Error())
	}
	errJson := json.Unmarshal(data, &tags)
	if errJson != nil {
		log.Panicln(err.Error())
	}
	tag.Results = append(tag.Results, tags.Results...)
	if tags.Next != "" {
		GetTags(tags.Next, tag)
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
	var tags Tag
	url := "https://hub.docker.com/v2/namespaces/" + repository + "/repositories/" + image + "/tags?page=1&page_size=100"
	GetTags(url, &tags)
	fmt.Println(len(tags.Results))
}
