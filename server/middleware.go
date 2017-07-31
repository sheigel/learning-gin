package server

import (
	"net/http"
	"net/http/httputil"
	"fmt"
)

func loggingMiddleware(handlerFunc http.HandlerFunc) func(w http.ResponseWriter, r *http.Request) {
	return func(res http.ResponseWriter, req *http.Request) {
		dump, _ := httputil.DumpRequest(req, true)
		fmt.Printf("received %s at %s with body %s\n", req.Method, req.RequestURI, string(dump))
		handlerFunc(res, req)
	}
}
