package openproject

import (
	"context"
	"io"
	"net/http"
)

func newRequestWithContext(ctx context.Context, method, url string, body io.Reader) (*http.Request, error) {
	return http.NewRequestWithContext(ctx, method, url, body)
}
