package fetcher

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/url"
	"strings"
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

type pp struct {
}

func (pp *pp) BuildParams(params url.Values) error {
	params.Add("id", "test")
	return nil
}
func TestCommonCommand(t *testing.T) {
	var err error
	f := New()
	u := "http://127.0.0.1/{{path}}/"
	err = URL(u).Exec(f)
	if err != nil {
		t.Fatal(err)
	}

	if f.URL.Path != "/{{path}}/" || f.URL.Host != "127.0.0.1" {
		t.Fatal(f)
	}

	err = Replace("{{path}}", "replacement").Exec(f)
	if err != nil {
		t.Fatal(err)
	}
	if f.URL.Path != "/replacement/" {
		t.Fatal(f)
	}
	err = Host("localhost").Exec(f)
	if err != nil {
		t.Fatal(err)
	}
	if f.URL.Host != "localhost" {
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
	if f.Header.Get("k1") != "" {
		t.Fatal(f)
	}
	err = SetHeader("k1", "v1").Exec(f)
	if err != nil {
		t.Fatal(err)
	}
	if f.Header.Get("k1") != "v1" {
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
	params := url.Values{}
	params.Set("k1", "v1")
	params.Set("k2", "v2")
	f = New()
	if f.URL.Query().Get("k1") != "" {
		t.Fatal(f)
	}
	if f.URL.Query().Get("k2") != "" {
		t.Fatal(f)
	}
	err = Params(params).Exec(f)
	if err != nil {
		t.Fatal(err)
	}
	if f.URL.Query().Get("k1") != "v1" {
		t.Fatal(f)
	}
	if f.URL.Query().Get("k2") != "v2" {
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
	if err != nil {
		panic(err)
	}
	if f.Header.Get("k1") != "v1" {
		t.Fatal(f)
	}
	f = New()
	if f.Method != "" {
		t.Fatal(f)
	}
	err = MethodBuilder(&mbp{}).Exec(f)
	if err != nil {
		panic(err)
	}
	if f.Method != "MethodBuilderProvider" {
		t.Fatal(f)
	}
	f = New()
	if f.URL.Query().Get("id") != "" {
		t.Fatal(f)
	}
	err = ParamsBuilder(&pp{}).Exec(f)
	if err != nil {
		panic(err)
	}
	if f.URL.Query().Get("id") != "test" {
		t.Fatal(f)
	}

	f = New()
	r, _, err := f.Raw()
	if err != nil {
		t.Fatal(err)
	}
	user, p, ok := r.BasicAuth()
	if ok || user != "" || p != "" {
		t.Fatal(u, p)
	}
	err = BasicAuth("user", "pw").Exec(f)
	if err != nil {
		panic(err)
	}
	r, _, err = f.Raw()
	if err != nil {
		t.Fatal(err)
	}
	user, p, ok = r.BasicAuth()
	if !ok || user != "user" || p != "pw" {
		t.Fatal(u, p)
	}
	f = New()
	pu, err := url.Parse("http://127.0.0.1")
	if err != nil {
		t.Fatal(err)
	}
	err = ParsedURL(pu).Exec(f)
	if err != nil {
		t.Fatal(err)
	}
	if f.URL.String() != pu.String() {
		t.Fatal()
	}
	f = New()
	f.URL.Path = "raw"
	err = PathPrefix("prefix").Exec(f)
	if err != nil {
		t.Fatal(err)
	}
	if f.URL.Path != "prefixraw" {
		t.Fatal(f)
	}
	f = New()
	f.URL.Path = "raw"
	err = PathSuffix("suffix").Exec(f)
	if err != nil {
		t.Fatal(err)
	}
	if f.URL.Path != "rawsuffix" {
		t.Fatal(f)
	}
	f = New()
	f.URL.Path = "raw"
	err = PathJoin("join").Exec(f)
	if err != nil {
		t.Fatal(err)
	}
	if f.URL.Path != "raw/join" {
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

func TestMultiPartWriter(t *testing.T) {
	w := NewMultiPartWriter()
	f := New()
	err := w.WriteFile("file", "filename", bytes.NewBufferString("content"))
	if err != nil {
		panic(err)
	}
	err = w.Close()
	if err != nil {
		panic(err)
	}
	err = w.Exec(f)
	if err != nil {
		panic(err)
	}
	if !strings.HasPrefix(f.Header.Get("Content-Type"), "multipart/form-data;") {
		t.Fatal(f)
	}
	form, err := multipart.NewReader(f.Body, w.Boundary()).ReadForm(20000)
	if err != nil {
		panic(err)
	}
	fs := form.File["file"]
	if len(fs) != 1 || fs[0].Filename != "filename" {
		t.Fatal(fs)
	}

}
