package middleware

import (
	"net/http"
	"strings"
)

// SetAccCookieMiddleware - middleware for setting account cookie to the request
func SetAccCookieMiddleware(cookie string) Middleware {
	return func(rt http.RoundTripper) http.RoundTripper {
		return internalRoundTripper(func(req *http.Request) (*http.Response, error) {
			cookies := strings.Split(cookie, "; ")
			for _, cookie := range cookies {
				parts := strings.SplitN(cookie, "=", 2)
				if len(parts) == 2 {
					req.AddCookie(&http.Cookie{
						Name:  parts[0],
						Value: parts[1],
					})
				}
			}

			return rt.RoundTrip(req)
		})
	}
}
