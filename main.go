package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"
)

type oldTags []struct {
	Layer string `json:"layer"`
	Name  string `json:"name"`
}

func getOldTags(url string) oldTags {
	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()

	r, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	tags := oldTags{}
	err = json.Unmarshal(r, &tags)
	if err != nil {
		panic(err)
	}

	return tags
}

func checkFileIsYaml(filePath string) {
	extension := strings.Split(filePath, ".")
	lastIndex := len(extension) - 1
	if len(extension) > 2 {
		panic(fmt.Sprintf("file extension is wrong. file path: %s", filePath))
	}
	if extension[lastIndex] != "yaml" && extension[lastIndex] != "yml" {
		panic(fmt.Sprintf("file is not yaml. file path: %s", filePath))
	}
}

func extractTag(rawYaml string) string {
	r := regexp.MustCompile(`image: ([^\s:]+):([^\s:]+)`)
	matched := r.FindStringSubmatch(rawYaml)

	if len(matched) < 3 {
		panic("'image: your-image-name:your-tag' phrase not exists in yaml.")
	}

	return matched[2]
}

func main() {
	filePath := flag.String(
		"f",
		"",
		`yaml file. it have to contain 'image: your-image-name:your-tag' phrase
			ex) mock-deployment.yml`,
	)
	tagsURL := flag.String(
		"tu",
		"",
		`registry server tags url
			ex) https://registry.hub.docker.com/v1/repositories/nginx/tags`,
	)
	flag.Parse()
	if *filePath == "" && *tagsURL == "" {
		flag.PrintDefaults()
		return
	}

	checkFileIsYaml(*filePath)

	data, err := ioutil.ReadFile(*filePath)
	if err != nil {
		panic(err)
	}

	tag := extractTag(string(data))

	oldTags := getOldTags(*tagsURL)
	for _, oldTag := range oldTags {
		if oldTag.Name == tag {
			panic(fmt.Sprintf("'%s' is already released.", tag))
		}
	}

	fmt.Print(tag)
}
