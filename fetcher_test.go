package fetcher

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/url"
	"testing"
)

func TestFetcher(t *testing.T) {
	c := &http.Client{}
	body := bytes.NewBufferString("123")
	var err error
	f := New()
	f.URL, err = url.Parse("127.0.0.01")
	if err != nil {
		t.Fatal(err)
	}
	f.Header.Set("k1", "v1")
	f.Doer = c
	f.Method = "TESTMDTOD"
	f.Body = body
	f.AppendBuilder(func(r *http.Request) error {
		r.SetBasicAuth("u", "p")
		return nil
	})
	req, doer, err := f.Raw()
	if err != nil {
		t.Fatal(f)
	}
	if doer != c {
		t.Fatal(doer)
	}

	if req.URL.String() != f.URL.String() {
		t.Fatal(req)
	}
	if req.Method != f.Method {
		t.Fatal(req)
	}
	bs, err := ioutil.ReadAll(req.Body)
	if err != nil {
		t.Fatal(f)
	}
	if string(bs) != "123" {
		t.Fatal(req)
	}
	u, p, ok := req.BasicAuth()
	if u != "u" || p != "p" || !ok {
		t.Fatal(u, p, ok)
	}
}
