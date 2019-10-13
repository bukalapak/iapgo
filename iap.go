package iapgo

import (
	"context"
	"net/http"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type credentialsFinderFn func(ctx context.Context, scopes ...string) (*google.Credentials, error)

var (
	credentialsFinder credentialsFinderFn = google.FindDefaultCredentials
)

type Transport struct {
	clientID string

	oauthTransport *oauth2.Transport
}

func (t *Transport) RoundTrip(r *http.Request) (*http.Response, error) {
	if t.oauthTransport == nil {
		creds, err := credentialsFinder(context.Background())
		if err != nil {
			return nil, err
		}

		conf, err := google.JWTConfigFromJSON(creds.JSON)
		if err != nil {
			return nil, err
		}

		conf.PrivateClaims = map[string]interface{}{
			"target_audience": t.clientID,
		}

		conf.UseIDToken = true

		t.oauthTransport = &oauth2.Transport{
			Source: conf.TokenSource(context.Background()),
			Base:   http.DefaultTransport,
		}
	}

	return t.oauthTransport.RoundTrip(r)
}
