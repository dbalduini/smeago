package smeago

import (
	"bytes"
	"testing"
)

func TestWriteSitemap(t *testing.T) {
	expected := `<?xml version="1.0" encoding="UTF-8"?>
<urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">
	<url>
		<loc>foo</loc>
	</url>
	<url>
		<loc>bar</loc>
	</url>
</urlset>`

	links := []string{"foo", "bar"}
	buff := new(bytes.Buffer)
	err := WriteSitemap(buff, links)

	if err != nil {
		t.Error(err.Error())
	}

	if buff.String() != expected {
		t.Error("expected\n", expected, "got\n", buff.String())
	}
}

func BenchmarkWriteSitemap(b *testing.B) {
	links := []string{"foo", "bar"}
	buff := new(bytes.Buffer)
	for i := 0; i < b.N; i++ {
		WriteSitemap(buff, links)
		buff.Reset()
	}
}
