// Package iapgo helps authenticating access to endpoints behind Google Cloud
// Identity-Aware Proxy (IAP).  It provides a Transport which implements
// http.RoundTripper.
package iapgo

import (
	"context"
	"errors"
	"net/http"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

var credentialsFinder = google.FindDefaultCredentials

var errUninitialized = errors.New("iapgo: unitialized Transport")

// Transport implements http.RoundTripper that can be used to access endpoints
// behind Google Cloud Identity-Aware Proxy.
type Transport struct {
	oauthTransport *oauth2.Transport
}

// NewTransport returns an initialized Transport.  It requires OAuth Client ID
// of the IAP resource target of the Transport.  It finds the service account
// key using Application Default Credentials (ADC) strategy described in
// https://cloud.google.com/docs/authentication/production.
func NewTransport(iapClientID string) (*Transport, error) {
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
		"target_audience": iapClientID,
	}

	conf.UseIDToken = true

	transport.oauthTransport = &oauth2.Transport{
		Source: conf.TokenSource(context.Background()),
		Base:   http.DefaultTransport,
	}

	return transport, nil
}

// RoundTrip authenticates an HTTP request using an ID token.  This ID token is
// retrieved using two-legged authentication with a Google endpoint defined in
// the service account key.
func (t *Transport) RoundTrip(r *http.Request) (*http.Response, error) {
	if t.oauthTransport == nil {
		return nil, errUninitialized
	}

	return t.oauthTransport.RoundTrip(r)
}
