package fetcher

import (
	"testing"
)

func TestMustPreset(t *testing.T) {

}
func TestPreset(t *testing.T) {
	s := newEchoServer()
	defer s.Close()
	var sc = &Server{
		ServerInfo: ServerInfo{
			URL: s.URL,
		},
	}
	preset, err := sc.CreatePreset()
	if err != nil {
		t.Fatal(err)
	}
	resp, err := preset.Fetch()
	if err != nil {
		t.Fatal(err)
	}
	resp.Body.Close()
	if resp.StatusCode != 200 {
		t.Fatal(resp)
	}
}
