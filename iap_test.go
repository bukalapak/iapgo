package iapgo_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/saifulwebid/iapgo"
)

var _ http.RoundTripper = new(iapgo.Transport)

func TestTransport_RoundTrip(t *testing.T) {
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.RequestURI != "/" {
			t.Fatal("request URI doesn't match")
		}
	}))

	client := http.Client{
		Transport: new(iapgo.Transport),
	}

	req, err := http.NewRequest("GET", svr.URL, nil)
	if err != nil {
		t.Fatal("error creating a request")
	}

	client.Do(req)
}
