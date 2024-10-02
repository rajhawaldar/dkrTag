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
type Response struct {
	Next string `json:"next"`
	Tags []Tag  `json:"results"`
}

func (s *Response) GetTags(url string) {
	var result Response
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
	s.Tags = append(s.Tags, result.Tags...)
	if result.Next != "" {
		s.GetTags(result.Next)
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
	var Response Response
	url := "https://hub.docker.com/v2/namespaces/" + repository + "/repositories/" + image + "/tags?page=1&page_size=100"
	Response.GetTags(url)
	fmt.Println(len(Response.Tags))
}
