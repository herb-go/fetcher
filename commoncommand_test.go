package fetcher

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"testing"
)

type rp struct {
}

func (rp *rp) BuildRequest(*http.Request) error {
	return nil
}

type hbp struct {
}

func (hbp *hbp) BuildHeader(h http.Header) error {
	h.Set("k1", "v1")
	return nil
}

type mbp struct {
}

func (mbp *mbp) RequestMethod() (string, error) {
	return "MethodBuilderProvider", nil
}
func TestCommonCommand(t *testing.T) {
	var err error
	f := New()
	url, err := url.Parse("http://127.0.0.1/{{path}}/")
	err = URL(url).Exec(f)
	if err != nil {
		t.Fatal(err)
	}
	if f.URL != url || f.URL.Path != "/{{path}}/" || f.URL.Host != "127.0.0.1" {
		t.Fatal(f)
	}
	err = Replace("{{path}}", "replacement").Exec(f)
	if err != nil {
		t.Fatal(err)
	}
	if f.URL != url || f.URL.Path != "/replacement/" {
		t.Fatal(f)
	}
	err = Host("localhost").Exec(f)
	if err != nil {
		t.Fatal(err)
	}
	if f.URL != url || f.URL.Host != "localhost" {
		t.Fatal(f)
	}

	f = New()
	err = Post.Exec(f)
	if err != nil {
		t.Fatal(err)
	}
	if f.Method != "POST" {
		t.Fatal(err)
	}
	header := http.Header{}
	header.Set("k1", "v1")
	header.Set("k2", "v2")
	header2 := http.Header{}
	header2.Set("k1", "newv1")
	header2.Set("k3", "v3")
	f = New()
	f.Header = header
	err = Header(header2).Exec(f)
	if err != nil {
		t.Fatal(err)
	}
	if f.Header.Get("k1") != "newv1" || f.Header.Get("k2") != "v2" || f.Header.Get("k3") != "v3" {
		t.Fatal(f)
	}
	f = New()
	if f.URL.Query().Get("k1") != "" {
		t.Fatal(f)
	}
	err = SetQuery("k1", "v1").Exec(f)
	if err != nil {
		t.Fatal(err)
	}
	if f.URL.Query().Get("k1") != "v1" {
		t.Fatal(f)
	}
	f = New()
	if len(f.Builders) != 0 {
		t.Fatal(f)
	}
	err = RequestBuilder(&rp{}).Exec(f)
	if err != nil {
		t.Fatal(err)
	}
	if len(f.Builders) != 1 {
		t.Fatal(f)
	}
	f = New()
	if f.Doer != nil {
		t.Fatal(err)
	}
	err = SetDoer(http.DefaultClient).Exec(f)
	if err != nil {
		t.Fatal(err)
	}
	if f.Doer != http.DefaultClient {
		t.Fatal(err)
	}
	f = New()
	if f.Header.Get("k1") != "" {
		t.Fatal(f)
	}
	err = HeaderBuilder(&hbp{}).Exec(f)
	if f.Header.Get("k1") != "v1" {
		t.Fatal(f)
	}
	f = New()
	if f.Method != "" {
		t.Fatal(f)
	}
	err = MethodBuilder(&mbp{}).Exec(f)
	if f.Method != "MethodBuilderProvider" {
		t.Fatal(f)
	}

}

func TestCommonCommantecho(t *testing.T) {
	var err error
	s := newEchoServer()
	defer s.Close()
	var sc = &Server{
		ServerInfo: ServerInfo{
			URL: s.URL,
		},
	}
	preset := MustPreset(sc)
	f := New()
	err = preset.Exec(f)
	if err != nil {
		t.Fatal(err)
	}
	err = JSONBody("12345").Exec(f)
	if err != nil {
		t.Fatal(err)
	}
	var data string
	resp, err := f.Fetch()
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	bs, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	err = json.Unmarshal(bs, &data)
	if err != nil {
		t.Fatal(err)
	}
	if data != "12345" {
		t.Fatal(data)
	}
	f = New()
	err = preset.Exec(f)
	if err != nil {
		t.Fatal(err)
	}
	err = JSONBody(nil).Exec(f)
	if err != nil {
		t.Fatal(err)
	}
	if f.Body != nil {
		t.Fatal(f)
	}
	f = New()
	body := bytes.NewBufferString("buf")
	err = preset.Exec(f)
	if err != nil {
		t.Fatal(err)
	}
	err = Body(body).Exec(f)
	if err != nil {
		t.Fatal(err)
	}
	if f.Body != body {
		t.Fatal(f)
	}

}
