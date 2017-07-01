package main

import (
	"bytes"
	"testing"
)

func TestWriteSitemap(t *testing.T) {
	expected := `<?xml version="1.0" encoding="UTF-8"?>
<urlset>
	<url>
		<loc>foo</loc>
		<loc>bar</loc>
	</url>
</urlset>`

	s := &Sitemap{}
	s.Urls = []string{"foo", "bar"}

	b := new(bytes.Buffer)
	err := s.Write(b)

	if err != nil {
		t.Error(err.Error())
	}

	if b.String() != expected {
		t.Error("expected\n", expected, "got\n", b.String())
	}
}
