package iapgo_test

import (
	"log"
	"net/http"

	"github.com/saifulwebid/iapgo"
)

func Example() {
	// Initialize Transport to be used. Define iapClientID with the OAuth Client
	// ID of the IAP that protects the endpoint.
	iapClientID := "12345678901-abcdefghijklmnopqrstuvwxyz123456.apps.googleusercontent.com"

	// Upon Transport creation, the service account key will be searched using
	// Application Default Credentials (ADC) strategy described in
	// https://cloud.google.com/docs/authentication/production.
	transport, err := iapgo.NewTransport(iapClientID)
	if err != nil {
		log.Fatal(err)
	}

	// Pair Transport with an http.Client.
	client := &http.Client{
		Transport: transport,
	}

	// Access endpoints behind IAP.
	client.Get("...")
}
