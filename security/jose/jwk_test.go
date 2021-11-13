package jose

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

const (
	JwkSetInvalidString = "Some string."
	JwkSetInvalidJson   = `{
		"field": "value"
	}`
	JwkSetValidJson = `{
		"keys": [
			{
				"alg": "RS256",
				"kty": "RSA",
				"use": "sig",
				"n": "...",
				"e": "AQAB",
				"kid": "...",
				"x5t": "...",
				"x5c": [
					"..."
				]
			},
			{
				"alg": "RS256",
				"kty": "RSA",
				"use": "sig",
				"n": "...",
				"e": "AQAB",
				"kid": "...",
				"x5t": "...",
				"x5c": [
					"..."
				]
			}
		]
	}`
)

func TestReadJwkSet(t *testing.T) {
	testCases := []struct {
		name  string
		jwks  string
		valid bool
	}{
		{"invalid string", JwkSetInvalidString, false},
		{"invalid json", JwkSetInvalidJson, false},
		{"valid json", JwkSetValidJson, true},
	}

	var want string
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, want)
	}))

	defer ts.Close()

	for _, tc := range testCases {
		want = tc.jwks
		t.Run(tc.name, func(t *testing.T) {
			res, err := readJwkSet(ts.URL)
			if tc.valid && err != nil {
				t.Fatalf(`readJwkSet() = %+v, %v, want match for %+v, nil`, res, err, want)
			}
		})
	}
}
