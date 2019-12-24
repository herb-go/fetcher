package fetcher

import (
	"bytes"
	"net/http"
	"net/url"
	"testing"
)

func TestMisc(t *testing.T) {
	url, err := url.Parse("abcs")
	if err != nil {
		t.Fatal(err)
	}
	url2 := CloneURL(url)
	if url2 == url {
		t.Fatal(url2)
	}
	if url2.String() != url.String() {
		t.Fatal(url2)
	}
	header := http.Header{}
	header.Set("abc", "123")
	header2 := CloneHeader(header)
	if &header == &header2 {
		t.Fatal(header2)
	}
	data := bytes.NewBuffer(nil)
	data2 := bytes.NewBuffer(nil)
	err = header.Write(data)
	if err != nil {
		t.Fatal(err)
	}
	err = header2.Write(data2)
	if err != nil {
		t.Fatal(err)
	}
	if bytes.Compare(data.Bytes(), data2.Bytes()) != 0 {
		t.Fatal(data2)
	}
	header3 := http.Header{}
	header3.Set("cde", "456")
	MergeHeader(header3, header)
	if header3.Get("abc") != "123" || header3.Get("cde") != "456" {
		t.Fatal(header3)
	}
	builders := []func(*http.Request) error{
		nil,
		func(*http.Request) error {
			return nil
		},
	}
	builders2 := CloneRequestBuilders(builders)
	if len(builders2) != 2 || builders2[0] != nil {
		t.Fatal(builders2)
	}
	builders[0] = func(*http.Request) error {
		return nil
	}
	if builders2[0] != nil {
		t.Fatal(builders2)
	}
}
