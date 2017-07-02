package smeago

import (
	"html"
	"io"
	"io/ioutil"
	"regexp"
)

var (
	anchorRE = regexp.MustCompile("<a.*href=\"([^\"]*)\"[^>]*>")
)

type Result struct {
	Links []string
}

func ParseRequestBody(rd io.Reader) (*Result, error) {
	r := &Result{}

	b, err := ioutil.ReadAll(rd)
	if err != nil {
		return r, err
	}

	s := string(b)
	links := getLinks(s)
	lc := links[:0]
	// Only internal links
	for _, l := range links {
		if l[0] == '/' {
			lc = append(lc, decodeURL(l))
		}
	}
	r.Links = lc
	return r, nil
}

func decodeURL(s string) string {
	return html.UnescapeString(s)
}

func getLinks(s string) []string {
	matches := anchorRE.FindAllStringSubmatch(s, -1)
	links := make([]string, 0)
	for _, a := range matches {
		links = append(links, a[1])
	}
	return links
}
