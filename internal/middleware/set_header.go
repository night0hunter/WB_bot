package middleware

import (
	"net/http"
)

// SetHeader - middleware for setting header to request
func SetHeader(key, value string) Middleware {
	return func(rt http.RoundTripper) http.RoundTripper {
		return internalRoundTripper(func(req *http.Request) (*http.Response, error) {
			header := req.Header
			if header == nil {
				header = make(http.Header)
			}

			header.Set(key, value)

			return rt.RoundTrip(req)
		})
	}
}
