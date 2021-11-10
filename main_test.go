package main

import (
	"testing"

	assert "github.com/go-playground/assert/v2"
)

func TestCheckFileIsYaml(t *testing.T) {
	fn := func(filePath string) {
		checkFileIsYaml(filePath)
	}
	assert.PanicMatches(t, func() { fn("asdfasdf") }, "file is not yaml. file path: asdfasdf")
	assert.PanicMatches(t, func() { fn("asdfasdf.csv") }, "file is not yaml. file path: asdfasdf.csv")
	assert.PanicMatches(t, func() { fn("asdfasdf..yml") }, "file extension is wrong. file path: asdfasdf..yml")

	checkFileIsYaml("asdfasdf.yml")
	checkFileIsYaml("asdfasdf.yaml")

}

func TestExtractTag(t *testing.T) {
	fn := func(rawString string) {
		extractTag(rawString)
	}
	panicPhrase := "'image: your-image-name:your-tag' phrase not exists in yaml."

	assert.PanicMatches(t, func() { fn("asdfasd") }, panicPhrase)
	assert.PanicMatches(t, func() { fn("image:asdfasdf") }, panicPhrase)
	assert.PanicMatches(t, func() { fn("image: asdfasdf") }, panicPhrase)
	assert.PanicMatches(t, func() { fn("image: asdfasdf::111") }, panicPhrase)
	assert.PanicMatches(t, func() { fn("image: asdfasdf:") }, panicPhrase)
	assert.PanicMatches(t, func() { fn("image: asdfasdf::111:dfasd") }, panicPhrase)
	assert.PanicMatches(t, func() { fn("image:asdfasdf:111") }, panicPhrase)

	tag := extractTag("image: asdfasdf:111")
	assert.Equal(t, tag, "111")
}

func TestGetOldTags(t *testing.T) {
	fn := func(url string) {
		getOldTags(url)
	}

	assert.PanicMatches(t, func() { fn("google.com") }, `Get "google.com": unsupported protocol scheme ""`)
	assert.PanicMatches(t, func() { fn("http://google.com") }, `invalid character '<' looking for beginning of value`)

	url := "https://registry.hub.docker.com/v1/repositories/nginx/tags"
	oldTags := getOldTags(url)
	assert.Equal(t, len(oldTags) > 0, true)
}
