package iapgo

import (
	"bytes"
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"golang.org/x/oauth2/google"
	"golang.org/x/oauth2/jws"
)

// Ensures Transport implements http.RoundTripper.
var _ http.RoundTripper = new(Transport)

func TestTransport_RoundTrip(t *testing.T) {
	clientID := "ABCD"

	authSvr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertAuthorizationRequest(t, r, clientID)
		w.Write(mockTokenResponse(t))
	}))

	origCredentialsFinder := credentialsFinder
	credentialsFinder = mockServiceAccountKey(t, authSvr.URL)
	defer func() {
		credentialsFinder = origCredentialsFinder
	}()

	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertAuthorizationHeader(t, r.Header.Get("Authorization"))
	}))

	client := &http.Client{
		Transport: &Transport{
			clientID: clientID,
		},
	}

	_, err := client.Get(svr.URL)
	if err != nil {
		t.Fatal("client.Get:", err)
	}
}

func TestTransport_RoundTrip_CredentialsError(t *testing.T) {
	origCredentialsFinder := credentialsFinder
	credentialsFinder = brokenCredentialsFn
	defer func() {
		credentialsFinder = origCredentialsFinder
	}()

	client := &http.Client{
		Transport: new(Transport),
	}

	_, err := client.Get("http://localhost")
	switch {
	case err == nil:
		t.Fatal("no error, want some")
	case !strings.HasSuffix(err.Error(), "brokenCredentialsFn"):
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestTransport_RoundTrip_JWTError(t *testing.T) {
	origCredentialsFinder := credentialsFinder
	credentialsFinder = mockUserCredentials(t, "http://localhost")
	defer func() {
		credentialsFinder = origCredentialsFinder
	}()

	client := &http.Client{
		Transport: new(Transport),
	}

	_, err := client.Get("http://localhost")
	switch {
	case err == nil:
		t.Fatal("no error, want some")
	case !strings.HasSuffix(err.Error(), `'type' field is "authorized_user" (expected "service_account")`):
		t.Fatalf("unexpected error: %v", err)
	}
}

func mockServiceAccountKey(t *testing.T, authURL string) credentialsFinderFn {
	t.Helper()

	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatalf("rsa.GenerateKey: %v", err)
	}

	enc := pem.EncodeToMemory(&pem.Block{
		Type:  "PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
	})

	jwtKey := struct {
		Type       string `json:"type"`
		TokenURI   string `json:"token_uri"`
		PrivateKey string `json:"private_key"`
	}{
		Type:       "service_account",
		TokenURI:   authURL,
		PrivateKey: string(enc),
	}

	jsonKey, err := json.Marshal(jwtKey)
	if err != nil {
		t.Fatalf("json.Marshal: %v", err)
	}

	return func(ctx context.Context, scopes ...string) (*google.Credentials, error) {
		return google.CredentialsFromJSON(context.Background(), jsonKey, scopes...)
	}
}

func mockUserCredentials(t *testing.T, authURL string) credentialsFinderFn {
	t.Helper()

	jwtKey := struct {
		Type string `json:"type"`
	}{
		Type: "authorized_user",
	}

	jsonKey, err := json.Marshal(jwtKey)
	if err != nil {
		t.Fatalf("json.Marshal: %v", err)
	}

	return func(ctx context.Context, scopes ...string) (*google.Credentials, error) {
		return google.CredentialsFromJSON(context.Background(), jsonKey, scopes...)
	}
}

func mockTokenResponse(t *testing.T) []byte {
	t.Helper()

	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatalf("rsa.GenerateKey: %v", err)
	}

	cs := &jws.ClaimSet{}
	header := &jws.Header{}

	idToken, err := jws.Encode(header, cs, privateKey)
	if err != nil {
		t.Fatalf("jws.Encode: %v", err)
	}

	retVal := struct {
		IDToken string `json:"id_token"`
	}{
		IDToken: idToken,
	}

	resp, err := json.Marshal(retVal)
	if err != nil {
		t.Fatalf("json.Marshal: %v", err)
	}

	return resp
}

func assertAuthorizationRequest(t *testing.T, r *http.Request, expectedClientID string) {
	t.Helper()

	if r.Body == nil {
		t.Fatal("r.Body is nil")
	}

	defer r.Body.Close()

	r.ParseForm()

	assertion := r.Form.Get("assertion")
	if assertion == "" {
		t.Fatal("empty assertion")
	}

	assertJWS(t, assertion)

	grantType := r.Form.Get("grant_type")
	expectedGrantType := "urn:ietf:params:oauth:grant-type:jwt-bearer"
	if grantType != expectedGrantType {
		t.Fatalf("grant_type = %s, want %s", grantType, expectedGrantType)
	}
}

func assertJWS(t *testing.T, payload string) {
	s := strings.Split(payload, ".")
	if len(s) < 2 {
		t.Fatal("jws: invalid token received")
	}

	decoded, err := base64.RawURLEncoding.DecodeString(s[1])
	if err != nil {
		t.Fatal(err)
	}

	var data map[string]interface{}
	err = json.NewDecoder(bytes.NewBuffer(decoded)).Decode(&data)
	if err != nil {
		t.Fatal("json.Decode:", err)
	}

	targetAudience := data["target_audience"]
	if targetAudience != "ABCD" {
		t.Fatalf("target_audience = %v, want %s", targetAudience, "ABCD")
	}
}

func assertAuthorizationHeader(t *testing.T, header string) {
	t.Helper()

	if header == "" {
		t.Fatal("empty Authorization header, want Bearer token")
	}

	if !strings.HasPrefix(header, "Bearer ") {
		t.Fatalf("header = %s, want Bearer token", header)
	}

	if len(header) <= len("Bearer ") {
		t.Fatal("empty Bearer token, want some")
	}
}

func brokenCredentialsFn(ctx context.Context, scopes ...string) (*google.Credentials, error) {
	return nil, fmt.Errorf("brokenCredentialsFn")
}
