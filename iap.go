package iapgo

import (
	"context"
	"errors"
	"net/http"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type credentialsFinderFn func(ctx context.Context, scopes ...string) (*google.Credentials, error)

var credentialsFinder credentialsFinderFn = google.FindDefaultCredentials

var errUninitialized = errors.New("iapgo: unitialized Transport")

type Transport struct {
	oauthTransport *oauth2.Transport
}

func newTransport(clientID string) (*Transport, error) {
	transport := &Transport{}

	creds, err := credentialsFinder(context.Background())
	if err != nil {
		return nil, err
	}

	conf, err := google.JWTConfigFromJSON(creds.JSON)
	if err != nil {
		return nil, err
	}

	conf.PrivateClaims = map[string]interface{}{
		"target_audience": clientID,
	}

	conf.UseIDToken = true

	transport.oauthTransport = &oauth2.Transport{
		Source: conf.TokenSource(context.Background()),
		Base:   http.DefaultTransport,
	}

	return transport, nil
}

func (t *Transport) RoundTrip(r *http.Request) (*http.Response, error) {
	if t.oauthTransport == nil {
		return nil, errUninitialized
	}

	return t.oauthTransport.RoundTrip(r)
}
