package fetcher

import (
	"bytes"
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
	resp, err = preset.FetchAndParse(Should200(nil))
	if err != nil {
		t.Fatal(err)
	}
	resp.Body.Close()
	var result string
	resp, err = preset.FetchWithBodyAndParse(bytes.NewBufferString("12345"), Should200(AsString(&result)))
	if err != nil {
		t.Fatal(err)
	}
	resp.Body.Close()
	if result != "12345" {
		t.Fatal(result)
	}
	result = ""
	resp, err = preset.FetchWithJSONBodyAndParse("12345", Should200(AsJSON(&result)))
	if err != nil {
		t.Fatal(err)
	}
	resp.Body.Close()
	if result != "12345" {
		t.Fatal(result)
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
	p = NewPreset().EndPoint("TESTMETHOD", "/pathprefix")
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

func TestOrder(t *testing.T) {
	cmd1 := CommandFunc(func(f *Fetcher) error { f.URL.Path = f.URL.Path + "cmd1"; return nil })
	cmd2 := CommandFunc(func(f *Fetcher) error { f.URL.Path = f.URL.Path + "cmd2"; return nil })
	cmd3 := CommandFunc(func(f *Fetcher) error { f.URL.Path = f.URL.Path + "cmd3"; return nil })
	cmd4 := CommandFunc(func(f *Fetcher) error { f.URL.Path = f.URL.Path + "cmd4"; return nil })
	cmd5 := CommandFunc(func(f *Fetcher) error { f.URL.Path = f.URL.Path + "cmd5"; return nil })
	p := NewPreset().CloneWith(cmd1, cmd2, cmd3, cmd4, cmd5)
	f := New()
	err := Exec(f, p)
	if err != nil {
		panic(err)
	}
	if f.URL.Path != "cmd1cmd2cmd3cmd4cmd5" {
		t.Fatal(f.URL.Path)
	}

	f = New()
	err = Exec(f, p.Commands()...)
	if err != nil {
		panic(err)
	}
	if f.URL.Path != "cmd1cmd2cmd3cmd4cmd5" {
		t.Fatal(f.URL.Path)
	}

	p = NewPreset().Concat(cmd1, cmd2).Concat(cmd3, cmd4, cmd5)
	f = New()
	err = Exec(f, p.Commands()...)
	if err != nil {
		panic(err)
	}
	if f.URL.Path != "cmd1cmd2cmd3cmd4cmd5" {
		t.Fatal(f.URL.Path)
	}

	p = NewPreset().Append(NewPreset().Concat(cmd1, cmd2)).Append(NewPreset().Concat(cmd3, cmd4, cmd5))
	f = New()
	err = Exec(f, p.Commands()...)
	if err != nil {
		panic(err)
	}
	if f.URL.Path != "cmd1cmd2cmd3cmd4cmd5" {
		t.Fatal(f.URL.Path)
	}
	p = NewPreset().Append(NewPreset().Concat(cmd1, cmd2), NewPreset().Concat(cmd3, cmd4, cmd5))
	f = New()
	err = Exec(f, p.Commands()...)
	if err != nil {
		panic(err)
	}
	if f.URL.Path != "cmd1cmd2cmd3cmd4cmd5" {
		t.Fatal(f.URL.Path)
	}

	p = NewPreset().Append(NewPreset().Concat(cmd1, cmd2)).Concat(cmd3, cmd4, cmd5)
	f = New()
	err = Exec(f, p.Commands()...)
	if err != nil {
		panic(err)
	}
	if f.URL.Path != "cmd1cmd2cmd3cmd4cmd5" {
		t.Fatal(f.URL.Path)
	}
}
