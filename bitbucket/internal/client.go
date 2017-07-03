package internal

import (
	"context"
	"net/http"
)

var HTTPClient ContextKeyHTTPClient

type ContextKeyHTTPClient struct{}

func ContextClient(ctx context.Context) *http.Client {
	if ctx != nil {
		if hc, ok := ctx.Value(HTTPClient).(*http.Client); ok {
			return hc
		}
	}
	return http.DefaultClient
}
