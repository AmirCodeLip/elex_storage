package http_handlers

import (
	"io"
	"net/http"
	"strings"
)

var hopByHopHeaders = map[string]bool{
	"Connection":          true,
	"Keep-Alive":          true,
	"Proxy-Authenticate":  true,
	"Proxy-Authorization": true,
	"Te":                  true,
	"Trailers":            true,
	"Transfer-Encoding":   true,
	"Upgrade":             true,
}

func HttpProxyFunc(handler *HttpHandler, fn func(r *http.Request) (*http.Response, error)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 1. Call method from api
		resp, err := fn(r)
		if err != nil {
			handler.logger.Error(err.Error())
			handler.httpErrorUtils.InternalServerError(w)
			return
		}
		defer resp.Body.Close()

		// 2. Copy response headers
		for k, v := range resp.Header {
			if hopByHopHeaders[k] {
				continue
			}
			// let your proxy own the CORS headers
			if strings.HasPrefix(k, "Access-Control-") {
				continue
			}
			for _, vv := range v {
				w.Header().Add(k, vv)
			}
		}
		w.WriteHeader(resp.StatusCode)
		io.Copy(w, resp.Body)
	}
}
