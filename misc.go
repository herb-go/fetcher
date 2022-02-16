package fetcher

import (
	"net/http"
)

//CloneHeader clone http header
func CloneHeader(h http.Header) http.Header {
	return h.Clone()
}

//MergeHeader merge src header to dst
func MergeHeader(dst http.Header, src http.Header) {
	for name := range src {
		for k := range src[name] {
			dst.Set(name, src[name][k])
		}
	}
}

//CloneRequestBuilders clone request builders
func CloneRequestBuilders(b []func(*http.Request) error) []func(*http.Request) error {
	builders := make([]func(*http.Request) error, len(b))
	copy(builders, b)
	return builders
}
