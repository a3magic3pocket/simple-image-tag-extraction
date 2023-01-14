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
	Name string `json:"name"`
}

type tagsResp struct {
	Count    int         `json:"count"`
	Next     interface{} `json:"next"`
	Previous interface{} `json:"previous"`
	Results  []struct {
		Name string `json:"name"`
	} `json:"results"`
}

func getRequest(url string) []byte {
	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()

	r, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	return r
}

func getOldTags(url string) oldTags {
	page := 1
	pageSize := 100
	originURL := url
	tags := oldTags{}

	for {
		url = fmt.Sprintf("%s?page=%d&page_size=%d", originURL, page, pageSize)

		resp := getRequest(url)

		tResp := tagsResp{}
		err := json.Unmarshal(resp, &tResp)
		if err != nil {
			panic(err)
		}

		if len(tResp.Results) <= 0 {
			break
		}

		tags = append(tags, tResp.Results...)

		page += 1
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

func refineTagsURL(tagsURL string) string {
	splited := strings.Split(tagsURL, "?")

	return splited[0]
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
			url: https://hub.docker.com/v2/namespaces/<namespace>/repositories/<repository>/tags
			ex) https://hub.docker.com/v2/namespaces/library/repositories/nginx/tags`,
	)
	flag.Parse()
	if *filePath == "" || *tagsURL == "" {
		if *filePath == "" {
			fmt.Println("filePath is empty")
		}
		if *tagsURL == "" {
			fmt.Println("tagsURL is empty")
		}

		flag.PrintDefaults()
		return
	}

	refinedTagsURL := refineTagsURL(*tagsURL)

	checkFileIsYaml(*filePath)

	data, err := ioutil.ReadFile(*filePath)
	if err != nil {
		panic(err)
	}

	tag := extractTag(string(data))

	oldTags := getOldTags(refinedTagsURL)
	for _, oldTag := range oldTags {
		if oldTag.Name == tag {
			panic(fmt.Sprintf("'%s' is already released.", tag))
		}
	}

	fmt.Print(tag)
}
