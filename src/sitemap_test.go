package smeago

import (
	"bytes"
	"testing"
)

func TestWrite(t *testing.T) {
	expected := `<?xml version="1.0" encoding="UTF-8"?>
<urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">
	<url>
		<loc>http://example.com/foo</loc>
	</url>
	<url>
		<loc>http://example.com/bar</loc>
	</url>
</urlset>`

	s := &Sitemap{
		Path:  "http://example.com",
		Links: []string{"/foo", "/bar"},
	}
	buff := new(bytes.Buffer)
	err := s.Write(buff)

	if err != nil {
		t.Error(err.Error())
	}

	if buff.String() != expected {
		t.Error("expected\n", expected, "got\n", buff.String())
	}
}

func BenchmarkWriteSitemap(b *testing.B) {
	s := &Sitemap{
		Path:  "http://example.com",
		Links: []string{"/foo", "/bar"},
	}
	buff := new(bytes.Buffer)
	for i := 0; i < b.N; i++ {
		s.Write(buff)
		buff.Reset()
	}
}
