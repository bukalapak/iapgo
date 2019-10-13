# iapgo

[![Build Status](https://travis-ci.com/saifulwebid/iapgo.svg?branch=master)](https://travis-ci.com/saifulwebid/iapgo)
[![codecov](https://codecov.io/gh/saifulwebid/iapgo/branch/master/graph/badge.svg)](https://codecov.io/gh/saifulwebid/iapgo)
[![GoDoc](https://godoc.org/github.com/saifulwebid/iapgo?status.svg)](https://godoc.org/github.com/saifulwebid/iapgo)

iapgo is a Go library to help authenticating access to endpoints behind Google [Cloud Identity-Aware Proxy](https://cloud.google.com/iap/).

This library is heavily using [`golang.org/x/oauth2/google`](https://godoc.org/golang.org/x/oauth2/google) to handle credentials parsing and [authentication](https://cloud.google.com/iap/docs/authentication-howto).

## Usage

```go
import (
    "log"
    "net/http"

    "github.com/saifulwebid/iapgo"
)

func main() {
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
```
