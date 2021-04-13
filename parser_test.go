package fetcher

import (
	"bytes"
	"testing"
)

func TestParser(t *testing.T) {
	var err error
	var resultstr string
	var result []byte
	s := newEchoServer()
	defer s.Close()
	var sc = &Server{
		ServerInfo: ServerInfo{
			URL: s.URL,
		},
	}
	preset := MustPreset(sc)
	_, err = FetchAndParse(preset, Should200(nil))
	if err != nil {
		t.Fatal(err)
	}
	_, err = FetchAndParse(preset.Concat(SetQuery("statuscode", "206")), Should200(nil))
	if err == nil {
		t.Fatal(err)
	}
	_, err = FetchAndParse(preset.Concat(SetQuery("statuscode", "206")), ShouldSuccess(nil))
	if err != nil {
		t.Fatal(err)
	}
	_, err = FetchAndParse(preset.Concat(SetQuery("statuscode", "404")), ShouldNoError(nil))
	if err != nil {
		t.Fatal(err)
	}
	_, err = FetchAndParse(preset.Concat(SetQuery("statuscode", "303")), ShouldSuccess(nil))
	if err == nil {
		t.Fatal(err)
	}
	_, err = FetchAndParse(preset.Concat(SetQuery("statuscode", "504")), ShouldNoError(nil))
	if err == nil {
		t.Fatal(err)
	}

	result = []byte{'a', 'b', 'c'}
	_, err = FetchAndParse(preset.Concat(SetQuery("statuscode", "200")), Should200(AsBytes(&result)))
	if err != nil {
		t.Fatal(err)
	}
	if len(result) != 0 {
		t.Fatal(result)
	}
	result = []byte{'a', 'b', 'c'}
	_, err = DoAndParse(nil, preset.Concat(SetQuery("statuscode", "306")), Should200(AsBytes(&result)))
	if err == nil {
		t.Fatal(err)
	}
	if len(result) != 3 {
		t.Fatal(result)
	}

	resultstr = ""
	_, err = FetchWithBodyAndParse(preset.Concat(SetQuery("statuscode", "206")), bytes.NewBufferString("statuscode206"), ShouldSuccess(AsString(&resultstr)))
	if err != nil {
		t.Fatal(err)
	}
	if resultstr != "statuscode206" {
		t.Fatal(resultstr)
	}
	resultstr = ""
	_, err = DoWithBodyAndParse(nil, preset.Concat(SetQuery("statuscode", "306")), bytes.NewBufferString("statuscode306"), ShouldSuccess(AsString(&resultstr)))
	if err == nil {
		t.Fatal(err)
	}
	if resultstr != "" {
		t.Fatal(resultstr)
	}
	resultstr = ""
	_, err = FetchWithJSONBodyAndParse(preset.Concat(SetQuery("statuscode", "406")), "statuscode406", ShouldNoError(AsJSON(&resultstr)))
	if err != nil {
		t.Fatal(err)
	}
	if resultstr != "statuscode406" {
		t.Fatal(resultstr)
	}
	resultstr = ""
	_, err = DoWithJSONBodyAndParse(nil, preset.Concat(SetQuery("statuscode", "504")), "statuscode504", ShouldNoError(AsJSON(&resultstr)))
	if err == nil {
		t.Fatal(err)
	}
	if resultstr != "" {
		t.Fatal(resultstr)
	}

}
