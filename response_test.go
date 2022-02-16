package fetcher

import (
	"bytes"
	"strings"
	"testing"
)

func TestAsError(t *testing.T) {
	s := newEchoServer()
	defer s.Close()
	var sc = &Server{
		ServerInfo: ServerInfo{
			URL: s.URL,
		},
	}
	preset := MustPreset(sc)
	resp, err := preset.FetchWithBody(bytes.NewBufferString("errbody"))
	if err != nil {
		t.Fatal(err)
	}
	errmsg := resp.Error()
	if !strings.Contains(errmsg, "errbody") || !strings.Contains(errmsg, "GET") || !strings.Contains(errmsg, s.URL) || !strings.Contains(errmsg, "200") {
		t.Fatal(errmsg)
	}
	if GetAPIErrCode(resp) != "" {
		t.Fatal(resp)
	}
	if GetAPIErrContent(resp) != "" {
		t.Fatal(resp)
	}
	if CompareAPIErrCode(resp, 200) != false {
		t.Fatal(resp)
	}
	if !IsResponseErr(resp) {
		t.Fatal(resp)
	}
	if !CompareResponseErrStatusCode(resp, 200) {
		t.Fatal(resp)
	}
	if CompareResponseErrStatusCode(resp, 404) {
		t.Fatal(resp)
	}
	errcode := resp.NewAPICodeErr("999")
	errcodemsg := errcode.Error()
	if !strings.Contains(errcodemsg, "errbody") || !strings.Contains(errcodemsg, "GET") || !strings.Contains(errcodemsg, s.URL) || !strings.Contains(errcodemsg, "999") {
		t.Fatal(errmsg)
	}
	if GetAPIErrCode(errcode) != "999" {
		t.Fatal(resp)
	}
	errcontent := GetAPIErrContent(errcode)
	if errcontent != "errbody" {
		t.Fatal(errcontent)
	}
	if CompareAPIErrCode(errcode, 999) != true {
		t.Fatal(resp)
	}
	tooLongBody := bytes.Repeat([]byte{'!'}, ErrMsgLengthLimit+1)
	resp, err = preset.FetchWithBody(bytes.NewBuffer(tooLongBody))
	if err != nil {
		t.Fatal(err)
	}
	errmsg = resp.Error()
	if !strings.Contains(errmsg, "GET") || !strings.Contains(errmsg, s.URL) || !strings.Contains(errmsg, "200") {
		t.Fatal(errmsg)
	}
	if len(errmsg) != ErrMsgLengthLimit {
		t.Fatal(errmsg)
	}
	errcode = resp.NewAPICodeErr("999")
	errcodemsg = errcode.Error()
	if !strings.Contains(errcodemsg, "GET") || !strings.Contains(errcodemsg, s.URL) || !strings.Contains(errcodemsg, "999") {
		t.Fatal(errmsg)
	}
	if len(errcodemsg) != ErrMsgLengthLimit {
		t.Fatal(errcodemsg)
	}
	errcode = resp.NewAPICodeErr("")
	errcodemsg = errcode.(*APICodeErr).ErrorPrivateRef()
	if errcodemsg != "" {
		t.Fatal(errcodemsg)
	}
	errcode = resp.NewAPICodeErr("999")
	errcodemsg = errcode.(*APICodeErr).ErrorPrivateRef()
	if !strings.Contains(errcodemsg, "999") {
		t.Fatal(errcodemsg)
	}
	if CompareResponseErrStatusCode(errcode, 200) {
		t.Fatal(resp)
	}
}
