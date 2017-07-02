package smeago

import (
	"bytes"
	"io"
)

// Common used bytes
var (
	header = []byte("<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n")
	root   = []byte("<urlset xmlns=\"http://www.sitemaps.org/schemas/sitemap/0.9\">")
	pre    = []byte("\n\t<url>\n\t\t<loc>")
	pos    = []byte("</loc>\n\t</url>")
)

// WriteSitemap writes links into writer with tab indentation
func WriteSitemap(w io.Writer, links []string) error {
	buff := new(bytes.Buffer)
	buff.Write(header)
	buff.Write(root)
	for _, loc := range links {
		buff.Write(pre)
		buff.WriteString(loc)
		buff.Write(pos)
	}
	buff.WriteString("\n</urlset>")
	_, err := w.Write(buff.Bytes())
	return err
}
