package smeago

import (
	"reflect"
	"testing"
)

func TestDecodeUrl(t *testing.T) {
	const (
		e = "/produto?id=256"
	)
	u := decodeURL("/produto?id&#x3D;256")

	if u != e {
		t.Error("expected", e)
	}
}

func TestGetLinks(t *testing.T) {
	html := `
<html>
<div><a href="/foo">Foo</a></div>
<div>
	<a href="/bar">Bar</a>
</div>
</html>
`
	expected := []string{"/foo", "/bar"}
	links := getLinks(html)

	if !reflect.DeepEqual(links, expected) {
		t.Error("unexpected links:", links)
	}
}
