package iapgo

import (
	"net/http"
)

type Transport struct{}

func (t *Transport) RoundTrip(r *http.Request) (*http.Response, error) {
	return nil, nil
}
