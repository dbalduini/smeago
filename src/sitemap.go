package smeago

import (
	"bufio"
	"bytes"
	"io"
	"os"
	"sort"
)

// Common used bytes
var (
	header = []byte("<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n")
	root   = []byte("<urlset xmlns=\"http://www.sitemaps.org/schemas/sitemap/0.9\">")
	pre    = []byte("\n\t<url>\n\t\t<loc>")
	pos    = []byte("</loc>\n\t</url>")
)

type Sitemap struct {
	Filename string
	Links    []string
	Path     string
}

// WriteToFile writes the sitemap into a file
func (s *Sitemap) WriteToFile(sortLinks bool) error {
	if sortLinks {
		sort.Strings(s.Links)
	}

	f, err := os.Create(s.Filename)
	if err != nil {
		return err
	}
	defer f.Close()

	wd := bufio.NewWriter(f)
	err = s.Write(wd)
	if err != nil {
		return err
	}
	wd.Flush()
	return nil
}

func (s *Sitemap) Write(w io.Writer) error {
	buff := new(bytes.Buffer)
	buff.Write(header)
	buff.Write(root)
	for _, loc := range s.Links {
		buff.Write(pre)
		buff.WriteString(s.Path)
		buff.WriteString(loc)
		buff.Write(pos)
	}
	buff.WriteString("\n</urlset>")
	_, err := w.Write(buff.Bytes())
	return err
}
