package fetcher

import (
	"testing"
)

func TestPreset(t *testing.T) {
	s := newEchoServer()
	defer s.Close()
	var sc = &Server{
		ServerInfo: ServerInfo{
			URL: s.URL,
		},
	}
	preset := MustPreset(sc)
	resp, err := preset.Fetch()
	if err != nil {
		t.Fatal(err)
	}
	resp.Body.Close()
	if resp.StatusCode != 200 {
		t.Fatal(resp)
	}
}

func TestPresetMethods(t *testing.T) {
	var err error
	cmd := CommandFunc(func(*Fetcher) error { return nil })
	p0 := NewPreset()
	if len(p0.Commands()) != 0 {
		t.Fatal(p0)
	}
	p := BuildPreset(cmd, cmd)
	if len(p.Commands()) != 2 {
		t.Fatal(p)
	}
	p2 := p.Append(p)
	if len(p.Commands()) != 2 || len(p2.Commands()) != 4 {
		t.Fatal(p, p2)
	}
	pnil := p.Append(BuildPreset(nil, nil))
	cmds := pnil.Commands()
	if len(cmds) != 4 || cmds[0] == nil || cmds[1] == nil || cmds[2] != nil || cmds[3] != nil {
		t.Fatal(pnil)
	}
	p = NewPreset().EndPoint("/pathprefix", "TESTMETHOD")
	f := New()
	f.URL.Path = "raw"
	err = p.Exec(f)
	if err != nil {
		t.Fatal(err)
	}
	if f.URL.Path != "raw/pathprefix" || f.Method != "TESTMETHOD" {
		t.Fatal(f)
	}
}
