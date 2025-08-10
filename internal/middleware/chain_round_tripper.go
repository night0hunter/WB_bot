package middleware

import "net/http"

// internalRoundTripper is a holder function to make the process of
// creating middleware a bit easier without requiring the consumer to
// implement the RoundTripper interface.
type internalRoundTripper func(*http.Request) (*http.Response, error)

func (rt internalRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	return rt(req)
}

// Middleware is our middleware creation functionality.
type Middleware func(http.RoundTripper) http.RoundTripper

// Chain is a handy function to wrap a base RoundTripper (optional)
// with the middlewares.
func Chain(rt http.RoundTripper, middlewares ...Middleware) http.RoundTripper {
	if rt == nil {
		rt = http.DefaultTransport
	}

	for _, m := range middlewares {
		rt = m(rt)
	}

	return rt
}
