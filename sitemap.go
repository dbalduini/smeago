package main

import (
	"encoding/xml"
	"io"
)

type Sitemap struct {
	XMLName xml.Name `xml:"urlset"`
	Urls    []string `xml:"url>loc"`
}

func (u *Sitemap) Write(w io.Writer) error {
	if _, err := w.Write([]byte(xml.Header)); err != nil {
		return err
	}
	output, err := xml.MarshalIndent(u, "", "\t")
	if err != nil {
		return err
	}
	_, err = w.Write(output)
	return err
}
