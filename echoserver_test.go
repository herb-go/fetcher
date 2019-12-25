package fetcher

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strconv"
)

func EchoAction(w http.ResponseWriter, r *http.Request) {
	for field := range r.Header {
		for _, v := range r.Header[field] {
			w.Header().Add(field, v)
		}
	}
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}
	statuscode := r.URL.Query().Get("statuscode")
	w.Header().Set("rawpath", r.URL.RawPath)
	code, err := strconv.Atoi(statuscode)
	if err != nil && code != 0 {
		w.WriteHeader(code)
	}
	w.Write(data)
}

func newEchoServer() *httptest.Server {
	s := httptest.NewServer(http.HandlerFunc(EchoAction))
	return s
}
