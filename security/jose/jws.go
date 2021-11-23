package jose

import (
	"crypto"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"strings"
	"time"
	"unicode/utf8"
)

// JSON Object Signing and Encryption (jose)

type JoseHeader struct {
	Alg string `json:"alg"`
	Kid string `json:"kid"`
	Typ string `json:"typ"`
}

type Jwk struct {
	Alg string   `json:"alg"`
	Kty string   `json:"kty"`
	Use string   `json:"use"`
	N   string   `json:"n"`
	E   string   `json:"e"`
	Kid string   `json:"kid"`
	X5T string   `json:"x5t"`
	X5C []string `json:"x5c"`
}

type JwkSet struct {
	Keys []Jwk `json:"keys"`
}

type Jwt struct {
	Iss string      `json:"iss"`
	Sub string      `json:"sub"`
	Aud interface{} `json:"aud"`
	Exp int64       `json:"exp"`
	Iat int64       `json:"iat"`
}

// ExpirationTime returns the local expiration time on or after
// which the JWT MUST NOT be accepted
func (jwt *Jwt) ExpirationTime() time.Time {
	return time.Unix(jwt.Exp, 0)
}

// IssuedAt returns the local time at which the JWT was issued
func (jwt *Jwt) IssuedAt() time.Time {
	return time.Unix(jwt.Iat, 0)
}

// VerifyCompact returns the verified state of a JWT using the
// JWS Compact Serialization format.
func VerifyCompact(token string, issuer string, audience string) error {
	if len(token) == 0 {
		return errors.New("missing token")
	}

	segments := strings.Split(token, ".")
	if len(segments) != 3 {
		return errors.New("incompatible token detected (not JWS compact)")
	}

	_, err := url.ParseRequestURI(issuer)
	if err != nil {
		return errors.New("improperyl formatted issuer")
	}

	jwksUrl := fmt.Sprintf("%s/.well-known/jwks.json", strings.TrimSuffix(issuer, "/"))

	decoder := base64.RawURLEncoding.DecodeString

	// JWS header

	header := segments[0]
	decodedHeader, err := decoder(header)
	if err != nil {
		return errors.New("unable to decode token header")
	}

	if !utf8.Valid(decodedHeader) {
		return errors.New("not a valid UTF-8 encoded sequence")
	}

	var joseHeader JoseHeader
	json.Unmarshal(decodedHeader, &joseHeader)
	if err != nil {
		return errors.New("unable to parse token header")
	}

	// JWS payload

	payload := segments[1]
	decodedPayload, err := decoder(payload)
	if err != nil {
		return errors.New("unable to decode token payload")
	}

	var jwt Jwt
	json.Unmarshal(decodedPayload, &jwt)
	if err != nil {
		return errors.New("unable to parse token payload")
	}

	// JWS signature

	signature := segments[2]
	decodedSignature, err := decoder(signature)
	if err != nil {
		return errors.New("unable to decode token signature")
	}

	kid := joseHeader.Kid
	encodedDer, err := fetchEncodedDer(jwksUrl, kid)
	if err != nil {
		return errors.New("unable to read JWKS der cert")
	}

	key, err := publicKeyFromEncodedDer(encodedDer)
	if err != nil {
		return errors.New("unable to read public key from der cert")
	}

	alg := joseHeader.Alg
	input := fmt.Sprintf("%s.%s", header, payload)
	err = verifySignature(key, alg, input, decodedSignature)
	if err != nil {
		return err
	}

	// JWT claims

	now := time.Now()
	if !now.Before(jwt.ExpirationTime()) {
		return errors.New("token expired")
	}

	if jwt.Iss != issuer {
		return errors.New("invalid issuer")
	}

	aud, err := parseAudience(jwt.Aud)
	if err != nil {
		return err
	}

	if !hasAudience(aud, audience) {
		return errors.New("invalid audience")
	}

	return nil
}

func fetchHash(alg string) (crypto.Hash, error) {
	switch alg {
	case "RS256":
		return crypto.SHA256, nil
	case "RS384":
		return crypto.SHA384, nil
	case "RS512":
		return crypto.SHA512, nil
	default:
		return crypto.SHA256, errors.New("unknown hash algorithm")
	}
}

func hasAudience(data []string, aud string) bool {
	has := false
	for _, val := range data {
		has = val == aud
		if has {
			break
		}
	}

	return has
}

func parseAudience(data interface{}) ([]string, error) {
	switch a := data.(type) {
	case string:
		return []string{a}, nil
	case []string:
		return a, nil
	case []interface{}:
		aud := make([]string, len(a))
		for i, v := range a {
			aud[i] = fmt.Sprint(v)
		}
		return aud, nil
	default:
		return nil, errors.New("unable to parse audience")
	}
}

func publicKeyFromEncodedDer(encodedDer string) (*rsa.PublicKey, error) {
	decoder := base64.StdEncoding.DecodeString
	decodedDer, err := decoder(encodedDer)
	if err != nil {
		return nil, err
	}

	cert, err := x509.ParseCertificate(decodedDer)
	if err != nil {
		return nil, err
	}

	key, ok := cert.PublicKey.(*rsa.PublicKey)
	if !ok {
		return nil, errors.New("not rsa public key")
	}

	return key, nil
}

func verifySignature(key *rsa.PublicKey, alg string, signingInput string, signature []byte) error {
	hash, err := fetchHash(alg)
	if err != nil {
		return err
	}

	hasher := hash.New()
	hasher.Write([]byte(signingInput))

	err = rsa.VerifyPKCS1v15(key, hash, hasher.Sum(nil), signature)

	return err
}

// func publicKeyFromExponentAndModulus(encodedE string, encodedN string) (*rsa.PublicKey, error) {
// 	decoder := base64.RawURLEncoding.DecodeString

// 	decodedE, err := decoder(encodedE)
// 	if err != nil {
// 		return nil, err
// 	}

// 	var buffer bytes.Buffer
// 	buffer.WriteByte(0)
// 	buffer.Write(decodedE)
// 	e := binary.BigEndian.Uint32(buffer.Bytes())

// 	decodedN, err := decoder(encodedN)
// 	if err != nil {
// 		return nil, err
// 	}

// 	n := new(big.Int)
// 	n.SetBytes(decodedN)

// 	var key = &rsa.PublicKey{
// 		E: int(e),
// 		N: n,
// 	}

// 	return key, nil
// }
