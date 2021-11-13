package jose

import (
	"crypto"
	"testing"
)

func TestFetchHashEmpty(t *testing.T) {
	alg := ""
	_, err := fetchHash(alg)
	if err == nil {
		t.Fatalf(`fetchHash("") = _, %v, want match for _, nil`, err)
	}
}

func TestFetchHashRS256(t *testing.T) {
	alg := "RS256"
	want := crypto.SHA256
	hash, err := fetchHash(alg)
	if hash != want || err != nil {
		t.Fatalf(`fetchHash("RS256") = %q, %v, want match for %#q, nil`, hash, err, want)
	}
}

func TestFetchHashRS384(t *testing.T) {
	alg := "RS384"
	want := crypto.SHA384
	hash, err := fetchHash(alg)
	if hash != want || err != nil {
		t.Fatalf(`fetchHash("RS256") = %q, %v, want match for %#q, nil`, hash, err, want)
	}
}

func TestFetchHashRS512(t *testing.T) {
	alg := "RS512"
	want := crypto.SHA512
	hash, err := fetchHash(alg)
	if hash != want || err != nil {
		t.Fatalf(`fetchHash("RS256") = %q, %v, want match for %#q, nil`, hash, err, want)
	}
}

func TestHasAudience(t *testing.T) {
	data := []string{"audience1", "audience2", "audience3"}
	aud := "audience2"
	want := true
	has := hasAudience(data, aud)
	if has != want {
		t.Fatalf(`hasAudience("[...], "audience2") = %t, want match for %t`, has, want)
	}
}

func TestHasAudienceNot(t *testing.T) {
	data := []string{"audience1", "audience2", "audience3"}
	aud := "audience4"
	want := false
	has := hasAudience(data, aud)
	if has != want {
		t.Fatalf(`hasAudience("[...], "audience4") = %t, want match for %t`, has, want)
	}
}

func TestParseAudience(t *testing.T) {
	data := "audience"
	want := []string{"audience"}
	audience, err := parseAudience(data)
	if (audience != nil && audience[0] != want[0]) || err != nil {
		t.Fatalf(`parseAudience("") = %q, %v, want match for %#q, nil`, audience, err, want)
	}
}

func TestParseAudienceArray(t *testing.T) {
	data := []string{"audience1", "audience2"}
	want := []string{"audience1", "audience2"}
	audience, err := parseAudience(data)
	if (audience != nil && audience[0] != want[0]) || err != nil {
		t.Fatalf(`parseAudience("") = %q, %v, want match for %#q, nil`, audience, err, want)
	}
}
